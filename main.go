package main

import (
	"fmt"
	"github.com/gmisail/glamlang/lexer"
)

func main() {
	l := lexer.ScanTokens("=> == !=")

	fmt.Println(l.Tokens)
}
