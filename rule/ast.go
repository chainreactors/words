package rule

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type Node interface {
	Pos() Position // position of first character belonging to the node
	End() Position // position of first character immediately after the node

	TokenLiteral() string
	String() string
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Expressions []Expression
}

func (p *Program) Pos() Position {
	if len(p.Expressions) > 0 {
		return p.Expressions[0].Pos()
	}
	return Position{}
}

func (p *Program) End() Position {
	aLen := len(p.Expressions)
	if aLen > 0 {
		return p.Expressions[aLen-1].End()
	}
	return Position{}
}

func (p *Program) TokenLiteral() string {
	if len(p.Expressions) > 0 {
		return p.Expressions[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Expressions {
		out.WriteString(s.String())
	}
	return out.String()
}

type RuleExpression struct {
	Functions []FunctionExpression
}

func (ie *RuleExpression) Pos() Position   { return ie.Functions[0].Pos() }
func (ie *RuleExpression) End() Position   { return ie.Functions[len(ie.Functions)-1].Pos() }
func (ie *RuleExpression) expressionNode() {}
func (ie *RuleExpression) TokenLiteral() string {
	var s strings.Builder
	for _, rule := range ie.Functions {
		s.WriteString(rule.TokenLiteral())
		s.WriteString(" ")
	}
	return s.String()
}

func (ie *RuleExpression) String() string {
	var s strings.Builder
	for _, rule := range ie.Functions {
		s.WriteString(rule.String())
	}
	return s.String()
}

type FunctionExpression struct {
	FunctionToken Token
	X             Token
	Y             Token
}

func (ie *FunctionExpression) Pos() Position { return ie.FunctionToken.Pos }
func (ie *FunctionExpression) End() Position {
	if ie.Y.NotNull() {
		return ie.Y.Pos
	} else if ie.X.NotNull() {
		return ie.X.Pos
	} else {
		return ie.FunctionToken.Pos
	}
}
func (ie *FunctionExpression) expressionNode() {}
func (ie *FunctionExpression) TokenLiteral() string {
	var s strings.Builder
	s.WriteString(ie.FunctionToken.Literal)
	if ie.X.NotNull() {
		s.WriteString(ie.X.Literal)
	}
	if ie.Y.NotNull() {
		s.WriteString(ie.Y.Literal)
	}
	return s.String()
}

func (ie *FunctionExpression) String() string {
	var s strings.Builder
	s.WriteString(ie.FunctionToken.String())
	s.WriteString(ie.X.String())
	s.WriteString(ie.Y.String())
	return s.String()
}

func (ie *FunctionExpression) Tokens() Tokens {
	return NewTokens(ie.FunctionToken, ie.X, ie.Y)
}

func (ie *FunctionExpression) IsValid() bool {
	if ie.FunctionToken.Type == TOKEN_NULL {
		return false
	}
	return true
}

type NumberLiteral struct {
	Token Token
	Value int64
}

func (nl *NumberLiteral) Pos() Position { return nl.Token.Pos }
func (nl *NumberLiteral) End() Position {
	length := utf8.RuneCountInString(nl.Token.Literal)
	pos := nl.Token.Pos
	return Position{Line: pos.Line, Col: pos.Col + length}
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }
