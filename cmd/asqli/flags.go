package main

import "flag"

// Flags holds all command-line flags for the application
type Flags struct {
	// Command flags
	Version bool
	Update  bool

	// AI Provider flags
	Provider string
	Model    string

	// Database type
	DBType string

	// Connection string
	Connection string

	// PostgreSQL/MySQL shared connection parameters
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	// PostgreSQL specific
	SSLMode string

	// MySQL specific
	ParseTime bool

	// SQLite specific
	File string

	// Timeout settings (in seconds)
	TimeoutConnection int
	TimeoutQuery      int
	TimeoutSchema     int
	TimeoutAI         int
}

// ParseFlags parses command-line flags and returns a Flags struct
func ParseFlags() *Flags {
	f := &Flags{}

	// Command flags
	flag.BoolVar(&f.Version, "version", false, "Print the version of asqli")
	flag.BoolVar(&f.Update, "update", false, "Check for updates and update to the latest version")

	// AI Provider
	flag.StringVar(&f.Provider, "provider", "openai", "AI provider (openai, claude, gemini, ollama)")
	flag.StringVar(&f.Model, "model", "", "AI model to use (defaults to provider's default model)")

	// Database type
	flag.StringVar(&f.DBType, "dbtype", "postgres", "Database type (postgres, mysql, sqlite)")

	// Connection string
	flag.StringVar(&f.Connection, "connection", "", "Database connection string (if provided, other connection params are ignored)")

	// PostgreSQL/MySQL shared connection parameters
	flag.StringVar(&f.Host, "host", "", "Database host")
	flag.IntVar(&f.Port, "port", 5432, "Database port")
	flag.StringVar(&f.User, "user", "", "Database username")
	flag.StringVar(&f.Password, "password", "", "Database password")
	flag.StringVar(&f.DBName, "db", "", "Database name")

	// PostgreSQL specific
	flag.StringVar(&f.SSLMode, "sslmode", "disable", "PostgreSQL SSL mode")

	// MySQL specific
	flag.BoolVar(&f.ParseTime, "parsetime", true, "MySQL: parse time values to Go time.Time")

	// SQLite specific
	flag.StringVar(&f.File, "file", "", "SQLite database file path")

	// Timeout settings (in seconds, 0 = use default)
	flag.IntVar(&f.TimeoutConnection, "timeout-connection", 0, "Database connection timeout in seconds (default: 10)")
	flag.IntVar(&f.TimeoutQuery, "timeout-query", 0, "Database query execution timeout in seconds (default: 30)")
	flag.IntVar(&f.TimeoutSchema, "timeout-schema", 0, "Schema fetch timeout in seconds (default: 30)")
	flag.IntVar(&f.TimeoutAI, "timeout-ai", 0, "AI generation timeout in seconds (default: 60)")

	flag.Parse()

	return f
}
