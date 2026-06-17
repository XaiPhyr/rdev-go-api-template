package email_test

import (
	"errors"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/email"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		svc     *email.MockEmailService
		wantErr bool
	}{
		{
			name: "successful email send",
			svc: &email.MockEmailService{
				SendEmailFunc: func(toEmail string) error {
					if toEmail != "to@local.com" {
						t.Errorf("expected to@local.com, got %s", toEmail)
					}
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "failed email send",
			svc: &email.MockEmailService{
				SendEmailFunc: func(toEmail string) error {
					return errors.New("smtp connection timeout")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.svc.SendEmail("to@local.com")

			if (err != nil) != tt.wantErr {
				t.Errorf("SendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
