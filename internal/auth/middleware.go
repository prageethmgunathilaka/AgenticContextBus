package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates HTTP middleware for JWT authentication
func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			if err == ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("agent_id", claims.AgentID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)

		c.Next()
	}
}

// RBACMiddleware creates HTTP middleware for RBAC authorization
func RBACMiddleware(rbac *RBAC, permission Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "no roles found"})
			c.Abort()
			return
		}

		roleStrs, ok := roles.([]string)
		if !ok || len(roleStrs) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid roles"})
			c.Abort()
			return
		}

		// Check if any role has the required permission
		authorized := false
		for _, roleStr := range roleStrs {
			if rbac.HasPermission(Role(roleStr), permission) {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates middleware that requires a specific role
func RequireRole(role Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "no roles found"})
			c.Abort()
			return
		}

		roleStrs, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid roles"})
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range roleStrs {
			if Role(r) == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("role %s required", role)})
			c.Abort()
			return
		}

		c.Next()
	}
}

