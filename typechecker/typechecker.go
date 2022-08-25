package typechecker

import (
	//	"fmt"
	"github.com/fatih/color"
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/context"
)

type TypeChecker struct {
	context *context.Context
}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{context: context.CreateContext()}
}

func (tc *TypeChecker) CheckAll(statements []ast.Statement) bool {
	isValid := true

	for _, s := range statements {
		err := tc.CheckStatement(s)

		if err != nil {
			color.Red(err.Error())
			//fmt.Printf("INVALID: %s\n", s.String())
		}
	}

	return isValid
}
