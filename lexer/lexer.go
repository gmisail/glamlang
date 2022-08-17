package lexer

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	current int
	line    int
	input   string
	index   int
	Tokens  []Token
}

type LexerError struct {
	line    int
	message string
}

type TokenPair struct {
	char      rune
	tokenType TokenType
}

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
}

func (l *LexerError) Error() string {
	return fmt.Sprintf("line %d: %s\n", l.line, l.message)
}

func (l *Lexer) IsAtEnd() bool {
	return len(l.input) <= l.current
}

func (l *Lexer) AdvanceChar() {
	l.index += 1
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

func (l *Lexer) ScanNumber() (*Token, *LexerError) {
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
			return nil, &LexerError{line: l.line, message: "Unexpected '.' after integer."}
		}

		tokenType = FLOAT
	}

	end := l.current
	literal := l.input[start:(end + 1)]

	return &Token{Type: tokenType, Literal: literal}, nil
}

func (l *Lexer) ScanConditional(options []TokenPair, fallback TokenType) TokenType {
	nextChar := l.PeekChar()

	for _, option := range options {
		if option.char == nextChar {
			l.AdvanceChar()

			return option.tokenType
		}
	}

	return fallback
}

func (l *Lexer) ScanString() string {
	// ignore opening quote
	l.AdvanceChar()

	start := l.current

	for {
		nextChar := l.PeekChar()

		if nextChar == '\n' {
			l.index = 0
			l.line += 1
		} else if nextChar == '"' || nextChar == 0 {
			break
		}

		l.AdvanceChar()
	}

	l.AdvanceChar()

	end := l.current
	literal := l.input[start:end]

	return literal
}

func (l *Lexer) SkipWhitespace() {
	for {
		switch l.CurrentChar() {
		case '\r', '\t', ' ':
			l.AdvanceChar()
		case '\n':
			l.index = 0
			l.line += 1
			l.AdvanceChar()
		default:
			return
		}
	}
}

func (l *Lexer) AddToken(tokenType TokenType, literal string) Token {
	t := Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Start:   l.index - len(literal),
		Length:  len(literal),
	}

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
		l.AddKeyword(l.ScanConditional([]TokenPair{{char: '>', tokenType: ARROW}}, SUB))
	case '*':
		l.AddKeyword(MULT)
	case '/':
		l.AddKeyword(DIV)
	case '=':
		l.AddKeyword(
			l.ScanConditional(
				[]TokenPair{{char: '=', tokenType: EQUALITY}, {char: '>', tokenType: THICK_ARROW}},
				EQUAL,
			),
		)
	case '>':
		l.AddKeyword(l.ScanConditional([]TokenPair{{char: '=', tokenType: GT_EQ}}, GT))
	case '<':
		l.AddKeyword(l.ScanConditional([]TokenPair{{char: '=', tokenType: LT_EQ}}, LT))
	case '!':
		l.AddKeyword(l.ScanConditional([]TokenPair{{char: '=', tokenType: NOT_EQUAL}}, BANG))
	case ':':
		l.AddKeyword(COLON)
	case '"':
		l.AddToken(STRING, l.ScanString())
	case '?':
		l.AddKeyword(QUESTION)
	case 0:
		return false
	default:
		if l.IsLetter(currentChar) {
			tokenType, literal := l.ScanKeyword()

			if tokenType == IDENTIFIER {
				l.AddToken(tokenType, literal)
			} else {
				l.AddKeyword(tokenType)
			}
		} else if unicode.IsDigit(currentChar) {
			token, err := l.ScanNumber()

			if err != nil {
				fmt.Print(err.Error())
				return false
			}

			l.AddToken(token.Type, token.Literal)
		} else {
			fmt.Printf("line %d: Unknown token: '%c'\n", l.line, currentChar)
			return false
		}
	}

	return true
}

func ScanTokens(input string) *Lexer {
	lexer := Lexer{current: -1, line: 1, index: -1, input: input, Tokens: make([]Token, 0)}

	for lexer.ScanToken() {
	}

	return &lexer
}
