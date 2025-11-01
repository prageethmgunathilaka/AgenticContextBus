package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/acb/internal/auth"
	"github.com/acb/internal/context"
	"github.com/acb/internal/registry"
	"github.com/gin-gonic/gin"
)

// HTTPServer wraps HTTP server
type HTTPServer struct {
	router         *gin.Engine
	registrySvc    *registry.Service
	contextMgr     *context.Manager
	jwtManager     *auth.JWTManager
	rbac           *auth.RBAC
	port           string
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(
	port string,
	registrySvc *registry.Service,
	contextMgr *context.Manager,
	jwtManager *auth.JWTManager,
	rbac *auth.RBAC,
) *HTTPServer {
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(requestIDMiddleware())

	srv := &HTTPServer{
		router:      router,
		registrySvc:  registrySvc,
		contextMgr:   contextMgr,
		jwtManager:   jwtManager,
		rbac:         rbac,
		port:         port,
	}

	srv.setupRoutes()
	return srv
}

// setupRoutes registers all routes
func (s *HTTPServer) setupRoutes() {
	// Health check (public)
	s.router.GET("/health", s.healthCheck)
	s.router.GET("/metrics", s.metrics)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Authentication routes (public)
		v1.POST("/auth/login", s.login)
		v1.POST("/auth/refresh", s.refreshToken)

		// Protected routes
		protected := v1.Group("")
		protected.Use(auth.AuthMiddleware(s.jwtManager))
		{
			// Agent routes
			agents := protected.Group("/agents")
			{
				agents.POST("", s.registerAgent)
				agents.GET("", s.listAgents)
				agents.GET("/:agent_id", s.getAgent)
				agents.DELETE("/:agent_id", s.unregisterAgent)
				agents.POST("/:agent_id/heartbeat", s.sendHeartbeat)
			}

			// Context routes
			contexts := protected.Group("/contexts")
			{
				contexts.POST("", s.createContext)
				contexts.GET("", s.listContexts)
				contexts.GET("/:context_id", s.getContext)
				contexts.PUT("/:context_id", s.updateContext)
				contexts.DELETE("/:context_id", s.deleteContext)
			}
		}
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.router,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("HTTP server starting on port %s\n", s.port)
	return srv.ListenAndServe()
}

// healthCheck handles health check requests
func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
		"time":    time.Now(),
	})
}

// metrics handles metrics requests
func (s *HTTPServer) metrics(c *gin.Context) {
	c.String(http.StatusOK, "# Metrics placeholder\n")
}

// Middleware functions
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// login and refreshToken handlers (simple implementations for MVP)
func (s *HTTPServer) login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// MVP: Simple authentication (in production, validate against database)
	// For MVP, generate token for any username
	accessToken, err := s.jwtManager.GenerateAccessToken(req.Username, "default", []string{"agent-full"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(req.Username, "default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    3600,
		"token_type":    "Bearer",
	})
}

func (s *HTTPServer) refreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := s.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(claims.AgentID, claims.TenantID, claims.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   3600,
	})
}

