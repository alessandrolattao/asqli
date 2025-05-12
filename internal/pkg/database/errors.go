package database

import "errors"

// Common database errors
var (
	// ErrUnsupportedDriver is returned when an unsupported database driver is specified
	ErrUnsupportedDriver = errors.New("unsupported database driver")
	
	// ErrInvalidConfig is returned when the database configuration is invalid
	ErrInvalidConfig = errors.New("invalid database configuration")
	
	// ErrConnectionFailed is returned when the database connection attempt fails
	ErrConnectionFailed = errors.New("failed to connect to database")
	
	// ErrQueryFailed is returned when a database query fails
	ErrQueryFailed = errors.New("database query failed")
	
	// ErrTableNotFound is returned when a specified table does not exist
	ErrTableNotFound = errors.New("table not found")
)