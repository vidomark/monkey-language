package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type InfixExpression struct {
	Token    *token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExpression *InfixExpression) expressionNode() {}

func (infixExpression *InfixExpression) TokenLiteral() string {
	return infixExpression.Token.Literal
}

func (infixExpression *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(infixExpression.Left.String())
	out.WriteString(" " + infixExpression.Operator + " ")
	out.WriteString(infixExpression.Right.String())
	out.WriteString(")")
	return out.String()
}
