package ui

import (
	"database/sql"

	"github.com/alessandrolattao/sqlai/internal/pkg/ai"
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	Bold        = "\033[1m"
)

// CLI represents a command-line interface for SQL AI
type CLI struct {
	db   *sql.DB
	ai   *ai.Client
	conn *database.Connection
}