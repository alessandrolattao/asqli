package config

import (
	"github.com/alessandrolattao/sqlai/internal/pkg/database"
)

// Config holds all configuration settings for the application
type Config struct {
	DB database.Config
}