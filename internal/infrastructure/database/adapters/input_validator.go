// Package adapters provides database identifier validation to prevent SQL injection attacks.
package adapters

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidTableName is returned when a table name contains invalid characters
	ErrInvalidTableName = errors.New("invalid table name: must contain only alphanumeric characters, underscores, and start with a letter")

	// ErrInvalidIndexName is returned when an index name contains invalid characters
	ErrInvalidIndexName = errors.New("invalid index name: must contain only alphanumeric characters and underscores")

	// ErrEmptyName is returned when a name is empty
	ErrEmptyName = errors.New("name cannot be empty")
)

// tableNameRegex is an internal regex used to validate SQL identifier names (tables, columns, indexes).
// Allowed characters: letters, numbers, underscores. Must start with letter or underscore.
// This strict pattern prevents SQL injection in PRAGMA and dynamic SQL statements.
var tableNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ValidateTableName validates a table name to prevent SQL injection
// Returns an error if the name is invalid
func ValidateTableName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	// Check for SQL keywords and dangerous patterns
	upperName := strings.ToUpper(name)
	dangerousKeywords := []string{
		"DROP", "DELETE", "INSERT", "UPDATE", "ALTER",
		"CREATE", "EXEC", "EXECUTE", "UNION", "SELECT",
		"--", "/*", "*/", ";", "'", "\"",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperName, keyword) {
			return ErrInvalidTableName
		}
	}

	// Validate against regex pattern
	if !tableNameRegex.MatchString(name) {
		return ErrInvalidTableName
	}

	return nil
}

// ValidateIndexName validates an index name to prevent SQL injection
// Uses the same rules as table names
func ValidateIndexName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	if !tableNameRegex.MatchString(name) {
		return ErrInvalidIndexName
	}

	return nil
}

// SanitizeIdentifier sanitizes an SQL identifier (table, column, index name)
// Returns the sanitized name or an error if validation fails
func SanitizeIdentifier(name string) (string, error) {
	if err := ValidateTableName(name); err != nil {
		return "", err
	}
	return name, nil
}
