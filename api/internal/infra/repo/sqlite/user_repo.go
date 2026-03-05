package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/domain"
)

// UserRepo implements repo.UserRepository using SQLite.
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user. Returns apperror.Conflict if the email already exists.
func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, email, password_hash, email_verified_at, verification_token, token_expires_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash,
		nullableTime(user.EmailVerifiedAt),
		user.VerificationToken,
		nullableTime(user.TokenExpiresAt),
		user.CreatedAt.UTC(), user.UpdatedAt.UTC(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return apperror.NewConflict("email already registered")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// FindByEmail returns a user by email, or nil if not found.
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findOne(ctx, "SELECT id, email, password_hash, email_verified_at, verification_token, token_expires_at, created_at, updated_at FROM users WHERE email = ?", email)
}

// FindByVerificationToken returns a user by verification token, or nil if not found.
func (r *UserRepo) FindByVerificationToken(ctx context.Context, token string) (*domain.User, error) {
	return r.findOne(ctx, "SELECT id, email, password_hash, email_verified_at, verification_token, token_expires_at, created_at, updated_at FROM users WHERE verification_token = ?", token)
}

// Update updates an existing user.
func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET email = ?, password_hash = ?, email_verified_at = ?, verification_token = ?, token_expires_at = ?, updated_at = ?
		 WHERE id = ?`,
		user.Email, user.PasswordHash,
		nullableTime(user.EmailVerifiedAt),
		user.VerificationToken,
		nullableTime(user.TokenExpiresAt),
		user.UpdatedAt.UTC(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepo) findOne(ctx context.Context, query string, args ...any) (*domain.User, error) {
	var u domain.User
	var emailVerifiedAt, tokenExpiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PasswordHash,
		&emailVerifiedAt, &u.VerificationToken, &tokenExpiresAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if emailVerifiedAt.Valid {
		t := emailVerifiedAt.Time
		u.EmailVerifiedAt = &t
	}
	if tokenExpiresAt.Valid {
		t := tokenExpiresAt.Time
		u.TokenExpiresAt = &t
	}

	return &u, nil
}

func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC()
}
