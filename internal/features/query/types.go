package query

// SQL represents a generated SQL query with metadata
type SQL struct {
	// The SQL query string
	Query string

	// Optional explanation of what the query does
	Explanation string

	// Metadata contains debug information (tokens, model, provider, etc.)
	Metadata map[string]any
}
