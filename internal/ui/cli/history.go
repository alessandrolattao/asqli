package cli

import (
	"os"
	"path/filepath"
	"strings"
)

// loadHistory loads command history from disk
func loadHistory() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	historyFile := filepath.Join(homeDir, ".sqlai_history")
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return []string{}
	}

	lines := strings.Split(string(data), "\n")
	// Filter empty lines
	history := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			history = append(history, line)
		}
	}

	return history
}

// saveHistory saves command history to disk
func saveHistory(history []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	historyFile := filepath.Join(homeDir, ".sqlai_history")
	data := strings.Join(history, "\n")
	return os.WriteFile(historyFile, []byte(data), 0600)
}
