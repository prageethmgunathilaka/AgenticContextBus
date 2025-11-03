package main

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	const key = "TEST_GET_ENV_KEY"
	os.Unsetenv(key)
	if v := getEnv(key, "default"); v != "default" {
		t.Fatalf("expected default, got %s", v)
	}
	os.Setenv(key, "value")
	defer os.Unsetenv(key)
	if v := getEnv(key, "default"); v != "value" {
		t.Fatalf("expected value, got %s", v)
	}
}
