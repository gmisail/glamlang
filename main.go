package main

import (
	"fmt"
	"os"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := ScanTokens(string(fileData))

	statements := Parse(l.Tokens)

	fmt.Println(len(statements))
}
