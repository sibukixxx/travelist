package repo

import (
	"context"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

// UserRepository defines the persistence interface for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
