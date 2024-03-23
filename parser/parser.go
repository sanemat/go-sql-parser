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
	var statements []Node
	for p.pos < len(p.tokens) {
		if p.tokens[p.pos].Type == tokens.TokenEOF {
			break
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, fmt.Errorf("parseStatement, err: %w", err)
		}
		statements = append(statements, stmt)
		if p.tokens[p.pos].Type == tokens.TokenSemicolon {
			p.pos++ // Skip the semicolon.
		}
	}
	return statements, nil
}
