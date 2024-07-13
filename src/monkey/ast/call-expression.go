package ast

import (
	"bytes"
	"strings"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type CallExpression struct {
	Function  Expression
	Token     token.Token
	Arguments []Expression
}

func (callExpression *CallExpression) expressionNode() {}

func (callExpression *CallExpression) TokenLiteral() string {
	return callExpression.Token.Literal
}

func (callExpression *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range callExpression.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(callExpression.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
