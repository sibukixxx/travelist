package email_test

import (
	"context"
	"testing"

	"github.com/sibukixxx/travelist/api/internal/infra/email"
)

func TestLogEmailSenderSendVerificationEmail(t *testing.T) {
	sender := &email.LogEmailSender{}

	err := sender.SendVerificationEmail(context.Background(), "test@example.com", "abc123token")
	if err != nil {
		t.Errorf("SendVerificationEmail() error = %v, want nil", err)
	}
}
