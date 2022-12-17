package typechecker

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
	"github.com/martinusso/inflect"
)

/*
Checks the type of the expression and, if valid, returns true and its type. If there
was an error while type checking, it will return false, nil.
*/
func (tc *TypeChecker) CheckExpression(expr ast.Expression) (ast.Type, error) {
	switch exprType := expr.(type) {
	case *ast.Literal:
		literalType := ast.CreateTypeFromLiteral(exprType.LiteralType)
		exprType.Type = literalType

		// literals always check successfully
		return literalType, nil
	case *ast.VariableExpression:
		targetExists, targetType := tc.context.FindVariable(exprType.Value)

		if !targetExists {
			return nil, CreateTypeError(
				fmt.Sprintf("Undefined variable '%s'.", exprType.Value),
				exprType.Line,
			)
		}

		return *targetType, nil
	case *ast.FunctionExpression:
		// validate that the body of the function is valid
		tc.context.EnterScope()

		parameters := make([]ast.Type, len(exprType.Parameters))

		// push the parameters into scope
		for i, param := range exprType.Parameters {
			paramType := param.Type
			parameters[i] = paramType

			if !tc.context.Add(param.Name, &paramType) {
				color.Red("Variable '%s' already exists in this scope.", param.Name)
			}
		}

		bodyErr := tc.CheckStatement(exprType.Body)

		if bodyErr != nil {
			tc.context.ExitScope()
			return nil, bodyErr
		}

		returnType := exprType.ReturnType
		hasReturn, returnErr := tc.checkLastReturnStatement(returnType, exprType.Body)

		if !hasReturn || returnErr != nil {
			tc.context.ExitScope()
			return nil, returnErr
		}

		if body, ok := exprType.Body.(*ast.BlockStatement); ok {
			if returnErr := tc.checkAllReturnStatements(returnType, body); returnErr != nil {
				tc.context.ExitScope()
				return nil, returnErr
			}
		}

		if returnErr != nil {
			tc.context.ExitScope()
			return nil, returnErr
		}

		tc.context.ExitScope()

		return &ast.FunctionType{Parameters: parameters, ReturnType: exprType.ReturnType}, nil
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
	case *ast.RecordInstance:
		return tc.checkRecordInstance(exprType)
	}

	return nil, nil
}

/*
 * Recursively expand every Record name into its expanded form.
 *
 * For instance, "User" --> { name: string, bio: string }
 */
func (tc *TypeChecker) resolve(record ast.RecordType) ast.Type {
	resolvedFields := make(map[string]ast.Type)

	for field, fieldType := range record.Fields {
		if subRecord, ok := fieldType.(*ast.VariableType); ok {
			isRecord, recordFields := tc.context.FindType(subRecord.Base)

			if isRecord {
				resolvedFields[field] = tc.resolve(*recordFields)

				continue
			}
		}

		// if the field isn't a record, keep the original type.
		resolvedFields[field] = fieldType
	}

	return &ast.RecordType{Fields: resolvedFields}
}

/**
 * Checks if two types match. Prior to comparing the types, it will
 * resolve the type in case it is a user-defined record.
 */
func (tc *TypeChecker) match(first ast.Type, second ast.Type) bool {
	firstType := first
	secondType := second

	if v, ok := first.(*ast.VariableType); ok {
		isRecord, recordFields := tc.context.FindType(v.Base)

		if isRecord {
			firstType = tc.resolve(*recordFields)
		}
	}

	if v, ok := second.(*ast.VariableType); ok {
		isRecord, recordFields := tc.context.FindType(v.Base)

		if isRecord {
			secondType = tc.resolve(*recordFields)
		}
	}

	return firstType.Equals(secondType)
}

func (tc *TypeChecker) checkRecordInstance(record *ast.RecordInstance) (ast.Type, error) {
	fields := make(map[string]ast.Type)

	for field, fieldValue := range record.Values {
		fieldType, typeErr := tc.CheckExpression(fieldValue)

		if typeErr != nil {
			return nil, typeErr
		}

		fields[field] = fieldType
	}

	return &ast.RecordType{Fields: fields}, nil
}

func (tc *TypeChecker) checkGetExpression(expr *ast.GetExpression) (ast.Type, error) {
	parentType, parentErr := tc.CheckExpression(expr.Parent)

	if parentErr != nil {
		return nil, parentErr
	}

	switch variableType := parentType.(type) {
	case *ast.VariableType:
		typeName := variableType.Base
		typeExists, typeMembers := tc.context.FindType(typeName)

		if !typeExists {
			return nil, CreateTypeError(
				fmt.Sprintf("Cannot access member variable from a non-existent type '%s'.", typeName),
				expr.Line,
			)
		}

		memberType, memberExists := typeMembers.Fields[expr.Name]

		if !memberExists {
			message := fmt.Sprintf(
				"Member variable '%s' does not exist on type '%s'.",
				expr.Name,
				typeName,
			)

			return nil, CreateTypeError(message, expr.Line)
		}

		expr.Type = memberType

		return memberType, nil
	case *ast.FunctionType:
		return nil, CreateTypeError(
			"Cannot access a member variable of a function type.",
			expr.Line,
		)
	}

	return nil, nil
}

