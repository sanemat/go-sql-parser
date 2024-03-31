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
	expressions, err := p.parseExpressions()
	if err != nil {
		return &SelectStatement{}, err
	}

	var tableNamePtr *string
	if p.peek().Type == tokens.TokenFrom {
		p.pos++ // Skip the FROM token
		tableName, err := p.parseSelectTableName()
		if err != nil {
			return &SelectStatement{}, err
		}
		tableNamePtr = &tableName
	}
	// Skip until the end of statement or start of the next statement
	for !(p.peek().Type == tokens.TokenSemicolon || p.peek().Type == tokens.TokenEOF) {
		p.pos++
	}
	// Optionally parse WHERE clause...
	return &SelectStatement{
		Expressions: expressions,
		Table:       tableNamePtr,
	}, nil
}

func (p *Parser) parseExpressions() ([]Expression, error) {
	var expressions []Expression

	for {
		token := p.next()

		var expr Expression
		var err error
		switch {
		case token.Type == tokens.TokenIdentifier:
			expr = &ColumnExpression{Name: token.Literal}
		case isLiteral(token.Type):
			expr, err = p.parseLiteral(token)
		default:
			err = fmt.Errorf("unexpected token in expression: %v", token.Literal)
		}
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)

		if p.peek().Type != tokens.TokenComma {
			break
		}
		p.pos++ // Skip the comma
	}
	return expressions, nil
}

func (p *Parser) parseLiteral(token tokens.Token) (Expression, error) {
	switch token.Type {
	case tokens.TokenNumericLiteral:
		value, err := strconv.ParseFloat(token.Literal, 64)
		if err != nil {
			return nil,
				fmt.Errorf("error parsing numeric literal: %s, err: %w",
					token.Literal, err)
		}
		return &NumericLiteral{Value: value}, nil
	case tokens.TokenStringLiteral:
		return &StringLiteral{Value: token.RawValue()}, nil
	case tokens.TokenBooleanLiteral:
		value, err := strconv.ParseBool(strings.ToLower(token.Literal))
		if err != nil {
			return nil,
				fmt.Errorf("error parsing boolean literal:str %s, err: %w",
					token.Literal, err)
		}
		return &BooleanLiteral{Value: value}, nil
	case tokens.TokenNull:
		return &NullValue{}, nil
	default:
		return nil, fmt.Errorf("unexpected literal type: %v", token.Type)
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
