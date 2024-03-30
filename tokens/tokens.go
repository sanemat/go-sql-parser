package tokens

import (
	"fmt"
	"strings"
)

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
	TokenSemicolon
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

func (t Token) RawValue() string {
	switch t.Type {
	case TokenStringLiteral:
		// Standard SQL single-quoted strings
		if strings.HasPrefix(t.Literal, "'") && strings.HasSuffix(t.Literal, "'") {
			return t.Literal[1 : len(t.Literal)-1]
		}
	}
	return t.Literal
}
