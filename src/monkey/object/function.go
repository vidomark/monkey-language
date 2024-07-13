package object

import (
	"bytes"
	"strings"
	"writing-in-interpreter-in-go/src/monkey/ast"
)

type Function struct {
	Parameters  []ast.Expression
	Body        *ast.BlockStatement
	Environment *Environment
}

func (function *Function) Type() Type {
	return FUNCTION
}

func (function *Function) Inspect() string {
	var out bytes.Buffer
	var params []string
	for _, parameter := range function.Parameters {
		params = append(params, parameter.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(function.Body.String())
	out.WriteString("\n}")
	return out.String()
}
