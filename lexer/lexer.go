package lexer

import (
	"unicode/utf8"

	"github.com/sanemat/go-sql-parser/tokens"
)

const eof = -1

// keywords defines SQL keywords to be recognized.
var keywords = map[string]tokens.TokenType{
	"SELECT": tokens.TokenSelect,
	"FROM":   tokens.TokenFrom,
	"NULL":   tokens.TokenNull,
	"TRUE":   tokens.TokenBooleanLiteral,
	"FALSE":  tokens.TokenBooleanLiteral,
	// Add more SQL keywords here...
}

var symbols = map[rune]tokens.TokenType{
	';': tokens.TokenSemicolon,
	',': tokens.TokenComma,
	'(': tokens.TokenSymbol,
	')': tokens.TokenSymbol,
	'=': tokens.TokenSymbol,
	'*': tokens.TokenSymbol,
	'+': tokens.TokenSymbol,
	'-': tokens.TokenSymbol,
	'/': tokens.TokenSymbol,
	// Add more symbols as needed.
}

// Lexer holds the state of the scanner.
type Lexer struct {
	input                  string         // Input string being scanned.
	start, position, width int            // Start position of this item, current position, and width of last rune.
	tokens                 []tokens.Token // Slice of tokens identified.
}

// NewLexer returns a new instance of Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Lex scans the input string and produces a slice of tokens.
func (l *Lexer) Lex() []tokens.Token {
	for state := lexText; state != nil; {
		state = state(l)
	}
	return l.tokens
}

func (l *Lexer) emit(t tokens.TokenType) {
	l.tokens = append(l.tokens, tokens.Token{Type: t, Literal: l.input[l.start:l.position]})
	l.start = l.position
}

func (l *Lexer) next() rune {
	if l.position >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.position:])
	l.width = w
	l.position += l.width
	return r
}

func (l *Lexer) ignore() {
	l.start = l.position
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) backup() {
	l.position -= l.width
}

func isSymbol(r rune) bool {
	_, exists := symbols[r]
	return exists
}

// emitToken is a helper to emit tokens with specific literals, simplifying token emission
func (l *Lexer) emitToken(t tokens.TokenType, literal string) {
	l.tokens = append(l.tokens, tokens.Token{Type: t, Literal: literal})
	l.start = l.position // Reset the start position for the next token
}

// peekAhead looks ahead 'n' runes in the input without changing the lexer's position, where 'n' is a positive integer.
// This function is useful for lookahead scenarios, such as detecting comment starts.
func (l *Lexer) peekAhead(n int) rune {
	pos := l.position
	for i := 0; i < n; i++ {
		_, width := utf8.DecodeRuneInString(l.input[pos:])
		pos += width
		if pos >= len(l.input) {
			return eof
		}
	}
	r, _ := utf8.DecodeRuneInString(l.input[pos:])
	return r
}
