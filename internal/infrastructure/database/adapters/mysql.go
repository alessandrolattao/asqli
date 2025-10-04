package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
)

// MySQLAdapter implements the Adapter interface for MySQL
type MySQLAdapter struct{}

// Ensure MySQLAdapter implements Adapter interface
var _ Adapter = (*MySQLAdapter)(nil)

// Connect establishes a connection to a MySQL database
func (a *MySQLAdapter) Connect(config Config) (*sql.DB, error) {
	if config.ConnectionString != "" {
		return sql.Open("mysql", config.ConnectionString)
	}

	// MySQL connection string: username:password@tcp(host:port)/dbname
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	// Add parameters like parseTime if needed
	params := []string{}
	if config.ParseTime {
		params = append(params, "parseTime=true")
	}

	if len(params) > 0 {
		connStr = connStr + "?" + strings.Join(params, "&")
	}

	return sql.Open("mysql", connStr)
}

// GetTableNames retrieves all table names from a MySQL database
func (a *MySQLAdapter) GetTableNames(ctx context.Context, db *sql.DB) ([]string, error) {
	query := "SHOW TABLES"

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

// GetTableDefinition retrieves the definition of a specific MySQL table
func (a *MySQLAdapter) GetTableDefinition(ctx context.Context, db *sql.DB, tableName string) (*TableDefinition, error) {
	// Get columns
	columnsQuery := `
		SELECT
			COLUMN_NAME,
			DATA_TYPE,
			IS_NULLABLE,
			COLUMN_DEFAULT,
			COLUMN_KEY = 'PRI' AS is_primary,
			EXTRA = 'auto_increment' AS is_autoincrement
		FROM
			INFORMATION_SCHEMA.COLUMNS
		WHERE
			TABLE_SCHEMA = DATABASE() AND
			TABLE_NAME = ?
		ORDER BY
			ORDINAL_POSITION
	`

	rows, err := db.QueryContext(ctx, columnsQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var columns []ColumnDefinition
	for rows.Next() {
		var column ColumnDefinition
		var isNullable string
		var defaultValue sql.NullString
		var isPrimary bool
		var isAutoIncr bool

		if err := rows.Scan(
			&column.Name,
			&column.Type,
			&isNullable,
			&defaultValue,
			&isPrimary,
			&isAutoIncr,
		); err != nil {
			return nil, err
		}

		column.Nullable = isNullable == "YES"
		if defaultValue.Valid {
			column.Default = defaultValue.String
		}
		column.IsPrimary = isPrimary
		column.IsAutoIncr = isAutoIncr

		columns = append(columns, column)
	}

	// Get constraints (foreign keys)
	constraintsQuery := `
		SELECT DISTINCT
			tc.CONSTRAINT_NAME,
			tc.CONSTRAINT_TYPE,
			'' AS constraint_definition,
			COALESCE(kcu.REFERENCED_TABLE_NAME, '') AS REFERENCED_TABLE_NAME
		FROM
			INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
		LEFT JOIN
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
			ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
			AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
			AND tc.TABLE_NAME = kcu.TABLE_NAME
		WHERE
			tc.TABLE_SCHEMA = DATABASE() AND
			tc.TABLE_NAME = ? AND
			tc.CONSTRAINT_TYPE != 'CHECK'
		ORDER BY
			tc.CONSTRAINT_TYPE
	`

	constraintRows, err := db.QueryContext(ctx, constraintsQuery, tableName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = constraintRows.Close() }()

	var constraints []ConstraintDefinition
	for constraintRows.Next() {
		var constraint ConstraintDefinition
		if err := constraintRows.Scan(
			&constraint.Name,
			&constraint.Type,
			&constraint.Definition,
			&constraint.ReferencedTable,
		); err != nil {
			return nil, err
		}

		// Get constraint definition for foreign keys
		if constraint.Type == "FOREIGN KEY" {
			keyQuery := `
				SELECT
					COLUMN_NAME,
					REFERENCED_COLUMN_NAME
				FROM
					INFORMATION_SCHEMA.KEY_COLUMN_USAGE
				WHERE
					TABLE_SCHEMA = DATABASE() AND
					TABLE_NAME = ? AND
					CONSTRAINT_NAME = ?
			`

			keyRows, err := db.QueryContext(ctx, keyQuery, tableName, constraint.Name)
			if err != nil {
				return nil, err
			}

			var columnNames []string
			var refColumnNames []string

			for keyRows.Next() {
				var columnName, refColumnName string
				if err := keyRows.Scan(&columnName, &refColumnName); err != nil {
					_ = keyRows.Close()
					return nil, err
				}

				columnNames = append(columnNames, columnName)
				refColumnNames = append(refColumnNames, refColumnName)
			}

			// Close rows immediately after use (not deferred in loop)
			_ = keyRows.Close()

			constraint.Definition = fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)",
				strings.Join(columnNames, ", "),
				constraint.ReferencedTable,
				strings.Join(refColumnNames, ", "))

			constraint.ReferencedColumns = refColumnNames
		}

		constraints = append(constraints, constraint)
	}

	return &TableDefinition{
		Name:        tableName,
		Columns:     columns,
		Constraints: constraints,
	}, nil
}

// GetDatabaseSchema retrieves schema information for all MySQL tables
func (a *MySQLAdapter) GetDatabaseSchema(ctx context.Context, db *sql.DB) (string, error) {
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
