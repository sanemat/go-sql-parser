package lexer

import (
	"reflect"
	"testing"

	"github.com/sanemat/go-sql-parser/tokens"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []tokens.Token
	}{
		{
			"select lower case",
			"select",
			[]tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"select upper case",
			"SELECT",
			[]tokens.Token{
				{Type: tokens.TokenSelect, Literal: "SELECT"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"table name",
			"tablename",
			[]tokens.Token{
				{Type: tokens.TokenIdentifier, Literal: "tablename"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"asterisk",
			"*",
			[]tokens.Token{
				{Type: tokens.TokenSymbol, Literal: "*"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"simple select sql",
			"select * from tablename;",
			[]tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenSymbol, Literal: "*"},
				{Type: tokens.TokenFrom, Literal: "from"},
				{Type: tokens.TokenIdentifier, Literal: "tablename"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"invalid syntax",
			"select # invalid syntax;",
			[]tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenError, Literal: ""},
			},
		},
		{
			"multiple columns",
			"select id, title from table1;",
			[]tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenIdentifier, Literal: "id"},
				{Type: tokens.TokenComma, Literal: ","},
				{Type: tokens.TokenIdentifier, Literal: "title"},
				{Type: tokens.TokenFrom, Literal: "from"},
				{Type: tokens.TokenIdentifier, Literal: "table1"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"null value",
			"null",
			[]tokens.Token{
				{Type: tokens.TokenNull, Literal: "null"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"true value",
			"true",
			[]tokens.Token{
				{Type: tokens.TokenBooleanLiteral, Literal: "true"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"false value",
			"false",
			[]tokens.Token{
				{Type: tokens.TokenBooleanLiteral, Literal: "false"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"number",
			"123",
			[]tokens.Token{
				{Type: tokens.TokenNumericLiteral, Literal: "123"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"big number",
			"1230000000000000000000000000000000000000000",
			[]tokens.Token{
				{Type: tokens.TokenNumericLiteral, Literal: "1230000000000000000000000000000000000000000"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"decimal points",
			"1.23",
			[]tokens.Token{
				{Type: tokens.TokenNumericLiteral, Literal: "1.23"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"multiple statements",
			"select 1; select 2;",
			[]tokens.Token{
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenNumericLiteral, Literal: "1"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenSelect, Literal: "select"},
				{Type: tokens.TokenNumericLiteral, Literal: "2"},
				{Type: tokens.TokenSemicolon, Literal: ";"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"string",
			"'text'",
			[]tokens.Token{
				{Type: tokens.TokenStringLiteral, Literal: "'text'"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"string with escaped quote",
			"'O''Reilly'",
			[]tokens.Token{
				{Type: tokens.TokenStringLiteral, Literal: "'O''Reilly'"},
				{Type: tokens.TokenEOF, Literal: ""},
			},
		},
		{
			"unterminated string",
			"'this string has no end",
			[]tokens.Token{
				{Type: tokens.TokenError, Literal: "'this string has no end"},
			},
		},
		{
			"empty string",
			"''",
			[]tokens.Token{
				{Type: tokens.TokenStringLiteral, Literal: "''"},
				{Type: tokens.TokenEOF, Literal: ""},
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
