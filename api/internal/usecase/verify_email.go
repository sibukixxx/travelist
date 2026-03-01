package usecase

import (
	"context"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
)

// EmailVerifier handles email verification.
type EmailVerifier struct {
	users UserRepo
	clock clock.Clock
}

// NewEmailVerifier creates a new EmailVerifier.
func NewEmailVerifier(users UserRepo, clk clock.Clock) *EmailVerifier {
	return &EmailVerifier{users: users, clock: clk}
}

// Verify verifies a user's email using the provided token.
func (v *EmailVerifier) Verify(ctx context.Context, token string) error {
	if token == "" {
		return apperror.NewBadRequestWithCode(apperror.InvalidToken, "verification token is required")
	}

	user, err := v.users.FindByVerificationToken(ctx, token)
	if err != nil {
		return apperror.NewInternal(err)
	}
	if user == nil {
		return apperror.NewBadRequestWithCode(apperror.InvalidToken, "invalid verification token")
	}

	// Check expiry
	now := v.clock.Now()
	if user.TokenExpiresAt != nil && now.After(*user.TokenExpiresAt) {
		return apperror.NewBadRequestWithCode(apperror.ExpiredToken, "verification token has expired")
	}

	// Mark as verified
	user.EmailVerifiedAt = &now
	user.VerificationToken = ""
	user.TokenExpiresAt = nil
	user.UpdatedAt = now

	if err := v.users.Update(ctx, user); err != nil {
		return apperror.NewInternal(err)
	}

	return nil
}
