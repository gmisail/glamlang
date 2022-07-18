package main

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
	STRUCT
	ARROW
	THICK_ARROW
	TRUE
	FALSE
	NULL
	FLOAT
	INT
)

func tokenTypeToString(token TokenType) string {
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
	case STRUCT:
		return "STRUCT"
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

func (t *Token) Print() {
	fmt.Printf("[type: %s, literal: %s]\n", tokenTypeToString(t.Type), t.Literal)
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}
