package server

import (
	"net/http"
	"time"

	"github.com/acb/internal/context"
	"github.com/acb/internal/models"
	"github.com/acb/internal/registry"
	"github.com/acb/internal/storage"
	"github.com/gin-gonic/gin"
)

// HTTP handlers implementation

func (s *HTTPServer) registerAgent(c *gin.Context) {
	var req registry.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant ID from context
	tenantID, _ := c.Get("tenant_id")
	if tenantIDStr, ok := tenantID.(string); ok && tenantIDStr != "" {
		req.TenantID = tenantIDStr
	} else {
		req.TenantID = "default" // MVP: single tenant
	}

	agent, err := s.registrySvc.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"agent": agent})
}

func (s *HTTPServer) listAgents(c *gin.Context) {
	filters := &storage.AgentFilters{
		Type:     c.Query("type"),
		Location: c.Query("location"),
		Status:   models.AgentStatus(c.Query("status")),
		Limit:    100,
		Offset:   0,
	}

	tenantID, _ := c.Get("tenant_id")
	if tenantIDStr, ok := tenantID.(string); ok {
		filters.TenantID = tenantIDStr
	} else {
		filters.TenantID = "default"
	}

	agents, err := s.registrySvc.Discover(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": agents,
		"total":  len(agents),
	})
}

func (s *HTTPServer) getAgent(c *gin.Context) {
	agentID := c.Param("agent_id")
	agent, err := s.registrySvc.Get(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

func (s *HTTPServer) unregisterAgent(c *gin.Context) {
	agentID := c.Param("agent_id")
	if err := s.registrySvc.Unregister(c.Request.Context(), agentID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *HTTPServer) sendHeartbeat(c *gin.Context) {
	agentID := c.Param("agent_id")
	if err := s.registrySvc.Heartbeat(c.Request.Context(), agentID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"acknowledged": true})
}

func (s *HTTPServer) createContext(c *gin.Context) {
	var req struct {
		Type          string                 `json:"type" binding:"required"`
		Payload       []byte                 `json:"payload" binding:"required"`
		Metadata      map[string]string      `json:"metadata"`
		Version       string                 `json:"version"`
		AccessControl models.AccessControl   `json:"access_control" binding:"required"`
		TTL           int                    `json:"ttl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agentID, _ := c.Get("agent_id")
	tenantID, _ := c.Get("tenant_id")

	createReq := &context.CreateRequest{
		Type:          req.Type,
		AgentID:       agentID.(string),
		TenantID:      tenantID.(string),
		Payload:       req.Payload,
		Metadata:      req.Metadata,
		Version:       req.Version,
		AccessControl: req.AccessControl,
	}

	if req.TTL > 0 {
		createReq.TTL = time.Duration(req.TTL) * time.Second
	}

	ctx, err := s.contextMgr.Create(c.Request.Context(), createReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"context": ctx})
}

func (s *HTTPServer) listContexts(c *gin.Context) {
	filters := &storage.ContextFilters{
		Type:    c.Query("type"),
		AgentID: c.Query("agent_id"),
		Limit:   100,
		Offset:  0,
	}

	tenantID, _ := c.Get("tenant_id")
	if tenantIDStr, ok := tenantID.(string); ok {
		filters.TenantID = tenantIDStr
	} else {
		filters.TenantID = "default"
	}

	contexts, err := s.contextMgr.List(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contexts": contexts,
		"total":    len(contexts),
	})
}

func (s *HTTPServer) getContext(c *gin.Context) {
	contextID := c.Param("context_id")
	ctx, err := s.contextMgr.Get(c.Request.Context(), contextID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "context not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"context": ctx})
}

func (s *HTTPServer) updateContext(c *gin.Context) {
	contextID := c.Param("context_id")
	var req struct {
		Payload       []byte               `json:"payload"`
		Metadata      map[string]string    `json:"metadata"`
		Version       string               `json:"version"`
		AccessControl models.AccessControl `json:"access_control"`
		TTL           int                  `json:"ttl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateReq := &context.UpdateRequest{
		Payload:       req.Payload,
		Metadata:      req.Metadata,
		Version:       req.Version,
		AccessControl: req.AccessControl,
	}

	if req.TTL > 0 {
		updateReq.TTL = time.Duration(req.TTL) * time.Second
	}

	ctx, err := s.contextMgr.Update(c.Request.Context(), contextID, updateReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"context": ctx})
}

func (s *HTTPServer) deleteContext(c *gin.Context) {
	contextID := c.Param("context_id")
	if err := s.contextMgr.Delete(c.Request.Context(), contextID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "context not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

