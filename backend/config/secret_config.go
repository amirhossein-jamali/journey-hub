// Package config provides configuration management using viper, supporting
// different environments, hot reloading, and secure configuration storage.
package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/amirhossein-jamali/journey-hub/internal/pkg/crypto"
)

const (
	// EncryptedPrefix is the prefix for encrypted values in config files
	EncryptedPrefix = "ENC:"

	// EnvKeyName is the environment variable that contains encryption key
	EnvKeyName = "JH_CONFIG_ENCRYPTION_KEY"
)

var (
	encryptor *crypto.Encryptor
)

// InitEncryption initializes the encryption system for config values
func InitEncryption() error {
	// Check for CI/CD environment
	if os.Getenv("JH_SKIP_CONFIG_ENCRYPTION") == "true" {
		fmt.Println("Config encryption disabled in CI/CD environment.")
		return nil
	}

	// Get encryption key from environment
	key := os.Getenv(EnvKeyName)
	if key == "" {
		// Check for specific environments
		env := os.Getenv("JH_ENV")
		if env == "production" || env == "staging" {
			return fmt.Errorf("encryption key must be provided in production/staging via %s", EnvKeyName)
		}

		// Generate a random key for development purposes
		// In production, the key should be properly managed and provided via environment
		key = base64.StdEncoding.EncodeToString([]byte("journeyhub-dev-encryption-key-32byte"))
		fmt.Printf("Warning: Using default development encryption key. For production, set %s environment variable.\n", EnvKeyName)
	}

	// Decode key if it's base64 encoded
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		// If not base64, use directly
		keyBytes = []byte(key)
	}

	// Ensure key has proper length (trim or pad as needed)
	keyBytes = adjustKeyLength(keyBytes)

	// Initialize encryptor
	encryptor, err = crypto.NewEncryptor(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to initialize config encryption: %w", err)
	}

	return nil
}

// adjustKeyLength ensures the key is 16, 24, or 32 bytes (for AES-128, AES-192, or AES-256)
func adjustKeyLength(key []byte) []byte {
	switch {
	case len(key) < 16:
		// Pad to 16 bytes
		newKey := make([]byte, 16)
		copy(newKey, key)
		return newKey
	case len(key) < 24:
		// Use as-is (AES-128)
		return key[:16]
	case len(key) < 32:
		// Use as-is (AES-192)
		return key[:24]
	default:
		// Use as-is (AES-256)
		return key[:32]
	}
}

// DecryptConfigValues decrypts any encrypted values in the configuration
func DecryptConfigValues(config *Config) error {
	// Skip in CI/CD environments
	if os.Getenv("JH_SKIP_CONFIG_ENCRYPTION") == "true" {
		return nil
	}

	if encryptor == nil {
		if err := InitEncryption(); err != nil {
			return err
		}
	}

	// Decrypt database password
	if strings.HasPrefix(config.Database.Password, EncryptedPrefix) {
		encValue := strings.TrimPrefix(config.Database.Password, EncryptedPrefix)
		decrypted, err := encryptor.Decrypt(encValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt database password: %w", err)
		}
		config.Database.Password = decrypted
	}

	// Decrypt Redis password
	if strings.HasPrefix(config.Redis.Password, EncryptedPrefix) {
		encValue := strings.TrimPrefix(config.Redis.Password, EncryptedPrefix)
		decrypted, err := encryptor.Decrypt(encValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt Redis password: %w", err)
		}
		config.Redis.Password = decrypted
	}

	// Decrypt MinIO secret key
	if strings.HasPrefix(config.Storage.SecretAccessKey, EncryptedPrefix) {
		encValue := strings.TrimPrefix(config.Storage.SecretAccessKey, EncryptedPrefix)
		decrypted, err := encryptor.Decrypt(encValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt MinIO secret key: %w", err)
		}
		config.Storage.SecretAccessKey = decrypted
	}

	// Decrypt JWT secret
	if strings.HasPrefix(config.JWT.Secret, EncryptedPrefix) {
		encValue := strings.TrimPrefix(config.JWT.Secret, EncryptedPrefix)
		decrypted, err := encryptor.Decrypt(encValue)
		if err != nil {
			return fmt.Errorf("failed to decrypt JWT secret: %w", err)
		}
		config.JWT.Secret = decrypted
	}

	return nil
}

// EncryptConfigValue encrypts a string value
func EncryptConfigValue(value string) (string, error) {
	// Skip in CI/CD environments
	if os.Getenv("JH_SKIP_CONFIG_ENCRYPTION") == "true" {
		// Return a mock encrypted value for CI/CD
		return EncryptedPrefix + "MOCK_ENCRYPTED_VALUE_FOR_CI", nil
	}

	if encryptor == nil {
		if err := InitEncryption(); err != nil {
			return "", err
		}
	}

	encrypted, err := encryptor.Encrypt(value)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt value: %w", err)
	}

	return EncryptedPrefix + encrypted, nil
}

// RegisterDecryptionHook registers a hook to decrypt config values when loaded
func RegisterDecryptionHook() {
	// Skip in CI/CD environments
	if os.Getenv("JH_SKIP_CONFIG_ENCRYPTION") == "true" {
		fmt.Println("Skipping config encryption/decryption hooks in CI/CD environment")
		return
	}

	// Initialize encryption
	if err := InitEncryption(); err != nil {
		fmt.Printf("Warning: Failed to initialize encryption: %s\n", err)
		return
	}

	// Register a hook to process values as they are read
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Config file changed, ensure decryption is applied
		if err := DecryptConfigValues(GlobalConfig); err != nil {
			fmt.Printf("Error decrypting config values after change: %s\n", err)
		}
	})
}
