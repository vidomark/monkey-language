package object

import (
	"bytes"
	"strings"
	"writing-in-interpreter-in-go/src/monkey/ast"
)

type Macro struct {
	Parameters []ast.Expression
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() Type { return MACRO }
func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range m.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n}")

	return out.String()
}
