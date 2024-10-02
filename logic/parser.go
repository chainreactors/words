package logic

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

const (
	_ int = iota
	LOWEST
	CONDOR  // ||
	CONDAND // &&
	PREFIX
)

var precedences = map[TokenType]int{
	TOKEN_OR:  CONDOR,
	TOKEN_AND: CONDAND,
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)

type Parser struct {
	l          *Lexer
	errors     []string //error messages
	errorLines []string //for using with wasm communication.

	curToken  Token
	peekToken Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:          l,
		errors:     []string{},
		errorLines: []string{},
	}

	p.registerAction()

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerAction() {
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(TOKEN_ILLEGAL, p.parsePrefixIllegalExpression)
	p.registerPrefix(TOKEN_IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(TOKEN_LPAREN, p.parseGroupedExpression)
	p.registerPrefix(TOKEN_BANG, p.parsePrefixExpression)
	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerPrefix(TOKEN_ILLEGAL, p.parseInfixIllegalExpression)
	p.registerInfix(TOKEN_AND, p.parseInfixExpression)
	p.registerInfix(TOKEN_OR, p.parseInfixExpression)
	p.registerPrefix(TOKEN_TRUE, p.parseBooleanLiteral)
	p.registerPrefix(TOKEN_FALSE, p.parseBooleanLiteral)
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}

	program.Expression = p.parseExpression(LOWEST)
	return program
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// Run the infix function until the next token has a higher precedence.
	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBooleanLiteral() Expression {
	return &BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(TOKEN_TRUE)}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixIllegalExpression() Expression {
	msg := fmt.Sprintf("Syntax Error:%v - Illegal token found. Literal: '%s'", p.curToken.Pos, p.curToken.Literal)
	p.errors = append(p.errors, msg)
	p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
	return nil
}

func (p *Parser) parseInfixIllegalExpression() Expression {
	msg := fmt.Sprintf("Syntax Error:%v - Illegal token found. Literal: '%s'", p.curToken.Pos, p.curToken.Literal)
	p.errors = append(p.errors, msg)
	p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
	return nil
}

func (p *Parser) parseNumber() Expression {
	lit := &NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v - could not parse %q as float", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	if t != TOKEN_EOF {
		msg := fmt.Sprintf("Syntax Error:%v- no prefix parse functions for '%s' found", p.curToken.Pos, t)
		p.errors = append(p.errors, msg)
		p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
	}
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	newPos := p.curToken.Pos
	newPos.Col = newPos.Col + utf8.RuneCountInString(p.curToken.Literal)
	msg := fmt.Sprintf("Syntax Error:%v- expected next token to be %s, got %s instead", newPos, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
	p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ErrorLines() []string {
	return p.errorLines
}

// DEBUG ONLY
func (p *Parser) debugToken(message string) {
	fmt.Printf("%s, curToken = %s, curToken.Pos = %d, peekToken = %s, peekToken.Pos=%d\n", message, p.curToken.Literal, p.curToken.Pos.Line, p.peekToken.Literal, p.peekToken.Pos.Line)
}

func (p *Parser) debugNode(message string, node Node) {
	fmt.Printf("%s, Node = %s\n", message, node.String())
}
