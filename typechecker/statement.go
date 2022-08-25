package typechecker

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

/*
Returns true if the statement could be type checked successfully.
*/
func (tc *TypeChecker) CheckStatement(statement ast.Statement) error {
	switch targetStatement := statement.(type) {
	case *ast.ExpressionStatement:
		_, err := tc.CheckExpression(targetStatement.Value)

		return err
	case *ast.VariableDeclaration:
		return tc.checkVariableDeclaration(targetStatement)
	case *ast.BlockStatement:
		tc.context.EnterScope()
		// check every statement within a block
		for _, innerStatement := range targetStatement.Statements {
			err := tc.CheckStatement(innerStatement)

			if err != nil {
				tc.context.ExitScope()

				return err
			}
		}

		tc.context.ExitScope()

		return nil
	case *ast.IfStatement:
		return tc.checkIfStatement(targetStatement)
	case *ast.WhileStatement:
		return tc.checkWhileStatement(targetStatement)
	case *ast.StructDeclaration:
		return tc.checkStructStatement(targetStatement)
	case *ast.ReturnStatement:
		return tc.checkReturnStatement(targetStatement)
	}

	return CreateTypeError(
		fmt.Sprintf("Failed to type check unknown statement: %T\n", statement),
		0,
	)
}

func (tc *TypeChecker) checkVariableDeclaration(v *ast.VariableDeclaration) error {
	/*
		let x : int = 100
				 |	   |
		verify that the l-value type is
		equal to the r-value type.
	*/
	variableType := v.Type

	isValidVariable := tc.context.Add(v.Name, &variableType)

	if !isValidVariable {
		message := fmt.Sprintf("Variable '%s' already in scope.", v.Name)
		return CreateTypeError(message, v.Line)
	}

	if v.Value == nil {
		return CreateTypeError("Variable declaration cannot have a value of null.", v.Line)
	}

	valueType, valueErr := tc.CheckExpression(v.Value)

	if valueErr != nil {
		return valueErr
	}

	isEqual := tc.match(variableType, valueType)

	if !isEqual {
		message := fmt.Sprintf(
			"Invalid type in variable declaration. Expected %s but got %s.",
			variableType.String(),
			valueType.String(),
		)

		return CreateTypeError(message, v.Line)
	}

	return nil
}

/*
Check if the condition of an if statement is a boolean, and that
its body type checks properly.
*/
func (tc *TypeChecker) checkIfStatement(stat *ast.IfStatement) error {
	conditionType, conditionErr := tc.CheckExpression(stat.Condition)

	if conditionErr != nil {
		return conditionErr
	}

	if !tc.match(conditionType, ast.CreateTypeFromLiteral(lexer.BOOL)) {
		message := fmt.Sprintf(
			"Expected condition in 'if' statement to be boolean, got %s.",
			conditionType.String(),
		)

		return CreateTypeError(message, stat.Line)
	}

	if statementErr := tc.CheckStatement(stat.Body); statementErr != nil {
		return statementErr
	}

	if stat.ElseBody != nil {
		if statementErr := tc.CheckStatement(stat.ElseBody); statementErr != nil {
			return statementErr
		}
	}

	return nil
}

/*
Check if the condition is a boolean and that the body type checks properly.
*/
func (tc *TypeChecker) checkWhileStatement(stat *ast.WhileStatement) error {
	conditionType, conditionErr := tc.CheckExpression(stat.Condition)

	if conditionErr != nil {
		return conditionErr
	}

	if !tc.match(conditionType, ast.CreateTypeFromLiteral(lexer.BOOL)) {
		message := fmt.Sprintf(
			"Expected condition in 'while' statement to be boolean, got %s.",
			conditionType.String(),
		)

		return CreateTypeError(message, stat.Line)
	}

	if statementErr := tc.CheckStatement(stat.Body); statementErr != nil {
		return statementErr
	}

	return nil
}

