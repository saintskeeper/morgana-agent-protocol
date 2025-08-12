package monitor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// EventConsumer monitors and consumes events from /tmp/morgana/events.jsonl
type EventConsumer struct {
	eventBus     events.EventBus
	eventFile    string
	watcher      *fsnotify.Watcher
	file         *os.File
	reader       *bufio.Reader
	offset       int64
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	pollInterval time.Duration
	running      bool
}

// ConsumerConfig provides configuration options for the EventConsumer
type ConsumerConfig struct {
	// EventFile is the path to the events JSONL file (default: /tmp/morgana/events.jsonl)
	EventFile string

	// PollInterval for fallback polling when fsnotify isn't available (default: 100ms)
	PollInterval time.Duration

	// BufferSize for the file reader (default: 64KB)
	BufferSize int
}

// DefaultConsumerConfig returns a default configuration
func DefaultConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		EventFile:    "/tmp/morgana/events.jsonl",
		PollInterval: 100 * time.Millisecond,
		BufferSize:   64 * 1024,
	}
}

// NewEventConsumer creates a new EventConsumer instance
func NewEventConsumer(eventBus events.EventBus, config ConsumerConfig) (*EventConsumer, error) {
	if config.EventFile == "" {
		config.EventFile = "/tmp/morgana/events.jsonl"
	}
	if config.PollInterval <= 0 {
		config.PollInterval = 100 * time.Millisecond
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 64 * 1024
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &EventConsumer{
		eventBus:     eventBus,
		eventFile:    config.EventFile,
		pollInterval: config.PollInterval,
		ctx:          ctx,
		cancel:       cancel,
		offset:       0,
		running:      false,
	}

	return consumer, nil
}

// Start begins consuming events from the JSONL file
func (ec *EventConsumer) Start() error {
	ec.mu.Lock()
	if ec.running {
		ec.mu.Unlock()
		return fmt.Errorf("consumer is already running")
	}
	ec.running = true
	ec.mu.Unlock()

	log.Printf("EventConsumer starting to monitor: %s", ec.eventFile)

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(ec.eventFile), 0755); err != nil {
		return fmt.Errorf("failed to create events directory: %w", err)
	}

	// Try to set up filesystem watcher first
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		ec.watcher = watcher
		// Watch both the file and its directory for rotation handling
		go ec.watchWithFsnotify()
	} else {
		log.Printf("Failed to create fsnotify watcher, falling back to polling: %v", err)
		go ec.watchWithPolling()
	}

	// Start the main reading loop
	go ec.readLoop()

	return nil
}

// Stop stops the event consumer
func (ec *EventConsumer) Stop() error {
	ec.mu.Lock()
	if !ec.running {
		ec.mu.Unlock()
		return nil
	}
	ec.running = false
	ec.mu.Unlock()

	ec.cancel()

	if ec.watcher != nil {
		ec.watcher.Close()
	}

	if ec.file != nil {
		ec.file.Close()
	}

	log.Printf("EventConsumer stopped")
	return nil
}

