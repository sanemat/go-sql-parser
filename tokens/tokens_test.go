package tokens

import (
	"testing"
)

func TestTokenRawValue(t *testing.T) {
	tests := []struct {
		name     string
		token    Token
		expected string
	}{
		{
			name: "boolean literal",
			token: Token{
				Type:    TokenBooleanLiteral,
				Literal: "true",
			},
			expected: "true",
		},
		{
			name: "string literal",
			token: Token{
				Type:    TokenStringLiteral,
				Literal: "'text'",
			},
			expected: "text",
		},
		{
			name: "string literal with quote",
			token: Token{
				Type:    TokenStringLiteral,
				Literal: "'O''Reilly'",
			},
			expected: "O'Reilly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.RawValue()
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
