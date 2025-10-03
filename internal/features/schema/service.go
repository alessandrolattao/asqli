// Package schema provides database schema extraction and caching functionality.
package schema

import (
	"context"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/database"
)

// Service handles database schema extraction and caching
type Service struct {
	conn  *database.Connection
	cache *Cache
}

// NewService creates a new schema service
func NewService(conn *database.Connection) *Service {
	return &Service{
		conn:  conn,
		cache: NewCache(),
	}
}

// Get retrieves the database schema (from cache if available)
func (s *Service) Get(ctx context.Context) (string, error) {
	// Check cache first
	if cached := s.cache.Get(); cached != "" {
		return cached, nil
	}

	// Extract schema from database
	schema, err := s.conn.GetDatabaseSchema(ctx)
	if err != nil {
		return "", err
	}

	// Store in cache
	s.cache.Set(schema)

	return schema, nil
}

// Invalidate clears the schema cache
func (s *Service) Invalidate() {
	s.cache.Clear()
}

// Refresh forces a schema refresh from the database
func (s *Service) Refresh(ctx context.Context) (string, error) {
	s.Invalidate()
	return s.Get(ctx)
}
