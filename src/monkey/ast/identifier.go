package ast

import "writing-in-interpreter-in-go/src/monkey/token"

type Identifier struct {
	Token *token.Token
	Value string
}

func (identifier *Identifier) expressionNode() {}

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}

func (identifier *Identifier) String() string {
	return identifier.Value
}
