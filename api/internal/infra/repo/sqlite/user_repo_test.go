package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/repo/sqlite"
)

func setupTestDB(t *testing.T) *sqlite.UserRepo {
	t.Helper()
	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return sqlite.NewUserRepo(db)
}

func newTestUser() *domain.User {
	now := time.Now().UTC().Truncate(time.Second)
	expires := now.Add(24 * time.Hour)
	return &domain.User{
		ID:                "usr_123",
		Email:             "test@example.com",
		PasswordHash:      "$2a$12$fakehash",
		VerificationToken: "abc123token",
		TokenExpiresAt:    &expires,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func TestUserRepoCreate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	user := newTestUser()

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify it was stored
	found, err := repo.FindByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.ID != "usr_123" {
		t.Errorf("ID = %q, want %q", found.ID, "usr_123")
	}
	if found.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", found.Email, "test@example.com")
	}
	if found.PasswordHash != "$2a$12$fakehash" {
		t.Errorf("PasswordHash = %q, want %q", found.PasswordHash, "$2a$12$fakehash")
	}
}

func TestUserRepoCreateDuplicate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	user := newTestUser()

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}

	user2 := newTestUser()
	user2.ID = "usr_456"
	err := repo.Create(ctx, user2)
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestUserRepoFindByEmailNotFound(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	found, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}
	if found != nil {
		t.Errorf("expected nil, got %v", found)
	}
}

func TestUserRepoFindByVerificationToken(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	user := newTestUser()

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByVerificationToken(ctx, "abc123token")
	if err != nil {
		t.Fatalf("FindByVerificationToken() error = %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.ID != "usr_123" {
		t.Errorf("ID = %q, want %q", found.ID, "usr_123")
	}
}

func TestUserRepoFindByVerificationTokenNotFound(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()

	found, err := repo.FindByVerificationToken(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("FindByVerificationToken() error = %v", err)
	}
	if found != nil {
		t.Errorf("expected nil, got %v", found)
	}
}

func TestUserRepoUpdate(t *testing.T) {
	repo := setupTestDB(t)
	ctx := context.Background()
	user := newTestUser()

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify email
	now := time.Now().UTC().Truncate(time.Second)
	user.EmailVerifiedAt = &now
	user.VerificationToken = ""
	user.TokenExpiresAt = nil

	if err := repo.Update(ctx, user); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	found, err := repo.FindByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("FindByEmail() error = %v", err)
	}
	if found.EmailVerifiedAt == nil {
		t.Fatal("expected EmailVerifiedAt to be set")
	}
	if found.VerificationToken != "" {
		t.Errorf("VerificationToken = %q, want empty", found.VerificationToken)
	}
}
