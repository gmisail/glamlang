package parser

import (
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (p *Parser) parseVariableDeclaration() (ast.Statement, error) {
	/*
		let <name> : <type> (= <expression>)?
	*/

	name, nameErr := p.Consume(lexer.IDENTIFIER, "Expected variable name.")

	if nameErr != nil {
		return nil, nameErr
	}

	_, colonErr := p.Consume(lexer.COLON, "Expected ':' after variable name in declaration.")

	if colonErr != nil {
		return nil, colonErr
	}

	variableType, typeErr := p.parseTypeDeclaration()

	if typeErr != nil {
		return nil, typeErr
	}

	var value ast.Expression = nil
	var exprErr error = nil

	if p.MatchToken(lexer.EQUAL) {
		value, exprErr = p.parseExpression()

		if exprErr != nil {
			return nil, exprErr
		}
	}

	return &ast.VariableDeclaration{Name: name.Literal, Type: variableType, Value: value}, nil
}

func (p *Parser) parseStructDeclaration() (ast.Statement, error) {
	identifier, identifierErr := p.Consume(lexer.IDENTIFIER, "Expected name after struct definition.")

	if identifierErr != nil {
		return nil, identifierErr
	}

	variables := make([]ast.VariableDeclaration, 0)

	_, leftBraceErr := p.Consume(lexer.L_BRACE, "Expected '{' when declaring struct.")

	if leftBraceErr != nil {
		return nil, leftBraceErr
	}

	for {
		variableName, variableErr := p.Consume(lexer.IDENTIFIER, "Expected variable name")

		if variableErr != nil {
			return nil, variableErr
		}

		_, colonErr := p.Consume(lexer.COLON, "Expected ':' after variable name.")

		if colonErr != nil {
			return nil, colonErr
		}

		variableType, variableTypeErr := p.parseTypeDeclaration()

		if variableTypeErr != nil {
			return nil, variableTypeErr
		}

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
		statement, statementErr := p.parseDeclaration()

		if statement == nil {
			break
		}

		if statementErr != nil {
			return nil, statementErr
		}

		statements = append(statements, statement)
	}

	_, rightBraceErr := p.Consume(lexer.R_BRACE, "Expected '}' after block")

	if rightBraceErr != nil {
		return nil, rightBraceErr
	}

	return &ast.BlockStatement{Statements: statements}, nil
}

func (p *Parser) parseIfStatement() (ast.Statement, error) {
	condition, conditionErr := p.parseExpression()

	if conditionErr != nil {
		return nil, conditionErr
	}

	ifBranch, ifBranchErr := p.parseStatement()

	if ifBranchErr != nil {
		return nil, ifBranchErr
	}

	var elseBranch ast.Statement = nil

	if p.MatchToken(lexer.ELSE) {
		elseBranchStatement, elseBranchStatementErr := p.parseStatement()

		if elseBranchStatementErr != nil {
			return nil, elseBranchStatementErr
		}

		elseBranch = elseBranchStatement
	}

	return &ast.IfStatement{Condition: condition, Body: ifBranch, ElseBody: elseBranch}, nil
}

func (p *Parser) parseWhileStatement() (ast.Statement, error) {
	condition, conditionErr := p.parseExpression()

	if conditionErr != nil {
		return nil, conditionErr
	}

	body, bodyErr := p.parseStatement()

	if bodyErr != nil {
		return nil, bodyErr
	}

	return &ast.WhileStatement{Condition: condition, Body: body}, nil
}

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	expression, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{Value: expression}, err
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	value, valueErr := p.parseExpression()

	if valueErr != nil {
		return nil, valueErr
	}

	return &ast.ReturnStatement{Value: value}, nil
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
	} else if p.MatchToken(lexer.RETURN) {
		return p.parseReturnStatement()
	}

	return p.parseExpressionStatement()
}
