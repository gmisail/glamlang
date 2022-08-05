package tests

import (
	"testing"

	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"github.com/stretchr/testify/assert"
)

func TestLogical(t *testing.T) {
	lex := lexer.ScanTokens(`
		((100 - 100) == 0) and false
		i == 0 or i != 100
		true and ((100 - 100) == 0)
		true or false
	`)

	ok, statements := parser.Parse(lex.Tokens)

	assert.True(t, ok)

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}
}

func TestEquality(t *testing.T) {
	lex := lexer.ScanTokens(`
		((100 - 100) == 0) == false
		i == (i != 100)
		true != ((100 - 100) == 0)
		true != false
	`)

	ok, statements := parser.Parse(lex.Tokens)

	assert.True(t, ok)

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}
}

func TestComparison(t *testing.T) {
	lex := lexer.ScanTokens(`
		((100 - 100) == 0) > false
		i < (i != 100)
		true <= ((100 - 100) == 0)
		true >= false
	`)

	ok, statements := parser.Parse(lex.Tokens)

	assert.True(t, ok)

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}
}
