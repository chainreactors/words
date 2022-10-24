package ast

import (
	"bytes"
	"github.com/chainreactors/words/mask/token"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Node interface {
	Pos() token.Position // position of first character belonging to the node
	End() token.Position // position of first character immediately after the node

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

func (p *Program) Pos() token.Position {
	if len(p.Expressions) > 0 {
		return p.Expressions[0].Pos()
	}
	return token.Position{}
}

func (p *Program) End() token.Position {
	aLen := len(p.Expressions)
	if aLen > 0 {
		return p.Expressions[aLen-1].End()
	}
	return token.Position{}
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

type MaskExpression struct {
	Start        token.Token
	CharacterSet []string
	RepeatToken  token.Token
	Repeat       int
}

func (ie *MaskExpression) Pos() token.Position { return ie.Start.Pos }
func (ie *MaskExpression) End() token.Position {
	pos := ie.RepeatToken.Pos
	pos.Offset = pos.Offset + len(strconv.Itoa(int(ie.Repeat)))
	return pos
}

func (ie *MaskExpression) expressionNode()      {}
func (ie *MaskExpression) TokenLiteral() string { return ie.Start.Literal }
func (ie *MaskExpression) String() string {
	var out bytes.Buffer

	out.WriteString("{?")
	out.WriteString(strings.Join(ie.CharacterSet, ""))
	out.WriteString("#")
	out.WriteString(strconv.Itoa(int(ie.Repeat)))
	out.WriteString("}")

	return out.String()
}

// 1 + 2 * 3
//type InfixExpression struct {
//	Start    token.Start
//	Operator string
//	Right    Expression
//	Left     Expression
//}
//
//func (ie *InfixExpression) Pos() token.Position { return ie.Start.Pos }
//func (ie *InfixExpression) End() token.Position { return ie.Right.End() }
//
//func (ie *InfixExpression) expressionNode()      {}
//func (ie *InfixExpression) TokenLiteral() string { return ie.Start.Literal }
//func (ie *InfixExpression) String() string {
//	var out bytes.Buffer
//
//	out.WriteString("{")
//	out.WriteString(ie.Left.String())
//	out.WriteString(" " + ie.Operator + " ")
//	out.WriteString(ie.Right.String())
//	out.WriteString("}")
//
//	return out.String()
//}

//// -2, -3
//type PrefixExpression struct {
//	Start    token.Start
//	Operator string
//	Right    Expression
//}
//
//func (pe *PrefixExpression) Pos() token.Position { return pe.Start.Pos }
//func (pe *PrefixExpression) End() token.Position { return pe.Right.End() }
//
//func (pe *PrefixExpression) expressionNode()      {}
//func (pe *PrefixExpression) TokenLiteral() string { return pe.Start.Literal }
//
//func (pe *PrefixExpression) String() string {
//	var out bytes.Buffer
//
//	out.WriteString("{")
//	out.WriteString(pe.Operator)
//	out.WriteString(pe.Right.String())
//	out.WriteString("}")
//
//	return out.String()
//}

type NumberLiteral struct {
	Token token.Token
	Value int64
}

func (nl *NumberLiteral) Pos() token.Position { return nl.Token.Pos }
func (nl *NumberLiteral) End() token.Position {
	length := utf8.RuneCountInString(nl.Token.Literal)
	pos := nl.Token.Pos
	return token.Position{Filename: pos.Filename, Line: pos.Line, Col: pos.Col + length}
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string       { return nl.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) Pos() token.Position { return i.Token.Pos }
func (i *Identifier) End() token.Position {
	length := utf8.RuneCountInString(i.Value)
	return token.Position{Filename: i.Token.Pos.Filename, Line: i.Token.Pos.Line, Col: i.Token.Pos.Col + length}
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
