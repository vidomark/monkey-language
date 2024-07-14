package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type ReturnStatement struct {
	Token       *token.Token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode() {}

func (returnStatement *ReturnStatement) TokenLiteral() string {
	return returnStatement.Token.Literal
}

func (returnStatement *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(returnStatement.TokenLiteral() + " ")
	if returnStatement.ReturnValue != nil {
		out.WriteString(returnStatement.ReturnValue.String())
	}
	return out.String()
}
