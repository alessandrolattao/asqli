// Package gemini implements the AI provider interface for Google's Gemini models.
package gemini

import (
	"context"
	"fmt"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	"google.golang.org/genai"
)

// ============================================
// Gemini Provider Implementation
// ============================================

// Client implements the ai.Provider interface for Google Gemini
type Client struct {
	client      *genai.Client
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
	ai.RegisterProvider(ai.ProviderGemini, New)
}

// New creates a new Gemini provider (implements ai.ProviderFactory)
func New(config ai.Config) (ai.Provider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required: %w", ai.ErrInvalidConfig)
	}

	ctx := context.Background()

	// Create client with API key
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Default model
	model := config.Model
	if model == "" {
		model = "gemini-2.5-flash"
	}

	// Default temperature
	temperature := config.Temperature
	if temperature == 0 {
		temperature = 0.0 // Deterministic for SQL
	}

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

// GenerateSQL generates a SQL query from natural language using Gemini
func (c *Client) GenerateSQL(ctx context.Context, req *ai.GenerateRequest) (*ai.GenerateResponse, error) {
	if req.Prompt == "" {
		return nil, ai.ErrEmptyPrompt
	}

	systemPrompt := buildSystemPrompt(req.Schema, req.DatabaseType, req.Context)

	// Build the full prompt with system instructions and user query
	fullPrompt := fmt.Sprintf("%s\n\nUser query: %s", systemPrompt, req.Prompt)

	// Create content parts
	parts := []*genai.Part{
		{Text: fullPrompt},
	}

	// Create generation config
	var generationConfig *genai.GenerateContentConfig
	if c.maxTokens > 0 || c.temperature != 0 {
		generationConfig = &genai.GenerateContentConfig{}
		if c.maxTokens > 0 {
			generationConfig.MaxOutputTokens = int32(c.maxTokens)
		}
		if c.temperature != 0 {
			generationConfig.Temperature = genai.Ptr(float32(c.temperature))
		}
	}

	// Generate content
	result, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		[]*genai.Content{{Parts: parts}},
		generationConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("gemini API error: %w", err)
	}

	// Extract text from response
	if len(result.Candidates) == 0 {
		return nil, ai.ErrGenerationFailed
	}

	var queryText string
	for _, part := range result.Candidates[0].Content.Parts {
		if part.Text != "" {
			queryText += part.Text
		}
	}

	if queryText == "" {
		return nil, ai.ErrGenerationFailed
	}

	query := cleanSQLResponse(queryText)

	// Build usage metadata
	usage := ai.UsageMetadata{
		Provider: "gemini",
		Model:    c.model,
	}

	if result.UsageMetadata != nil {
		usage.PromptTokens = int(result.UsageMetadata.PromptTokenCount)
		usage.ResponseTokens = int(result.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(result.UsageMetadata.TotalTokenCount)
		if result.UsageMetadata.CachedContentTokenCount > 0 {
			usage.CachedTokens = int(result.UsageMetadata.CachedContentTokenCount)
		}
	}

	return &ai.GenerateResponse{
		Query:      query,
		Confidence: 1.0, // Gemini doesn't provide confidence scores
		Usage:      usage,
	}, nil
}

// Name returns the provider name
func (c *Client) Name() string {
	return "gemini"
}

// Close releases any resources held by the provider
func (c *Client) Close() error {
	// Gemini client doesn't need explicit cleanup
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
