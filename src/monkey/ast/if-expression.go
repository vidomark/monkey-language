package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type IfExpression struct {
	Token       *token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpression *IfExpression) expressionNode() {}

func (ifExpression *IfExpression) TokenLiteral() string {
	return ifExpression.Token.Literal
}

func (ifExpression *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ifExpression.Condition.String())
	out.WriteString(" ")
	out.WriteString(ifExpression.Consequence.String())
	if ifExpression.Alternative != nil {
		out.WriteString(ifExpression.Alternative.String())
	}
	return out.String()
}
