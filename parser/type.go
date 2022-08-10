package parser

import (
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (p *Parser) parseTypeDeclaration() (ast.Type, error) {
	if p.MatchToken(lexer.L_PAREN) {
		arguments := make([]ast.Type, 0)

		for {
			if !p.MatchToken(lexer.COMMA) && p.MatchToken(lexer.R_PAREN) {
				break
			}

			argumentType, argumentErr := p.parseTypeDeclaration()

			if argumentErr != nil {
				return nil, argumentErr
			}

			arguments = append(arguments, argumentType)
		}

		p.Consume(lexer.ARROW, "Expected '->' after argument type declaration.")

		returnType, _ := p.parseTypeDeclaration()

		return &ast.FunctionType{Parameters: arguments, ReturnType: returnType}, nil
	}

	name, nameErr := p.Consume(lexer.IDENTIFIER, "Expected type name.")

	if nameErr != nil {
		return nil, nameErr
	}

	isOptional := p.MatchToken(lexer.QUESTION)

	return &ast.VariableType{Base: name.Literal, SubType: nil, Optional: isOptional}, nil
}
