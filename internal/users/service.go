package users

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/audit_logs"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	ReadOne(ctx context.Context, uuid string) (*models.User, error)
	ReadAll(ctx context.Context, q dto.Query) ([]models.User, int, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, uuid string) error
}

type UserService interface {
	Create(ctx context.Context, req UserRequest) error
	ReadOne(ctx context.Context, uuid string) (*models.User, error)
	ReadAll(ctx context.Context, req dto.BaseFilters) ([]models.User, int, error)
	Update(ctx context.Context, uuid string, req UserRequest) error
	Delete(ctx context.Context, uuid string) error
}

type service struct {
	r        UserRepository
	es       email.EmailService
	redis    *redis.Client
	auditLog audit_logs.AuditLogService
}

func NewUserService(r UserRepository, es email.EmailService, redis *redis.Client, auditLog audit_logs.AuditLogService) *service {
	return &service{r: r, es: es, redis: redis, auditLog: auditLog}
}

func (s *service) Create(ctx context.Context, req UserRequest) error {
	// hash password req.Password =

	err := s.r.Create(ctx, req.ToModel(&models.User{}))

	return err
}

func (s *service) ReadOne(ctx context.Context, uuid string) (*models.User, error) {
	user, err := s.r.ReadOne(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) ReadAll(ctx context.Context, req dto.BaseFilters) ([]models.User, int, error) {
	filters := req.SanitizeQuery([]string{"name", "price"})

	users, total, err := s.r.ReadAll(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *service) Update(ctx context.Context, uuid string, req UserRequest) error {
	user, err := s.r.ReadOne(ctx, uuid)
	if err != nil {
		return err
	}

	if err := s.r.Update(ctx, req.ToModel(user)); err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	err := s.r.Delete(ctx, uuid)

	return err
}
