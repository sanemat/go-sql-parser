package lexer

import (
	"strings"

	"github.com/sanemat/go-sql-parser/tokens"
)

// stateFn represents the state of the lexer as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// lexText is the entry point for the lexing process. It determines the next state based on the current character.
func lexText(l *Lexer) stateFn {
	for {
		switch {
		case l.peek() == '-' && l.peekAhead(1) == '-':
			return lexComment
		case isLetter(l.peek()):
			return lexIdentifier
		case isDigit(l.peek()):
			return lexNumeric
		case isWhitespace(l.peek()):
			return lexWhitespace
		case l.peek() == '\'': // Handle string literals
			return lexString
		case l.peek() == eof:
			l.emit(tokens.TokenEOF)
			return nil
		case isSymbol(l.peek()):
			return lexSymbol
		default:
			l.emit(tokens.TokenError)
			return nil
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
		l.emit(tokens.TokenIdentifier)
	}
	return lexText
}

// lexNumeric scans a numeric identifier.
func lexNumeric(l *Lexer) stateFn {
	seenDecimal := false
	for {
		r := l.peek()
		if isDigit(r) {
			l.next()
		} else if r == '.' && !seenDecimal {
			seenDecimal = true
			l.next()
		} else {
			break
		}
	}
	num := l.input[l.start:l.position]
	l.emitToken(tokens.TokenNumericLiteral, num)
	return lexText
}

// lexWhitespace consumes a run of whitespace characters.
func lexWhitespace(l *Lexer) stateFn {
	for isWhitespace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexText
}

// lexSymbol handles the lexing of symbols.
func lexSymbol(l *Lexer) stateFn {
	symbol := l.next()
	if tokenType, exists := symbols[symbol]; exists {
		l.emitToken(tokenType, string(symbol))
	}
	return lexText
}

// lexComment captures the entire line of a comment and emits it as a TokenComment.
func lexComment(l *Lexer) stateFn {
	l.position += 2 // Skip "--"
	for {
		r := l.next()
		if r == '\n' || r == eof {
			break
		}
	}
	commentText := l.input[l.start:l.position]
	l.emitToken(tokens.TokenComment, commentText)
	return lexText
}

// lexString scans a string literal enclosed in single quotes.
func lexString(l *Lexer) stateFn {
	l.next() // Skip the initial single quote
	for {
		switch r := l.peek(); {
		case r == eof:
			// EOF before closing quote
			l.emit(tokens.TokenError)
			return nil
		case r == '\'' && l.peekAhead(1) != '\'':
			// End of string literal
			l.next() // Consume the closing quote
			stringText := l.input[l.start:l.position]
			l.emitToken(tokens.TokenStringLiteral, stringText)
			return lexText
		case r == '\'' && l.peekAhead(1) == '\'':
			// Escaped quote
			l.next() // Consume the first quote
			l.next() // Consume the second quote
		default:
			l.next() // Consume the next character
		}
	}
}
