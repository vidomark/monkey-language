package ast

import (
	"bytes"
	"fmt"
	"strings"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type FunctionLiteral struct {
	Token      token.Token
	Parameters []Expression
	Body       *BlockStatement
	Name       string
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
	if functionLiteral.Name != "" {
		out.WriteString(fmt.Sprintf("<%s>", functionLiteral.Name))
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(functionLiteral.Body.String())
	return out.String()
}
