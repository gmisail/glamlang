package parser

import (
	"fmt"

	"github.com/gmisail/glamlang/io"
	"github.com/gmisail/glamlang/lexer"
	"github.com/gmisail/glamlang/util"
)

type ParseError struct {
	source  *io.SourceFile
	token   *lexer.Token
	message string
}

func (p *ParseError) Error() string {
	if p.token.Line == 0 {
		return fmt.Sprintf("EOF: %s", p.message)
	}

	hint := util.Hint(p.source, p.token, p.message)

	return fmt.Sprintf("Error at line %d:\n%s", p.token.Line, hint)
}

func CreateParseError(source *io.SourceFile, token *lexer.Token, message string) *ParseError {
	return &ParseError{source, token, message}
}
