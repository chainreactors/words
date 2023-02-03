package rule

import (
	"fmt"
	"testing"
)

//var rules = ":\nc\nu\nC\n##append numbers 1...5\n$1\n$2\n$3\n$4\n$5\n$6\n$7\n$8\n$9\n$0\n$1 $2 $3\n$1 $2 $3 $4\n$1 $2 $3 $4 $5\n$1 $2 $3 $4 $5 $6\n##current year 2018-2022\n$2 $0 $1 $8\n$2 $0 $1 $9\n$2 $0 $2 $0\n$2 $0 $2 $1\n$2 $0 $2 $2\n#years or month\n$0 $1\n$0 $2\n$0 $3\n$0 $4\n$0 $5\n$0 $6\n$0 $7\n$0 $8\n$0 $9\n$1 $0\n$1 $1\n$1 $2\n$1 $3\n$1 $4\n$1 $5\n$1 $6\n$1 $7\n$1 $8\n$1 $9\n$2 $0\n$2 $1\n$2 $2\n$2 $3\n$2 $4\n$2 $5\n$2 $6\n$2 $7\n$2 $8\n$2 $9\n$3 $0\n$3 $1\n##append sepcial chars\n$!\n$@\n$#\n$$\n$! $@\n$! $@ $#\n$! $@ $# $$\n##special chars + numbers\n$1 $2 $3 $!\n$! $1 $2 $3\n$1 $@ !#\n$! $@ 1#\n##special chars + years\n$2 $0 $1 $8 $!\n$2 $0 $1 $9 $!\n$2 $0 $2 $0 $!\n$2 $0 $2 $1 $!\n$2 $0 $2 $2 $!\n$! $2 $0 $1 $8\n$! $2 $0 $1 $9\n$! $2 $0 $2 $0\n$! $2 $0 $2 $1\n$! $2 $0 $2 $2\n$2 $0 $1 $8 $! $@ $#\n$2 $0 $1 $9 $! $@ $#\n$2 $0 $2 $0 $! $@ $#\n$2 $0 $2 $1 $! $@ $#\n$2 $0 $2 $2 $! $@ $#\n$0 $1 $! \n$0 $2 $!\n$0 $3 $!\n$0 $4 $!\n$0 $5 $!\n$0 $6 $!\n$0 $7 $!\n$0 $8 $!\n$0 $9 $!\n$1 $0 $!\n$1 $1 $!\n$1 $2 $!\n$1 $3 $!\n$1 $4 $!\n$1 $5 $!\n$1 $6 $!\n$1 $7 $!\n$1 $8 $!\n$1 $9 $!\n$2 $0 $!\n$2 $1 $!\n$2 $2 $!\n$2 $3 $!\n$2 $4 $!\n$2 $5 $!\n$2 $6 $!\n$2 $7 $!\n$2 $8 $!\n$2 $9 $!\n$3 $0 $!\n$3 $1 $!\n#all above cap\nc $1\nc $2\nc $3\nc $4\nc $5\nc $6\nc $7\nc $8\nc $9\nc $0\nc $1 $2 $3\nc $1 $2 $3 $4\nc $1 $2 $3 $4 $5\nc $1 $2 $3 $4 $5 $6\nc $2 $0 $1 $8\nc $2 $0 $1 $9\nc $2 $0 $2 $0\nc $2 $0 $2 $1\nc $2 $0 $2 $2\nc $!\nc $@\nc $#\nc $$\nc $! $@\nc $! $@ $#\nc $! $@ $# $$\nc $1 $2 $3 $!\nc $! $1 $2 $3\nc $1 $@ !#\nc $! $@ 1#\nc $2 $0 $1 $8 $!\nc $2 $0 $1 $9 $!\nc $2 $0 $2 $0 $!\nc $2 $0 $2 $1 $!\nc $2 $0 $2 $2 $!\nc $! $2 $0 $1 $8\nc $! $2 $0 $1 $9\nc $! $2 $0 $2 $0\nc $! $2 $0 $2 $1\nc $! $2 $0 $2 $2\nc $2 $0 $1 $8 $! $@ $#\nc $2 $0 $1 $9 $! $@ $#\nc $2 $0 $2 $0 $! $@ $#\nc $2 $0 $2 $1 $! $@ $#\nc $2 $0 $2 $2 $! $@ $#\nc $0 $1 $! \nc $0 $2 $!\nc $0 $3 $!\nc $0 $4 $!\nc $0 $5 $!\nc $0 $6 $!\nc $0 $7 $!\nc $0 $8 $!\nc $0 $9 $!\nc $1 $0 $!\nc $1 $1 $!\nc $1 $2 $!\nc $1 $3 $!\nc $1 $4 $!\nc $1 $5 $!\nc $1 $6 $!\nc $1 $7 $!\nc $1 $8 $!\nc $1 $9 $!\nc $2 $0 $!\nc $2 $1 $!\nc $2 $2 $!\nc $2 $3 $!\nc $2 $4 $!\nc $2 $5 $!\nc $2 $6 $!\nc $2 $7 $!\nc $2 $8 $!\nc $2 $9 $!\nc $3 $0 $!\nc $3 $1 $!\nc $0 $1\nc $0 $2\nc $0 $3\nc $0 $4\nc $0 $5\nc $0 $6\nc $0 $7\nc $0 $8\nc $0 $9\nc $1 $0\nc $1 $1\nc $1 $2\nc $1 $3\nc $1 $4\nc $1 $5\nc $1 $6\nc $1 $7\nc $1 $8\nc $1 $9\nc $2 $0\nc $2 $1\nc $2 $2\nc $2 $3\nc $2 $4\nc $2 $5\nc $2 $6\nc $2 $7\nc $2 $8\nc $2 $9\nc $3 $0\nc $3 $1"

