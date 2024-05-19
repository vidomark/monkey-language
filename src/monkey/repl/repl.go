package repl

import (
	"bufio"
	"fmt"
	"io"
	"writing-in-interpreter-in-go/src/monkey/lexer"
	"writing-in-interpreter-in-go/src/monkey/token"
)

const PROMPT = ">>"

func Start(in io.Reader) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := lexer.NewLexer(scanner.Text())
		for tok := line.NextToken(); tok.Type != token.EOF; tok = line.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
