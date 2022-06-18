package main

import (
	"github.com/gmisail/glamlang/lexer"
)

func main() {
	lexer.ScanTokens("123 456.789 432.1")
}
