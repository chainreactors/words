package logic

import (
	"bytes"
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
	Expression Expression
}

func (p *Program) Pos() Position {
	return p.Expression.Pos()
}

func (p *Program) End() Position {
	return p.Expression.End()
}

func (p *Program) TokenLiteral() string {
	return p.Expression.TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer

	out.WriteString(p.Expression.String())
	return out.String()
}

// 1 + 2 * 3
type InfixExpression struct {
	Token    Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) Pos() Position { return ie.Token.Pos }
func (ie *InfixExpression) End() Position { return ie.Right.End() }

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// -2, -3
type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Pos() Position { return pe.Token.Pos }
func (pe *PrefixExpression) End() Position { return pe.Right.End() }

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type NumberLiteral struct {
	Token Token
	Value float64
}

func (nl *NumberLiteral) Pos() Position { return nl.Token.Pos }
func (nl *NumberLiteral) End() Position {
	length := utf8.RuneCountInString(nl.Token.Literal)
	pos := nl.Token.Pos
	return Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }

type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) Pos() Position { return i.Token.Pos }
func (i *Identifier) End() Position {
	length := utf8.RuneCountInString(i.Value)
	return Position{Filename: i.Token.Pos.Filename, Line: i.Token.Pos.Line, Col: i.Token.Pos.Col + length}
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type BooleanLiteral struct {
	Token Token
	Value bool
}

func (b *BooleanLiteral) Pos() Position {
	return b.Token.Pos
}

func (b *BooleanLiteral) End() Position {
	length := utf8.RuneCountInString(b.Token.Literal)
	pos := b.Token.Pos
	return Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }
