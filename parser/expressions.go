package parser

import (
	"fmt"
)

type ColumnExpression struct {
	Name string
}

func (c *ColumnExpression) String() string {
	return fmt.Sprintf("ColumnExpression(%s)", c.Name)
}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

type NumericLiteral struct {
	Value float64
}

func (n *NumericLiteral) String() string {
	return fmt.Sprintf("NumericLiteral(%f)", n.Value)
}

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral('%s')", s.Value)
}

type NullValue struct{}

func (n *NullValue) String() string {
	return "NullValue(NULL)"
}

type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("BooleanLiteral(%t)", b.Value)
}
