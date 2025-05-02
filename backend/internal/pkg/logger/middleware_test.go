package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMiddleware(t *testing.T) {
	// Create an in-memory logger for testing
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	zapLogger := zap.New(core)

	// Create a test handler that will be wrapped by our middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	})

	// Create middleware
	middleware := Middleware(zapLogger)

	// Apply middleware to our test handler
	handler := middleware(testHandler)

	// Create test server with the middleware-wrapped handler
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a test request
	resp, err := http.Get(server.URL + "/test?param=value")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Parse the log output to verify request was logged
	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify logged fields
	assert.Equal(t, "HTTP request", logEntry["message"])
	assert.Equal(t, "INFO", logEntry["level"])
	assert.Equal(t, "GET", logEntry["method"])
	assert.Equal(t, "/test", logEntry["path"])
	assert.Equal(t, "param=value", logEntry["query"])
	assert.Equal(t, float64(200), logEntry["status"]) // JSON numbers are float64
}

func TestResponseWriterMethods(t *testing.T) {
	// Create a test implementation of responseWriter
	origWriter := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: origWriter,
		statusCode:     0,
	}

	// Test WriteHeader
	rw.WriteHeader(http.StatusCreated)
	assert.Equal(t, http.StatusCreated, rw.statusCode)
	assert.Equal(t, http.StatusCreated, origWriter.Code)

	// Test Write when status is already set
	n, err := rw.Write([]byte("test"))
	require.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, "test", origWriter.Body.String())
	assert.Equal(t, http.StatusCreated, rw.statusCode) // Status should remain the same

	// Test with a new responseWriter (reset for this test)
	origWriter = httptest.NewRecorder()
	rw = &responseWriter{
		ResponseWriter: origWriter,
		statusCode:     0,
	}

	// Test Write when status is not set (should default to 200 OK)
	n, err = rw.Write([]byte("test"))
	require.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, "test", origWriter.Body.String())
	assert.Equal(t, http.StatusOK, rw.statusCode) // Status should be set to 200 OK
}

func TestMiddlewareWithErrorResponse(t *testing.T) {
	// Create an in-memory logger
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		TimeKey:     "time",
		EncodeLevel: zapcore.CapitalLevelEncoder,
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	zapLogger := zap.New(core)

	// Create a test handler that returns an error status
	errorHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	})

	// Apply middleware
	handler := Middleware(zapLogger)(errorHandler)

	// Create test request
	req := httptest.NewRequest("POST", "/api/data", bytes.NewBufferString("test data"))
	req.Header.Set("User-Agent", "test-agent")

	// Record the response
	recorder := httptest.NewRecorder()

	// Process the request
	handler.ServeHTTP(recorder, req)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	// Parse the log
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify logged fields
	assert.Equal(t, "HTTP request", logEntry["message"])
	assert.Equal(t, "POST", logEntry["method"])
	assert.Equal(t, "/api/data", logEntry["path"])
	assert.Equal(t, "test-agent", logEntry["user_agent"])
	assert.Equal(t, float64(500), logEntry["status"]) // Status should be 500
}
