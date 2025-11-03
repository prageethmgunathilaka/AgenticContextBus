package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/acb/pkg/acb-sdk"
)

func main() {
	// Create client
	client := acb.NewClient(
		acb.WithEndpoint("http://localhost:8080"),
		acb.WithTimeout(30*time.Second),
	)
	defer client.Close()

	// Login to get token
	// TODO: Implement login in SDK
	log.Println("Agent A: Registering...")

	// Register agent
	req := &acb.RegisterAgentRequest{
		ID:           "agent-a",
		Type:         "hello-world",
		Location:     "local",
		Capabilities: []string{"greeting"},
		Metadata: map[string]string{
			"version": "1.0",
		},
	}

	agent, err := client.RegisterAgent(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to register agent: %v", err)
	}

	log.Printf("Agent A registered: %s", agent.ID)

	// Share context
	payload := []byte(`{"message": "Hello from Agent A!"}`)
	err = client.ShareContext(context.Background(), "greeting", payload,
		acb.WithScope(acb.ScopePublic),
		acb.WithTTL(time.Hour),
	)
	if err != nil {
		log.Fatalf("Failed to share context: %v", err)
	}

	log.Println("Agent A: Context shared successfully")

	// Send message to Agent B
	err = client.SendTo(context.Background(), "agent-b", "greetings", map[string]string{
		"from":    "agent-a",
		"message": "Hello Agent B!",
	})
	if err != nil {
		log.Printf("Warning: Failed to send message: %v", err)
	}

	log.Println("Agent A: Message sent to Agent B")

	// Keep running
	fmt.Println("Agent A running... Press Ctrl+C to stop")
	select {}
}
