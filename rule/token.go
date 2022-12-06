package rule

import (
	"errors"
	"fmt"
	"strconv"
)

// token
type TokenType int

const (
	TOKEN_NULL TokenType = iota
	TOKEN_IDENTIFIER
	TOKEN_EOF
	TOKEN_LINEEOF
	TOKEN_FUNCTION
	TOKEN_SPLIT
	TOKEN_NUMBER //10 or 10.1
)

//for debug & testing
func (tt TokenType) String() string {
	switch tt {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_LINEEOF:
		return "CL"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_FUNCTION:
		return "FUNCTION"
	case TOKEN_SPLIT:
		return "SPILT"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Pos     Position
	Type    TokenType
	Literal string
}

//Stringer method for Token
func (t Token) String() string {
	return fmt.Sprintf("Position: %s, Type: %s, Literal: %s", t.Pos, t.Type, t.Literal)
}

func (t Token) NotNull() bool {
	if t.Type == TOKEN_NULL {
		return false
	} else {
		return true
	}
}

var (
	NaN = errors.New("Not a number")
	EOF = errors.New("EOF")
)

func NewTokens(toks ...Token) Tokens {
	var tokens Tokens
	for _, t := range toks {
		if t.NotNull() {
			tokens.toks = append(tokens.toks, t)
		}
	}
	return tokens
}

type Tokens struct {
	toks []Token
	cur  int
}

func (t *Tokens) NextInt() (int, error) {
	if t.cur+1 > len(t.toks) {
		return 0, EOF
	}
	t.cur++
	if tok := t.toks[t.cur-1]; tok.Type == TOKEN_NUMBER {
		return strconv.Atoi(tok.Literal)
	} else {
		return 0, NaN
	}
}

func (t *Tokens) MustNextInt() int {
	i, err := t.NextInt()
	if err != nil {
		if err == EOF {
			panic(fmt.Sprintf("%s %s: %s", t.toks[t.cur-1].Literal, err.Error(), t.toks[t.cur-1].String()))
		} else {
			panic(fmt.Sprintf("%s %s: %s", t.toks[t.cur].Literal, err.Error(), t.toks[t.cur].String()))
		}
	}
	return i
}

func (t *Tokens) NextString() (string, error) {
	if t.cur+1 > len(t.toks) {
		return "", EOF
	}
	t.cur++
	return t.toks[t.cur-1].Literal, nil
}

func (t *Tokens) MustNext() string {
	s, err := t.NextString()
	if err != nil {
		if err == EOF {
			panic(fmt.Sprintf("%s %s: %s", t.toks[t.cur-1].Literal, err.Error(), t.toks[t.cur-1].String()))
		} else {
			panic(fmt.Sprintf("%s %s: %s", t.toks[t.cur].Literal, err.Error(), t.toks[t.cur].String()))
		}
	}
	return s
}

//Position is the location of a code point in the source
type Position struct {
	Filename string
	Offset   int //offset relative to entire file
	Line     int
	Col      int //offset relative to each line
}

//Stringer method for Position
func (p Position) String() string {
	var msg string
	if p.Filename == "" {
		msg = fmt.Sprint(" <", p.Line, ":", p.Col, "> ")
	} else {
		msg = fmt.Sprint(" <", p.Filename, ":", p.Line, ":", p.Col, "> ")
	}

	return msg
}

//We could not use `Line()` as function name, because `Line` is the struct's field
func (p Position) Sline() string { //String line
	var msg string
	if p.Filename == "" {
		msg = fmt.Sprint(p.Line)
	} else {
		msg = fmt.Sprint(" <", p.Filename, ":", p.Line, "> ")
	}
	return msg
}