var rules = `l
^/ ^. ^/ #CVE-2010-3863
^/ ^. ^. ^/ ^s ^j  #CVE-2014-0074
^/ ^; ^. ^. ^/ ^s ^j  #CVE-2020-1957
^;
^; ^/ ^/
^/ ^a ^; ^/
$% $2 $f $1
$% $2 $f $. $j $s
$/ $% $3 $b $1
$/ $% $3 $b $. $j $s
$% $3 $b $1
$% $3 $b $. $j $s
$% $0 $d $% $0 $a
$/ $a $% $2 $5 $% $3 $2 $% $6 $6 #CVE-2020-11989
$/ $% $3 $b $a #CVE-2020-13933
#CVE-2020-17510
$/ $.
$/ $% $2 $e
$/ $% $2 $e $/
$/ $% $2 $e $% $2 $e
$/ $% $2 $e $% $2 $e $/

$/ $% $2 $0 #CVE-2020-17523
# i1a i10 i1% #CVE-2022-32532

# java common
$. $a
$. $j $s
$/ $~
$#
$; $/
$/ $. $.
$/ $. $/
$; $. $j $s
$/ $/
$. $. $; $/
^/ ^. ^.
^/ ^/
^/ ^/ ^; ^/
^/ ^f ^2 ^% ^/
^/ ^; ^/
^/ ^; ^. ^/
^/ ^e ^2 ^% ^/
`

func TestProcessRule(t *testing.T) {
	word := "123444"
	toks := Tokens{
		toks: []Token{
			newToken(TOKEN_IDENTIFIER, 'x'),
			newToken(TOKEN_NUMBER, '0'),
			newToken(TOKEN_NUMBER, '4'),
		},
	}
	println(ProcessFunction(word, toks))
}

func TestEval(t *testing.T) {
	word := "adm1"
	input := "$! $@"
	//fmt.Printf("Input = %s\n", input)
	program := Compile(input, "")
	ss := Eval(program, word)
	for i, s := range ss {
		fmt.Printf("%s : %s\n", program.Expressions[i].String(), s)
	}
}

func TestRun(t *testing.T) {
	word := "admin"
	input := rules
	ss, err := RunWithString(input, word)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i, s := range ss {
		fmt.Printf("%d : %s\n", i, s)
	}
}

func TestLexer(t *testing.T) {
	input := rules
	//fmt.Printf("Input = %s\n", input)

	l := NewLexer(input)
	for {
		tok := l.NextToken()
		if tok.Type == TOKEN_EOF {
			break
		}
		println(tok.String())
	}
}

func TestLine(t *testing.T) {
	input := "1@\n##special chars + years\n$2 $0 $1 $8 $!"
	//fmt.Printf("Input = %s\n", input)
	l := NewLexer(input)
	p := NewParser(l)
	for {
		tokens := p.nextLine()
		fmt.Println(tokens)
		if p.peekToken.Type == TOKEN_EOF {
			break
		}
	}
}

func TestParser(t *testing.T) {
	input := rules
	//fmt.Printf("Input = %s\n", input)
	l := NewLexer(input)
	p := NewParser(l)
	programs := p.ParseProgram(nil)
	for _, expr := range programs.Expressions {
		fmt.Println(expr.TokenLiteral())
	}
}
