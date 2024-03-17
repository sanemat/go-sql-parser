package sqlparser

import (
	"fmt"
	"strings"
)

// Node represents a node in the Abstract Syntax Tree (AST)
type Node interface{}

type Expression interface {
	// Expression is a marker interface for all expression types
}

type ColumnExpression struct {
	Name string
}

type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

type NumericLiteral struct {
	Value float64
}

type StringLiteral struct {
	Value string
}

// SelectStatement represents a parsed SELECT statement
type SelectStatement struct {
	Expressions []Expression
	Table       string
	Where       *Condition
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
	switch strings.ToUpper(p.tokens[p.pos].Literal) {
	case "SELECT":
		return p.parseSelect()
	case "INSERT":
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
	expressions, err := p.parseSelectExpressions()
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
		Expressions: expressions,
		Table:       tableName,
	}, nil
}

func (p *Parser) parseSelectExpressions() ([]Expression, error) {
	var expressions []Expression

	// Loop until we reach the "FROM" keyword
	for !(p.tokens[p.pos].Type == TokenKeyword && strings.ToUpper(p.tokens[p.pos].Literal) == "FROM") {
		// Parse an expression
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)

		// Check if the next token is a comma. If so, skip it.
		if p.peekToken().Literal == "," {
			p.pos += 2 // Skip the comma and move to the next expression
			continue
		} else if p.peekToken().Type == TokenKeyword && strings.ToUpper(p.peekToken().Literal) == "FROM" {
			p.pos++ // Move past the last expression before the FROM keyword
			break
		} else {
			// Handle unexpected tokens
			return nil, fmt.Errorf("unexpected token %s", p.peekToken().Literal)
		}
	}

	if len(expressions) == 0 {
		// If no expressions were found, return an error.
		return nil, fmt.Errorf("no expressions found in select statement")
	}
	return expressions, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	// This is a simplified placeholder. You'll need to replace this with actual logic
	// to parse different types of expressions based on your tokens.
	token := p.tokens[p.pos]
	switch token.Type {
	case TokenIdentifier:
		// This could be a column name or the beginning of a function call
		p.pos++
		return &ColumnExpression{Name: token.Literal}, nil
	// Add cases for other types of expressions: NumericLiteral, BinaryExpression, etc.
	default:
		return nil, fmt.Errorf("unexpected token %s", token.Literal)
	}
}

func (p *Parser) parseSelectTableName() (string, error) {
	if p.tokens[p.pos].Type != TokenIdentifier {
		return "", fmt.Errorf("unexpected token %s", p.tokens[p.pos].Literal)
	}
	return p.tokens[p.pos].Literal, nil
}

func (p *Parser) peekToken() Token {
	peekPos := p.pos + 1
	if peekPos >= len(p.tokens) {
		return Token{Type: TokenEOF, Literal: ""}
	}
	return p.tokens[peekPos]
}
