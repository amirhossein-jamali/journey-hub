package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Storage    StorageConfig    `mapstructure:"storage"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	Swagger    SwaggerConfig    `mapstructure:"swagger"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
}

// AppConfig holds application specific configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig holds server specific configuration
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	AllowOrigins []string      `mapstructure:"allow_origins"`
}

// DatabaseConfig holds database specific configuration
type DatabaseConfig struct {
	URI                  string        `mapstructure:"uri"`
	Name                 string        `mapstructure:"name"`
	ConnectTimeout       time.Duration `mapstructure:"connect_timeout"`
	MaxConnIdleTime      time.Duration `mapstructure:"max_conn_idle_time"`
	MaxPoolSize          uint64        `mapstructure:"max_pool_size"`
	User                 string        `mapstructure:"user"`
	Password             string        `mapstructure:"password"`
	EnableCommandMonitor bool          `mapstructure:"enable_command_monitor"`
}

// RedisConfig holds redis specific configuration
type RedisConfig struct {
	URI               string        `mapstructure:"uri"`
	Password          string        `mapstructure:"password"`
	DB                int           `mapstructure:"db"`
	PoolSize          int           `mapstructure:"pool_size"`
	ConnectTimeout    time.Duration `mapstructure:"connect_timeout"`
	MaxRetries        int           `mapstructure:"max_retries"`
	MinIdleConns      int           `mapstructure:"min_idle_conns"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	EnableCachePrefix bool          `mapstructure:"enable_cache_prefix"`
	CachePrefix       string        `mapstructure:"cache_prefix"`
}

// StorageConfig holds MinIO/S3 storage configuration
type StorageConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
	Location        string `mapstructure:"location"`
	BasePath        string `mapstructure:"base_path"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl"`
	Issuer           string        `mapstructure:"issuer"`
	Audience         string        `mapstructure:"audience"`
	SigningAlgorithm string        `mapstructure:"signing_algorithm"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	EnableFile bool   `mapstructure:"enable_file"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`    // megabytes
	MaxBackups int    `mapstructure:"max_backups"` // number of backups
	MaxAge     int    `mapstructure:"max_age"`     // days
	Compress   bool   `mapstructure:"compress"`
}

// SwaggerConfig holds Swagger API documentation configuration
type SwaggerConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Host        string `mapstructure:"host"`
	BasePath    string `mapstructure:"base_path"`
	Schemes     string `mapstructure:"schemes"`
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
}

// MonitoringConfig holds monitoring configuration
type MonitoringConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	PrometheusPort  int    `mapstructure:"prometheus_port"`
	HealthCheckPath string `mapstructure:"health_check_path"`
}

// Global instance of configuration that can be accessed from other packages
var GlobalConfig Config

// configChangeCallbacks holds functions to be called when config changes
var configChangeCallbacks []func()

// RegisterConfigChangeCallback registers a function to be called when config changes
func RegisterConfigChangeCallback(callback func()) {
	configChangeCallbacks = append(configChangeCallbacks, callback)
}

// LoadConfig reads config from file or environment
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	viper.SetConfigName("config")   // Name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(configPath) // Path to look for the config file in

	// Set default values
	setDefaults()

	// Check if encryption should be skipped (for CI/CD environments)
	if os.Getenv("JH_SKIP_CONFIG_ENCRYPTION") != "true" {
		// Initialize encryption for secure config values
		if err := InitEncryption(); err != nil {
			return nil, fmt.Errorf("failed to initialize config encryption: %w", err)
		}
	}

	// Try to read configuration from file
	if err := viper.ReadInConfig(); err != nil {
		// Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; using defaults and environment variables
			fmt.Println("Config file not found, using defaults and environment variables")
		} else {
			// Config file was found but another error was produced
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	// Override with environment variables
	viper.SetEnvPrefix("JH") // Prefix for environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // Read in environment variables that match

	// Unmarshal config into struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Decrypt sensitive configuration values only if encryption is not skipped
	if os.Getenv("JH_SKIP_CONFIG_ENCRYPTION") != "true" {
		if err := DecryptConfigValues(config); err != nil {
			return nil, fmt.Errorf("failed to decrypt config values: %w", err)
		}
	}

	// Store the config in the global variable
	GlobalConfig = *config

	// Setup hot-reload of configuration
	setupConfigHotReload()

	return config, nil
}

// setupConfigHotReload configures viper to watch for config file changes
func setupConfigHotReload() {
	// Watch for changes in the config file and automatically reread
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)

		// Create a new config instance
		newConfig := &Config{}

		// Unmarshal into the new config
		if err := viper.Unmarshal(newConfig); err != nil {
			fmt.Printf("Error unmarshaling updated config: %s\n", err)
			return
		}

		// Decrypt sensitive values
		if err := DecryptConfigValues(newConfig); err != nil {
			fmt.Printf("Error decrypting updated config values: %s\n", err)
			return
		}

		// Update the global config
		GlobalConfig = *newConfig
		fmt.Println("Config reloaded successfully")

		// Call all registered callbacks
		for _, callback := range configChangeCallbacks {
			callback()
		}
	})
}

// setDefaults sets default values for configuration
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "Journey Hub")
	viper.SetDefault("app.version", "0.1.0")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)

	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")
	viper.SetDefault("server.allow_origins", []string{"*"})

	// Database defaults
	viper.SetDefault("database.uri", "mongodb://localhost:27017")
	viper.SetDefault("database.name", "journeyhub")
	viper.SetDefault("database.connect_timeout", "10s")
	viper.SetDefault("database.max_conn_idle_time", "60s")
	viper.SetDefault("database.max_pool_size", 100)
	viper.SetDefault("database.enable_command_monitor", false)

	// Redis defaults
	viper.SetDefault("redis.uri", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.connect_timeout", "5s")
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.min_idle_conns", 2)
	viper.SetDefault("redis.idle_timeout", "5m")
	viper.SetDefault("redis.enable_cache_prefix", true)
	viper.SetDefault("redis.cache_prefix", "jh:")

	// Storage defaults
	viper.SetDefault("storage.endpoint", "localhost:9000")
	viper.SetDefault("storage.use_ssl", false)
	viper.SetDefault("storage.bucket_name", "journeyhub")
	viper.SetDefault("storage.location", "us-east-1")
	viper.SetDefault("storage.base_path", "uploads")

	// JWT defaults
	viper.SetDefault("jwt.access_token_ttl", "15m")
	viper.SetDefault("jwt.refresh_token_ttl", "7d")
	viper.SetDefault("jwt.issuer", "journeyhub")
	viper.SetDefault("jwt.audience", "journeyhub-users")
	viper.SetDefault("jwt.signing_algorithm", "HS256")

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.enable_file", false)
	viper.SetDefault("log.file_path", "logs/app.log")
	viper.SetDefault("log.max_size", 10)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.compress", true)

	// Swagger defaults
	viper.SetDefault("swagger.enabled", true)
	viper.SetDefault("swagger.host", "localhost:8080")
	viper.SetDefault("swagger.base_path", "/api/v1")
	viper.SetDefault("swagger.schemes", "http,https")
	viper.SetDefault("swagger.title", "Journey Hub API")
	viper.SetDefault("swagger.description", "API for Journey Hub travel platform")
	viper.SetDefault("swagger.version", "1.0")

	// Monitoring defaults
	viper.SetDefault("monitoring.enabled", true)
	viper.SetDefault("monitoring.prometheus_port", 9090)
	viper.SetDefault("monitoring.health_check_path", "/health")
}

// GetConfigFilePath returns the path of the config file
func GetConfigFilePath() string {
	// Get environment
	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // Default to development
	}

	// Get config directory based on environment
	baseDir := "config"
	configDir := fmt.Sprintf("%s/%s", baseDir, env)

	// Check if directory exists
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// Fallback to config directory
		configDir = baseDir
	}

	return configDir
}
