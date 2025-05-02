package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptor(t *testing.T) {
	// Test key must be 16, 24, or 32 bytes (AES-128, AES-192, or AES-256)
	key := []byte("0123456789ABCDEF") // 16 bytes for AES-128

	// Test with invalid key length
	invalidKey := []byte("too-short")
	_, err := NewEncryptor(invalidKey)
	assert.Error(t, err, "Should error with invalid key length")

	// Create encryptor
	encryptor, err := NewEncryptor(key)
	assert.NoError(t, err, "Should create encryptor without error")
	assert.NotNil(t, encryptor, "Encryptor should not be nil")

	// Test encryption and decryption
	testCases := []struct {
		name     string
		input    string
		wantSame bool
	}{
		{
			name:     "Simple string",
			input:    "hello world",
			wantSame: true,
		},
		{
			name:     "Empty string",
			input:    "",
			wantSame: true,
		},
		{
			name:     "Special characters",
			input:    "!@#$%^&*()_+{}:<>?[];\",./",
			wantSame: true,
		},
		{
			name:     "Unicode characters",
			input:    "こんにちは世界",
			wantSame: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := encryptor.Encrypt(tc.input)
			assert.NoError(t, err, "Encryption should not error")

			// Check if encrypted text is different from original
			assert.NotEqual(t, tc.input, encrypted, "Encrypted text should be different from original")

			// Check if IsEncrypted recognizes this as encrypted
			assert.True(t, IsEncrypted(encrypted), "IsEncrypted should recognize encrypted text")

			// Decrypt
			decrypted, err := encryptor.Decrypt(encrypted)
			assert.NoError(t, err, "Decryption should not error")

			// Compare decrypted with original
			if tc.wantSame {
				assert.Equal(t, tc.input, decrypted, "Decrypted text should match original")
			} else {
				assert.NotEqual(t, tc.input, decrypted, "Decrypted text should not match original")
			}
		})
	}
}

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Too short",
			input: "abc",
			want:  false,
		},
		{
			name:  "Not base64",
			input: "not!base@64",
			want:  false,
		},
		{
			name:  "Valid base64",
			input: "VGhpcyBpcyBhIHRlc3Q=", // "This is a test" in base64
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEncrypted(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
