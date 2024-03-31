package parser

import (
	"fmt"
	"testing"

	"github.com/sanemat/go-sql-parser/tokens"
)

func TestNodeStringMethods(t *testing.T) {
	tests := []struct {
		name     string
		node     fmt.Stringer
		expected string
	}{
		{
			name: "SelectStatement with expressions",
			node: &SelectStatement{
				Expressions: []Expression{
					&NumericLiteral{Value: 123},
				},
			},
			expected: "SelectStatement(Expressions: [NumericLiteral(123.000000)], Table: nil)",
		},
		{
			name: "SelectStatement with expressions with table",
			node: &SelectStatement{
				Expressions: []Expression{
					&ColumnExpression{Name: "column1"},
					&NumericLiteral{Value: 123},
				},
				Table: addr("table1"),
			},
			expected: "SelectStatement(Expressions: [ColumnExpression(column1), NumericLiteral(123.000000)], Table: table1)",
		},
		{
			name:     "ColumnExpression",
			node:     &ColumnExpression{Name: "column2"},
			expected: "ColumnExpression(column2)",
		},
		{
			name:     "NumericLiteral",
			node:     &NumericLiteral{Value: 456.78},
			expected: "NumericLiteral(456.780000)",
		},
		{
			name:     "StringLiteral",
			node:     &StringLiteral{Value: "test string"},
			expected: "StringLiteral('test string')",
		},
		{
			name:     "NullValue",
			node:     &NullValue{},
			expected: "NullValue(NULL)",
		},
		{
			name:     "BooleanLiteral True",
			node:     &BooleanLiteral{Value: true},
			expected: "BooleanLiteral(true)",
		},
		{
			name:     "BooleanLiteral False",
			node:     &BooleanLiteral{Value: false},
			expected: "BooleanLiteral(false)",
		},
		{
			name: "BinaryExpression",
			node: &BinaryExpression{
				Left:     &NumericLiteral{Value: 123},
				Operator: tokens.TokenGreaterThan,
				Right:    &NumericLiteral{Value: 234},
			},
			expected: "BinaryExpression(NumericLiteral(123.000000) > NumericLiteral(234.000000))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.String()
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
