package sqlparser

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TokenType defines the type of lexed tokens.
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenIdentifier
	TokenKeyword
	TokenSymbol
	// Extend with more token types as needed (e.g., TokenString, TokenNumber)
)

// Token represents a lexed piece of the input.
type Token struct {
	Type    TokenType
	Literal string
}

// keywords defines SQL keywords to be recognized.
var keywords = map[string]TokenType{
	"SELECT": TokenKeyword,
	"FROM":   TokenKeyword,
	// Add more SQL keywords here...
}

// Lexer holds the state of the scanner.
type Lexer struct {
	input                  string  // Input string being scanned.
	start, position, width int     // Start position of this item, current position, and width of last rune.
	tokens                 []Token // Slice of tokens identified.
}

// NewLexer returns a new instance of Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Lex scans the input string and produces a slice of tokens.
func (l *Lexer) Lex() []Token {
	for state := lexText; state != nil; {
		state = state(l)
	}
	return l.tokens
}

// lexText is the lexer function for general text.
func lexText(l *Lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.position:], " ") {
			if l.position > l.start {
				l.emit(TokenIdentifier)
			}
			l.ignore()
		} else if isLetter(l.peek()) {
			return lexIdentifier
		} else if l.next() == eof {
			break
		}
	}
	// Emit an EOF token when done.
	l.emit(TokenEOF)
	return nil
}

// lexIdentifier scans an alphanumeric identifier.
func lexIdentifier(l *Lexer) stateFn {
	for isLetter(l.peek()) || isDigit(l.peek()) {
		l.next()
	}
	word := l.input[l.start:l.position]
	if tokType, isKeyword := keywords[strings.ToUpper(word)]; isKeyword {
		l.emit(tokType)
	} else {
		l.emit(TokenIdentifier)
	}
	return lexText
}

// stateFn represents the state of the lexer as a function that returns the next state.
type stateFn func(*Lexer) stateFn

func (l *Lexer) emit(t TokenType) {
	l.tokens = append(l.tokens, Token{Type: t, Literal: l.input[l.start:l.position]})
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

const eof = -1

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}
