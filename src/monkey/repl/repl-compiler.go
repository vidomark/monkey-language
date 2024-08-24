package repl

import (
	"bufio"
	"fmt"
	"io"
	"writing-in-interpreter-in-go/src/monkey/compiler"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	"writing-in-interpreter-in-go/src/monkey/object"
	"writing-in-interpreter-in-go/src/monkey/parser"
	"writing-in-interpreter-in-go/src/monkey/vm"
)

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	var constants []object.Object
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()
	for index, builtin := range object.Builtins {
		symbolTable.DefineBuiltin(index, builtin.Name)
	}

	for {
		_, _ = fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		input := scanner.Text()
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		c := compiler.NewWithState(symbolTable, constants)
		err := c.Compile(program)
		if err != nil {
			return
		}
		virtualMachine := vm.NewWithGlobalsStore(c.ByteCode(), globals)
		err = virtualMachine.Run()
		if err != nil {
			_, _ = fmt.Fprintf(out, "Executing bytecode failed: \n %s\n", err)
		}
		lastPopped := virtualMachine.LastPopped()
		_, _ = io.WriteString(out, lastPopped.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}
