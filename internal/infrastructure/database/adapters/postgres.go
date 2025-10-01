package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// PostgresAdapter implements the Adapter interface for PostgreSQL
type PostgresAdapter struct{}

// Ensure PostgresAdapter implements Adapter interface
var _ Adapter = (*PostgresAdapter)(nil)

// Connect establishes a connection to a PostgreSQL database
func (a *PostgresAdapter) Connect(config Config) (*sql.DB, error) {
	if config.ConnectionString != "" {
		return sql.Open("postgres", config.ConnectionString)
	}

	sslMode := "disable"
	if config.SSLMode != "" {
		sslMode = config.SSLMode
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, sslMode)

	return sql.Open("postgres", connStr)
}

// GetTableNames retrieves all table names from a PostgreSQL database
func (a *PostgresAdapter) GetTableNames(ctx context.Context, db *sql.DB) ([]string, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
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

// GetTableDefinition retrieves the definition of a specific PostgreSQL table
func (a *PostgresAdapter) GetTableDefinition(ctx context.Context, db *sql.DB, tableName string) (*TableDefinition, error) {
	// Get columns
	columnsQuery := `
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			CASE WHEN pk.column_name IS NOT NULL THEN true ELSE false END AS is_primary,
			CASE WHEN c.column_default LIKE '%nextval%' THEN true ELSE false END AS is_autoincrement
		FROM
			information_schema.columns c
		LEFT JOIN (
			SELECT kcu.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu
				ON tc.constraint_name = kcu.constraint_name
			WHERE tc.constraint_type = 'PRIMARY KEY'
			AND tc.table_name = $1
		) pk ON c.column_name = pk.column_name
		WHERE
			c.table_schema = 'public' AND
			c.table_name = $1
		ORDER BY
			c.ordinal_position
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

	// Get constraints
	constraintsQuery := `
		SELECT
			c.conname AS constraint_name,
			CASE
				WHEN c.contype = 'p' THEN 'PRIMARY KEY'
				WHEN c.contype = 'f' THEN 'FOREIGN KEY'
				WHEN c.contype = 'u' THEN 'UNIQUE'
				WHEN c.contype = 'c' THEN 'CHECK'
				ELSE c.contype::text
			END AS constraint_type,
			pg_get_constraintdef(c.oid) AS constraint_definition,
			CASE
				WHEN c.contype = 'f' THEN
					(SELECT relname FROM pg_class WHERE oid = c.confrelid)
				ELSE ''
			END AS referenced_table
		FROM
			pg_constraint c
		JOIN
			pg_class cl ON cl.oid = c.conrelid
		WHERE
			cl.relname = $1 AND
			cl.relkind = 'r'
		ORDER BY
			c.contype
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
		constraints = append(constraints, constraint)
	}

	return &TableDefinition{
		Name:        tableName,
		Columns:     columns,
		Constraints: constraints,
	}, nil
}

// GetDatabaseSchema retrieves schema information for all PostgreSQL tables
func (a *PostgresAdapter) GetDatabaseSchema(ctx context.Context, db *sql.DB) (string, error) {
	tables, err := a.GetTableNames(ctx, db)
	if err != nil {
		return "", err
	}

	var schemaBuilder strings.Builder
	schemaBuilder.WriteString("DATABASE SCHEMA:\n\n")

	for _, tableName := range tables {
		tableDef, err := a.GetTableDefinition(ctx, db, tableName)
		if err != nil {
			return "", err
		}

		schemaBuilder.WriteString(fmt.Sprintf("TABLE: %s\n", tableDef.Name))

		// Columns
		schemaBuilder.WriteString("Columns:\n")
		for _, col := range tableDef.Columns {
			nullable := "NOT NULL"
			if col.Nullable {
				nullable = "NULL"
			}

			defaultVal := ""
			if col.Default != "" {
				defaultVal = fmt.Sprintf(" DEFAULT %s", col.Default)
			}

			primaryKey := ""
			if col.IsPrimary {
				primaryKey = " PRIMARY KEY"
			}

			autoIncr := ""
			if col.IsAutoIncr {
				autoIncr = " AUTO INCREMENT"
			}

			schemaBuilder.WriteString(fmt.Sprintf("  %s %s %s%s%s%s\n",
				col.Name, col.Type, nullable, defaultVal, primaryKey, autoIncr))
		}

		// Constraints
		if len(tableDef.Constraints) > 0 {
			schemaBuilder.WriteString("Constraints:\n")
			for _, constraint := range tableDef.Constraints {
				schemaBuilder.WriteString(fmt.Sprintf("  %s: %s\n",
					constraint.Type, constraint.Definition))

				if constraint.Type == "FOREIGN KEY" && constraint.ReferencedTable != "" {
					schemaBuilder.WriteString(fmt.Sprintf("    REFERENCES: %s\n",
						constraint.ReferencedTable))
				}
			}
		}

		schemaBuilder.WriteString("\n")
	}

	return schemaBuilder.String(), nil
}
