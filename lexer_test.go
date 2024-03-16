package sqlparser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			"select lower case",
			"select",
			[]Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"select upper case",
			"SELECT",
			[]Token{
				{Type: TokenKeyword, Literal: "SELECT"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"table name",
			"tablename",
			[]Token{
				{Type: TokenIdentifier, Literal: "tablename"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"asterisk",
			"*",
			[]Token{
				{Type: TokenSymbol, Literal: "*"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"simple select sql",
			"select * from tablename;",
			[]Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenSymbol, Literal: "*"},
				{Type: TokenKeyword, Literal: "from"},
				{Type: TokenIdentifier, Literal: "tablename"},
				{Type: TokenSymbol, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"invalid syntax",
			"select # invalid syntax;",
			[]Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenError, Literal: ""},
			},
		},
		{
			"multiple columns",
			"select id, title from table1;",
			[]Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenIdentifier, Literal: "id"},
				{Type: TokenSymbol, Literal: ","},
				{Type: TokenIdentifier, Literal: "title"},
				{Type: TokenKeyword, Literal: "from"},
				{Type: TokenIdentifier, Literal: "table1"},
				{Type: TokenSymbol, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Lex()

			if !reflect.DeepEqual(tokens, tt.expected) {
				t.Errorf("unexpected tokens. expected=%+v, got=%+v", tt.expected, tokens)
			}
		})
	}
}
