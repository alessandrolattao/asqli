// Package database provides database connection management and adapter abstraction.
package database

import (
	"database/sql"
	"errors"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/database/adapters"
)

// Re-export types from adapters for convenience
type (
	DriverType           = adapters.DriverType
	Config               = adapters.Config
	TableDefinition      = adapters.TableDefinition
	ColumnDefinition     = adapters.ColumnDefinition
	ConstraintDefinition = adapters.ConstraintDefinition
)

// Re-export constants
const (
	PostgreSQL = adapters.PostgreSQL
	MySQL      = adapters.MySQL
	SQLite     = adapters.SQLite
)

// Errors
var (
	ErrUnsupportedDriver = errors.New("unsupported database driver")
	ErrInvalidConfig     = errors.New("invalid database configuration")
	ErrConnectionFailed  = errors.New("database connection failed")
)

// ExecuteQuery runs a SQL query and returns the result in a tabular format
func ExecuteQuery(db *sql.DB, query string) ([]map[string]any, []string, error) {
	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rows.Close() }()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	// Create a slice of interface{} to hold the values
	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))

	// Create a slice to store the result
	var result []map[string]any

	// Iterate through the result set
	for rows.Next() {
		// Initialize the pointers
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the valuePtrs
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, nil, err
		}

		// Create a map for this row
		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]

			b, ok := val.([]byte)
			if !ok {
				row[col] = val
				continue
			}
			row[col] = string(b)
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return result, columns, nil
}
