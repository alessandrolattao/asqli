package cli

import (
	"fmt"
	"strings"

	"github.com/alessandrolattao/asqli/internal/features/execution"
	"github.com/atotto/clipboard"
)

// Table represents a navigable table component
type Table struct {
	result *execution.Result

	// Navigation
	selectedRow int
	selectedCol int
	offsetRow   int // Vertical scroll offset
	offsetCol   int // Horizontal scroll offset

	// Dimensions
	width  int
	height int

	// Column widths (calculated)
	colWidths []int
}

// NewTable creates a new table component
func NewTable(result *execution.Result, width, height int) *Table {
	if result == nil {
		return &Table{width: width, height: height}
	}

	// Calculate column widths
	colWidths := make([]int, len(result.Columns))
	for i, col := range result.Columns {
		colWidths[i] = len(col)
		for _, row := range result.Rows {
			if val, ok := row[col]; ok {
				strVal := fmt.Sprintf("%v", val)
				if len(strVal) > colWidths[i] {
					colWidths[i] = len(strVal)
				}
			}
		}
		// Limit max width
		if colWidths[i] > MaxColumnWidth {
			colWidths[i] = MaxColumnWidth
		}
	}

	return &Table{
		result:      result,
		selectedRow: 0,
		selectedCol: 0,
		offsetRow:   0,
		offsetCol:   0,
		width:       width,
		height:      height,
		colWidths:   colWidths,
	}
}

// SetSize updates the table dimensions
func (t *Table) SetSize(width, height int) {
	t.width = width
	t.height = height
	// Ensure selected column and row are still visible after resize
	t.ensureColumnVisible()
	t.ensureRowVisible()
}

// MoveUp moves the selection up
func (t *Table) MoveUp() {
	if t.result == nil || len(t.result.Rows) == 0 {
		return
	}

	if t.selectedRow > 0 {
		t.selectedRow--

		// Adjust scroll if needed - selected row should be first visible when going up
		if t.selectedRow < t.offsetRow {
			t.offsetRow = t.selectedRow
		}
	}
}

// MoveDown moves the selection down
func (t *Table) MoveDown() {
	if t.result == nil || len(t.result.Rows) == 0 {
		return
	}

	if t.selectedRow < len(t.result.Rows)-1 {
		t.selectedRow++

		// Adjust scroll if needed - selected row should be last visible when going down
		visibleRows := t.height - 5 // Account for header, separators, borders, indicators
		lastVisibleRow := t.offsetRow + visibleRows - 1

		if t.selectedRow > lastVisibleRow {
			t.offsetRow++
		}
	}
}

// MoveLeft moves the selection left
func (t *Table) MoveLeft() {
	if t.result == nil || len(t.result.Columns) == 0 {
		return
	}

	if t.selectedCol > 0 {
		t.selectedCol--
		t.ensureColumnVisible()
	}
}

// MoveRight moves the selection right
func (t *Table) MoveRight() {
	if t.result == nil || len(t.result.Columns) == 0 {
		return
	}

	if t.selectedCol < len(t.result.Columns)-1 {
		t.selectedCol++
		t.ensureColumnVisible()
	}
}

// ensureColumnVisible adjusts offsetCol to make sure selectedCol is visible
func (t *Table) ensureColumnVisible() {
	if t.result == nil || len(t.result.Columns) == 0 {
		return
	}

	// If selected column is before offset, move offset left
	if t.selectedCol < t.offsetCol {
		t.offsetCol = t.selectedCol
		return
	}

	// Calculate if selected column is visible with current offset
	totalWidth := 0
	for i := t.offsetCol; i <= t.selectedCol; i++ {
		colWidth := t.colWidths[i] + 3 // padding + separator
		totalWidth += colWidth
	}

	// If selected column doesn't fit, adjust offset to show it
	if totalWidth > t.width {
		// Move offset right until selected column fits
		for t.offsetCol < t.selectedCol {
			// Try removing the leftmost column
			firstColWidth := t.colWidths[t.offsetCol] + 3
			totalWidth -= firstColWidth
			t.offsetCol++

			if totalWidth <= t.width {
				break
			}
		}
	}
}

// ensureRowVisible adjusts offsetRow to make sure selectedRow is visible
func (t *Table) ensureRowVisible() {
	if t.result == nil || len(t.result.Rows) == 0 {
		return
	}

	visibleRows := t.height - 5 // Account for header, separators, borders, indicators

	// If selected row is before offset, move offset up
	if t.selectedRow < t.offsetRow {
		t.offsetRow = t.selectedRow
		return
	}

	// If selected row is after last visible, move offset down
	lastVisibleRow := t.offsetRow + visibleRows - 1
	if t.selectedRow > lastVisibleRow {
		t.offsetRow = t.selectedRow - visibleRows + 1
		if t.offsetRow < 0 {
			t.offsetRow = 0
		}
	}
}

