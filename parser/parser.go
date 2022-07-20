package parser

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

type Parser struct {
	current int
	Tokens  []lexer.Token
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

func (p *Parser) Check(tokenType lexer.TokenType) bool {
	next := p.PeekToken()

	return next != nil && next.Type == tokenType
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

func Parse(tokens []lexer.Token) []ast.Statement {
	parser := &Parser{current: 0, Tokens: tokens}

	statements := make([]ast.Statement, 0)

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
