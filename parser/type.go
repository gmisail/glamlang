package parser

import (
	"github.com/gmisail/glamlang/ast"
	"github.com/gmisail/glamlang/lexer"
)

func (p *Parser) parseTypeDeclaration() (ast.TypeDefinition, error) {
	if p.MatchToken(lexer.L_PAREN) {
		arguments := make([]ast.TypeDefinition, 0)

		for {
			argumentType, _ := p.parseTypeDeclaration()
			arguments = append(arguments, argumentType)

			// TODO: handle me

			if !p.MatchToken(lexer.COMMA) && p.MatchToken(lexer.R_PAREN) {
				break
			}
		}

		p.Consume(lexer.ARROW, "Expected '->' after argument type declaration.")

		returnType, _ := p.parseTypeDeclaration()

		return &ast.FunctionType{ArgumentTypes: arguments, ReturnType: returnType}, nil
	}

	name, _ := p.Consume(lexer.IDENTIFIER, "Expected type name.")
	isOptional := p.MatchToken(lexer.QUESTION)

	return &ast.VariableType{Base: name.Literal, SubType: nil, Optional: isOptional}, nil
}
