package main

import (
	_ "github.com/alessandrolattao/asqli/internal/infrastructure/ai/claude" // Register Claude provider
	_ "github.com/alessandrolattao/asqli/internal/infrastructure/ai/gemini" // Register Gemini provider
	_ "github.com/alessandrolattao/asqli/internal/infrastructure/ai/ollama" // Register Ollama provider
	_ "github.com/alessandrolattao/asqli/internal/infrastructure/ai/openai" // Register OpenAI provider
)

func main() {
	// Parse command-line flags
	flags := ParseFlags()

	// Handle version command
	if flags.Version {
		handleVersion()
		return
	}

	// Handle update command
	if flags.Update {
		handleUpdate()
		return
	}

	// Build configurations
	dbConfig := buildDatabaseConfig(flags)
	timeoutConfig := buildTimeoutConfig(flags)

	// Start query session
	runQuerySession(dbConfig, timeoutConfig, flags.Provider, flags.Model)
}
