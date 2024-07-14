package ast

type ModifierFunc func(node Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, statement := range node.Statements {
			node.Statements[i] = Modify(statement, modifier).(Statement)
		}
	case *BlockStatement:
		for i, statement := range node.Statements {
			node.Statements[i] = Modify(statement, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Right = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Operand = Modify(node.Operand, modifier).(Expression)
	case *IfExpression:
		node.Condition = Modify(node.Condition, modifier).(Expression)
		node.Consequence = Modify(node.Consequence, modifier).(*BlockStatement)
		node.Alternative = Modify(node.Alternative, modifier).(*BlockStatement)
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i] = Modify(node.Parameters[i], modifier).(Expression)
		}
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i], _ = Modify(node.Elements[i], modifier).(Expression)
		}
	}

	return modifier(node)
}
