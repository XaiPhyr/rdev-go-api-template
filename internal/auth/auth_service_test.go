package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/XaiPhyr/rdev-go-api-template/internal/config"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
	"golang.org/x/crypto/bcrypt"
)

var testAuthRepo = &auth.MockAuthRepository{}
var testAuthSvc = NewTestAuth(testAuthRepo)

func NewTestAuth(testAuthRepo *auth.MockAuthRepository) auth.AuthService {
	config := &config.Config{JWTSecretKey: "rdev-go-api-template_jwt_key"}

	return auth.NewAuthService(testAuthRepo, config, nil, nil)
}

func TestAuthServiceLogin(t *testing.T) {
	t.Run("login credentials", func(t *testing.T) {
		testAuthRepo.GetUsernameOrEmailFunc = func(ctx context.Context, username string) (*models.User, error) {
			passwordHash, _ := bcrypt.GenerateFromPassword([]byte("!Abc1234"), bcrypt.DefaultCost)
			user := &models.User{
				ID:       1,
				Username: "test_account",
				Email:    "test_account@local.com",
				Password: string(passwordHash),
			}

			return user, nil
		}

		req := auth.LoginRequest{
			Username: "rdev",
			Password: "!Abc1234",
		}

		_, err := testAuthSvc.Login(context.Background(), req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestAuthServiceRegister(t *testing.T) {
	t.Run("user registration", func(t *testing.T) {
		testAuthRepo.RegisterFunc = func(ctx context.Context, user *models.User) error {
			return nil
		}

		req := auth.RegisterRequest{
			FirstName: "John",
			LastName:  "Nuñez",
			Email:     "rdev@test.com",
			Username:  "rdev",
			Password:  "!Abc1234",
		}

		err := testAuthSvc.Register(context.Background(), req)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

func TestAuthServiceLoginRateLimit(t *testing.T) {
	t.Run("login rate limit", func(t *testing.T) {
		callCount := 0
		testAuthRepo.GetUsernameOrEmailFunc = func(ctx context.Context, username string) (*models.User, error) {
			callCount++

			passwordHash, _ := bcrypt.GenerateFromPassword([]byte("!Abc1234"), bcrypt.DefaultCost)
			user := &models.User{
				ID:       1,
				Username: "test_account",
				Email:    "test_account@local.com",
				Password: string(passwordHash),
			}

			if callCount > 4 {
				return nil, errors.New("too many requests.")
			}

			return user, nil
		}

		var err error
		for range 3 {
			req := auth.LoginRequest{
				Username: "rdev",
				Password: "!Abc1234",
			}

			_, err = testAuthSvc.Login(context.Background(), req)
		}

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
