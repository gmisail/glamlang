package main

import "fmt"

type Parser struct {
	current int
	Tokens  []Token
}

type ParseError struct {
	line    int
	message string
}

func (p *ParseError) Error() string {
	if p.line == 0 {
		return fmt.Sprintf("EOF: %s", p.message)
	}

	return fmt.Sprintf("line %d: %s", p.line, p.message)
}

func (p *Parser) AdvanceToken() {
	p.current += 1
}

func (p *Parser) CurrentToken() *Token {
	if p.current >= len(p.Tokens) {
		return nil
	}

	return &(p.Tokens[p.current])
}

func (p *Parser) PreviousToken() *Token {
	if p.current-1 < 0 {
		return nil
	}

	return &(p.Tokens[p.current-1])
}

func (p *Parser) PeekToken() *Token {
	if p.current+1 < len(p.Tokens) {
		return &(p.Tokens[p.current+1])
	}

	return nil
}

func (p *Parser) MatchToken(types ...TokenType) bool {
	next := p.CurrentToken()

	if next != nil {
		for _, tokenType := range types {
			if next.Type == tokenType {
				p.AdvanceToken()
				return true
			}
		}
	}

	return false
}

func (p *Parser) Check(tokenType TokenType) bool {
	next := p.PeekToken()

	return next != nil && next.Type == tokenType
}

func (p *Parser) Consume(tokenType TokenType, message string) (*Token, error) {
	if !p.MatchToken(tokenType) {
		currentToken := p.CurrentToken()

		if currentToken != nil {
			return nil, &ParseError{message: message, line: currentToken.Line}
		}

		return nil, &ParseError{message: message, line: 0}
	}

	return p.PreviousToken(), nil
}

func (p *Parser) parsePrimary() (Expression, error) {
	if p.MatchToken(FALSE) {
		return &Literal{Value: false}, nil
	} else if p.MatchToken(TRUE) {
		return &Literal{Value: true}, nil
	} else if p.MatchToken(NULL) {
		return &Literal{Value: nil}, nil
	} else if p.MatchToken(STRING, INT, FLOAT) {
		return &Literal{Value: p.PreviousToken().Literal}, nil
	} else if p.MatchToken(IDENTIFIER) {
		return &VariableExpression{Value: p.PreviousToken().Literal}, nil
	} else if p.MatchToken(L_PAREN) {
		expr, _ := p.parseExpression()
		_, err := p.Consume(R_PAREN, "Expected closing parenthesis for group expression.")

		if err != nil {
			return nil, err
		}

		return &Group{Value: expr}, nil
	}

	return nil, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	if p.MatchToken(BANG, SUB) {
		op := p.PreviousToken()
		expr, err := p.parsePrimary()

		if err != nil {
			return nil, err
		}

		return &Unary{Value: expr, Operator: op.Type}, nil
	}

	return p.parsePrimary()
}

func (p *Parser) parseFactor() (Expression, error) {
	expr, err := p.parseUnary()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(MULT, DIV) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseUnary()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Expression, error) {
	expr, err := p.parseFactor()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(ADD, SUB) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseFactor()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (Expression, error) {
	expr, err := p.parseTerm()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(GT, GT_EQ, LT, LT_EQ) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseTerm()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseEquality() (Expression, error) {
	expr, err := p.parseComparison()

	if err != nil {
		return nil, err
	}

	for p.MatchToken(EQUALITY, NOT_EQUAL) {
		op := p.PreviousToken()
		rightExpr, rightErr := p.parseComparison()

		if rightErr != nil {
			return nil, rightErr
		}

		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseEquality()
}

func (p *Parser) parseTypeDeclaration() (TypeDefinition, error) {
	if p.MatchToken(L_PAREN) {
		arguments := make([]TypeDefinition, 0)

		for {
			argumentType, _ := p.parseTypeDeclaration()
			arguments = append(arguments, argumentType)

			// TODO: handle me

			if !p.MatchToken(COMMA) && p.MatchToken(R_PAREN) {
				break
			}
		}

		p.Consume(ARROW, "Expected '->' after argument type declaration.")

		returnType, _ := p.parseTypeDeclaration()

		return &FunctionType{ArgumentTypes: arguments, ReturnType: returnType}, nil
	}

	name, _ := p.Consume(IDENTIFIER, "Expected type name.")
	isOptional := p.MatchToken(QUESTION)

	return &VariableType{Base: name.Literal, SubType: nil, Optional: isOptional}, nil
}

func (p *Parser) parseVariableDeclaration() (Statement, error) {
	/*
		let <name> : <type> (= <expression>)?
	*/

	name, _ := p.Consume(IDENTIFIER, "Expected variable name.")

	p.Consume(COLON, "Expected ':' after variable name in declaration.")

	variableType, _ := p.parseTypeDeclaration()

	var value Expression = nil
	if p.MatchToken(EQUAL) {
		value, _ = p.parseExpression()
	}

	return &VariableDeclaration{Name: name.Literal, Type: variableType, Value: value}, nil
}

func (p *Parser) parseDeclaration() (Statement, error) {
	if p.MatchToken(LET) {
		return p.parseVariableDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseExpressionStatement() (Statement, error) {
	expression, err := p.parseExpression()

	// TODO: handle me
	if err != nil {
		fmt.Println(err)
	}

	return &ExpressionStatement{Value: expression}, err
}

func (p *Parser) parseStatement() (Statement, error) {
	if p.CurrentToken() == nil {
		return nil, nil
	}

	// if p.MatchToken() {

	// }

	return p.parseExpressionStatement()
}

func Parse(tokens []Token) []Statement {
	parser := &Parser{current: 0, Tokens: tokens}

	statements := make([]Statement, 0)

	for {
		statement, err := parser.parseDeclaration()

		// no more statements
		if statement == nil && err == nil {
			break
		} else if statement == nil && err != nil {
			// handle error

			// skip tokens until we find another statement to parse
			// if not, break
		}

		statements = append(statements, statement)
	}

	return statements
}
