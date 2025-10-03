package cli

// QueryHistory represents a completed query with its prompt, SQL, and debug information
type QueryHistory struct {
	Prompt   string         // User's natural language prompt
	SQL      string         // Generated/executed SQL query
	Metadata map[string]any // Debug metadata (tokens, model, etc.)
}
