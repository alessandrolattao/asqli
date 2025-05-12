package ai

// Client represents an AI client for generating SQL queries
type Client struct {
	openaiClient OpenAIClientInterface
	schemaCache  string // Cache the schema to avoid repeated fetching
}

// QueryGenerator defines the interface for SQL query generation
type QueryGenerator interface {
	GenerateQuery(prompt string) (string, error)
	GenerateQueryWithSchema(prompt, schema string) (string, error)
	ClearSchemaCache()
}