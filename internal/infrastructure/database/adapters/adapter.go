// Package adapters provides database-specific implementations for PostgreSQL, MySQL, and SQLite.
package adapters

import (
	"context"
	"database/sql"
)

// Adapter defines the interface for database-specific operations
type Adapter interface {
	// Connect establishes a connection to the database
	Connect(config Config) (*sql.DB, error)

	// GetTableNames retrieves all table names from the database
	GetTableNames(ctx context.Context, db *sql.DB) ([]string, error)

	// GetTableDefinition retrieves the definition of a specific table
	GetTableDefinition(ctx context.Context, db *sql.DB, tableName string) (*TableDefinition, error)

	// GetDatabaseSchema retrieves schema information for all tables
	GetDatabaseSchema(ctx context.Context, db *sql.DB) (string, error)
}
