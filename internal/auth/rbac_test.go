package auth

import (
    "context"
    "testing"
)

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

