// Package ai provides an abstraction layer for AI providers with a plugin-based architecture.
package ai

import (
	"context"
	"fmt"
	"maps"
	"slices"
)

// ============================================
// Provider Interface (Contract)
// ============================================

// Provider defines the interface that all AI providers must implement
type Provider interface {
	// GenerateSQL generates a SQL query from natural language prompt
	GenerateSQL(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// Name returns the provider name (e.g., "openai", "claude")
	Name() string

	// Close releases any resources held by the provider
	Close() error
}

// ============================================
// Request/Response Types
// ============================================

// GenerateRequest contains the input for SQL generation
type GenerateRequest struct {
	// User's natural language prompt
	Prompt string

	// Database schema context
	Schema string

	// Optional: Database type (postgres, mysql, sqlite)
	DatabaseType string

	// Optional: Additional context or examples
	Context string
}

// GenerateResponse contains the AI-generated SQL
type GenerateResponse struct {
	// Generated SQL query
	Query string

	// Confidence score (0.0-1.0) if supported by provider
	Confidence float64

	// Explanation of what the query does (optional)
	Explanation string

	// Usage metadata
	Usage UsageMetadata
}

// UsageMetadata contains standardized usage information from AI providers
type UsageMetadata struct {
	// AI provider name (e.g., "openai", "gemini")
	Provider string

	// Model name used for generation
	Model string

	// Number of tokens in the prompt
	PromptTokens int

	// Number of tokens in the response
	ResponseTokens int

	// Total tokens used (prompt + response)
	TotalTokens int

	// Number of cached tokens (if applicable)
	CachedTokens int
}

// ============================================
// Provider Type Registry
// ============================================

// ProviderType represents available AI providers
type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
	ProviderClaude ProviderType = "claude"
	ProviderGemini ProviderType = "gemini"
	ProviderOllama ProviderType = "ollama"
)

// ============================================
// Configuration
// ============================================

// Config contains provider configuration
type Config struct {
	// Provider type
	Type ProviderType

	// API Key (for cloud providers)
	APIKey string

	// Model name (e.g., "gpt-4-turbo", "claude-3-opus")
	Model string

	// Base URL (for custom endpoints)
	BaseURL string

	// Temperature (0.0-1.0)
	Temperature float64

	// Max tokens for response
	MaxTokens int

	// Provider-specific options
	Options map[string]any
}

// ============================================
// Factory Pattern
// ============================================

// ProviderFactory is a function that creates a new provider instance
type ProviderFactory func(config Config) (Provider, error)

var providerRegistry = make(map[ProviderType]ProviderFactory)

// RegisterProvider registers a new AI provider
func RegisterProvider(providerType ProviderType, factory ProviderFactory) {
	providerRegistry[providerType] = factory
}

// NewProvider creates a new AI provider based on config
func NewProvider(config Config) (Provider, error) {
	factory, exists := providerRegistry[config.Type]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedProvider, config.Type)
	}

	return factory(config)
}

// ListProviders returns a list of registered provider types
func ListProviders() []ProviderType {
	// Use maps.Keys + slices.Collect (Go 1.23+) for idiomatic iteration
	return slices.Collect(maps.Keys(providerRegistry))
}
