package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/pkg/ai"
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
	"github.com/alessandrolattao/sqlai/internal/pkg/database/drivers"
	"github.com/chzyer/readline"
)

// NewCLI creates a new CLI instance
func NewCLI(dbConfig database.Config) (*CLI, error) {
	// Initialize AI client
	aiClient, err := ai.NewClient()
	if err != nil {
		return nil, err
	}

	// Create the appropriate adapter based on driver type
	var adapter database.Adapter

	switch dbConfig.DriverType {
	case database.PostgreSQL:
		adapter = &drivers.PostgreSQLAdapter{}
	case database.MySQL:
		adapter = &drivers.MySQLAdapter{}
	case database.SQLite:
		adapter = &drivers.SQLiteAdapter{}
	default:
		// Default to PostgreSQL for backward compatibility
		adapter = &drivers.PostgreSQLAdapter{}
	}

	// Connect to database
	conn, err := database.Open(dbConfig, adapter)
	if err != nil {
		return nil, err
	}

	return &CLI{
		db:   conn.DB,
		ai:   aiClient,
		conn: conn,
	}, nil
}

// Start begins the CLI interaction loop
func (c *CLI) Start() {
	fmt.Printf("%s%sWelcome to SQL AI!%s Enter your natural language query or type '%sexit%s' to quit.\n",
		Bold, ColorCyan, ColorReset, ColorYellow, ColorReset)
	fmt.Printf("%s%sWARNING: This software is in early development and not ready for production use.%s\n",
		Bold, ColorRed, ColorReset)
	fmt.Println(ColorBlue + "---------------------------------------------------------------------" + ColorReset)

	// Configure readline instance
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("%s%ssqlai%s %s>%s ", Bold, ColorGreen, ColorReset, ColorYellow, ColorReset),
		HistoryFile:     os.ExpandEnv("$HOME/.sqlai_history"),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	
	if err != nil {
		fmt.Printf("%s%sError initializing readline: %v%s\n", Bold, ColorRed, err, ColorReset)
		return
	}
	defer rl.Close()
	
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		
		// Trim leading/trailing whitespace
		line = strings.TrimSpace(line)
		
		if line == "exit" {
			fmt.Println(ColorCyan + "Goodbye!" + ColorReset)
			break
		}
		
		if line == "" {
			continue
		}
		
		// Process the query
		c.processQuery(line)
	}
}