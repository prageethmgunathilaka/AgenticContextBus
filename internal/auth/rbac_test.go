package auth

import (
    "context"
    "testing"
)

func TestRBAC_HasPermissionMatrix(t *testing.T) {
    r := NewRBAC()

    // spot-check a few representative permissions per role
    cases := []struct{
        role Role
        perm Permission
        allow bool
    }{
        {RoleAdmin, PermissionAgentRegister, true},
        {RoleAdmin, PermissionStreamRead, true},
        {RoleObserver, PermissionContextRead, true},
        {RoleObserver, PermissionContextWrite, false},
        {RoleAgentProducer, PermissionMessageSend, true},
        {RoleAgentProducer, PermissionMessageReceive, false},
        {RoleAgentConsumer, PermissionMessageReceive, true},
        {RoleAgentConsumer, PermissionMessageSend, false},
        {RoleAgentFull, PermissionStreamCreate, true},
        {Role("unknown"), PermissionAgentRead, false},
    }

    for _, c := range cases {
        got := r.HasPermission(c.role, c.perm)
        if got != c.allow {
            t.Fatalf("role %s perm %s expected %v, got %v", c.role, c.perm, c.allow, got)
        }
    }
}

func TestRBAC_HasAnyPermission(t *testing.T) {
    r := NewRBAC()
    if !r.HasAnyPermission(RoleAgentProducer, PermissionContextRead, PermissionMessageSend) {
        t.Fatalf("expected HasAnyPermission to be true")
    }
    if r.HasAnyPermission(RoleObserver, PermissionContextDelete, PermissionMessageSend) {
        t.Fatalf("expected HasAnyPermission to be false")
    }
}

func TestRBAC_RequirePermission_WithContext(t *testing.T) {
    r := NewRBAC()
    // context provides roles when role arg is empty
    ctx := context.WithValue(context.Background(), "roles", []string{"agent-full"})
    if err := r.RequirePermission(ctx, "", PermissionMessageSend); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // insufficient permissions
    ctx2 := context.WithValue(context.Background(), "roles", []string{"observer"})
    if err := r.RequirePermission(ctx2, "", PermissionContextDelete); err == nil {
        t.Fatalf("expected forbidden error")
    }

    // no roles in context
    if err := r.RequirePermission(context.Background(), "", PermissionAgentRead); err == nil {
        t.Fatalf("expected unauthorized error when no roles present")
    }
}

func TestRBAC_GetPermissionsAndParseRoles(t *testing.T) {
    r := NewRBAC()
    perms := r.GetPermissions(RoleAgentFull)
    if len(perms) == 0 {
        t.Fatalf("expected permissions for agent-full")
    }

    parsed := ParseRoles([]string{"Admin", "observer", "nope"})
    if len(parsed) != 2 {
        t.Fatalf("expected 2 valid roles, got %d", len(parsed))
    }
    if parsed[0] != RoleAdmin && parsed[1] != RoleAdmin {
        t.Fatalf("expected admin to be parsed (case-insensitive)")
    }
}


