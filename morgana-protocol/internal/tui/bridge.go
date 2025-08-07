package tui

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// EventBridge connects the event system to bubbletea
type EventBridge struct {
	eventBus      events.EventBus
	program       *tea.Program
	config        TUIConfig
	subscriptions []string
	stopped       bool
	mu            sync.RWMutex

	// Performance metrics
	eventsProcessed int64
	lastFPSUpdate   time.Time
	fpsCounter      int64
}

// NewEventBridge creates a new bridge between events and tea
func NewEventBridge(eventBus events.EventBus, config TUIConfig) *EventBridge {
	return &EventBridge{
		eventBus:      eventBus,
		config:        config,
		subscriptions: make([]string, 0, 10),
		lastFPSUpdate: time.Now(),
	}
}

// Start begins the bridge operation and connects to the event bus
func (b *EventBridge) Start(program *tea.Program) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.stopped {
		return nil
	}

	b.program = program

	// Subscribe to all events to forward them to the TUI
	subID := b.eventBus.SubscribeAll(b.handleEvent)
	b.subscriptions = append(b.subscriptions, subID)

	// Start the ticker for periodic updates and FPS control
	go b.startTicker()

	return nil
}

// Stop shuts down the bridge and cleans up subscriptions
func (b *EventBridge) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.stopped {
		return nil
	}

	b.stopped = true

	// Unsubscribe from all events
	for _, subID := range b.subscriptions {
		b.eventBus.Unsubscribe(subID)
	}
	b.subscriptions = nil

	return nil
}

// handleEvent processes events from the event bus and sends them to bubbletea
func (b *EventBridge) handleEvent(event events.Event) {
	b.mu.RLock()
	program := b.program
	stopped := b.stopped
	b.mu.RUnlock()

	if stopped || program == nil {
		return
	}

	// Convert event to tea message
	msg := EventMessage{
		Event:     event,
		Timestamp: time.Now(),
	}

	// Send to bubbletea program
	go func() {
		program.Send(msg)
		b.eventsProcessed++
	}()
}

// startTicker starts the periodic update ticker for FPS control
func (b *EventBridge) startTicker() {
	ticker := time.NewTicker(b.config.RefreshRate)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			b.mu.RLock()
			program := b.program
			stopped := b.stopped
			b.mu.RUnlock()

			if stopped || program == nil {
				return
			}

			// Send periodic tick message for smooth updates
			go func() {
				program.Send(TickMessage(t))
				b.updateFPS()
			}()
		}
	}
}

// updateFPS calculates and tracks rendering FPS
func (b *EventBridge) updateFPS() {
	b.fpsCounter++
	now := time.Now()

	// Update FPS calculation every second
	if now.Sub(b.lastFPSUpdate) >= time.Second {
		b.lastFPSUpdate = now
		b.fpsCounter = 0
	}
}

// GetStats returns performance statistics for the bridge
func (b *EventBridge) GetStats() BridgeStats {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return BridgeStats{
		EventsProcessed: b.eventsProcessed,
		FPS:             float64(b.fpsCounter),
		IsRunning:       !b.stopped,
		Subscriptions:   len(b.subscriptions),
	}
}

// BridgeStats contains statistics about the event bridge
type BridgeStats struct {
	EventsProcessed int64
	FPS             float64
	IsRunning       bool
	Subscriptions   int
}

// EventProcessor handles specific event processing logic
type EventProcessor struct {
	config TUIConfig
}

// NewEventProcessor creates a new event processor
func NewEventProcessor(config TUIConfig) *EventProcessor {
	return &EventProcessor{
		config: config,
	}
}

// ProcessEvent converts events to appropriate TUI state updates
func (p *EventProcessor) ProcessEvent(event events.Event, state *ModelState) {
	switch e := event.(type) {
	case *events.TaskStartedEvent:
		p.processTaskStarted(e, state)
	case *events.TaskProgressEvent:
		p.processTaskProgress(e, state)
	case *events.TaskCompletedEvent:
		p.processTaskCompleted(e, state)
	case *events.TaskFailedEvent:
		p.processTaskFailed(e, state)
	case *events.OrchestratorStartedEvent:
		p.processOrchestratorStarted(e, state)
	case *events.OrchestratorCompletedEvent:
		p.processOrchestratorCompleted(e, state)
	case *events.OrchestratorFailedEvent:
		p.processOrchestratorFailed(e, state)
	case *events.AdapterValidationEvent:
		p.processAdapterValidation(e, state)
	case *events.AdapterPromptLoadEvent:
		p.processAdapterPromptLoad(e, state)
	case *events.AdapterExecutionEvent:
		p.processAdapterExecution(e, state)
	}

	// Add log entry for all events
	p.addLogEntry(event, state)

	// Update event counter
	state.EventCount++
}

