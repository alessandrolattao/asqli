// Package config provides shared configuration constants and utilities.
package config

import "time"

// TimeoutConfig holds timeout values for various operations to prevent indefinite blocking.
// These values are tuned for typical usage scenarios while avoiding false timeouts.
type TimeoutConfig struct {
	// DatabaseConnection is the maximum time to wait for database connection and ping
	DatabaseConnection time.Duration

	// DatabaseQuery is the maximum time to wait for a database query to execute
	DatabaseQuery time.Duration

	// SchemaFetch is the maximum time to wait for database schema extraction
	SchemaFetch time.Duration

	// AIGeneration is the maximum time to wait for AI providers to generate SQL
	AIGeneration time.Duration
}

// DefaultTimeouts returns the default timeout configuration
func DefaultTimeouts() TimeoutConfig {
	return TimeoutConfig{
		DatabaseConnection: 10 * time.Second,
		DatabaseQuery:      30 * time.Second,
		SchemaFetch:        30 * time.Second,
		AIGeneration:       60 * time.Second,
	}
}
