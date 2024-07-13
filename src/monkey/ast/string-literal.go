package ast

import "writing-in-interpreter-in-go/src/monkey/token"

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
