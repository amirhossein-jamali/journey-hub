# Journey Hub Configuration System

This package implements a flexible configuration system for the Journey Hub application with the following features:

## Features

- **Multiple Environment Support**: Separate configuration files for development, test, and production environments
- **Environment Variable Override**: All configuration settings can be overridden by environment variables
- **Secure Storage**: Sensitive configuration values can be encrypted in configuration files
- **Hot Reload**: Configuration changes are automatically detected and applied without application restart
- **Validation**: Configuration values are validated when loaded to ensure correctness

## Directory Structure

```
config/
├── development/      # Development environment configuration
│   └── config.yaml
├── test/             # Test environment configuration
│   └── config.yaml
├── production/       # Production environment configuration
│   └── config.yaml
├── config.go         # Main configuration loader
├── secret_config.go  # Encryption support for sensitive values
└── README.md         # This file
```

## Usage

### Loading Configuration

```go
import "github.com/amirhossein-jamali/journey-hub/config"

func main() {
    // Load configuration from the appropriate environment directory
    configPath := config.GetConfigFilePath()
    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Configuration is also available via global variable
    port := config.GlobalConfig.Server.Port
}
```

### Environment Variables

All configuration settings can be overridden using environment variables with the `JH_` prefix.
Nested settings use underscore to separate levels:

```
JH_SERVER_PORT=9090
JH_DATABASE_NAME=custom_db
JH_JWT_SECRET=my-secret-key
```

### Hot Reload

Configuration is automatically reloaded when the config file changes. To register a callback function that will be called after configuration reload:

```go
config.RegisterConfigChangeCallback(func() {
    // This function will be called when config changes
    log.Println("Configuration changed, updating services...")
    
    // Update services with new configuration
    updateLogLevel(config.GlobalConfig.Log.Level)
    reconnectDatabase(config.GlobalConfig.Database)
})
```

### Secure Configuration Values

Sensitive configuration values like passwords and secrets can be encrypted in the configuration files. The encryption key can be provided via the `JH_CONFIG_ENCRYPTION_KEY` environment variable.

To encrypt a value, use the encrypt utility:

```bash
# Using command-line flag
go run cmd/encrypt/main.go -value "my-secret-password"

# Using interactive mode (more secure)
go run cmd/encrypt/main.go -i
```

Then add the encrypted value to your config file with the `ENC:` prefix:

```yaml
database:
  password: "ENC:AES-GCM-ENCRYPTED-VALUE"
```

The system will automatically decrypt these values when loading the configuration.

## Security Considerations

- The encryption key (`JH_CONFIG_ENCRYPTION_KEY`) should be managed securely and not checked into version control
- For production, consider using a secure secrets management solution like HashiCorp Vault or AWS Secrets Manager
- Always use encrypted values for sensitive information in configuration files

## Configuration Reference

### App Configuration

```yaml
app:
  name: "Journey Hub"          # Application name
  version: "0.1.0"             # Application version
  environment: "development"   # Environment (development, test, production)
  debug: true                  # Enable debug mode
```

### Server Configuration

```yaml
server:
  host: "0.0.0.0"              # Host to bind to
  port: 8080                   # Port to listen on
  read_timeout: "15s"          # HTTP read timeout
  write_timeout: "15s"         # HTTP write timeout
  idle_timeout: "60s"          # HTTP idle timeout
  allow_origins: ["*"]         # CORS allowed origins
```

### Database Configuration

...and so on for each configuration section... 