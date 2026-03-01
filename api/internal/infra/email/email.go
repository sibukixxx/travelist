package email

import (
	"context"
	"log"
)

// EmailSender is the interface for sending emails.
type EmailSender interface {
	SendVerificationEmail(ctx context.Context, to, token string) error
}

// LogEmailSender logs verification emails instead of sending them.
type LogEmailSender struct{}

func (s *LogEmailSender) SendVerificationEmail(_ context.Context, to, token string) error {
	log.Printf("[EMAIL] Verification email to=%s token=%s url=http://localhost:8080/api/users/verify?token=%s", to, token, token)
	return nil
}
