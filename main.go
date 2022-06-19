package main

import (
	"fmt"
	"os"

	"github.com/gmisail/glamlang/lexer"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := lexer.ScanTokens(string(fileData))

	fmt.Println(l.Tokens)
}
