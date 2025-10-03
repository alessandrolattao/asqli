// Package database provides database connection management and adapter abstraction.
package database

import (
	"context"
	"database/sql"

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

// ExecuteQuery runs a SQL query with the given context and returns the result in a tabular format.
// Returns a slice of row maps, column names, and any error encountered.
func ExecuteQuery(ctx context.Context, db *sql.DB, query string) ([]map[string]any, []string, error) {
	// Execute the query with context
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer closeRows(rows)

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
		// Initialize the pointers (Go 1.22+ range style)
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
