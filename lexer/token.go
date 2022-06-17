package lexer

type TokenType int64

const (
	END_OF_FILE TokenType = iota
	IDENTIFIER	
	NUMBER
	BOOL
	EQUAL
	ADD
	SUB
	MULT
	DIV
	GT
	GT_EQ
	LT
	LT_EQ
	BANG
	EQUALITY
	NOT_EQUAL
	COMMA
	PERIOD
	COLON
	L_PAREN
	R_PAREN
	L_BRACE
	R_BRACE
	L_BRACKET
	R_BRACKET
	QUOTE
	STRING
	LET
	WHILE
	FOR
	IF
	ELSE
	ARROW
	THICK_ARROW
	TRUE
	FALSE
	NULL
)

type Token struct {
	Type TokenType
	Literal string
}