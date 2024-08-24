package compiler

import (
	"fmt"
	"writing-in-interpreter-in-go/src/monkey/ast"
	"writing-in-interpreter-in-go/src/monkey/code"
	"writing-in-interpreter-in-go/src/monkey/object"
)

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIndex  int
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type CompilationScope struct {
	instructions           code.Instructions
	lastInstruction        EmittedInstruction
	penultimateInstruction EmittedInstruction
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:           code.Instructions{},
		lastInstruction:        EmittedInstruction{},
		penultimateInstruction: EmittedInstruction{},
	}
	symbolTable := NewSymbolTable()
	for index, builtin := range object.Builtins {
		symbolTable.DefineBuiltin(index, builtin.Name)
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (compiler *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := compiler.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.BlockStatement:
		for _, statement := range node.Statements {
			err := compiler.Compile(statement)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := compiler.Compile(node.Expression)
		if err != nil {
			return err
		}
		compiler.emit(code.OpPop)
	case *ast.LetStatement:
		err := compiler.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := compiler.symbolTable.Define(node.Identifier.Value)
		if symbol.Scope == GlobalScope {
			compiler.emit(code.OpSetGlobal, symbol.Index)
		}

		if symbol.Scope == LocalScope {
			compiler.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := compiler.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		compiler.loadSymbol(symbol)
	case *ast.PrefixExpression:
		err := compiler.Compile(node.Operand)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			compiler.emit(code.OpBang)
		case "-":
			compiler.emit(code.OpNegate)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := compiler.Compile(node.Right)
			if err != nil {
				return err
			}
			err = compiler.Compile(node.Left)
			if err != nil {
				return err
			}
			compiler.emit(code.OpGreaterThan)
			return nil
		}
		err := compiler.Compile(node.Left)
		if err != nil {
			return err
		}
		err = compiler.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			compiler.emit(code.OpAdd)
		case "-":
			compiler.emit(code.OpSub)
		case "*":
			compiler.emit(code.OpMul)
		case "/":
			compiler.emit(code.OpDiv)
		case ">":
			compiler.emit(code.OpGreaterThan)
		case "==":
			compiler.emit(code.OpEqual)
		case "!=":
			compiler.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		err := compiler.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpIfFalsePosition := compiler.emit(code.OpJumpIfFalse, -1)

		err = compiler.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if compiler.lastInstructionIs(code.OpPop) {
			compiler.removeLastInstruction()
		}

		jumpPosition := compiler.emit(code.OpJump, -1)

		afterConsequencePosition := len(compiler.currentInstructions())
		compiler.replaceOperand(jumpIfFalsePosition, afterConsequencePosition)

		if node.Alternative != nil {
			err := compiler.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if compiler.lastInstructionIs(code.OpPop) {
				compiler.removeLastInstruction()
			}
		} else {
			compiler.emit(code.OpNull)
		}

		alternativePosition := len(compiler.currentInstructions())
		compiler.replaceOperand(jumpPosition, alternativePosition)
	case *ast.BooleanExpression:
		boolean := object.Boolean{Value: node.Value}
		if boolean.Value {
			compiler.emit(code.OpTrue)
		} else {
			compiler.emit(code.OpFalse)
		}
	case *ast.IntegerLiteral:
		integer := object.Integer{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(&integer))
	case *ast.StringLiteral:
		stringLiteral := object.String{Value: node.Value}
		compiler.emit(code.OpConstant, compiler.addConstant(&stringLiteral))
	case *ast.ReturnStatement:
		err := compiler.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		compiler.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := compiler.Compile(node.Function)
		if err != nil {
			return err
		}
		for _, argument := range node.Arguments {
			err := compiler.Compile(argument)
			if err != nil {
				return err
			}
		}

		compiler.emit(code.OpCall, len(node.Arguments))
	case *ast.FunctionLiteral:
		compiler.enterScope()

		for _, parameter := range node.Parameters {
			compiler.symbolTable.Define(parameter.(*ast.Identifier).Value)
		}

		err := compiler.Compile(node.Body)
		if err != nil {
			return err
		}

		if compiler.lastInstructionIs(code.OpPop) {
			compiler.replaceLastPopWithReturn()
		}
		if !compiler.lastInstructionIs(code.OpReturnValue) {
			compiler.emit(code.OpReturnVoid)
		}

		symbolTable, instruction := compiler.leaveScope()

		for _, symbol := range symbolTable.FreeSymbols {
			compiler.loadSymbol(symbol)
		}

		function := &object.CompiledFunction{
			Instructions:       instruction,
			LocalVariableArity: symbolTable.numDefinitions,
			ParameterArity:     len(node.Parameters),
		}
		compiler.emit(code.OpClosure, compiler.addConstant(function), len(symbolTable.FreeSymbols))
	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := compiler.Compile(el)
			if err != nil {
				return err
			}
		}
		compiler.emit(code.OpArray, len(node.Elements))
	case *ast.IndexExpression:
		err := compiler.Compile(node.Expression)
		if err != nil {
			return err
		}

		err = compiler.Compile(node.Index)
		if err != nil {
			return err
		}

		compiler.emit(code.OpIndex)
	}
	return nil
}

