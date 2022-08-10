package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/parser"
)

func TestVariableDeclarations(t *testing.T) {
	lex := lexer.ScanTokens(`
		let x : int = 100
		let y : string? = "hello"
		let z : float?
		let h : bool
	`)

	_, statements := parser.Parse(lex.Tokens)

	if len(statements) != 4 {
		t.Errorf("Expected %d statements, got %d.", 4, len(statements))
	}

	x := statements[0].(*ast.VariableDeclaration)
	xType := x.Type.(*ast.VariableType)
	xValue := x.Value.(*ast.Literal)

	if x.Name != "x" {
		t.Errorf("Expected \"%s\", got \"%s\".", "x", x.Name)
	}

	if xType.Base != "int" {
		t.Errorf("Expected \"%s\", got \"%s\".", "int", xType.Base)
	}

	if xValue.Value != "100" {
		t.Errorf("Expected \"%s\", got \"%s\".", "100", xValue.Value)
	}

	y := statements[1].(*ast.VariableDeclaration)
	yType := y.Type.(*ast.VariableType)
	yValue := y.Value.(*ast.Literal)

	if !yType.Optional {
		t.Errorf("Expected y to be optional.")
	}

	if yValue.Value != "hello" {
		t.Errorf("Expected value of y to be \"hello\", got \"%s\"", yValue.Value)
	}
}

func TestFunctionType(t *testing.T) {
	lex := lexer.ScanTokens(`
		let square : (int, int) -> int
	`)

	_, statements := parser.Parse(lex.Tokens)

	if len(statements) != 1 {
		t.Errorf("Expected %d statements, got %d.", 1, len(statements))
	}

	variable := statements[0].(*ast.VariableDeclaration)
	functionType := variable.Type.(*ast.FunctionType)

	if len(functionType.Parameters) != 2 {
		t.Errorf("Expected %d statements, got %d.", 2, len(functionType.Parameters))
	}
}

func TestStructDeclaration(t *testing.T) {
	lex := lexer.ScanTokens(`
		struct Account {
			balance: int
			credit_limit: int
		}
	`)

	_, statements := parser.Parse(lex.Tokens)

	if len(statements) != 1 {
		t.Errorf("Expected %d statements, got %d.", 1, len(statements))
	}

	structDec := statements[0].(*ast.StructDeclaration)

	if len(structDec.Variables) != 2 {
		t.Errorf("Expected %d variables, got %d.", 2, len(structDec.Variables))
	}

	assert.Equal(t, "Account", structDec.Name)
	assert.Equal(t, "balance", structDec.Variables[0].Name)
	assert.Equal(t, "credit_limit", structDec.Variables[1].Name)
}

func TestBlockStatement(t *testing.T) {
	lex := lexer.ScanTokens(`
		{
			let x : int = 100
			let y : int = 5
		}
	`)

	_, statements := parser.Parse(lex.Tokens)

	assert.IsType(t, &ast.BlockStatement{}, statements[0])

	body := statements[0].(*ast.BlockStatement)

	assert.Len(t, body.Statements, 2)
}

func TestIfStatement(t *testing.T) {
	lex := lexer.ScanTokens(`
		if true {

		} else {

		}
	`)

	_, statements := parser.Parse(lex.Tokens)

	assert.IsType(t, &ast.IfStatement{}, statements[0])

	ifStatement := statements[0].(*ast.IfStatement)

	assert.NotNil(t, ifStatement.Body)
	assert.NotNil(t, ifStatement.ElseBody)
}

func TestWhileStatement(t *testing.T) {
	lex := lexer.ScanTokens(`
		while true {
			print()
		}
	`)

	_, statements := parser.Parse(lex.Tokens)

	assert.IsType(t, &ast.WhileStatement{}, statements[0])

	whileStat := statements[0].(*ast.WhileStatement)

	assert.NotNil(t, whileStat.Body)
	assert.IsType(t, &ast.BlockStatement{}, whileStat.Body)
}

func TestReturnStatement(t *testing.T) {
	lex := lexer.ScanTokens(`
		return false
	`)

	_, statements := parser.Parse(lex.Tokens)

	assert.IsType(t, &ast.ReturnStatement{}, statements[0])

	returnStat := statements[0].(*ast.ReturnStatement)

	assert.NotNil(t, returnStat.Value)
	assert.IsType(t, &ast.Literal{}, returnStat.Value)
}
