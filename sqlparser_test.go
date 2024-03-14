package sqlparser

import (
	"reflect"
	"testing"
)

func TestLexSingleKeyword(t *testing.T) {
	input := "SELECT"
	expected := []Token{
		{Type: TokenKeyword, Literal: "SELECT"},
		{Type: TokenEOF, Literal: ""},
	}

	lexer := NewLexer(input)
	tokens := lexer.Lex()

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("TestLexSingleKeyword: unexpected tokens. expected=%+v, got=%+v", expected, tokens)
	}
}
