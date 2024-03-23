package parser

import (
	"fmt"
	"strings"
)

// Node represents a node in the Abstract Syntax Tree (AST).
type Node interface{}

// Expression interface for expressions in the AST.
type Expression interface {
	String() string
}

// Condition represents a condition in a WHERE clause.
type Condition struct {
	Column   string
	Operator string
	Value    string
}

// SelectStatement represents a parsed SELECT statement.
type SelectStatement struct {
	Expressions []Expression
	Table       *string
	Where       *Condition
}

func (s *SelectStatement) String() string {
	expressions := make([]string, len(s.Expressions))
	for i, expr := range s.Expressions {
		expressions[i] = expr.String()
	}
	tableName := "nil"
	if s.Table != nil {
		tableName = *s.Table
	}
	return fmt.Sprintf(
		"SelectStatement(Expressions: [%s], Table: %s)",
		strings.Join(expressions, ", "), tableName)
}

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
