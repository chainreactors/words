package mask

import (
	"fmt"
)

// token
type TokenType int

const (
	TOKEN_ILLEGAL TokenType = (iota - 1) // Illegal token
	TOKEN_EOF                            //End Of File

	TOKEN_START      // ?
	TOKEN_REPEAT     // #
	TOKEN_SPLIT      // ,
	TOKEN_LPAREN     // {
	TOKEN_RPAREN     // }
	TOKEN_ESCAPE     // \
	TOKEN_NUMBER     //10 or 10.1
	TOKEN_IDENTIFIER //identifier
)

// for debug & testing
func (tt TokenType) String() string {
	switch tt {
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_START:
		return "START"
	case TOKEN_REPEAT:
		return "#"
	case TOKEN_ESCAPE:
		return "ESCAPE"
	case TOKEN_LPAREN:
		return "{"
	case TOKEN_RPAREN:
		return "}"
	case TOKEN_SPLIT:
		return "SPLIT"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	default:
		return "UNKNOWN"
	}
}

//var keywords = map[string]TokenType{}

type Token struct {
	Pos     Position
	Type    TokenType
	Literal string
}

// Stringer method for Token
func (t Token) String() string {
	return fmt.Sprintf("Position: %s, Type: %s, Literal: %s", t.Pos, t.Type, t.Literal)
}

// Position is the location of a code point in the source
type Position struct {
	Filename string
	Offset   int //offset relative to entire file
	Line     int
	Col      int //offset relative to each line
}

// Stringer method for Position
func (p Position) String() string {
	var msg string
	if p.Filename == "" {
		msg = fmt.Sprint(" <", p.Line, ":", p.Col, "> ")
	} else {
		msg = fmt.Sprint(" <", p.Filename, ":", p.Line, ":", p.Col, "> ")
	}

	return msg
}

// We could not use `Line()` as function name, because `Line` is the struct's field
func (p Position) Sline() string { //String line
	var msg string
	if p.Filename == "" {
		msg = fmt.Sprint(p.Line)
	} else {
		msg = fmt.Sprint(" <", p.Filename, ":", p.Line, "> ")
	}
	return msg
}
