package main

import (
	"os"

	"github.com/gmisail/glamlang/lexer"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := lexer.ScanTokens(string(fileData))

	for _, tok := range l.Tokens {
		tok.Print()
	}
}
