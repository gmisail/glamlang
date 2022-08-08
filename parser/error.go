package parser

import "fmt"

type ParseError struct {
	line    int
	message string
}

func (p *ParseError) Error() string {
	if p.line == 0 {
		return fmt.Sprintf("EOF: %s", p.message)
	}

	return fmt.Sprintf("line %d: %s", p.line, p.message)
}

func CreateParseError(line int, message string) *ParseError {
	return &ParseError{line, message}
}
