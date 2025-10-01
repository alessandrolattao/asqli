package query

import "errors"

// SQL represents a generated SQL query
type SQL struct {
	// The SQL query string
	Query string

	// Optional explanation of what the query does
	Explanation string
}

// Errors
var (
	ErrEmptyPrompt = errors.New("prompt cannot be empty")
	ErrInvalidSQL  = errors.New("invalid SQL query")
)
