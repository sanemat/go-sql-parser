package lexer

import "unicode"

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// Helper functions to classify characters
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
