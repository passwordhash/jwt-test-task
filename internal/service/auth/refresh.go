package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var _ RefreshTokenGenerator = (*RefreshTokenManager)(nil)

// RefreshTokenManager is a struct that provides methods to manage refresh tokens.
// Can be mocked for testing purposes.
type RefreshTokenManager struct{}

// Generate is a method that creates a new random refresh token of the
// specified length encoded in base64 URL format.
func (RefreshTokenManager) Generate(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Hash is a method that hashes the provided refresh token using bcrypt.
func (RefreshTokenManager) Hash(token string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash refresh token: %w", err)
	}

	return string(h), nil
}
