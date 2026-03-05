package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/infra/repo"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

type stubUserRepo struct {
	users map[string]domain.User
}

func (r *stubUserRepo) Create(_ context.Context, user *domain.User) error {
	if _, ok := r.users[user.Email]; ok {
		return repo.ErrUserAlreadyExists
	}
	r.users[user.Email] = *user
	return nil
}

func (r *stubUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	user, ok := r.users[email]
	if !ok {
		return nil, nil
	}
	u := user
	return &u, nil
}

func TestUserRegistrarRegister(t *testing.T) {
	t.Run("registers normalized email", func(t *testing.T) {
		repository := &stubUserRepo{users: map[string]domain.User{}}
		fixed := clock.FixedClock{Time: time.Date(2026, 3, 1, 1, 2, 3, 0, time.UTC)}
		registrar := usecase.NewUserRegistrar(repository, fixed)

		user, err := registrar.Register(context.Background(), "  Foo.Bar+test@Example.COM ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if user.Email != "foo.bar+test@example.com" {
			t.Fatalf("Email = %q, want normalized", user.Email)
		}
		if user.CreatedAt != fixed.Time.UTC() {
			t.Fatalf("CreatedAt = %v, want %v", user.CreatedAt, fixed.Time.UTC())
		}
	})

	t.Run("rejects invalid email", func(t *testing.T) {
		repository := &stubUserRepo{users: map[string]domain.User{}}
		registrar := usecase.NewUserRegistrar(repository, clock.FixedClock{Time: time.Now()})

		_, err := registrar.Register(context.Background(), "invalid-email")
		if err == nil {
			t.Fatal("expected error for invalid email")
		}
	})

	t.Run("rejects duplicated email", func(t *testing.T) {
		repository := &stubUserRepo{
			users: map[string]domain.User{
				"user@example.com": {Email: "user@example.com"},
			},
		}
		registrar := usecase.NewUserRegistrar(repository, clock.FixedClock{Time: time.Now()})

		_, err := registrar.Register(context.Background(), "user@example.com")
		if !usecase.IsEmailAlreadyExistsError(err) {
			t.Fatalf("expected already exists error, got: %v", err)
		}
	})
}
