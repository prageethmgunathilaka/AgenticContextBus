package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/acb/internal/errors"
)

// Role represents a user role
type Role string

const (
	RoleAdmin          Role = "admin"
	RoleAgentProducer  Role = "agent-producer"
	RoleAgentConsumer  Role = "agent-consumer"
	RoleAgentFull      Role = "agent-full"
	RoleObserver       Role = "observer"
)

// Permission represents a permission
type Permission string

const (
	PermissionAgentRegister Permission = "agent:register"
	PermissionAgentRead      Permission = "agent:read"
	PermissionAgentWrite     Permission = "agent:write"
	PermissionContextCreate  Permission = "context:create"
	PermissionContextRead    Permission = "context:read"
	PermissionContextWrite   Permission = "context:write"
	PermissionContextDelete  Permission = "context:delete"
	PermissionMessageSend    Permission = "message:send"
	PermissionMessageReceive Permission = "message:receive"
	PermissionStreamCreate   Permission = "stream:create"
	PermissionStreamRead     Permission = "stream:read"
)

// rolePermissions maps roles to their permissions
var rolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermissionAgentRegister, PermissionAgentRead, PermissionAgentWrite,
		PermissionContextCreate, PermissionContextRead, PermissionContextWrite, PermissionContextDelete,
		PermissionMessageSend, PermissionMessageReceive,
		PermissionStreamCreate, PermissionStreamRead,
	},
	RoleAgentProducer: {
		PermissionAgentRead,
		PermissionContextCreate, PermissionContextRead, PermissionContextWrite,
		PermissionMessageSend,
		PermissionStreamCreate,
	},
	RoleAgentConsumer: {
		PermissionAgentRead,
		PermissionContextRead,
		PermissionMessageReceive,
		PermissionStreamRead,
	},
	RoleAgentFull: {
		PermissionAgentRead,
		PermissionContextCreate, PermissionContextRead, PermissionContextWrite,
		PermissionMessageSend, PermissionMessageReceive,
		PermissionStreamCreate, PermissionStreamRead,
	},
	RoleObserver: {
		PermissionAgentRead,
		PermissionContextRead,
	},
}

// RBAC provides role-based access control
type RBAC struct{}

// NewRBAC creates a new RBAC instance
func NewRBAC() *RBAC {
	return &RBAC{}
}

// HasPermission checks if a role has a specific permission
func (r *RBAC) HasPermission(role Role, permission Permission) bool {
	permissions, ok := rolePermissions[role]
	if !ok {
		return false
	}

	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if a role has any of the specified permissions
func (r *RBAC) HasAnyPermission(role Role, permissions ...Permission) bool {
	for _, perm := range permissions {
		if r.HasPermission(role, perm) {
			return true
		}
	}
	return false
}

// RequirePermission checks permission and returns error if not authorized
func (r *RBAC) RequirePermission(ctx context.Context, role Role, permission Permission) error {
	// Get role from context if not provided
	if role == "" {
		roles, ok := ctx.Value("roles").([]string)
		if !ok || len(roles) == 0 {
			return errors.Unauthorized("no role found in context")
		}
		role = Role(roles[0])
	}

	if !r.HasPermission(role, permission) {
		return errors.Forbidden(fmt.Sprintf("role %s does not have permission %s", role, permission))
	}

	return nil
}

// GetPermissions returns all permissions for a role
func (r *RBAC) GetPermissions(role Role) []Permission {
	return rolePermissions[role]
}

// ParseRoles parses role strings into Role slice
func ParseRoles(roleStrs []string) []Role {
	roles := make([]Role, 0, len(roleStrs))
	for _, rs := range roleStrs {
		role := Role(strings.ToLower(rs))
		if _, ok := rolePermissions[role]; ok {
			roles = append(roles, role)
		}
	}
	return roles
}

