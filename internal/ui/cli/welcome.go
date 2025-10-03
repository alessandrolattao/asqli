package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderWelcomeScreen renders the welcome screen with logo and helpful information
func (m Model) renderWelcomeScreen(height int) string {
	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB6C1")).
		Bold(true)
	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E0E0E0"))
	exampleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#98C379"))
	subtleTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#808080"))

	// Build content
	var content strings.Builder

	// Logo with tagline (centered)
	content.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(GetLogo()))
	content.WriteString("\n\n")

	// Connection info (centered)
	dbType := string(m.dbConfig.DriverType)
	providerName := string(m.aiConfig.Type)
	if m.aiProvider != nil {
		providerName = m.aiProvider.Name()
	}

	connectionInfo := subtleTextStyle.Render(fmt.Sprintf("Connected to %s • AI: %s", dbType, providerName))
	if m.aiConfig.Model != "" {
		connectionInfo += subtleTextStyle.Render(fmt.Sprintf(" (%s)", m.aiConfig.Model))
	}
	content.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(connectionInfo))
	content.WriteString("\n\n")

	// Main description (centered)
	description := textStyle.Render("Ask questions about your database in natural language") + "\n" +
		subtleTextStyle.Render("or use ") + exampleStyle.Render("#") + subtleTextStyle.Render(" prefix for raw SQL queries")
	content.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(description))
	content.WriteString("\n\n")

	// Examples section (left-aligned within a centered box)
	var examples strings.Builder
	examples.WriteString(titleStyle.Render("Examples:"))
	examples.WriteString("\n")
	examples.WriteString(textStyle.Render("  • "))
	examples.WriteString(exampleStyle.Render("show me all users from Italy"))
	examples.WriteString("\n")
	examples.WriteString(textStyle.Render("  • "))
	examples.WriteString(exampleStyle.Render("count active subscriptions by plan"))
	examples.WriteString("\n")
	examples.WriteString(textStyle.Render("  • "))
	examples.WriteString(exampleStyle.Render("list top 10 customers by revenue"))
	examples.WriteString("\n")
	examples.WriteString(textStyle.Render("  • "))
	examples.WriteString(exampleStyle.Render("# SELECT * FROM users WHERE created_at > NOW() - INTERVAL '7 days'"))

	// Center the examples box, but keep text left-aligned inside
	examplesBox := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Render(examples.String())
	content.WriteString(lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(examplesBox))
	content.WriteString("\n")

	return lipgloss.Place(m.width, height, lipgloss.Center, lipgloss.Center, content.String())
}
