package main

import (
	"os"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := ScanTokens(string(fileData))

	for _, tok := range l.Tokens {
		tok.Print()
	}
}