// View renders the table
func (t *Table) View() string {
	if t.result == nil || len(t.result.Rows) == 0 {
		return tableEmptyStyle.Render("No results to display")
	}

	var output strings.Builder

	visibleCols := t.getVisibleColumns()

	// Render top border
	output.WriteString(t.renderTopBorder(visibleCols))
	output.WriteString("\n")

	// Render header
	output.WriteString(t.renderHeader(visibleCols))
	output.WriteString("\n")

	// Render separator
	output.WriteString(t.renderSeparator(visibleCols))
	output.WriteString("\n")

	// Render rows
	visibleRows := t.height - 5 // Account for header, separators, borders, indicators
	endRow := t.offsetRow + visibleRows
	if endRow > len(t.result.Rows) {
		endRow = len(t.result.Rows)
	}

	for i := t.offsetRow; i < endRow; i++ {
		output.WriteString(t.renderRow(i, visibleCols))
		output.WriteString("\n")
	}

	// Render bottom border
	output.WriteString(t.renderBottomBorder(visibleCols))

	// Add scroll indicators if needed
	var indicators []string

	// Row scroll indicator
	if t.offsetRow > 0 || endRow < len(t.result.Rows) {
		rowInfo := fmt.Sprintf("Rows %d-%d of %d", t.offsetRow+1, endRow, len(t.result.Rows))
		indicators = append(indicators, rowInfo)
	}

	// Column scroll indicator
	hasMoreLeft := t.offsetCol > 0
	hasMoreRight := len(visibleCols) > 0 && visibleCols[len(visibleCols)-1] < len(t.result.Columns)-1

	if hasMoreLeft || hasMoreRight {
		lastVisible := 0
		if len(visibleCols) > 0 {
			lastVisible = visibleCols[len(visibleCols)-1]
		}
		colInfo := fmt.Sprintf("Cols %d-%d of %d", t.offsetCol+1, lastVisible+1, len(t.result.Columns))
		indicators = append(indicators, colInfo)
	}

	if len(indicators) > 0 {
		output.WriteString("\n")

		// Calculate table width
		tableWidth := 2 // Left and right borders
		for _, colIdx := range visibleCols {
			tableWidth += t.colWidths[colIdx] + 2 // Column width + padding
		}
		tableWidth += len(visibleCols) - 1 // Separators between columns

		// Build center text
		centerText := strings.Join(indicators, " • ")

		// Calculate padding
		leftArrowChar := ""
		rightArrowChar := ""
		if hasMoreLeft {
			leftArrowChar = "←"
		}
		if hasMoreRight {
			rightArrowChar = "→"
		}

		// Available space for content (excluding arrows)
		availableWidth := tableWidth
		if hasMoreLeft {
			availableWidth--
		}
		if hasMoreRight {
			availableWidth--
		}

		// Calculate padding for centering
		textLen := len(centerText)
		if textLen < availableWidth {
			leftPad := (availableWidth - textLen) / 2
			rightPad := availableWidth - textLen - leftPad
			centerText = strings.Repeat(" ", leftPad) + centerText + strings.Repeat(" ", rightPad)
		}

		// Build indicator line
		indicatorLine := leftArrowChar + centerText + rightArrowChar
		output.WriteString(tableSubtleStyle.Render(indicatorLine))
	}

	return output.String()
}

// getVisibleColumns calculates which columns can fit in the current width
func (t *Table) getVisibleColumns() []int {
	if len(t.result.Columns) == 0 {
		return []int{}
	}

	visible := []int{}
	totalWidth := 0

	// Start from horizontal offset
	for i := t.offsetCol; i < len(t.result.Columns); i++ {
		colWidth := t.colWidths[i] + 3 // Add padding and separator
		if totalWidth+colWidth > t.width {
			break
		}
		visible = append(visible, i)
		totalWidth += colWidth
	}

	// Ensure at least one column is visible
	if len(visible) == 0 && len(t.result.Columns) > 0 {
		// Show the offset column even if it doesn't fit perfectly
		visible = append(visible, t.offsetCol)
	}

	return visible
}

// renderTopBorder renders the top border of the table
func (t *Table) renderTopBorder(visibleCols []int) string {
	var border strings.Builder

	hasMoreLeft := t.offsetCol > 0
	hasMoreRight := len(visibleCols) > 0 && visibleCols[len(visibleCols)-1] < len(t.result.Columns)-1

	// Left corner - use ┬ if there are more columns to the left
	if hasMoreLeft {
		border.WriteString(tableSeparatorStyle.Render("┬"))
	} else {
		border.WriteString(tableSeparatorStyle.Render("╭"))
	}

	for idx, colIdx := range visibleCols {
		border.WriteString(tableSeparatorStyle.Render(strings.Repeat("─", t.colWidths[colIdx]+2)))

		if idx != len(visibleCols)-1 {
			border.WriteString(tableSeparatorStyle.Render("┬"))
		}
	}

	// Right corner - use ┬ if there are more columns to the right
	if hasMoreRight {
		border.WriteString(tableSeparatorStyle.Render("┬"))
	} else {
		border.WriteString(tableSeparatorStyle.Render("╮"))
	}

	return border.String()
}

