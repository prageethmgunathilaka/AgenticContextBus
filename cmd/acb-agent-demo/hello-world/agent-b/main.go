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

	log.Println("Agent B: Registering...")

	// Register agent
	req := &acb.RegisterAgentRequest{
		ID:           "agent-b",
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

	log.Printf("Agent B registered: %s", agent.ID)

	// Subscribe to greetings
	subscription := client.Subscribe("greetings", func(ctx context.Context, c *acb.Context) error {
		log.Printf("Agent B: Received context: %s", c.ID)
		log.Printf("Agent B: Message: %s", string(c.Payload))
		return nil
	})
	defer func() {
		_ = subscription.Unsubscribe()
	}()

	log.Println("Agent B: Subscribed to greetings")

	// Send response
	err = client.SendTo(context.Background(), "agent-a", "greetings", map[string]string{
		"from":    "agent-b",
		"message": "Hello Agent A!",
	})
	if err != nil {
		log.Printf("Warning: Failed to send message: %v", err)
	}

	log.Println("Agent B: Message sent to Agent A")

	// Keep running
	fmt.Println("Agent B running... Press Ctrl+C to stop")
	select {}
}
