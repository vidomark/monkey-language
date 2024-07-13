package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	return ""
}
