package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/pkg/database"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// SQLiteAdapter implements the Adapter interface for SQLite
type SQLiteAdapter struct{}

// Connect establishes a connection to a SQLite database
func (a *SQLiteAdapter) Connect(config database.Config) (*sql.DB, error) {
	var connStr string
	
	if config.ConnectionString != "" {
		connStr = config.ConnectionString
	} else if config.FilePath != "" {
		connStr = config.FilePath
	} else if config.DBName != "" {
		// If only DBName is provided, use it as the file path
		connStr = config.DBName
	} else {
		return nil, database.ErrInvalidConfig
	}
	
	return sql.Open("sqlite3", connStr)
}

// GetTableNames retrieves all table names from a SQLite database
func (a *SQLiteAdapter) GetTableNames(db *sql.DB) ([]string, error) {
	query := `
		SELECT name FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
func (a *SQLiteAdapter) GetTableDefinition(db *sql.DB, tableName string) (*database.TableDefinition, error) {
	// Get pragma info for columns
	pragmaQuery := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	
	rows, err := db.Query(pragmaQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []database.ColumnDefinition
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, isPrimary int
		var defaultValue sql.NullString
		
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &isPrimary); err != nil {
			return nil, err
		}

		column := database.ColumnDefinition{
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
				createTableQuery := fmt.Sprintf("SELECT sql FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
				var createTableSQL string
				err := db.QueryRow(createTableQuery).Scan(&createTableSQL)
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
	fkeyQuery := fmt.Sprintf("PRAGMA foreign_key_list(%s)", tableName)
	
	fkeyRows, err := db.Query(fkeyQuery)
	if err != nil {
		return nil, err
	}
	defer fkeyRows.Close()

	var constraints []database.ConstraintDefinition
	fkeyMap := make(map[string]*database.ConstraintDefinition)
	
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
			constraint = &database.ConstraintDefinition{
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
		constraints = append(constraints, database.ConstraintDefinition{
			Name:        "pk_" + tableName,
			Type:        "PRIMARY KEY",
			Definition:  "PRIMARY KEY (" + strings.Join(primaryCols, ", ") + ")",
		})
	}
	
	// Get index information which can indicate UNIQUE constraints
	indexQuery := fmt.Sprintf("PRAGMA index_list(%s)", tableName)
	
	indexRows, err := db.Query(indexQuery)
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()
	
	for indexRows.Next() {
		var seq int
		var name string
		var unique, origin, partial int
		
		if err := indexRows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			return nil, err
		}
		
		// Only process unique constraints
		if unique == 1 && origin == 'u' {
			// Get the columns in this index
			indexInfoQuery := fmt.Sprintf("PRAGMA index_info(%s)", name)
			indexInfoRows, err := db.Query(indexInfoQuery)
			if err != nil {
				return nil, err
			}
			defer indexInfoRows.Close()
			
			var indexCols []string
			for indexInfoRows.Next() {
				var seqno, cid int
				var colName string
				
				if err := indexInfoRows.Scan(&seqno, &cid, &colName); err != nil {
					return nil, err
				}
				
				indexCols = append(indexCols, colName)
			}
			
			if len(indexCols) > 0 {
				constraints = append(constraints, database.ConstraintDefinition{
					Name:        name,
					Type:        "UNIQUE",
					Definition:  "UNIQUE (" + strings.Join(indexCols, ", ") + ")",
				})
			}
		}
	}

	return &database.TableDefinition{
		Name:        tableName,
		Columns:     columns,
		Constraints: constraints,
	}, nil
}

// GetDatabaseSchema retrieves schema information for all SQLite tables
func (a *SQLiteAdapter) GetDatabaseSchema(db *sql.DB) (string, error) {
	tables, err := a.GetTableNames(db)
	if err != nil {
		return "", err
	}

	var schemaBuilder strings.Builder
	schemaBuilder.WriteString("DATABASE SCHEMA:\n\n")

	for _, tableName := range tables {
		tableDef, err := a.GetTableDefinition(db, tableName)
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
				autoIncr = " AUTOINCREMENT"
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