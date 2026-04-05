package repl

import (
	"bufio"
	"bytelang/src/bytelang/compiler"
	"bytelang/src/bytelang/lexer"
	"bytelang/src/bytelang/object"
	"bytelang/src/bytelang/parser"
	"bytelang/src/bytelang/vm"
	"fmt"
	"io"
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
		byteCode := c.ByteCode()
		constants = byteCode.Constants
		virtualMachine := vm.NewWithGlobalsStore(byteCode, globals)
		err = virtualMachine.Run()
		if err != nil {
			_, _ = fmt.Fprintf(out, "Executing bytecode failed: \n %s\n", err)
		}
		lastPopped := virtualMachine.LastPopped()
		_, _ = io.WriteString(out, lastPopped.Inspect())
		_, _ = io.WriteString(out, "\n")
	}
}
