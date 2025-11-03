package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/acb/internal/auth"
	ctxmgr "github.com/acb/internal/context"
	"github.com/acb/internal/models"
	"github.com/acb/internal/registry"
	"github.com/acb/internal/storage"
	"github.com/gin-gonic/gin"
)

type inMemoryAgentStore struct{ m map[string]*models.Agent }

func (s *inMemoryAgentStore) Create(ctx context.Context, a *models.Agent) error {
	if s.m == nil {
		s.m = map[string]*models.Agent{}
	}
	if _, ok := s.m[a.ID]; ok {
		return nil
	}
	cp := *a
	s.m[a.ID] = &cp
	return nil
}
func (s *inMemoryAgentStore) Get(ctx context.Context, id string) (*models.Agent, error) {
	if a, ok := s.m[id]; ok {
		cp := *a
		return &cp, nil
	}
	return nil, fmt.Errorf("agent not found")
}
func (s *inMemoryAgentStore) Update(ctx context.Context, a *models.Agent) error {
	if _, ok := s.m[a.ID]; !ok {
		return fmt.Errorf("agent not found")
	}
	cp := *a
	s.m[a.ID] = &cp
	return nil
}
func (s *inMemoryAgentStore) Delete(ctx context.Context, id string) error { delete(s.m, id); return nil }
func (s *inMemoryAgentStore) List(ctx context.Context, f *storage.AgentFilters) ([]*models.Agent, error) {
	res := []*models.Agent{}
	for _, a := range s.m {
		cp := *a
		res = append(res, &cp)
	}
	return res, nil
}
func (s *inMemoryAgentStore) UpdateLastSeen(ctx context.Context, id string) error {
	if a, ok := s.m[id]; ok {
		a.LastSeen = time.Now()
		a.Status = models.AgentStatusOnline
		return nil
	}
	return fmt.Errorf("agent not found")
}

type inMemoryContextStore struct{ m map[string]*models.Context }

func (s *inMemoryContextStore) Create(ctx context.Context, c *models.Context) error {
	if s.m == nil {
		s.m = map[string]*models.Context{}
	}
	cp := *c
	s.m[c.ID] = &cp
	return nil
}
func (s *inMemoryContextStore) Get(ctx context.Context, id string) (*models.Context, error) {
	if c, ok := s.m[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, fmt.Errorf("context not found")
}
func (s *inMemoryContextStore) Update(ctx context.Context, c *models.Context) error {
	if _, ok := s.m[c.ID]; !ok {
		return fmt.Errorf("context not found")
	}
	cp := *c
	s.m[c.ID] = &cp
	return nil
}
func (s *inMemoryContextStore) Delete(ctx context.Context, id string) error {
	if _, ok := s.m[id]; !ok {
		return fmt.Errorf("context not found")
	}
	delete(s.m, id)
	return nil
}
func (s *inMemoryContextStore) List(ctx context.Context, f *storage.ContextFilters) ([]*models.Context, error) {
	res := []*models.Context{}
	for _, c := range s.m {
		cp := *c
		res = append(res, &cp)
	}
	return res, nil
}
func (s *inMemoryContextStore) DeleteExpired(ctx context.Context) (int, error) { return 0, nil }

func makeServerForHandlersTest(t *testing.T) *HTTPServer {
	t.Helper()
	gin.SetMode(gin.TestMode)
	jwtManager := auth.NewJWTManager("test-secret")
	rsvc := registry.NewService(&inMemoryAgentStore{})
	cmgr := ctxmgr.NewManager(&inMemoryContextStore{})
	return NewHTTPServer("8080", rsvc, cmgr, jwtManager, auth.NewRBAC())
}

func authHeader(t *testing.T, jwt *auth.JWTManager) string {
	t.Helper()
	tok, err := jwt.GenerateAccessToken("agent-1", "default", []string{"agent-full"})
	if err != nil {
		t.Fatalf("token error: %v", err)
	}
	return "Bearer " + tok
}

func TestRefreshTokenEndpoint(t *testing.T) {
	httpSrv := makeServerForHandlersTest(t)

	// invalid payload
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// valid refresh flow
	refresh, err := httpSrv.jwtManager.GenerateRefreshToken("agent-1", "default")
	if err != nil {
		t.Fatalf("gen refresh: %v", err)
	}
	body, _ := json.Marshal(map[string]string{"refresh_token": refresh})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestLoginEndpoint(t *testing.T) {
	httpSrv := makeServerForHandlersTest(t)

	body, _ := json.Marshal(map[string]string{"username": "u", "password": "p"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestContextHandlers_CRUD(t *testing.T) {
	httpSrv := makeServerForHandlersTest(t)
	jwt := httpSrv.jwtManager
	hdr := authHeader(t, jwt)

	// create context
	create := map[string]any{
		"type":           "greeting",
		"payload":        []byte("hello"),
		"metadata":       map[string]string{"k": "v"},
		"version":        "1",
		"access_control": map[string]any{"scope": "public"},
		"ttl":            60,
	}
	body, _ := json.Marshal(create)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/contexts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", hdr)
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create expected 201, got %d", w.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	ctxObj := resp["context"].(map[string]any)
	ctxID := ctxObj["id"].(string)

	// list contexts
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/contexts", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list expected 200, got %d", w.Code)
	}

	// get context
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/contexts/"+ctxID, nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get expected 200, got %d", w.Code)
	}

	// update context
	update := map[string]any{
		"metadata": map[string]string{"k": "v2"},
		"version":  "2",
		"ttl":      120,
	}
	body, _ = json.Marshal(update)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/contexts/"+ctxID, bytes.NewBuffer(body))
	req.Header.Set("Authorization", hdr)
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update expected 200, got %d", w.Code)
	}

	// delete context
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/contexts/"+ctxID, nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete expected 204, got %d", w.Code)
	}
}

