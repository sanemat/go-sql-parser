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
	}

	for i, tt := range tests {
		lexer := NewLexer(tt.input)
		tokens := lexer.Lex()

		if !reflect.DeepEqual(tokens, tt.expected) {
			t.Errorf("Test case %d: unexpected tokens. expected=%+v, got=%+v", i, tt.expected, tokens)
		}
	}
}

// func TestLexSymbol(t *testing.T) {
// 	input := "*"
// 	expected := []Token{
// 		{Type: TokenSymbol, Literal: "*"},
// 		{Type: TokenEOF, Literal: ""},
// 	}

// 	lexer := NewLexer(input)
// 	tokens := lexer.Lex()

// 	if !reflect.DeepEqual(tokens, expected) {
// 		t.Errorf("TestLexSymbol: unexpected tokens. expected=%+v, got=%+v", expected, tokens)
// 	}
// }
