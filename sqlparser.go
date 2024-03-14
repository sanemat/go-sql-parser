package sqlparser

import "strings"

// TokenType defines the type of tokens
type TokenType int

const (
	TokenIdentifier TokenType = iota // default type
	TokenKeyword                     // reserved keywords
	// Add more token types as needed
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
}

// keywords maps reserved SQL keywords to their token type
var keywords = map[string]TokenType{
	"SELECT": TokenKeyword,
	"FROM":   TokenKeyword,
	// Add all reserved SQL keywords here
}

// isKeyword checks if a given identifier is a reserved keyword
func isKeyword(identifier string) TokenType {
	if tokType, ok := keywords[strings.ToUpper(identifier)]; ok {
		return tokType
	}
	return TokenIdentifier
}

// tokenize simulates tokenizing an input SQL query (simplified)
func tokenize(input string) []Token {
	// This example splits input by spaces, a real lexer would be more sophisticated
	words := strings.Split(input, " ")
	tokens := make([]Token, 0, len(words))

	for _, word := range words {
		tokType := isKeyword(word)
		tokens = append(tokens, Token{Type: tokType, Literal: word})
	}

	return tokens
}
