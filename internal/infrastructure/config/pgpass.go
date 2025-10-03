package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// PgPassEntry represents a single line in .pgpass file
type PgPassEntry struct {
	Hostname string
	Port     string
	Database string
	Username string
	Password string
}

// GetPgPassFilePath returns the path to the .pgpass file
func GetPgPassFilePath() (string, error) {
	// Check PGPASSFILE environment variable first
	if pgpassFile := os.Getenv("PGPASSFILE"); pgpassFile != "" {
		return pgpassFile, nil
	}

	// Default location depends on OS
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		return filepath.Join(appData, "postgresql", "pgpass.conf"), nil
	}

	// Unix/Linux/macOS
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("HOME environment variable not set")
	}
	return filepath.Join(home, ".pgpass"), nil
}

// ParsePgPass reads and parses the .pgpass file
func ParsePgPass() ([]PgPassEntry, error) {
	pgpassPath, err := GetPgPassFilePath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	info, err := os.Stat(pgpassPath)
	if os.IsNotExist(err) {
		return nil, nil // File doesn't exist, not an error
	}
	if err != nil {
		return nil, fmt.Errorf("error accessing .pgpass file: %w", err)
	}

	// Check permissions on Unix-like systems
	if runtime.GOOS != "windows" {
		mode := info.Mode()
		if mode.Perm()&0077 != 0 {
			return nil, fmt.Errorf(".pgpass file has incorrect permissions (must be 0600)")
		}
	}

	file, err := os.Open(pgpassPath)
	if err != nil {
		return nil, fmt.Errorf("error opening .pgpass file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: error closing .pgpass file: %v\n", closeErr)
		}
	}()

	var entries []PgPassEntry
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entry, err := parsePgPassLine(line)
		if err != nil {
			// Log warning but continue parsing
			fmt.Fprintf(os.Stderr, "Warning: invalid .pgpass entry at line %d: %v\n", lineNum, err)
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .pgpass file: %w", err)
	}

	return entries, nil
}

// parsePgPassLine parses a single line from .pgpass file
func parsePgPassLine(line string) (PgPassEntry, error) {
	// Split by unescaped colons
	parts := splitPgPassLine(line)

	if len(parts) != 5 {
		return PgPassEntry{}, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	return PgPassEntry{
		Hostname: parts[0],
		Port:     parts[1],
		Database: parts[2],
		Username: parts[3],
		Password: parts[4],
	}, nil
}

// splitPgPassLine splits a .pgpass line by unescaped colons
func splitPgPassLine(line string) []string {
	var parts []string
	var current strings.Builder
	escaped := false

	for i, r := range line {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			// Check if this is escaping a colon or backslash
			if i+1 < len(line) && (line[i+1] == ':' || line[i+1] == '\\') {
				escaped = true
				continue
			}
			// Otherwise, keep the backslash
			current.WriteRune(r)
			continue
		}

		if r == ':' {
			parts = append(parts, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	// Add the last part
	parts = append(parts, current.String())

	return parts
}

// FindPassword searches for a matching password in .pgpass file
func FindPassword(hostname, port, database, username string) (string, error) {
	entries, err := ParsePgPass()
	if err != nil {
		return "", err
	}

	if entries == nil {
		return "", nil // No .pgpass file
	}

	// Convert port to string if needed
	portStr := port

	// Match entries (first match wins)
	for _, entry := range entries {
		if matchField(entry.Hostname, hostname) &&
			matchField(entry.Port, portStr) &&
			matchField(entry.Database, database) &&
			matchField(entry.Username, username) {
			return entry.Password, nil
		}
	}

	return "", nil // No match found
}

// matchField checks if a field matches, supporting wildcards
func matchField(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	return pattern == value
}

// GetPasswordForConfig attempts to find password in .pgpass for given connection parameters
func GetPasswordForConfig(host string, port int, database, username string) string {
	portStr := strconv.Itoa(port)
	password, err := FindPassword(host, portStr, database, username)
	if err != nil {
		// Log error but don't fail
		fmt.Fprintf(os.Stderr, "Warning: error reading .pgpass: %v\n", err)
		return ""
	}
	return password
}
