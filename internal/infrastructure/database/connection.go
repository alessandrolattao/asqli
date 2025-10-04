package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alessandrolattao/asqli/internal/infrastructure/database/adapters"
)

// init registers all built-in database adapters
func init() {
	RegisterAdapter(adapters.PostgreSQL, func() adapters.Adapter {
		return &adapters.PostgresAdapter{}
	})
	RegisterAdapter(adapters.MySQL, func() adapters.Adapter {
		return &adapters.MySQLAdapter{}
	})
	RegisterAdapter(adapters.SQLite, func() adapters.Adapter {
		return &adapters.SQLiteAdapter{}
	})
}

// Connection represents a database connection
type Connection struct {
	// The underlying database connection
	DB *sql.DB

	// The driver type for this connection
	DriverType adapters.DriverType

	// Driver-specific adapter
	adapter adapters.Adapter
}

// Open establishes a connection to the database using the specified configuration
func Open(config adapters.Config) (*Connection, error) {
	// Create adapter based on driver type
	adapter, err := NewAdapter(config.DriverType)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %w", err)
	}

	// Connect to the database
	db, err := adapter.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	// Ensure db is not nil before testing connection
	if db == nil {
		return nil, fmt.Errorf("%w: database connection is nil", ErrConnectionFailed)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		// Safe to close: db is guaranteed to be non-nil at this point
		_ = db.Close()
		return nil, fmt.Errorf("%w: ping failed: %v", ErrConnectionFailed, err)
	}

	return &Connection{
		DB:         db,
		DriverType: config.DriverType,
		adapter:    adapter,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// ExecuteQuery runs a SQL query with the given context and returns the result.
// The context is used for cancellation and timeout control.
func (c *Connection) ExecuteQuery(ctx context.Context, query string) ([]map[string]any, []string, error) {
	return ExecuteQuery(ctx, c.DB, query)
}

// GetTableNames retrieves all table names from the database using the given context.
func (c *Connection) GetTableNames(ctx context.Context) ([]string, error) {
	return c.adapter.GetTableNames(ctx, c.DB)
}

// GetTableDefinition retrieves the definition of a specific table using the given context.
func (c *Connection) GetTableDefinition(ctx context.Context, tableName string) (*adapters.TableDefinition, error) {
	return c.adapter.GetTableDefinition(ctx, c.DB, tableName)
}

// GetDatabaseSchema retrieves schema information for all tables using the given context.
func (c *Connection) GetDatabaseSchema(ctx context.Context) (string, error) {
	return c.adapter.GetDatabaseSchema(ctx, c.DB)
}
