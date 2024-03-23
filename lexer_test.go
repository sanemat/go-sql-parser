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
				{Type: TokenSelect, Literal: "select"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"select upper case",
			"SELECT",
			[]Token{
				{Type: TokenSelect, Literal: "SELECT"},
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
				{Type: TokenSelect, Literal: "select"},
				{Type: TokenSymbol, Literal: "*"},
				{Type: TokenFrom, Literal: "from"},
				{Type: TokenIdentifier, Literal: "tablename"},
				{Type: TokenSemicolon, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"invalid syntax",
			"select # invalid syntax;",
			[]Token{
				{Type: TokenSelect, Literal: "select"},
				{Type: TokenError, Literal: ""},
			},
		},
		{
			"multiple columns",
			"select id, title from table1;",
			[]Token{
				{Type: TokenSelect, Literal: "select"},
				{Type: TokenIdentifier, Literal: "id"},
				{Type: TokenComma, Literal: ","},
				{Type: TokenIdentifier, Literal: "title"},
				{Type: TokenFrom, Literal: "from"},
				{Type: TokenIdentifier, Literal: "table1"},
				{Type: TokenSemicolon, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"null value",
			"null",
			[]Token{
				{Type: TokenNull, Literal: "null"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"true value",
			"true",
			[]Token{
				{Type: TokenBooleanLiteral, Literal: "true"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"false value",
			"false",
			[]Token{
				{Type: TokenBooleanLiteral, Literal: "false"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"number",
			"123",
			[]Token{
				{Type: TokenNumericLiteral, Literal: "123"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"big number",
			"1230000000000000000000000000000000000000000",
			[]Token{
				{Type: TokenNumericLiteral, Literal: "1230000000000000000000000000000000000000000"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"decimal points",
			"1.23",
			[]Token{
				{Type: TokenNumericLiteral, Literal: "1.23"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		{
			"multiple statements",
			"select 1; select 2;",
			[]Token{
				{Type: TokenSelect, Literal: "select"},
				{Type: TokenNumericLiteral, Literal: "1"},
				{Type: TokenSemicolon, Literal: ";"},
				{Type: TokenSelect, Literal: "select"},
				{Type: TokenNumericLiteral, Literal: "2"},
				{Type: TokenSemicolon, Literal: ";"},
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
