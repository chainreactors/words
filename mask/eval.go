package mask

import (
	"fmt"
)

func Eval(node Node) (val Object) {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node)
	case *MaskExpression:
		return evalMask(node)
	}

	return nil
}

func evalProgram(program *Program) *GENERATOR {
	var results *GENERATOR
	for _, expr := range program.Expressions {
		var g *GENERATOR
		switch expr.(type) {
		case *Identifier:
			g = NewGeneratorSingle(expr.(*Identifier).String())
		case *MaskExpression:
			g = evalMask(expr.(*MaskExpression)).(*GENERATOR)
		default:
		}
		if results == nil {
			results = g
		} else {
			results.Cross(g)
		}
	}

	return results
}

func evalMask(n *MaskExpression) Object {
	if n.Start.Literal == "?" {
		return NewGenerator(n.CharacterSet, n.Repeat, false)
	} else {
		return NewGenerator(n.CharacterSet, n.Repeat, true)
	}
}

func Run(code string, params [][]string, keywords map[string][]string) ([]string, error) {
	l := NewLexer(code)
	p := NewParser(l, params, keywords)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		return nil, fmt.Errorf("compile error")
	}

	return Eval(program).(*GENERATOR).All(), nil
}

func RunToStream(code string, params [][]string, keywords map[string][]string) (chan string, error) {
	l := NewLexer(code)
	p := NewParser(l, params, keywords)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		return nil, fmt.Errorf("compile error")
	}

	return Eval(program).(*GENERATOR).Streamer, nil
}
