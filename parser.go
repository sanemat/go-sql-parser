package sqlparser

import (
	"fmt"
	"strings"
)

// Node represents a node in the Abstract Syntax Tree (AST)
type Node interface{}

// SelectStatement represents a parsed SELECT statement
type SelectStatement struct {
	Columns []string
	Table   string
	Where   *Condition
}

// Condition represents a condition in a WHERE clause
type Condition struct {
	Column   string
	Operator string
	Value    string
}

// Parser holds the state of the parser
type Parser struct {
	tokens []Token
	pos    int // current position in the token slice
}

// NewParser creates a new Parser instance
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parse starts the parsing process and returns the AST
func (p *Parser) Parse() (Node, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	// Check the first significant token to determine the statement type
	switch strings.ToLower(p.tokens[p.pos].Literal) {
	case "select":
		return p.parseSelect()
	case "insert":
		// return p.parseInsert() // Assuming you have a parseInsert method
		return nil, fmt.Errorf("parseInsert is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported statement type: %v", p.tokens[p.pos].Literal)
	}
}

// parseSelect parses a Query rule from the grammar
func (p *Parser) parseSelect() (*SelectStatement, error) {
	// For simplicity, this function assumes the tokens match the expected pattern.
	// In practice, you would check token types and handle errors.
	p.pos++ // Skip the SELECT token
	columnList, err := p.parseSelectColumnList()
	if err != nil {
		return &SelectStatement{}, err
	}
	p.pos++ // Skip the FROM token
	tableName, err := p.parseSelectTableName()
	if err != nil {
		return &SelectStatement{}, err
	}
	// Optionally parse WHERE clause...
	return &SelectStatement{
		Columns: columnList,
		Table:   tableName,
	}, nil
}

// Additional parse functions (parseColumnList, parseTableName, parseWhereClause) need to be implemented
func (p *Parser) parseSelectColumnList() ([]string, error) {
	var columns []string
	// Loop until we reach the "FROM" keyword
	for p.tokens[p.pos].Type != TokenKeyword || strings.ToLower(p.tokens[p.pos].Literal) != "from" {
		if p.tokens[p.pos].Type == TokenIdentifier {
			// Add the column name to the list
			columns = append(columns, p.tokens[p.pos].Literal)
		} else if p.tokens[p.pos].Type == TokenSymbol && p.tokens[p.pos].Literal != "," {
			// If the token is not an identifier or a comma, it's unexpected
			return nil, fmt.Errorf("unexpected token %s", p.tokens[p.pos].Literal)
		}
		p.pos++ // Move to the next token
	}

	if len(columns) == 0 {
		// If no columns were found, return an error
		return nil, fmt.Errorf("no columns found in select statement")
	}
	return columns, nil
}

func (p *Parser) parseSelectTableName() (string, error) {
	if p.tokens[p.pos].Type != TokenIdentifier {
		return "", fmt.Errorf("unexpected token %s", p.tokens[p.pos].Literal)
	}
	return p.tokens[p.pos].Literal, nil
}
