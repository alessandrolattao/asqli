// Package database provides a generic interface for different database engines
package database

import (
	"database/sql"
)

// DriverType represents the supported database driver types
type DriverType string

const (
	// PostgreSQL driver
	PostgreSQL DriverType = "postgres"
	// MySQL driver
	MySQL DriverType = "mysql"
	// SQLite driver
	SQLite DriverType = "sqlite3"
)

// Config represents generic database configuration parameters
type Config struct {
	// DriverType specifies which database driver to use
	DriverType DriverType
	
	// ConnectionString is the full connection string. If provided, it takes precedence 
	// over individual connection parameters
	ConnectionString string
	
	// Individual connection parameters
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	
	// Additional parameters specific to certain databases
	SSLMode   string // For PostgreSQL
	FilePath  string // For SQLite
	ParseTime bool   // For MySQL
}

// Connection represents a database connection
type Connection struct {
	// The underlying database connection
	DB *sql.DB
	
	// The driver type for this connection
	DriverType DriverType
	
	// Driver-specific adapter
	Adapter Adapter
}

// Adapter defines the interface for database-specific operations
type Adapter interface {
	// Connect establishes a connection to the database
	Connect(config Config) (*sql.DB, error)
	
	// GetTableNames retrieves all table names from the database
	GetTableNames(db *sql.DB) ([]string, error)
	
	// GetTableDefinition retrieves the definition of a specific table
	GetTableDefinition(db *sql.DB, tableName string) (*TableDefinition, error)
	
	// GetDatabaseSchema retrieves schema information for all tables
	GetDatabaseSchema(db *sql.DB) (string, error)
}

// Open establishes a connection to the database using the specified configuration
func Open(config Config, adapter Adapter) (*Connection, error) {
	if adapter == nil {
		return nil, ErrUnsupportedDriver
	}
	
	// Connect to the database
	db, err := adapter.Connect(config)
	if err != nil {
		return nil, err
	}
	
	return &Connection{
		DB:         db,
		DriverType: config.DriverType,
		Adapter:    adapter,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// ExecuteQuery runs a SQL query and returns the result
func (c *Connection) ExecuteQuery(query string) ([]map[string]any, []string, error) {
	return ExecuteQuery(c.DB, query)
}

// GetTableNames retrieves all table names from the database
func (c *Connection) GetTableNames() ([]string, error) {
	return c.Adapter.GetTableNames(c.DB)
}

// GetTableDefinition retrieves the definition of a specific table
func (c *Connection) GetTableDefinition(tableName string) (*TableDefinition, error) {
	return c.Adapter.GetTableDefinition(c.DB, tableName)
}

// GetDatabaseSchema retrieves schema information for all tables
func (c *Connection) GetDatabaseSchema() (string, error) {
	return c.Adapter.GetDatabaseSchema(c.DB)
}