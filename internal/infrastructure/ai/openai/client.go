// Package openai implements the AI provider interface for OpenAI's GPT models.
package openai

import (
	"context"
	"fmt"
	"strings"

	"github.com/alessandrolattao/asqli/internal/infrastructure/ai"
	"github.com/sashabaranov/go-openai"
)

// ============================================
// OpenAI Provider Implementation
// ============================================

// Client implements the ai.Provider interface for OpenAI
type Client struct {
	client      *openai.Client
	model       string
	temperature float64
	maxTokens   int
}

// Ensure Client implements ai.Provider interface
var _ ai.Provider = (*Client)(nil)

// ============================================
// Factory Registration
// ============================================

func init() {
	// Auto-register this provider on package import
	ai.RegisterProvider(ai.ProviderOpenAI, New)
}

// New creates a new OpenAI provider (implements ai.ProviderFactory)
func New(config ai.Config) (ai.Provider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required: %w", ai.ErrInvalidConfig)
	}

	// Default model
	model := config.Model
	if model == "" {
		model = openai.GPT5Mini
	}

	// Default temperature
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.0 // Deterministic for SQL
	}

	client := openai.NewClient(config.APIKey)

	return &Client{
		client:      client,
		model:       model,
		temperature: temperature,
		maxTokens:   config.MaxTokens,
	}, nil
}

// ============================================
// Interface Implementation
// ============================================

// GenerateSQL generates a SQL query from natural language using OpenAI
func (c *Client) GenerateSQL(ctx context.Context, req *ai.GenerateRequest) (*ai.GenerateResponse, error) {
	if req.Prompt == "" {
		return nil, ai.ErrEmptyPrompt
	}

	systemPrompt := buildSystemPrompt(req.Schema, req.DatabaseType, req.Context)

	chatReq := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Prompt,
			},
		},
		// Temperature and TopP are omitted - reasoning models optimize these internally
	}

	if c.maxTokens > 0 {
		chatReq.MaxTokens = c.maxTokens
	}

	resp, err := c.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, ai.ErrGenerationFailed
	}

	query := cleanSQLResponse(resp.Choices[0].Message.Content)

	return &ai.GenerateResponse{
		Query:      query,
		Confidence: 1.0, // OpenAI doesn't provide confidence scores
		Usage: ai.UsageMetadata{
			Provider:       "openai",
			Model:          resp.Model,
			PromptTokens:   resp.Usage.PromptTokens,
			ResponseTokens: resp.Usage.CompletionTokens,
			TotalTokens:    resp.Usage.TotalTokens,
		},
	}, nil
}

// Name returns the provider name
func (c *Client) Name() string {
	return "openai"
}

// Close releases any resources held by the provider
func (c *Client) Close() error {
	// OpenAI client doesn't need cleanup
	return nil
}

// ============================================
// Helper Functions
// ============================================

// buildSystemPrompt constructs the system prompt with schema information
func buildSystemPrompt(schema, dbType, context string) string {
	prompt := `You are a helpful assistant that generates SQL queries based on natural language descriptions.

You'll receive database schema information that includes tables, their columns, data types, constraints,
and relationships between tables. Use this information to generate accurate SQL queries.

Respond ONLY with the SQL query without any explanation or markdown formatting. Do not include any comments
in the SQL or any additional text.`

	if dbType != "" {
		prompt += fmt.Sprintf("\n\nTarget database: %s", dbType)
	}

	if schema != "" {
		prompt += fmt.Sprintf("\n\n%s", schema)
	}

	if context != "" {
		prompt += fmt.Sprintf("\n\n%s", context)
	}

	return prompt
}

// cleanSQLResponse removes markdown formatting from AI response
func cleanSQLResponse(response string) string {
	// Remove markdown code blocks
	query := strings.TrimSpace(response)

	// Check for opening code block with language identifier
	if strings.HasPrefix(query, "```") {
		lines := strings.Split(query, "\n")
		if len(lines) > 2 {
			// Remove first line (```sql or ```)
			query = strings.Join(lines[1:], "\n")
		}
	}

	// Remove closing code block - TrimSuffix is idempotent
	query = strings.TrimSuffix(query, "```")

	return strings.TrimSpace(query)
}
