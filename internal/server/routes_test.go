package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/XaiPhyr/rdev-go-api-template/internal/middleware"
	"github.com/gin-gonic/gin"
)

func TestRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	mockAuth := &auth.MockAuthService{
		ParseTokenFunc: func(token string) (int64, error) {
			return 42, nil
		},
		CanAccessFunc: func(ctx context.Context, userID int64, role string) (bool, error) {
			return true, nil
		},
	}

	protectedGroup := r.Group("/api/v1")
	protectedGroup.Use(middleware.AuthRequired(mockAuth))

	protectedGroup.GET(
		"/service_types",
		middleware.PermissionRequired(mockAuth, "service_types:view"),
		middleware.RateLimiter(),
		func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		},
	)

	for i := range 15 {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/service_types", nil)
		req.Header.Set("Content-Type", "application/json")

		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if i+1 > 4 {
			if w.Code != http.StatusTooManyRequests {
				t.Errorf("Request %d: Expected status 429 (Too Many Requests), got %d", i+1, w.Code)
			}
		} else {
			if w.Code != http.StatusOK {
				t.Errorf("Request %d: Expected status 200 (OK) before limit hit, got %d", i+1, w.Code)
			}
		}
	}
}
