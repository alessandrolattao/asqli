package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// renderResultsArea renders the results display area
func (m Model) renderResultsArea(height int) string {
	// Create padding style
	paddingStyle := lipgloss.NewStyle().Padding(1, 2)

	if m.table != nil {
		// Show table with navigation and padding
		tableView := m.table.View()
		return paddingStyle.Render(tableView)
	}

	// Show error or success message if no table
	// Check for generation errors (API errors, quota, etc.)
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		msg := errorStyle.Render("✗ AI Generation Error: " + m.err.Error())
		return paddingStyle.Render(msg)
	}

	// Check for execution errors (SQL errors)
	if m.currentError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		msg := errorStyle.Render("✗ Query Execution Error: " + m.currentError.Error())
		return paddingStyle.Render(msg)
	}

	if m.currentResult != nil {
		successStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true)
		msg := successStyle.Render("✓ Query executed successfully") + "\n" +
			subtleStyle.Render(fmt.Sprintf("(%d rows affected)", len(m.currentResult.Rows)))
		return paddingStyle.Render(msg)
	}

	// Empty state - show welcome screen with logo
	return m.renderWelcomeScreen(height)
}