func (tc *TypeChecker) checkFunctionCall(expr *ast.FunctionCall) (ast.Type, error) {
	calleeType, calleeErr := tc.CheckExpression(expr.Callee)

	if calleeErr != nil {
		return nil, calleeErr
	}

	switch calleeVariableType := calleeType.(type) {
	case *ast.VariableType:
		return nil, CreateTypeError(
			"Cannot call instance of non-function.",
			expr.Line,
		)
	case *ast.FunctionType:
		var functionInstance ast.FunctionType = *calleeVariableType

		if len(functionInstance.Parameters) != len(expr.Arguments) {
			// TODO: diff the two functions, i.e. what was expected vs. what it got

			message := fmt.Sprintf(
				"Function call has %d arguments, got %d.",
				len(functionInstance.Parameters),
				len(expr.Arguments),
			)

			return nil, CreateTypeError(message, expr.Line)
		}

		for i, param := range functionInstance.Parameters {
			argType, argErr := tc.CheckExpression(expr.Arguments[i])

			if argErr != nil {
				return nil, argErr
			}

			if !tc.match(param, argType) {
				message := fmt.Sprintf(
					"Expected type of %d%s argument to be %s, got %s.",
					i+1,
					inflect.Ordinal(i+1),
					param.String(),
					argType.String(),
				)

				return nil, CreateTypeError(message, expr.Line)
			}
		}

		expr.Type = functionInstance.ReturnType

		return functionInstance.ReturnType, nil
	}

	return nil, nil
}

func (tc *TypeChecker) checkBinary(expr *ast.Binary) (ast.Type, error) {
	leftType, leftErr := tc.CheckExpression(expr.Left)

	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := tc.CheckExpression(expr.Right)

	if rightErr != nil {
		return nil, rightErr
	}

	isEqual := tc.match(leftType, rightType)

	if !isEqual {
		message := fmt.Sprintf(
			"Types do not match in binary expression. Left type is %s while the right type is %s.",
			leftType.String(),
			rightType.String(),
		)

		return nil, CreateTypeError(message, expr.Line)
	}

	isValid := HasBinaryRule(expr.Operator, leftType.String())

	if !isValid {
		message := fmt.Sprintf(
			"Cannot apply operation '%s' to type %s.",
			lexer.GetSymbol(expr.Operator),
			leftType.String(),
		)

		return nil, CreateTypeError(message, expr.Line)
	}

	switch expr.Operator {
	case lexer.LT,
		lexer.LT_EQ,
		lexer.GT,
		lexer.GT_EQ,
		lexer.EQUALITY,
		lexer.NOT_EQUAL:
		boolType := ast.CreateTypeFromLiteral(lexer.BOOL)
		expr.Type = boolType

		return boolType, nil
	case lexer.ADD, lexer.SUB, lexer.MULT, lexer.DIV:
		expr.Type = leftType

		return leftType, nil
	}

	return nil, nil
}

func (tc *TypeChecker) checkUnary(expr *ast.Unary) (ast.Type, error) {
	valueType, valueErr := tc.CheckExpression(expr.Value)

	if valueErr != nil {
		return nil, valueErr
	}

	isValid := HasUnaryRule(expr.Operator, valueType.String())

	if !isValid {
		message := fmt.Sprintf(
			"Cannot apply operation '%s' to type %s.",
			lexer.GetSymbol(expr.Operator),
			valueType.String(),
		)

		return nil, CreateTypeError(message, expr.Line)
	}

	switch expr.Operator {
	case lexer.BANG:
		// !(boolean expression)
		if !tc.match(valueType, ast.CreateTypeFromLiteral(lexer.BOOL)) {
			message := fmt.Sprintf(
				"Expected type in 'not' expression to be bool, instead got incompatible type %s.",
				valueType.String(),
			)

			return nil, CreateTypeError(message, expr.Line)
		}

		expr.Type = valueType

		return valueType, nil
	case lexer.SUB:
		// -(number)
		if !tc.match(valueType, ast.CreateTypeFromLiteral(lexer.INT)) &&
			!tc.match(valueType, ast.CreateTypeFromLiteral(lexer.FLOAT)) {
			message := fmt.Sprintf(
				"Expected type in negation to be int or float, instead got incompatible type %s.",
				valueType.String(),
			)

			return nil, CreateTypeError(message, expr.Line)
		}

		expr.Type = valueType

		return valueType, nil
	}

	return nil, nil
}

func (tc *TypeChecker) checkLogical(expr *ast.Logical) (ast.Type, error) {
	leftType, leftErr := tc.CheckExpression(expr.Left)

	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := tc.CheckExpression(expr.Right)

	if rightErr != nil {
		return nil, rightErr
	}

	boolType := ast.CreateTypeFromLiteral(lexer.BOOL)

	if !tc.match(leftType, boolType) {
		return nil, CreateTypeError(
			fmt.Sprintf(
				"Expected the left side of logical statement to be of type bool, got %s.",
				leftType.String(),
			),
			expr.Line,
		)
	}

	if !tc.match(rightType, boolType) {
		return nil, CreateTypeError(
			fmt.Sprintf(
				"Expected the right side of logical statement to be of type bool, got %s.",
				rightType.String(),
			),
			expr.Line,
		)
	}

	expr.Type = boolType

	return boolType, nil
}
