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

/*
	Checks the type of the expression and, if valid, returns true and its type. If there
	was an error while type checking, it will return false, nil.
*/
func (tc *TypeChecker) CheckExpression(expr ast.Expression) (bool, Type) {
	switch exprType := expr.(type) {
	case *ast.Literal:
		return true, CreateTypeFromLiteral(exprType.Type) // literals always check successfully
	case *ast.VariableExpression:
		targetExists, targetType := tc.context.Find(exprType.Value)

		if !targetExists {
			color.Red(fmt.Sprintf("[type] Can't find variable %s.\n", exprType.String()))

			return false, nil
		}

		return true, *targetType
	case *ast.FunctionExpression:
		// validate that the body of the function is valid
		tc.context.EnterScope()

		parameters := make([]Type, len(exprType.Parameters))

		// push the parameters into scope
		for i, param := range exprType.Parameters {
			paramType := CreateTypeFrom(param.Type)
			parameters[i] = paramType

			if !tc.context.Add(param.Name, &paramType) {
				color.Red("[type] Variable '%s' already exists in this scope.", param.Name)
			}
		}

		validBody := tc.CheckStatement(exprType.Body)

		// validate that the valid being return matches the return type

		tc.context.ExitScope()

		if !validBody {
			return false, nil
		}

		return true, &FunctionType{Parameters: parameters, ReturnType: CreateTypeFrom(exprType.ReturnType)}
	case *ast.FunctionCall:
		targetExists, targetType := tc.checkFunctionCall(exprType)

		if !targetExists {
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
	case *ast.Unary:
		targetExists, targetType := tc.checkUnary(exprType)

		if !targetExists {
			return false, nil
		}

		return true, targetType
	case *ast.Logical:
		targetExists, targetType := tc.checkLogical(exprType)

		if !targetExists {
			return false, nil
		}

		return true, targetType
	case *ast.GetExpression:
		return tc.checkGetExpression(exprType)
		/*
			- FUNCTION CALL
		*/
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

func (tc *TypeChecker) checkGetExpression(expr *ast.GetExpression) (bool, Type) {
	isParentValid, parentType := tc.CheckExpression(expr.Parent)

	if !isParentValid {
		return false, nil
	}

	switch variableType := parentType.(type) {
	case *VariableType:
		typeName := variableType.Name
		typeExists, typeMembers := tc.context.environment.GetType(typeName)

		if !typeExists {
			color.Red("[type] Cannot access member variable from a non-existent struct type.")

			return false, nil
		}

		memberExists, memberType := typeMembers.Find(expr.Name)

		if !memberExists {
			color.Red(
				"[type] Member variable '%s' does not exist on type '%s'.",
				expr.Name,
				typeName,
			)

			return false, nil
		}

		return true, *memberType
	case *FunctionType:
		color.Red("[type] Cannot access a member variable of a function type.")
	}

	return false, nil
}

func (tc *TypeChecker) checkFunctionCall(expr *ast.FunctionCall) (bool, Type) {
	isCalleeValid, calleeType := tc.CheckExpression(expr.Callee)

	if !isCalleeValid {
		return false, nil
	}

	switch calleeVariableType := calleeType.(type) {
	case *VariableType:
		color.Red("[type] Cannot call variable instance as function.")
		return false, nil
	case *FunctionType:
		var functionInstance FunctionType = *calleeVariableType

		if len(functionInstance.Parameters) != len(expr.Arguments) {
			color.Red(
				"[type] Function call only has %d arguments, expected %d.",
				len(functionInstance.Parameters),
				len(expr.Arguments),
			)

			return false, nil
		}

		for i, param := range functionInstance.Parameters {
			argValid, argType := tc.CheckExpression(expr.Arguments[i])

			if !argValid {
				return false, nil
			}

			if !param.Equals(argType) {
				color.Red(
					"[type] Expected type %s, got %s.",
					param.String(),
					argType.String(),
				)

				return false, nil
			}
		}

		return true, functionInstance.ReturnType
	}

	return false, nil
}

func (tc *TypeChecker) checkBinary(expression *ast.Binary) (bool, Type) {
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
			leftType.String(),
			rightType.String(),
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

func (tc *TypeChecker) checkUnary(expr *ast.Unary) (bool, Type) {
	switch expr.Operator {
	case lexer.BANG:
		// !(boolean expression)
		validType, valueType := tc.CheckExpression(expr.Value)

		if !validType || !valueType.Equals(CreateTypeFromLiteral(lexer.BOOL)) {
			return false, nil
		}

		return true, valueType
	case lexer.SUB:
		// -(number)
		validType, valueType := tc.CheckExpression(expr.Value)

		if !validType || !valueType.Equals(CreateTypeFromLiteral(lexer.INT)) || !valueType.Equals(CreateTypeFromLiteral(lexer.FLOAT)) {
			return false, nil
		}

		return true, valueType
	}

	return false, nil
}

func (tc *TypeChecker) checkLogical(expr *ast.Logical) (bool, Type) {
	isLeftValid, leftType := tc.CheckExpression(expr.Left)
	isRightValid, rightType := tc.CheckExpression(expr.Right)

	if !isLeftValid || !isRightValid {
		color.Red("[type] l-value or r-value failed type checking.")

		return false, nil
	}

	boolType := CreateTypeFromLiteral(lexer.BOOL)

	if !(leftType.Equals(boolType) && rightType.Equals(boolType)) {
		return false, nil
	}

	return true, boolType
}
