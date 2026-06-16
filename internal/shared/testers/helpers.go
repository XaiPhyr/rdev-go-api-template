package testers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/XaiPhyr/rdev-go-api-template/internal/auth"
	"github.com/gin-gonic/gin"
)

type AuthFunc func(authSvc *auth.MockAuthService) gin.HandlerFunc
type PermissionFunc func(authSvc *auth.MockAuthService, requiredRole string) gin.HandlerFunc

type GinTestConfig struct {
	Method         string
	RoutePath      string
	URL            string
	Body           any
	TokenHeader    string
	ExpectedStatus int
	RequiredRole   string
	HandlerFunc    gin.HandlerFunc
	MockParseToken func(token string) (int64, error)
	MockCanAccess  func(ctx context.Context, userID int64, role string) (bool, error)
	AuthFunc       AuthFunc
	PermissionFunc PermissionFunc
}

func JSONCtx(t *testing.T, method, url string, body any, tokenHeader ...string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	var jsonBytes []byte
	if strBody, ok := body.(string); ok {
		jsonBytes = []byte(strBody)
	} else if body != nil {
		var err error
		jsonBytes, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal json body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if len(tokenHeader) > 0 && tokenHeader[0] != "" {
		req.Header.Set("Authorization", tokenHeader[0])
	}

	ctx.Request = req

	return ctx, w
}

func RunGin(t *testing.T, config GinTestConfig, defaultAuth AuthFunc, defaultPerm PermissionFunc) {
	t.Helper()

	r := gin.New()

	mockAuth := &auth.MockAuthService{
		ParseTokenFunc: config.MockParseToken,
		CanAccessFunc:  config.MockCanAccess,
	}

	finalAuth := defaultAuth
	if config.AuthFunc != nil {
		finalAuth = config.AuthFunc
	}

	finalPerm := defaultPerm
	if config.PermissionFunc != nil {
		finalPerm = config.PermissionFunc
	}

	r.Handle(
		config.Method,
		config.RoutePath,
		finalAuth(mockAuth),
		finalPerm(mockAuth, config.RequiredRole),
		config.HandlerFunc,
	)

	ctx, w := JSONCtx(t, config.Method, config.URL, config.Body, config.TokenHeader)

	r.ServeHTTP(w, ctx.Request)

	if w.Code != config.ExpectedStatus {
		t.Errorf("failed: expected status %d, got %d. Body response: %s", config.ExpectedStatus, w.Code, w.Body.String())
	}
}
