package object

import "fmt"

type Integer struct {
	Value int64
}

func (integer *Integer) Type() Type {
	return INTEGER
}

func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}
