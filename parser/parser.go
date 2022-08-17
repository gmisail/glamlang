package parser

import (
	"github.com/fatih/color"
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

type Parser struct {
	current int
	Tokens  []lexer.Token
}

func (p *Parser) AdvanceToken() {
	p.current += 1
}

func (p *Parser) CurrentToken() *lexer.Token {
	if p.current >= len(p.Tokens) {
		return nil
	}

	return &(p.Tokens[p.current])
}

func (p *Parser) PreviousToken() *lexer.Token {
	if p.current-1 < 0 {
		return nil
	}

	return &(p.Tokens[p.current-1])
}

func (p *Parser) PeekToken() *lexer.Token {
	if p.current+1 < len(p.Tokens) {
		return &(p.Tokens[p.current+1])
	}

	return nil
}

func (p *Parser) MatchToken(types ...lexer.TokenType) bool {
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

func (p *Parser) Check(types ...lexer.TokenType) bool {
	next := p.PeekToken()
	isNext := false

	for _, token := range types {
		if next.Type == token {
			isNext = true
			break
		}
	}

	return next != nil && isNext
}

func (p *Parser) Consume(tokenType lexer.TokenType, message string) (*lexer.Token, error) {
	if !p.MatchToken(tokenType) {
		currentToken := p.CurrentToken()

		if currentToken != nil {
			return nil, &ParseError{message: message, line: currentToken.Line}
		}

		return nil, &ParseError{message: message, line: 0}
	}

	return p.PreviousToken(), nil
}

func (p *Parser) Calibrate() {
	p.AdvanceToken()

	for p.PeekToken() != nil && !p.Check(lexer.WHILE, lexer.IF, lexer.WHILE, lexer.LET, lexer.TYPE, lexer.L_BRACE) {
		p.AdvanceToken()
	}
}

func Parse(tokens []lexer.Token) (bool, []ast.Statement) {
	parser := &Parser{current: 0, Tokens: tokens}
	statements := make([]ast.Statement, 0)
	isValid := true

	for {
		statement, err := parser.parseDeclaration()

		// no more statements
		if statement == nil && err == nil {
			break
		} else if statement == nil && err != nil {
			isValid = false

			color.Red(err.Error())
			parser.Calibrate()

			continue
		}

		statements = append(statements, statement)
	}

	return isValid, statements
}
