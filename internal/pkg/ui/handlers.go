package ui

import (
	"fmt"
	"strings"
)

// processQuery handles a single query from start to finish
func (c *CLI) processQuery(prompt string) {
	fmt.Printf("\n%s%sGenerating SQL query...%s\n", Bold, ColorBlue, ColorReset)

	// Get database schema for better query generation
	schema, err := c.conn.GetDatabaseSchema()
	if err != nil {
		fmt.Printf("%sWarning: Could not fetch database schema: %v%s\n", ColorYellow, err, ColorReset)
	}

	// Generate query using AI with schema information
	query, err := c.ai.GenerateQueryWithSchema(prompt, schema)
	if err != nil {
		fmt.Printf("%s%sError generating query: %v%s\n", Bold, ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("\n%s%sGenerated SQL:%s %s%s%s\n\n",
		Bold, ColorPurple, ColorReset, ColorCyan, query, ColorReset)

	// Check if query is non-SELECT and prompt for confirmation if needed
	if shouldConfirmQuery(query) {
		if !confirmQuery() {
			return
		}
	}

	fmt.Printf("%s%sExecuting query...%s\n", Bold, ColorBlue, ColorReset)

	// Execute query
	data, columns, err := c.conn.ExecuteQuery(query)
	if err != nil {
		fmt.Printf("%s%sError executing query: %v%s\n", Bold, ColorRed, err, ColorReset)
		return
	}

	// Display appropriate output based on query type
	displayQueryResults(query, data, columns)
}

// shouldConfirmQuery checks if a query requires confirmation before execution
func shouldConfirmQuery(query string) bool {
	trimmedQuery := strings.TrimSpace(strings.ToUpper(query))
	return !strings.HasPrefix(trimmedQuery, "SELECT")
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

// confirmQuery prompts the user for confirmation and returns their choice
func confirmQuery() bool {
	fmt.Printf("%s%sWARNING: This query might modify data.%s\n", Bold, ColorYellow, ColorReset)
	fmt.Printf("Do you want to proceed? (y/N): ")
	
	var confirm string
	fmt.Scanln(&confirm)
	
	confirm = strings.ToLower(strings.TrimSpace(confirm))
	if confirm != "y" && confirm != "yes" {
		fmt.Printf("%s%sQuery execution cancelled.%s\n", Bold, ColorRed, ColorReset)
		return false
	}
	
	return true
}

// displayQueryResults formats and displays query results
func displayQueryResults(query string, data []map[string]any, columns []string) {
	queryType := getQueryType(query)
	isSelectQuery := queryType == "SELECT"
	
	if isSelectQuery {
		// For SELECT queries, format and display result table
		result := formatAsTable(data, columns)
		fmt.Println(result)
	} else {
		// For non-SELECT queries, show appropriate success message
		fmt.Print(formatSuccessMessage(queryType))
		
		// If we have any result data (like affected rows), show it
		if len(data) > 0 {
			result := formatAsTable(data, columns)
			fmt.Println(result)
		}
	}
}