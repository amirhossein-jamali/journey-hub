package logger

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/amirhossein-jamali/journey-hub/config"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Config holds the configuration for the logger
type Config struct {
	Level      string `mapstructure:"level"`
	Encoding   string `mapstructure:"encoding"`
	OutputPath string `mapstructure:"output_path"`
	DevMode    bool   `mapstructure:"dev_mode"`
	// Log rotation settings
	MaxSize    int  `mapstructure:"max_size"`    // maximum size in megabytes before log is rotated
	MaxBackups int  `mapstructure:"max_backups"` // maximum number of old log files to retain
	MaxAge     int  `mapstructure:"max_age"`     // maximum number of days to retain old log files
	Compress   bool `mapstructure:"compress"`    // whether to compress old log files
}

// NewLoggerConfigFromGlobal creates a logger config from the global application config
func NewLoggerConfigFromGlobal() Config {
	logCfg := config.GlobalConfig.Log

	// Convert from app's LogConfig to logger's Config
	return Config{
		Level:      logCfg.Level,
		Encoding:   logCfg.Format, // Map format to encoding
		OutputPath: determineOutputPath(logCfg),
		DevMode:    config.GlobalConfig.App.Debug,
		MaxSize:    logCfg.MaxSize,
		MaxBackups: logCfg.MaxBackups,
		MaxAge:     logCfg.MaxAge,
		Compress:   logCfg.Compress,
	}
}

// determineOutputPath decides the output path based on log configuration
func determineOutputPath(logCfg config.LogConfig) string {
	if logCfg.EnableFile {
		return logCfg.FilePath
	}
	return logCfg.Output
}

// InitLogger initializes the global logger with the given configuration
func InitLogger(cfg Config) *zap.Logger {
	once.Do(func() {
		globalLogger = newLogger(cfg)
	})
	return globalLogger
}

// GetLogger returns the globally initialized zap logger
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		// Return a default logger if not initialized
		return newLogger(Config{
			Level:      "info",
			Encoding:   "console",
			OutputPath: "stdout",
			DevMode:    true,
			// Default rotation settings
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		})
	}
	return globalLogger
}

// newLogger creates a new zap logger with the given configuration
func newLogger(cfg Config) *zap.Logger {
	// Determine log level
	var level zapcore.Level

	// Parse log level from config string
	parsedLevel, err := zapcore.ParseLevel(cfg.Level)
	if err == nil {
		level = parsedLevel
	} else {
		// Default to info level if invalid level provided
		level = zapcore.InfoLevel
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create encoder based on configuration
	var encoder zapcore.Encoder
	if cfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output path with rotation support
	var writeSyncer zapcore.WriteSyncer
	switch cfg.OutputPath {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	default:
		// Use lumberjack for log rotation when writing to files
		writeSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    cfg.MaxSize,    // megabytes
			MaxBackups: cfg.MaxBackups, // number of backups
			MaxAge:     cfg.MaxAge,     // days
			Compress:   cfg.Compress,   // compress old files
		})
	}

	// Create the core
	core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(level))

	// Create logger
	var logger *zap.Logger
	if cfg.DevMode {
		// Development mode logger with caller and stacktrace
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		// Production mode logger
		logger = zap.New(core)
	}

	// Add a timestamp to log initialization
	logger.Info("Logger initialized with rotation",
		zap.Time("time", time.Now()),
		zap.String("path", cfg.OutputPath),
		zap.Int("maxSize", cfg.MaxSize),
		zap.Int("maxBackups", cfg.MaxBackups),
		zap.Int("maxAge", cfg.MaxAge),
		zap.Bool("compress", cfg.Compress))

	return logger
}

// Field is a wrapper for zap.Field to make logger usage more intuitive
type Field = zap.Field

// Common field creation functions for use with the logger
var (
	String  = zap.String
	Int     = zap.Int
	Int64   = zap.Int64
	Float64 = zap.Float64
	Bool    = zap.Bool
	Error   = zap.Error
	Any     = zap.Any
	Time    = zap.Time
)
