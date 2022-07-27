package typechecker

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"

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
	switch targetStatement := statement.(type) {
	case *ast.ExpressionStatement:
		isValid, _ := tc.CheckExpression(targetStatement)
		return isValid
	case *ast.VariableDeclaration:
		return tc.checkVariableDeclaration(targetStatement)
	case *ast.BlockStatement:
		// check every statement within a block
		for _, innerStatement := range targetStatement.Statements {
			if !tc.CheckStatement(innerStatement) {
				color.Green("failed on statement: %s", innerStatement.String())

				return false
			}
		}

		return true
	case *ast.IfStatement:
		return tc.checkIfStatement(targetStatement)
	case *ast.WhileStatement:
		return tc.checkWhileStatement(targetStatement)
	case *ast.StructDeclaration:
		return tc.checkStructStatement(targetStatement)
	}

	color.Green(fmt.Sprintf("Failed to type check unknown statement: %T\n", statement))

	// if the switch fails, then this is an unknown statement.
	return false
}

/*
	Check if the condition of an if statement is a boolean, and that
	its body type checks properly.
*/
func (tc *TypeChecker) checkIfStatement(stat *ast.IfStatement) bool {
	isValidType, conditionType := tc.CheckExpression(stat.Condition)

	if !isValidType {
		return false
	}

	if !conditionType.Equals(CreateTypeFromLiteral(lexer.BOOL)) {
		color.Red("[type] Expected condition in 'if' statement to be boolean, got %s.", conditionType.Name)

		return false
	}

	if !tc.CheckStatement(stat.Body) {
		return false
	}

	if stat.ElseBody != nil {
		if !tc.CheckStatement(stat.ElseBody) {
			return false
		}
	}

	return true
}

/*
	Check if the condition is a boolean and that the body type checks properly.
*/
func (tc *TypeChecker) checkWhileStatement(stat *ast.WhileStatement) bool {
	isValidType, conditionType := tc.CheckExpression(stat.Condition)

	if !isValidType {
		return false
	}

	if !conditionType.Equals(CreateTypeFromLiteral(lexer.BOOL)) {
		color.Red("[type] Expected condition in 'while' statement to be boolean, got %s.", conditionType.Name)

		return false
	}

	if !tc.CheckStatement(stat.Body) {
		return false
	}

	return true
}

func (tc *TypeChecker) checkStructStatement(stat *ast.StructDeclaration) bool {
	isUnique, structEnv := tc.context.environment.AddType(stat.Name)

	if !isUnique {
		color.Red("[type] Struct '%s' already defined.", stat.Name)

		return false
	}

	for _, structVariable := range stat.Variables {
		variableName := structVariable.Name
		variableType := CreateTypeFrom(structVariable.Type)

		if !tc.context.environment.CustomTypeExists(variableType.Name) {
			isPrimitive, _ := IsInternalType(structVariable.Type)

			if !isPrimitive {
				color.Red("[type] Type '%s' does not exist in this context.", variableType.Name)
				return false
			}
		}

		structEnv.Add(variableName, variableType)
	}

	return true
}

/*
	Checks the type of the expression and, if valid, returns true and its type. If there
	was an error while type checking, it will return false, nil.
*/
func (tc *TypeChecker) CheckExpression(expr ast.Expression) (bool, *Type) {
	switch exprType := expr.(type) {
	case *ast.Literal:
		return true, CreateTypeFromLiteral(exprType.Type) // literals always check successfully
	case *ast.VariableExpression:
		targetExists, targetType := tc.context.Find(exprType.Value)

		if !targetExists {
			color.Red(fmt.Sprintf("[type] Can't find variable %s.\n", exprType.Value))

			return false, nil
		}

		return true, targetType
	case *ast.Group:
		return tc.CheckExpression(exprType.Value)
	case *ast.Binary:
		targetExists, targetType := tc.checkBinary(exprType)

		if !targetExists {
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
		color.Red("[type] error in nested expression")

		return false
	}

	isEqual := variableType.Equals(valueType)

	if !isEqual {
		color.Red("[type] Invalid type in variable declaration. Expected %s but got %s.\n", valueType.Name, variableType.Name)

		return false
	}

	return isEqual
}

func (tc *TypeChecker) checkBinary(expression *ast.Binary) (bool, *Type) {
	isLeftValid, leftType := tc.CheckExpression(expression.Left)
	isRightValid, rightType := tc.CheckExpression(expression.Right)

	if !isLeftValid || !isRightValid {
		color.Red("[type] l-value or r-value failed type checking.")

		return false, nil
	}

	isEqual := leftType.Equals(rightType)

	if !isEqual {
		color.Red(
			"[type] Types do not match in binary expression. Left type is %s while the right type is %s.",
			leftType.Name,
			rightType.Name,
		)

		return false, nil
	}

	switch expression.Operator {
	case lexer.LT,
		lexer.LT_EQ,
		lexer.GT,
		lexer.GT_EQ,
		lexer.EQUALITY,
		lexer.NOT_EQUAL:
		return true, CreateTypeFromLiteral(lexer.BOOL)
	case lexer.ADD, lexer.SUB, lexer.MULT, lexer.DIV:
		return true, leftType
	}

	return true, nil
}
