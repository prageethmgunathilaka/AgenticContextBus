package models

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidAgentID     = errors.New("invalid agent ID")
	ErrInvalidAgentType   = errors.New("invalid agent type")
	ErrInvalidContextID   = errors.New("invalid context ID")
	ErrInvalidContextType = errors.New("invalid context type")
	ErrPayloadTooLarge    = errors.New("payload too large for direct context")
	ErrInvalidScope       = errors.New("invalid context scope")
	ErrInvalidMessageType = errors.New("invalid message type")
)

const (
	MaxDirectContextSize = 1 * 1024 * 1024 // 1MB
)

// ValidateAgent validates an Agent struct
func (a *Agent) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("%w: agent ID cannot be empty", ErrInvalidAgentID)
	}
	if a.Type == "" {
		return fmt.Errorf("%w: agent type cannot be empty", ErrInvalidAgentType)
	}
	if a.Status == "" {
		a.Status = AgentStatusUnknown
	}
	return nil
}

// ValidateContext validates a Context struct
func (c *Context) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("%w: context type cannot be empty", ErrInvalidContextType)
	}
	if c.AgentID == "" {
		return errors.New("agent ID cannot be empty")
	}
	if len(c.Payload) > MaxDirectContextSize {
		return fmt.Errorf("%w: payload size %d exceeds max %d", ErrPayloadTooLarge, len(c.Payload), MaxDirectContextSize)
	}
	if err := c.AccessControl.Validate(); err != nil {
		return err
	}
	return nil
}

// ValidateAccessControl validates AccessControl struct
func (ac *AccessControl) Validate() error {
	switch ac.Scope {
	case ScopePublic, ScopePrivate, ScopeGroup, ScopeShared:
		// Valid scope
		if ac.Scope == ScopeGroup || ac.Scope == ScopeShared {
			if len(ac.AllowedIDs) == 0 {
				return fmt.Errorf("%w: %s scope requires allowed_ids", ErrInvalidScope, ac.Scope)
			}
		}
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidScope, ac.Scope)
	}
}

// ValidateMessage validates a Message struct
func (m *Message) Validate() error {
	if m.From == "" {
		return errors.New("message sender cannot be empty")
	}
	if m.Topic == "" {
		return errors.New("message topic cannot be empty")
	}
	if m.Type == "" {
		return fmt.Errorf("%w: message type cannot be empty", ErrInvalidMessageType)
	}
	switch m.Type {
	case MessageTypeEvent, MessageTypeCommand, MessageTypeQuery, MessageTypeResponse:
		// Valid type
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidMessageType, m.Type)
	}
}

// IsExpired checks if context is expired
func (c *Context) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(c.ExpiresAt)
}

// CalculateExpiration calculates expiration time from TTL
func (c *Context) CalculateExpiration() time.Time {
	if c.TTL <= 0 {
		return time.Time{}
	}
	return time.Now().Add(c.TTL)
}

