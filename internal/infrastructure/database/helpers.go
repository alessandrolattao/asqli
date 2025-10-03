package database

import (
	"database/sql"
	"fmt"
	"os"
)

// closeRows is a helper function to close sql.Rows and log any errors.
// This is intended for use in defer statements where we want to ensure cleanup
// but cannot return the error. Errors are logged to stderr.
//
// Usage:
//
//	defer closeRows(rows)
func closeRows(rows *sql.Rows) {
	if rows == nil {
		return
	}
	if err := rows.Close(); err != nil {
		// Log to stderr since we're in a defer and can't return the error
		fmt.Fprintf(os.Stderr, "Warning: Failed to close database rows: %v\n", err)
	}
}
