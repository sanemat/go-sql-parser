package parser

import (
	"github.com/sanemat/go-sql-parser/tokens"
)

// Node represents a node in the Abstract Syntax Tree (AST).
type Node interface{}

// Expression interface for expressions in the AST.
type Expression interface {
	String() string
}

// SelectStatement represents a parsed SELECT statement.
type SelectStatement struct {
	Expressions []Expression
	Table       *string
	Where       *Condition
}

// Condition represents a condition in a WHERE clause.
type Condition struct {
	Column   string
	Operator string
	Value    string
}

// Parser holds the state of the parser.
type Parser struct {
	tokens []tokens.Token
	pos    int // current position in the token slice
}

// NewParser creates a new Parser instance.
func NewParser(tokens []tokens.Token) *Parser {
	return &Parser{tokens: tokens}
}
