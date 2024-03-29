package rule

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

type (
	prefixParseFn func() Expression
)

type Parser struct {
	l          *Lexer
	errors     []string //error messages
	errorLines []string //for using with wasm communication.

	curToken       Token
	peekToken      Token
	prefixParseFns map[TokenType]prefixParseFn
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:          l,
		errors:     []string{},
		errorLines: []string{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram(filterExpr *RuleExpression) *Program {
	program := &Program{}
	for p.curToken.Type != TOKEN_EOF {
		program.Expressions = append(program.Expressions, p.parseExpression(filterExpr))
	}

	return program
}

func (p *Parser) parseExpression(filterExpr *RuleExpression) Expression {
	expr := p.parseRuleExpression(p.nextLine()).(*RuleExpression)
	if filterExpr != nil {
		expr.Functions = append(expr.Functions, filterExpr.Functions...)
	}
	return expr
}

func (p *Parser) parseRuleExpression(tokens []Token) Expression {
	var pos int
	expr := RuleExpression{}
	functionExpr := FunctionExpression{}
	for _, tok := range tokens {
		t := tok
		if tok.Type == TOKEN_SPLIT || tok.Type == TOKEN_EOF || tok.Type == TOKEN_NULL {
			expr.Functions = append(expr.Functions, functionExpr)
			pos = 0
			functionExpr = FunctionExpression{}
			continue
		}
		if pos == 0 {
			functionExpr.FunctionToken = t
		} else if pos == 1 {
			functionExpr.X = t
		} else if pos == 2 {
			functionExpr.Y = t
		}
		pos++
	}
	return &expr
}

func (p *Parser) parseNumber() Expression {
	lit := &NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 32)
	if err != nil {
		msg := fmt.Sprintf("Syntax Error:%v - could not parse %q as float", p.curToken.Pos, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
		return nil
	}
	lit.Value = value
	return lit
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

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) nextLine() []Token {
	var tokens []Token
	for p.curToken.Type == TOKEN_LINEEOF {
		p.nextToken()
	}

	for p.curToken.Type != TOKEN_LINEEOF && p.curToken.Type != TOKEN_EOF {
		if p.curToken.Type == TOKEN_SPLIT && p.peekToken.Type != TOKEN_FUNCTION {

		} else {
			tokens = append(tokens, p.curToken)
		}

		p.nextToken()
	}

	// 跳过换行符
	p.nextToken()
	tokens = append(tokens, newToken(TOKEN_EOF, ' '))
	return tokens
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

//DEBUG ONLY
func (p *Parser) debugToken(message string) {
	fmt.Printf("%s, curToken = %s, curToken.Pos = %d, peekToken = %s, peekToken.Pos=%d\n", message, p.curToken.Literal, p.curToken.Pos.Line, p.peekToken.Literal, p.peekToken.Pos.Line)
}

func (p *Parser) debugNode(message string, node Node) {
	fmt.Printf("%s, Node = %s\n", message, node.String())
}
