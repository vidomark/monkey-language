package main

import (
	"fmt"
	"os"
	"os/user"
	"bytelang/src/bytelang/repl"
)

func main() {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Bytelang programming language!\n",
		currentUser.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.StartRepl(os.Stdout, os.Stdin)
}
