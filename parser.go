package main

type Parser struct {
	current int
	Tokens  []*Token
}

func (p *Parser) AdvanceToken() {
	p.current += 1
}

func (p *Parser) CurrentToken() *Token {
	if p.current >= len(p.Tokens) {
		return nil
	}

	return p.Tokens[p.current]
}

func (p *Parser) PreviousToken() *Token {
	if p.current-1 < 0 {
		return nil
	}

	return p.Tokens[p.current-1]
}

func (p *Parser) PeekToken() *Token {
	if p.current+1 < len(p.Tokens) {
		return p.Tokens[p.current+1]
	}

	return nil
}

func (p *Parser) MatchToken(types ...TokenType) bool {
	next := p.PeekToken()

	for _, tokenType := range types {
		if next.Type == tokenType {
			p.AdvanceToken()
			return true
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
		// parse parenthesis expr

		expr := parseExpression()

		return &Group{Value: expr}, nil
	}

	return nil, nil
}

func (p *Parser) parseUnary() {

}

func (p *Parser) parseFactor() {

}

func (p *Parser) parseTerm() {

}

func (p *Parser) parseComparison() (Expression, error) {

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

func Parse(tokens []*Token) *Parser {
	parser := &Parser{current: -1, Tokens: tokens}

	return parser
}
