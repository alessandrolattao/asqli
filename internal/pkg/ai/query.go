package ai

import (
	"context"
	"errors"
	"strings"
)

// GenerateQuery generates a SQL query based on a natural language prompt
func (c *Client) GenerateQuery(prompt string) (string, error) {
	return c.GenerateQueryWithSchema(prompt, "")
}

// GenerateQueryWithSchema generates a SQL query based on a natural language prompt and database schema
func (c *Client) GenerateQueryWithSchema(prompt, schema string) (string, error) {
	systemPrompt := getSystemPrompt(schema, c.schemaCache)
	
	// Update schema cache if new schema is provided
	if schema != "" {
		c.schemaCache = schema
	}

	// Create and send the request to the OpenAI API
	response, err := c.openaiClient.CreateCompletion(context.Background(), systemPrompt, prompt)
	if err != nil {
		return "", err
	}

	if response == "" {
		return "", errors.New("no response from AI")
	}

	// Process the response to extract the SQL query
	query := cleanQueryResponse(response)
	return query, nil
}

// ClearSchemaCache clears the cached schema
func (c *Client) ClearSchemaCache() {
	c.schemaCache = ""
}

// cleanQueryResponse removes any markdown formatting from the AI response
func cleanQueryResponse(query string) string {
	// Remove any markdown code block syntax (```sql or ``` at beginning and end)
	if len(query) > 3 {
		// Remove opening code block markers (```sql or ```)
		if query[0:3] == "```" {
			// Find the end of the first line
			firstLineEnd := 0
			for i, char := range query {
				if char == '\n' {
					firstLineEnd = i
					break
				}
			}
			if firstLineEnd > 0 {
				query = query[firstLineEnd+1:]
			} else {
				query = query[3:] // No newline found, just remove ```
			}
		}

		// Remove closing code block markers (```)
		if len(query) >= 3 && query[len(query)-3:] == "```" {
			query = query[:len(query)-3]
		}
	}

	// Trim any extra whitespace at beginning and end
	return strings.TrimSpace(query)
}