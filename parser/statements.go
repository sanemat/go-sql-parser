package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sanemat/go-sql-parser/tokens"
)

// parseStatement is the entry point for parsing a single SQL statement.
func (p *Parser) parseStatement() (Node, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}
	// Check the first significant token to determine the statement type
	switch {
	case p.peek().Type == tokens.TokenSelect:
		return p.parseSelect()
	case strings.ToUpper(p.peek().Literal) == "INSERT":
		// return p.parseInsert() // Assuming you have a parseInsert method
		return nil, fmt.Errorf("parseInsert is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported statement type: %v, literal: %s, at position %d", p.peek().Type, p.peek().Literal, p.pos)
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
	case tokens.TokenStringLiteral:
		// Handle string literals by removing quotes and creating a StringLiteral expression
		p.pos++
		rawValue := token.RawValue()
		return &StringLiteral{Value: rawValue}, nil
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
	// Ensure the current token is an identifier (e.g., table name).
	if p.tokens[p.pos].Type != tokens.TokenIdentifier {
		return "", fmt.Errorf("expected table name, found %s", p.tokens[p.pos].Literal)
	}
	tableName := p.tokens[p.pos].Literal
	p.pos++ // Move past the table name.
	return tableName, nil
}
