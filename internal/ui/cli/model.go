package cli

import (
	"github.com/alessandrolattao/asqli/internal/features/execution"
	"github.com/alessandrolattao/asqli/internal/features/query"
	"github.com/alessandrolattao/asqli/internal/features/schema"
	"github.com/alessandrolattao/asqli/internal/infrastructure/ai"
	"github.com/alessandrolattao/asqli/internal/infrastructure/config"
	"github.com/alessandrolattao/asqli/internal/infrastructure/database"
	"github.com/alessandrolattao/asqli/internal/infrastructure/database/adapters"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the Bubble Tea model for the CLI application
type Model struct {
	// Configuration
	dbConfig      adapters.Config
	aiConfig      ai.Config
	timeoutConfig config.TimeoutConfig

	// Services (initialized after connection)
	queryService     *query.Service
	executionService *execution.Service
	schemaService    *schema.Service

	// Database connection (to close on exit)
	dbConn     *database.Connection
	aiProvider ai.Provider

	// UI components
	spinner   spinner.Model
	textInput textinput.Model
	list      list.Model
	table     *Table

	// Application state
	state         state
	schema        string
	err           error
	currentPrompt string
	generatedSQL  string
	currentUsage  ai.UsageMetadata

	// Current result display
	currentResult *execution.Result
	currentError  error

	// Status message
	statusMessage string

	// History (for display and navigation)
	history      []string
	historyIndex int

	// Query history (for AI context with prompt + SQL)
	queryHistory []QueryHistory

	// Terminal dimensions
	width  int
	height int
}

// NewModel creates a new Bubble Tea model with configuration
func NewModel(
	dbConfig adapters.Config,
	aiConfig ai.Config,
	timeoutConfig config.TimeoutConfig,
) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	ti := textinput.New()
	ti.Placeholder = "Enter your query..."
	ti.Focus()
	ti.CharLimit = 500
	ti.Width = 80

	// Create history list
	delegate := simpleDelegate{}
	historyList := list.New([]list.Item{}, delegate, 80, 20)
	historyList.Title = "Query History"
	historyList.SetShowStatusBar(false)
	historyList.SetShowHelp(false)
	historyList.SetFilteringEnabled(false)
	historyList.Styles.Title = logoStyle

	return Model{
		dbConfig:      dbConfig,
		aiConfig:      aiConfig,
		timeoutConfig: timeoutConfig,
		state:         stateConnecting,
		spinner:       s,
		textInput:     ti,
		list:          historyList,
		history:       loadHistory(),
		historyIndex:  -1,
	}
}

// Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		connectDatabaseCmd(m.dbConfig, m.aiConfig, m.timeoutConfig),
		m.spinner.Tick,
	)
}
