package cli

import "github.com/alessandrolattao/asqli/internal/infrastructure/ai"

// QueryHistory represents a completed query with its prompt, SQL, and debug information
type QueryHistory struct {
	Prompt string           // User's natural language prompt
	SQL    string           // Generated/executed SQL query
	Usage  ai.UsageMetadata // Usage metadata (tokens, model, provider, etc.)
}
