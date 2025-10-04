package main

import (
	"fmt"
	"os"

	"github.com/alessandrolattao/asqli/internal/infrastructure/config"
	"github.com/alessandrolattao/asqli/internal/infrastructure/database/adapters"
)

// buildDatabaseConfig creates a database configuration from flags
func buildDatabaseConfig(flags *Flags) adapters.Config {
	// Determine database driver type
	var dbType adapters.DriverType
	switch flags.DBType {
	case "postgres", "postgresql":
		dbType = adapters.PostgreSQL
	case "mysql", "mariadb":
		dbType = adapters.MySQL
	case "sqlite", "sqlite3":
		dbType = adapters.SQLite
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported database type '%s'. Supported types: postgres, mysql, sqlite\n", flags.DBType)
		os.Exit(1)
	}

	// Create the appropriate configuration
	cfg := adapters.Config{
		DriverType: dbType,
	}

	// If connection string is provided, use it directly
	if flags.Connection != "" {
		cfg.ConnectionString = flags.Connection
		return cfg
	}

	// Otherwise, build config based on database type
	switch dbType {
	case adapters.PostgreSQL:
		cfg.Host = flags.Host
		cfg.Port = flags.Port
		cfg.User = flags.User
		cfg.Password = flags.Password
		cfg.DBName = flags.DBName
		cfg.SSLMode = flags.SSLMode

		// Try to get password from .pgpass if not provided
		if cfg.Password == "" && cfg.Host != "" && cfg.User != "" && cfg.DBName != "" {
			if pgpassPassword := config.GetPasswordForConfig(cfg.Host, cfg.Port, cfg.DBName, cfg.User); pgpassPassword != "" {
				cfg.Password = pgpassPassword
			}
		}

	case adapters.MySQL:
		cfg.Host = flags.Host
		cfg.Port = flags.Port
		cfg.User = flags.User
		cfg.Password = flags.Password
		cfg.DBName = flags.DBName
		cfg.ParseTime = flags.ParseTime

		// Try to get password from .pgpass if not provided (works for MySQL too)
		if cfg.Password == "" && cfg.Host != "" && cfg.User != "" && cfg.DBName != "" {
			if pgpassPassword := config.GetPasswordForConfig(cfg.Host, cfg.Port, cfg.DBName, cfg.User); pgpassPassword != "" {
				cfg.Password = pgpassPassword
			}
		}

	case adapters.SQLite:
		file := flags.File
		if file == "" && flags.DBName != "" {
			// Use dbName as file path for SQLite if file is not specified
			file = flags.DBName
		}
		if file == "" {
			fmt.Fprintf(os.Stderr, "Error: SQLite database file path not specified. Use --file parameter.\n")
			os.Exit(1)
		}
		cfg.FilePath = file
	}

	return cfg
}
