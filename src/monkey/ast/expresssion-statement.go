package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type ExpressionStatement struct {
	Token      *token.Token
	Expression Expression
}

func (statement *ExpressionStatement) statementNode() {}

func (statement *ExpressionStatement) TokenLiteral() string {
	return statement.Token.Literal
}

func (statement *ExpressionStatement) String() string {
	var out bytes.Buffer
	if statement.Expression != nil {
		out.WriteString(statement.Expression.String())
	}
	return out.String()
}