func (compiler *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: compiler.currentInstructions(),
		Constants:    compiler.constants,
	}
}

func (compiler *Compiler) emit(opcode code.Opcode, operands ...int) int {
	instructions := code.Make(opcode, operands...)
	position := compiler.addInstruction(instructions)
	compiler.setLastInstruction(opcode, position)
	return position
}

func (compiler *Compiler) addInstruction(instructions []byte) int {
	position := len(compiler.currentInstructions())
	updatedInstructions := append(compiler.currentInstructions(), instructions...)
	compiler.setCurrentInstruction(updatedInstructions)
	return position
}

func (compiler *Compiler) addConstant(node object.Object) int {
	position := len(compiler.constants)
	compiler.constants = append(compiler.constants, node)
	return position
}

func (compiler *Compiler) currentScope() CompilationScope {
	return compiler.scopes[compiler.scopeIndex]
}

func (compiler *Compiler) currentInstructions() code.Instructions {
	return compiler.currentScope().instructions
}

func (compiler *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(compiler.currentInstructions()) == 0 {
		return false
	}
	return compiler.currentScope().lastInstruction.Opcode == op
}

func (compiler *Compiler) removeLastInstruction() {
	last := compiler.currentScope().lastInstruction
	previous := compiler.currentScope().penultimateInstruction

	newInstructions := compiler.currentInstructions()[:last.Position]

	compiler.setCurrentInstruction(newInstructions)
	compiler.setLastEmittedInstruction(previous)
}

func (compiler *Compiler) replaceOperand(opcodePosition int, operand int) {
	opcode := code.Opcode(compiler.currentInstructions()[opcodePosition])
	instructions := code.Make(opcode, operand)
	compiler.replaceInstruction(opcodePosition, instructions)
}

func (compiler *Compiler) replaceInstruction(position int, instruction []byte) {
	instructions := compiler.currentInstructions()
	for i := 0; i < len(instruction); i++ {
		instructions[position+i] = instruction[i]
	}
}

func (compiler *Compiler) setCurrentInstruction(instructions code.Instructions) {
	compiler.scopes[compiler.scopeIndex].instructions = instructions
}

func (compiler *Compiler) setLastEmittedInstruction(instructions EmittedInstruction) {
	compiler.scopes[compiler.scopeIndex].lastInstruction = instructions
}

func (compiler *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := compiler.scopes[compiler.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	compiler.scopes[compiler.scopeIndex].penultimateInstruction = previous
	compiler.scopes[compiler.scopeIndex].lastInstruction = last
}

func (compiler *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:           code.Instructions{},
		lastInstruction:        EmittedInstruction{},
		penultimateInstruction: EmittedInstruction{},
	}
	compiler.symbolTable = NewEnclosedSymbolTable(compiler.symbolTable)
	compiler.scopes = append(compiler.scopes, scope)
	compiler.scopeIndex++
}
func (compiler *Compiler) leaveScope() (*SymbolTable, code.Instructions) {
	instructions := compiler.currentInstructions()
	localSymbolTable := compiler.symbolTable
	compiler.symbolTable = localSymbolTable.Enclosing
	compiler.scopes = compiler.scopes[:len(compiler.scopes)-1]
	compiler.scopeIndex--
	return localSymbolTable, instructions
}

func (compiler *Compiler) replaceLastPopWithReturn() {
	lastPos := compiler.currentScope().lastInstruction.Position
	compiler.replaceInstruction(lastPos, code.Make(code.OpReturnValue))
	compiler.scopes[compiler.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (compiler *Compiler) loadSymbol(symbol Symbol) {
	if symbol.Scope == GlobalScope {
		compiler.emit(code.OpGetGlobal, symbol.Index)
	}

	if symbol.Scope == LocalScope {
		compiler.emit(code.OpGetLocal, symbol.Index)
	}

	if symbol.Scope == BuiltinScope {
		compiler.emit(code.OpGetBuiltin, symbol.Index)
	}

	if symbol.Scope == FreeScope {
		compiler.emit(code.OpGetFree, symbol.Index)
	}
}
