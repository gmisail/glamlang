package typechecker

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

/*
	Returns true if the statement could be type checked successfully.
*/
func (tc *TypeChecker) CheckStatement(statement ast.Statement) bool {
	switch targetStatement := statement.(type) {
	case *ast.ExpressionStatement:
		isValid, _ := tc.CheckExpression(targetStatement.Value)

		return isValid
	case *ast.VariableDeclaration:
		return tc.checkVariableDeclaration(targetStatement)
	case *ast.BlockStatement:
		tc.context.EnterScope()
		// check every statement within a block
		for _, innerStatement := range targetStatement.Statements {
			if !tc.CheckStatement(innerStatement) {
				color.Green("failed on statement: %s", innerStatement.String())

				tc.context.ExitScope()

				return false
			}
		}

		tc.context.ExitScope()

		return true
	case *ast.IfStatement:
		return tc.checkIfStatement(targetStatement)
	case *ast.WhileStatement:
		return tc.checkWhileStatement(targetStatement)
	case *ast.StructDeclaration:
		return tc.checkStructStatement(targetStatement)
	case *ast.ReturnStatement:
		return tc.checkReturnStatement(targetStatement)
	}

	color.Green(fmt.Sprintf("Failed to type check unknown statement: %T\n", statement))

	// if the switch fails, then this is an unknown statement.
	return false
}

func (tc *TypeChecker) checkVariableDeclaration(v *ast.VariableDeclaration) bool {
	/*
		let x : int = 100
				 |	   |
		verify that the l-value type is
		equal to the r-value type.
	*/
	variableType := CreateTypeFrom(v.Type)
	isValidVariable := tc.context.Add(v.Name, &variableType)

	if !isValidVariable {
		color.Red(fmt.Sprintf("[type] Variable %s already in scope.\n", v.String()))

		return false
	}

	validType, valueType := tc.CheckExpression(v.Value)

	if !validType {
		color.Red("[type] error in nested expression")

		return false
	}

	isEqual := variableType.Equals(valueType)

	if !isEqual {
		color.Red("[type] Invalid type in variable declaration. Expected %s but got %s.\n", variableType.String(), valueType.String())

		return false
	}

	return isEqual
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
		color.Red("[type] Expected condition in 'if' statement to be boolean, got %s.", conditionType.String())

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
		color.Red("[type] Expected condition in 'while' statement to be boolean, got %s.", conditionType.String())

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
		color.Red("[type] Struct '%s' already defined.", stat.String())

		return false
	}

	for _, structVariable := range stat.Variables {
		variableName := structVariable.Name
		variableType := CreateTypeFrom(structVariable.Type)

		switch innerType := variableType.(type) {
		case *VariableType:
			if !tc.context.environment.CustomTypeExists(innerType.Name) {
				isPrimitive, _ := IsInternalType(structVariable.Type)

				if !isPrimitive {
					color.Red("[type] Type '%s' does not exist in this context.", variableType.String())
					return false
				}
			}
		case *FunctionType:
			continue
		}

		structEnv.Add(variableName, &variableType)
	}

	return true
}

func (tc *TypeChecker) checkReturnStatement(stat *ast.ReturnStatement) bool {
	if stat.Value == nil {
		return true
	}

	ok, _ := tc.CheckExpression(stat.Value)

	if !ok {
		return false
	}

	return ok
}

func (tc *TypeChecker) hasReturnStatement(body *ast.BlockStatement) bool {
	if len(body.Statements) <= 0 {
		return false
	}

	for _, statement := range body.Statements {
		switch statement.(type) {
		case *ast.ReturnStatement:
			return true
		}
	}

	return false
}

/*
	Returns if there is a return statement and if there is an error.
*/
func (tc *TypeChecker) checkLastReturnStatement(expectedType Type, body ast.Statement) (bool, error) {
	/*
		TODO: clean up this code, error messages, etc...
	*/

	/*
			There are two options:

			fun (x, y): int => x * y
		                 ^^^^^^
						ExpressionStatement

			fun (x, y): int => {
				...
				return x * y
			}
			^^^^^^^^^^^^^^^^^^
			BlockStatement
	*/
	switch stat := body.(type) {
	case *ast.BlockStatement:
		if len(stat.Statements) <= 0 {
			return false, errors.New("[type] Body does not have a return statement.")
		}

		lastStatement := stat.Statements[len(stat.Statements)-1]

		switch statement := lastStatement.(type) {
		case *ast.ReturnStatement:
			_, returnType := tc.CheckExpression(statement.Value)

			if !expectedType.Equals(returnType) {
				return false, errors.New(
					fmt.Sprintf(
						"[type] Expected function to return value of type %s, but instead returned %s.",
						expectedType.String(),
						returnType.String(),
					))
			}

			return true, nil
		}

		return false, errors.New("[type] The last statement of a function body must be a return statement.")
	case *ast.ExpressionStatement:
		isValid, expressionType := tc.CheckExpression(stat.Value)

		if !isValid {
			return false, errors.New("[type] Error in returned expression.")
		}

		if !expressionType.Equals(expressionType) {
			return false, errors.New(
				fmt.Sprintf(
					"[type] Expected function to return value of type %s, but instead returned %s.",
					expectedType.String(),
					expressionType.String(),
				))
		}

		return true, nil
	}

	return false, errors.New("Checking for return statement on invalid statement.")
}

func (tc *TypeChecker) checkStatementForReturns(expectedType Type, statement ast.Statement) bool {
	switch statementType := statement.(type) {
	case *ast.ReturnStatement:
		_, returnType := tc.CheckExpression(statementType.Value)

		if !expectedType.Equals(returnType) {
			return false
		}
	case *ast.IfStatement:
		if !tc.checkStatementForReturns(expectedType, statementType.Body) {
			return false
		}

		if statementType.ElseBody != nil {
			if !tc.checkStatementForReturns(expectedType, statementType.ElseBody) {
				return false
			}
		}
	case *ast.WhileStatement:
		if !tc.checkStatementForReturns(expectedType, statementType.Body) {
			return false
		}
	case *ast.BlockStatement:
		if !tc.checkAllReturnStatements(expectedType, statementType) {
			return false
		}
	}

	return true
}

func (tc *TypeChecker) checkAllReturnStatements(expectedType Type, body *ast.BlockStatement) bool {
	// if block has a return, it must be the LAST statement.
	if tc.hasReturnStatement(body) {
		isValid, err := tc.checkLastReturnStatement(expectedType, body)

		if err != nil && !isValid {
			color.Red(err.Error())
			// TODO: handle error
			return false
		}
	}

	// check in any blocks for return statements, make sure that they are expectedType
	for _, statement := range body.Statements {
		if !tc.checkStatementForReturns(expectedType, statement) {
			return false
		}
	}

	return true
}
