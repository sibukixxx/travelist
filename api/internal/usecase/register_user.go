package usecase

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
)

// UserRegistrar handles user registration.
type UserRegistrar struct {
	repo  repo.UserRepository
	clock clock.Clock
}

func NewUserRegistrar(repository repo.UserRepository, clk clock.Clock) *UserRegistrar {
	return &UserRegistrar{
		repo:  repository,
		clock: clk,
	}
}

func (ur *UserRegistrar) Register(ctx context.Context, email string) (*domain.User, error) {
	normalizedEmail, err := normalizeAndValidateEmail(email)
	if err != nil {
		return nil, err
	}

	now := ur.clock.Now().UTC()
	user := &domain.User{
		ID:        fmt.Sprintf("usr_%d", now.UnixNano()),
		Email:     normalizedEmail,
		CreatedAt: now,
	}

	if err := ur.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func normalizeAndValidateEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", fmt.Errorf("email is required")
	}

	addr, err := mail.ParseAddress(normalized)
	if err != nil || addr.Address != normalized {
		return "", fmt.Errorf("invalid email format")
	}
	return normalized, nil
}

func IsEmailAlreadyExistsError(err error) bool {
	return err == repo.ErrUserAlreadyExists
}
