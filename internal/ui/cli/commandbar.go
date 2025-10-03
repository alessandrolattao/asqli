package cli

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// CommandBar represents the 6-line command bar component
type CommandBar struct {
	width         int
	state         state
	spinner       spinner.Model
	textInput     textinput.Model
	statusMessage string
	generatedSQL  string
}

// NewCommandBar creates a new command bar component
func NewCommandBar(width int, currentState state, sp spinner.Model, ti textinput.Model, statusMsg string, sql string) CommandBar {
	return CommandBar{
		width:         width,
		state:         currentState,
		spinner:       sp,
		textInput:     ti,
		statusMessage: statusMsg,
		generatedSQL:  sql,
	}
}

// View renders the 6-line command bar
// Line 1: Generated SQL (when available)
// Line 2: Divider
// Line 3: Status (with spinner when active, empty when idle)
// Line 4: Text input
// Line 5: Divider
// Line 6: Help text
func (c CommandBar) View() string {
	dividerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB6C1"))
	divider := dividerStyle.Render(strings.Repeat("─", c.width))

	// SQL line (1st line) - shows the generated SQL query when available
	var sqlLine string
	if c.generatedSQL != "" {
		sqlStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#98C379"))
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB6C1")).Bold(true)

		// Normalize SQL to single line (replace newlines with spaces)
		displaySQL := strings.ReplaceAll(c.generatedSQL, "\n", " ")
		// Remove multiple spaces
		displaySQL = strings.Join(strings.Fields(displaySQL), " ")

		// Truncate SQL if too long
		maxSQLWidth := c.width - 6 // Reserve space for "SQL: " label
		if len(displaySQL) > maxSQLWidth {
			displaySQL = displaySQL[:maxSQLWidth-3] + "..."
		}

		sqlLine = labelStyle.Render("SQL: ") + sqlStyle.Render(displaySQL)
	}

	// Status line (3rd line) - shows spinner when loading, help when idle, or status message
	var statusLine string

	// Show status message if available (takes priority)
	if c.statusMessage != "" {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
		statusLine = successStyle.Render("✓ " + c.statusMessage)
	} else {
		switch c.state {
		case stateConnecting:
			statusLine = c.spinner.View() + " " + subtleStyle.Render("Connecting")
		case stateLoadingSchema:
			statusLine = c.spinner.View() + " " + subtleStyle.Render("Loading schema")
		case stateThinking:
			statusLine = c.spinner.View() + " " + subtleStyle.Render("Thinking")
		case stateExecuting:
			statusLine = c.spinner.View() + " " + subtleStyle.Render("Executing query")
		case stateConfirming:
			statusLine = subtleStyle.Render("⚠ DANGEROUS QUERY - Proceed? (y/n)")
		case stateReady:
			statusLine = subtleStyle.Render("Use # for raw SQL or ask me anything • Type 'exit' to quit")
		default:
			statusLine = ""
		}
	}

	// Help line (6th line)
	helpText := subtleStyle.Render("↑↓←→: table navigation • Ctrl+↑↓: history • Ctrl+r: history list • Ctrl+p: prompt info • Ctrl+c: copy as TSV • Esc: prompt clear • Ctrl+q: quit")

	return sqlLine + "\n" +
		divider + "\n" +
		statusLine + "\n" +
		c.textInput.View() + "\n" +
		divider + "\n" +
		helpText
}
