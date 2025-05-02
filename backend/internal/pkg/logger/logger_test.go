package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"runtime"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// captureOutput redirects standard output and returns a function to restore it
// and get the captured output
func captureOutput(t *testing.T) func() string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w

	return func() string {
		w.Close()
		out, err := io.ReadAll(r)
		require.NoError(t, err)
		os.Stdout = originalStdout
		return string(out)
	}
}

func TestNewLoggerWithConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		checkOutput func(t *testing.T, output string)
	}{
		{
			name: "debug level console encoder",
			config: Config{
				Level:      "debug",
				Encoding:   "console",
				OutputPath: "stdout",
				DevMode:    true,
			},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "Logger initialized")
				// Don't check for DEBUG in the output since the log message itself is at INFO level
			},
		},
		{
			name: "info level json encoder",
			config: Config{
				Level:      "info",
				Encoding:   "json",
				OutputPath: "stdout",
				DevMode:    false,
			},
			checkOutput: func(t *testing.T, output string) {
				var logEntry map[string]interface{}
				err := json.Unmarshal([]byte(strings.TrimSpace(output)), &logEntry)
				require.NoError(t, err)

				// Changed to check for partial message since the full message now includes "with rotation"
				assert.Contains(t, logEntry["message"].(string), "Logger initialized")
				assert.Equal(t, "INFO", logEntry["level"])
			},
		},
		{
			name: "invalid level defaults to info",
			config: Config{
				Level:      "invalid_level",
				Encoding:   "console",
				OutputPath: "stdout",
				DevMode:    false,
			},
			checkOutput: func(t *testing.T, output string) {
				assert.Contains(t, output, "INFO")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getOutput := captureOutput(t)

			// Reset globalLogger to ensure consistent test environment
			globalLogger = nil

			// Create logger with test config using the private function
			logger := newLogger(tt.config)
			require.NotNil(t, logger)

			// Check output
			output := getOutput()
			tt.checkOutput(t, output)
		})
	}
}

func TestGetLogger(t *testing.T) {
	// Reset global logger for test
	origLogger := globalLogger
	defer func() { globalLogger = origLogger }()

	// Test 1: When globalLogger is not initialized
	globalLogger = nil

	// GetLogger should return a new logger but not set globalLogger
	logger1 := GetLogger()
	require.NotNil(t, logger1)

	// globalLogger should still be nil
	assert.Nil(t, globalLogger, "globalLogger should still be nil after GetLogger call")

	// Test 2: When globalLogger is initialized
	mockLogger := zap.NewNop()
	globalLogger = mockLogger

	// GetLogger should return the globalLogger
	logger2 := GetLogger()
	assert.Same(t, mockLogger, logger2, "GetLogger should return the globalLogger when set")
}

func TestInitLogger(t *testing.T) {
	// Reset global logger
	origLogger := globalLogger
	defer func() { globalLogger = origLogger }()

	globalLogger = nil

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "logger_test_*.log")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	config := Config{
		Level:      "info",
		Encoding:   "console",
		OutputPath: tmpFile.Name(),
		DevMode:    false,
	}

	// Initialize logger with the temp file output
	logger := InitLogger(config)
	require.NotNil(t, logger)

	// Force logger to sync to ensure log is written
	err = logger.Sync()
	require.NoError(t, err)

	// Read the log file content
	_, err = tmpFile.Seek(0, 0)
	require.NoError(t, err)
	content, err := io.ReadAll(tmpFile)
	require.NoError(t, err)

	// Verify log message
	assert.Contains(t, string(content), "Logger initialized")

	// Get the original logger for comparison
	oldLogger := globalLogger

	// This should not create a new logger
	sameLogger := InitLogger(Config{Level: "debug"})

	// Verify the singleton behavior
	assert.Same(t, oldLogger, globalLogger)
	assert.Same(t, oldLogger, sameLogger)
}

func TestLoggerWrapper(t *testing.T) {
	// Create a memory syncer for testing log output
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "time",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	zapLogger := zap.New(core)

	// Create our wrapper
	log := NewLogger(zapLogger)

	// Test all logging methods
	log.Debug("debug message") // Should not appear - below threshold
	log.Info("info message")
	log.Warn("warning message")
	log.Error("error message")

	// Check output
	output := buf.String()
	assert.NotContains(t, output, "debug message")
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warning message")
	assert.Contains(t, output, "error message")

	// Test With method
	contextLogger := log.With(String("context", "test"))
	buf.Reset()
	contextLogger.Info("context test")

	contextOutput := buf.String()
	assert.Contains(t, contextOutput, "context")
	assert.Contains(t, contextOutput, "test")
}

