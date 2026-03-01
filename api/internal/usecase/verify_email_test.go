package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
	"github.com/sibukixxx/travelist/api/internal/usecase"
)

func TestEmailVerifierVerify(t *testing.T) {
	fixedTime := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	clk := clock.FixedClock{Time: fixedTime}

	t.Run("verifies email successfully", func(t *testing.T) {
		expires := fixedTime.Add(24 * time.Hour)
		user := &domain.User{
			ID:                "usr_123",
			Email:             "user@example.com",
			VerificationToken: "validtoken123",
			TokenExpiresAt:    &expires,
		}
		repo := &stubUserRepo{findByTokenResult: user}
		verifier := usecase.NewEmailVerifier(repo, clk)

		err := verifier.Verify(context.Background(), "validtoken123")

		if err != nil {
			t.Fatalf("Verify() error = %v", err)
		}
		if repo.updatedUser == nil {
			t.Fatal("expected user to be updated")
		}
		if repo.updatedUser.EmailVerifiedAt == nil {
			t.Fatal("expected EmailVerifiedAt to be set")
		}
		if repo.updatedUser.VerificationToken != "" {
			t.Errorf("VerificationToken = %q, want empty", repo.updatedUser.VerificationToken)
		}
		if repo.updatedUser.TokenExpiresAt != nil {
			t.Errorf("TokenExpiresAt = %v, want nil", repo.updatedUser.TokenExpiresAt)
		}
	})

	t.Run("returns INVALID_TOKEN for unknown token", func(t *testing.T) {
		repo := &stubUserRepo{findByTokenResult: nil}
		verifier := usecase.NewEmailVerifier(repo, clk)

		err := verifier.Verify(context.Background(), "unknowntoken")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected AppError, got %T", err)
		}
		if appErr.ErrCode != apperror.InvalidToken {
			t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.InvalidToken)
		}
	})

	t.Run("returns EXPIRED_TOKEN for expired token", func(t *testing.T) {
		expired := fixedTime.Add(-1 * time.Hour) // Already expired
		user := &domain.User{
			ID:                "usr_123",
			VerificationToken: "expiredtoken",
			TokenExpiresAt:    &expired,
		}
		repo := &stubUserRepo{findByTokenResult: user}
		verifier := usecase.NewEmailVerifier(repo, clk)

		err := verifier.Verify(context.Background(), "expiredtoken")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected AppError, got %T", err)
		}
		if appErr.ErrCode != apperror.ExpiredToken {
			t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.ExpiredToken)
		}
	})

	t.Run("returns INVALID_TOKEN for empty token", func(t *testing.T) {
		repo := &stubUserRepo{}
		verifier := usecase.NewEmailVerifier(repo, clk)

		err := verifier.Verify(context.Background(), "")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			t.Fatalf("expected AppError, got %T", err)
		}
		if appErr.ErrCode != apperror.InvalidToken {
			t.Errorf("ErrCode = %q, want %q", appErr.ErrCode, apperror.InvalidToken)
		}
	})
}
