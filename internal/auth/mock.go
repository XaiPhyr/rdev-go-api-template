package auth

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
)

type MockAuthService struct {
	ParseTokenFunc func(token string) (int64, error)
	CanAccessFunc  func(ctx context.Context, userID int64, role string) (bool, error)
}

func (m *MockAuthService) GenerateToken(userID int64) (string, error) {
	return "", nil
}

func (m *MockAuthService) ParseToken(token string) (int64, error) {
	if m.ParseTokenFunc != nil {
		return m.ParseTokenFunc(token)
	}

	return 0, nil
}

func (m *MockAuthService) CanAccess(ctx context.Context, userID int64, requiredRole string) (bool, error) {
	if m.CanAccessFunc != nil {
		return m.CanAccessFunc(ctx, userID, requiredRole)
	}

	return false, nil
}

func (m *MockAuthService) Login(ctx context.Context, req LoginRequest) (string, error) {
	return "", nil
}

func (m *MockAuthService) Register(ctx context.Context, req RegisterRequest) error {
	return nil
}

type MockAuthRepository struct {
	GetUsernameOrEmailFunc  func(ctx context.Context, username string) (*models.User, error)
	RegisterFunc            func(ctx context.Context, user *models.User) error
	CheckUserPermissionFunc func(ctx context.Context, userID int64, roleName string) ([]string, error)
}

func (m *MockAuthRepository) GetUsernameOrEmail(ctx context.Context, username string) (*models.User, error) {
	if m.GetUsernameOrEmailFunc != nil {
		return m.GetUsernameOrEmailFunc(ctx, username)
	}

	return &models.User{}, nil
}

func (m *MockAuthRepository) Register(ctx context.Context, user *models.User) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, user)
	}

	return nil
}

func (m *MockAuthRepository) CheckUserPermission(ctx context.Context, userID int64, requiredRole string) ([]string, error) {
	if m.CheckUserPermissionFunc != nil {
		return m.CheckUserPermissionFunc(ctx, userID, requiredRole)
	}

	return nil, nil
}