// renderBottomBorder renders the bottom border of the table
func (t *Table) renderBottomBorder(visibleCols []int) string {
	var border strings.Builder

	hasMoreLeft := t.offsetCol > 0
	hasMoreRight := len(visibleCols) > 0 && visibleCols[len(visibleCols)-1] < len(t.result.Columns)-1

	// Left corner - use ┴ if there are more columns to the left
	if hasMoreLeft {
		border.WriteString(tableSeparatorStyle.Render("┴"))
	} else {
		border.WriteString(tableSeparatorStyle.Render("╰"))
	}

	for idx, colIdx := range visibleCols {
		border.WriteString(tableSeparatorStyle.Render(strings.Repeat("─", t.colWidths[colIdx]+2)))

		if idx != len(visibleCols)-1 {
			border.WriteString(tableSeparatorStyle.Render("┴"))
		}
	}

	// Right corner - use ┴ if there are more columns to the right
	if hasMoreRight {
		border.WriteString(tableSeparatorStyle.Render("┴"))
	} else {
		border.WriteString(tableSeparatorStyle.Render("╯"))
	}

	return border.String()
}

// renderHeader renders the table header
func (t *Table) renderHeader(visibleCols []int) string {
	var header strings.Builder

	header.WriteString(tableSeparatorStyle.Render("│"))

	for _, colIdx := range visibleCols {
		col := t.result.Columns[colIdx]
		cellContent := fmt.Sprintf("%-*s", t.colWidths[colIdx], col)

		if colIdx == t.selectedCol {
			header.WriteString(tableHeaderSelectedStyle.Render(cellContent))
		} else {
			header.WriteString(tableHeaderStyle.Render(cellContent))
		}

		header.WriteString(tableSeparatorStyle.Render("│"))
	}

	return header.String()
}

// renderSeparator renders the separator line
func (t *Table) renderSeparator(visibleCols []int) string {
	var separator strings.Builder

	separator.WriteString(tableSeparatorStyle.Render("├"))

	for idx, colIdx := range visibleCols {
		separator.WriteString(tableSeparatorStyle.Render(strings.Repeat("─", t.colWidths[colIdx]+2)))

		if idx != len(visibleCols)-1 {
			separator.WriteString(tableSeparatorStyle.Render("┼"))
		}
	}

	separator.WriteString(tableSeparatorStyle.Render("┤"))

	return separator.String()
}

// CopyToClipboard copies the entire table to clipboard as TSV
func (t *Table) CopyToClipboard() error {
	if t.result == nil || len(t.result.Rows) == 0 {
		return fmt.Errorf("no data to copy")
	}

	var output strings.Builder

	// Header row
	for i, col := range t.result.Columns {
		output.WriteString(col)
		if i < len(t.result.Columns)-1 {
			output.WriteString("\t")
		}
	}
	output.WriteString("\n")

	// Data rows
	for _, row := range t.result.Rows {
		for i, col := range t.result.Columns {
			val := row[col]
			output.WriteString(fmt.Sprintf("%v", val))
			if i < len(t.result.Columns)-1 {
				output.WriteString("\t")
			}
		}
		output.WriteString("\n")
	}

	return clipboard.WriteAll(output.String())
}

// GetSelectedColumn returns the name of the currently selected column
func (t *Table) GetSelectedColumn() string {
	if t.result == nil || len(t.result.Columns) == 0 {
		return ""
	}
	if t.selectedCol < 0 || t.selectedCol >= len(t.result.Columns) {
		return ""
	}
	return t.result.Columns[t.selectedCol]
}

// GetSelectedValue returns the value of the currently selected cell
func (t *Table) GetSelectedValue() any {
	if t.result == nil || len(t.result.Rows) == 0 {
		return nil
	}
	if t.selectedRow < 0 || t.selectedRow >= len(t.result.Rows) {
		return nil
	}
	if t.selectedCol < 0 || t.selectedCol >= len(t.result.Columns) {
		return nil
	}

	col := t.result.Columns[t.selectedCol]
	return t.result.Rows[t.selectedRow][col]
}

// renderRow renders a single data row
func (t *Table) renderRow(rowIdx int, visibleCols []int) string {
	var row strings.Builder
	resultRow := t.result.Rows[rowIdx]

	row.WriteString(tableDividerStyle.Render("│"))

	for _, colIdx := range visibleCols {
		col := t.result.Columns[colIdx]
		val := resultRow[col]
		strVal := fmt.Sprintf("%v", val)

		// Truncate if too long
		if len(strVal) > t.colWidths[colIdx] {
			strVal = strVal[:t.colWidths[colIdx]-3] + "..."
		}

		cellContent := fmt.Sprintf("%-*s", t.colWidths[colIdx], strVal)

		// Highlight selected cell
		if rowIdx == t.selectedRow && colIdx == t.selectedCol {
			row.WriteString(tableCellSelectedStyle.Render(cellContent))
		} else if rowIdx == t.selectedRow {
			row.WriteString(tableRowSelectedStyle.Render(cellContent))
		} else {
			row.WriteString(tableCellStyle.Render(cellContent))
		}

		row.WriteString(tableDividerStyle.Render("│"))
	}

	return row.String()
}