// processTaskStarted handles task started events
func (p *EventProcessor) processTaskStarted(event *events.TaskStartedEvent, state *ModelState) {
	if state.TaskStates == nil {
		state.TaskStates = make(map[string]*TaskState)
	}

	state.TaskStates[event.TaskID()] = &TaskState{
		ID:         event.TaskID(),
		AgentType:  event.AgentType,
		Status:     StatusRunning,
		Progress:   0.0,
		Stage:      "starting",
		Message:    "Task started",
		StartTime:  event.Timestamp(),
		RetryCount: event.RetryCount,
	}

	if state.StatusInfo != nil {
		state.StatusInfo.ActiveTasks++
	}
}

// processTaskProgress handles task progress events
func (p *EventProcessor) processTaskProgress(event *events.TaskProgressEvent, state *ModelState) {
	if task, exists := state.TaskStates[event.TaskID()]; exists {
		task.Status = StatusRunning
		task.Progress = event.Progress
		task.Stage = event.Stage
		task.Message = event.Message
		task.Duration = event.Duration
	}
}

// processTaskCompleted handles task completed events
func (p *EventProcessor) processTaskCompleted(event *events.TaskCompletedEvent, state *ModelState) {
	if task, exists := state.TaskStates[event.TaskID()]; exists {
		task.Status = StatusCompleted
		task.Progress = 1.0
		task.Stage = "completed"
		task.Message = "Task completed successfully"
		task.Duration = event.Duration
		task.Model = event.Model
		task.Output = event.Output
	}

	if state.StatusInfo != nil {
		state.StatusInfo.ActiveTasks--
		state.StatusInfo.CompletedTasks++
	}
}

// processTaskFailed handles task failed events
func (p *EventProcessor) processTaskFailed(event *events.TaskFailedEvent, state *ModelState) {
	if task, exists := state.TaskStates[event.TaskID()]; exists {
		task.Status = StatusFailed
		task.Stage = "failed"
		task.Message = "Task failed"
		task.Duration = event.Duration
		task.Error = event.Error
		task.RetryCount = event.RetryCount
	}

	if state.StatusInfo != nil {
		state.StatusInfo.ActiveTasks--
		state.StatusInfo.FailedTasks++
	}
}

// processOrchestratorStarted handles orchestrator started events
func (p *EventProcessor) processOrchestratorStarted(event *events.OrchestratorStartedEvent, state *ModelState) {
	// Update system status with orchestrator info
	if state.StatusInfo == nil {
		state.StatusInfo = &StatusInfo{}
	}
}

// processOrchestratorCompleted handles orchestrator completed events
func (p *EventProcessor) processOrchestratorCompleted(event *events.OrchestratorCompletedEvent, state *ModelState) {
	// Update final orchestrator status
}

// processOrchestratorFailed handles orchestrator failed events
func (p *EventProcessor) processOrchestratorFailed(event *events.OrchestratorFailedEvent, state *ModelState) {
	// Update orchestrator failure status
}

// processAdapterValidation handles adapter validation events
func (p *EventProcessor) processAdapterValidation(event *events.AdapterValidationEvent, state *ModelState) {
	// Update validation status in relevant task
	if task, exists := state.TaskStates[event.TaskID()]; exists {
		if event.Valid {
			task.Stage = "validated"
			task.Message = "Agent validated"
		} else {
			task.Stage = "validation_failed"
			task.Message = "Validation failed: " + event.Error
		}
	}
}

// processAdapterPromptLoad handles adapter prompt load events
func (p *EventProcessor) processAdapterPromptLoad(event *events.AdapterPromptLoadEvent, state *ModelState) {
	// Update prompt load status
	if task, exists := state.TaskStates[event.TaskID()]; exists {
		if event.Success {
			task.Stage = "prompt_loaded"
			task.Message = "Prompt loaded successfully"
		} else {
			task.Stage = "prompt_load_failed"
			task.Message = "Prompt load failed: " + event.Error
		}
	}
}

