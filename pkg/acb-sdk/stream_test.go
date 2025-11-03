package acb

import (
    "bytes"
    "context"
    "testing"
)

func TestStreamBuilderChain(t *testing.T) {
    c := NewClient()
    sb := c.StreamContext(context.Background(), "file").
        FromReader(bytes.NewReader([]byte("data"))).
        WithChunkSize(1024).
        OnProgress(func(p float64) {})

    if sb == nil {
        t.Fatal("expected non-nil stream builder")
    }
    if err := sb.Send(); err == nil {
        t.Fatal("expected Send() to return error in MVP")
    }
}

func TestRequestAndSubscribe(t *testing.T) {
    c := NewClient()
    req := c.Request(context.Background(), "agent", "topic", 1)
    if req == nil {
        t.Fatal("expected non-nil request")
    }
    if _, err := req.Wait(); err == nil {
        t.Fatal("expected Wait() error in MVP")
    }

    sub := c.Subscribe("topic", func(ctx context.Context, cxt *Context) error { return nil })
    if sub == nil {
        t.Fatal("expected non-nil subscribe")
    }
    if sub.WithFilter("k", "v") == nil {
        t.Fatal("expected WithFilter to return self")
    }
    if err := sub.Unsubscribe(); err == nil {
        t.Fatal("expected Unsubscribe error in MVP")
    }
}


