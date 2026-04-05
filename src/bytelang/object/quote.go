package object

import "bytelang/src/bytelang/ast"

type Quote struct {
	Node ast.Node
}

func (quote *Quote) Type() Type {
	return QUOTE
}

func (quote *Quote) Inspect() string {
	return "QUOTE(" + quote.Node.String() + ")"
}
