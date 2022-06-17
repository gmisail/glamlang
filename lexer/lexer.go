package lexer

import (
	"container/list"
)

type Lexer struct {
	current int
	line    int
	input   string
	Tokens  *list.List
}

// func (l *Lexer) CurrentToken() *Token {
// 	token := l.Tokens.Front()

// 	if token == nil {
// 		return nil
// 	}

// 	return token.Value.(*Token)
// }

// func (l *Lexer) PreviousToken() *Token {
// 	token := l.Tokens.Front()

// 	if token == nil || token.Prev() == nil {
// 		return nil
// 	}

// 	return (token.Prev().Value).(*Token)
// }

// func (l *Lexer) NextToken() *Token {
// 	l.current += 1

// 	return l.CurrentToken()
// }

func (l *Lexer) IsAtEnd() bool {
	return len(l.input) <= l.current
}

func (l *Lexer) AdvanceChar() {
	l.current += 1
}

func (l *Lexer) CurrentChar() rune {
	if l.IsAtEnd() {
		return rune(0)
	}

	return rune(l.input[l.current])
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

func (l *Lexer) AddToken(tokenType TokenType) *Token {
	t := &Token{Type: tokenType}
	l.Tokens.PushBack(t)
	return t
}

func (l *Lexer) ScanToken() bool {
	if l.IsAtEnd() {
		return false
	}

	l.AdvanceChar()
	l.SkipWhitespace()

	currentChar := l.CurrentChar()

	switch currentChar {
	case '(':
		l.AddToken(L_PAREN)
	case ')':
		l.AddToken(R_PAREN)
	case '{':
		l.AddToken(L_BRACE)
	case '}':
		l.AddToken(R_BRACE)
	case '[':
		l.AddToken(L_BRACKET)
	case ']':
		l.AddToken(R_BRACKET)
	case ',':
		l.AddToken(COMMA)
	case '.':
		l.AddToken(PERIOD)
	case '+':
		l.AddToken(ADD)
	case '-':
		l.AddToken(SUB)
	case '*':
		l.AddToken(MULT)
	case '/':
		l.AddToken(DIV)
	case '=':
		l.AddToken(EQUAL)
	case '>':
		l.AddToken(GT)
	case '<':
		l.AddToken(LT)
	case '!':
		l.AddToken(BANG)
	case ':':
		l.AddToken(COLON)
	case '"':
		l.AddToken(QUOTE)
	default:
		return false
	}

	return true
}

func ScanTokens(input string) *Lexer {
	lexer := Lexer{current: -1, line: 0, input: input, Tokens: list.New()}

	for lexer.ScanToken() {
	}

	return &lexer
}
