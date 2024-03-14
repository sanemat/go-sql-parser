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
			"SELECT * FROM table;",
			[]Token{
				{Type: TokenKeyword, Literal: "SELECT"},
				{Type: TokenSymbol, Literal: "*"},
				{Type: TokenKeyword, Literal: "FROM"},
				{Type: TokenIdentifier, Literal: "table"},
				{Type: TokenSymbol, Literal: ";"},
				{Type: TokenEOF, Literal: ""},
			},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		lexer := NewLexer(tt.input)
		tokens := lexer.Lex()

		if !reflect.DeepEqual(tokens, tt.expected) {
			t.Errorf("unexpected tokens. expected=%+v, got=%+v", tt.expected, tokens)
		}
	}
}
