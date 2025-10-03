// Package cli contains the Bubble Tea UI state machine.
package cli

// state represents the current state of the CLI application's finite state machine.
// The application transitions between these states during normal operation.
type state int

const (
	// stateConnecting indicates the app is connecting to the database
	stateConnecting state = iota

	// stateLoadingSchema indicates the app is loading the database schema
	stateLoadingSchema

	// stateReady indicates the app is ready for user input
	stateReady

	// stateThinking indicates the AI is generating SQL from user prompt
	stateThinking

	// stateExecuting indicates a SQL query is being executed
	stateExecuting

	// stateConfirming indicates the app is waiting for user confirmation of a dangerous query
	stateConfirming

	// stateHistory indicates the app is displaying query history
	stateHistory

	// stateInfo indicates the app is displaying information about the last query
	stateInfo
)
