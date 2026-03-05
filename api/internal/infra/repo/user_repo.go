package repo

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/sibukixxx/travelist/api/internal/domain"
)

var ErrUserAlreadyExists = errors.New("user already exists")

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

// InMemoryUserRepository is a simple in-memory user store.
type InMemoryUserRepository struct {
	mu      sync.RWMutex
	byEmail map[string]domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		byEmail: make(map[string]domain.User),
	}
}

func (r *InMemoryUserRepository) Create(_ context.Context, user *domain.User) error {
	key := normalizeEmail(user.Email)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byEmail[key]; exists {
		return ErrUserAlreadyExists
	}
	r.byEmail[key] = *user
	return nil
}

func (r *InMemoryUserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	key := normalizeEmail(email)

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byEmail[key]
	if !ok {
		return nil, nil
	}
	u := user
	return &u, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
