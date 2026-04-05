package ast

import "bytelang/src/bytelang/token"

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
