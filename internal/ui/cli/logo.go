package cli

import (
	"github.com/charmbracelet/lipgloss"
)

// GetLogo returns the ASQLI ASCII art logo with tagline
func GetLogo() string {
	logo := `
 ░▒▓██████▓▒░ ░▒▓███████▓▒░░▒▓██████▓▒░░▒▓█▓▒░      ░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░
░▒▓████████▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓████████▓▒░▒▓█▓▒░
                             ░▒▓█▓▒░`

	taglineSymbol := "                              ░▒▓██▓▒░"
	taglineText := " AI-POWERED SQL CLIENT"

	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB6C1"))
	taglineStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D3D3D3"))

	return logoStyle.Render(logo) + "\n" + logoStyle.Render(taglineSymbol) + taglineStyle.Render(taglineText)
}