func (tc *TypeChecker) checkStructStatement(stat *ast.StructDeclaration) error {
	isDefined, _ := tc.context.FindType(stat.Name)
	fields := make(map[string]ast.Type)

	if isDefined {
		message := fmt.Sprintf("Struct '%s' already defined.", stat.String())
		return CreateTypeError(message, stat.Line)
	}

	for variableName, variableType := range stat.Record.Variables {
		switch innerType := variableType.(type) {
		case *ast.VariableType:
			if !tc.context.TypeExists(innerType.Base) {
				isPrimitive, _ := ast.IsInternalType(variableType)

				if !isPrimitive {
					message := fmt.Sprintf(
						"Type '%s' does not exist in this context.",
						variableType.String(),
					)

					// TODO: update to use line number from the type field
					return CreateTypeError(message, 0)
				}
			}
		case *ast.FunctionType:
			continue
		case *ast.RecordType:
			continue
		}

		fields[variableName] = variableType
	}

	tc.context.AddType(stat.Name, ast.RecordType{Variables: fields})

	return nil
}

func (tc *TypeChecker) checkReturnStatement(stat *ast.ReturnStatement) error {
	if stat.Value == nil {
		return CreateTypeError("Return statement must have value.", stat.Line)
	}

	returnType, err := tc.CheckExpression(stat.Value)

	if err != nil {
		return err
	}

	stat.Type = returnType

	return nil
}

/*
Check if block has a return statement.
*/
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
Returns if there is a return statement and an error if it exists.
*/
func (tc *TypeChecker) checkLastReturnStatement(
	expectedType ast.Type,
	body ast.Statement,
) (bool, error) {
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
			return false, CreateTypeError(
				"Body does not have a return statement.",
				stat.Line,
			)
		}

		lastStatement := stat.Statements[len(stat.Statements)-1]

		if statement, ok := lastStatement.(*ast.ReturnStatement); ok {
			returnType := statement.Value.GetType().(ast.Type)

			if !tc.match(expectedType, returnType) {
				return false, CreateTypeError(
					fmt.Sprintf(
						"Expected function to return value of type %s, but instead returned %s.",
						expectedType.String(),
						returnType.String(),
					),
					statement.Line,
				)
			}

			return true, nil
		}

		return false, CreateTypeError(
			"The last statement of a function body must be a return statement.",
			stat.Line,
		)
	case *ast.ExpressionStatement:
		expressionType, err := tc.CheckExpression(stat.Value)

		if err != nil {
			return false, err
		}

		if !tc.match(expressionType, expressionType) {
			return false, CreateTypeError(
				fmt.Sprintf(
					"Expected function to return value of type %s, but instead returned %s.",
					expectedType.String(),
					expressionType.String(),
				),
				stat.Line,
			)
		}

		return true, nil
	}

	return false, CreateTypeError(
		"Checking for return statement on invalid statement.",
		0,
	)
}

/*
Validate type of return statement in a single statement. Check for nested
return statements as well.
*/
func (tc *TypeChecker) checkStatementForReturns(
	expectedType ast.Type,
	statement ast.Statement,
) error {
	switch statementType := statement.(type) {
	case *ast.ReturnStatement:
		returnType := statementType.Value.GetType()

		if returnType == nil {
			return CreateTypeError(
				fmt.Sprintf("Invalid return type."),
				statementType.Value.GetLine(),
			)
		}

		if !tc.match(expectedType, returnType) {
			return CreateTypeError(
				fmt.Sprintf("Expected incorrect return type."),
				statementType.Value.GetLine(),
			)
		}
	case *ast.IfStatement:
		if err := tc.checkStatementForReturns(expectedType, statementType.Body); err != nil {
			return err
		}

		if statementType.ElseBody != nil {
			if err := tc.checkStatementForReturns(expectedType, statementType.ElseBody); err != nil {
				return err
			}
		}
	case *ast.WhileStatement:
		if err := tc.checkStatementForReturns(expectedType, statementType.Body); err != nil {
			return err
		}
	case *ast.BlockStatement:
		if err := tc.checkAllReturnStatements(expectedType, statementType); err != nil {
			return err
		}
	}

	return nil
}

/*
Validate type of *every* return statement in a block.
*/
func (tc *TypeChecker) checkAllReturnStatements(
	expectedType ast.Type,
	body *ast.BlockStatement,
) error {
	// if block has a return, it must be the LAST statement.
	if tc.hasReturnStatement(body) {
		isValid, err := tc.checkLastReturnStatement(expectedType, body)

		if err != nil && !isValid {
			return err
		}
	}

	// check in any blocks for return statements, make sure that they are expectedType
	for _, statement := range body.Statements {
		if err := tc.checkStatementForReturns(expectedType, statement); err != nil {
			return err
		}
	}

	return nil
}
