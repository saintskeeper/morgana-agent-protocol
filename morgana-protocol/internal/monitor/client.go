package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// IPCClient forwards events from a local EventBus to a remote monitor daemon
type IPCClient struct {
	socketPath     string
	eventBus       events.EventBus
	conn           net.Conn
	encoder        *json.Encoder
	subscriptionID string
	mu             sync.Mutex
	connected      bool
}

// NewIPCClient creates a new IPC client for event forwarding
func NewIPCClient(socketPath string, eventBus events.EventBus) *IPCClient {
	return &IPCClient{
		socketPath: socketPath,
		eventBus:   eventBus,
	}
}

// Connect establishes connection to the monitor daemon and starts forwarding events
func (c *IPCClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Try to connect to the Unix domain socket
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return fmt.Errorf("failed to connect to monitor daemon: %w", err)
	}

	c.conn = conn
	c.encoder = json.NewEncoder(conn)
	c.connected = true

	// Subscribe to all events and forward them
	c.subscriptionID = c.eventBus.SubscribeAll(c.forwardEvent)

	log.Printf("Connected to Morgana Monitor daemon")
	return nil
}

// Disconnect stops forwarding events and closes the connection
func (c *IPCClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	// Unsubscribe from events
	if c.subscriptionID != "" {
		c.eventBus.Unsubscribe(c.subscriptionID)
		c.subscriptionID = ""
	}

	// Close connection
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.encoder = nil
		c.connected = false

		if err != nil {
			return fmt.Errorf("error closing connection: %w", err)
		}
	}

	log.Printf("Disconnected from Morgana Monitor daemon")
	return nil
}

// Close is an alias for Disconnect for consistency with other interfaces
func (c *IPCClient) Close() error {
	return c.Disconnect()
}

// IsConnected returns whether the client is currently connected
func (c *IPCClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

// forwardEvent converts an Event to an IPCMessage and sends it to the monitor
func (c *IPCClient) forwardEvent(event events.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected || c.encoder == nil {
		return
	}

	// Create IPC message from event
	msg := IPCMessage{
		Type:      string(event.Type()),
		TaskID:    event.TaskID(),
		Timestamp: event.Timestamp(),
		Data:      event, // Send the entire event as data
	}

	// Send message to monitor daemon
	if err := c.encoder.Encode(msg); err != nil {
		log.Printf("Error forwarding event to monitor: %v", err)
		// On error, mark as disconnected to prevent further attempts
		c.connected = false
	}
}

// TryConnect attempts to connect to the monitor daemon, returns true if successful
func (c *IPCClient) TryConnect() bool {
	err := c.Connect()
	return err == nil
}
