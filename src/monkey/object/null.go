package object

import "fmt"

var (
	NULL = &Null{}
)

type Null struct{}

func (n *Null) Type() Type { return NULL_OBJ }

func (n *Null) Inspect() string { return "null" }

func (function *CompiledFunction) Type() Type { return COMPILED_FUNCTION }

func (function *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", function)
}
