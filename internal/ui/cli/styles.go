package cli

import "github.com/charmbracelet/lipgloss"

var (
	// Logo style with pink color
	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFB6C1"))

	// Subtle text style (help, hints)
	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFB6C1")).
				Bold(true).
				Padding(0, 1)

	tableHeaderSelectedStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FFB6C1")).
					Background(lipgloss.Color("#2D2D2D")).
					Bold(true).
					Padding(0, 1)

	tableSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666"))

	tableDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#666666"))

	tableCellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0E0E0")).
			Padding(0, 1)

	tableRowSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#1A1A1A")).
				Padding(0, 1)

	tableCellSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#FFB6C1")).
				Bold(true).
				Padding(0, 1)

	tableEmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true)

	tableSubtleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	// Spinner style
	spinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB6C1"))

	// Danger/warning style for dangerous operations
	dangerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E87580")).
			Bold(true)

	// Error style for error messages
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E06C75")).
			Bold(true)

	// Success style for success messages
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98C379")).
			Bold(true)
)
