package object

import "fmt"

type String struct {
	Value string
}

func (integer *String) Type() Type {
	return STRING
}

func (integer *String) Inspect() string {
	return fmt.Sprintf("%s", integer.Value)
}
