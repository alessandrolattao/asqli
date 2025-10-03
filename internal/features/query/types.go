package query

import "github.com/alessandrolattao/sqlai/internal/infrastructure/ai"

// SQL represents a generated SQL query with metadata
type SQL struct {
	// The SQL query string
	Query string

	// Optional explanation of what the query does
	Explanation string

	// Usage metadata (tokens, model, provider, etc.)
	Usage ai.UsageMetadata
}
