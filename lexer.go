package sqlparser

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// keywords defines SQL keywords to be recognized.
var keywords = map[string]TokenType{
	"SELECT": TokenSelect,
	"FROM":   TokenKeyword,
	"NULL":   TokenNull,
	"TRUE":   TokenBooleanLiteral,
	"FALSE":  TokenBooleanLiteral,
	// Add more SQL keywords here...
}

var symbols = map[rune]TokenType{
	';': TokenSymbol,
	',': TokenSymbol,
	'(': TokenSymbol,
	')': TokenSymbol,
	'=': TokenSymbol,
	'*': TokenSymbol,
	'+': TokenSymbol,
	'-': TokenSymbol,
	'/': TokenSymbol,
	// Add more symbols as needed.
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

// lexText is the entry point for the lexing process. It determines the next state based on the current character.
func lexText(l *Lexer) stateFn {
	for {
		switch {
		case l.peek() == '-' && l.peekAhead(1) == '-':
			// Detected start of a comment
			if l.position > l.start {
				l.emit(TokenIdentifier) // Emit any pending identifier before the comment
			}
			return lexComment // Transition to comment handling
		case isLetter(l.peek()):
			// If we encounter a letter and we're not currently processing a token,
			// it's time to handle an identifier or keyword.
			if l.position > l.start {
				l.emit(TokenIdentifier) // This may not be necessary depending on your overall design.
			}
			return lexIdentifier
		case isDigit(l.peek()):
			if l.position > l.start {
				l.emit(TokenIdentifier) // This may not be necessary depending on your overall design.
			}
			return lexNumeric
		case isWhitespace(l.peek()):
			if l.position > l.start {
				l.emit(TokenIdentifier) // Emit identifier if whitespace follows directly after
			}
			return lexWhitespace
		case l.peek() == eof:
			// End of file; if there's any unemitted token, emit it as an identifier.
			if l.position > l.start {
				l.emit(TokenIdentifier)
			}
			l.emit(TokenEOF) // Emit EOF token
			return nil       // No next state, lexing is complete
		case isSymbol(l.peek()):
			if l.position > l.start {
				l.emit(TokenIdentifier)
			}
			return lexSymbol // Transition to a symbol handling state
			// Handle other cases...
		default:
			// The lexer encounters a character that does not match any known patterns
			l.emit(TokenError)
			return nil // Transition to an error state or halt lexing.
		}
	}
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
// Each state function processes part of the input string and transitions the lexer to the next appropriate state.
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

func lexSymbol(l *Lexer) stateFn {
	symbol := l.next()
	if tokenType, exists := symbols[symbol]; exists {
		l.emitToken(tokenType, string(symbol))
	}
	return lexText
}

func isSymbol(r rune) bool {
	_, exists := symbols[r]
	return exists
}

// lexComment captures the entire line of a comment and emits it as a TokenComment.
func lexComment(l *Lexer) stateFn {
	l.position += 2              // Skip the initial "--"
	startOfComment := l.position // Start capturing the comment text after "--"

	// Consume characters until the end of the line or file
	for {
		r := l.next()
		if r == '\n' || r == eof {
			break // Stop at the end of the line or file
		}
	}

	// Capture the entire comment text, including spaces
	commentText := l.input[startOfComment:l.position]
	// Optionally trim whitespace from the start and end if desired
	// commentText = strings.TrimSpace(commentText)

	// Emit the captured comment as a single token
	l.emitToken(TokenComment, commentText)

	if l.peek() == eof {
		l.emit(TokenEOF) // Ensure EOF is emitted if at end of file
	}
	return lexText // Return to the main lexing function
}

// emitToken is a helper to emit tokens with specific literals, simplifying token emission
func (l *Lexer) emitToken(t TokenType, literal string) {
	l.tokens = append(l.tokens, Token{Type: t, Literal: literal})
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

func lexWhitespace(l *Lexer) stateFn {
	for isWhitespace(l.peek()) {
		l.next()
	}
	l.ignore() // Ignore whitespace span
	return lexText
}

// Helper functions to classify characters
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

// lexNumeric scans a numeric identifier.
func lexNumeric(l *Lexer) stateFn {
	// Initial state allows digits and a single decimal point.
	seenDecimal := false
	for {
		r := l.peek()
		if isDigit(r) {
			l.next()
		} else if r == '.' && !seenDecimal {
			// Ensure only a single decimal point is read.
			seenDecimal = true
			l.next()
		} else {
			break // Exit on any non-numeric character
		}
	}

	num := l.input[l.start:l.position]
	l.emitToken(TokenNumericLiteral, num)
	return lexText
}
