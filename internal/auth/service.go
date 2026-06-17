package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/XaiPhyr/rdev-go-api-template/internal/config"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/email"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/helpers"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetUsernameOrEmail(ctx context.Context, username string) (*models.User, error)
	Register(ctx context.Context, user *models.User) error
	CheckUserPermission(ctx context.Context, userID int64, roleName string) ([]string, error)
}

type AuthService interface {
	GenerateToken(userID int64) (string, error)
	ParseToken(token string) (int64, error)
	CanAccess(ctx context.Context, userID int64, requiredRole string) (bool, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
	Register(ctx context.Context, req RegisterRequest) error
}

type service struct {
	r     AuthRepository
	c     *config.Config
	es    email.EmailService
	redis *redis.Client
}

func NewAuthService(r AuthRepository, c *config.Config, es email.EmailService, redis *redis.Client) *service {
	return &service{r: r, c: c, es: es, redis: redis}
}

func (s *service) GenerateToken(userID int64) (string, error) {
	jwtKey := []byte(s.c.JWTSecretKey)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (s *service) ParseToken(token string) (int64, error) {
	jwtKey := []byte(s.c.JWTSecretKey)

	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !t.Valid {
		return 0, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			return int64(userID), nil
		}
	}

	return 0, jwt.ErrTokenInvalidClaims
}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, error) {
	err := helpers.ValidateStruct(req)
	if err != nil {
		return "", fmt.Errorf("field missing %v", err)
	}

	username := helpers.CleanSpecialChars(strings.TrimSpace(req.Username))
	user, err := s.r.GetUsernameOrEmail(ctx, username)
	if err != nil {
		return "", fmt.Errorf("no user found, %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("incorrect password, %v", err)
	}

	return s.GenerateToken(user.ID)
}

func (s *service) Register(ctx context.Context, req RegisterRequest) error {
	err := helpers.ValidateStruct(req)
	if err != nil {
		return fmt.Errorf("field missing %v", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hash %v", err)
	}

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(passwordHash),
	}

	err = s.r.Register(ctx, user)
	if err != nil {
		return errors.New("unable to register")
	}

	go func(email string) {
		if err := s.es.SendEmail(email); err != nil {
			log.Printf("Failed to send email: %v", err)
		}
	}(req.Email)

	return nil
}

func (s *service) CanAccess(ctx context.Context, userID int64, requiredRole string) (bool, error) {
	cacheKey := fmt.Sprintf("user:perms:%d", userID)

	existCount, err := s.redis.Exists(ctx, cacheKey).Result()
	if err == nil && existCount > 0 {
		isSuperAdmin, _ := s.redis.SIsMember(ctx, cacheKey, "super_admin").Result()
		if isSuperAdmin {
			return true, nil
		}

		hasRole, _ := s.redis.SIsMember(ctx, cacheKey, requiredRole).Result()
		return hasRole, nil
	}

	allPerms, err := s.r.CheckUserPermission(ctx, userID, requiredRole)
	if err != nil {
		log.Println(fmt.Errorf("user permission error: %w", err))
		return false, err
	}

	if len(allPerms) > 0 {
		pipe := s.redis.Pipeline()
		pipe.SAdd(ctx, cacheKey, allPerms)
		pipe.Expire(ctx, cacheKey, 1*time.Hour)
		_, err := pipe.Exec(ctx)
		if err != nil {
			log.Printf("failed to update redis: %v", err)
		}
	} else {
		s.redis.SAdd(ctx, cacheKey, "NONE")
		s.redis.Expire(ctx, cacheKey, 5*time.Minute)
	}

	return slices.Contains(allPerms, requiredRole) || slices.Contains(allPerms, "super_admin"), nil
}
