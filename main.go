package main

import (
	//	"os"
	"fmt"
)

func main() {
	/*	fileData, err := os.ReadFile("./demo.gl")

		if err != nil {
			panic(err)
		}
	*/
	l := ScanTokens("5 * (100 - 5)")

	ast, err := Parse(l.Tokens)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ast.String())
}
