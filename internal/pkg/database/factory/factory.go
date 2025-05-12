package factory

import (
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
	"github.com/alessandrolattao/sqlai/internal/pkg/database/drivers"
)

// NewPostgreSQLAdapter creates a new PostgreSQL adapter
func NewPostgreSQLAdapter() database.Adapter {
	return &drivers.PostgreSQLAdapter{}
}

// NewMySQLAdapter creates a new MySQL adapter
func NewMySQLAdapter() database.Adapter {
	return &drivers.MySQLAdapter{}
}

// NewSQLiteAdapter creates a new SQLite adapter
func NewSQLiteAdapter() database.Adapter {
	return &drivers.SQLiteAdapter{}
}