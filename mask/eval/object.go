package eval

import (
	"fmt"
)

//object
type ObjectType string

const (
	NUMBER_OBJ    = "NUMBER"
	GENERATOR_OBJ = "GENERATOR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Number struct {
	Value int64
}

func (n *Number) Inspect() string {
	return fmt.Sprintf("%g", n.Value)
}

func (n *Number) Type() ObjectType { return NUMBER_OBJ }

func NewNumber(f int64) *Number {
	return &Number{Value: f}
}
