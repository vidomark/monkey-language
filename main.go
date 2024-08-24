package main

import (
	"fmt"
	"os"
	"os/user"
	"writing-in-interpreter-in-go/src/monkey/repl"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n",
		currentUser.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.StartRepl(os.Stdout, os.Stdin)
}
