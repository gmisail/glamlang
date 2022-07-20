package parser

import (
	"fmt"

	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (p *Parser) parseVariableDeclaration() (ast.Statement, error) {
	/*
		let <name> : <type> (= <expression>)?
	*/

	name, _ := p.Consume(lexer.IDENTIFIER, "Expected variable name.")

	p.Consume(lexer.COLON, "Expected ':' after variable name in declaration.")

	variableType, _ := p.parseTypeDeclaration()

	var value ast.Expression = nil
	if p.MatchToken(lexer.EQUAL) {
		value, _ = p.parseExpression()
	}

	return &ast.VariableDeclaration{Name: name.Literal, Type: variableType, Value: value}, nil
}

func (p *Parser) parseStructDeclaration() (ast.Statement, error) {
	identifier, _ := p.Consume(lexer.IDENTIFIER, "Expected name after struct definition.")
	variables := make([]ast.VariableDeclaration, 0)

	p.Consume(lexer.L_BRACE, "Expected '{' when declaring struct.")

	for {
		variableName, _ := p.Consume(lexer.IDENTIFIER, "Expected variable name")

		p.Consume(lexer.COLON, "Expected ':' after variable name.")

		variableType, _ := p.parseTypeDeclaration()
		variables = append(variables, ast.VariableDeclaration{Name: variableName.Literal, Type: variableType, Value: nil})

		if !p.MatchToken(lexer.COMMA) && p.MatchToken(lexer.R_BRACE) {
			break
		}
	}

	return &ast.StructDeclaration{Name: identifier.Literal, Variables: variables}, nil
}

func (p *Parser) parseDeclaration() (ast.Statement, error) {
	if p.MatchToken(lexer.LET) {
		return p.parseVariableDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseBlockStatement() (ast.Statement, error) {
	statements := make([]ast.Statement, 0)

	// TODO: make function for "check"
	for p.CurrentToken().Type != lexer.R_BRACE /*&& is not at end */ {
		statement, _ := p.parseDeclaration()

		if statement == nil {
			break
		}

		// TODO: handle error

		statements = append(statements, statement)
	}

	p.Consume(lexer.R_BRACE, "Expected '}' after block")

	return &ast.BlockStatement{Statements: statements}, nil
}

func (p *Parser) parseIfStatement() (ast.Statement, error) {
	condition, _ := p.parseExpression()
	ifBranch, _ := p.parseStatement()

	var elseBranch ast.Statement = nil

	if p.MatchToken(lexer.ELSE) {
		elseBranchStatement, _ := p.parseStatement()
		elseBranch = elseBranchStatement
	}

	return &ast.IfStatement{Condition: condition, Body: ifBranch, ElseBody: elseBranch}, nil
}

func (p *Parser) parseWhileStatement() (ast.Statement, error) {
	condition, _ := p.parseExpression()
	body, _ := p.parseStatement()

	return &ast.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	expression, err := p.parseExpression()

	// TODO: handle me
	if err != nil {
		fmt.Println(err)
	}

	return &ast.ExpressionStatement{Value: expression}, err
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	if p.CurrentToken() == nil {
		return nil, nil
	}

	if p.MatchToken(lexer.L_BRACE) {
		return p.parseBlockStatement()
	} else if p.MatchToken(lexer.IF) {
		return p.parseIfStatement()
	} else if p.MatchToken(lexer.WHILE) {
		return p.parseWhileStatement()
	} else if p.MatchToken(lexer.STRUCT) {
		return p.parseStructDeclaration()
	}

	return p.parseExpressionStatement()
}
