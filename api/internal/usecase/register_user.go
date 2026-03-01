package usecase

import (
	"context"
	"fmt"

	"github.com/sibukixxx/travelist/api/internal/apperror"
	"github.com/sibukixxx/travelist/api/internal/domain"
	"github.com/sibukixxx/travelist/api/internal/infra/clock"
)

// UserRepo is the repository interface used by user usecases.
type UserRepo interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

// EmailSender is the interface for sending emails.
type EmailSender interface {
	SendVerificationEmail(ctx context.Context, to, token string) error
}

// RegisterRequest is the input for user registration.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResult is the output of user registration.
type RegisterResult struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// UserRegistrar handles user registration.
type UserRegistrar struct {
	users UserRepo
	email EmailSender
	clock clock.Clock
}

// NewUserRegistrar creates a new UserRegistrar.
func NewUserRegistrar(users UserRepo, email EmailSender, clk clock.Clock) *UserRegistrar {
	return &UserRegistrar{users: users, email: email, clock: clk}
}

// Register creates a new user with email verification.
func (r *UserRegistrar) Register(ctx context.Context, req RegisterRequest) (*RegisterResult, error) {
	// Validate
	if err := domain.ValidateEmail(req.Email); err != nil {
		return nil, apperror.NewValidation(err.Error())
	}
	if err := domain.ValidatePassword(req.Password); err != nil {
		return nil, apperror.NewValidation(err.Error())
	}

	// Check duplicate
	existing, err := r.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Errorf("find by email: %w", err))
	}
	if existing != nil {
		return nil, apperror.NewConflict("email already registered")
	}

	// Hash password
	hash, err := domain.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.NewInternal(fmt.Errorf("hash password: %w", err))
	}

	// Generate verification token
	token, err := domain.GenerateVerificationToken()
	if err != nil {
		return nil, apperror.NewInternal(fmt.Errorf("generate token: %w", err))
	}

	now := r.clock.Now()
	expires := now.Add(24 * 60 * 60 * 1e9) // 24 hours

	user := &domain.User{
		ID:                domain.NewUserID(),
		Email:             req.Email,
		PasswordHash:      hash,
		VerificationToken: token,
		TokenExpiresAt:    &expires,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := r.users.Create(ctx, user); err != nil {
		return nil, err
	}

	if err := r.email.SendVerificationEmail(ctx, user.Email, user.VerificationToken); err != nil {
		return nil, apperror.NewInternal(fmt.Errorf("send verification email: %w", err))
	}

	return &RegisterResult{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}