func TestFieldHelpers(t *testing.T) {
	// Setup in-memory logger
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "time",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	zapLogger := zap.New(core)
	log := NewLogger(zapLogger)

	// Test field helpers
	log.Info("test message",
		String("string_field", "value"),
		Int("int_field", 42),
		Int64("int64_field", 42),
		Float64("float64_field", 3.14),
		Bool("bool_field", true),
		Any("any_field", map[string]string{"key": "value"}),
	)

	// Parse JSON
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Validate fields
	assert.Equal(t, "value", logEntry["string_field"])
	assert.Equal(t, float64(42), logEntry["int_field"]) // JSON numbers are float64
	assert.Equal(t, float64(42), logEntry["int64_field"])
	assert.Equal(t, 3.14, logEntry["float64_field"])
	assert.Equal(t, true, logEntry["bool_field"])

	// Check the Any field
	anyField, ok := logEntry["any_field"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "value", anyField["key"])
}

func TestLogRotation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping log rotation test in short mode")
	}

	// Create a temporary directory for test log files
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a log file path in the temp directory
	logFilePath := filepath.Join(tempDir, "test_rotation.log")

	// Configure logger with rotation settings
	config := Config{
		Level:      "info",
		Encoding:   "json",
		OutputPath: logFilePath,
		DevMode:    false,
		MaxSize:    1,    // 1MB max size to trigger rotation easily
		MaxBackups: 3,    // Keep 3 backups
		MaxAge:     1,    // 1 day retention
		Compress:   true, // Compress old logs
	}

	// Reset globalLogger to ensure consistent test environment
	origLogger := globalLogger
	defer func() { globalLogger = origLogger }()
	globalLogger = nil

	// Initialize logger with rotation settings
	zapLogger := newLogger(config) // Use newLogger directly instead of InitLogger to avoid singleton issues
	require.NotNil(t, zapLogger)

	// Generate enough logs to trigger rotation (at least 1MB)
	largeMessage := strings.Repeat("a", 10000) // 10KB per message
	log := NewLogger(zapLogger)

	// Write logs in batches to help force rotation
	t.Log("Writing logs to trigger rotation...")
	for i := 0; i < 200; i++ {
		// Write 10 logs per batch
		for j := 0; j < 10; j++ {
			log.Info(fmt.Sprintf("Test message %d-%d: %s", i, j, largeMessage))
		}

		// Force sync after each batch
		err = zapLogger.Sync()
		require.NoError(t, err)

		// Check if rotation has occurred
		files, _ := filepath.Glob(filepath.Join(tempDir, "test_rotation.log.*"))
		if len(files) > 0 {
			t.Logf("Rotation occurred after %d batches, found %d rotated files", i+1, len(files))
			break
		}
	}

	// Force sync to ensure all logs are written
	err = zapLogger.Sync()
	require.NoError(t, err)

	// Verify that the log file exists
	_, err = os.Stat(logFilePath)
	assert.NoError(t, err)

	// Get file size
	mainLogInfo, err := os.Stat(logFilePath)
	require.NoError(t, err)
	t.Logf("Main log file size: %d bytes", mainLogInfo.Size())

	// Wait a moment for rotation to complete
	time.Sleep(500 * time.Millisecond)

	// Check for rotated log files
	files, err := filepath.Glob(filepath.Join(tempDir, "test_rotation.log.*"))
	assert.NoError(t, err)

	// List all files in the directory for debugging
	allFiles, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	for _, file := range allFiles {
		info, err := file.Info()
		if err != nil {
			t.Logf("Found file: %s (error getting size: %v)", file.Name(), err)
		} else {
			t.Logf("Found file: %s (size: %d)", file.Name(), info.Size())
		}
	}

	// Skip the rotation test on Windows if no rotation occurred
	// This is due to Windows file locking that can sometimes prevent rotation
	if len(files) == 0 && runtime.GOOS == "windows" {
		t.Skip("Skipping rotation assertion on Windows, as file locking may prevent rotation during tests")
	} else {
		// We should have at least one rotated file
		assert.Greater(t, len(files), 0, "Expected at least one rotated log file")
	}
}
