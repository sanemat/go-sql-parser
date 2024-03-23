package parser

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/sanemat/go-sql-parser/tokens"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		input   []tokens.Token
		want    Node
		wantErr error
	}{
		{
			name: "simple select",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenIdentifier, Literal: "column1"},
				{Type: tokens.TokenFrom, Literal: "from"},
				{Type: tokens.TokenIdentifier, Literal: "tablea"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&ColumnExpression{Name: "column1"},
					},
					Table: addr("tablea"),
					Where: nil,
				},
			},
		},
		{
			name: "select multiple columns",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenIdentifier, Literal: "id"},
				{Type: tokens.TokenComma, Literal: ","},
				{Type: tokens.TokenIdentifier, Literal: "title"},
				{Type: tokens.TokenFrom, Literal: "from"},
				{Type: tokens.TokenIdentifier, Literal: "table1"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&ColumnExpression{Name: "id"},
						&ColumnExpression{Name: "title"},
					},
					Table: addr("table1"),
					Where: nil,
				},
			},
		},
		{
			name: "select 1;",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenNumericLiteral, Literal: "1"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&NumericLiteral{Value: 1},
					},
					Table: nil,
					Where: nil,
				},
			},
		},
		{
			name: "select null",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenNull, Literal: "null"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&NullValue{},
					},
					Table: nil,
					Where: nil,
				},
			},
		},
		{
			name: "select bool",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenBooleanLiteral, Literal: "false"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&BooleanLiteral{Value: false},
					},
					Table: nil,
					Where: nil,
				},
			},
		},
		{
			name: "select bool case insensitive",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenBooleanLiteral, Literal: "tRue"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&BooleanLiteral{Value: true},
					},
					Table: nil,
					Where: nil,
				},
			},
		},
		{
			name: "multiple statements",
			input: []tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenNumericLiteral, Literal: "1"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenNumericLiteral, Literal: "2"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
			want: []Node{
				&SelectStatement{
					Expressions: []Expression{
						&NumericLiteral{Value: 1},
					},
				},
				&SelectStatement{
					Expressions: []Expression{
						&NumericLiteral{Value: 2},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)
			got, err := p.Parse()
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Praser.Parse() error = %v, at position %d, want %v", err, p.pos, tt.want)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parser.Parse() = %v, want %v", got, tt.want)
				}
			} else {
				if err == nil {
					t.Errorf("Parser.Parse() = %v, wantErr %v", got, tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Parser.Parse() error = %v, at position %d, wantErr %v", err, p.pos, tt.wantErr)
				}
			}
		})
	}
}

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
