package logic

import (
	"fmt"
)

// token
type TokenType int

const (
	TOKEN_ILLEGAL TokenType = (iota - 1) // Illegal token
	TOKEN_EOF                            //End Of File
	TOKEN_AND                            // &&
	TOKEN_OR                             // OR
	TOKEN_BANG                           // !
	TOKEN_LPAREN                         // (
	TOKEN_RPAREN                         // )

	TOKEN_IDENTIFIER //identifier
	TOKEN_TRUE       //true
	TOKEN_FALSE      //false
)

// for debug & testing
func (tt TokenType) String() string {
	switch tt {
	case TOKEN_ILLEGAL:
		return "ILLEGAL"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_BANG:
		return "!"
	case TOKEN_AND:
		return "AND"
	case TOKEN_OR:
		return "OR"

	case TOKEN_LPAREN:
		return "("
	case TOKEN_RPAREN:
		return ")"
	case TOKEN_TRUE:
		return "TRUE"
	case TOKEN_FALSE:
		return "FALSE"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	default:
		return "UNKNOWN"
	}
}

var keywords = map[string]TokenType{
	"true":  TOKEN_TRUE,
	"false": TOKEN_FALSE,
}

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

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TOKEN_IDENTIFIER
}
