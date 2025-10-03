package cli

import (
	"github.com/charmbracelet/lipgloss"
)

// renderLoadingScreen renders the loading screen with centered logo and spinner
func (m Model) renderLoadingScreen() string {
	var statusMsg string
	if m.state == stateConnecting {
		statusMsg = m.spinner.View() + " Connecting to database..."
	} else {
		statusMsg = m.spinner.View() + " Loading schema..."
	}

	// Center everything vertically and horizontally
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		GetLogo(),
		"",
		spinnerStyle.Render(statusMsg),
	)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
