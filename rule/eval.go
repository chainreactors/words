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

func Run(rules []Expression, word string) (ss []string, evalErr error) {
	defer func() {
		if err := recover(); err != nil {
			evalErr = fmt.Errorf("run error: %v", err)
		}
	}()
	ss = make([]string, len(rules))
	for i, rule := range rules {
		ss[i] = evalRuleExpression(rule.(*RuleExpression), word)
	}
	return ss, evalErr
}

func RunSkipError(rules []Expression, word string) []string {
	ss := make([]string, len(rules))
	for i, rule := range rules {
		ss[i] = evalRuleExpressionSkipError(rule.(*RuleExpression), word)
	}
	return ss
}

func RunAsStream(rules []Expression, word string) chan string {
	ch := make(chan string)
	go func() {
		for _, expr := range rules {
			ch <- evalRuleExpressionSkipError(expr.(*RuleExpression), word)
		}
		close(ch)
	}()
	return ch
}

func RunWithString(rules, word string) (ss []string, evalErr error) {
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

func Compile(rules string) []Expression {
	l := NewLexer(rules)
	p := NewParser(l)
	return p.ParseProgram().Expressions
}
