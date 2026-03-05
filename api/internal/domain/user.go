package domain

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/mail"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a registered user.
type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty"`
	VerificationToken string     `json:"-"`
	TokenExpiresAt    *time.Time `json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// NewUserID generates a unique user ID with "usr_" prefix.
func NewUserID() string {
	return fmt.Sprintf("usr_%d", time.Now().UnixNano())
}

// ValidateEmail checks that the email is well-formed.
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}
	return nil
}

// ValidatePassword checks that the password is between 8 and 72 characters.
// The 72-byte limit is enforced because bcrypt silently truncates beyond that.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters")
	}
	return nil
}

// HashPassword hashes a password using bcrypt with cost 12.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword compares a bcrypt hash with a plaintext password.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateVerificationToken generates a 32-byte (64 hex chars) cryptographically random token.
func GenerateVerificationToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// TokenMatches compares two tokens using constant-time comparison.
func TokenMatches(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
