package main

import (
	"fmt"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"os"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := lexer.ScanTokens(string(fileData))

	statements := parser.Parse(l.Tokens)

	for _, s := range statements {
		//fmt.Println(reflect.TypeOf(s))

		fmt.Println(s.String())
	}

	fmt.Println(len(statements))
}
