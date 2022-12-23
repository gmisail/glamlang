package main

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gmisail/glamlang/codegen"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"github.com/gmisail/glamlang/typechecker"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must provide file name!")

		return
	}

	fileName := os.Args[1]

	start := time.Now()
	l := lexer.ScanTokens(fileName)
	color.Blue("[glam] Done lexing in %s.", time.Since(start))

	start = time.Now()
	ok, statements := parser.Parse(l, l.Tokens)
	color.Blue("[glam] Done parsing in %s.", time.Since(start))

	if !ok {
		return
	}

	start = time.Now()
	checker := typechecker.CreateTypeChecker()
	checker.CheckAll(statements)

	color.Blue("[glam] Done type checking in %s.", time.Since(start))

	start = time.Now()

	codegen.Compile(codegen.GetNativeBackend(), statements)

	color.Blue("[glam] Done compiling in %s.", time.Since(start))
}
