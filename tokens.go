package sqlparser

import "fmt"

// TokenType defines the type of lexed tokens.
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenIdentifier
	TokenKeyword
	TokenSymbol
	TokenComment
	TokenStringLiteral
	TokenNumericLiteral
	TokenDateAndTimeLiteral
	TokenHexadecimalLiteral
	TokenBitValueLiteral
	TokenBooleanLiteral
	TokenNull

	TokenSelect
	TokenFrom
	TokenComma
	// Extend with more token types as needed (e.g., TokenString, TokenNumber)
)

// Token represents a lexed piece of the input.
type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("token Type: %d, Literal: %s", t.Type, t.Literal)
}
