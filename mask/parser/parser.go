package parser

import (
	"fmt"
	"github.com/chainreactors/words/mask/ast"
	"github.com/chainreactors/words/mask/lexer"
	"github.com/chainreactors/words/mask/token"
	strconv "strconv"
	"unicode/utf8"
)

const (
	_ int = iota
	LOWEST

	SUM
	PRODUCT
	PREFIX
)

var (
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letter       = Lowercase + Uppercase
	Digit        = "0123456789"
	LowercaseHex = Digit + "abcdef"
	UppercaseHex = Digit + "ABCDEF"
	Hex          = Digit + "abcdefABCDEF"
	Punctuation  = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Printable    = Letter + Digit + Punctuation
	Whitespace   = "\t\n\r\x0b\x0c"
)

var MetawordMap = map[string]string{
	"l": Lowercase,
	"u": Uppercase,
	"w": Letter,
	"d": Digit,
	"h": LowercaseHex,
	"H": UppercaseHex,
	"x": Hex,
	"p": Punctuation,
	"P": Printable,
	"s": Whitespace,
}

func string2Bytes(s string) []string {
	ss := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		ss[i] = string(s[i])
	}
	return ss
}

func ParseCharacterSet(s string) []string {
	var cs []string
	for i := 0; i < len(s); i++ {
		cs = append(cs, string2Bytes(MetawordMap[string(s[i])])...)
	}
	return cs
}

var precedences = map[token.TokenType]int{

	//token.TOKEN_PLUS:     SUM,
	//token.TOKEN_MINUS:    SUM,
	//token.TOKEN_MULTIPLY: PRODUCT,
	//token.TOKEN_DIVIDE:   PRODUCT,
	//token.TOKEN_MOD:      PRODUCT,
	//token.TOKEN_POWER:    PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l          *lexer.Lexer
	errors     []string //error messages
	errorLines []string //for using with wasm communication.

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewParser(l *lexer.Lexer) *Parser {
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
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	//p.registerPrefix(token.TOKEN_ILLEGAL, p.parsePrefixIllegalExpression)
	p.registerPrefix(token.TOKEN_NUMBER, p.parseNumber)
	p.registerPrefix(token.TOKEN_IDENTIFIER, p.parseIdentifier)
	//p.registerPrefix(token.TOKEN_PLUS, p.parsePrefixExpression)
	//p.registerPrefix(token.TOKEN_MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TOKEN_LPAREN, p.parseMaskExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	//p.registerPrefix(token.TOKEN_ILLEGAL, p.parseInfixIllegalExpression)
	//p.registerInfix(token.TOKEN_PLUS, p.parseInfixExpression)
	//p.registerInfix(token.TOKEN_MINUS, p.parseInfixExpression)
	//p.registerInfix(token.TOKEN_MULTIPLY, p.parseInfixExpression)
	//p.registerInfix(token.TOKEN_DIVIDE, p.parseInfixExpression)
	//p.registerInfix(token.TOKEN_MOD, p.parseInfixExpression)
	//p.registerInfix(token.TOKEN_POWER, p.parseInfixExpression)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	for p.curToken.Type != token.TOKEN_EOF {
		program.Expressions = append(program.Expressions, p.parseExpression(LOWEST))
		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// Run the infix function until the next token has a higher precedence.
	//for precedence < p.peekPrecedence() {
	//	infix := p.infixParseFns[p.peekToken.Type]
	//	if infix == nil {
	//		return leftExp
	//	}
	//	p.nextToken()
	//	leftExp = infix(leftExp)
	//}

	return leftExp
}

func (p *Parser) parseMaskExpression() ast.Expression {
	expression := &ast.MaskExpression{Start: p.peekToken}
	p.nextToken()
	if p.expectPeek(token.TOKEN_IDENTIFIER) {
		expression.CharacterSet = ParseCharacterSet(p.curToken.Literal)
	}

	if p.peekToken.Type == token.TOKEN_REPEAT {
		expression.RepeatToken = p.peekToken
		p.nextToken()
		if p.peekToken.Type == token.TOKEN_NUMBER {
			expression.Repeat, _ = strconv.Atoi(p.peekToken.Literal)
		}
	}

	p.nextToken()
	if !p.expectPeek(token.TOKEN_RPAREN) {
		return nil
	}
	return expression
}

//func (p *Parser) parsePrefixExpression() ast.Expression {
//	expression := &ast.PrefixExpression{Start: p.curToken, Operator: p.curToken.Literal}
//	p.nextToken()
//	expression.Right = p.parseExpression(PREFIX)
//
//	return expression
//}

//func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
//	expression := &ast.InfixExpression{
//		Start:    p.curToken,
//		Operator: p.curToken.Literal,
//		Left:     left,
//	}
//	precedence := p.curPrecedence()
//
//	// if the token is '**', we process it specially. e.g. 3 ** 2 ** 3 = 3 ** (2 ** 3)
//	// i.e. Exponent operator '**'' has right-to-left associativity
//	//if p.curTokenIs(token.TOKEN_POWER) {
//	//	precedence--
//	//}
//
//	p.nextToken()
//	expression.Right = p.parseExpression(precedence)
//
//	return expression
//}

//func (p *Parser) parseGroupedExpression() ast.Expression {
//	p.nextToken()
//
//	exp := p.parseExpression(LOWEST)
//
//	if !p.expectPeek(token.TOKEN_RPAREN) {
//		return nil
//	}
//
//	return exp
//}
//
//func (p *Parser) parsePrefixIllegalExpression() ast.Expression {
//	msg := fmt.Sprintf("Syntax Error:%v - Illegal token found. Literal: '%s'", p.curToken.Pos, p.curToken.Literal)
//	p.errors = append(p.errors, msg)
//	p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
//	return nil
//}
//
//func (p *Parser) parseInfixIllegalExpression() ast.Expression {
//	msg := fmt.Sprintf("Syntax Error:%v - Illegal token found. Literal: '%s'", p.curToken.Pos, p.curToken.Literal)
//	p.errors = append(p.errors, msg)
//	p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
//	return nil
//}

func (p *Parser) parseNumber() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}

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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	if t != token.TOKEN_EOF {
		msg := fmt.Sprintf("Syntax Error:%v- no prefix parse functions for '%s' found", p.curToken.Pos, t)
		p.errors = append(p.errors, msg)
		p.errorLines = append(p.errorLines, p.curToken.Pos.Sline())
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
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

func (p *Parser) debugNode(message string, node ast.Node) {
	fmt.Printf("%s, Node = %s\n", message, node.String())
}
