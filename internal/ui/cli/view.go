package cli

// View renders the UI
func (m Model) View() string {
	// History view takes over entire screen
	if m.state == stateHistory {
		return m.list.View() + "\n" + subtleStyle.Render("↑↓ navigate • Enter select • Ctrl+d clear history • Esc back")
	}

	// Info view takes over entire screen
	if m.state == stateInfo {
		return m.renderInfoView()
	}

	// Loading screen (connecting or loading schema)
	if m.state == stateConnecting || m.state == stateLoadingSchema {
		return m.renderLoadingScreen()
	}

	// Main screen with results and command bar
	return m.renderMainScreen()
}
