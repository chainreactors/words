package mask

import (
	"fmt"
	"github.com/chainreactors/words/mask/eval"
	"github.com/chainreactors/words/mask/lexer"
	"github.com/chainreactors/words/mask/parser"
	"github.com/chainreactors/words/mask/token"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	input := "test{?lu#3}"
	fmt.Printf("Input = %s\n", input)

	l := lexer.NewLexer(input)
	for {
		tok := l.NextToken()
		fmt.Printf("Type: %s, Literal = %s\n", tok.Type, tok.Literal)
		if tok.Type == token.TOKEN_EOF {
			break
		}
	}
}

func TestParser(t *testing.T) {
	input := "test{?lu#3}"
	expected := "test{?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#3}"
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
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
		{"test{?lu#3}", "test{?abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ#3}"},
		//{"1 + 2", "3"},
		//{"2 + (3 * 4) / ( 6 - 3 ) + 10", "16"},
		//{"2 + 3 * 4 / 6 - 3  + 10", "11"},
		//{"(5 + 2) * (4 - 2) + 6", "20"},
		//{"5 + 2 * 4 - 2 + 6", "17"},
		//{"5 + 2.1 * 4 - 2 + 6.2", "17.6"},
		//{"2 + 2 ** 2 ** 3", "258"},
		//{"10", "10"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := parser.NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Println(err)
			}
			break
		}

		evaluated := eval.Eval(program)
		if evaluated != nil {
			if evaluated.Inspect() != tt.expected {
				fmt.Printf("%s\n", evaluated.Inspect())
			} else {
				fmt.Printf("%s = %s\n", tt.input, tt.expected)
			}
		}
	}
}

func TestGenerator(t *testing.T) {
	gen := eval.NewGenerator([]string{"a", "b", "c", "d"}, 3)
	println(gen.Strings)
}
