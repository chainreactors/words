package rule

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
	statusToken  int
	line         int
	col          int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: []rune(input),
	}
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
	if l.ch == '#' && l.statusToken == 0 {
		l.skipComment()
	}
	if l.statusToken != 0 {
		l.skipSpace()
	}
	pos := l.getPos()

	switch l.ch {
	case 0:
		tok.Literal = "<EOF>"
		tok.Type = TOKEN_EOF
	case '\n':
		for l.peek() == '\n' {
			l.readNext()
		}
		tok = newToken(TOKEN_LINEEOF, l.ch)
		tok.Literal = "\\n"
		l.statusToken = 0
	case ' ':
		tok = newToken(TOKEN_SPLIT, l.ch)
		tok.Pos = pos
		l.skipSpace()
		l.statusToken = 0
		return tok
	default:
		if l.statusToken == 0 {
			if isSingleFunc(l.ch) {
				tok = newToken(TOKEN_FUNCTION, l.ch)
			} else if isDoubleFunc(l.ch) {
				tok = newToken(TOKEN_FUNCTION, l.ch)
				l.statusToken = 1
			} else if isTernaryFunc(l.ch) {
				tok = newToken(TOKEN_FUNCTION, l.ch)
				l.statusToken = 2
			} else if isFilterFunc(l.ch) {
				tok = newToken(TOKEN_FILTER, l.ch)
				l.statusToken = 1
			} else if isTernaryFilterFunc(l.ch) {
				tok = newToken(TOKEN_FILTER, l.ch)
				l.statusToken = 2
			} else {
				tok = Token{
					Type:    TOKEN_NULL,
					Literal: string(l.ch),
				}
			}
		} else {
			l.statusToken--
			if isDigit(l.ch) {
				tok.Literal = l.readNumber()
				tok.Type = TOKEN_NUMBER
				tok.Pos = pos
				return tok
			} else if isLetter(l.ch) {
				tok = newToken(TOKEN_IDENTIFIER, l.ch)
			}
		}
	}

	tok.Pos = pos
	l.readNext()
	return tok
}

func (l *Lexer) readNumber() string {
	var ret []rune

	ch := l.ch
	ret = append(ret, ch)
	l.readNext()

	for isDigit(l.ch) || l.ch == '.' {
		ret = append(ret, l.ch)
		l.readNext()
	}

	return string(ret)
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readNext()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) && l.ch != '\n' && l.ch != ' ' {
		l.readNext()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readNext()
	}
}

func (l *Lexer) skipSpace() {
	for l.ch == ' ' {
		l.readNext()
	}
}

func (l *Lexer) allTokens() []Token {
	if len(l.input) == 0 {
		return nil
	}
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}
	return tokens
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

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

var singleFuncTokens = []rune{':', 'l', 'u', 'c', 'C', 't', 'r', 'd', 'f', '{', '}', '[', ']', 'k', 'K', 'E'}
var unaryFuncTokens = []rune{'T', 'p', '$', '^', 'D', 'O', '\'', '@', 'Z', 'L', 'R', '+', '-', '.', ',', 'y', 'Y', 'e'}
var ternaryFuncTokens = []rune{'*', 's', 'i', 'o', 'x'}
var filterFuncTokens = []rune{'>', '<', '_', '!', '/', '(', ')'}
var ternaryFilterFuncTokens = []rune{'=', '%'}

func isSingleFunc(ch rune) bool {
	for _, f := range singleFuncTokens {
		if f == ch {
			return true
		}
	}
	return false
}

func isDoubleFunc(ch rune) bool {
	for _, f := range unaryFuncTokens {
		if f == ch {
			return true
		}
	}
	return false
}

func isTernaryFunc(ch rune) bool {
	for _, f := range ternaryFuncTokens {
		if f == ch {
			return true
		}
	}
	return false
}

func isFilterFunc(ch rune) bool {
	for _, f := range filterFuncTokens {
		if f == ch {
			return true
		}
	}
	return false
}

func isTernaryFilterFunc(ch rune) bool {
	for _, f := range ternaryFilterFuncTokens {
		if f == ch {
			return true
		}
	}
	return false
}

func isLetter(ch rune) bool {
	return unicode.IsPrint(ch)
}
