package mask

func Eval(node Node) (val Object) {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node)
	//case *ast.NumberLiteral:
	//	return evalNumber(node)
	case *MaskExpression:
		return evalMask(node)
		//return evalMask(node)
		//case *ast.PrefixExpression:
		//	right := Eval(node.Right)
		//	return evalPrefixExpression(node, right)
		//case *ast.InfixExpression:
		//	left := Eval(node.Left)
		//
		//	right := Eval(node.Right)
		//	return evalInfixExpression(node, left, right)
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

	if results == nil {
		return nil
	}
	return results
}

//func evalNumber(n *ast.NumberLiteral) Object {
//	return NewNumber(n.Value)
//}

func evalMask(n *MaskExpression) Object {
	return NewGenerator(n.CharacterSet, n.Repeat)
}
