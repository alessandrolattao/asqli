// Package ai defines errors related to AI provider operations and configuration.
package ai

import "errors"

// Sentinel errors returned by AI providers and the provider factory.
var (
	// ErrUnsupportedProvider is returned when an unsupported AI provider is requested
	ErrUnsupportedProvider = errors.New("unsupported AI provider")

	// ErrInvalidConfig is returned when provider configuration is invalid
	ErrInvalidConfig = errors.New("invalid provider configuration")

	// ErrGenerationFailed is returned when SQL generation fails
	ErrGenerationFailed = errors.New("SQL generation failed")

	// ErrEmptyPrompt is returned when the prompt is empty
	ErrEmptyPrompt = errors.New("prompt cannot be empty")
)
