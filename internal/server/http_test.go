package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acb/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer_HealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	httpSrv := NewHTTPServer(
		"8080",
		nil, // registrySvc
		nil, // contextMgr
		auth.NewJWTManager("test-secret"),
		auth.NewRBAC(),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	httpSrv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestHTTPServer_Metrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	httpSrv := NewHTTPServer(
		"8080",
		nil, // registrySvc
		nil, // contextMgr
		auth.NewJWTManager("test-secret"),
		auth.NewRBAC(),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/metrics", nil)
	httpSrv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHTTPServer_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	httpSrv := NewHTTPServer(
		"8080",
		nil,
		nil,
		auth.NewJWTManager("test-secret"),
		auth.NewRBAC(),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", nil)
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)

	// Should fail without body
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHTTPServer_AuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	jwtManager := auth.NewJWTManager("test-secret")
	httpSrv := NewHTTPServer(
		"8080",
		nil,
		nil,
		jwtManager,
		auth.NewRBAC(),
	)

	// Generate token
	token, err := jwtManager.GenerateAccessToken("agent-1", "default", []string{"agent-full"})
	require.NoError(t, err)

	// Test protected endpoint with token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	httpSrv.router.ServeHTTP(w, req)

	// Should not return 401 (may return other status based on handler)
	assert.NotEqual(t, http.StatusUnauthorized, w.Code)
}

func TestHTTPServer_AuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	httpSrv := NewHTTPServer(
		"8080",
		nil,
		nil,
		auth.NewJWTManager("test-secret"),
		auth.NewRBAC(),
	)

	// Test protected endpoint without token
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/agents", nil)
	httpSrv.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCorsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(corsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
}

func TestCorsPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(corsMiddleware())
	router.OPTIONS("/test", func(c *gin.Context) { c.String(200, "ok") })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(requestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

func TestRequestIDMiddleware_Passthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(requestIDMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "abc123")
	router.ServeHTTP(w, req)

	assert.Equal(t, "abc123", w.Header().Get("X-Request-ID"))
}
