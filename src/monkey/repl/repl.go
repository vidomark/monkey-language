package repl

import (
	"bufio"
	"fmt"
	"io"
	"writing-in-interpreter-in-go/src/monkey/evaluator"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	"writing-in-interpreter-in-go/src/monkey/object"
	"writing-in-interpreter-in-go/src/monkey/parser"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environment := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		evaluated := evaluator.Eval(&program, environment)
		_, _ = io.WriteString(out, evaluated.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}
