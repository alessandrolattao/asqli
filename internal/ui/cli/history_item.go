package cli

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// historyItem represents an item in the history list
type historyItem struct {
	query string
}

// FilterValue implements list.Item
func (i historyItem) FilterValue() string {
	return i.query
}

// simpleDelegate is a minimal list item delegate
type simpleDelegate struct{}

// Height implements list.ItemDelegate
func (d simpleDelegate) Height() int {
	return 1
}

// Spacing implements list.ItemDelegate
func (d simpleDelegate) Spacing() int {
	return 0
}

// Update implements list.ItemDelegate
func (d simpleDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

// Render implements list.ItemDelegate
func (d simpleDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	histItem, ok := item.(historyItem)
	if !ok {
		return
	}

	str := histItem.query
	if len(str) > m.Width()-4 {
		str = str[:m.Width()-7] + "..."
	}

	// Selected item style
	if index == m.Index() {
		selectedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB6C1")).
			Bold(true)
		str = "→ " + selectedStyle.Render(str)
	} else {
		str = "  " + str
	}

	_, _ = w.Write([]byte(str))
}
