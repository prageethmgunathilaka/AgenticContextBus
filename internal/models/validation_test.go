package models

import (
	"testing"
	"time"
)

func TestAgent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		agent   *Agent
		wantErr bool
	}{
		{
			name: "valid agent",
			agent: &Agent{
				ID:   "agent-1",
				Type: "ml",
			},
			wantErr: false,
		},
		{
			name: "missing agent ID",
			agent: &Agent{
				Type: "ml",
			},
			wantErr: true,
		},
		{
			name: "missing agent type",
			agent: &Agent{
				ID: "agent-1",
			},
			wantErr: true,
		},
		{
			name: "sets default status",
			agent: &Agent{
				ID:   "agent-1",
				Type: "ml",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.agent.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Agent.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.agent.Status == "" {
				t.Errorf("Agent.Validate() should set default status")
			}
		})
	}
}

func TestContext_Validate(t *testing.T) {
	tests := []struct {
		name    string
		context *Context
		wantErr bool
	}{
		{
			name: "valid context",
			context: &Context{
				Type:    "user-profile",
				AgentID: "agent-1",
				Payload: []byte("small payload"),
				AccessControl: AccessControl{
					Scope: ScopePublic,
				},
			},
			wantErr: false,
		},
		{
			name: "payload too large",
			context: &Context{
				Type:    "user-profile",
				AgentID: "agent-1",
				Payload: make([]byte, MaxDirectContextSize+1),
				AccessControl: AccessControl{
					Scope: ScopePublic,
				},
			},
			wantErr: true,
		},
		{
			name: "missing context type",
			context: &Context{
				AgentID: "agent-1",
				Payload: []byte("data"),
				AccessControl: AccessControl{
					Scope: ScopePublic,
				},
			},
			wantErr: true,
		},
		{
			name: "missing agent ID",
			context: &Context{
				Type:    "user-profile",
				Payload: []byte("data"),
				AccessControl: AccessControl{
					Scope: ScopePublic,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.context.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Context.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccessControl_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ac      AccessControl
		wantErr bool
	}{
		{
			name: "valid public scope",
			ac: AccessControl{
				Scope: ScopePublic,
			},
			wantErr: false,
		},
		{
			name: "valid private scope",
			ac: AccessControl{
				Scope: ScopePrivate,
			},
			wantErr: false,
		},
		{
			name: "valid group scope with allowed IDs",
			ac: AccessControl{
				Scope:      ScopeGroup,
				AllowedIDs: []string{"agent-1", "agent-2"},
			},
			wantErr: false,
		},
		{
			name: "invalid group scope without allowed IDs",
			ac: AccessControl{
				Scope: ScopeGroup,
			},
			wantErr: true,
		},
		{
			name: "invalid scope",
			ac: AccessControl{
				Scope: ContextScope("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ac.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AccessControl.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{
		{
			name: "valid event message",
			message: &Message{
				From:  "agent-1",
				Topic: "events",
				Type:  MessageTypeEvent,
			},
			wantErr: false,
		},
		{
			name: "missing sender",
			message: &Message{
				Topic: "events",
				Type:  MessageTypeEvent,
			},
			wantErr: true,
		},
		{
			name: "missing topic",
			message: &Message{
				From: "agent-1",
				Type: MessageTypeEvent,
			},
			wantErr: true,
		},
		{
			name: "invalid message type",
			message: &Message{
				From:  "agent-1",
				Topic: "events",
				Type:  MessageType("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Message.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContext_IsExpired(t *testing.T) {
	tests := []struct {
		name    string
		context *Context
		want    bool
	}{
		{
			name: "not expired",
			context: &Context{
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "expired",
			context: &Context{
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			want: true,
		},
		{
			name: "zero time (no expiration)",
			context: &Context{
				ExpiresAt: time.Time{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.context.IsExpired(); got != tt.want {
				t.Errorf("Context.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContext_CalculateExpiration(t *testing.T) {
	ctx := &Context{
		TTL: 3600 * time.Second, // 1 hour
	}
	expiresAt := ctx.CalculateExpiration()
	if expiresAt.IsZero() {
		t.Error("CalculateExpiration() should return non-zero time")
	}
	if expiresAt.Before(time.Now()) {
		t.Error("CalculateExpiration() should return future time")
	}
}

