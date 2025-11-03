package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/p", AuthMiddleware(NewJWTManager("s")), func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/p", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/p", AuthMiddleware(NewJWTManager("s")), func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/p", nil)
	req.Header.Set("Authorization", "Token abc")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwt := NewJWTManager("s")
	token, err := jwt.GenerateAccessToken("a", "default", []string{"agent-full"})
	if err != nil {
		t.Fatalf("token error: %v", err)
	}

	r := gin.New()
	r.GET("/p", AuthMiddleware(jwt), func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/p", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwt := NewJWTManager("s")
	token, err := jwt.GenerateAccessToken("a", "default", []string{"observer"})
	if err != nil {
		t.Fatalf("token error: %v", err)
	}

	r := gin.New()
	r.GET("/p",
		AuthMiddleware(jwt),
		RequireRole(Role("observer")),
		func(c *gin.Context) { c.String(200, "ok") },
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/p", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}


