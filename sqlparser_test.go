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

func TestLexIdentifier(t *testing.T) {
	input := "tablename"
	expected := []Token{
		{Type: TokenIdentifier, Literal: "tablename"},
		{Type: TokenEOF, Literal: ""},
	}

	lexer := NewLexer(input)
	tokens := lexer.Lex()

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("TestLexIdentifier: unexpected tokens. expected=%+v, got=%+v", expected, tokens)
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
