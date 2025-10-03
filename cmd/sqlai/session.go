package main

import (
	"fmt"
	"os"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database/adapters"
	"github.com/alessandrolattao/sqlai/internal/ui/cli"
)

// runQuerySession starts a query session with the specified database and AI provider
func runQuerySession(dbConfig adapters.Config, providerStr string, modelStr string) {
	// Determine AI provider type
	var providerType ai.ProviderType
	var apiKeyEnvVar string

	switch providerStr {
	case "openai":
		providerType = ai.ProviderOpenAI
		apiKeyEnvVar = "OPENAI_API_KEY"
	case "claude":
		providerType = ai.ProviderClaude
		apiKeyEnvVar = "ANTHROPIC_API_KEY"
	case "gemini":
		providerType = ai.ProviderGemini
		apiKeyEnvVar = "GEMINI_API_KEY"
	case "ollama":
		providerType = ai.ProviderOllama
		apiKeyEnvVar = "" // Ollama doesn't require an API key
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported AI provider '%s'. Supported providers: openai, claude, gemini, ollama\n", providerStr)
		os.Exit(1)
	}

	// Check AI Provider API key (skip for Ollama)
	var apiKey string
	if apiKeyEnvVar != "" {
		apiKey = os.Getenv(apiKeyEnvVar)
		if apiKey == "" {
			fmt.Fprintf(os.Stderr, "Error: %s environment variable is not set\n", apiKeyEnvVar)
			os.Exit(1)
		}
	}

	aiConfig := ai.Config{
		Type:        providerType,
		APIKey:      apiKey,
		Model:       modelStr, // Use specified model or default
		Temperature: 0.0,
	}

	// Start CLI - it will handle connection and initialization
	cliApp := cli.NewApp(dbConfig, aiConfig)

	if err := cliApp.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
