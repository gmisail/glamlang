package main

import (
	"fmt"

	"github.com/gmisail/glamlang/lexer"
)

func main() {
	l := lexer.ScanTokens(`hello "this is a very long test"`)

	fmt.Println(l.Tokens)
}
