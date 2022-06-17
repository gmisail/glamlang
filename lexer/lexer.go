package lexer

import (
	"container/list"
	"unicode"
	"fmt"
)

type Lexer struct {
	current int
	line    int
	input   string
	Tokens  *list.List
}

var keywords = map[string] TokenType {
		
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

	fmt.Println(literal)

	if tokenType, ok := keywords[literal]; ok {
		return tokenType, literal
	}

	// if the keyword isn't registered, then it is an ID
	return IDENTIFIER, literal
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

func (l *Lexer) AddToken(tokenType TokenType, literal string) *Token {
	t := &Token{Type: tokenType, Literal: literal }
	l.Tokens.PushBack(t)
	return t
}

func (l *Lexer) AddKeyword(tokenType TokenType) *Token {
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
		} else {
			panic("Unknown token: " + string(currentChar)) 
		}
	}

	return true
}

func ScanTokens(input string) *Lexer {
	lexer := Lexer{current: -1, line: 0, input: input, Tokens: list.New()}

	for lexer.ScanToken() {
	}

	return &lexer
}
