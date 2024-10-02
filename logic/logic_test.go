package logic

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	inputs := []string{
		"(true || b) && c",
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

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"a && true", true},
		{"false", false},
		{"true && false", false},
		{"true || false", true},
		{"true && true", true},
		{"false || false", false},
		{"true && (true || false)", true},
		{"false || (true && false)", false},
		{"true && (false || true)", true},
		{"true && (false || true) && false", false},
		{"(true || false) && true", true},
		{"(true || false) && false", false},
		{"(true || false) && (true && false)", false},
		{"(true || false) && (true || false)", true},
		{"(true || false) && (true || false) && false", false},
		{"(true || false) && (true || false) && true", true},
		{"(true || false) && (true || false) && (true && false)", false},
		{"(true || false) && (true || false) && (true || false)", true},
		{"(true || false) && (true || false) && (true || false) && false", false},
		{"(true || false) && (true || false) && (true || false) && true", true},
		{"(true || false) && (true || false) && (true || false) && (true && false)", false},
		{"(true || false) && (true || false) && (true || false) && (true || false)", true},
		{"(true || false) && (true || false) && (true || false) && (true || false) && false", false},
		{"(true || false) && (true || false) && (true || false) && (true || false) && true", true},
		{"(true || false) && (true || false) && (true || false) && (true || false) && (true && false)", false},
		{"(true || false) && (true || false) && (true || false) && (true || false) && (true || false)", true},
	}

	for _, tt := range tests {
		l := NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Println(err)
			}
			break
		}
		env := map[string]bool{"a": false}
		evaluated := EvalLogic(program, env)
		fmt.Println(tt.input, "result: ", evaluated)
	}

}
