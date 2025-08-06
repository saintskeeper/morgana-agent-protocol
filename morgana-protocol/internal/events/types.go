package events

import (
	"context"
	"time"
)

// EventType represents the type of event
type EventType string

const (
	// Task lifecycle events
	EventTaskStarted   EventType = "task.started"
	EventTaskProgress  EventType = "task.progress"
	EventTaskCompleted EventType = "task.completed"
	EventTaskFailed    EventType = "task.failed"

	// Orchestrator events
	EventOrchestratorStarted   EventType = "orchestrator.started"
	EventOrchestratorCompleted EventType = "orchestrator.completed"
	EventOrchestratorFailed    EventType = "orchestrator.failed"

	// Adapter events
	EventAdapterValidation EventType = "adapter.validation"
	EventAdapterPromptLoad EventType = "adapter.prompt_load"
	EventAdapterExecution  EventType = "adapter.execution"
)

// Event is the base interface for all events
type Event interface {
	Type() EventType
	Timestamp() time.Time
	TaskID() string
	Context() context.Context
}

// BaseEvent provides common event fields
type BaseEvent struct {
	EventType EventType       `json:"event_type"`
	Time      time.Time       `json:"timestamp"`
	ID        string          `json:"task_id"`
	Ctx       context.Context `json:"-"`
}

func (e BaseEvent) Type() EventType          { return e.EventType }
func (e BaseEvent) Timestamp() time.Time     { return e.Time }
func (e BaseEvent) TaskID() string           { return e.ID }
func (e BaseEvent) Context() context.Context { return e.Ctx }

// TaskStartedEvent is emitted when a task begins execution
type TaskStartedEvent struct {
	BaseEvent
	AgentType  string                 `json:"agent_type"`
	Prompt     string                 `json:"prompt,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
	RetryCount int                    `json:"retry_count,omitempty"`
	ModelHint  string                 `json:"model_hint,omitempty"`
	Complexity string                 `json:"complexity,omitempty"`
	Timeout    time.Duration          `json:"timeout"`
}

// TaskProgressEvent is emitted during task execution to show progress
type TaskProgressEvent struct {
	BaseEvent
	AgentType string        `json:"agent_type"`
	Stage     string        `json:"stage"`    // validation, prompt_load, execution
	Message   string        `json:"message"`  // Human-readable progress message
	Progress  float64       `json:"progress"` // 0.0 - 1.0
	Duration  time.Duration `json:"duration"` // Time elapsed since task start
}

// TaskCompletedEvent is emitted when a task completes successfully
type TaskCompletedEvent struct {
	BaseEvent
	AgentType    string        `json:"agent_type"`
	Output       string        `json:"output,omitempty"`
	OutputLength int           `json:"output_length"`
	Duration     time.Duration `json:"duration"`
	Model        string        `json:"model"`
}

// TaskFailedEvent is emitted when a task fails
type TaskFailedEvent struct {
	BaseEvent
	AgentType  string        `json:"agent_type"`
	Error      string        `json:"error"`
	Duration   time.Duration `json:"duration"`
	Stage      string        `json:"stage"` // Where the failure occurred
	RetryCount int           `json:"retry_count"`
}

// OrchestratorStartedEvent is emitted when orchestrator begins execution
type OrchestratorStartedEvent struct {
	BaseEvent
	Mode           string `json:"mode"` // sequential, parallel
	TaskCount      int    `json:"task_count"`
	MaxConcurrency int    `json:"max_concurrency,omitempty"`
}

// OrchestratorCompletedEvent is emitted when orchestrator completes
type OrchestratorCompletedEvent struct {
	BaseEvent
	Mode         string        `json:"mode"`
	TaskCount    int           `json:"task_count"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
	Duration     time.Duration `json:"duration"`
}

// OrchestratorFailedEvent is emitted when orchestrator fails
type OrchestratorFailedEvent struct {
	BaseEvent
	Mode           string        `json:"mode"`
	TaskCount      int           `json:"task_count"`
	CompletedCount int           `json:"completed_count"`
	Error          string        `json:"error"`
	Duration       time.Duration `json:"duration"`
}

// AdapterValidationEvent is emitted during adapter validation
type AdapterValidationEvent struct {
	BaseEvent
	AgentType string `json:"agent_type"`
	Valid     bool   `json:"valid"`
	Error     string `json:"error,omitempty"`
}

// AdapterPromptLoadEvent is emitted during prompt loading
type AdapterPromptLoadEvent struct {
	BaseEvent
	AgentType    string `json:"agent_type"`
	Success      bool   `json:"success"`
	PromptLength int    `json:"prompt_length,omitempty"`
	Error        string `json:"error,omitempty"`
}

// AdapterExecutionEvent is emitted during task execution phases
type AdapterExecutionEvent struct {
	BaseEvent
	AgentType string        `json:"agent_type"`
	Phase     string        `json:"phase"` // start, model_selection, task_call, complete
	Model     string        `json:"model,omitempty"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
}

// NewTaskStartedEvent creates a new task started event
func NewTaskStartedEvent(ctx context.Context, taskID string, agentType string, prompt string, options map[string]interface{}, retryCount int, modelHint string, complexity string, timeout time.Duration) *TaskStartedEvent {
	return &TaskStartedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskStarted,
			Time:      time.Now(),
			ID:        taskID,
			Ctx:       ctx,
		},
		AgentType:  agentType,
		Prompt:     prompt,
		Options:    options,
		RetryCount: retryCount,
		ModelHint:  modelHint,
		Complexity: complexity,
		Timeout:    timeout,
	}
}

// NewTaskProgressEvent creates a new task progress event
func NewTaskProgressEvent(ctx context.Context, taskID string, agentType string, stage string, message string, progress float64, duration time.Duration) *TaskProgressEvent {
	return &TaskProgressEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskProgress,
			Time:      time.Now(),
			ID:        taskID,
			Ctx:       ctx,
		},
		AgentType: agentType,
		Stage:     stage,
		Message:   message,
		Progress:  progress,
		Duration:  duration,
	}
}

// NewTaskCompletedEvent creates a new task completed event
func NewTaskCompletedEvent(ctx context.Context, taskID string, agentType string, output string, duration time.Duration, model string) *TaskCompletedEvent {
	return &TaskCompletedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskCompleted,
			Time:      time.Now(),
			ID:        taskID,
			Ctx:       ctx,
		},
		AgentType:    agentType,
		Output:       output,
		OutputLength: len(output),
		Duration:     duration,
		Model:        model,
	}
}

// NewTaskFailedEvent creates a new task failed event
func NewTaskFailedEvent(ctx context.Context, taskID string, agentType string, err string, duration time.Duration, stage string, retryCount int) *TaskFailedEvent {
	return &TaskFailedEvent{
		BaseEvent: BaseEvent{
			EventType: EventTaskFailed,
			Time:      time.Now(),
			ID:        taskID,
			Ctx:       ctx,
		},
		AgentType:  agentType,
		Error:      err,
		Duration:   duration,
		Stage:      stage,
		RetryCount: retryCount,
	}
}
