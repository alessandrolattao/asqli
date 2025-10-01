package adapters

// DriverType represents the supported database driver types
type DriverType string

const (
	// PostgreSQL driver
	PostgreSQL DriverType = "postgres"
	// MySQL driver
	MySQL DriverType = "mysql"
	// SQLite driver
	SQLite DriverType = "sqlite3"
)

// Config represents database configuration parameters
type Config struct {
	// DriverType specifies which database driver to use
	DriverType DriverType

	// ConnectionString is the full connection string. If provided, it takes precedence
	// over individual connection parameters
	ConnectionString string

	// Individual connection parameters
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	// Additional parameters specific to certain databases
	SSLMode   string // For PostgreSQL
	FilePath  string // For SQLite
	ParseTime bool   // For MySQL
}

// TableDefinition contains information about a database table
type TableDefinition struct {
	Name        string
	Columns     []ColumnDefinition
	Constraints []ConstraintDefinition
}

// ColumnDefinition contains information about a table column
type ColumnDefinition struct {
	Name       string
	Type       string
	Nullable   bool
	Default    string
	IsPrimary  bool
	IsAutoIncr bool
}

// ConstraintDefinition contains information about table constraints
type ConstraintDefinition struct {
	Name              string
	Type              string // PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK, etc.
	Definition        string
	ReferencedTable   string
	ReferencedColumns []string
}
