package main

import "github.com/gmisail/glamlang/lexer"

func main() {
	lexer.ScanTokens("[ ] { }")
}
