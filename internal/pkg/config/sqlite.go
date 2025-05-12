package config

import (
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
)

// NewSQLiteConfig creates a new configuration for SQLite
func NewSQLiteConfig(filePath string) *Config {
	return &Config{
		DB: database.Config{
			DriverType: database.SQLite,
			FilePath:   filePath,
		},
	}
}