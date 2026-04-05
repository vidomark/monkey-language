package ast

import (
	"strconv"
	"bytelang/src/bytelang/token"
)

type BooleanExpression struct {
	Token *token.Token
	Value bool
}

func (booleanExpression *BooleanExpression) expressionNode() {}

func (booleanExpression *BooleanExpression) TokenLiteral() string {
	return booleanExpression.Token.Literal
}

func (booleanExpression *BooleanExpression) String() string {
	return strconv.FormatBool(booleanExpression.Value)
}
