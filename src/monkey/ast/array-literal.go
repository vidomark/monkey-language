package ast

import (
	"bytes"
	"strings"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type ArrayLiteral struct {
	Token    *token.Token
	Elements []Expression
}

func (arrayLiteral *ArrayLiteral) expressionNode() {}

func (arrayLiteral *ArrayLiteral) TokenLiteral() string { return arrayLiteral.Token.Literal }

func (arrayLiteral *ArrayLiteral) String() string {
	var out bytes.Buffer
	var elements []string
	for _, el := range arrayLiteral.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
