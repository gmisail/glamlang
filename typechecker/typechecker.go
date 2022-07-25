package typechecker

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"

	"github.com/fatih/color"
)

type TypeChecker struct {
	context *Context
}

func CreateTypeChecker() *TypeChecker {
	return &TypeChecker{context: CreateContext()}
}

/*
	Returns true if the statement could be type checked successfully.
*/
func (tc *TypeChecker) CheckStatement(statement ast.Statement) bool {
	switch statementType := statement.(type) {
	case *ast.VariableDeclaration:
		return tc.checkVariableDeclaration(statementType)
	}

	color.Green(fmt.Sprintf("Failed to type check unknown statement.\n"))

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
		return true, CreateTypeFromLiteral(exprType.Type) // literals always check successfully default:
	case *ast.VariableExpression:
		targetExists, targetType := tc.context.Find(exprType.Value)

		if !targetExists {
			color.Red(fmt.Sprintf("[type] Can't find variable %s.\n", exprType.Value))

			return false, nil
		}

		return true, targetType
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
		color.Red(fmt.Sprintf("[type] Variable %s already in scope.\n", v.Name))

		return false
	}

	validType, valueType := tc.CheckExpression(v.Value)

	if !validType {
		color.Red(fmt.Sprintf("[type] Invalid type in variable declaration. Got %s but expected %s.\n", valueType, variableType))

		return false
	}

	isEqual := variableType.Equals(valueType)

	return isEqual
}