// processAdapterExecution handles adapter execution events
func (p *EventProcessor) processAdapterExecution(event *events.AdapterExecutionEvent, state *ModelState) {
	// Update execution phase
	if task, exists := state.TaskStates[event.TaskID()]; exists {
		task.Stage = event.Phase
		task.Duration = event.Duration
		if event.Model != "" {
			task.Model = event.Model
		}
		if !event.Success && event.Error != "" {
			task.Error = event.Error
		}
	}
}

// addLogEntry adds a log entry for the event
func (p *EventProcessor) addLogEntry(event events.Event, state *ModelState) {
	level := p.getLogLevel(event)
	message := p.formatEventMessage(event)
	agentType := p.extractAgentType(event)
	stage := p.extractStage(event)

	entry := &LogEntry{
		ID:        event.TaskID(),
		Timestamp: event.Timestamp(),
		Level:     level,
		TaskID:    event.TaskID(),
		AgentType: agentType,
		Stage:     stage,
		Message:   message,
		Event:     event,
	}

	// Add to log entries with circular buffer behavior
	state.LogEntries = append(state.LogEntries, entry)

	// Keep only the most recent entries for performance
	if len(state.LogEntries) > p.config.MaxLogLines {
		// Remove oldest entries
		copy(state.LogEntries, state.LogEntries[len(state.LogEntries)-p.config.MaxLogLines:])
		state.LogEntries = state.LogEntries[:p.config.MaxLogLines]
	}
}

// getLogLevel determines the log level for an event
func (p *EventProcessor) getLogLevel(event events.Event) LogLevel {
	switch event.Type() {
	case events.EventTaskFailed, events.EventOrchestratorFailed:
		return LogLevelError
	case events.EventTaskStarted, events.EventOrchestratorStarted:
		return LogLevelInfo
	case events.EventTaskProgress:
		return LogLevelDebug
	case events.EventAdapterValidation:
		if e, ok := event.(*events.AdapterValidationEvent); ok && !e.Valid {
			return LogLevelError
		}
		return LogLevelInfo
	default:
		return LogLevelInfo
	}
}

// formatEventMessage creates a human-readable message for the event
func (p *EventProcessor) formatEventMessage(event events.Event) string {
	switch e := event.(type) {
	case *events.TaskStartedEvent:
		return "Task started: " + e.AgentType
	case *events.TaskProgressEvent:
		return e.Message
	case *events.TaskCompletedEvent:
		return "Task completed successfully"
	case *events.TaskFailedEvent:
		return "Task failed: " + e.Error
	case *events.OrchestratorStartedEvent:
		return "Orchestrator started"
	case *events.OrchestratorCompletedEvent:
		return "Orchestrator completed"
	case *events.OrchestratorFailedEvent:
		return "Orchestrator failed: " + e.Error
	case *events.AdapterValidationEvent:
		if e.Valid {
			return "Agent validation successful"
		}
		return "Agent validation failed: " + e.Error
	case *events.AdapterPromptLoadEvent:
		if e.Success {
			return "Prompt loaded successfully"
		}
		return "Prompt load failed: " + e.Error
	case *events.AdapterExecutionEvent:
		return "Execution phase: " + e.Phase
	default:
		return string(event.Type())
	}
}

// extractAgentType extracts agent type from event if available
func (p *EventProcessor) extractAgentType(event events.Event) string {
	switch e := event.(type) {
	case *events.TaskStartedEvent:
		return e.AgentType
	case *events.TaskProgressEvent:
		return e.AgentType
	case *events.TaskCompletedEvent:
		return e.AgentType
	case *events.TaskFailedEvent:
		return e.AgentType
	case *events.AdapterValidationEvent:
		return e.AgentType
	case *events.AdapterPromptLoadEvent:
		return e.AgentType
	case *events.AdapterExecutionEvent:
		return e.AgentType
	default:
		return ""
	}
}

// extractStage extracts stage/phase from event if available
func (p *EventProcessor) extractStage(event events.Event) string {
	switch e := event.(type) {
	case *events.TaskProgressEvent:
		return e.Stage
	case *events.TaskFailedEvent:
		return e.Stage
	case *events.AdapterExecutionEvent:
		return e.Phase
	default:
		return string(event.Type())
	}
}
