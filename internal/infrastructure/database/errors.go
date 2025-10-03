// Package database defines errors related to database connections and operations.
package database

import "errors"

// Sentinel errors returned by database operations and the adapter factory.
var (
	// ErrUnsupportedDriver is returned when an unsupported database driver is requested
	ErrUnsupportedDriver = errors.New("unsupported database driver")

	// ErrInvalidConfig is returned when database configuration is invalid
	ErrInvalidConfig = errors.New("invalid database configuration")

	// ErrConnectionFailed is returned when database connection fails
	ErrConnectionFailed = errors.New("database connection failed")
)
