package logger

import "go.uber.org/zap"

// Logger interface defines a minimal subset of zap.Logger methods
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
}

// zapLogger implements the Logger interface using zap
type zapLogger struct {
	logger *zap.Logger
}

// Debug logs a message at debug level
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs a message at info level
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a message at warn level
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs a message at error level
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with the given fields
func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

// NewLogger creates a new logger from a zap.Logger
func NewLogger(logger *zap.Logger) Logger {
	return &zapLogger{
		logger: logger,
	}
}

// Default returns a default logger
func Default() Logger {
	return NewLogger(GetLogger())
}
