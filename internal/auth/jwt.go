package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/acb/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// JWTClaims represents JWT claims
type JWTClaims struct {
	AgentID  string   `json:"agent_id"`
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secretKey     []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	signingMethod jwt.SigningMethod
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string) *JWTManager {
	if secretKey == "" {
		// Generate a random key for development (NOT for production!)
		secretKey = "development-secret-key-change-in-production"
	}

	return &JWTManager{
		secretKey:     []byte(secretKey),
		accessTTL:     time.Duration(constants.DefaultAccessTokenTTL) * time.Second,
		refreshTTL:    time.Duration(constants.DefaultRefreshTokenTTL) * time.Second,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// GenerateAccessToken generates a new access token
func (m *JWTManager) GenerateAccessToken(agentID, tenantID string, roles []string) (string, error) {
	claims := JWTClaims{
		AgentID:  agentID,
		TenantID: tenantID,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	return token.SignedString(m.secretKey)
}

// GenerateRefreshToken generates a new refresh token
func (m *JWTManager) GenerateRefreshToken(agentID, tenantID string) (string, error) {
	claims := JWTClaims{
		AgentID:  agentID,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(m.signingMethod, claims)
	return token.SignedString(m.secretKey)
}

// ValidateToken validates and parses a token
func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// GenerateSecretKey generates a random secret key (for development)
func GenerateSecretKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", key), nil
}

