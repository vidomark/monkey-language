package evaluator

import "writing-in-interpreter-in-go/src/monkey/object"

type BuiltinFunction func(args ...object.Object) object.Object

type Builtin struct {
	Function BuiltinFunction
}

func (builtin *Builtin) Type() object.Type { return object.BUILT_IN }

func (builtin *Builtin) Inspect() string { return "builtin function" }
