package config

import (
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
)

// NewMySQLConfig creates a new configuration for MySQL
func NewMySQLConfig(dbHost, dbUser, dbPassword, dbName string, dbPort int, parseTime bool) *Config {
	return &Config{
		DB: database.Config{
			DriverType: database.MySQL,
			Host:       dbHost,
			Port:       dbPort,
			User:       dbUser,
			Password:   dbPassword,
			DBName:     dbName,
			ParseTime:  parseTime,
		},
	}
}