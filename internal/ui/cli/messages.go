package cli

import (
	"github.com/alessandrolattao/sqlai/internal/features/execution"
	"github.com/alessandrolattao/sqlai/internal/features/query"
	"github.com/alessandrolattao/sqlai/internal/features/schema"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/ai"
	"github.com/alessandrolattao/sqlai/internal/infrastructure/database"
)

// connectionMsg is sent when database connection completes
type connectionMsg struct {
	dbConn           *database.Connection
	aiProvider       ai.Provider
	schemaService    *schema.Service
	queryService     *query.Service
	executionService *execution.Service
	err              error
}

// schemaMsg is sent when schema fetch completes
type schemaMsg struct {
	schema string
	err    error
}

// sqlGeneratedMsg is sent when SQL generation completes
type sqlGeneratedMsg struct {
	sql *query.SQL
	err error
}

// queryExecutedMsg is sent when query execution completes
type queryExecutedMsg struct {
	result *execution.Result
	err    error
}
