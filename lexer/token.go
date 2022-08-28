package lexer

import "fmt"

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
	QUESTION
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
	RETURN
	FUNCTION
	TYPE
	MODULE
	AND
	OR
	NEW
	ARROW
	THICK_ARROW
	TRUE
	FALSE
	NULL
	FLOAT
	INT
)

func TokenTypeToString(token TokenType) string {
	switch token {
	case END_OF_FILE:
		return "END_OF_FILE"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NUMBER:
		return "NUMBER"
	case BOOL:
		return "BOOL"
	case EQUAL:
		return "EQUAL"
	case ADD:
		return "ADD"
	case SUB:
		return "SUB"
	case MULT:
		return "MULT"
	case DIV:
		return "DIV"
	case GT:
		return "GT"
	case GT_EQ:
		return "GT_EQ"
	case LT:
		return "LT"
	case LT_EQ:
		return "LT_EQ"
	case BANG:
		return "BANG"
	case EQUALITY:
		return "EQUALITY"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case COMMA:
		return "COMMA"
	case PERIOD:
		return "PERIOD"
	case COLON:
		return "COLON"
	case L_PAREN:
		return "L_PAREN"
	case R_PAREN:
		return "R_PAREN"
	case L_BRACE:
		return "L_BRACE"
	case R_BRACE:
		return "R_BRACE"
	case L_BRACKET:
		return "L_BRACKET"
	case R_BRACKET:
		return "R_BRACKET"
	case QUOTE:
		return "QUOTE"
	case STRING:
		return "STRING"
	case LET:
		return "LET"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case RETURN:
		return "RETURN"
	case TYPE:
		return "TYPE"
	case MODULE:
		return "MODULE"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NEW:
		return "NEW"
	case ARROW:
		return "ARROW"
	case THICK_ARROW:
		return "THICK_ARROW"
	case QUESTION:
		return "QUESTION"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case FLOAT:
		return "FLOAT"
	case INT:
		return "INT"
	}

	return "?"
}

func GetSymbol(token TokenType) string {
	switch token {
	case EQUAL:
		return "="
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MULT:
		return "*"
	case DIV:
		return "/"
	case GT:
		return ">"
	case GT_EQ:
		return ">="
	case LT:
		return "<"
	case LT_EQ:
		return "<="
	case BANG:
		return "!"
	case EQUALITY:
		return "=="
	case NOT_EQUAL:
		return "!="
	case COMMA:
		return ","
	case PERIOD:
		return "."
	case COLON:
		return ":"
	case L_PAREN:
		return "("
	case R_PAREN:
		return ")"
	case L_BRACE:
		return "{"
	case R_BRACE:
		return "}"
	case L_BRACKET:
		return "["
	case R_BRACKET:
		return "]"
	case QUOTE:
		return "\""
	case ARROW:
		return "->"
	case THICK_ARROW:
		return "=>"
	case QUESTION:
		return "?"
	case NULL:
		return "none"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case LET:
		return "let"
	case WHILE:
		return "while"
	case FOR:
		return "for"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case RETURN:
		return "return"
	case TYPE:
		return "type"
	case MODULE:
		return "module"
	case AND:
		return "and"
	case OR:
		return "or"
	}

	return ""
}

func (t *Token) Print() {
	fmt.Printf("[type: %s, literal: %s]\n", TokenTypeToString(t.Type), t.Literal)
}

type Token struct {
	Type     TokenType
	Literal  string
	Line     int
	Relative int
	Absolute int
	Length   int
}
