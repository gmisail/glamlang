package main

import (
	"fmt"
	"os"

	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"github.com/gmisail/glamlang/typechecker"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := lexer.ScanTokens(string(fileData))

	statements := parser.Parse(l.Tokens)
	checker := typechecker.CreateTypeChecker()

	for _, s := range statements {
		ok := checker.CheckStatement(s)

		if ok {
			fmt.Printf("VALID: %s\n", s.String())
		} else {
			fmt.Printf("INVALID: %s\n", s.String())
		}
	}
}
