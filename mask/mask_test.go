package mask

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLexer(t *testing.T) {
	inputs := []string{
		"test-{@aaa|bbb#3}",
		"test-{@aaa|bbb#3}+{@ccc|ddd#3}",
		"test{{1iohoi",
		"test{aaa}",
	}
	for _, input := range inputs {
		fmt.Printf("Input = %s\n", input)
		l := NewLexer(input)
		for {
			tok := l.NextToken()
			fmt.Printf("Type: %s, Literal = %s\n", tok.Type, tok.Literal)
			if tok.Type == TOKEN_EOF {
				break
			}
		}
	}
}

func TestParser(t *testing.T) {
	dicts := [][]string{
		[]string{"aaa", "bbb", "ccc"},
	}
	input := "test{@month}"
	expected := "test{?a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z#3}"
	l := NewLexer(input)
	p := NewParser(l, dicts, nil)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	println(program.String())
	if program.String() != expected {
		fmt.Printf("Syntax error: expected %s, got %s\n", expected, program.String())
		os.Exit(1)
	}

	fmt.Printf("input  = %s\n", input)
	fmt.Printf("output = %s\n", program.String())
}

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		//{"test{?lu#3}", "test{?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#3}"},
		{"test{?123}", ""},
		//{"1 + 2", "3"},
		//{"2 + (3 * 4) / ( 6 - 3 ) + 10", "16"},
		//{"2 + 3 * 4 / 6 - 3  + 10", "11"},
		//{"(5 + 2) * (4 - 2) + 6", "20"},
		//{"5 + 2 * 4 - 2 + 6", "17"},
		//{"5 + 2.1 * 4 - 2 + 6.2", "17.6"},
		//{"2 + 2 ** 2 ** 3", "258"},
		//{"10", "10"},
	}
	dicts := [][]string{
		[]string{"aaa", "bbb", "ccc"},
	}
	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l, dicts, nil)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Println(err)
			}
			break
		}

		evaluated := Eval(program)
		if evaluated != nil {
			if evaluated.Inspect() != tt.expected {
				fmt.Printf("%s\n", evaluated.Inspect())
			} else {
				fmt.Printf("%s = %s\n", tt.input, tt.expected)
			}
		}
	}
}

func TestRun(t *testing.T) {
	var err error
	words, err := Run("/{?l}", nil, nil)
	fmt.Printf("%v,%v", words, err)
	keywords := map[string][]string{"test": []string{"a", "b", "c", "d"}}
	stream, err := RunToStream("{$l#2}.oocl.com{@test}", nil, keywords)
	for w := range stream {
		fmt.Println(w)
	}
}

func TestProduct(t *testing.T) {
	words := Product(wrapStream([]string{"a", "b", "c", "d"}), []string{"a", "b", "c", "d"})
	for w := range words {
		fmt.Println(w)
	}
}

func TestGenerator(t *testing.T) {
	start := time.Now()
	gen := NewGenerator(ParseCharacterSetWithIdent("l"), 5, false)
	for s := range gen.Streamer {
		s = s
		continue
	}
	println(time.Since(start).String())
	//for s := range gen.Streamer {
	//	fmt.Println(s)
	//}
}

func TestNewGeneratorSingle(t *testing.T) {
	gen := NewGeneratorSingle("a")
	for s := range gen.Streamer {
		fmt.Println(s)
	}
}

func TestCross(t *testing.T) {
	gen1 := NewGenerator([]string{"a", "b", "c", "d"}, 2, true)
	//gen2 := NewGenerator([]string{"eee", "fff"}, 2, true)
	gen3 := NewGeneratorSingle("z")
	gen1.Cross(gen3)
	for w := range gen1.Streamer {
		fmt.Println(w)
	}
}
