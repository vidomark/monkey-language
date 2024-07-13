package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type ReturnStatement struct {
	Token      *token.Token
	Expression Expression
}

func (returnStatement *ReturnStatement) statementNode() {}

func (returnStatement *ReturnStatement) TokenLiteral() string {
	return returnStatement.Token.Literal
}

func (returnStatement *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(returnStatement.TokenLiteral() + " ")
	if returnStatement.Expression != nil {
		out.WriteString(returnStatement.Expression.String())
	}
	return out.String()
}
