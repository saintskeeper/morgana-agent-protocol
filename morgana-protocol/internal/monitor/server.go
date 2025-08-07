package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// IPCServer handles Unix domain socket communication for event forwarding
type IPCServer struct {
	socketPath    string
	eventBus      events.EventBus
	listener      net.Listener
	mu            sync.Mutex
	clients       map[net.Conn]bool
	ctx           context.Context
	cancel        context.CancelFunc
	eventBuffer   []IPCMessage // Circular buffer for event history
	bufferMutex   sync.RWMutex // Thread-safe access to event buffer
	maxBufferSize int          // Maximum size of the circular buffer
	bufferIndex   int          // Current position in circular buffer
	bufferCount   int          // Number of events currently in buffer
}

// NewIPCServer creates a new IPC server instance
func NewIPCServer(socketPath string, eventBus events.EventBus) *IPCServer {
	return &IPCServer{
		socketPath:    socketPath,
		eventBus:      eventBus,
		clients:       make(map[net.Conn]bool),
		eventBuffer:   make([]IPCMessage, DefaultEventBufferSize),
		maxBufferSize: DefaultEventBufferSize,
		bufferIndex:   0,
		bufferCount:   0,
	}
}

// Start begins listening for client connections on the Unix domain socket (blocking)
func (s *IPCServer) Start(ctx context.Context) error {
	if err := s.StartNonBlocking(ctx); err != nil {
		return err
	}
	// Wait for context cancellation
	<-s.ctx.Done()
	return nil
}

// StartNonBlocking begins listening for client connections without blocking
func (s *IPCServer) StartNonBlocking(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	// Remove existing socket file if it exists
	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing socket: %w", err)
	}

	// Create Unix domain socket listener
	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create Unix socket listener: %w", err)
	}

	s.listener = listener
	log.Printf("IPC server listening on %s", s.socketPath)

	// Accept connections in a separate goroutine
	go s.acceptConnections()

	return nil
}

// Stop gracefully shuts down the IPC server
func (s *IPCServer) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}

	// Close all client connections
	s.mu.Lock()
	for client := range s.clients {
		client.Close()
	}
	s.clients = make(map[net.Conn]bool)
	s.mu.Unlock()

	// Remove socket file
	if err := os.Remove(s.socketPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Error removing socket file: %v", err)
	}

	return nil
}

// acceptConnections continuously accepts new client connections
func (s *IPCServer) acceptConnections() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				// Check if we're shutting down
				select {
				case <-s.ctx.Done():
					return
				default:
					log.Printf("Error accepting connection: %v", err)
					continue
				}
			}

			// Track the client connection
			s.mu.Lock()
			s.clients[conn] = true
			s.mu.Unlock()

			// Handle client in separate goroutine
			go s.handleConnection(conn)
		}
	}
}

// handleConnection processes messages from a single client connection
func (s *IPCServer) handleConnection(conn net.Conn) {
	defer func() {
		// Remove client from tracking
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()

		// Close connection
		conn.Close()
		log.Printf("Client disconnected")
	}()

	log.Printf("New client connected")
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// Try to decode as IPCMessage first
			var rawMsg json.RawMessage
			if err := decoder.Decode(&rawMsg); err != nil {
				// Connection closed or decode error
				if err.Error() != "EOF" {
					log.Printf("Error decoding message from client: %v", err)
				}
				return
			}

			// Check if it's a request message
			var req IPCRequest
			if err := json.Unmarshal(rawMsg, &req); err == nil && req.Type == MessageTypeRequest {
				// Handle request
				if reqType, ok := req.Data["request"].(string); ok && reqType == RequestHistory {
					log.Printf("Client requested event history")
					// Send buffered events
					events := s.getBufferedEvents()
					replay := IPCReplay{
						Events: events,
						Count:  len(events),
					}

					// Wrap in a message
					replayMsg := IPCMessage{
						Type:      MessageTypeReplay,
						Timestamp: time.Now(),
						Data:      replay,
					}

					if err := encoder.Encode(replayMsg); err != nil {
						log.Printf("Error sending replay to client: %v", err)
					} else {
						log.Printf("Sent %d buffered events to client", len(events))
					}
				}
			} else {
				// Regular event message
				var msg IPCMessage
				if err := json.Unmarshal(rawMsg, &msg); err != nil {
					log.Printf("Error parsing message: %v", err)
					continue
				}

				// Add the message to the circular buffer before processing
				s.addToBuffer(msg)

				// Reconstruct and publish the event to the monitor's event bus
				event := s.reconstructEvent(msg)
				if event != nil {
					s.eventBus.PublishAsync(event)
				}
			}
		}
	}
}

