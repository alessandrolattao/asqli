// Package claude implements the AI provider interface for Anthropic's Claude models.
package claude

import (
	"context"
	"fmt"
	"strings"

	"github.com/alessandrolattao/asqli/internal/infrastructure/ai"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// ============================================
// Claude Provider Implementation
// ============================================

// Client implements the ai.Provider interface for Anthropic Claude
type Client struct {
	client      anthropic.Client
	model       anthropic.Model
	temperature float64
	maxTokens   int64
}

// Ensure Client implements ai.Provider interface
var _ ai.Provider = (*Client)(nil)

// ============================================
// Factory Registration
// ============================================

func init() {
	// Auto-register this provider on package import
	ai.RegisterProvider(ai.ProviderClaude, New)
}

// New creates a new Claude provider (implements ai.ProviderFactory)
func New(config ai.Config) (ai.Provider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("claude API key is required: %w", ai.ErrInvalidConfig)
	}

	// Default model
	model := anthropic.Model(config.Model)
	if config.Model == "" {
		model = "claude-sonnet-4-5"
	}

	// Default temperature
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.0 // Deterministic for SQL
	}

	// Default max tokens if not specified
	maxTokens := int64(config.MaxTokens)
	if maxTokens == 0 {
		maxTokens = 4096
	}

	client := anthropic.NewClient(
		option.WithAPIKey(config.APIKey),
	)

	return &Client{
		client:      client,
		model:       model,
		temperature: temperature,
		maxTokens:   maxTokens,
	}, nil
}

// ============================================
// Interface Implementation
// ============================================

// GenerateSQL generates a SQL query from natural language using Claude
func (c *Client) GenerateSQL(ctx context.Context, req *ai.GenerateRequest) (*ai.GenerateResponse, error) {
	if req.Prompt == "" {
		return nil, ai.ErrEmptyPrompt
	}

	systemPrompt := buildSystemPrompt(req.Schema, req.DatabaseType, req.Context)

	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       c.model,
		MaxTokens:   c.maxTokens,
		Temperature: param.NewOpt(c.temperature),
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	if len(message.Content) == 0 {
		return nil, ai.ErrGenerationFailed
	}

	// Extract text from the first content block
	var responseText string
	firstBlock := message.Content[0]

	// Use tagged switch for better readability
	switch firstBlock.Type {
	case "thinking":
		// Skip thinking blocks, get the next one
		if len(message.Content) > 1 {
			responseText = message.Content[1].AsText().Text
		} else {
			return nil, ai.ErrGenerationFailed
		}
	case "text":
		responseText = firstBlock.AsText().Text
	default:
		return nil, ai.ErrGenerationFailed
	}

	query := cleanSQLResponse(responseText)

	return &ai.GenerateResponse{
		Query:      query,
		Confidence: 1.0, // Claude doesn't provide confidence scores
		Usage: ai.UsageMetadata{
			Provider:       "claude",
			Model:          string(message.Model),
			PromptTokens:   int(message.Usage.InputTokens),
			ResponseTokens: int(message.Usage.OutputTokens),
			TotalTokens:    int(message.Usage.InputTokens + message.Usage.OutputTokens),
			CachedTokens:   int(message.Usage.CacheReadInputTokens + message.Usage.CacheCreationInputTokens),
		},
	}, nil
}

// Name returns the provider name
func (c *Client) Name() string {
	return "claude"
}

// Close releases any resources held by the provider
func (c *Client) Close() error {
	// Claude client doesn't need cleanup
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