// watchWithFsnotify uses fsnotify to monitor file changes
func (ec *EventConsumer) watchWithFsnotify() {
	defer ec.watcher.Close()

	// Watch the events directory for file creation/rotation
	dir := filepath.Dir(ec.eventFile)
	if err := ec.watcher.Add(dir); err != nil {
		log.Printf("Failed to watch directory %s: %v", dir, err)
		// Fall back to polling
		go ec.watchWithPolling()
		return
	}

	// Also watch the file itself if it exists
	if _, err := os.Stat(ec.eventFile); err == nil {
		if err := ec.watcher.Add(ec.eventFile); err != nil {
			log.Printf("Failed to watch file %s: %v", ec.eventFile, err)
		}
	}

	for {
		select {
		case <-ec.ctx.Done():
			return
		case event, ok := <-ec.watcher.Events:
			if !ok {
				return
			}

			// Handle different types of file events
			if event.Name == ec.eventFile {
				if event.Op&fsnotify.Write == fsnotify.Write {
					// File was written to - trigger read
					ec.triggerRead()
				} else if event.Op&fsnotify.Create == fsnotify.Create {
					// File was created (possibly after rotation) - reopen
					ec.reopenFile()
				} else if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					// File was removed or renamed (rotation) - prepare for new file
					ec.handleFileRotation()
				}
			}

		case err, ok := <-ec.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

// watchWithPolling uses polling as a fallback when fsnotify is unavailable
func (ec *EventConsumer) watchWithPolling() {
	ticker := time.NewTicker(ec.pollInterval)
	defer ticker.Stop()

	var lastModTime time.Time
	var lastSize int64

	for {
		select {
		case <-ec.ctx.Done():
			return
		case <-ticker.C:
			stat, err := os.Stat(ec.eventFile)
			if err != nil {
				if !os.IsNotExist(err) {
					log.Printf("Error stating file %s: %v", ec.eventFile, err)
				}
				continue
			}

			// Check if file has been modified or size changed
			if stat.ModTime().After(lastModTime) || stat.Size() != lastSize {
				if stat.Size() < lastSize {
					// File was truncated - likely rotated
					ec.handleFileRotation()
				} else {
					// File was appended to
					ec.triggerRead()
				}
				lastModTime = stat.ModTime()
				lastSize = stat.Size()
			}
		}
	}
}

// triggerRead signals the read loop to read new data
func (ec *EventConsumer) triggerRead() {
	// In this simple implementation, the read loop continuously reads
	// This method could be enhanced with channels for more precise control
}

// reopenFile reopens the events file (useful after rotation)
func (ec *EventConsumer) reopenFile() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.file != nil {
		ec.file.Close()
		ec.file = nil
		ec.reader = nil
	}
	ec.offset = 0

	// The read loop will reopen the file on next iteration
}

// handleFileRotation handles file rotation scenarios
func (ec *EventConsumer) handleFileRotation() {
	log.Printf("Detected file rotation for %s", ec.eventFile)
	ec.reopenFile()
}

// readLoop continuously reads from the events file
func (ec *EventConsumer) readLoop() {
	for {
		select {
		case <-ec.ctx.Done():
			return
		default:
			if err := ec.ensureFileOpen(); err != nil {
				log.Printf("Error opening events file: %v", err)
				time.Sleep(time.Second) // Wait before retrying
				continue
			}

			// Try to read a line
			line, err := ec.reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					// No more data available, wait a bit
					time.Sleep(10 * time.Millisecond)
					continue
				}
				log.Printf("Error reading from events file: %v", err)
				ec.reopenFile() // Force reopen on error
				time.Sleep(time.Second)
				continue
			}

			// Update offset
			ec.mu.Lock()
			ec.offset += int64(len(line))
			ec.mu.Unlock()

			// Process the event line
			if err := ec.processEventLine(line); err != nil {
				log.Printf("Error processing event line: %v", err)
			}
		}
	}
}

// ensureFileOpen ensures the events file is open for reading
func (ec *EventConsumer) ensureFileOpen() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if ec.file != nil {
		return nil
	}

	file, err := os.OpenFile(ec.eventFile, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, create an empty one
			if f, createErr := os.OpenFile(ec.eventFile, os.O_CREATE|os.O_RDONLY, 0644); createErr == nil {
				f.Close()
				// Now open it normally
				file, err = os.OpenFile(ec.eventFile, os.O_RDONLY, 0644)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to open events file: %w", err)
		}
	}

	// Seek to the last known position
	if ec.offset > 0 {
		if _, seekErr := file.Seek(ec.offset, 0); seekErr != nil {
			log.Printf("Warning: failed to seek to offset %d: %v", ec.offset, seekErr)
			ec.offset = 0 // Reset offset on seek failure
		}
	}

	ec.file = file
	ec.reader = bufio.NewReader(file)

	return nil
}

// processEventLine parses and publishes a single event line
func (ec *EventConsumer) processEventLine(line []byte) error {
	// Trim whitespace and skip empty lines
	line = []byte(string(line))
	if len(line) == 0 {
		return nil
	}

	// Try parsing as raw event data first (more common format)
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(line, &rawEvent); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Try to convert raw event to IPCMessage format
	var msg IPCMessage
	if err := ec.convertRawEventToIPCMessage(rawEvent, &msg); err != nil {
		// If conversion fails, try parsing directly as IPCMessage
		if jsonErr := json.Unmarshal(line, &msg); jsonErr != nil {
			return fmt.Errorf("failed to convert raw event (%v) and parse as IPCMessage (%v)", err, jsonErr)
		}
	}

	// Convert IPCMessage to proper Event and publish
	event := ec.reconstructEvent(msg)
	if event != nil {
		// Use async publish for better performance
		if !ec.eventBus.PublishAsync(event) {
			log.Printf("Warning: failed to publish event (queue full): %s", msg.Type)
		}
		return nil
	}

	return fmt.Errorf("failed to reconstruct event from message type: %s", msg.Type)
}

