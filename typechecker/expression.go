package typechecker

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

/*
	Checks the type of the expression and, if valid, returns true and its type. If there
	was an error while type checking, it will return false, nil.
*/
func (tc *TypeChecker) CheckExpression(expr ast.Expression) (Type, error) {
	switch exprType := expr.(type) {
	case *ast.Literal:
		return CreateTypeFromLiteral(exprType.Type), nil // literals always check successfully
	case *ast.VariableExpression:
		targetExists, targetType := tc.context.Find(exprType.Value)

		if !targetExists {
			return nil, CreateTypeError(fmt.Sprintf("[type] Can't find variable %s.\n", exprType.String()))
		}

		return *targetType, nil
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
		returnType := CreateTypeFrom(exprType.ReturnType)
		hasReturn, returnErr := tc.checkLastReturnStatement(returnType, exprType.Body)

		if !hasReturn {
			return nil, CreateTypeError("Function missing return statement.")
		}

		if body, ok := exprType.Body.(*ast.BlockStatement); ok {
			allReturnsValid := tc.checkAllReturnStatements(returnType, body)

			if !allReturnsValid {
				return nil, CreateTypeError("Some returns in function expression are invalid.")
			}
		}

		if returnErr != nil {
			color.Red(returnErr.Error())

			return nil, returnErr
		}

		tc.context.ExitScope()

		if !validBody {
			return nil, CreateTypeError("Invalid function body.")
		}

		return &FunctionType{Parameters: parameters, ReturnType: CreateTypeFrom(exprType.ReturnType)}, nil
	case *ast.FunctionCall:
		return tc.checkFunctionCall(exprType)
	case *ast.Group:
		return tc.CheckExpression(exprType.Value)
	case *ast.Binary:
		return tc.checkBinary(exprType)
	case *ast.Unary:
		return tc.checkUnary(exprType)
	case *ast.Logical:
		return tc.checkLogical(exprType)
	case *ast.GetExpression:
		return tc.checkGetExpression(exprType)
	}

	return nil, nil
}

func (tc *TypeChecker) checkGetExpression(expr *ast.GetExpression) (Type, error) {
	parentType, parentErr := tc.CheckExpression(expr.Parent)

	if parentErr != nil {
		return nil, parentErr
	}

	switch variableType := parentType.(type) {
	case *VariableType:
		typeName := variableType.Name
		typeExists, typeMembers := tc.context.environment.GetType(typeName)

		if !typeExists {
			return nil, CreateTypeError("Cannot access member variable from a non-existent struct type.")
		}

		memberExists, memberType := typeMembers.Find(expr.Name)

		if !memberExists {
			message := fmt.Sprintf(
				"[type] Member variable '%s' does not exist on type '%s'.",
				expr.Name,
				typeName,
			)

			return nil, CreateTypeError(message)
		}

		return *memberType, nil
	case *FunctionType:
		return nil, CreateTypeError("Cannot access a member variable of a function type.")
	}

	return nil, nil
}

func (tc *TypeChecker) checkFunctionCall(expr *ast.FunctionCall) (Type, error) {
	calleeType, calleeErr := tc.CheckExpression(expr.Callee)

	if calleeErr != nil {
		return nil, calleeErr
	}

	switch calleeVariableType := calleeType.(type) {
	case *VariableType:
		return nil, CreateTypeError("Cannot call instance of non-function.")
	case *FunctionType:
		var functionInstance FunctionType = *calleeVariableType

		if len(functionInstance.Parameters) != len(expr.Arguments) {
			message := fmt.Sprintf(
				"Function call only has %d arguments, expected %d.",
				len(functionInstance.Parameters),
				len(expr.Arguments),
			)

			return nil, CreateTypeError(message)
		}

		for i, param := range functionInstance.Parameters {
			argType, argErr := tc.CheckExpression(expr.Arguments[i])

			if argErr != nil {
				return nil, argErr
			}

			if !param.Equals(argType) {
				message := fmt.Sprintf(
					"[type] Expected type %s, got %s.",
					param.String(),
					argType.String(),
				)

				return nil, CreateTypeError(message)
			}
		}

		return functionInstance.ReturnType, nil
	}

	return nil, nil
}

func (tc *TypeChecker) checkBinary(expr *ast.Binary) (Type, error) {
	leftType, leftErr := tc.CheckExpression(expr.Left)

	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := tc.CheckExpression(expr.Right)

	if rightErr != nil {
		return nil, rightErr
	}

	isEqual := leftType.Equals(rightType)

	if !isEqual {
		message := fmt.Sprintf(
			"Types do not match in binary expression. Left type is %s while the right type is %s.",
			leftType.String(),
			rightType.String(),
		)

		return nil, CreateTypeError(message)
	}

	switch expr.Operator {
	case lexer.LT,
		lexer.LT_EQ,
		lexer.GT,
		lexer.GT_EQ,
		lexer.EQUALITY,
		lexer.NOT_EQUAL:
		return CreateTypeFromLiteral(lexer.BOOL), nil
	case lexer.ADD, lexer.SUB, lexer.MULT, lexer.DIV:
		return leftType, nil
	}

	return nil, nil
}

func (tc *TypeChecker) checkUnary(expr *ast.Unary) (Type, error) {
	switch expr.Operator {
	case lexer.BANG:
		// !(boolean expression)
		valueType, valueErr := tc.CheckExpression(expr.Value)

		if valueErr != nil {
			return nil, valueErr
		}

		if !valueType.Equals(CreateTypeFromLiteral(lexer.BOOL)) {
			message := fmt.Sprintf("[type] Expected type in 'not' expression to be bool, instead got incompatible type %s.", valueType.String())
			return nil, CreateTypeError(message)
		}

		return valueType, nil
	case lexer.SUB:
		// -(number)
		valueType, valueErr := tc.CheckExpression(expr.Value)

		if valueErr != nil {
			return nil, valueErr
		}

		if !valueType.Equals(CreateTypeFromLiteral(lexer.INT)) && !valueType.Equals(CreateTypeFromLiteral(lexer.FLOAT)) {
			message := fmt.Sprintf("[type] Expected type in negation to be int or float, instead got incompatible type %s.", valueType.String())
			return nil, CreateTypeError(message)
		}

		return valueType, nil
	}

	return nil, nil
}

func (tc *TypeChecker) checkLogical(expr *ast.Logical) (Type, error) {
	leftType, leftErr := tc.CheckExpression(expr.Left)

	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := tc.CheckExpression(expr.Right)

	if rightErr != nil {
		return nil, rightErr
	}

	boolType := CreateTypeFromLiteral(lexer.BOOL)

	if !leftType.Equals(boolType) {
		return nil, CreateTypeError(fmt.Sprintf("Expected the left side of logical statement to be of type bool, got %s.", leftType.String()))
	}

	if !rightType.Equals(boolType) {
		return nil, CreateTypeError(fmt.Sprintf("Expected the right side of logical statement to be of type bool, got %s.", rightType.String()))
	}

	return boolType, nil
}
