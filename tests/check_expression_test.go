package tests

import (
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
	"github.com/gmisail/glamlang/typechecker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnary(t *testing.T) {
	lex := lexer.ScanTokens(`
		let x : int = -100
		let y : int = !500
		let z : bool = -false
		let h : bool = !false
	`)

	ok, statements := parser.Parse(lex.Tokens)
	states := []bool{true, false, false, true}

	assert.True(t, ok)

	tc := typechecker.CreateTypeChecker()

	for i, statement := range statements {
		switch s := statement.(type) {
		case *ast.VariableDeclaration:
			_, err := tc.CheckExpression(s.Value)

			assert.Equal(t, states[i], err == nil)
		default:
			t.Errorf("Expected statement to be variable declaration.")
		}
	}

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}
}

func TestBinary(t *testing.T) {
	lex := lexer.ScanTokens(`
		let x : int = -100 + 50
		let y : int = 500 * true
		let z : bool = false + true 
		let h : bool = 100 - 5 
	`)

	ok, statements := parser.Parse(lex.Tokens)
	states := []bool{true, false, false, true}

	assert.True(t, ok)

	tc := typechecker.CreateTypeChecker()

	for i, statement := range statements {
		switch s := statement.(type) {
		case *ast.VariableDeclaration:
			_, err := tc.CheckExpression(s.Value)

			assert.Equal(t, states[i], err == nil)
		default:
			t.Errorf("Expected statement to be variable declaration.")
		}
	}

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}
}
