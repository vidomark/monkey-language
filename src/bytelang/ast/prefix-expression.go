package ast

import (
	"bytes"
	"bytelang/src/bytelang/token"
)

type PrefixExpression struct {
	Token    *token.Token
	Operator string
	Operand  Expression
}

func (prefixExpression *PrefixExpression) expressionNode() {}

func (prefixExpression *PrefixExpression) TokenLiteral() string {
	return prefixExpression.Token.Literal
}

func (prefixExpression *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(prefixExpression.Operator)
	out.WriteString(prefixExpression.Operand.String())
	out.WriteString(")")
	return out.String()
}
