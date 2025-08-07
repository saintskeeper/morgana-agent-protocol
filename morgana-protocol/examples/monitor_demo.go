package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
)

func main() {
	fmt.Println("Morgana Monitor IPC Demo")
	fmt.Println("========================")

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run monitor_demo.go [server|client]")
		os.Exit(1)
	}

	mode := os.Args[1]
	socketPath := "/tmp/morgana_demo.sock"

	switch mode {
	case "server":
		runServer(socketPath)
	case "client":
		runClient(socketPath)
	default:
		fmt.Printf("Unknown mode: %s. Use 'server' or 'client'\n", mode)
		os.Exit(1)
	}
}

func runServer(socketPath string) {
	fmt.Println("Starting Morgana Monitor server...")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Subscribe to all events for demonstration
	eventBus.SubscribeAll(func(event events.Event) {
		fmt.Printf("[SERVER] Received event: %s (Task: %s) at %s\n",
			event.Type(), event.TaskID(), event.Timestamp().Format(time.RFC3339))
	})

	// Create and start IPC server
	server := monitor.NewIPCServer(socketPath, eventBus)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle cleanup
	defer func() {
		fmt.Println("Stopping server...")
		if err := server.Stop(); err != nil {
			log.Printf("Error stopping server: %v", err)
		}
	}()

	// Start server
	go func() {
		if err := server.Start(ctx); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	fmt.Printf("Server listening on %s\n", socketPath)
	fmt.Println("Press Ctrl+C to stop...")

	// Wait for interrupt
	select {}
}

func runClient(socketPath string) {
	fmt.Println("Starting Morgana Monitor client...")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Create IPC client
	client := monitor.NewIPCClient(socketPath, eventBus)

	// Try to connect
	if !client.TryConnect() {
		fmt.Printf("Failed to connect to server at %s\n", socketPath)
		fmt.Println("Make sure the server is running first.")
		os.Exit(1)
	}
	defer client.Close()

	fmt.Printf("Connected to server at %s\n", socketPath)

	// Simulate sending various events
	events_to_send := []struct {
		name  string
		event events.Event
	}{
		{
			name: "Task Started",
			event: events.NewTaskStartedEvent(
				context.Background(),
				"demo-task-1",
				"demo-agent",
				"Demo prompt for testing",
				map[string]interface{}{"demo": true},
				0,
				"gpt-4",
				"simple",
				time.Minute,
			),
		},
		{
			name: "Task Progress",
			event: events.NewTaskProgressEvent(
				context.Background(),
				"demo-task-1",
				"demo-agent",
				"validation",
				"Validating agent configuration",
				0.25,
				time.Second*5,
			),
		},
		{
			name: "Task Progress",
			event: events.NewTaskProgressEvent(
				context.Background(),
				"demo-task-1",
				"demo-agent",
				"execution",
				"Executing agent task",
				0.75,
				time.Second*20,
			),
		},
		{
			name: "Task Completed",
			event: events.NewTaskCompletedEvent(
				context.Background(),
				"demo-task-1",
				"demo-agent",
				"Demo task completed successfully!",
				time.Second*30,
				"gpt-4",
			),
		},
	}

	for i, eventData := range events_to_send {
		fmt.Printf("[CLIENT] Sending %s event (%d/%d)\n", eventData.name, i+1, len(events_to_send))

		// Publish event to local bus (will be forwarded to server)
		eventBus.Publish(eventData.event)

		// Wait between events for demo purposes
		time.Sleep(time.Second * 2)
	}

	fmt.Println("[CLIENT] All events sent successfully!")
	fmt.Println("[CLIENT] Disconnecting...")
}
