package audit_logs

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
)

type MockAuditLogService struct {
	// code here!
}

func (m *MockAuditLogService) QueAuditLog(ctx context.Context) {
	// code here!
}

func (m *MockAuditLogService) AuditLog(audit AuditLogRequest, module_id, module string, beforeChange, afterChange any, err error) {
	// code here!
}

type MockAuditLogRepository struct {
	CreateFunc func(ctx context.Context, audit_log models.AuditLog) error
}

func (m *MockAuditLogRepository) Create(ctx context.Context, audit_log models.AuditLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, audit_log)
	}

	return nil
}
