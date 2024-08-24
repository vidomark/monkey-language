package object

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Function BuiltinFunction
}

func (builtin *Builtin) Type() Type { return BUILT_IN }

func (builtin *Builtin) Inspect() string { return "builtin function" }
