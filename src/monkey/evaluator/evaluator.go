package evaluator

import (
	"fmt"
	"writing-in-interpreter-in-go/src/monkey/ast"
	"writing-in-interpreter-in-go/src/monkey/object"
)

var Builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	"push":  object.GetBuiltinByName("push"),
}

func Eval(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, environment)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, environment)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Expression, environment)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, environment)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.BlockStatement:
		return evaluateBlockStatement(node, environment)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.LetStatement:
		expression := Eval(node.Value, environment)
		if isError(expression) {
			return expression
		}
		environment.Set(node.Identifier.Value, expression)
		return expression
	case *ast.ReturnStatement:
		returnValue := Eval(node.ReturnValue, environment)
		if isError(returnValue) {
			return returnValue
		}
		return &object.ReturnValue{Value: returnValue}
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters:  node.Parameters,
			Body:        node.Body,
			Environment: environment,
		}
	case *ast.CallExpression:
		if node.Function.TokenLiteral() == "quote" {
			return quote(node.Arguments[0], environment)
		}
		function := Eval(node.Function, environment)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyArguments(function, args)
	case *ast.PrefixExpression:
		operand := Eval(node.Operand, environment)
		if isError(operand) {
			return operand
		}
		return evalPrefixOperator(operand, node.Operator)
	case *ast.InfixExpression:
		leftOperand := Eval(node.Left, environment)
		if isError(leftOperand) {
			return leftOperand
		}

		rightOperand := Eval(node.Right, environment)
		if isError(rightOperand) {
			return rightOperand
		}
		return evalInfixExpression(node.Operator, leftOperand, rightOperand)
	}
	return nil
}

func evalStatements(stmts []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, environment)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalIdentifier(identifier *ast.Identifier, environment *object.Environment) object.Object {
	if value, ok := environment.Get(identifier.Value); ok {
		return value
	}

	builtin := object.GetBuiltinByName(identifier.Value)
	if builtin != nil {
		return builtin
	}

	return newError("Identifier not found: %s", identifier.Value)
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maximum := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > maximum {
		return object.NULL
	}
	return arrayObject.Elements[idx]
}

func evalPrefixOperator(operand object.Object, operator string) object.Object {
	switch operator {
	case `!`:
		return evalBangOperator(operand)
	case `-`:
		return evalNegateOperator(operand)
	default:
		return newError("unknown operator: %s%s", operator, operand.Type())
	}
}

func evalInfixExpression(
	operator string,
	leftOperand, rightOperand object.Object,
) object.Object {
	switch {
	case leftOperand.Type() == object.INTEGER && rightOperand.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, leftOperand, rightOperand)
	case leftOperand.Type() == object.STRING && rightOperand.Type() == object.STRING:
		return evalStringInfixExpression(operator, leftOperand, rightOperand)
	case operator == "==":
		return nativeBoolToBooleanObject(leftOperand == rightOperand)
	case operator == "!=":
		return nativeBoolToBooleanObject(leftOperand != rightOperand)
	default:
		return newError("type mismatch: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

func evalStringInfixExpression(operator string, leftOperand object.Object, rightOperand object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
	left := leftOperand.(*object.String)
	right := rightOperand.(*object.String)
	return &object.String{Value: left.Value + right.Value}
}

func evalBangOperator(operand object.Object) object.Object {
	switch operand {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	}
	return object.FALSE
}

func evalNegateOperator(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER {
		return newError("unknown operator: -%s", operand.Type())
	}

	integer, ok := operand.(*object.Integer)
	if !ok {
		panic("Expected *object.Integer.")
	}

	return &object.Integer{Value: -integer.Value}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluateBlockStatement(node *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object
	for _, statement := range node.Statements {
		result = Eval(statement, environment)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return result
			}
		}
	}
	return result
}

func evalIfExpression(node *ast.IfExpression, environment *object.Environment) object.Object {
	condition := Eval(node.Condition, environment)

	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Consequence, environment)
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, environment)
	}

	return object.NULL
}

func evalExpressions(expressions []ast.Expression, environment *object.Environment) []object.Object {
	var result []object.Object
	for _, expression := range expressions {
		evaluated := Eval(expression, environment)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyArguments(function object.Object, args []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		env := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, env)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if result := fn.Function(args...); result != nil {
			return result
		}
		return object.NULL
	default:
		return newError("unsupported function type: %T", function)
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	returnValue, ok := obj.(*object.ReturnValue)
	if ok {
		return returnValue.Value
	}
	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Environment)
	for paramIdx, param := range fn.Parameters {
		identifier := param.(*ast.Identifier)
		env.Set(identifier.Value, args[paramIdx])
	}
	return env
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}
