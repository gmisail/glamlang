package parser

import (
	"fmt"
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (p *Parser) parsePrimary() (ast.Expression, error) {
	token := p.CurrentToken()

	if p.MatchToken(lexer.FALSE) {
		return &ast.Literal{Value: false, Type: lexer.BOOL, Line: token.Line}, nil
	} else if p.MatchToken(lexer.TRUE) {
		return &ast.Literal{Value: true, Type: lexer.BOOL, Line: token.Line}, nil
	} else if p.MatchToken(lexer.NULL) {
		return &ast.Literal{Value: nil, Type: lexer.NULL, Line: token.Line}, nil
	} else if p.MatchToken(lexer.STRING, lexer.INT, lexer.FLOAT) {
		return &ast.Literal{Value: token.Literal, Type: token.Type, Line: token.Line}, nil
	} else if p.MatchToken(lexer.IDENTIFIER) {
		return &ast.VariableExpression{Value: p.PreviousToken().Literal, Line: token.Line}, nil
	} else if p.MatchToken(lexer.L_PAREN) {
		expr, _ := p.parseExpression()
		//benabenabenabenabenabenabenabenabenabenabenabenabenabena
		//yayayayayayayayayayaayayayayayaya
		_, err := p.Consume(lexer.R_PAREN, "Expected closing parenthesis for group expression.")

		if err != nil {
			return nil, err
		}

		return &ast.Group{Value: expr, Line: token.Line}, nil
	}

	return nil, &ParseError{message: fmt.Sprintf("Unexpected token: %s", lexer.TokenTypeToString(p.CurrentToken().Type)), line: p.CurrentToken().Line}
}

func (p *Parser) finishParseCall(startLine int, callee ast.Expression) (ast.Expression, error) {
	arguments := make([]ast.Expression, 0)

	for p.CurrentToken().Type != lexer.R_PAREN {
		argument, argumentErr := p.parseExpression()

		if argumentErr != nil {
			return nil, argumentErr
		}

		arguments = append(arguments, argument)

		// if the next token is ), then there are no more arguments
		// and we shouldn't expect to see another comma.
		if p.CurrentToken().Type != lexer.R_PAREN {
			_, commaErr := p.Consume(lexer.COMMA, "Expected comma after arguments in function call.")

			if commaErr != nil {
				return nil, commaErr
			}
		}
	}

	_, rightParenErr := p.Consume(lexer.R_PAREN, "Expected ')' after function call arguments.")

	if rightParenErr != nil {
		return nil, rightParenErr
	}

	return &ast.FunctionCall{Callee: callee, Arguments: arguments, Line: startLine}, nil
}

func (p *Parser) parseCall() (ast.Expression, error) {
	expr, primaryErr := p.parsePrimary()

	if primaryErr != nil {
		return nil, primaryErr
	}

	var callErr error = nil

	for {
		if p.MatchToken(lexer.L_PAREN) {
			line := p.PreviousToken().Line
			expr, callErr = p.finishParseCall(line, expr)

			if callErr != nil {
				return nil, callErr
			}
		} else if p.MatchToken(lexer.PERIOD) {
			// TODO: handle me
			line := p.PreviousToken().Line
			name, _ := p.Consume(lexer.IDENTIFIER, "Expected identifier after '.'")
			expr = &ast.GetExpression{Name: name.Literal, Parent: expr, Line: line}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) parseFunction() (ast.Expression, error) {
	if p.MatchToken(lexer.FUNCTION) {
		parameters := make([]ast.VariableDeclaration, 0)

		leftParen, leftParenErr := p.Consume(lexer.L_PAREN, "Expected '('")

		if leftParenErr != nil {
			return nil, leftParenErr
		}

		line := leftParen.Line

		// if there's a right parenthesis, that means the function doesn't have any parameters.
		if !p.MatchToken(lexer.R_PAREN) {
			for {
				// no more parameters :(
				if !p.MatchToken(lexer.COMMA) && p.MatchToken(lexer.R_PAREN) {
					break
				}

				parameter, parameterErr := p.Consume(lexer.IDENTIFIER, "Expected parameter name.")

				if parameterErr != nil {
					return nil, parameterErr
				}

				_, colonErr := p.Consume(lexer.COLON, "Expected ':' before type.")

				if colonErr != nil {
					return nil, colonErr
				}

				parameterType, parameterTypeErr := p.parseTypeDeclaration()

				if parameterTypeErr != nil {
					return nil, parameterTypeErr
				}

				parameters = append(parameters, ast.VariableDeclaration{
					Name: parameter.Literal, Type: parameterType, Value: nil,
				})
			}
		}

		_, colonErr := p.Consume(lexer.COLON, "Expected ':' before function return type.")

		if colonErr != nil {
			return nil, colonErr
		}

		returnType, returnTypeErr := p.parseTypeDeclaration()

		if returnTypeErr != nil {
			return nil, returnTypeErr
		}

		_, thickArrowErr := p.Consume(lexer.THICK_ARROW, "Expected '=>' after parameter defintion.")

		if thickArrowErr != nil {
			return nil, thickArrowErr
		}

		body, statErr := p.parseStatement()

		if statErr != nil {
			return nil, statErr
		}

		return &ast.FunctionExpression{Parameters: parameters, Body: body, ReturnType: returnType, Line: line}, nil
	}

	return p.parseCall()
}

func (p *Parser) parseUnary() (ast.Expression, error) {
	if p.MatchToken(lexer.BANG, lexer.SUB) {
		op := p.PreviousToken()
		expr, err := p.parsePrimary()

		if err != nil {
			return nil, err
		}

		return &ast.Unary{Value: expr, Operator: op.Type, Line: op.Line}, nil
	}

	return p.parseFunction()
}

func (p *Parser) parseFactor() (ast.Expression, error) {
	expr, err := p.parseUnary()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(lexer.MULT, lexer.DIV) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseUnary()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type, Line: op.Line}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (ast.Expression, error) {
	expr, err := p.parseFactor()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(lexer.ADD, lexer.SUB) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseFactor()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type, Line: op.Line}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (ast.Expression, error) {
	expr, err := p.parseTerm()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(lexer.GT, lexer.GT_EQ, lexer.LT, lexer.LT_EQ) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseTerm()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type, Line: op.Line}
	}

	return expr, nil
}

func (p *Parser) parseEquality() (ast.Expression, error) {
	expr, err := p.parseComparison()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(lexer.EQUALITY, lexer.NOT_EQUAL) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseComparison()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type, Line: op.Line}
	}

	return expr, nil
}

func (p *Parser) parseLogicalAnd() (ast.Expression, error) {
	expr, eqErr := p.parseEquality()

	if eqErr != nil {
		return nil, eqErr
	}

	for p.MatchToken(lexer.AND) {
		op := p.PreviousToken()
		right, rightErr := p.parseEquality()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Logical{Left: expr, Right: right, Operator: op.Type, Line: op.Line}
	}

	return expr, nil
}

func (p *Parser) parseLogicalOr() (ast.Expression, error) {
	expr, andErr := p.parseLogicalAnd()

	if andErr != nil {
		return nil, andErr
	}

	for p.MatchToken(lexer.OR) {
		op := p.PreviousToken()
		right, rightErr := p.parseLogicalAnd()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Logical{Left: expr, Right: right, Operator: op.Type, Line: op.Line}
	}

	return expr, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseLogicalOr()
}
