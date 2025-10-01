// Package query provides SQL query generation from natural language using AI providers.
package query

import (
	"context"
	"fmt"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
)

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

// Generate generates a SQL query from a natural language prompt
func (s *Service) Generate(ctx context.Context, prompt, schema string) (*SQL, error) {
	// Validate input
	if prompt == "" {
		return nil, ErrEmptyPrompt
	}

	// Create request for AI provider
	req := &ai.GenerateRequest{
		Prompt: prompt,
		Schema: schema,
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
