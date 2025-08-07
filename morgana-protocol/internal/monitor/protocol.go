package monitor

import (
	"time"
)

const (
	// DefaultEventBufferSize is the default size of the circular event buffer
	DefaultEventBufferSize = 1000

	// Message types for IPC communication
	MessageTypeEvent   = "event"   // Regular event message
	MessageTypeRequest = "request" // Client request (e.g., history)
	MessageTypeReplay  = "replay"  // Server replay of buffered events

	// Request types
	RequestHistory = "history" // Request event history
)

// IPCMessage represents a message sent over the IPC channel
type IPCMessage struct {
	Type      string      `json:"type"`      // Event type as string
	TaskID    string      `json:"task_id"`   // Task identifier
	Timestamp time.Time   `json:"timestamp"` // Event timestamp
	Data      interface{} `json:"data"`      // Event-specific data payload
}

// IPCResponse represents a response message (for future use)
type IPCResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// IPCRequest represents a client request to the server
type IPCRequest struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data,omitempty"`
}

// IPCReplay represents a batch of historical events
type IPCReplay struct {
	Events []IPCMessage `json:"events"`
	Count  int          `json:"count"`
}
