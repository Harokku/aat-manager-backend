package db

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
)

// hexDecodeString is a helper function to decode a hexadecimal encoded string.
func hexDecodeString(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func TestEncryptToken(t *testing.T) {
	key := make([]byte, 32) // Create a 32 byte length slice for random secret
	_, _ = rand.Read(key)   // Randomize key

	// table driven tests
	tests := []struct {
		name        string
		plainText   string
		expectedErr error
	}{
		{"empty string", "", nil},
		{"short string", "test", nil},
		{"long string", "this is a very long string", nil},
		{"special characters", "!@#$%^&*()", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encryptToken(tt.plainText, key)

			if (err != nil) != (tt.expectedErr != nil) {
				t.Errorf("encryptToken() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestEncryptDecryptToken(t *testing.T) {
	key := make([]byte, 32)
	_, _ = rand.Read(key)

	tests := []struct {
		name      string
		plainText string
	}{
		{"empty string", ""},
		{"short string", "test"},
		{"long string", "this is a very long string"},
		{"special characters", "!@#$%^&*()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipherText, err := encryptToken(tt.plainText, key)
			if err != nil {
				t.Errorf("encryption failed: %v", err)
			}

			decryptedText, err := decryptToken(cipherText, key)
			if err != nil {
				t.Errorf("decryption failed: %v", err)
			}

			if tt.plainText != decryptedText {
				t.Errorf("decrypted text doesn't match the original one, got %v, want %v", decryptedText, tt.plainText)
			}
		})
	}
}
