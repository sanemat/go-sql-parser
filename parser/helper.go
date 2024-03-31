package parser

import "github.com/sanemat/go-sql-parser/tokens"

// get pointer of value
// https://github.com/golang/go/issues/45624#issuecomment-1843832599
//
//lint:ignore U1000 for testing
func addr[T any](v T) *T { return &v }

// Helper function to convert an operator token type to its string representation
func operatorToString(op tokens.TokenType) string {
	switch op {
	case tokens.TokenGreaterThan:
		return ">"
	case tokens.TokenGreaterThanOrEqual:
		return ">="
	case tokens.TokenLessThan:
		return "<"
	case tokens.TokenLessThanOrEqual:
		return "<="
	// Add more cases as necessary
	default:
		return "unknown_operator"
	}
}

// Helper function to check if a token type is a literal
func isLiteral(tokenType tokens.TokenType) bool {
	return tokenType == tokens.TokenNumericLiteral || tokenType == tokens.TokenStringLiteral ||
		tokenType == tokens.TokenBooleanLiteral || tokenType == tokens.TokenNull
}
