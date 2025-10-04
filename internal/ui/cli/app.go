// Package cli provides an interactive command-line interface for ASQLI using Bubble Tea.
package cli

import (
	"github.com/alessandrolattao/asqli/internal/infrastructure/ai"
	"github.com/alessandrolattao/asqli/internal/infrastructure/config"
	"github.com/alessandrolattao/asqli/internal/infrastructure/database/adapters"
	tea "github.com/charmbracelet/bubbletea"
)

// App represents the CLI application
type App struct {
	dbConfig      adapters.Config
	aiConfig      ai.Config
	timeoutConfig config.TimeoutConfig
}

// NewApp creates a new CLI application
func NewApp(
	dbConfig adapters.Config,
	aiConfig ai.Config,
	timeoutConfig config.TimeoutConfig,
) *App {
	return &App{
		dbConfig:      dbConfig,
		aiConfig:      aiConfig,
		timeoutConfig: timeoutConfig,
	}
}

// Start begins the Bubble Tea interactive loop
func (a *App) Start() error {
	// Create Bubble Tea model
	m := NewModel(a.dbConfig, a.aiConfig, a.timeoutConfig)

	// Create program WITH alternate screen for full UI rendering
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Enable alternate screen mode
		tea.WithMouseCellMotion(), // Enable mouse support for table navigation
	)

	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Check if the model has an error (e.g., connection failed)
	if m, ok := finalModel.(Model); ok && m.err != nil {
		return m.err
	}

	return nil
}
