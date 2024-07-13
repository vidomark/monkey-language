package object

import "fmt"

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Boolean struct{ Value bool }

func (b *Boolean) Type() Type { return BOOLEAN }

func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }
