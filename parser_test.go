package main

import (
	"testing"
)

func TestVariableDeclarations(t *testing.T) {
	lex := ScanTokens(`
		let x : int = 100
		let y : string? = "hello"
		let z : float?
		let h : bool
	`)

	statements := Parse(lex.Tokens)

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}

	x := statements[0].(*VariableDeclaration)
	xType := x.Type.(*VariableType)
	xValue := x.Value.(*Literal)

	if x.Name != "x" {
		t.Errorf("Expected \"%s\", got \"%s\".", "x", x.Name)
	}

	if xType.Base != "int" {
		t.Errorf("Expected \"%s\", got \"%s\".", "int", xType.Base)
	}

	if xValue.Value != "100" {
		t.Errorf("Expected \"%s\", got \"%s\".", "100", xValue.Value)
	}

	y := statements[1].(*VariableDeclaration)
	yType := y.Type.(*VariableType)
	yValue := y.Value.(*Literal)

	if !yType.Optional {
		t.Errorf("Expected y to be optional.")
	}

	if yValue.Value != "hello" {
		t.Errorf("Expected value of y to be \"hello\", got \"%s\"", yValue.Value)
	}
}
