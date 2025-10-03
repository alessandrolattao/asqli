// Package query defines errors related to SQL query generation and validation.
package query

import "errors"

// Sentinel errors returned by the query service.
var (
	// ErrEmptyPrompt is returned when a prompt is empty
	ErrEmptyPrompt = errors.New("prompt cannot be empty")

	// ErrInvalidSQL is returned when generated SQL is invalid
	ErrInvalidSQL = errors.New("invalid SQL query")
)
