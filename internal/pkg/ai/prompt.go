package ai

import (
	"fmt"
)

// Default system prompt template for SQL generation
const defaultSystemPrompt = `You are a helpful assistant that generates PostgreSQL queries based on natural language descriptions.

You'll receive database schema information that includes tables, their columns, data types, constraints,
and relationships between tables. Use this information to generate accurate SQL queries.

Respond ONLY with the SQL query without any explanation or markdown formatting. Do not include any comments
in the SQL or any additional text.`

// getSystemPrompt builds the system prompt with optional schema information
func getSystemPrompt(newSchema, cachedSchema string) string {
	systemPrompt := defaultSystemPrompt

	// First try to use the new schema if provided
	if newSchema != "" {
		return fmt.Sprintf("%s\n\nDatabase Schema:\n%s", systemPrompt, newSchema)
	}
	
	// Fall back to cached schema if available
	if cachedSchema != "" {
		return fmt.Sprintf("%s\n\nDatabase Schema:\n%s", systemPrompt, cachedSchema)
	}
	
	// Return default prompt if no schema is available
	return systemPrompt
}