// Package query provides SQL query generation from natural language using AI providers.
package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
)

// History represents a previous query execution with prompt and SQL
type History struct {
	Prompt string // User's natural language prompt
	SQL    string // Generated/executed SQL query
}

// Service handles SQL query generation from natural language
type Service struct {
	aiProvider ai.Provider
}

// NewService creates a new query generation service
func NewService(aiProvider ai.Provider) *Service {
	return &Service{
		aiProvider: aiProvider,
	}
}

// Generate generates a SQL query from a natural language prompt with contextual awareness.
// It considers the database schema, query history, and currently selected table cell to generate
// contextually relevant SQL queries that understand follow-up requests and references.
func (s *Service) Generate(ctx context.Context, prompt, schema string, queryHistory []History, selectedColumn string, selectedValue any) (*SQL, error) {
	// Validate input
	if prompt == "" {
		return nil, ErrEmptyPrompt
	}

	// Build context from previous queries and selected cell
	var contextStr string
	var sb strings.Builder

	// Add selected cell context if available
	if selectedColumn != "" && selectedValue != nil {
		sb.WriteString("Currently selected cell:\n")
		sb.WriteString(fmt.Sprintf("Column: %s\n", selectedColumn))
		sb.WriteString(fmt.Sprintf("Value: %v\n", selectedValue))
		sb.WriteString("\nIf the user refers to 'selected', 'this', or similar terms, they likely mean this value.\n")
		sb.WriteString("Use this information to filter or reference specific data in your query.\n\n")
	}

	// Add previous queries context with both prompts and SQL
	if len(queryHistory) > 0 {
		sb.WriteString("Recent conversation history:\n")
		sb.WriteString("The user has previously asked the following questions and received these SQL queries:\n\n")
		for i, qh := range queryHistory {
			sb.WriteString(fmt.Sprintf("%d. User asked: \"%s\"\n", i+1, qh.Prompt))
			sb.WriteString(fmt.Sprintf("   Generated SQL: %s\n\n", qh.SQL))
		}
		sb.WriteString("Use this conversation context to understand what the user is referring to.\n")
		sb.WriteString("If the user's current request is a follow-up (e.g., \"show only the last 10\", \"filter by that user\", \"add a limit\"),\n")
		sb.WriteString("base your query on the most recent SQL but apply the requested modification.\n\n")
	}

	contextStr = sb.String()

	// Create request for AI provider
	req := &ai.GenerateRequest{
		Prompt:  prompt,
		Schema:  schema,
		Context: contextStr,
	}

	// Generate SQL via AI provider
	resp, err := s.aiProvider.GenerateSQL(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate SQL: %w", err)
	}

	// Validate generated SQL
	if err := s.Validate(resp.Query); err != nil {
		return nil, err
	}

	return &SQL{
		Query:       resp.Query,
		Explanation: resp.Explanation,
		Metadata:    resp.Metadata,
	}, nil
}

// Validate validates a SQL query
func (s *Service) Validate(query string) error {
	trimmed := strings.TrimSpace(query)

	if trimmed == "" {
		return ErrInvalidSQL
	}

	// Basic SQL validation (can be extended)
	upper := strings.ToUpper(trimmed)
	validStarts := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "WITH", "CREATE", "ALTER", "DROP"}

	for _, start := range validStarts {
		if strings.HasPrefix(upper, start) {
			return nil
		}
	}

	return ErrInvalidSQL
}

// IsDangerous checks if a query might modify data
func (s *Service) IsDangerous(query string) bool {
	trimmed := strings.TrimSpace(strings.ToUpper(query))
	dangerousStarts := []string{"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "CREATE", "TRUNCATE"}

	for _, start := range dangerousStarts {
		if strings.HasPrefix(trimmed, start) {
			return true
		}
	}

	return false
}
