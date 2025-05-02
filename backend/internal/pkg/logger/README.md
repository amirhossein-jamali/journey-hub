# Logger Package

This package provides a complete logging system using the `zap` library for recording events throughout the application.

## Package Structure

The `logger` package consists of several core files:

- `logger.go`: Main logger implementation with singleton pattern
- `wrapper.go`: Abstraction layer for ease of use
- `middleware.go`: HTTP middleware for logging requests
- `*_test.go`: Test files for each component

## Singleton Pattern

This package uses the singleton design pattern to ensure that only one instance of the logger exists throughout the application:

```go
var (
    globalLogger *zap.Logger
    once         sync.Once
)

func InitLogger(cfg Config) *zap.Logger {
    once.Do(func() {
        globalLogger = newLogger(cfg)
    })
    return globalLogger
}

func GetLogger() *zap.Logger {
    if globalLogger == nil {
        // Return a default logger if not initialized
        return newLogger(Config{...})
    }
    return globalLogger
}
```

Important aspects of the singleton pattern:

1. The `globalLogger` variable is kept private at the package level
2. `sync.Once` is used to ensure initialization happens only once
3. The `GetLogger()` function provides safe access to the singleton instance
4. If `globalLogger` is not initialized, a default instance is created

## Interface

In `wrapper.go`, a `Logger` interface is defined to simplify logger usage:

```go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    With(fields ...Field) Logger
}
```

This interface allows users to use logging capabilities without direct dependency on `zap`.

## HTTP Middleware

The `middleware.go` file provides an HTTP middleware that can be used to log all HTTP requests:

```go
func Middleware(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Log HTTP request
        })
    }
}
```

## Log Rotation

The logger supports automatic log rotation using the `lumberjack` package. This is especially useful for production environments where log files can grow very large. Log rotation allows you to:

- Limit log file size
- Keep a specific number of backup files
- Automatically remove old log files
- Compress old log files to save space

Log rotation is configured through these settings in the `Config` struct:

```go
type Config struct {
    // ... other fields ...
    
    // Log rotation settings
    MaxSize    int  `mapstructure:"max_size"`    // maximum size in megabytes
    MaxBackups int  `mapstructure:"max_backups"` // maximum number of old log files to retain
    MaxAge     int  `mapstructure:"max_age"`     // maximum number of days to retain old files
    Compress   bool `mapstructure:"compress"`    // whether to compress old log files
}
```

When logs are written to a file (not stdout/stderr), the rotation settings are applied automatically.

## Testing Strategy

The tests for the `logger` package are kept in the same directory as the main code, which is in line with Go standards. For each main file, there is a corresponding test file:

- `logger_test.go`: Tests for `logger.go`
- `wrapper_test.go`: Tests for `wrapper.go`
- `middleware_test.go`: Tests for `middleware.go`

### Benefits of keeping tests alongside the main code:

1. **Access to private variables and functions**: Tests can access package-private variables (like `globalLogger`)
2. **Clear location**: Tests are located right next to the code they're testing
3. **Adherence to Go standards**: This approach follows Go's conventional practices
4. **Ease of execution**: Tests can be run with the simple command `go test ./...`

## Usage

### Initial Setup

To use this package, you need to initialize the logger at the entry point of your application:

```go
import "github.com/amirhossein-jamali/journey-hub/internal/pkg/logger"

func main() {
    // Initialize logger with custom configuration
    cfg := logger.Config{
        Level:      "debug",
        Encoding:   "json",
        OutputPath: "logs/app.log",
        DevMode:    true,
        // Rotation settings
        MaxSize:    100,    // 100MB per file
        MaxBackups: 5,      // Keep 5 old files
        MaxAge:     30,     // 30 days retention
        Compress:   true,   // Compress old files
    }
    logger.InitLogger(cfg)
    
    // Or use global configuration
    // logger.InitLogger(logger.NewLoggerConfigFromGlobal())
}
```

### Logging with the Interface

It's recommended to use the `Logger` interface for logging:

```go
import "github.com/amirhossein-jamali/journey-hub/internal/pkg/logger"

func DoSomething() {
    log := logger.Default()
    
    // Simple logging
    log.Info("Operation completed")
    
    // Structured logging with fields
    log.Error("Error occurred", 
        logger.String("service", "payment"),
        logger.Int("status", 500),
        logger.Error(err),
    )
    
    // Create a logger with context
    contextLogger := log.With(logger.String("component", "authentication"))
    contextLogger.Debug("Validating token")
}
```

### Use in HTTP

For logging HTTP requests:

```go
import (
    "net/http"
    "github.com/amirhossein-jamali/journey-hub/internal/pkg/logger"
)

func SetupRoutes() http.Handler {
    router := http.NewServeMux()
    
    // Add logging middleware
    logMiddleware := logger.Middleware(logger.GetLogger())
    
    return logMiddleware(router)
}
```

## Enhancement and Maintenance

To add new capabilities to the logger package:

1. Implement the new feature
2. Write related tests in the corresponding `*_test.go` file
3. Update the documentation in this README 