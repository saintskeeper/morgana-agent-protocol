package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// IPCServer handles Unix domain socket communication for event forwarding
type IPCServer struct {
	socketPath string
	eventBus   events.EventBus
	listener   net.Listener
	mu         sync.Mutex
	clients    map[net.Conn]bool
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewIPCServer creates a new IPC server instance
func NewIPCServer(socketPath string, eventBus events.EventBus) *IPCServer {
	return &IPCServer{
		socketPath: socketPath,
		eventBus:   eventBus,
		clients:    make(map[net.Conn]bool),
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

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			var msg IPCMessage
			if err := decoder.Decode(&msg); err != nil {
				// Connection closed or decode error
				if err.Error() != "EOF" {
					log.Printf("Error decoding message from client: %v", err)
				}
				return
			}

			// Reconstruct and publish the event to the monitor's event bus
			event := s.reconstructEvent(msg)
			if event != nil {
				s.eventBus.PublishAsync(event)
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
