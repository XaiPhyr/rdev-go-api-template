package email_test

import (
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/email"
)

func TestEmail(t *testing.T) {
	mockSvc := email.NewEmailService("localhost", "1234", "from@local.com")

	t.Run("test email", func(t *testing.T) {
		err := mockSvc.SendEmail("to@local.com")

		t.Errorf("expected no error, got %v", err)
	})
}
