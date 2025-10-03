package cli

import (
	"github.com/charmbracelet/lipgloss"
)

// renderMainScreen renders the main screen with results area and command bar
func (m Model) renderMainScreen() string {
	// Calculate areas
	resultsHeight := m.height - CommandBarHeight

	// Render results area - fill the available space
	resultsArea := m.renderResultsArea(resultsHeight)

	// Ensure results area fills the height
	resultsAreaPadded := lipgloss.NewStyle().
		Height(resultsHeight).
		Render(resultsArea)

	// Render command bar
	commandBar := NewCommandBar(m.width, m.state, m.spinner, m.textInput, m.statusMessage, m.generatedSQL)
	commandBarView := commandBar.View()

	// Combine vertically - results area fills space, command bar at bottom
	return lipgloss.JoinVertical(
		lipgloss.Left,
		resultsAreaPadded,
		commandBarView,
	)
}
