package auth

import (
	"context"
	"testing"
)

func TestRBAC_RequirePermission(t *testing.T) {
	r := NewRBAC()

	// explicit role provided
	if err := r.RequirePermission(context.Background(), Role("agent-full"), PermissionMessageSend); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// insufficient permissions with explicit role
	if err := r.RequirePermission(context.Background(), Role("observer"), PermissionContextDelete); err == nil {
		t.Fatalf("expected forbidden error")
	}

	// missing role and no roles in context
	if err := r.RequirePermission(context.Background(), "", PermissionAgentRead); err == nil {
		t.Fatalf("expected unauthorized error when no roles present")
	}
}
