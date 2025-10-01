package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alessandrolattao/sqlai/internal/features/execution"
	"github.com/alessandrolattao/sqlai/internal/features/query"
	"github.com/alessandrolattao/sqlai/internal/features/schema"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	_ "github.com/alessandrolattao/sqlai/internal/infrastructure/ai/openai" // Register OpenAI provider
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database/adapters"
	"github.com/alessandrolattao/sqlai/internal/ui/cli"
)

const (
	// Version is the current version of sqlai
	Version = "0.1.0"
)

func main() {
	// Command flags
	versionFlag := flag.Bool("version", false, "Print the version of sqlai")

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
		runQuerySession(cfg)
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
	case adapters.MySQL:
		cfg.Host = *dbHost
		cfg.Port = *dbPort
		cfg.User = *dbUser
		cfg.Password = *dbPassword
		cfg.DBName = *dbName
		cfg.ParseTime = *parseTime
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
	runQuerySession(cfg)
}

func runQuerySession(dbConfig adapters.Config) {
	// 1. Initialize infrastructure

	// Database connection
	dbConn, err := database.Open(dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close database connection: %v\n", err)
		}
	}()

	// AI Provider
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: OPENAI_API_KEY environment variable is not set\n")
		os.Exit(1)
	}

	aiProvider, err := ai.NewProvider(ai.Config{
		Type:        ai.ProviderOpenAI,
		APIKey:      apiKey,
		Model:       "", // Use default
		Temperature: 0.0,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing AI provider: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := aiProvider.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close AI provider: %v\n", err)
		}
	}()

	// 2. Initialize features
	schemaService := schema.NewService(dbConn)
	queryService := query.NewService(aiProvider)
	executionService := execution.NewService(dbConn)

	// 3. Initialize UI
	cliApp := cli.NewApp(queryService, executionService, schemaService)

	// 4. Start CLI
	cliApp.Start()
}
