package acb

import (
    "crypto/tls"
    "testing"
    "time"
    "context"
)

func TestNewClientOptions(t *testing.T) {
    c := NewClient(
        WithEndpoint("http://example"),
        WithCredentials("token"),
        WithTimeout(5*time.Second),
        WithTLS(&tls.Config{InsecureSkipVerify: true}),
    )

    if c.endpoint != "http://example" {
        t.Fatalf("endpoint not set")
    }
    if c.token != "token" {
        t.Fatalf("token not set")
    }
    if c.httpClient == nil || c.httpClient.Timeout != 5*time.Second {
        t.Fatalf("timeout not set")
    }
}

func TestClientNotImplementedAPIs(t *testing.T) {
    c := NewClient()
    if _, err := c.RegisterAgent(context.Background(), &RegisterAgentRequest{}); err == nil {
        t.Fatal("expected error for RegisterAgent")
    }
    if _, err := c.GetAgent(context.Background(), "id"); err == nil {
        t.Fatal("expected error for GetAgent")
    }
    if err := c.ShareContext(context.Background(), "t", []byte("x")); err == nil {
        t.Fatal("expected error for ShareContext")
    }
    if _, err := c.GetContext(context.Background(), "id"); err == nil {
        t.Fatal("expected error for GetContext")
    }
    if err := c.SendTo(context.Background(), "to", "topic", 1); err == nil {
        t.Fatal("expected error for SendTo")
    }
    if err := c.Broadcast(context.Background(), "topic", 1); err == nil {
        t.Fatal("expected error for Broadcast")
    }
}


