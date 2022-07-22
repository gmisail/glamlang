package typechecker

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
)

type TypeChecker struct {
	context *Context
}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{context: CreateGlobal()}
}

/*
	Returns true if the statement could be type checked successfully.
*/
func (tc *TypeChecker) CheckStatement(statement ast.Statement) bool {
	switch statementType := statement.(type) {
	case *ast.VariableDeclaration:
		return tc.checkVariableDeclaration(statementType)
	}

	fmt.Println("Failed to type check unknown statement.")

	// if the switch fails, then this is an unknown statement.
	return false
}

/*
	Checks the type of the expression and, if valid, returns true and its type. If there
	was an error while type checking, it will return false, nil.
*/
func (tc *TypeChecker) CheckExpression(expr ast.Expression) (bool, *Type) {
	switch exprType := expr.(type) {
	case *ast.Literal:
		return true, CreateType("int", false) // literals always check successfully
	default:
		fmt.Println(exprType)
	}

	return false, nil
}

func (tc *TypeChecker) checkVariableDeclaration(v *ast.VariableDeclaration) bool {
	/*
		let x : int = 100
				 |	   |
		verify that the l-value type is
		equal to the r-value type.
	*/
	variableType := CreateTypeFrom(v.Type)
	isValidVariable := tc.context.Add(v.Name, variableType)

	if !isValidVariable {
		// TODO: more graceful error handling
		fmt.Printf("[type] Variable %s already in scope.\n", v.Name)

		return false
	}

	validType, valueType := tc.CheckExpression(v.Value)

	if !validType {
		fmt.Printf("[type] Invalid type in variable declaration. Got %s but expected %s.\n", valueType, variableType)

		return false
	}

	isEqual := variableType.Equals(valueType)

	return isEqual
}
