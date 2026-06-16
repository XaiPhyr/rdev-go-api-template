package audit_logs_test

import (
	"context"
	"testing"
	"time"

	"github.com/XaiPhyr/rdev-go-api-template/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
)

func TestAuditLog(t *testing.T) {
	t.Run("auditlog func successfully pushes to channel and saves via worker", func(t *testing.T) {
		ctx := t.Context()

		repo := &audit_logs.MockAuditLogRepository{}
		svc := audit_logs.NewAuditLogService(repo)

		done := make(chan struct{})

		repo.CreateFunc = func(ctx context.Context, audit_log models.AuditLog) error {
			if audit_log.Module != "Auth" {
				t.Errorf("expected module to be 'Auth', got %s", audit_log.Module)
			}

			close(done)
			return nil
		}

		go svc.QueAuditLog(ctx)

		req := audit_logs.AuditLogRequest{
			UserID:    1,
			Path:      "/login",
			Action:    "POST",
			IPAddress: "127.0.0.1",
			UserAgent: "Mozilla",
		}
		svc.AuditLog(req, "mod_01", "Auth", nil, nil, nil)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("test timed out: QueAuditLog did not process the log within 2 seconds")
		}
	})
}
