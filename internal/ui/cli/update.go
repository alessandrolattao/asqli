package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case spinner.TickMsg:
		// Update spinner only when loading
		if m.state == stateConnecting || m.state == stateLoadingSchema || m.state == stateThinking || m.state == stateExecuting {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tea.KeyMsg:
		// Clear status message on any key press (except copy)
		if m.statusMessage != "" && msg.String() != "ctrl+c" {
			m.statusMessage = ""
		}

		// Handle history view separately
		if m.state == stateHistory {
			switch msg.String() {
			case "esc":
				m.state = stateReady
				return m, nil
			case "enter":
				if selectedItem := m.list.SelectedItem(); selectedItem != nil {
					histItem := selectedItem.(historyItem)
					m.textInput.SetValue(histItem.query)
					m.state = stateReady
					return m, nil
				}
			case "ctrl+d":
				// Reset history
				m.history = []string{}
				if err := saveHistory(m.history); err != nil {
					m.statusMessage = "Warning: Failed to save cleared history"
				} else {
					m.statusMessage = "History cleared!"
				}
				m.list.SetItems([]list.Item{})
				m.state = stateReady
				return m, nil
			case "up", "down", "k", "j":
				// Only allow navigation keys
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			default:
				// Ignore all other keys in history view
				return m, nil
			}
		}

		// Handle info view separately
		if m.state == stateInfo {
			switch msg.String() {
			case "esc":
				// Exit info view
				m.state = stateReady
				return m, nil
			default:
				// Ignore all other keys in info view
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+q":
			// Save history before quitting (best effort, don't block quit)
			if err := saveHistory(m.history); err != nil {
				// Log to stderr since we're quitting anyway
				fmt.Fprintf(os.Stderr, "Warning: Failed to save history: %v\n", err)
			}
			// Clean up resources before quitting (best effort)
			if m.dbConn != nil {
				if err := m.dbConn.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to close database: %v\n", err)
				}
			}
			if m.aiProvider != nil {
				if err := m.aiProvider.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to close AI provider: %v\n", err)
				}
			}
			return m, tea.Quit

		case "ctrl+c":
			// Copy table to clipboard
			if m.state == stateReady && m.table != nil {
				if err := m.table.CopyToClipboard(); err != nil {
					m.statusMessage = "Failed to copy: " + err.Error()
					return m, nil
				}
				m.statusMessage = "Table copied to clipboard!"
				return m, nil
			}

		case "ctrl+r":
			// Open history view (only when ready)
			if m.state == stateReady {
				return m.openHistoryView()
			}

		case "ctrl+p":
			// Open info view (only when ready and we have query history)
			if m.state == stateReady && len(m.queryHistory) > 0 {
				m.state = stateInfo
				return m, nil
			}

		case "ctrl+up":
			// Navigate history up (more recent)
			if m.state == stateReady {
				return m.navigateHistory(1)
			}

		case "ctrl+down":
			// Navigate history down (older)
			if m.state == stateReady {
				return m.navigateHistory(-1)
			}

		case "up":
			// Navigate table up when result is shown
			if m.state == stateReady && m.table != nil {
				m.table.MoveUp()
				return m, nil
			}

		case "down":
			// Navigate table down when result is shown
			if m.state == stateReady && m.table != nil {
				m.table.MoveDown()
				return m, nil
			}

		case "left":
			// Navigate table left when result is shown
			if m.state == stateReady && m.table != nil {
				m.table.MoveLeft()
				return m, nil
			}

		case "right":
			// Navigate table right when result is shown
			if m.state == stateReady && m.table != nil {
				m.table.MoveRight()
				return m, nil
			}

		case "esc":
			// Clear input field when ready
			if m.state == stateReady {
				m.textInput.SetValue("")
				m.historyIndex = -1
				return m, nil
			}
			// Cancel confirmation
			if m.state == stateConfirming {
				m.state = stateReady
				m.generatedSQL = ""
				return m, nil
			}

		case "enter":
			// Submit query when ready
			if m.state == stateReady {
				return m.handleSubmit()
			}

		case "y":
			// Confirm dangerous query
			if m.state == stateConfirming {
				return m.handleConfirmYes()
			}

		case "n":
			// Cancel dangerous query
			if m.state == stateConfirming {
				m.state = stateReady
				m.generatedSQL = ""
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = msg.Width - TablePaddingHorizontal
		m.list.SetWidth(msg.Width - TablePaddingHorizontal)
		m.list.SetHeight(msg.Height - TablePaddingVertical)

		// Update table size if exists
		if m.table != nil {
			tableHeight := m.height - CommandBarHeight - TablePaddingVertical
			tableWidth := m.width - TablePaddingHorizontal
			m.table.SetSize(tableWidth, tableHeight)
		}

		return m, nil

	case connectionMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		// Connection successful, store services
		m.dbConn = msg.dbConn
		m.aiProvider = msg.aiProvider
		m.schemaService = msg.schemaService
		m.queryService = msg.queryService
		m.executionService = msg.executionService
		m.state = stateLoadingSchema
		// Now fetch schema
		return m, tea.Batch(
			fetchSchemaCmd(m.schemaService),
			m.spinner.Tick,
		)

	case schemaMsg:
		if msg.err != nil {
			m.err = msg.err
			// Clean up before quitting (best effort on error path)
			if m.dbConn != nil {
				if err := m.dbConn.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to close database: %v\n", err)
				}
			}
			if m.aiProvider != nil {
				if err := m.aiProvider.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to close AI provider: %v\n", err)
				}
			}
			return m, tea.Quit
		}
		m.schema = msg.schema
		m.state = stateReady
		return m, nil

	case sqlGeneratedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.currentError = nil
			m.currentResult = nil
			m.table = nil
			m.statusMessage = "✗ " + msg.err.Error()
			m.state = stateReady
			return m, nil
		}

		m.generatedSQL = msg.sql.Query
		m.currentUsage = msg.sql.Usage
		m.err = nil // Clear any previous generation errors

		// Check if query is dangerous
		if m.queryService.IsDangerous(m.generatedSQL) {
			m.state = stateConfirming
			return m, nil
		}

		// Execute directly if not dangerous
		m.state = stateExecuting
		return m, tea.Batch(
			executeQueryCmd(m.executionService, m.generatedSQL),
			m.spinner.Tick,
		)

	case queryExecutedMsg:
		// Store for history (avoid consecutive duplicates)
		if m.currentPrompt != "" {
			// Only add if it's different from the last entry
			if len(m.history) == 0 || m.history[len(m.history)-1] != m.currentPrompt {
				m.history = append(m.history, m.currentPrompt)
				if err := saveHistory(m.history); err != nil {
					// Log error but don't interrupt the flow
					fmt.Fprintf(os.Stderr, "Warning: Failed to save history: %v\n", err)
				}
			}
		}

		// Store complete query history (prompt + SQL + usage) for AI context and debug
		if m.currentPrompt != "" && m.generatedSQL != "" {
			m.queryHistory = append(m.queryHistory, QueryHistory{
				Prompt: m.currentPrompt,
				SQL:    m.generatedSQL,
				Usage:  m.currentUsage,
			})
			// Keep only last N queries for context
			if len(m.queryHistory) > MaxQueryHistory {
				m.queryHistory = m.queryHistory[len(m.queryHistory)-MaxQueryHistory:]
			}
		}

		// Store result in model for rendering
		m.currentResult = msg.result
		m.currentError = msg.err
		m.err = nil // Clear any previous generation errors

		// Set status message based on result
		if msg.err != nil {
			m.statusMessage = "✗ " + msg.err.Error()
		} else if msg.result != nil {
			m.statusMessage = fmt.Sprintf("✓ Query executed successfully (%d rows)", len(msg.result.Rows))
		}

		// Create table if result has rows
		if msg.result != nil && len(msg.result.Rows) > 0 {
			// Reserve space for command bar and padding
			tableHeight := m.height - CommandBarHeight - TablePaddingVertical
			tableWidth := m.width - TablePaddingHorizontal
			m.table = NewTable(msg.result, tableWidth, tableHeight)
		} else {
			m.table = nil
		}

		// Clear current state (keep generatedSQL to display in command bar)
		m.currentPrompt = ""
		m.err = nil
		m.historyIndex = -1
		m.state = stateReady

		return m, nil
	}

	// Update text input when ready
	if m.state == stateReady {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// navigateHistory navigates through command history
func (m Model) navigateHistory(direction int) (Model, tea.Cmd) {
	if len(m.history) == 0 {
		return m, nil
	}

	newIndex := m.historyIndex + direction
	if newIndex < -1 {
		newIndex = -1
	}
	if newIndex >= len(m.history) {
		newIndex = len(m.history) - 1
	}

	m.historyIndex = newIndex
	if m.historyIndex == -1 {
		m.textInput.SetValue("")
	} else {
		m.textInput.SetValue(m.history[len(m.history)-1-m.historyIndex])
	}

	return m, nil
}

// openHistoryView opens the history list view
func (m Model) openHistoryView() (Model, tea.Cmd) {
	// Convert history slice to list items (reverse order - most recent first)
	items := make([]list.Item, len(m.history))
	for i, query := range m.history {
		items[len(m.history)-1-i] = historyItem{query: query}
	}

	m.list.SetItems(items)
	m.state = stateHistory

	return m, nil
}

// handleSubmit processes query submission
func (m Model) handleSubmit() (Model, tea.Cmd) {
	query := strings.TrimSpace(m.textInput.Value())

	if query == "" {
		return m, nil
	}

	// Check for exit commands
	if query == "exit" || query == "quit" {
		// Save history before quitting (best effort)
		if err := saveHistory(m.history); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save history: %v\n", err)
		}
		// Clean up resources (best effort)
		if m.dbConn != nil {
			if err := m.dbConn.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to close database: %v\n", err)
			}
		}
		if m.aiProvider != nil {
			if err := m.aiProvider.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to close AI provider: %v\n", err)
			}
		}
		return m, tea.Quit
	}

	// Clear input
	m.textInput.SetValue("")
	m.historyIndex = -1

	// Check for raw SQL (# prefix)
	if strings.HasPrefix(query, "#") {
		m.generatedSQL = strings.TrimSpace(strings.TrimPrefix(query, "#"))
		m.currentPrompt = query

		// Check if dangerous
		if m.queryService.IsDangerous(m.generatedSQL) {
			m.state = stateConfirming
			return m, nil
		}

		// Execute directly
		m.state = stateExecuting
		return m, tea.Batch(
			executeQueryCmd(m.executionService, m.generatedSQL),
			m.spinner.Tick,
		)
	}

	// Generate SQL with AI
	m.currentPrompt = query
	m.state = stateThinking

	// Get selected column and value if table exists
	var selectedColumn string
	var selectedValue any
	if m.table != nil {
		selectedColumn = m.table.GetSelectedColumn()
		selectedValue = m.table.GetSelectedValue()
	}

	return m, tea.Batch(
		generateSQLCmd(m.queryService, query, m.schema, m.queryHistory, selectedColumn, selectedValue),
		m.spinner.Tick,
	)
}

// handleConfirmYes proceeds with dangerous query execution
func (m Model) handleConfirmYes() (Model, tea.Cmd) {
	m.state = stateExecuting
	return m, tea.Batch(
		executeQueryCmd(m.executionService, m.generatedSQL),
		m.spinner.Tick,
	)
}
