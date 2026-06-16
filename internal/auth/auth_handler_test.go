package auth_test

import (
	"net/http"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/testers"
	"github.com/gin-gonic/gin"
)

var h = auth.NewAuthHandler(&auth.MockAuthService{})

func TestAuthHandlerLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		url            string
		body           any
		handlerFunc    gin.HandlerFunc
		expectedStatus int
	}{
		{
			name: "complete login credentials",
			url:  "/login",
			body: auth.LoginRequest{
				Username: "rdev",
				Password: "!Abc1234",
			},
			handlerFunc:    h.Login,
			expectedStatus: http.StatusOK,
		},
		{
			name: "login no password",
			url:  "/login",
			body: auth.LoginRequest{
				Username: "rdev",
				Password: "",
			},
			handlerFunc:    h.Login,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "login password less than minimum required",
			url:  "/login",
			body: auth.LoginRequest{
				Username: "rdev",
				Password: "1234",
			},
			handlerFunc:    h.Login,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "correct json form",
			url:  "/register",
			body: auth.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "test@local.co",
				Username:  "jdoe",
				Password:  "12341234",
			},
			handlerFunc:    h.Register,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "empty data",
			url:  "/register",
			body: auth.RegisterRequest{
				FirstName: "",
				LastName:  "",
				Email:     "",
				Username:  "",
				Password:  "",
			},
			handlerFunc:    h.Register,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "minimum password required",
			url:  "/register",
			body: auth.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "rdev",
				Username:  "rdev",
				Password:  "1234",
			},
			handlerFunc:    h.Register,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, w := testers.JSONCtx(t, http.MethodPost, tt.url, tt.body)

			tt.handlerFunc(ctx)

			if w.Code != tt.expectedStatus {
				t.Errorf("Login() status = %d, resp = %v, want %d", w.Code, w.Body, tt.expectedStatus)
			}
		})
	}
}
