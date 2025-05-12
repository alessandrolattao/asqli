package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alessandrolattao/sqlai/internal/pkg/config"
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
	"github.com/alessandrolattao/sqlai/internal/pkg/ui"
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
	var dbType database.DriverType
	switch *dbTypeStr {
	case "postgres", "postgresql":
		dbType = database.PostgreSQL
	case "mysql", "mariadb":
		dbType = database.MySQL
	case "sqlite", "sqlite3":
		dbType = database.SQLite
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported database type '%s'. Supported types: postgres, mysql, sqlite\n", *dbTypeStr)
		os.Exit(1)
	}

	// Create the appropriate configuration
	var cfg *config.Config

	if *connStr != "" {
		// If connection string is provided, use it directly
		cfg = config.NewCustomConfig(dbType, *connStr)
	} else {
		// Otherwise, use the appropriate config builder based on database type
		switch dbType {
		case database.PostgreSQL:
			cfg = config.NewConfig(*dbHost, *dbUser, *dbPassword, *dbName, *dbSSLMode, *dbPort)
		case database.MySQL:
			cfg = config.NewMySQLConfig(*dbHost, *dbUser, *dbPassword, *dbName, *dbPort, *parseTime)
		case database.SQLite:
			if *dbFile == "" && *dbName != "" {
				// Use dbName as file path for SQLite if dbFile is not specified
				*dbFile = *dbName
			}
			if *dbFile == "" {
				fmt.Fprintf(os.Stderr, "Error: SQLite database file path not specified. Use --file parameter.\n")
				os.Exit(1)
			}
			cfg = config.NewSQLiteConfig(*dbFile)
		}
	}

	// Default to query command
	runQuerySession(cfg.DB)
}

func runQuerySession(dbConfig database.Config) {
	// Create a new CLI interface
	cli, err := ui.NewCLI(dbConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Start the CLI
	cli.Start()
}

