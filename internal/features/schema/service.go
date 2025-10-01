// Package schema provides database schema extraction and caching functionality.
package schema

import (
	"context"
	"sync"

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

// Cache provides thread-safe schema caching
type Cache struct {
	mu     sync.RWMutex
	schema string
}

// NewCache creates a new schema cache
func NewCache() *Cache {
	return &Cache{}
}

// Get retrieves the cached schema
func (c *Cache) Get() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.schema
}

// Set stores the schema in cache
func (c *Cache) Set(schema string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.schema = schema
}

// Clear removes the cached schema
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.schema = ""
}
