package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"github.com/gmisail/glamlang/typechecker"
)

func main() {
	fileData, err := os.ReadFile("./typecheck.gl")

	if err != nil {
		panic(err)
	}

	l := lexer.ScanTokens(string(fileData))

	color.Blue("[glam] Done lexing.")
	statements := parser.Parse(l.Tokens)

	color.Blue("[glam] Done parsing.")
	checker := typechecker.CreateTypeChecker()

	for _, s := range statements {
		ok := checker.CheckStatement(s)

		if ok {
			fmt.Printf("VALID: %s\n", s.String())
		} else {
			fmt.Printf("INVALID: %s\n", s.String())
		}
	}

	color.Blue("[glam] Done type checking.")
}