// convertRawEventToIPCMessage converts a raw JSON event to IPCMessage format
func (ec *EventConsumer) convertRawEventToIPCMessage(rawEvent map[string]interface{}, msg *IPCMessage) error {
	// Extract basic fields
	if eventType, ok := rawEvent["event_type"].(string); ok {
		msg.Type = eventType
	} else if eventType, ok := rawEvent["type"].(string); ok {
		msg.Type = eventType
	} else {
		return fmt.Errorf("missing event type")
	}

	if taskID, ok := rawEvent["task_id"].(string); ok {
		msg.TaskID = taskID
	} else if taskID, ok := rawEvent["id"].(string); ok {
		msg.TaskID = taskID
	}

	// Parse timestamp with multiple formats
	if timestampStr, ok := rawEvent["timestamp"].(string); ok {
		// Try various timestamp formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05.999999999Z07:00", // Extended RFC3339 with nanoseconds
			"2006-01-02T15:04:05Z07:00",           // RFC3339 without fractions
			time.RFC1123,
			time.RFC1123Z,
		}

		for _, format := range formats {
			if t, err := time.Parse(format, timestampStr); err == nil {
				msg.Timestamp = t
				break
			}
		}

		// If all parsing failed, use current time
		if msg.Timestamp.IsZero() {
			msg.Timestamp = time.Now()
		}
	} else {
		msg.Timestamp = time.Now()
	}

	// Set the data payload to the entire raw event
	msg.Data = rawEvent

	return nil
}

