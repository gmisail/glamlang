package main

import (
	"fmt"
	"os"
	"reflect"
)

func main() {
	fileData, err := os.ReadFile("./demo.gl")

	if err != nil {
		panic(err)
	}

	l := ScanTokens(string(fileData))

	statements := Parse(l.Tokens)

	for _, s := range statements {
		fmt.Println(reflect.TypeOf(s))
	}

	fmt.Println(len(statements))
}
