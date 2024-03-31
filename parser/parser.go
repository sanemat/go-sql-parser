package parser

import (
	"fmt"

	"github.com/sanemat/go-sql-parser/tokens"
)

// Parser holds the state of the parser.
type Parser struct {
	tokens []tokens.Token
	pos    int // current position in the token slice
}

// NewParser creates a new Parser instance.
func NewParser(tokens []tokens.Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parse starts the parsing process and returns the ASTs
func (p *Parser) Parse() ([]Node, error) {
	var nodes []Node
	for p.peek().Type != tokens.TokenEOF {
		node, err := p.parseStatement()
		if err != nil {
			return nil, fmt.Errorf("parseStatement, err: %w", err)
		}
		nodes = append(nodes, node)

		if p.peek().Type == tokens.TokenSemicolon {
			p.pos++ // Advance past the semicolon only if it's present
		}
	}
	return nodes, nil
}

func (p *Parser) peek() tokens.Token {
	if p.pos >= len(p.tokens) {
		// Return an EOF token if we're at or beyond the end of the tokens slice
		return tokens.Token{Type: tokens.TokenEOF, Literal: ""}
	}
	return p.tokens[p.pos]
}
