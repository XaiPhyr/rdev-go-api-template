package users

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
)

type MockUserService struct {
	// code here!
}

func (m *MockUserService) Create(ctx context.Context, req UserRequest) error {
	return nil
}

func (m *MockUserService) ReadOne(ctx context.Context, uuid string) (*models.User, error) {
	return nil, nil
}

func (m *MockUserService) ReadAll(ctx context.Context, req dto.BaseFilters) ([]models.User, int, error) {
	return nil, 0, nil
}

func (m *MockUserService) Update(ctx context.Context, uuid string, req UserRequest) error {
	return nil
}

func (m *MockUserService) Delete(ctx context.Context, uuid string) error {
	return nil
}

type MockUserRepository struct {
	CreateFunc  func(ctx context.Context, user *models.User) error
	ReadOneFunc func(ctx context.Context, uuid string) (*models.User, error)
	ReadAllFunc func(ctx context.Context, q dto.Query) ([]models.User, int, error)
	UpdateFunc  func(ctx context.Context, user *models.User) error
	DeleteFunc  func(ctx context.Context, uuid string) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}

	return nil
}

func (m *MockUserRepository) ReadOne(ctx context.Context, uuid string) (*models.User, error) {
	if m.ReadOneFunc != nil {
		return m.ReadOneFunc(ctx, uuid)
	}

	return nil, nil
}

func (m *MockUserRepository) ReadAll(ctx context.Context, q dto.Query) ([]models.User, int, error) {
	if m.ReadAllFunc != nil {
		return m.ReadAllFunc(ctx, q)
	}

	return nil, 0, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}

	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, uuid string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uuid)
	}

	return nil
}
