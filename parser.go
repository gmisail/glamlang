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

func (p *ParseError) Error() {
	fmt.Errorf("line %d: %d\n", p.line, p.message)
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

func (p *Parser) Check(tokenType TokenType) bool {
	next := p.PeekToken()

	return next != nil && next.Type == tokenType
}

func (p *Parser) Consume(tokenType TokenType, message string) (bool, error) {
	if !p.MatchToken(tokenType) {
		//	currentToken := p.CurrentToken()
		return false, nil
	}

	return true, nil
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

func (p *Parser) parsePrimary() (Expression, error) {
	if p.MatchToken(FALSE) {
		return &Literal{Value: false}, nil
	} else if p.MatchToken(TRUE) {
		return &Literal{Value: true}, nil
	} else if p.MatchToken(NULL) {
		return &Literal{Value: nil}, nil
	} else if p.MatchToken(STRING, INT, FLOAT) {
		return &Literal{Value: p.PreviousToken().Literal}, nil
	} else if p.MatchToken(L_PAREN) {
		expr, _ := p.parseExpression()

		p.Consume(R_PAREN, "Expected closing parenthesis.")

		return &Group{Value: expr}, nil
	}

	return nil, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	if p.MatchToken(BANG, SUB) {
		op := p.PreviousToken()
		expr, _ := p.parsePrimary()
		return &Unary{Value: expr, Operator: op.Type}, nil
	}

	return p.parsePrimary()
}

func (p *Parser) parseFactor() (Expression, error) {
	expr, _ := p.parseUnary()

	for p.MatchToken(MULT, DIV) {
		op := p.PreviousToken()
		rightExpr, _ := p.parseUnary()

		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseTerm() (Expression, error) {
	expr, _ := p.parseFactor()

	for p.MatchToken(ADD, SUB) {
		op := p.PreviousToken()
		rightExpr, _ := p.parseFactor()

		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseComparison() (Expression, error) {
	expr, _ := p.parseTerm()

	for p.MatchToken(GT, GT_EQ, LT, LT_EQ) {
		op := p.PreviousToken()
		rightExpr, _ := p.parseTerm()
		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseEquality() (Expression, error) {
	expr, _ := p.parseComparison()

	for p.MatchToken(EQUALITY, NOT_EQUAL) {
		op := p.PreviousToken()
		rightExpr, _ := p.parseComparison()
		expr = &Binary{Left: expr, Right: rightExpr, Operator: op.Type}
	}

	return expr, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseEquality()
}

func Parse(tokens []Token) (Expression, error) {
	parser := &Parser{current: 0, Tokens: tokens}

	return parser.parseExpression()
}
