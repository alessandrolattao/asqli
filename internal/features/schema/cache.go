// Package schema provides database schema extraction and caching functionality.
package schema

import "sync"

// Cache provides thread-safe schema caching with optimized read performance.
// Uses RWMutex to allow multiple concurrent readers while ensuring exclusive write access.
// This is beneficial since schema is read frequently but modified rarely.
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
