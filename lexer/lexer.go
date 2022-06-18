package lexer

import (
	"unicode"
)

type Lexer struct {
	current int
	line    int
	input   string
	Tokens  []Token
}

var keywords = map[string]TokenType{
	"let":   LET,
	"while": WHILE,
	"for":   FOR,
	"if":    IF,
	"else":  ELSE,
	"true":  TRUE,
	"false": FALSE,
}

func (l *Lexer) IsAtEnd() bool {
	return len(l.input) <= l.current
}

func (l *Lexer) AdvanceChar() {
	l.current += 1
}

func (l *Lexer) CharAt(i int) rune {
	if i < len(l.input) && i >= 0 {
		return rune(l.input[i])
	}

	return rune(0)
}

func (l *Lexer) CurrentChar() rune {
	return l.CharAt(l.current)
}

func (l *Lexer) PeekChar() rune {
	return l.CharAt(l.current + 1)
}

func (l *Lexer) IsLetter(char rune) bool {
	if unicode.IsLetter(char) || char == '_' {
		return true
	}

	return false
}

func (l *Lexer) ScanKeyword() (TokenType, string) {
	start := l.current

	for l.IsLetter(l.PeekChar()) || unicode.IsDigit(l.PeekChar()) {
		l.AdvanceChar()
	}

	end := l.current
	literal := l.input[start:(end + 1)]

	if tokenType, ok := keywords[literal]; ok {
		return tokenType, literal
	}

	// if the keyword isn't registered, then it is an ID
	return IDENTIFIER, literal
}

func (l *Lexer) ScanNumber() (TokenType, string) {
	start := l.current
	tokenType := INT

	for unicode.IsDigit(l.PeekChar()) {
		l.AdvanceChar()
	}

	// get numbers after decimal point
	if l.PeekChar() == '.' {
		l.AdvanceChar()

		initialDecimal := l.current

		for unicode.IsDigit(l.PeekChar()) {
			l.AdvanceChar()
		}

		if l.current == initialDecimal {
			panic("Unexpected '.' after integer.")
		}

		tokenType = FLOAT
	}

	end := l.current
	literal := l.input[start:(end + 1)]

	return tokenType, literal
}

func (l *Lexer) SkipWhitespace() {
	for {
		switch l.CurrentChar() {
		case '\r', '\t', ' ':
			l.AdvanceChar()
		case '\n':
			l.line += 1
			l.AdvanceChar()
		default:
			return
		}
	}
}

func (l *Lexer) AddToken(tokenType TokenType, literal string) Token {
	t := Token{Type: tokenType, Literal: literal}
	l.Tokens = append(l.Tokens, t)
	return t
}

func (l *Lexer) AddKeyword(tokenType TokenType) Token {
	return l.AddToken(tokenType, "")
}

func (l *Lexer) ScanToken() bool {
	l.AdvanceChar()

	if l.IsAtEnd() {
		return false
	}

	l.SkipWhitespace()

	currentChar := l.CurrentChar()

	switch currentChar {
	case '(':
		l.AddKeyword(L_PAREN)
	case ')':
		l.AddKeyword(R_PAREN)
	case '{':
		l.AddKeyword(L_BRACE)
	case '}':
		l.AddKeyword(R_BRACE)
	case '[':
		l.AddKeyword(L_BRACKET)
	case ']':
		l.AddKeyword(R_BRACKET)
	case ',':
		l.AddKeyword(COMMA)
	case '.':
		l.AddKeyword(PERIOD)
	case '+':
		l.AddKeyword(ADD)
	case '-':
		l.AddKeyword(SUB)
	case '*':
		l.AddKeyword(MULT)
	case '/':
		l.AddKeyword(DIV)
	case '=':
		l.AddKeyword(EQUAL)
	case '>':
		l.AddKeyword(GT)
	case '<':
		l.AddKeyword(LT)
	case '!':
		l.AddKeyword(BANG)
	case ':':
		l.AddKeyword(COLON)
	case '"':
		l.AddKeyword(QUOTE)
	default:
		if l.IsLetter(currentChar) {
			tokenType, literal := l.ScanKeyword()

			if tokenType == IDENTIFIER {
				l.AddToken(tokenType, literal)
			} else {
				l.AddKeyword(tokenType)
			}
		} else if unicode.IsDigit(currentChar) {
			tokenType, literal := l.ScanNumber()
			l.AddToken(tokenType, literal)
		} else {
			panic("Unknown token: " + string(currentChar))
		}
	}

	return true
}

func ScanTokens(input string) *Lexer {
	lexer := Lexer{current: -1, line: 0, input: input, Tokens: make([]Token, 0)}

	for lexer.ScanToken() {
	}

	return &lexer
}
