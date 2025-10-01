package execution

// Result represents the result of a query execution
type Result struct {
	// Rows contains the query result data
	Rows []map[string]any

	// Columns contains the column names
	Columns []string
}
