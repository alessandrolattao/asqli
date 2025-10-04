// Package execution provides SQL query execution functionality.
package execution

import (
	"context"
	"fmt"

	"github.com/alessandrolattao/asqli/internal/infrastructure/database"
)

// Service handles SQL query execution
type Service struct {
	conn *database.Connection
}

// NewService creates a new execution service
func NewService(conn *database.Connection) *Service {
	return &Service{
		conn: conn,
	}
}

// Execute executes a SQL query with the given context and returns the results as a structured Result.
// The context is used for cancellation and timeout control.
func (s *Service) Execute(ctx context.Context, query string) (*Result, error) {
	// Execute the query with context
	data, columns, err := s.conn.ExecuteQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &Result{
		Rows:    data,
		Columns: columns,
	}, nil
}
