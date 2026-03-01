package domain_test

import (
	"strings"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"empty string", "", true},
		{"missing @", "userexample.com", true},
		{"missing domain", "user@", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail(%q) error = %v, wantErr = %v", tt.email, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		pw      string
		wantErr bool
	}{
		{"valid 8 chars", "abcd1234", false},
		{"valid 72 chars", strings.Repeat("a", 72), false},
		{"too short 7 chars", "abcd123", true},
		{"too long 73 chars", strings.Repeat("a", 73), true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidatePassword(tt.pw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword(%q) error = %v, wantErr = %v", tt.pw, err, tt.wantErr)
			}
		})
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "securepassword"
	hash, err := domain.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == password {
		t.Error("hash should not equal plaintext")
	}

	if err := domain.CheckPassword(hash, password); err != nil {
		t.Errorf("CheckPassword() should succeed for correct password, got error: %v", err)
	}

	if err := domain.CheckPassword(hash, "wrongpassword"); err == nil {
		t.Error("CheckPassword() should fail for wrong password")
	}
}

func TestGenerateVerificationToken(t *testing.T) {
	token, err := domain.GenerateVerificationToken()
	if err != nil {
		t.Fatalf("GenerateVerificationToken() error = %v", err)
	}

	// 32 bytes = 64 hex chars
	if len(token) != 64 {
		t.Errorf("token length = %d, want 64", len(token))
	}

	// Should be different each time
	token2, _ := domain.GenerateVerificationToken()
	if token == token2 {
		t.Error("two generated tokens should not be identical")
	}
}

func TestTokenMatches(t *testing.T) {
	t.Run("returns true for matching tokens", func(t *testing.T) {
		token := "abc123def456"
		if !domain.TokenMatches(token, token) {
			t.Error("identical tokens should match")
		}
	})

	t.Run("returns false for different tokens", func(t *testing.T) {
		if domain.TokenMatches("abc", "xyz") {
			t.Error("different tokens should not match")
		}
	})

	t.Run("returns false for different length tokens", func(t *testing.T) {
		if domain.TokenMatches("short", "longer_token") {
			t.Error("different length tokens should not match")
		}
	})
}

func TestNewUserID(t *testing.T) {
	id := domain.NewUserID()
	if !strings.HasPrefix(id, "usr_") {
		t.Errorf("NewUserID() = %q, want prefix \"usr_\"", id)
	}
}
