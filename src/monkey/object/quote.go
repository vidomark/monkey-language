package object

import "writing-in-interpreter-in-go/src/monkey/ast"

type Quote struct {
	Node ast.Node
}

func (quote *Quote) Type() Type {
	return QUOTE
}

func (quote *Quote) Inspect() string {
	return "QUOTE(" + quote.Node.String() + ")"
}
