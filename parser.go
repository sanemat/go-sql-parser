package sqlparser

import (
	"fmt"
	"strconv"
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

type NullValue struct {
}

type BooleanLiteral struct {
	Value bool
}

// SelectStatement represents a parsed SELECT statement
type SelectStatement struct {
	Expressions []Expression
	Table       *string
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
	token := p.tokens[p.pos]
	// Check the first significant token to determine the statement type
	switch {
	case token.Type == TokenSelect:
		return p.parseSelect()
	case strings.ToUpper(token.Literal) == "INSERT":
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

	token := p.tokens[p.pos]
	var tableNamePtr *string
	if token.Type == TokenFrom {
		p.pos++ // Skip the FROM token
		tableName, err := p.parseSelectTableName()
		if err != nil {
			return &SelectStatement{}, err
		}
		tableNamePtr = &tableName
	}
	// Optionally parse WHERE clause...
	return &SelectStatement{
		Expressions: expressions,
		Table:       tableNamePtr,
	}, nil
}

func (p *Parser) parseSelectExpressions() ([]Expression, error) {
	var expressions []Expression

	// Continue looping until "FROM" is encountered or until semicolon
	for !(p.tokens[p.pos].Type == TokenFrom || p.tokens[p.pos].Type == TokenSemicolon) {
		expr, err := p.parseSelectExpression()
		if err != nil {
			return nil, fmt.Errorf("parseSelectExpression: %w", err)
		}
		expressions = append(expressions, expr)

		// If next token is a comma, skip it
		nextToken := p.tokens[p.pos]
		if nextToken.Type == TokenComma {
			p.pos++ // Move forward the expression token
		}
	}

	// Handle edge case: No expressions
	if len(expressions) == 0 {
		return nil, fmt.Errorf("no expressions found in select statement")
	}

	return expressions, nil
}

func (p *Parser) parseSelectExpression() (Expression, error) {
	// This is a simplified placeholder. You'll need to replace this with actual logic
	// to parse different types of expressions based on your tokens.
	token := p.tokens[p.pos]
	switch token.Type {
	case TokenIdentifier:
		// This could be a column name or the beginning of a function call
		p.pos++
		return &ColumnExpression{Name: token.Literal}, nil
	case TokenNumericLiteral:
		p.pos++
		numFloat, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("parseSelectExpression strconv.ParseFloat num %s, err: %w", token.Literal, err)
		}
		return &NumericLiteral{Value: numFloat}, nil
	case TokenNull:
		p.pos++
		return &NullValue{}, nil
	case TokenBooleanLiteral:
		p.pos++
		bl, err := strconv.ParseBool(token.Literal)
		if err != nil {
			return nil, fmt.Errorf("parseSelectExpression strconv.ParseBool str %s, err: %w", token.Literal, err)
		}
		return &BooleanLiteral{Value: bl}, nil
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
