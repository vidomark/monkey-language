package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (blockStatement *BlockStatement) statementNode() {}

func (blockStatement *BlockStatement) TokenLiteral() string {
	return blockStatement.Token.Literal
}

func (blockStatement *BlockStatement) String() string {
	var out bytes.Buffer
	for _, statement := range blockStatement.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}
