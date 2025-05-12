package config

import (
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
)

// NewCustomConfig creates a configuration using a custom connection string
func NewCustomConfig(driverType database.DriverType, connectionString string) *Config {
	return &Config{
		DB: database.Config{
			DriverType:       driverType,
			ConnectionString: connectionString,
		},
	}
}