package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	_ "github.com/alessandrolattao/sqlai/internal/infrastructure/ai/claude" // Register Claude provider
	_ "github.com/alessandrolattao/sqlai/internal/infrastructure/ai/gemini" // Register Gemini provider
	_ "github.com/alessandrolattao/sqlai/internal/infrastructure/ai/ollama" // Register Ollama provider
	_ "github.com/alessandrolattao/sqlai/internal/infrastructure/ai/openai" // Register OpenAI provider
	"github.com/alessandrolattao/sqlai/internal/infrastructure/config"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database/adapters"
	"github.com/alessandrolattao/sqlai/internal/ui/cli"
)

func main() {
	// Command flags
	versionFlag := flag.Bool("version", false, "Print the version of sqlai")

	// AI Provider
	providerStr := flag.String("provider", "openai", "AI provider (openai, claude, gemini, ollama)")
	modelStr := flag.String("model", "", "AI model to use (defaults to provider's default model)")

	// Database type
	dbTypeStr := flag.String("dbtype", "postgres", "Database type (postgres, mysql, sqlite)")

	// Connection string
	connStr := flag.String("connection", "", "Database connection string (if provided, other connection params are ignored)")

	// PostgreSQL/MySQL shared connection parameters
	dbHost := flag.String("host", "", "Database host")
	dbPort := flag.Int("port", 5432, "Database port")
	dbUser := flag.String("user", "", "Database username")
	dbPassword := flag.String("password", "", "Database password")
	dbName := flag.String("db", "", "Database name")

	// PostgreSQL specific
	dbSSLMode := flag.String("sslmode", "disable", "PostgreSQL SSL mode")

	// MySQL specific
	parseTime := flag.Bool("parsetime", true, "MySQL: parse time values to Go time.Time")

	// SQLite specific
	dbFile := flag.String("file", "", "SQLite database file path")

	// Parse flags
	flag.Parse()

	// Version command takes precedence
	if *versionFlag {
		fmt.Printf("SQL AI v%s\n", Version)
		fmt.Println("WARNING: This software is in early development and not ready for production use.")
		return
	}

	// Determine database driver type
	var dbType adapters.DriverType
	switch *dbTypeStr {
	case "postgres", "postgresql":
		dbType = adapters.PostgreSQL
	case "mysql", "mariadb":
		dbType = adapters.MySQL
	case "sqlite", "sqlite3":
		dbType = adapters.SQLite
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported database type '%s'. Supported types: postgres, mysql, sqlite\n", *dbTypeStr)
		os.Exit(1)
	}

	// Create the appropriate configuration
	var cfg adapters.Config
	cfg.DriverType = dbType

	if *connStr != "" {
		// If connection string is provided, use it directly
		cfg.ConnectionString = *connStr
		runQuerySession(cfg, *providerStr, *modelStr)
		return
	}

	// Otherwise, use the appropriate config builder based on database type
	switch dbType {
	case adapters.PostgreSQL:
		cfg.Host = *dbHost
		cfg.Port = *dbPort
		cfg.User = *dbUser
		cfg.Password = *dbPassword
		cfg.DBName = *dbName
		cfg.SSLMode = *dbSSLMode

		// Try to get password from .pgpass if not provided
		if cfg.Password == "" && cfg.Host != "" && cfg.User != "" && cfg.DBName != "" {
			if pgpassPassword := config.GetPasswordForConfig(cfg.Host, cfg.Port, cfg.DBName, cfg.User); pgpassPassword != "" {
				cfg.Password = pgpassPassword
			}
		}
	case adapters.MySQL:
		cfg.Host = *dbHost
		cfg.Port = *dbPort
		cfg.User = *dbUser
		cfg.Password = *dbPassword
		cfg.DBName = *dbName
		cfg.ParseTime = *parseTime

		// Try to get password from .pgpass if not provided (works for MySQL too)
		if cfg.Password == "" && cfg.Host != "" && cfg.User != "" && cfg.DBName != "" {
			if pgpassPassword := config.GetPasswordForConfig(cfg.Host, cfg.Port, cfg.DBName, cfg.User); pgpassPassword != "" {
				cfg.Password = pgpassPassword
			}
		}
	case adapters.SQLite:
		if *dbFile == "" && *dbName != "" {
			// Use dbName as file path for SQLite if dbFile is not specified
			*dbFile = *dbName
		}
		if *dbFile == "" {
			fmt.Fprintf(os.Stderr, "Error: SQLite database file path not specified. Use --file parameter.\n")
			os.Exit(1)
		}
		cfg.FilePath = *dbFile
	}

	// Default to query command
	runQuerySession(cfg, *providerStr, *modelStr)
}

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
