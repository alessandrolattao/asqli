package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/alessandrolattao/sqlai/internal/features/execution"
	"github.com/alessandrolattao/sqlai/internal/features/query"
	"github.com/alessandrolattao/sqlai/internal/features/schema"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database/adapters"
	tea "github.com/charmbracelet/bubbletea"
)

// connectDatabaseCmd connects to the database asynchronously
func connectDatabaseCmd(dbConfig adapters.Config, aiConfig ai.Config) tea.Cmd {
	return func() tea.Msg {
		// Connect to database
		dbConn, err := database.Open(dbConfig)
		if err != nil {
			return connectionMsg{err: err}
		}

		// Initialize AI provider
		aiProvider, err := ai.NewProvider(aiConfig)
		if err != nil {
			// Cleanup: attempt to close database connection (best effort on error path)
			if closeErr := dbConn.Close(); closeErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to close database during cleanup: %v\n", closeErr)
			}
			return connectionMsg{err: err}
		}

		// Create services
		schemaService := schema.NewService(dbConn)
		queryService := query.NewService(aiProvider)
		executionService := execution.NewService(dbConn)

		return connectionMsg{
			dbConn:           dbConn,
			aiProvider:       aiProvider,
			schemaService:    schemaService,
			queryService:     queryService,
			executionService: executionService,
			err:              nil,
		}
	}
}

// fetchSchemaCmd fetches the database schema asynchronously
func fetchSchemaCmd(s *schema.Service) tea.Cmd {
	return func() tea.Msg {
		schema, err := s.Get(context.Background())
		return schemaMsg{schema: schema, err: err}
	}
}

// generateSQLCmd generates SQL from natural language prompt asynchronously
func generateSQLCmd(s *query.Service, prompt, schema string, queryHistory []QueryHistory, selectedColumn string, selectedValue any) tea.Cmd {
	return func() tea.Msg {
		// Convert QueryHistory to query.History
		history := make([]query.History, len(queryHistory))
		for i, qh := range queryHistory {
			history[i] = query.History{
				Prompt: qh.Prompt,
				SQL:    qh.SQL,
			}
		}

		sql, err := s.Generate(context.Background(), prompt, schema, history, selectedColumn, selectedValue)
		return sqlGeneratedMsg{sql: sql, err: err}
	}
}

// executeQueryCmd executes a SQL query asynchronously
func executeQueryCmd(s *execution.Service, query string) tea.Cmd {
	return func() tea.Msg {
		result, err := s.Execute(context.Background(), query)
		return queryExecutedMsg{result: result, err: err}
	}
}
