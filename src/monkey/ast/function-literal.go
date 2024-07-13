package ast

import (
	"bytes"
	"strings"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type FunctionLiteral struct {
	Token      token.Token
	Parameters []Expression
	Body       *BlockStatement
}

func (functionLiteral *FunctionLiteral) expressionNode() {}

func (functionLiteral *FunctionLiteral) TokenLiteral() string {
	return functionLiteral.Token.Literal
}

func (functionLiteral *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range functionLiteral.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(functionLiteral.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(functionLiteral.Body.String())
	return out.String()
}
