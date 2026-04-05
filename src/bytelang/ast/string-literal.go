package ast

import "bytelang/src/bytelang/token"

type StringLiteral struct {
	Token *token.Token
	Value string
}

func (stringLiteral *StringLiteral) expressionNode() {}

func (stringLiteral *StringLiteral) TokenLiteral() string {
	return stringLiteral.Token.Literal
}

func (stringLiteral *StringLiteral) String() string {
	return stringLiteral.Token.Literal
}
