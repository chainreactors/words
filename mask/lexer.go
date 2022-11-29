package mask

import (
	"unicode"
)

// Lexer
type Lexer struct {
	input        []rune
	ch           rune //current character
	position     int  //character offset
	readPosition int  //reading offset
	tokenCache   []Token
	curCache     int
	line         int
	col          int
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
	if l.tokenCache == nil || l.curCache+1 > len(l.tokenCache) {
		l.tokenCache = l.ReadTokens()
		l.curCache = 0
	}
	tok := l.tokenCache[l.curCache]
	l.curCache++
	return tok
}

func (l *Lexer) ReadTokens() []Token {
	var toks []Token
	//l.skipWhitespace()

	pos := l.getPos()

	switch l.ch {
	case '{':
		toks = l.readMask()
		return toks
	case 0:
		var tok Token
		tok.Literal = "<EOF>"
		tok.Type = TOKEN_EOF
		tok.Pos = pos
		toks = append(toks, tok)
	default:
		var tok Token
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier(false)
			tok.Pos = pos
			tok.Type = TOKEN_IDENTIFIER
			toks = append(toks, tok)
			return toks
		} else {
			tok = newToken(TOKEN_ILLEGAL, l.ch, pos)
			toks = append(toks, tok)
		}
	}

	l.readNext()
	return toks
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

func (l *Lexer) readMask() []Token {

	startPos := l.getPos()
	lparenTok := newToken(TOKEN_LPAREN, l.ch, l.getPos())
	//toks = append(toks, newToken(TOKEN_LPAREN, l.ch, l.getPos()))
	l.readNext()
	if isStart(l.ch) {
		var toks []Token
		toks = append(toks, lparenTok)
		toks = append(toks, newToken(TOKEN_START, l.ch, l.getPos()))
		l.readNext()
		for l.ch != '#' && l.ch != '}' {
			if isDigit(l.ch) {
				toks = append(toks, Token{
					l.getPos(),
					TOKEN_NUMBER,
					l.readIdentifier(true),
				})
			} else if isLetter(l.ch) {
				toks = append(toks, Token{
					l.getPos(),
					TOKEN_IDENTIFIER,
					l.readIdentifier(true),
				})
			} else if isSplit(l.ch) {
				l.readNext()
			}
		}

		if l.ch == '#' {
			toks = append(toks, newToken(TOKEN_REPEAT, l.ch, l.getPos()))
			l.readNext()
			toks = append(toks, Token{
				l.getPos(),
				TOKEN_NUMBER,
				l.readNumber(),
			})
		}

		toks = append(toks, newToken(TOKEN_RPAREN, l.ch, l.getPos()))
		l.readNext()
		return toks
	} else {
		return []Token{Token{
			Pos:     startPos,
			Type:    TOKEN_IDENTIFIER,
			Literal: "{" + l.readIdentifier(false),
		}}
	}
}

func (l *Lexer) readIdentifier(safe bool) string {
	position := l.position
	if safe {
		for (isLetter(l.ch) || isDigit(l.ch)) && !(l.ch == '#' || l.ch == '}') {
			l.readNext()
		}
	} else {
		for isLetter(l.ch) || isDigit(l.ch) {
			l.readNext()
		}
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
		Offset: l.position,
		Line:   l.line,
		Col:    l.col,
	}
}

func newToken(tokenType TokenType, ch rune, pos Position) Token {
	return Token{Type: tokenType, Literal: string(ch), Pos: pos}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isSplit(ch rune) bool {
	if ch == '|' || ch == '-' {
		return true
	}
	return false
}

func isStart(ch rune) bool {
	if ch == '?' || ch == '@' || ch == '$' {
		return true
	}
	return false
}

func isLetter(ch rune) bool {
	return ch != '{' && ch != 0 && !isDigit(ch) && !isStart(ch) && !isSplit(ch) && ch != '}'
}
