package auth

import (
	"context"
	"errors"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
	"github.com/uptrace/bun"
)

type repository struct {
	db *bun.DB
}

func NewAuthRepository(db *bun.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetUsernameOrEmail(ctx context.Context, username string) (*models.User, error) {
	return nil, errors.New("user not found")
}

func (r *repository) Register(ctx context.Context, user *models.User) error {
	return nil
}

func (r *repository) CheckUserPermission(ctx context.Context, userID int64, roleName string) ([]string, error) {
	return nil, nil
}
