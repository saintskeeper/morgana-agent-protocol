package monitor

import (
	"time"
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
