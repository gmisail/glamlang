package main

import (
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"github.com/gmisail/glamlang/typechecker"
)

func main() {
	fileData, err := os.ReadFile("./functions.gl")

	if err != nil {
		panic(err)
	}

	start := time.Now()
	l := lexer.ScanTokens(string(fileData))
	color.Blue("[glam] Done lexing in %s.", time.Since(start))

	start = time.Now()
	ok, statements := parser.Parse(l.Tokens)
	color.Blue("[glam] Done parsing in %s.", time.Since(start))

	if !ok {
		return
	}

	checker := typechecker.CreateTypeChecker()
	start = time.Now()

	checker.CheckAll(statements)

	color.Blue("[glam] Done type checking in %s.", time.Since(start))
}
