package ast

import (
	"bytes"
	"writing-in-interpreter-in-go/src/monkey/token"
)

type IndexExpression struct {
	Token      *token.Token
	Expression Expression
	Index      Expression
}

func (indexExpression *IndexExpression) expressionNode() {}

func (indexExpression *IndexExpression) TokenLiteral() string { return indexExpression.Token.Literal }

func (indexExpression *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(indexExpression.Expression.String())
	out.WriteString("[")
	out.WriteString(indexExpression.Index.String())
	out.WriteString("])")
	return out.String()
}
