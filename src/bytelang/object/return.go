package object

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() Type { return RETURN }

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }
