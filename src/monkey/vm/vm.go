package vm

import (
	"encoding/binary"
	"fmt"
	"writing-in-interpreter-in-go/src/monkey/code"
	"writing-in-interpreter-in-go/src/monkey/compiler"
	"writing-in-interpreter-in-go/src/monkey/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

type VirtualMachine struct {
	sp          int
	constants   []object.Object
	globals     []object.Object
	stack       []object.Object
	frames      []*Frame
	framesIndex int
}

func New(byteCode *compiler.ByteCode) *VirtualMachine {
	mainFn := &object.CompiledFunction{Instructions: byteCode.Instructions}
	mainClosure := &object.Closure{Function: mainFn}
	mainFrame := NewFrame(mainClosure, 0)
	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	return &VirtualMachine{
		sp:          0,
		stack:       make([]object.Object, StackSize),
		constants:   byteCode.Constants,
		globals:     make([]object.Object, GlobalsSize),
		frames:      frames,
		framesIndex: 1,
	}
}

func NewWithGlobalsStore(bytecode *compiler.ByteCode, s []object.Object) *VirtualMachine {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (virtualMachine *VirtualMachine) currentFrame() *Frame {
	return virtualMachine.frames[virtualMachine.framesIndex-1]
}

func (virtualMachine *VirtualMachine) pushFrame(frame *Frame) {
	virtualMachine.frames[virtualMachine.framesIndex] = frame
	virtualMachine.framesIndex++
}

func (virtualMachine *VirtualMachine) popFrame() *Frame {
	virtualMachine.framesIndex--
	return virtualMachine.frames[virtualMachine.framesIndex]
}

func (virtualMachine *VirtualMachine) LastPopped() object.Object {
	return virtualMachine.stack[virtualMachine.sp]
}

func (virtualMachine *VirtualMachine) Run() error {
	for virtualMachine.currentFrame().ip < len(virtualMachine.currentFrame().Instructions())-1 {
		instructions := virtualMachine.currentFrame().Instructions()
		virtualMachine.currentFrame().ip++
		ip := virtualMachine.currentFrame().ip
		opcode := code.Opcode(instructions[ip])
		switch opcode {
		case code.OpConstant:
			constantIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			frame := virtualMachine.currentFrame()
			frame.ip += 2
			err := virtualMachine.push(virtualMachine.constants[constantIndex])
			if err != nil {
				return err
			}
		case code.OpGetGlobal:
			globalIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			virtualMachine.currentFrame().ip += 2
			global := virtualMachine.globals[globalIndex]
			err := virtualMachine.push(global)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			virtualMachine.currentFrame().ip += 2
			virtualMachine.globals[globalIndex] = virtualMachine.pop()
		case code.OpGetLocal:
			localIndex := instructions[ip+1]
			virtualMachine.currentFrame().ip += 1
			frame := virtualMachine.currentFrame()
			err := virtualMachine.push(virtualMachine.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := instructions[ip+1]
			virtualMachine.currentFrame().ip += 1
			frame := virtualMachine.currentFrame()
			virtualMachine.stack[frame.basePointer+int(localIndex)] = virtualMachine.pop()
		case code.OpGetFree:
			index := instructions[ip+1]
			virtualMachine.currentFrame().ip += 1
			closure := virtualMachine.currentFrame().closure
			err := virtualMachine.push(closure.FreeVariables[index])
			if err != nil {
				return err
			}
		case code.OpGetBuiltin:
			builtinIndex := instructions[ip+1]
			virtualMachine.currentFrame().ip += 1
			builtinDefinition := object.Builtins[builtinIndex]
			err := virtualMachine.push(builtinDefinition.Builtin)
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := virtualMachine.executeBinaryOperation(opcode)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := virtualMachine.executeComparison(opcode)
			if err != nil {
				return err
			}
		case code.OpArray:
			numElements := int(binary.BigEndian.Uint16(instructions[ip+1:]))
			virtualMachine.currentFrame().ip += 2
			array := virtualMachine.buildArray(virtualMachine.sp-numElements, virtualMachine.sp)
			virtualMachine.sp = virtualMachine.sp - numElements
			err := virtualMachine.push(array)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := virtualMachine.pop()
			expression := virtualMachine.pop()
			err := virtualMachine.executeIndexExpression(expression, index)
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			returnValue := virtualMachine.pop()
			frame := virtualMachine.popFrame()
			virtualMachine.sp = frame.basePointer
			virtualMachine.pop() // pop compiled function
			err := virtualMachine.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturnVoid:
			frame := virtualMachine.popFrame()
			virtualMachine.sp = frame.basePointer
			virtualMachine.pop() // pop compiled function
			err := virtualMachine.push(object.NULL)
			if err != nil {
				return err
			}
		case code.OpClosure:
			constantIndex := binary.BigEndian.Uint16(instructions[ip+1:])
			freeVariableArity := instructions[ip+3]
			virtualMachine.currentFrame().ip += 3
			function := virtualMachine.constants[constantIndex]
			freeVariables := make([]object.Object, freeVariableArity)
			for i := 0; i < int(freeVariableArity); i++ {
				freeVariables[i] = virtualMachine.stack[virtualMachine.sp-int(freeVariableArity)+i]
			}
			virtualMachine.sp = virtualMachine.sp - int(freeVariableArity)
			err := virtualMachine.push(&object.Closure{Function: function.(*object.CompiledFunction), FreeVariables: freeVariables})
			if err != nil {
				return err
			}
		case code.OpCall:
			argumentArity := instructions[ip+1]
			virtualMachine.currentFrame().ip += 1
			function := virtualMachine.stack[virtualMachine.sp-1-int(argumentArity)]
			switch callee := function.(type) {
			case *object.Closure:
				err := virtualMachine.callClosure(callee, int(argumentArity))
				if err != nil {
					return err
				}
			case *object.Builtin:
				err := virtualMachine.callBuiltin(callee, int(argumentArity))
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("calling non-function")
			}
		case code.OpCurrentClosure:
			err := virtualMachine.push(virtualMachine.currentFrame().closure)
			if err != nil {
				return err
			}
		case code.OpNegate, code.OpBang:
			err := virtualMachine.executeUnaryOperation(opcode)
			if err != nil {
				return err
			}
		case code.OpJump:
			jumpPosition := int(binary.BigEndian.Uint16(instructions[ip+1:]))
			virtualMachine.currentFrame().ip = jumpPosition - 1
		case code.OpJumpIfFalse:
			jumpPosition := int(binary.BigEndian.Uint16(instructions[ip+1:]))
			virtualMachine.currentFrame().ip += 2
			condition := virtualMachine.pop()
			if !isTruthy(condition) {
				virtualMachine.currentFrame().ip = jumpPosition - 1
			}
		case code.OpTrue:
			err := virtualMachine.push(object.TRUE)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := virtualMachine.push(object.FALSE)
			if err != nil {
				return err
			}
		case code.OpNull:
			err := virtualMachine.push(object.NULL)
			if err != nil {
				return err
			}
		case code.OpPop:
			virtualMachine.pop()
		default:
			return fmt.Errorf("unknown opcode %b", opcode)
		}
	}
	return nil
}

func (virtualMachine *VirtualMachine) push(object object.Object) error {
	if virtualMachine.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	virtualMachine.stack[virtualMachine.sp] = object
	virtualMachine.sp++
	return nil
}

func (virtualMachine *VirtualMachine) pop() object.Object {
	value := virtualMachine.stack[virtualMachine.sp-1]
	virtualMachine.sp--
	return value
}

func (virtualMachine *VirtualMachine) executeBinaryOperation(op code.Opcode) error {
	right := virtualMachine.pop()
	left := virtualMachine.pop()
	leftType := left.Type()
	rightType := right.Type()
	if leftType == object.INTEGER && rightType == object.INTEGER {
		return virtualMachine.executeBinaryIntegerOperation(op, *left.(*object.Integer), *right.(*object.Integer))
	}
	if leftType == object.STRING && rightType == object.STRING {
		return virtualMachine.executeBinaryStringOperation(op, *left.(*object.String), *right.(*object.String))
	}
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (virtualMachine *VirtualMachine) executeBinaryIntegerOperation(op code.Opcode, left, right object.Integer) error {
	var result int64
	switch op {
	case code.OpAdd:
		result = left.Value + right.Value
	case code.OpSub:
		result = left.Value - right.Value
	case code.OpMul:
		result = left.Value * right.Value
	case code.OpDiv:
		result = left.Value / right.Value
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}
	return virtualMachine.push(&object.Integer{Value: result})
}

func (virtualMachine *VirtualMachine) executeBinaryStringOperation(op code.Opcode, left, right object.String) error {
	var result string
	switch op {
	case code.OpAdd:
		result = left.Value + right.Value
	default:
		return fmt.Errorf("unknown string operator %T", op)
	}
	return virtualMachine.push(&object.String{Value: result})
}

func (virtualMachine *VirtualMachine) executeUnaryOperation(opcode code.Opcode) error {
	operand := virtualMachine.pop()
	switch opcode {
	case code.OpBang:
		return virtualMachine.executeBangOperation(operand)
	case code.OpNegate:
		return virtualMachine.executeNegateOperation(operand)
	default:
		return fmt.Errorf("unexpected opcode %T", opcode)
	}
}

func (virtualMachine *VirtualMachine) executeComparison(op code.Opcode) error {
	right := virtualMachine.pop()
	left := virtualMachine.pop()
	if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
		return virtualMachine.executeIntegerComparison(op, left, right)
	}
	switch op {
	case code.OpEqual:
		return virtualMachine.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return virtualMachine.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (virtualMachine *VirtualMachine) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch op {
	case code.OpEqual:
		return virtualMachine.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return virtualMachine.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return virtualMachine.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (virtualMachine *VirtualMachine) executeBangOperation(operand object.Object) error {
	switch operand {
	case object.TRUE:
		return virtualMachine.push(object.FALSE)
	case object.FALSE:
		return virtualMachine.push(object.TRUE)
	case object.NULL:
		return virtualMachine.push(object.TRUE)
	default:
		return virtualMachine.push(object.FALSE)
	}
}

func (virtualMachine *VirtualMachine) executeNegateOperation(operand object.Object) error {
	value, ok := operand.(*object.Integer)
	if !ok {
		return fmt.Errorf("expected *object.Integer. Got %T", value)
	}

	return virtualMachine.push(&object.Integer{Value: -value.Value})
}

func nativeBoolToBooleanObject(boolean bool) object.Object {
	if boolean {
		return object.TRUE
	}

	return object.FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (virtualMachine *VirtualMachine) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)
	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = virtualMachine.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (virtualMachine *VirtualMachine) executeIndexExpression(expression, index object.Object) error {
	switch {
	case expression.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return virtualMachine.executeArrayIndex(expression, index)
	default:
		return fmt.Errorf("index operator not supported: %s", expression.Type())
	}
}

func (virtualMachine *VirtualMachine) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	maximum := int64(len(arrayObject.Elements) - 1)
	if i < 0 || i > maximum {
		return virtualMachine.push(object.NULL)
	}
	return virtualMachine.push(arrayObject.Elements[i])
}

func (virtualMachine *VirtualMachine) callClosure(closure *object.Closure, arity int) error {
	if arity != closure.Function.ParameterArity {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			closure.Function.ParameterArity, arity)
	}
	frame := NewFrame(closure, virtualMachine.sp-arity)
	virtualMachine.pushFrame(frame)
	virtualMachine.sp = frame.basePointer + closure.Function.LocalVariableArity
	return nil
}

func (virtualMachine *VirtualMachine) callBuiltin(builtin *object.Builtin, arity int) error {
	args := virtualMachine.stack[virtualMachine.sp-arity : virtualMachine.sp]
	result := builtin.Function(args...)
	virtualMachine.sp = virtualMachine.sp - arity - 1
	if result != nil {
		return virtualMachine.push(result)
	}
	return virtualMachine.push(object.NULL)
}
