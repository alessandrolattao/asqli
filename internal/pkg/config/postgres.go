package config

import (
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
)

// NewConfig creates a new configuration for PostgreSQL
func NewConfig(dbHost, dbUser, dbPassword, dbName, dbSSLMode string, dbPort int) *Config {
	return &Config{
		DB: database.Config{
			DriverType: database.PostgreSQL,
			Host:       dbHost,
			Port:       dbPort,
			User:       dbUser,
			Password:   dbPassword,
			DBName:     dbName,
			SSLMode:    dbSSLMode,
		},
	}
}