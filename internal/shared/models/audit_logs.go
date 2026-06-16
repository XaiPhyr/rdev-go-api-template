package models

import "time"

type AuditLog struct {
	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	Status    string     `bun:"status,default:'A'" json:"status"`
	UUID      string     `bun:"uuid,notnull,unique,type:uuid,default:gen_random_uuid()" json:"uuid"`
	CreatedAt time.Time  `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`

	UserID         int64  `bun:"user_id" json:"user_id"`
	Path           string `bun:"path" json:"path"`
	Action         string `bun:"action" json:"action"`
	ResponseStatus int    `bun:"response_status" json:"response_status"`
	ModuleID       string `bun:"module_id" json:"module_id"`
	Module         string `bun:"module" json:"module"`
	BeforeChange   any    `bun:"before_change" json:"before_change"`
	AfterChange    any    `bun:"after_change" json:"after_change"`
	IPAddress      string `bun:"ip_address" json:"ip_address"`
	UserAgent      string `bun:"user_agent" json:"user_agent"`
	ErrorMessage   string `bun:"error_message" json:"error_message"`
}
