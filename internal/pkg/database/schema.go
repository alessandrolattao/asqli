package database

import (
	"database/sql"
)

// TableDefinition contains information about a database table
type TableDefinition struct {
	Name        string
	Columns     []ColumnDefinition
	Constraints []ConstraintDefinition
}

// ColumnDefinition contains information about a table column
type ColumnDefinition struct {
	Name       string
	Type       string
	Nullable   bool
	Default    string
	IsPrimary  bool
	IsAutoIncr bool
}

// ConstraintDefinition contains information about table constraints
type ConstraintDefinition struct {
	Name           string
	Type           string  // PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK, etc.
	Definition     string
	ReferencedTable string
	ReferencedColumns []string
}

// Common query execution code that works across database types
// ExecuteQuery runs a SQL query and returns the result in a tabular format
func ExecuteQuery(db *sql.DB, query string) ([]map[string]any, []string, error) {
	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

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
			var v any
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return result, columns, nil
}