package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encryptor provides methods to encrypt and decrypt sensitive data
type Encryptor struct {
	key []byte
}

// NewEncryptor creates a new Encryptor with the provided key
// Key should be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("encryption key must be 16, 24, or 32 bytes long")
	}
	return &Encryptor{key: key}, nil
}

// Encrypt encrypts plaintext using AES-GCM
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and append nonce
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Convert to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func (e *Encryptor) Decrypt(encryptedText string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce from the beginning of ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// IsEncrypted checks if a string appears to be encrypted (base64 encoded)
func IsEncrypted(text string) bool {
	// Check if the string starts with the encrypted prefix
	if len(text) < 8 {
		return false
	}

	// Try to decode as base64
	_, err := base64.StdEncoding.DecodeString(text)
	return err == nil
}
