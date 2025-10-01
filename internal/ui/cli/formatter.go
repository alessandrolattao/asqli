package cli

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	Bold        = "\033[1m"
)

// formatAsTable formats query results as an ASCII table with colors
func formatAsTable(data []map[string]any, columns []string) string {
	if len(data) == 0 {
		return ColorYellow + "No results found" + ColorReset
	}

	// Determine maximum width for each column
	colWidths := make(map[string]int)
	for _, col := range columns {
		colWidths[col] = len(col)
	}

	for _, row := range data {
		for col, val := range row {
			strVal := fmt.Sprintf("%v", val)
			if len(strVal) > colWidths[col] {
				colWidths[col] = len(strVal)
			}
		}
	}

	// Build the table
	var sb strings.Builder

	// Border color
	borderColor := ColorBlue

	// Header
	sb.WriteString(borderColor + "+")
	for _, col := range columns {
		sb.WriteString(strings.Repeat("-", colWidths[col]+2))
		sb.WriteString("+")
	}
	sb.WriteString(ColorReset + "\n")

	// Column headers with bold and cyan
	sb.WriteString(borderColor + "|" + ColorReset)
	for _, col := range columns {
		fmt.Fprintf(&sb, " %s%s%-*s%s %s|", Bold, ColorCyan, colWidths[col], col, ColorReset, borderColor)
	}
	sb.WriteString(ColorReset + "\n")

	// Separator
	sb.WriteString(borderColor + "+")
	for _, col := range columns {
		sb.WriteString(strings.Repeat("-", colWidths[col]+2))
		sb.WriteString("+")
	}
	sb.WriteString(ColorReset + "\n")

	// Rows with alternating colors
	for i, row := range data {
		rowColor := ColorGray // Light gray for odd rows
		if i%2 == 0 {
			rowColor = "" // Default color for even rows
		}

		sb.WriteString(borderColor + "|" + ColorReset)
		for _, col := range columns {
			val := row[col]
			strVal := fmt.Sprintf("%v", val)
			fmt.Fprintf(&sb, " %s%-*s%s %s|", rowColor, colWidths[col], strVal, ColorReset, borderColor)
		}
		sb.WriteString(ColorReset + "\n")
	}

	// Footer
	sb.WriteString(borderColor + "+")
	for _, col := range columns {
		sb.WriteString(strings.Repeat("-", colWidths[col]+2))
		sb.WriteString("+")
	}
	sb.WriteString(ColorReset + "\n")

	return sb.String()
}

// formatSuccessMessage returns an appropriate success message for different query types
func formatSuccessMessage(queryType string) string {
	switch queryType {
	case "INSERT":
		return fmt.Sprintf("%s%sSuccess: Data inserted successfully.%s\n", Bold, ColorGreen, ColorReset)
	case "UPDATE":
		return fmt.Sprintf("%s%sSuccess: Data updated successfully.%s\n", Bold, ColorGreen, ColorReset)
	case "DELETE":
		return fmt.Sprintf("%s%sSuccess: Data deleted successfully.%s\n", Bold, ColorGreen, ColorReset)
	default:
		return fmt.Sprintf("%s%sSuccess: Query executed successfully.%s\n", Bold, ColorGreen, ColorReset)
	}
}
