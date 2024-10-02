package logic

import (
	"unicode"
)

// Lexer
type Lexer struct {
	filename     string
	input        []rune
	ch           rune //current character
	position     int  //character offset
	readPosition int  //reading offset

	line int
	col  int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.ch = ' '
	l.position = 0
	l.readPosition = 0

	l.line = 1
	l.col = 0

	l.readNext()
	//0xFEFF: BOM(byte order mark), only permitted as very first character
	if l.ch == 0xFEFF {
		l.readNext() //ignore BOM at file beginning
	}

	return l
}

func (l *Lexer) readNext() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
		if l.ch == '\n' {
			l.col = 0
			l.line++
		} else {
			l.col += 1
		}
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peek() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	pos := l.getPos()

	switch l.ch {
	case '(':
		tok = newToken(TOKEN_LPAREN, l.ch)
	case ')':
		tok = newToken(TOKEN_RPAREN, l.ch)
	case '!':
		tok = newToken(TOKEN_BANG, l.ch)
	case '&':
		if l.peek() == '&' {
			tok = Token{Type: TOKEN_AND, Literal: string(l.ch) + string(l.peek())}
			l.readNext()
		}
	case '|':
		if l.peek() == '|' {
			tok = Token{Type: TOKEN_OR, Literal: string(l.ch) + string(l.peek())}
			l.readNext()
		}
	case 0:
		tok.Literal = "<EOF>"
		tok.Type = TOKEN_EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Pos = pos
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch)
		}
	}

	tok.Pos = pos
	l.readNext()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readNext()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readNext()
	}
}

func (l *Lexer) getPos() Position {
	return Position{
		Filename: l.filename,
		Offset:   l.position,
		Line:     l.line,
		Col:      l.col,
	}
}

func newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}
