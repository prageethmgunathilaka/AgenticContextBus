//go:build cgo
// +build cgo

package router

import (
	"context"
	"testing"
)

func TestGetTopicName(t *testing.T) {
	if got := GetTopicName("tenant", "topic"); got != "acb.tenant.topic" {
		t.Fatalf("unexpected topic name: %s", got)
	}
}

func TestEnsureTopic_Noop(t *testing.T) {
	tm := NewTopicManager(nil)
	if err := tm.EnsureTopic(context.Background(), "topic"); err != nil {
		t.Fatalf("EnsureTopic returned error: %v", err)
	}
}