func TestHandlers_BadBodiesAndNotFound(t *testing.T) {
	httpSrv := makeServerForHandlersTest(t)
	jwt := httpSrv.jwtManager
	hdr := authHeader(t, jwt)

	// registerAgent bad body
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/agents", bytes.NewBufferString("{}"))
	req.Header.Set("Authorization", hdr)
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("register bad body expected 400, got %d", w.Code)
	}

	// createContext bad body
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/contexts", bytes.NewBufferString("{}"))
	req.Header.Set("Authorization", hdr)
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("create ctx bad body expected 400, got %d", w.Code)
	}

	// getAgent not found
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/agents/missing", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("get missing expected 404, got %d", w.Code)
	}

	// deleteContext not found
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/contexts/missing", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("delete missing expected 404, got %d", w.Code)
	}
}

func TestHandlers_NilContextManager(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := auth.NewJWTManager("s")
	httpSrv := NewHTTPServer("8080", nil, nil, jwtManager, auth.NewRBAC())
	hdr := authHeader(t, jwtManager)

	w := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/contexts", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestAgentHandlers_CRUD(t *testing.T) {
	httpSrv := makeServerForHandlersTest(t)
	jwt := httpSrv.jwtManager
	hdr := authHeader(t, jwt)

	// register
	reg := map[string]any{
		"id":           "agent-x",
		"type":         "ml",
		"location":     "us",
		"capabilities": []string{"a"},
		"metadata":     map[string]string{"k": "v"},
	}
	body, _ := json.Marshal(reg)
	w := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/agents", bytes.NewBuffer(body))
	req.Header.Set("Authorization", hdr)
	req.Header.Set("Content-Type", "application/json")
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("register expected 201, got %d", w.Code)
	}

	// list
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/agents", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list expected 200, got %d", w.Code)
	}

	// get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/agents/agent-x", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get expected 200, got %d", w.Code)
	}

	// heartbeat
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/agents/agent-x/heartbeat", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("heartbeat expected 200, got %d", w.Code)
	}

	// delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/agents/agent-x", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete expected 204, got %d", w.Code)
	}

	// get after delete -> 404
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/agents/agent-x", nil)
	req.Header.Set("Authorization", hdr)
	httpSrv.router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("get after delete expected 404, got %d", w.Code)
	}
}


