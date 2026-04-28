package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"ims-intro/pkg/domain"
)

func buildTestToken(t *testing.T, key string, userID int64, role string, expiresAt time.Time) string {
	t.Helper()

	claims := &domain.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	return tokenString
}

func TestAuthMiddleware_MissingToken_ShouldFail(t *testing.T) {
	t.Setenv("JWT_KEY", "unit-test-secret")
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	called := false

	err := AuthMiddleware(func(c echo.Context) error {
		called = true
		return nil
	})(ctx)

	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if called {
		t.Fatalf("expected next handler not to be called")
	}
}

func TestAuthMiddleware_ValidToken_ShouldAllow(t *testing.T) {
	key := "unit-test-secret"
	t.Setenv("JWT_KEY", key)
	e := echo.New()
	token := buildTestToken(t, key, 41, "Admin", time.Now().Add(1*time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	called := false

	err := AuthMiddleware(func(c echo.Context) error {
		called = true
		if got := c.Get("user_id"); got != int64(41) {
			t.Fatalf("expected user_id 41, got %v", got)
		}
		if got := c.Get("role"); got != "admin" {
			t.Fatalf("expected role admin, got %v", got)
		}
		return nil
	})(ctx)

	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if !called {
		t.Fatalf("expected next handler to be called")
	}
}

func TestAuthMiddleware_ExpiredToken_ShouldFail(t *testing.T) {
	key := "unit-test-secret"
	t.Setenv("JWT_KEY", key)
	e := echo.New()
	token := buildTestToken(t, key, 41, "user", time.Now().Add(-1*time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	called := false

	err := AuthMiddleware(func(c echo.Context) error {
		called = true
		return nil
	})(ctx)

	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if called {
		t.Fatalf("expected next handler not to be called for expired token")
	}
}

func TestAuthMiddleware_InvalidSignature_ShouldFail(t *testing.T) {
	t.Setenv("JWT_KEY", "unit-test-secret")
	e := echo.New()
	token := buildTestToken(t, "wrong-secret", 41, "user", time.Now().Add(1*time.Hour))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	called := false

	err := AuthMiddleware(func(c echo.Context) error {
		called = true
		return nil
	})(ctx)

	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if called {
		t.Fatalf("expected next handler not to be called for invalid token")
	}
}

func TestAuthorize_AllowedRole_ShouldPass(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("role", "admin")
	called := false

	err := Authorize([]string{"admin"})(func(c echo.Context) error {
		called = true
		return nil
	})(ctx)

	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if !called {
		t.Fatalf("expected next handler to be called for allowed role")
	}
}

func TestAuthorize_ForbiddenRole_ShouldFail(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.Set("role", "user")
	called := false

	err := Authorize([]string{"admin"})(func(c echo.Context) error {
		called = true
		return nil
	})(ctx)

	if err != nil {
		t.Fatalf("expected no handler error, got %v", err)
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
	if called {
		t.Fatalf("expected next handler not to be called for forbidden role")
	}
}
