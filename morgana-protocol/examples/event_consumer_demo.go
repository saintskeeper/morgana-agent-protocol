package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/monitor"
)

func main() {
	log.SetPrefix("[event-consumer-demo] ")

	// Create event bus
	eventBus := events.NewEventBus(events.DefaultBusConfig())
	defer eventBus.Close()

	// Subscribe to all events for demonstration
	subscriptionID := eventBus.SubscribeAll(func(event events.Event) {
		log.Printf("Received event: %s (Task: %s, Time: %s)",
			event.Type(),
			event.TaskID(),
			event.Timestamp().Format(time.RFC3339))
	})
	defer eventBus.Unsubscribe(subscriptionID)

	// Create consumer configuration
	config := monitor.DefaultConsumerConfig()
	config.EventFile = "/tmp/morgana/events.jsonl"
	config.PollInterval = 50 * time.Millisecond // More frequent polling for demo

	// Create and start event consumer
	consumer, err := monitor.NewEventConsumer(eventBus, config)
	if err != nil {
		log.Fatalf("Failed to create event consumer: %v", err)
	}

	log.Printf("Starting event consumer for file: %s", config.EventFile)
	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start event consumer: %v", err)
	}
	defer consumer.Stop()

	// Create a sample event file for demonstration
	go createSampleEvents(config.EventFile)

	log.Printf("Event consumer started. Monitoring: %s", consumer.GetEventFile())
	log.Printf("Press Ctrl+C to stop")

	// Wait for interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Printf("Shutting down...")
}

// createSampleEvents creates some sample events for demonstration
func createSampleEvents(eventFile string) {
	// Wait a bit for the consumer to start
	time.Sleep(2 * time.Second)

	// Create the events directory if it doesn't exist
	os.MkdirAll("/tmp/morgana", 0755)

	file, err := os.OpenFile(eventFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to create sample events file: %v", err)
		return
	}
	defer file.Close()

	// Sample events in JSONL format
	sampleEvents := []string{
		`{"event_type":"task.started","task_id":"demo-001","timestamp":"` + time.Now().Format(time.RFC3339) + `","agent_type":"code-implementer","prompt":"Create a simple function","complexity":"low"}`,
		`{"event_type":"task.progress","task_id":"demo-001","timestamp":"` + time.Now().Add(1*time.Second).Format(time.RFC3339) + `","agent_type":"code-implementer","stage":"execution","message":"Processing request","progress":0.5}`,
		`{"event_type":"task.completed","task_id":"demo-001","timestamp":"` + time.Now().Add(3*time.Second).Format(time.RFC3339) + `","agent_type":"code-implementer","output":"Function created successfully","duration":"3s","model":"claude-sonnet-4"}`,
	}

	log.Printf("Creating sample events...")
	for i, event := range sampleEvents {
		time.Sleep(2 * time.Second) // Stagger the events
		if _, err := file.WriteString(event + "\n"); err != nil {
			log.Printf("Failed to write sample event %d: %v", i, err)
			return
		}
		file.Sync() // Force write to disk
		log.Printf("Wrote sample event %d", i+1)
	}

	log.Printf("Sample events completed")
}
