package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// ReceiverClient connects to a monitor daemon and receives events (for TUI display)
type ReceiverClient struct {
	socketPath string
	eventBus   events.EventBus
	conn       net.Conn
	decoder    *json.Decoder
	encoder    *json.Encoder
	mu         sync.Mutex
	connected  bool
	stopCh     chan struct{}
}

// NewReceiverClient creates a new receiver client
func NewReceiverClient(socketPath string, eventBus events.EventBus) *ReceiverClient {
	return &ReceiverClient{
		socketPath: socketPath,
		eventBus:   eventBus,
		stopCh:     make(chan struct{}),
	}
}

// Connect establishes connection and starts receiving events
func (r *ReceiverClient) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Connect to the Unix domain socket
	conn, err := net.Dial("unix", r.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to monitor daemon: %w", err)
	}

	r.conn = conn
	r.decoder = json.NewDecoder(conn)
	r.encoder = json.NewEncoder(conn)
	r.connected = true

	// Request event history
	historyReq := IPCRequest{
		Type: MessageTypeRequest,
		Data: map[string]interface{}{
			"request": RequestHistory,
		},
	}

	if err := r.encoder.Encode(historyReq); err != nil {
		log.Printf("Failed to request event history: %v", err)
	} else {
		log.Printf("Requested event history from monitor")
	}

	// Start receiving events in background
	go r.receiveLoop()

	log.Printf("Receiver client connected to monitor")
	return nil
}

// receiveLoop continuously receives and processes messages
func (r *ReceiverClient) receiveLoop() {
	for {
		select {
		case <-r.stopCh:
			return
		default:
			var msg IPCMessage
			if err := r.decoder.Decode(&msg); err != nil {
				if err.Error() != "EOF" {
					log.Printf("Error decoding message: %v", err)
				}
				r.mu.Lock()
				r.connected = false
				r.mu.Unlock()
				return
			}

			// Handle different message types
			switch msg.Type {
			case MessageTypeReplay:
				// Process historical events
				if replay, ok := msg.Data.(map[string]interface{}); ok {
					if eventsData, ok := replay["events"].([]interface{}); ok {
						log.Printf("Received %d historical events", len(eventsData))
						for _, eventData := range eventsData {
							if eventMap, ok := eventData.(map[string]interface{}); ok {
								// Convert map back to IPCMessage
								var histMsg IPCMessage
								if msgBytes, err := json.Marshal(eventMap); err == nil {
									if err := json.Unmarshal(msgBytes, &histMsg); err == nil {
										// Reconstruct and publish the historical event
										r.processMessage(histMsg)
									}
								}
							}
						}
					}
				}
			default:
				// Regular event message
				r.processMessage(msg)
			}
		}
	}
}

// processMessage reconstructs an event and publishes it to the event bus
func (r *ReceiverClient) processMessage(msg IPCMessage) {
	// Create a generic event for TUI display
	// The TUI will handle the display based on event type
	genericEvent := &events.GenericEvent{
		BaseEvent: events.BaseEvent{
			EventType: events.EventType(msg.Type),
			Time:      msg.Timestamp,
			ID:        msg.TaskID,
			Ctx:       nil, // Context not needed for display
		},
		Data: msg.Data,
	}

	// Publish to event bus for TUI consumption
	r.eventBus.PublishAsync(genericEvent)
}

// Disconnect closes the connection
func (r *ReceiverClient) Disconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.connected {
		return nil
	}

	close(r.stopCh)

	if r.conn != nil {
		err := r.conn.Close()
		r.conn = nil
		r.decoder = nil
		r.encoder = nil
		r.connected = false
		return err
	}

	return nil
}

// IsConnected returns whether the client is connected
func (r *ReceiverClient) IsConnected() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.connected
}