// reconstructEvent converts an IPCMessage back to a proper Event type
// This is similar to the server.go implementation but adapted for the consumer
func (ec *EventConsumer) reconstructEvent(msg IPCMessage) events.Event {
	ctx := context.Background()

	baseEvent := events.BaseEvent{
		EventType: events.EventType(msg.Type),
		Time:      msg.Timestamp,
		ID:        msg.TaskID,
		Ctx:       ctx,
	}

	// Handle both IPCMessage format and direct event format
	var data map[string]interface{}
	var ok bool

	if data, ok = msg.Data.(map[string]interface{}); !ok {
		// If Data is not a map, create a generic event
		return &events.GenericEvent{
			BaseEvent: baseEvent,
			Data:      msg.Data,
		}
	}

	// Reconstruct specific event types
	switch events.EventType(msg.Type) {
	case events.EventTaskStarted:
		event := &events.TaskStartedEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if prompt, ok := data["prompt"].(string); ok {
			event.Prompt = prompt
		}
		if modelHint, ok := data["model_hint"].(string); ok {
			event.ModelHint = modelHint
		}
		if complexity, ok := data["complexity"].(string); ok {
			event.Complexity = complexity
		}
		if retryCount, ok := data["retry_count"].(float64); ok {
			event.RetryCount = int(retryCount)
		}
		if timeoutStr, ok := data["timeout"].(string); ok {
			if timeout, err := time.ParseDuration(timeoutStr); err == nil {
				event.Timeout = timeout
			}
		}
		if options, ok := data["options"].(map[string]interface{}); ok {
			event.Options = options
		}
		return event

	case events.EventTaskProgress:
		event := &events.TaskProgressEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if stage, ok := data["stage"].(string); ok {
			event.Stage = stage
		}
		if message, ok := data["message"].(string); ok {
			event.Message = message
		}
		if progress, ok := data["progress"].(float64); ok {
			event.Progress = progress
		}
		if durationStr, ok := data["duration"].(string); ok {
			if duration, err := time.ParseDuration(durationStr); err == nil {
				event.Duration = duration
			}
		}
		return event

	case events.EventTaskCompleted:
		event := &events.TaskCompletedEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if output, ok := data["output"].(string); ok {
			event.Output = output
			event.OutputLength = len(output)
		}
		if outputLength, ok := data["output_length"].(float64); ok {
			event.OutputLength = int(outputLength)
		}
		if model, ok := data["model"].(string); ok {
			event.Model = model
		}
		if durationStr, ok := data["duration"].(string); ok {
			if duration, err := time.ParseDuration(durationStr); err == nil {
				event.Duration = duration
			}
		}
		return event

	case events.EventTaskFailed:
		event := &events.TaskFailedEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if errorMsg, ok := data["error"].(string); ok {
			event.Error = errorMsg
		}
		if stage, ok := data["stage"].(string); ok {
			event.Stage = stage
		}
		if retryCount, ok := data["retry_count"].(float64); ok {
			event.RetryCount = int(retryCount)
		}
		if durationStr, ok := data["duration"].(string); ok {
			if duration, err := time.ParseDuration(durationStr); err == nil {
				event.Duration = duration
			}
		}
		return event

	case events.EventOrchestratorStarted:
		event := &events.OrchestratorStartedEvent{BaseEvent: baseEvent}
		if mode, ok := data["mode"].(string); ok {
			event.Mode = mode
		}
		if taskCount, ok := data["task_count"].(float64); ok {
			event.TaskCount = int(taskCount)
		}
		if maxConcurrency, ok := data["max_concurrency"].(float64); ok {
			event.MaxConcurrency = int(maxConcurrency)
		}
		return event

	case events.EventOrchestratorCompleted:
		event := &events.OrchestratorCompletedEvent{BaseEvent: baseEvent}
		if mode, ok := data["mode"].(string); ok {
			event.Mode = mode
		}
		if taskCount, ok := data["task_count"].(float64); ok {
			event.TaskCount = int(taskCount)
		}
		if successCount, ok := data["success_count"].(float64); ok {
			event.SuccessCount = int(successCount)
		}
		if failureCount, ok := data["failure_count"].(float64); ok {
			event.FailureCount = int(failureCount)
		}
		if durationStr, ok := data["duration"].(string); ok {
			if duration, err := time.ParseDuration(durationStr); err == nil {
				event.Duration = duration
			}
		}
		return event

	case events.EventOrchestratorFailed:
		event := &events.OrchestratorFailedEvent{BaseEvent: baseEvent}
		if mode, ok := data["mode"].(string); ok {
			event.Mode = mode
		}
		if taskCount, ok := data["task_count"].(float64); ok {
			event.TaskCount = int(taskCount)
		}
		if completedCount, ok := data["completed_count"].(float64); ok {
			event.CompletedCount = int(completedCount)
		}
		if errorMsg, ok := data["error"].(string); ok {
			event.Error = errorMsg
		}
		if durationStr, ok := data["duration"].(string); ok {
			if duration, err := time.ParseDuration(durationStr); err == nil {
				event.Duration = duration
			}
		}
		return event

	case events.EventAdapterValidation:
		event := &events.AdapterValidationEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if valid, ok := data["valid"].(bool); ok {
			event.Valid = valid
		}
		if errorMsg, ok := data["error"].(string); ok {
			event.Error = errorMsg
		}
		return event

	case events.EventAdapterPromptLoad:
		event := &events.AdapterPromptLoadEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if success, ok := data["success"].(bool); ok {
			event.Success = success
		}
		if promptLength, ok := data["prompt_length"].(float64); ok {
			event.PromptLength = int(promptLength)
		}
		if errorMsg, ok := data["error"].(string); ok {
			event.Error = errorMsg
		}
		return event

	case events.EventAdapterExecution:
		event := &events.AdapterExecutionEvent{BaseEvent: baseEvent}
		if agentType, ok := data["agent_type"].(string); ok {
			event.AgentType = agentType
		}
		if phase, ok := data["phase"].(string); ok {
			event.Phase = phase
		}
		if model, ok := data["model"].(string); ok {
			event.Model = model
		}
		if success, ok := data["success"].(bool); ok {
			event.Success = success
		}
		if errorMsg, ok := data["error"].(string); ok {
			event.Error = errorMsg
		}
		if durationStr, ok := data["duration"].(string); ok {
			if duration, err := time.ParseDuration(durationStr); err == nil {
				event.Duration = duration
			}
		}
		return event
	}

	// Return a generic event if we can't reconstruct the specific type
	return &events.GenericEvent{
		BaseEvent: baseEvent,
		Data:      data,
	}
}

// IsRunning returns whether the consumer is currently running
func (ec *EventConsumer) IsRunning() bool {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return ec.running
}

// GetEventFile returns the path to the events file being monitored
func (ec *EventConsumer) GetEventFile() string {
	return ec.eventFile
}

// GetCurrentOffset returns the current read offset in the events file
func (ec *EventConsumer) GetCurrentOffset() int64 {
	ec.mu.RLock()
	defer ec.mu.RUnlock()
	return ec.offset
}
