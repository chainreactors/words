package logic

func EvalLogic(node Node) bool {
	return IsTrue(Eval(node))
}

func Eval(node Node) (val Object) {
	switch node := node.(type) {
	case *Program:
		return evalProgram(node)
	//case *NumberLiteral:
	//	return evalNumber(node)
	case *PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node, right)
	case *InfixExpression:
		left := Eval(node.Left)

		right := Eval(node.Right)
		return evalInfixExpression(node, left, right)
	case *BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

func evalProgram(program *Program) (results Object) {
	// for _, expr := range program.Expressions {
	// 	results = Eval(expr)
	// }
	results = Eval(program.Expression)
	return results
}

func evalPrefixExpression(node *PrefixExpression, right Object) Object {
	switch node.Operator {
	//case "+":
	//	return evalPlusPrefixOperatorExpression(node, right)
	//case "-":
	//	return evalMinusPrefixOperatorExpression(node, right)
	case "!":
		return evalBangOperatorExpression(node, right)
	default:
		return nil
	}
}

func evalBangOperatorExpression(node *PrefixExpression, right Object) Object {
	return nativeBoolToBooleanObject(!IsTrue(right))
}

func evalInfixExpression(node *InfixExpression, left, right Object) Object {
	operator := node.Operator
	switch {
	case operator == "&&":
		leftCond := objectToNativeBoolean(left)
		if !leftCond {
			return FALSE
		}

		rightCond := objectToNativeBoolean(right)
		return nativeBoolToBooleanObject(leftCond && rightCond)
	case operator == "||":
		leftCond := objectToNativeBoolean(left)
		if leftCond {
			return TRUE
		}

		rightCond := objectToNativeBoolean(right)
		return nativeBoolToBooleanObject(leftCond || rightCond)
	//case left.Type() == NUMBER_OBJ && right.Type() == NUMBER_OBJ:
	//	return evalNumberInfixExpression(node, left, right)
	default:
		return nil
	}
}
