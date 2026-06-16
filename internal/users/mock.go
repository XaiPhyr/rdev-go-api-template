package users

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
)

type MockUsersService struct {
	// code here!
}

func (m *MockUsersService) Create(ctx context.Context, req UserRequest) error {
	return nil
}

func (m *MockUsersService) ReadOne(ctx context.Context, uuid string) (*models.User, error) {
	return nil, nil
}

func (m *MockUsersService) ReadAll(ctx context.Context, req dto.BaseFilters) ([]models.User, int, error) {
	return nil, 0, nil
}

func (m *MockUsersService) Update(ctx context.Context, uuid string, req UserRequest) error {
	return nil
}

func (m *MockUsersService) Delete(ctx context.Context, uuid string) error {
	return nil
}

type MockUsersRepository struct {
	CreateFunc  func(ctx context.Context, user *models.User) error
	ReadOneFunc func(ctx context.Context, uuid string) (*models.User, error)
	ReadAllFunc func(ctx context.Context, q dto.Query) ([]models.User, int, error)
	UpdateFunc  func(ctx context.Context, user *models.User) error
	DeleteFunc  func(ctx context.Context, uuid string) error
}

func (m *MockUsersRepository) Create(ctx context.Context, user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}

	return nil
}

func (m *MockUsersRepository) ReadOne(ctx context.Context, uuid string) (*models.User, error) {
	if m.ReadOneFunc != nil {
		return m.ReadOneFunc(ctx, uuid)
	}

	return nil, nil
}

func (m *MockUsersRepository) ReadAll(ctx context.Context, q dto.Query) ([]models.User, int, error) {
	if m.ReadAllFunc != nil {
		return m.ReadAllFunc(ctx, q)
	}

	return nil, 0, nil
}

func (m *MockUsersRepository) Update(ctx context.Context, user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}

	return nil
}

func (m *MockUsersRepository) Delete(ctx context.Context, uuid string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, uuid)
	}

	return nil
}
