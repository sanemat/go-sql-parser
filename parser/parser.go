package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sanemat/go-sql-parser/tokens"
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
	tokens []tokens.Token
	pos    int // current position in the token slice
}

// NewParser creates a new Parser instance
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

// parseStatement is the entry point for parsing a single SQL statement.
func (p *Parser) parseStatement() (Node, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}
	token := p.tokens[p.pos]
	// Check the first significant token to determine the statement type
	switch {
	case token.Type == tokens.TokenSelect:
		return p.parseSelect()
	case strings.ToUpper(token.Literal) == "INSERT":
		// return p.parseInsert() // Assuming you have a parseInsert method
		return nil, fmt.Errorf("parseInsert is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported statement type: %v, literal: %s, at position %d", token.Type, token.Literal, p.pos)
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
	if token.Type == tokens.TokenFrom {
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
	for !(p.tokens[p.pos].Type == tokens.TokenFrom || p.tokens[p.pos].Type == tokens.TokenSemicolon) {
		expr, err := p.parseSelectExpression()
		if err != nil {
			return nil, fmt.Errorf("parseSelectExpression: %w", err)
		}
		expressions = append(expressions, expr)

		// If next token is a comma, skip it
		nextToken := p.tokens[p.pos]
		if nextToken.Type == tokens.TokenComma {
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
	case tokens.TokenIdentifier:
		// This could be a column name or the beginning of a function call
		p.pos++
		return &ColumnExpression{Name: token.Literal}, nil
	case tokens.TokenNumericLiteral:
		p.pos++
		numFloat, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("parseSelectExpression strconv.ParseFloat num %s, err: %w", token.Literal, err)
		}
		return &NumericLiteral{Value: numFloat}, nil
	case tokens.TokenNull:
		p.pos++
		return &NullValue{}, nil
	case tokens.TokenBooleanLiteral:
		p.pos++
		bl, err := strconv.ParseBool(strings.ToLower(token.Literal))
		if err != nil {
			return nil, fmt.Errorf("parseSelectExpression strconv.ParseBool str %s, err: %w", token.Literal, err)
		}
		return &BooleanLiteral{Value: bl}, nil
	default:
		return nil, fmt.Errorf("unexpected token %s", token.Literal)
	}
}

func (p *Parser) parseSelectTableName() (string, error) {
	if p.tokens[p.pos].Type != tokens.TokenIdentifier {
		return "", fmt.Errorf("unexpected token %s", p.tokens[p.pos].Literal)
	}
	return p.tokens[p.pos].Literal, nil
}
