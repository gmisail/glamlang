package lexer_test

import (
	"testing"

	"github.com/gmisail/glamlang/lexer"
)

func TestNumberOfTokens(t *testing.T) {
	lex := lexer.ScanTokens("[](){}.,    +-*/")

	numTokens := lex.Tokens.Len()

	if numTokens != 12 {
		t.Errorf("Found %d tokens, expected 13.", numTokens)
	}
}
