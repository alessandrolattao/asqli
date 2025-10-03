package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"errors"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// SQLiteAdapter implements the Adapter interface for SQLite
type SQLiteAdapter struct{}

// Ensure SQLiteAdapter implements Adapter interface
var _ Adapter = (*SQLiteAdapter)(nil)

// Connect establishes a connection to a SQLite database
func (a *SQLiteAdapter) Connect(config Config) (*sql.DB, error) {
	if config.ConnectionString != "" {
		return sql.Open("sqlite3", config.ConnectionString)
	}

	if config.FilePath != "" {
		return sql.Open("sqlite3", config.FilePath)
	}

	if config.DBName != "" {
		// If only DBName is provided, use it as the file path
		return sql.Open("sqlite3", config.DBName)
	}

	return nil, errors.New("invalid database configuration")
}

// GetTableNames retrieves all table names from a SQLite database
func (a *SQLiteAdapter) GetTableNames(ctx context.Context, db *sql.DB) ([]string, error) {
	query := `
		SELECT name FROM sqlite_master
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// GetTableDefinition retrieves the definition of a specific SQLite table
func (a *SQLiteAdapter) GetTableDefinition(ctx context.Context, db *sql.DB, tableName string) (*TableDefinition, error) {
	// Validate table name to prevent SQL injection
	if err := ValidateTableName(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %w", err)
	}

	// Get pragma info for columns
	// Note: PRAGMA statements don't support parameter binding, so we validate the input strictly
	pragmaQuery := fmt.Sprintf("PRAGMA table_info(%s)", tableName)

	rows, err := db.QueryContext(ctx, pragmaQuery)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var columns []ColumnDefinition
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, isPrimary int
		var defaultValue sql.NullString

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &isPrimary); err != nil {
			return nil, err
		}

		column := ColumnDefinition{
			Name:      name,
			Type:      dataType,
			Nullable:  notNull == 0,
			IsPrimary: isPrimary == 1,
		}

		if defaultValue.Valid {
			column.Default = defaultValue.String
		}

		// Check if it's an autoincrement column
		if isPrimary == 1 {
			// In SQLite, autoincrement is only applicable to INTEGER PRIMARY KEY columns
			if strings.ToUpper(dataType) == "INTEGER" {
				// Get the SQL that created the table
				// tableName is already validated, safe to use in query
				createTableQuery := "SELECT sql FROM sqlite_master WHERE type='table' AND name=?"
				var createTableSQL string
				err := db.QueryRowContext(ctx, createTableQuery, tableName).Scan(&createTableSQL)
				if err != nil {
					return nil, err
				}

				// Check for AUTOINCREMENT keyword - simplified approach
				column.IsAutoIncr = strings.Contains(createTableSQL, "AUTOINCREMENT") &&
					column.IsPrimary && strings.ToUpper(dataType) == "INTEGER"
			}
		}

		columns = append(columns, column)
	}

	// Get foreign key constraints
	// tableName is already validated, safe to use in PRAGMA
	fkeyQuery := fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName)

	fkeyRows, err := db.QueryContext(ctx, fkeyQuery)
	if err != nil {
		return nil, err
	}
	defer func() { _ = fkeyRows.Close() }()

	var constraints []ConstraintDefinition
	fkeyMap := make(map[string]*ConstraintDefinition)

	for fkeyRows.Next() {
		var id, seq int
		var refTable, from, to string
		var onUpdate, onDelete string
		var match string

		if err := fkeyRows.Scan(&id, &seq, &refTable, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			return nil, err
		}

		// For SQLite, we need to group together the columns that belong to the same constraint
		fkeyID := fmt.Sprintf("fk_%d", id)
		constraint, exists := fkeyMap[fkeyID]

		if !exists {
			constraint = &ConstraintDefinition{
				Name:              fkeyID,
				Type:              "FOREIGN KEY",
				ReferencedTable:   refTable,
				ReferencedColumns: []string{},
			}
			fkeyMap[fkeyID] = constraint
		}

		constraint.ReferencedColumns = append(constraint.ReferencedColumns, to)

		// Build a definition string
		columns := strings.Join(constraint.ReferencedColumns, ", ")
		constraint.Definition = fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
			from, refTable, columns)
	}

	// Convert the map to a slice
	for _, c := range fkeyMap {
		constraints = append(constraints, *c)
	}

	// Add primary key constraint if we have primary key columns
	var primaryCols []string
	for _, col := range columns {
		if col.IsPrimary {
			primaryCols = append(primaryCols, col.Name)
		}
	}

	if len(primaryCols) > 0 {
		constraints = append(constraints, ConstraintDefinition{
			Name:       "pk_" + tableName,
			Type:       "PRIMARY KEY",
			Definition: "PRIMARY KEY (" + strings.Join(primaryCols, ", ") + ")",
		})
	}

	// Get index information which can indicate UNIQUE constraints
	// tableName is already validated, safe to use in PRAGMA
	indexQuery := fmt.Sprintf("PRAGMA index_list(%s)", tableName)

	indexRows, err := db.QueryContext(ctx, indexQuery)
	if err != nil {
		return nil, err
	}
	defer func() { _ = indexRows.Close() }()

	for indexRows.Next() {
		var seq int
		var name string
		var unique, origin, partial int

		if err := indexRows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			return nil, err
		}

		// Only process unique constraints
		if unique == 1 && origin == 'u' {
			// Validate index name before using it in query
			if err := ValidateIndexName(name); err != nil {
				// Skip invalid index names instead of failing
				continue
			}

			// Get the columns in this index
			// index name is validated, safe to use in PRAGMA
			indexInfoQuery := fmt.Sprintf("PRAGMA index_info(%s)", name)
			indexInfoRows, err := db.QueryContext(ctx, indexInfoQuery)
			if err != nil {
				return nil, err
			}

			var indexCols []string
			for indexInfoRows.Next() {
				var seqno, cid int
				var colName string

				if err := indexInfoRows.Scan(&seqno, &cid, &colName); err != nil {
					_ = indexInfoRows.Close()
					return nil, err
				}

				indexCols = append(indexCols, colName)
			}

			// Close rows immediately after use (not deferred in loop)
			_ = indexInfoRows.Close()

			if len(indexCols) > 0 {
				constraints = append(constraints, ConstraintDefinition{
					Name:       name,
					Type:       "UNIQUE",
					Definition: "UNIQUE (" + strings.Join(indexCols, ", ") + ")",
				})
			}
		}
	}

	return &TableDefinition{
		Name:        tableName,
		Columns:     columns,
		Constraints: constraints,
	}, nil
}

// GetDatabaseSchema retrieves schema information for all SQLite tables
func (a *SQLiteAdapter) GetDatabaseSchema(ctx context.Context, db *sql.DB) (string, error) {
	tables, err := a.GetTableNames(ctx, db)
	if err != nil {
		return "", err
	}

	// Get all table definitions
	tableDefs := make([]*TableDefinition, 0, len(tables))
	for _, tableName := range tables {
		tableDef, err := a.GetTableDefinition(ctx, db, tableName)
		if err != nil {
			return "", err
		}
		tableDefs = append(tableDefs, tableDef)
	}

	// Use shared formatter
	return FormatDatabaseSchema(tableDefs), nil
}
