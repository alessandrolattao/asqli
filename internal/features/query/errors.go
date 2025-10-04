// Package query defines errors related to SQL query generation and validation.
package query

import (
	"errors"
	"fmt"
)

// Sentinel errors returned by the query service.
var (
	// ErrEmptyPrompt is returned when a prompt is empty
	ErrEmptyPrompt = errors.New("prompt cannot be empty")

	// ErrInvalidSQL is returned when generated SQL is invalid
	ErrInvalidSQL = errors.New("invalid SQL query")
)

// ValidationError is returned when SQL validation fails.
// It includes the invalid query for debugging purposes.
type ValidationError struct {
	Query string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%v: %s", e.Err, e.Query)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
