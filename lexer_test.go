package sqlparser

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []Token
	}{
		{
			"select",
			[]Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"SELECT",
			[]Token{
				{Type: TokenKeyword, Literal: "SELECT"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"tablename",
			[]Token{
				{Type: TokenIdentifier, Literal: "tablename"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"*",
			[]Token{
				{Type: TokenSymbol, Literal: "*"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
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
			"select # invalid sintax;",
			[]Token{
				{Type: TokenKeyword, Literal: "select"},
				{Type: TokenError, Literal: ""},
			},
		},
		{
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

	for i, tt := range tests {
		lexer := NewLexer(tt.input)
		tokens := lexer.Lex()

		if !reflect.DeepEqual(tokens, tt.expected) {
			t.Errorf("Test case %d: unexpected tokens. expected=%+v, got=%+v", i, tt.expected, tokens)
		}
	}
}
