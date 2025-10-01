// Package execution provides SQL query execution functionality.
package execution

import (
	"context"
	"fmt"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/database"
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

// Execute executes a SQL query and returns the results
func (s *Service) Execute(ctx context.Context, query string) (*Result, error) {
	// Execute the query
	data, columns, err := s.conn.ExecuteQuery(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &Result{
		Rows:    data,
		Columns: columns,
	}, nil
}
