package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")

	token, err := mgr.GenerateAccessToken("agent-1", "tenant-1", []string{"agent-full"})
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")

	token, err := mgr.GenerateRefreshToken("agent-1", "tenant-1")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTManager_ValidateToken(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")

	// Generate token
	token, err := mgr.GenerateAccessToken("agent-1", "tenant-1", []string{"agent-full"})
	require.NoError(t, err)

	// Validate token
	claims, err := mgr.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "agent-1", claims.AgentID)
	assert.Equal(t, "tenant-1", claims.TenantID)
	assert.Equal(t, []string{"agent-full"}, claims.Roles)
}

func TestJWTManager_ValidateToken_Invalid(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")

	// Invalid token
	_, err := mgr.ValidateToken("invalid-token")
	assert.Error(t, err)
}

func TestJWTManager_ValidateToken_Expired(t *testing.T) {
	mgr := NewJWTManager("test-secret-key")
	mgr.accessTTL = -1 * time.Hour // Set negative TTL for testing

	// Generate expired token
	token, err := mgr.GenerateAccessToken("agent-1", "tenant-1", []string{"agent-full"})
	require.NoError(t, err)

	// Should fail validation
	_, err = mgr.ValidateToken(token)
	assert.Error(t, err)
}

func TestGenerateSecretKey(t *testing.T) {
	key, err := GenerateSecretKey()
	require.NoError(t, err)
	assert.NotEmpty(t, key)
	assert.Len(t, key, 64) // Hex string of 32 bytes
}

func TestRBAC_HasPermission(t *testing.T) {
	rbac := NewRBAC()

	tests := []struct {
		name       string
		role       Role
		permission Permission
		want       bool
	}{
		{
			name:       "admin has all permissions",
			role:       RoleAdmin,
			permission: PermissionAgentRegister,
			want:       true,
		},
		{
			name:       "agent-producer can create context",
			role:       RoleAgentProducer,
			permission: PermissionContextCreate,
			want:       true,
		},
		{
			name:       "agent-producer cannot receive messages",
			role:       RoleAgentProducer,
			permission: PermissionMessageReceive,
			want:       false,
		},
		{
			name:       "agent-consumer can receive messages",
			role:       RoleAgentConsumer,
			permission: PermissionMessageReceive,
			want:       true,
		},
		{
			name:       "observer can only read",
			role:       RoleObserver,
			permission: PermissionContextRead,
			want:       true,
		},
		{
			name:       "observer cannot write",
			role:       RoleObserver,
			permission: PermissionContextWrite,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rbac.HasPermission(tt.role, tt.permission)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRBAC_HasAnyPermission(t *testing.T) {
	rbac := NewRBAC()

	// Admin has all permissions
	got := rbac.HasAnyPermission(RoleAdmin, PermissionAgentRegister, PermissionContextCreate)
	assert.True(t, got)

	// Agent-producer has some but not all
	got = rbac.HasAnyPermission(RoleAgentProducer, PermissionContextCreate, PermissionMessageReceive)
	assert.True(t, got)

	// Observer has none
	got = rbac.HasAnyPermission(RoleObserver, PermissionContextWrite, PermissionMessageSend)
	assert.False(t, got)
}

func TestRBAC_GetPermissions(t *testing.T) {
	rbac := NewRBAC()

	permissions := rbac.GetPermissions(RoleAdmin)
	assert.Greater(t, len(permissions), 0)

	permissions = rbac.GetPermissions(RoleAgentProducer)
	assert.Greater(t, len(permissions), 0)
}

func TestParseRoles(t *testing.T) {
	tests := []struct {
		name     string
		roleStrs []string
		want     int
	}{
		{
			name:     "valid roles",
			roleStrs: []string{"admin", "agent-full"},
			want:     2,
		},
		{
			name:     "invalid roles filtered",
			roleStrs: []string{"admin", "invalid-role"},
			want:     1,
		},
		{
			name:     "case insensitive",
			roleStrs: []string{"ADMIN", "Agent-Full"},
			want:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roles := ParseRoles(tt.roleStrs)
			assert.Len(t, roles, tt.want)
		})
	}
}

