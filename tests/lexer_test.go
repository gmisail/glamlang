package tests

import (
	"testing"
)

func TestNumberOfTokens(t *testing.T) {
	lex := ScanTokens("[](){}.,    +-*/")

	numTokens := len(lex.Tokens)

	if numTokens != 12 {
		t.Errorf("Found %d tokens, expected 12.", numTokens)
	}
}

func TestKeywords(t *testing.T) {
	lex := ScanTokens("hello let while for if else true false")
	expected := []TokenType{
		IDENTIFIER, LET, WHILE, FOR, IF, ELSE, TRUE, FALSE,
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
	lex := ScanTokens("100 123456 12.14 5000.00")
	expected := []TokenType{
		INT, INT, FLOAT, FLOAT,
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

func TestString(t *testing.T) {
	lex := ScanTokens("hello \"from way up here\"")
	expected := []TokenType{
		IDENTIFIER, STRING,
	}

	numTokens := len(lex.Tokens)

	if numTokens != 2 {
		t.Errorf("Found %d tokens, expected 2.", numTokens)
	}

	for i, tok := range lex.Tokens {
		if tok.Type != expected[i] {
			t.Errorf("Token '%s' has type %d, expecting %d.", tok.Literal, tok.Type, expected[i])
		}
	}
}

func TestConditionalTokens(t *testing.T) {
	lex := ScanTokens("=> == != ->")
	expected := []TokenType{
		THICK_ARROW, EQUALITY, NOT_EQUAL, ARROW,
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
