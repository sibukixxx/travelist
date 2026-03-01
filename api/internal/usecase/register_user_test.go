package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

// --- Test Doubles ---

type stubUserRepo struct {
	users             []*domain.User
	createdUser       *domain.User
	updatedUser       *domain.User
	createErr         error
	findByEmailResult *domain.User
	findByTokenResult *domain.User
}

func (s *stubUserRepo) Create(_ context.Context, user *domain.User) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.createdUser = user
	s.users = append(s.users, user)
	return nil
}

func (s *stubUserRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	return s.findByEmailResult, nil
}

func (s *stubUserRepo) FindByVerificationToken(_ context.Context, token string) (*domain.User, error) {
	return s.findByTokenResult, nil
}

func (s *stubUserRepo) Update(_ context.Context, user *domain.User) error {
	s.updatedUser = user
	return nil
}

type stubEmailSender struct {
	sentTo    string
	sentToken string
	err       error
}

func (s *stubEmailSender) SendVerificationEmail(_ context.Context, to, token string) error {
	s.sentTo = to
	s.sentToken = token
	return s.err
}

// --- Register Tests ---

func TestUserRegistrarRegister(t *testing.T) {
	fixedTime := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	clk := clock.FixedClock{Time: fixedTime}

	t.Run("registers user successfully", func(t *testing.T) {
		repo := &stubUserRepo{}
		emailSender := &stubEmailSender{}
		registrar := usecase.NewUserRegistrar(repo, emailSender, clk)

		result, err := registrar.Register(context.Background(), usecase.RegisterRequest{
			Email:    "user@example.com",
			Password: "securepass",
		})

		if err != nil {
			t.Fatalf("Register() error = %v", err)
		}
		if result.Email != "user@example.com" {
			t.Errorf("Email = %q, want %q", result.Email, "user@example.com")
		}
		if !strings.HasPrefix(result.UserID, "usr_") {
			t.Errorf("UserID = %q, want prefix \"usr_\"", result.UserID)
		}

		// User was saved
		if repo.createdUser == nil {
			t.Fatal("expected user to be created")
		}
		if repo.createdUser.Email != "user@example.com" {
			t.Errorf("created email = %q, want %q", repo.createdUser.Email, "user@example.com")
		}
		// Password was hashed
		if repo.createdUser.PasswordHash == "securepass" {
			t.Error("password should be hashed, not stored as plaintext")
		}
		// Token was generated
		if len(repo.createdUser.VerificationToken) != 64 {
			t.Errorf("token length = %d, want 64", len(repo.createdUser.VerificationToken))
		}
		// Token expires in 24h
		if repo.createdUser.TokenExpiresAt == nil {
			t.Fatal("TokenExpiresAt should be set")
		}
		wantExpiry := fixedTime.Add(24 * time.Hour)
		if !repo.createdUser.TokenExpiresAt.Equal(wantExpiry) {
			t.Errorf("TokenExpiresAt = %v, want %v", repo.createdUser.TokenExpiresAt, wantExpiry)
		}
		// Verification email was sent
		if emailSender.sentTo != "user@example.com" {
			t.Errorf("email sent to %q, want %q", emailSender.sentTo, "user@example.com")
		}
	})

	t.Run("returns validation error for invalid email", func(t *testing.T) {
		registrar := usecase.NewUserRegistrar(&stubUserRepo{}, &stubEmailSender{}, clk)

		_, err := registrar.Register(context.Background(), usecase.RegisterRequest{
			Email:    "invalid",
			Password: "securepass",
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected AppError, got %T", err)
		}
		if appErr.ErrCode != apperror.ValidationError {
			t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.ValidationError)
		}
	})

	t.Run("returns validation error for short password", func(t *testing.T) {
		registrar := usecase.NewUserRegistrar(&stubUserRepo{}, &stubEmailSender{}, clk)

		_, err := registrar.Register(context.Background(), usecase.RegisterRequest{
			Email:    "user@example.com",
			Password: "short",
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected AppError, got %T", err)
		}
		if appErr.ErrCode != apperror.ValidationError {
			t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.ValidationError)
		}
	})

	t.Run("returns conflict error for duplicate email", func(t *testing.T) {
		existing := &domain.User{ID: "usr_existing", Email: "existing@example.com"}
		repo := &stubUserRepo{findByEmailResult: existing}
		registrar := usecase.NewUserRegistrar(repo, &stubEmailSender{}, clk)

		_, err := registrar.Register(context.Background(), usecase.RegisterRequest{
			Email:    "existing@example.com",
			Password: "securepass",
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected AppError, got %T", err)
		}
		if appErr.ErrCode != apperror.Conflict {
			t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.Conflict)
		}
	})

	t.Run("returns internal error when email sending fails", func(t *testing.T) {
		repo := &stubUserRepo{}
		emailSender := &stubEmailSender{err: errors.New("SMTP down")}
		registrar := usecase.NewUserRegistrar(repo, emailSender, clk)

		_, err := registrar.Register(context.Background(), usecase.RegisterRequest{
			Email:    "user@example.com",
			Password: "securepass",
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
