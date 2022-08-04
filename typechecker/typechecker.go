package typechecker

import (
	"fmt"
	"github.com/gmisail/glamlang/ast"
)

type TypeChecker struct {
	context *Context
}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{context: CreateContext()}
}

func (tc *TypeChecker) CheckAll(statements []ast.Statement) bool {
	isValid := true

	for _, s := range statements {
		ok := tc.CheckStatement(s)

		if ok {
			fmt.Printf("VALID: %s\n", s.String())
		} else {
			isValid = false
			fmt.Printf("INVALID: %s\n", s.String())
		}
	}

	return isValid
}
