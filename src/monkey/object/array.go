package object

import (
	"bytes"
	"strings"
)

type Array struct{ Elements []Object }

func (ao *Array) Type() Type { return ARRAY }

func (ao *Array) Inspect() string {
	var out bytes.Buffer
	var elements []string
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
