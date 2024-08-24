package object

import "fmt"

type Closure struct {
	Function      *CompiledFunction
	FreeVariables []Object
}

func (closure *Closure) Type() Type {
	return CLOSURE
}

func (closure *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", closure)
}
