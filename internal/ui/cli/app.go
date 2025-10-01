// Package cli provides an interactive command-line interface for SQL AI.
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/alessandrolattao/sqlai/internal/features/execution"
	"github.com/alessandrolattao/sqlai/internal/features/query"
	"github.com/alessandrolattao/sqlai/internal/features/schema"
	"github.com/chzyer/readline"
)

// App represents the CLI application
type App struct {
	queryService     *query.Service
	executionService *execution.Service
	schemaService    *schema.Service
}

// NewApp creates a new CLI application
func NewApp(
	queryService *query.Service,
	executionService *execution.Service,
	schemaService *schema.Service,
) *App {
	return &App{
		queryService:     queryService,
		executionService: executionService,
		schemaService:    schemaService,
	}
}

// Start begins the CLI interaction loop
func (a *App) Start() {
	fmt.Printf("%s%sWelcome to SQL AI!%s\n", Bold, ColorCyan, ColorReset)
	fmt.Printf("Enter your natural language query, or prefix with '%s#%s' for raw SQL.\n", ColorYellow, ColorReset)
	fmt.Printf("Type '%sexit%s' to quit.\n\n", ColorYellow, ColorReset)
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
	defer func() { _ = rl.Close() }()

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

		// Check if it's a raw SQL query (prefixed with #)
		if strings.HasPrefix(line, "#") {
			rawSQL := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if rawSQL != "" {
				a.processRawQuery(rawSQL)
			}
			continue
		}

		// Process the query with AI
		a.processQuery(line)
	}
}

// processRawQuery executes a raw SQL query without AI generation
func (a *App) processRawQuery(sqlQuery string) {
	ctx := context.Background()

	fmt.Printf("\n%s%sExecuting raw SQL query...%s\n", Bold, ColorPurple, ColorReset)
	fmt.Printf("%s%s%s\n\n", ColorCyan, sqlQuery, ColorReset)

	// Check if query requires confirmation
	if a.queryService.IsDangerous(sqlQuery) {
		if !a.confirmQuery() {
			return
		}
	}

	fmt.Printf("%s%sExecuting...%s\n", Bold, ColorBlue, ColorReset)

	// Execute query
	result, err := a.executionService.Execute(ctx, sqlQuery)
	if err != nil {
		fmt.Printf("%s%sError executing query: %v%s\n", Bold, ColorRed, err, ColorReset)
		return
	}

	// Display results
	a.displayResults(sqlQuery, result)
}

// processQuery handles a single query from start to finish
func (a *App) processQuery(prompt string) {
	ctx := context.Background()

	fmt.Printf("\n%s%sGenerating SQL query...%s\n", Bold, ColorBlue, ColorReset)

	// 1. Get database schema
	schemaStr, err := a.schemaService.Get(ctx)
	if err != nil {
		fmt.Printf("%sWarning: Could not fetch database schema: %v%s\n", ColorYellow, err, ColorReset)
	}

	// 2. Generate SQL query
	sql, err := a.queryService.Generate(ctx, prompt, schemaStr)
	if err != nil {
		fmt.Printf("%s%sError generating query: %v%s\n", Bold, ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("\n%s%sGenerated SQL:%s %s%s%s\n\n",
		Bold, ColorPurple, ColorReset, ColorCyan, sql.Query, ColorReset)

	// 3. Check if query requires confirmation
	if a.queryService.IsDangerous(sql.Query) {
		if !a.confirmQuery() {
			return
		}
	}

	fmt.Printf("%s%sExecuting query...%s\n", Bold, ColorBlue, ColorReset)

	// 4. Execute query
	result, err := a.executionService.Execute(ctx, sql.Query)
	if err != nil {
		fmt.Printf("%s%sError executing query: %v%s\n", Bold, ColorRed, err, ColorReset)
		return
	}

	// 5. Display results
	a.displayResults(sql.Query, result)
}

// confirmQuery prompts the user for confirmation
func (a *App) confirmQuery() bool {
	fmt.Printf("%s%sWARNING: This query might modify data.%s\n", Bold, ColorYellow, ColorReset)
	fmt.Printf("Do you want to proceed? (y/N): ")

	var confirm string
	_, _ = fmt.Scanln(&confirm)

	confirm = strings.ToLower(strings.TrimSpace(confirm))
	if confirm != "y" && confirm != "yes" {
		fmt.Printf("%s%sQuery execution cancelled.%s\n", Bold, ColorRed, ColorReset)
		return false
	}

	return true
}

// displayResults formats and displays query results
func (a *App) displayResults(queryStr string, result *execution.Result) {
	queryType := getQueryType(queryStr)
	isSelectQuery := queryType == "SELECT"

	if isSelectQuery {
		// For SELECT queries, format and display result table
		formatted := formatAsTable(result.Rows, result.Columns)
		fmt.Println(formatted)
		return
	}

	// For non-SELECT queries, show appropriate success message
	fmt.Print(formatSuccessMessage(queryType))

	// If we have any result data (like affected rows), show it
	if len(result.Rows) > 0 {
		formatted := formatAsTable(result.Rows, result.Columns)
		fmt.Println(formatted)
	}
}

// getQueryType extracts the type of SQL query (SELECT, INSERT, etc.)
func getQueryType(query string) string {
	trimmedQuery := strings.TrimSpace(strings.ToUpper(query))

	// Extract the first word as the query type
	parts := strings.Fields(trimmedQuery)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
