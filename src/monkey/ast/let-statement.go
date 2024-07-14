package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type LetStatement struct {
	Token      *token.Token
	Identifier *Identifier
	Value      Expression
}

func (statement *LetStatement) statementNode() {}

func (statement *LetStatement) TokenLiteral() string {
	return (*statement.Token).Literal
}

func (statement *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(statement.TokenLiteral() + " ")
	out.WriteString((*statement.Identifier).String())
	out.WriteString(" = ")
	if statement.Value != nil {
		out.WriteString(statement.Value.String())
	}
	return out.String()
}
