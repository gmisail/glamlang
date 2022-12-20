package lexer

var keywords = map[string]TokenType{
	"let":    LET,
	"while":  WHILE,
	"for":    FOR,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"fn":     FUNCTION,
	"and":    AND,
	"or":     OR,
	"type":   TYPE,
	"mod":    MODULE,
	"return": RETURN,
	"new":    NEW,
}

func LookupKeyword(literal string) TokenType {
	if tokenType, ok := keywords[literal]; ok {
		return tokenType
	}

	return IDENTIFIER
}
