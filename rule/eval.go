package rule

import "fmt"

func Eval(node Node, word string) []string {
	switch n := node.(type) {
	case *Program:
		return evalProgram(n, word)
	case *RuleExpression:
		return []string{evalRuleExpression(n, word)}
	}
	return nil
}

func evalProgram(program *Program, word string) []string {
	ss := make([]string, len(program.Expressions))
	for i, expr := range program.Expressions {
		fmt.Println(expr.TokenLiteral())
		ss[i] = evalRuleExpression(expr.(*RuleExpression), word)
	}
	return ss
}

func evalRuleExpression(r *RuleExpression, word string) string {
	var err error
	for _, f := range r.Functions {
		word, err = evalFunctionExpression(&f, word)
		if err != nil {
			panic(r.TokenLiteral() + ", token: " + err.Error())
		}
	}
	return word
}

func evalRuleExpressionSkipError(r *RuleExpression, word string) string {
	var err error
	for _, f := range r.Functions {
		word, err = evalFunctionExpression(&f, word)
		if err != nil {
			return word
		}
	}
	return word
}

func evalFunctionExpression(f *FunctionExpression, word string) (string, error) {
	if !f.IsValid() {
		return "", fmt.Errorf("%s is illegel, %s", f.TokenLiteral(), f.String())
	}
	return ProcessFunction(word, f.Tokens()), nil
}

func Run(rules, word string) (ss []string, evalErr error) {
	defer func() {
		if err := recover(); err != nil {
			evalErr = fmt.Errorf("run error: %v", err)
		}
	}()
	l := NewLexer(rules)
	p := NewParser(l)
	programs := p.ParseProgram()
	return Eval(programs, word), evalErr
}

func RunSkipError(rules, word string) []string {
	l := NewLexer(rules)
	p := NewParser(l)
	programs := p.ParseProgram()
	ss := make([]string, len(programs.Expressions))
	for i, expr := range programs.Expressions {
		fmt.Println(expr.TokenLiteral())
		ss[i] = evalRuleExpressionSkipError(expr.(*RuleExpression), word)
	}
	return ss
}
