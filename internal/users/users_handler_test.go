package users_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/XaiPhyr/rdev-go-api-template/internal/middleware"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/testers"
	"github.com/XaiPhyr/rdev-go-api-template/internal/users"
	"github.com/gin-gonic/gin"
)

func TestUserHandlerw(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &users.MockUsersService{}
	h := users.NewUserHandler(mockSvc)

	successParseToken := func(token string) (int64, error) {
		return 42, nil
	}
	successCanAccess := func(ctx context.Context, userID int64, role string) (bool, error) {
		return true, nil
	}
	authFunc := func(authSvc *auth.MockAuthService) gin.HandlerFunc {
		return middleware.AuthRequired(authSvc)
	}
	permissionFunc := func(authSvc *auth.MockAuthService, requiredRole string) gin.HandlerFunc {
		return middleware.PermissionRequired(authSvc, requiredRole)
	}

	tests := []struct {
		name           string
		method         string
		routePath      string
		url            string
		requiredRole   string
		tokenHeader    string
		authFunc       func(authSvc *auth.MockAuthService) gin.HandlerFunc
		permissionFunc func(authSvc *auth.MockAuthService, requiredRole string) gin.HandlerFunc
		mockParseToken func(token string) (int64, error)
		mockCanAccess  func(ctx context.Context, userID int64, role string) (bool, error)
		body           any
		handlerFunc    gin.HandlerFunc
		expectedStatus int
	}{
		{
			name:           "Create - Success",
			method:         http.MethodPost,
			routePath:      "/api/v1/users",
			url:            "/api/v1/users",
			requiredRole:   "users:create",
			tokenHeader:    "Bearer valid-token",
			mockParseToken: successParseToken,
			mockCanAccess:  successCanAccess,
			body: users.UserRequest{
				FirstName: new("test first_name"),
				LastName:  new("test last_name"),
				Email:     new("test@local.com"),
				Username:  new("test"),
				Password:  new("12345678"),
			},
			handlerFunc:    h.Create,
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testers.RunGin(t, testers.GinTestConfig{
				Method:         tt.method,
				RoutePath:      tt.routePath,
				URL:            tt.url,
				RequiredRole:   tt.requiredRole,
				TokenHeader:    tt.tokenHeader,
				AuthFunc:       tt.authFunc,
				PermissionFunc: tt.permissionFunc,
				MockParseToken: tt.mockParseToken,
				MockCanAccess:  tt.mockCanAccess,
				Body:           tt.body,
				HandlerFunc:    tt.handlerFunc,
				ExpectedStatus: tt.expectedStatus,
			}, authFunc, permissionFunc)
		})
	}
}
