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
	var exprErr error

	if p.MatchToken(lexer.EQUAL) {
		value, exprErr = p.parseExpression()

		if exprErr != nil {
			return nil, exprErr
		}
	}

	return &ast.VariableDeclaration{
		Name:         name.Literal,
		Type:         variableType,
		Value:        value,
		NodeMetadata: ast.CreateMetadata(name.Line),
	}, nil
}

func (p *Parser) parseRecord() (ast.Type, error) {
	_, leftBraceErr := p.Consume(lexer.L_BRACE, "Expected '{' when declaring record.")

	if leftBraceErr != nil {
		return nil, leftBraceErr
	}

	if p.MatchToken(lexer.R_BRACE) {
		return &ast.RecordType{Fields: make(map[string]ast.Type)}, nil
	}

	fields := make(map[string]ast.Type)

	for hasComma := true; hasComma; hasComma = p.MatchToken(lexer.COMMA) {
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

		fields[variableName.Literal] = variableType
	}

	_, rightBraceErr := p.Consume(lexer.R_BRACE, "Expected '}' to close record declaration.")

	if rightBraceErr != nil {
		return nil, rightBraceErr
	}

	return &ast.RecordType{Fields: fields}, nil
}

func (p *Parser) parseRecordDeclaration() (ast.Statement, error) {
	identifier, identifierErr := p.Consume(
		lexer.IDENTIFIER,
		"Expected name after type definition.",
	)

	if identifierErr != nil {
		return nil, identifierErr
	}

	var inheritsFrom string

	if p.MatchToken(lexer.L_PAREN) {
		inherits, inheritsErr := p.Consume(lexer.IDENTIFIER, "Expected base type to inherit from.")

		if inheritsErr != nil {
			return nil, inheritsErr
		}

		_, parenErr := p.Consume(lexer.R_PAREN, "Expected closing parenthesis.")

		if parenErr != nil {
			return nil, parenErr
		}

		inheritsFrom = inherits.Literal
	}

	record, recordErr := p.parseRecord()

	if recordErr != nil {
		return nil, recordErr
	}

	recordValue, isRecord := record.(*ast.RecordType)

	if !isRecord {
		return nil, CreateParseError(identifier.Line, "Expected record declaration.")
	}

	return &ast.RecordDeclaration{
		Name:         identifier.Literal,
		Record:       *recordValue,
		Inherits:     inheritsFrom,
		NodeMetadata: ast.CreateMetadata(identifier.Line),
	}, nil
}

func (p *Parser) parseDeclaration() (ast.Statement, error) {
	if p.MatchToken(lexer.LET) {
		return p.parseVariableDeclaration()
	}

	return p.parseStatement()
}

func (p *Parser) parseBlockStatement() (ast.Statement, error) {
	statements := make([]ast.Statement, 0)

	openParen := p.PreviousToken()

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

	return &ast.BlockStatement{
		Statements:   statements,
		NodeMetadata: ast.CreateMetadata(openParen.Line),
	}, nil
}

func (p *Parser) parseIfStatement() (ast.Statement, error) {
	line := p.PreviousToken().Line

	_, openParenErr := p.Consume(lexer.L_PAREN, "Expected open parenthesis.")

	if openParenErr != nil {
		return nil, openParenErr
	}

	condition, conditionErr := p.parseExpression()

	if conditionErr != nil {
		return nil, conditionErr
	}

	_, closeParenErr := p.Consume(lexer.R_PAREN, "Expected closing parenthesis.")

	if closeParenErr != nil {
		return nil, closeParenErr
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

	return &ast.IfStatement{
		Condition:    condition,
		Body:         ifBranch,
		ElseBody:     elseBranch,
		NodeMetadata: ast.CreateMetadata(line),
	}, nil
}

func (p *Parser) parseWhileStatement() (ast.Statement, error) {
	line := p.PreviousToken().Line

	_, openParenErr := p.Consume(lexer.L_PAREN, "Expected open parenthesis.")

	if openParenErr != nil {
		return nil, openParenErr
	}

	condition, conditionErr := p.parseExpression()

	if conditionErr != nil {
		return nil, conditionErr
	}

	_, closeParenErr := p.Consume(lexer.R_PAREN, "Expected closing parenthesis.")

	if closeParenErr != nil {
		return nil, closeParenErr
	}

	body, bodyErr := p.parseStatement()

	if bodyErr != nil {
		return nil, bodyErr
	}

	return &ast.WhileStatement{
		Condition:    condition,
		Body:         body,
		NodeMetadata: ast.CreateMetadata(line),
	}, nil
}

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	line := p.CurrentToken().Line
	expression, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{Value: expression, NodeMetadata: ast.CreateMetadata(line)}, err
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	line := p.PreviousToken().Line
	value, valueErr := p.parseExpression()

	if valueErr != nil {
		return nil, valueErr
	}

	return &ast.ReturnStatement{Value: value, NodeMetadata: ast.CreateMetadata(line)}, nil
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
	} else if p.MatchToken(lexer.TYPE) {
		return p.parseRecordDeclaration()
	} else if p.MatchToken(lexer.RETURN) {
		return p.parseReturnStatement()
	}

	return p.parseExpressionStatement()
}
