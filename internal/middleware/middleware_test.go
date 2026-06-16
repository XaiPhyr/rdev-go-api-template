package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/XaiPhyr/rdev-go-api-template/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestAuthAndPermissionRequired(t *testing.T) {
	tests := []struct {
		name           string
		authSvc        auth.AuthService
		setupContext   func(req *http.Request, r *gin.Engine)
		expectedStatus int
	}{
		{
			name:           "Missing Authorization Header",
			authSvc:        &auth.MockAuthService{},
			setupContext:   func(req *http.Request, r *gin.Engine) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:    "Empty Token in Authorization Header",
			authSvc: &auth.MockAuthService{},
			setupContext: func(req *http.Request, r *gin.Engine) {
				req.Header.Set("Authorization", "Bearer ")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:    "Authenticated But Lacks Permission",
			authSvc: &auth.MockAuthService{},
			setupContext: func(req *http.Request, r *gin.Engine) {
				req.Header.Set("Authorization", "Bearer valid-token-without-permission")
				r.Use(func(c *gin.Context) {
					c.Next()
				})
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "Authenticated With Correct Permission",
			authSvc: &auth.MockAuthService{
				ParseTokenFunc: func(token string) (int64, error) {
					return 123, nil
				},
				CanAccessFunc: func(ctx context.Context, userID int64, role string) (bool, error) {
					if role == "users:create" && userID == 123 {
						return true, nil
					}

					return false, nil
				},
			},
			setupContext: func(req *http.Request, r *gin.Engine) {
				req.Header.Set("Authorization", "Bearer valid-token-with-permission")
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			r := gin.New()

			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			tt.setupContext(req, r)

			r.GET("/",
				middleware.AuthRequired(tt.authSvc),
				middleware.PermissionRequired(tt.authSvc, "users:create"),
				func(c *gin.Context) {
					c.Status(http.StatusOK)
				},
			)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: status = %d, want %d", tt.name, w.Code, tt.expectedStatus)
			}
		})
	}
}
