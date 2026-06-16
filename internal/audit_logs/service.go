package audit_logs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
)

type AuditLogRepository interface {
	Create(ctx context.Context, audit_log models.AuditLog) error
}

type AuditLogService interface {
	QueAuditLog(ctx context.Context)
	AuditLog(audit AuditLogRequest, module_id, module string, beforeChange, afterChange any, err error)
}

type service struct {
	r        AuditLogRepository
	auditQue chan models.AuditLog
}

func NewAuditLogService(r AuditLogRepository) AuditLogService {
	return &service{r: r, auditQue: make(chan models.AuditLog, 1000)}
}

func (s *service) QueAuditLog(ctx context.Context) {
	for {
		select {
		case log := <-s.auditQue:
			fmt.Println("QUE START...")
			if err := s.r.Create(ctx, log); err != nil {
				s.handleDeadLetter(log, err)
			}
		case <-ctx.Done():
			fmt.Println("QUE DONE...")
			return
		}
	}
}

func (s *service) handleDeadLetter(al models.AuditLog, err error) {
	_ = al
	fmt.Printf("⚠️ DEAD LETTER TRIGGERED: %v\n", err)
}

func (s *service) AuditLog(audit AuditLogRequest, module_id, module string, beforeChange, afterChange any, err error) {
	errMsg := ""
	status := http.StatusOK
	if err != nil {
		status = http.StatusBadRequest
		errMsg = err.Error()
	}

	al := models.AuditLog{
		UserID:         audit.UserID,
		Path:           audit.Path,
		Action:         audit.Action,
		ResponseStatus: status,
		ModuleID:       module_id,
		Module:         module,
		BeforeChange:   beforeChange,
		AfterChange:    afterChange,
		IPAddress:      audit.IPAddress,
		UserAgent:      audit.UserAgent,
		ErrorMessage:   errMsg,
	}

	s.auditQue <- al
}