// reconstructEvent converts an IPCMessage back to an Event
func (s *IPCServer) reconstructEvent(msg IPCMessage) events.Event {
	// Parse the timestamp
	timestamp := msg.Timestamp

	// Create a context (simplified for monitoring)
	ctx := context.Background()

	// Create base event
	baseEvent := events.BaseEvent{
		EventType: events.EventType(msg.Type),
		Time:      timestamp,
		ID:        msg.TaskID,
		Ctx:       ctx,
	}

	// Reconstruct specific event types based on the message type
	switch events.EventType(msg.Type) {
	case events.EventTaskStarted:
		if data, ok := msg.Data.(map[string]interface{}); ok {
			event := &events.TaskStartedEvent{
				BaseEvent: baseEvent,
			}
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
			return event
		}

	case events.EventTaskProgress:
		if data, ok := msg.Data.(map[string]interface{}); ok {
			event := &events.TaskProgressEvent{
				BaseEvent: baseEvent,
			}
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
			return event
		}

	case events.EventTaskCompleted:
		if data, ok := msg.Data.(map[string]interface{}); ok {
			event := &events.TaskCompletedEvent{
				BaseEvent: baseEvent,
			}
			if agentType, ok := data["agent_type"].(string); ok {
				event.AgentType = agentType
			}
			if model, ok := data["model"].(string); ok {
				event.Model = model
			}
			if output, ok := data["output"].(string); ok {
				event.Output = output
				event.OutputLength = len(output)
			}
			return event
		}

	case events.EventTaskFailed:
		if data, ok := msg.Data.(map[string]interface{}); ok {
			event := &events.TaskFailedEvent{
				BaseEvent: baseEvent,
			}
			if agentType, ok := data["agent_type"].(string); ok {
				event.AgentType = agentType
			}
			if errorMsg, ok := data["error"].(string); ok {
				event.Error = errorMsg
			}
			if stage, ok := data["stage"].(string); ok {
				event.Stage = stage
			}
			return event
		}
	}

	// Return a generic event if we can't reconstruct the specific type
	return &baseEvent
}

// GetClientCount returns the number of connected clients
func (s *IPCServer) GetClientCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.clients)
}

// addToBuffer adds an event to the circular buffer in a thread-safe manner
func (s *IPCServer) addToBuffer(msg IPCMessage) {
	s.bufferMutex.Lock()
	defer s.bufferMutex.Unlock()

	// Add the message to the current buffer position
	s.eventBuffer[s.bufferIndex] = msg

	// Update buffer index (circular wraparound)
	s.bufferIndex = (s.bufferIndex + 1) % s.maxBufferSize

	// Update buffer count, capped at maxBufferSize
	if s.bufferCount < s.maxBufferSize {
		s.bufferCount++
	}
}

// getBufferedEvents returns a copy of all buffered events in chronological order
func (s *IPCServer) getBufferedEvents() []IPCMessage {
	s.bufferMutex.RLock()
	defer s.bufferMutex.RUnlock()

	if s.bufferCount == 0 {
		return []IPCMessage{}
	}

	// Create a slice to hold the events in chronological order
	events := make([]IPCMessage, s.bufferCount)

	if s.bufferCount < s.maxBufferSize {
		// Buffer is not yet full, events are from 0 to bufferIndex-1
		copy(events, s.eventBuffer[:s.bufferCount])
	} else {
		// Buffer is full, need to wrap around
		// Oldest event is at bufferIndex, newest is at bufferIndex-1
		oldestIndex := s.bufferIndex

		// Copy from oldest to end of buffer
		remainingSlots := s.maxBufferSize - oldestIndex
		copy(events, s.eventBuffer[oldestIndex:])

		// Copy from beginning of buffer to current position
		if oldestIndex > 0 {
			copy(events[remainingSlots:], s.eventBuffer[:oldestIndex])
		}
	}

	return events
}

// GetBufferedEventCount returns the number of events currently in the buffer
func (s *IPCServer) GetBufferedEventCount() int {
	s.bufferMutex.RLock()
	defer s.bufferMutex.RUnlock()
	return s.bufferCount
}
