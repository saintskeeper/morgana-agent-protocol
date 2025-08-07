package events

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var taskIDCounter int64

// GenerateTaskID generates a unique task identifier
func GenerateTaskID() string {
	counter := atomic.AddInt64(&taskIDCounter, 1)

	// Generate a random component for additional uniqueness
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to just counter if random generation fails
		return fmt.Sprintf("task_%d", counter)
	}

	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("task_%d_%s", counter, randomHex)
}

// GetTaskIDFromContext extracts task ID from context, or generates one if not present
func GetTaskIDFromContext(ctx context.Context) string {
	if taskID, ok := ctx.Value("task_id").(string); ok {
		return taskID
	}
	return GenerateTaskID()
}

// SetTaskIDInContext stores task ID in context
func SetTaskIDInContext(ctx context.Context, taskID string) context.Context {
	return context.WithValue(ctx, "task_id", taskID)
}

// EventMetrics provides common event metrics and timing utilities
type EventMetrics struct {
	mu            sync.RWMutex
	eventCount    map[EventType]int64
	lastEventTime map[EventType]int64 // Unix timestamp in nanoseconds
}

// NewEventMetrics creates a new EventMetrics instance
func NewEventMetrics() *EventMetrics {
	return &EventMetrics{
		eventCount:    make(map[EventType]int64),
		lastEventTime: make(map[EventType]int64),
	}
}

// RecordEvent records an event occurrence for metrics
func (em *EventMetrics) RecordEvent(eventType EventType) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.eventCount[eventType]++
	em.lastEventTime[eventType] = timeNow()
}

// GetEventCount returns the count for a specific event type
func (em *EventMetrics) GetEventCount(eventType EventType) int64 {
	em.mu.RLock()
	defer em.mu.RUnlock()

	return em.eventCount[eventType]
}

// timeNow returns current time in nanoseconds (for atomic operations)
func timeNow() int64 {
	return int64(time.Now().UnixNano())
}
