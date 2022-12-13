package mask

import (
	"fmt"
	strconv "strconv"
	"unicode/utf8"
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

var SpecialWords map[string][]string = map[string][]string{}

//func AddCustomWord(s []string) {
//	CustomWords = append(CustomWords, s)
//}

func string2Bytes(s string) []string {
	ss := make([]string, len(s))
	for i := 0; i < len(s); i++ {
		ss[i] = string(s[i])
	}
	return ss
}

func ParseCharacterSetWithSpecial(s string) []string {
	if ss, ok := SpecialWords[s]; ok {
		return ss
	} else {
		return nil
	}
}

func ParseCharacterSetWithIdent(s string) []string {
	var cs []string
	for i := 0; i < len(s); i++ {
		cs = append(cs, string2Bytes(MetawordMap[string(s[i])])...)
	}
	return cs
}

func ParseCharacterSetWithNumber(s string, custom [][]string) []string {
	var cs []string
	for i := 0; i < len(s); i++ {
		if len(custom) >= i+1 {
			cs = append(cs, custom[i]...)
		} else {
			fmt.Printf("index %d out of dicts, not enough dict\n", i)
		}
	}
	return cs
}

type (
	prefixParseFn func() Expression
)

func NewParser(l *Lexer, params [][]string) *Parser {
	p := &Parser{
		l:          l,
		errors:     []string{},
		errorLines: []string{},
		params:     params,
	}

	p.registerAction()

	p.nextToken()
	p.nextToken()
	return p
}

type Parser struct {
	l          *Lexer
	errors     []string //error messages
	errorLines []string //for using with wasm communication.

	curToken       Token
	peekToken      Token
	tokenCache     []Token
	curCache       int
	params         [][]string
	prefixParseFns map[TokenType]prefixParseFn
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerAction() {
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(TOKEN_NUMBER, p.parseNumber)
	p.registerPrefix(TOKEN_IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(TOKEN_LPAREN, p.parseMaskExpression)
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	for p.curToken.Type != TOKEN_EOF {
		program.Expressions = append(program.Expressions, p.parseExpression())
		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpression() Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseMaskExpression() Expression {
	expression := &MaskExpression{Start: p.peekToken}

	p.nextToken()
	for p.peekToken.Type != TOKEN_REPEAT && p.peekToken.Type != TOKEN_RPAREN {
		if p.peekToken.Type == TOKEN_IDENTIFIER {
			if expression.Start.Literal == "@" {
				expression.CharacterSet = append(expression.CharacterSet, ParseCharacterSetWithSpecial(p.peekToken.Literal)...)
			} else {
				expression.CharacterSet = ParseCharacterSetWithIdent(p.peekToken.Literal)
			}
		} else if p.peekToken.Type == TOKEN_NUMBER {
			expression.CharacterSet = ParseCharacterSetWithNumber(p.peekToken.Literal, p.params)
		}
		p.nextToken()
	}

	if p.peekToken.Type == TOKEN_REPEAT {
		expression.RepeatToken = p.peekToken
		p.nextToken()
		if p.peekToken.Type == TOKEN_NUMBER {
			expression.Repeat, _ = strconv.Atoi(p.peekToken.Literal)
		}
		p.nextToken()
	}

	if !p.expectPeek(TOKEN_RPAREN) {
		return nil
	}
	return expression
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

//DEBUG ONLY
func (p *Parser) debugToken(message string) {
	fmt.Printf("%s, curToken = %s, curToken.Pos = %d, peekToken = %s, peekToken.Pos=%d\n", message, p.curToken.Literal, p.curToken.Pos.Line, p.peekToken.Literal, p.peekToken.Pos.Line)
}

func (p *Parser) debugNode(message string, node Node) {
	fmt.Printf("%s, Node = %s\n", message, node.String())
}
