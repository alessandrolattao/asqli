package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/alessandrolattao/asqli/internal/features/execution"
	"github.com/alessandrolattao/asqli/internal/features/query"
	"github.com/alessandrolattao/asqli/internal/features/schema"
	"github.com/alessandrolattao/asqli/internal/infrastructure/ai"
	"github.com/alessandrolattao/asqli/internal/infrastructure/config"
	"github.com/alessandrolattao/asqli/internal/infrastructure/database"
	"github.com/alessandrolattao/asqli/internal/infrastructure/database/adapters"
	tea "github.com/charmbracelet/bubbletea"
)

// connectDatabaseCmd connects to the database asynchronously
func connectDatabaseCmd(dbConfig adapters.Config, aiConfig ai.Config, timeoutConfig config.TimeoutConfig) tea.Cmd {
	return func() tea.Msg {
		// Connect to database
		dbConn, err := database.Open(dbConfig, timeoutConfig)
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
func fetchSchemaCmd(s *schema.Service, timeoutConfig config.TimeoutConfig) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutConfig.SchemaFetch)
		defer cancel()

		schema, err := s.Get(ctx)
		return schemaMsg{schema: schema, err: err}
	}
}

// generateSQLCmd generates SQL from natural language prompt asynchronously
func generateSQLCmd(s *query.Service, timeoutConfig config.TimeoutConfig, prompt, schema string, queryHistory []QueryHistory, selectedColumn string, selectedValue any) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutConfig.AIGeneration)
		defer cancel()

		// Convert QueryHistory to query.History
		history := make([]query.History, len(queryHistory))
		for i, qh := range queryHistory {
			history[i] = query.History{
				Prompt: qh.Prompt,
				SQL:    qh.SQL,
			}
		}

		sql, err := s.Generate(ctx, prompt, schema, history, selectedColumn, selectedValue)
		return sqlGeneratedMsg{sql: sql, err: err}
	}
}

// executeQueryCmd executes a SQL query asynchronously
func executeQueryCmd(s *execution.Service, timeoutConfig config.TimeoutConfig, query string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutConfig.DatabaseQuery)
		defer cancel()

		result, err := s.Execute(ctx, query)
		return queryExecutedMsg{result: result, err: err}
	}
}
