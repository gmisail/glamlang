package lexer_test

import (
	"testing"

	"github.com/gmisail/glamlang/lexer"
)

func TestNumberOfTokens(t *testing.T) {
	lex := lexer.ScanTokens("[](){}.,    +-*/")

	numTokens := len(lex.Tokens)

	if numTokens != 12 {
		t.Errorf("Found %d tokens, expected 12.", numTokens)
	}
}

func TestKeywords(t *testing.T) {
	lex := lexer.ScanTokens("hello let while for if else true false")
	expected := []lexer.TokenType{
		lexer.IDENTIFIER, lexer.LET, lexer.WHILE, lexer.FOR, lexer.IF, lexer.ELSE, lexer.TRUE, lexer.FALSE,
	}

	numTokens := len(lex.Tokens)

	if numTokens != len(expected) {
		t.Errorf("Found %d tokens, expected %d.", numTokens, len(expected))
	}

	for i, tok := range lex.Tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token '%s' has type %d, expecting %d.", tok.Literal, tok.Type, expected[i])
		}
	}
}

func TestNumbers(t *testing.T) {
	lex := lexer.ScanTokens("100 123456 12.14 5000.00")
	expected := []lexer.TokenType{
		lexer.INT, lexer.INT, lexer.FLOAT, lexer.FLOAT,
	}

	numTokens := len(lex.Tokens)

	if numTokens != 4 {
		t.Errorf("Found %d tokens, expected 4.", numTokens)
	}

	for i, tok := range lex.Tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token '%s' has type %d, expecting %d.", tok.Literal, tok.Type, expected[i])
		}
	}
}
