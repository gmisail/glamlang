package main

import "github.com/gmisail/glamlang/lexer"

func main() {
	lexer.ScanTokens("[ ] { }           longer_keyword hello world _123_keyword_is_valid")
}
