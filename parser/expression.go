package parser

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (p *Parser) parsePrimary() (ast.Expression, error) {
	if p.MatchToken(lexer.FALSE) {
		return &ast.Literal{Value: false, Type: lexer.BOOL}, nil
	} else if p.MatchToken(lexer.TRUE) {
		return &ast.Literal{Value: true, Type: lexer.BOOL}, nil
	} else if p.MatchToken(lexer.NULL) {
		return &ast.Literal{Value: nil, Type: lexer.NULL}, nil
	} else if p.MatchToken(lexer.STRING, lexer.INT, lexer.FLOAT) {
		previousToken := p.PreviousToken()

		return &ast.Literal{Value: previousToken.Literal, Type: previousToken.Type}, nil
	} else if p.MatchToken(lexer.IDENTIFIER) {
		return &ast.VariableExpression{Value: p.PreviousToken().Literal}, nil
	} else if p.MatchToken(lexer.L_PAREN) {
		expr, _ := p.parseExpression()

		_, err := p.Consume(lexer.R_PAREN, "Expected closing parenthesis for group expression.")

		if err != nil {
			return nil, err
		}

		return &ast.Group{Value: expr}, nil
	}

	return nil, &ParseError{message: fmt.Sprintf("Unexpected token: %s", lexer.TokenTypeToString(p.CurrentToken().Type)), line: p.CurrentToken().Line}
}

func (p *Parser) finishParseCall(callee ast.Expression) (ast.Expression, error) {
	arguments := make([]ast.Expression, 0)

	for p.CurrentToken().Type != lexer.R_PAREN {
		argument, _ := p.parseExpression()
		arguments = append(arguments, argument)

		_, commaErr := p.Consume(lexer.COMMA, "Expected comma after arguments in function call.")

		if commaErr != nil {
			return nil, commaErr
		}
	}

	_, rightParenErr := p.Consume(lexer.R_PAREN, "Expected ')' after function call arguments.")

	if rightParenErr != nil {
		return nil, rightParenErr
	}

	return &ast.FunctionCall{Callee: callee, Arguments: arguments}, nil
}

func (p *Parser) parseCall() (ast.Expression, error) {
	expr, primaryErr := p.parsePrimary()

	if primaryErr != nil {
		return nil, primaryErr
	}

	var callErr error = nil

	for {
		if p.MatchToken(lexer.L_PAREN) {
			expr, callErr = p.finishParseCall(expr)

			if callErr != nil {
				return nil, callErr
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) parseFunction() (ast.Expression, error) {
	if p.MatchToken(lexer.FUNCTION) {
		parameters := make([]string, 0)

		_, leftParenErr := p.Consume(lexer.L_PAREN, "Expected '('")

		if leftParenErr != nil {
			return nil, leftParenErr
		}

		// if there's a right parenthesis, that means the function doesn't have any parameters.
		if !p.MatchToken(lexer.R_PAREN) {
			for {
				parameter, parameterErr := p.Consume(lexer.IDENTIFIER, "Expected parameter name.")

				if parameterErr != nil {
					return nil, parameterErr
				}

				parameters = append(parameters, parameter.Literal)

				// no more parameters :(
				if !p.MatchToken(lexer.COMMA) && p.MatchToken(lexer.R_PAREN) {
					break
				}
			}
		}

		_, thickArrowErr := p.Consume(lexer.THICK_ARROW, "Expected '=>' after parameter defintion.")

		if thickArrowErr != nil {
			return nil, thickArrowErr
		}

		body, statErr := p.parseStatement()

		if statErr != nil {
			return nil, statErr
		}

		return &ast.FunctionExpression{Parameters: parameters, Body: body}, nil
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

		return &ast.Unary{Value: expr, Operator: op.Type}, nil
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

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type}
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

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type}
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

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type}
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

		expr = &ast.Binary{Left: expr, Right: rightExpr, Operator: op.Type}
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

		expr = &ast.Logical{Left: expr, Right: right, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseLogicalOr() (ast.Expression, error) {
	expr, andErr := p.parseLogicalAnd()

	if andErr != nil {
		return nil, andErr
	}

	for p.MatchToken(lexer.AND) {
		op := p.PreviousToken()
		right, rightErr := p.parseLogicalAnd()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &ast.Logical{Left: expr, Right: right, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseExpression() (ast.Expression, error) {
	return p.parseLogicalOr()
}
