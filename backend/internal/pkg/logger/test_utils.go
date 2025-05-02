package logger

import "go.uber.org/zap"

// The following functions are for testing purposes only

// ResetForTests resets the global logger for testing purposes
func ResetForTests() {
	globalLogger = nil
}

// SetGlobalLoggerForTests sets the global logger for testing
func SetGlobalLoggerForTests(logger *zap.Logger) {
	globalLogger = logger
}

// GetGlobalLoggerForTests returns the current globalLogger for testing
func GetGlobalLoggerForTests() *zap.Logger {
	return globalLogger
}

// NewLoggerForTests creates a new logger with the given config for testing
func NewLoggerForTests(cfg Config) *zap.Logger {
	return newLogger(cfg)
}
