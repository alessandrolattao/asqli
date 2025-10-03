// Package ollama implements the AI provider interface for Ollama's local models.
package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	"github.com/ollama/ollama/api"
)

// ============================================
// Ollama Provider Implementation
// ============================================

// Client implements the ai.Provider interface for Ollama
type Client struct {
	client  *api.Client
	model   string
	baseURL string
}

// Ensure Client implements ai.Provider interface
var _ ai.Provider = (*Client)(nil)

// ============================================
// Factory Registration
// ============================================

func init() {
	// Auto-register this provider on package import
	ai.RegisterProvider(ai.ProviderOllama, New)
}

// New creates a new Ollama provider (implements ai.ProviderFactory)
func New(config ai.Config) (ai.Provider, error) {
	var client *api.Client
	var err error

	// If custom base URL is provided, use NewClient
	if config.BaseURL != "" {
		baseURL, parseErr := url.Parse(config.BaseURL)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid base URL: %w", parseErr)
		}
		client = api.NewClient(baseURL, http.DefaultClient)
	} else {
		// Otherwise use ClientFromEnvironment (uses OLLAMA_HOST env var or default localhost:11434)
		client, err = api.ClientFromEnvironment()
		if err != nil {
			return nil, fmt.Errorf("failed to create Ollama client: %w", err)
		}
	}

	// Determine which model to use
	model := config.Model
	if model == "" {
		// Auto-detect model from running or available models
		model, err = detectModel(client)
		if err != nil {
			return nil, fmt.Errorf("failed to detect Ollama model: %w (use --model to specify a model)", err)
		}
	}

	return &Client{
		client:  client,
		model:   model,
		baseURL: config.BaseURL,
	}, nil
}

// detectModel tries to find an available model to use
func detectModel(client *api.Client) (string, error) {
	ctx := context.Background()

	// First, try to get running models
	runningResp, err := client.ListRunning(ctx)
	if err == nil && len(runningResp.Models) > 0 {
		// Use the first running model
		return runningResp.Models[0].Model, nil
	}

	// If no running models, try to get locally available models
	listResp, err := client.List(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot list models: %w", err)
	}

	if len(listResp.Models) == 0 {
		return "", fmt.Errorf("no models available locally, please pull a model first")
	}

	// Use the first available model
	return listResp.Models[0].Model, nil
}

// ============================================
// Interface Implementation
// ============================================

// GenerateSQL generates a SQL query from natural language using Ollama
func (c *Client) GenerateSQL(ctx context.Context, req *ai.GenerateRequest) (*ai.GenerateResponse, error) {
	if req.Prompt == "" {
		return nil, ai.ErrEmptyPrompt
	}

	systemPrompt := buildSystemPrompt(req.Schema, req.DatabaseType, req.Context)

	// Prepare messages
	messages := []api.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: req.Prompt,
		},
	}

	// Prepare chat request
	chatReq := &api.ChatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   new(bool), // Disable streaming (false)
	}

	var fullResponse string
	var promptTokens, responseTokens int

	// Execute chat request
	err := c.client.Chat(ctx, chatReq, func(resp api.ChatResponse) error {
		// Accumulate response content
		fullResponse += resp.Message.Content

		if resp.Done {
			// Capture token usage when done
			promptTokens = int(resp.PromptEvalCount)
			responseTokens = int(resp.EvalCount)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ollama API error: %w", err)
	}

	if fullResponse == "" {
		return nil, ai.ErrGenerationFailed
	}

	query := cleanSQLResponse(fullResponse)

	return &ai.GenerateResponse{
		Query:      query,
		Confidence: 1.0, // Ollama doesn't provide confidence scores
		Usage: ai.UsageMetadata{
			Provider:       "ollama",
			Model:          c.model,
			PromptTokens:   promptTokens,
			ResponseTokens: responseTokens,
			TotalTokens:    promptTokens + responseTokens,
		},
	}, nil
}

// Name returns the provider name
func (c *Client) Name() string {
	return "ollama"
}

// Close releases any resources held by the provider
func (c *Client) Close() error {
	// Ollama client doesn't need cleanup
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
