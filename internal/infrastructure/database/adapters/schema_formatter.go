// Package adapters provides shared schema formatting utilities.
package adapters

import (
	"fmt"
	"strings"
)

// FormatTableDefinition formats a table definition into a human-readable string
// This is shared across all database adapters to ensure consistent formatting
func FormatTableDefinition(tableDef *TableDefinition) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("TABLE: %s\n", tableDef.Name))

	// Columns
	sb.WriteString("Columns:\n")
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
			// Note: Different databases use different syntax (AUTO_INCREMENT, AUTOINCREMENT, SERIAL)
			// For schema display purposes, we use a generic AUTO_INCREMENT label
			autoIncr = " AUTO_INCREMENT"
		}

		sb.WriteString(fmt.Sprintf("  %s %s %s%s%s%s\n",
			col.Name, col.Type, nullable, defaultVal, primaryKey, autoIncr))
	}

	// Constraints
	if len(tableDef.Constraints) > 0 {
		sb.WriteString("Constraints:\n")
		for _, constraint := range tableDef.Constraints {
			sb.WriteString(fmt.Sprintf("  %s: %s\n",
				constraint.Type, constraint.Definition))

			if constraint.Type == "FOREIGN KEY" && constraint.ReferencedTable != "" {
				sb.WriteString(fmt.Sprintf("    REFERENCES: %s\n",
					constraint.ReferencedTable))
			}
		}
	}

	sb.WriteString("\n")

	return sb.String()
}

// FormatDatabaseSchema formats all table definitions into a complete schema string
// This is shared across all database adapters to ensure consistent formatting
func FormatDatabaseSchema(tables []*TableDefinition) string {
	var sb strings.Builder
	sb.WriteString("DATABASE SCHEMA:\n\n")

	for _, tableDef := range tables {
		sb.WriteString(FormatTableDefinition(tableDef))
	}

	return sb.String()
}
