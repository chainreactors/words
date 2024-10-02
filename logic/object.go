package logic

import (
	"fmt"
)

// object
type ObjectType string

const (
	//NUMBER_OBJ  = "NUMBER"
	BOOLEAN_OBJ = "BOOLEAN"
)

var (
	TRUE  = &Boolean{Bool: true}
	FALSE = &Boolean{Bool: false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Boolean struct {
	Bool bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%v", b.Bool)
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

func nativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func IsTrue(obj Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		//	switch obj.Type() {
		//	case NUMBER_OBJ:
		//		if obj.(*Number).Value == 0.0 {
		//			return false
		//		}
		//	}
		return false
	}
}

func objectToNativeBoolean(o Object) bool {
	switch obj := o.(type) {
	case *Boolean:
		return obj.Bool
	//case *Number:
	//	if obj.Value == 0.0 {
	//		return false
	//	}
	//	return true
	default:
		return true
	}
}
