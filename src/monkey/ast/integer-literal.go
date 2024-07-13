package ast

import "writing-in-interpreter-in-go/src/monkey/token"

type IntegerLiteral struct {
	Token *token.Token
	Value int64
}

func (integerLiteral *IntegerLiteral) expressionNode() {}

func (integerLiteral *IntegerLiteral) TokenLiteral() string {
	return integerLiteral.Token.Literal
}

func (integerLiteral *IntegerLiteral) String() string {
	return integerLiteral.Token.Literal
}
