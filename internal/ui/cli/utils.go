package cli

// wrapText wraps text to fit within the given width, attempting to break at natural boundaries.
// It prefers breaking at spaces, commas, or parentheses to maintain SQL readability.
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	for len(text) > width {
		// Find a good breaking point (space, comma, etc.)
		breakPoint := width
		for i := width; i > 0; i-- {
			if text[i] == ' ' || text[i] == ',' || text[i] == '(' || text[i] == ')' {
				breakPoint = i + 1
				break
			}
		}

		lines = append(lines, text[:breakPoint])
		text = text[breakPoint:]
	}

	if len(text) > 0 {
		lines = append(lines, text)
	}

	return lines
}
