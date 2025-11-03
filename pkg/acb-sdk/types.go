package acb

import "github.com/acb/internal/models"

// Types exports SDK types
type (
	Agent         = models.Agent
	Context       = models.Context
	Message       = models.Message
	AgentStatus   = models.AgentStatus
	MessageType   = models.MessageType
	ContextScope  = models.ContextScope
	AccessControl = models.AccessControl
)

// Constants exports SDK constants
const (
	ScopePublic  = models.ScopePublic
	ScopePrivate = models.ScopePrivate
	ScopeGroup   = models.ScopeGroup
	ScopeShared  = models.ScopeShared

	AgentStatusOnline  = models.AgentStatusOnline
	AgentStatusOffline = models.AgentStatusOffline
	AgentStatusUnknown = models.AgentStatusUnknown

	MessageTypeEvent    = models.MessageTypeEvent
	MessageTypeCommand  = models.MessageTypeCommand
	MessageTypeQuery    = models.MessageTypeQuery
	MessageTypeResponse = models.MessageTypeResponse
)
