package adapter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/prompt"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/telemetry"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/pkg/task"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Task represents an agent task request
type Task struct {
	AgentType  string                 `json:"agent_type"`
	Prompt     string                 `json:"prompt"`
	Options    map[string]interface{} `json:"options,omitempty"`
	RetryCount int                    `json:"retry_count,omitempty"`
	ModelHint  string                 `json:"model_hint,omitempty"`
	Complexity string                 `json:"complexity,omitempty"`
}

// Result represents the output from an agent task
type Result struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// Adapter bridges specialized agent types with the general-purpose Task tool
type Adapter struct {
	promptLoader   *prompt.Loader
	taskClient     *task.Client
	tracer         trace.Tracer
	defaultTimeout time.Duration
	timeouts       map[string]time.Duration
	modelSelector  *ModelSelector
	mu             sync.RWMutex
}

// New creates a new Adapter instance
func New(promptLoader *prompt.Loader, taskClient *task.Client, tracer trace.Tracer) *Adapter {
	return &Adapter{
		promptLoader:   promptLoader,
		taskClient:     taskClient,
		tracer:         tracer,
		defaultTimeout: 2 * time.Minute,
		timeouts:       make(map[string]time.Duration),
		modelSelector:  NewModelSelector(),
	}
}

// SetTimeouts configures timeouts from config
func (a *Adapter) SetTimeouts(defaultTimeout time.Duration, agentTimeouts map[string]time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.defaultTimeout = defaultTimeout
	a.timeouts = agentTimeouts
}

// getTimeout returns the timeout for a specific agent type
func (a *Adapter) getTimeout(agentType string) time.Duration {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if timeout, ok := a.timeouts[agentType]; ok {
		return timeout
	}
	return a.defaultTimeout
}

// Execute runs a task with the appropriate agent type
func (a *Adapter) Execute(ctx context.Context, task Task) Result {
	// Create timeout context
	timeout := a.getTimeout(task.AgentType)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Start main span
	attrs := telemetry.AgentAttributes(task.AgentType, fmt.Sprintf("%p", &task))
	attrs = append(attrs, attribute.String("timeout", timeout.String()))
	if task.RetryCount > 0 {
		attrs = append(attrs, attribute.Int("retry_count", task.RetryCount))
	}
	if task.ModelHint != "" {
		attrs = append(attrs, attribute.String("model_hint", task.ModelHint))
	}
	if task.Complexity != "" {
		attrs = append(attrs, attribute.String("complexity", task.Complexity))
	}
	ctx, span := a.tracer.Start(ctx, "agent.execute",
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	startTime := time.Now()
	// Validate agent type
	_, validateSpan := a.tracer.Start(ctx, "agent.validate")
	validAgents := []string{"code-implementer", "sprint-planner", "test-specialist", "validation-expert"}
	if !contains(validAgents, task.AgentType) {
		err := fmt.Errorf("unknown agent type: %s. Available types: %v", task.AgentType, validAgents)
		validateSpan.RecordError(err)
		validateSpan.SetStatus(codes.Error, "Invalid agent type")
		validateSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "Validation failed")
		return Result{
			Error: err.Error(),
		}
	}
	validateSpan.SetAttributes(attribute.Bool("validation.passed", true))
	validateSpan.End()

	// Load agent prompt
	ctx, loadSpan := a.tracer.Start(ctx, "agent.load_prompt")
	agentPrompt, err := a.promptLoader.Load(task.AgentType)
	if err != nil {
		loadSpan.RecordError(err)
		loadSpan.SetStatus(codes.Error, "Failed to load prompt")
		loadSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "Prompt loading failed")
		return Result{
			Error: fmt.Sprintf("loading agent prompt: %v", err),
		}
	}
	loadSpan.SetAttributes(
		attribute.Int("prompt.system_length", len(agentPrompt)),
	)
	loadSpan.End()

	// Select appropriate model based on task context
	selectedModel := a.modelSelector.SelectModel(task)
	modelCapabilities := a.modelSelector.GetModelCapabilities(selectedModel)

	// Add model information to options for Task tool
	options := task.Options
	if options == nil {
		options = make(map[string]interface{})
	}
	options["model"] = selectedModel
	options["model_capabilities"] = modelCapabilities

	// Combine agent system prompt with task prompt
	fullPrompt := fmt.Sprintf("%s\n\nTask: %s", agentPrompt, task.Prompt)
	span.SetAttributes(
		telemetry.PromptAttributes(len(fullPrompt), false)...,
	)
	span.SetAttributes(
		attribute.String("model.selected", selectedModel),
		attribute.Bool("model.token_efficient", modelCapabilities["token_efficient"].(bool)),
		attribute.String("model.reasoning_level", modelCapabilities["reasoning_level"].(string)),
		attribute.String("model.cost_tier", modelCapabilities["cost_tier"].(string)),
	)

	// Execute via Task tool with general-purpose type
	ctx, execSpan := a.tracer.Start(ctx, "agent.task_execution",
		trace.WithAttributes(
			attribute.String("task.type", "general-purpose"),
			attribute.String("task.model", selectedModel),
		),
	)
	result, err := a.taskClient.RunWithContext(ctx, "general-purpose", fullPrompt, options)
	execTime := time.Since(startTime)
	execSpan.SetAttributes(
		attribute.Int64("execution.duration_ms", execTime.Milliseconds()),
	)

	if err != nil {
		execSpan.RecordError(err)
		execSpan.SetStatus(codes.Error, "Task execution failed")
		execSpan.End()
		span.RecordError(err)
		span.SetStatus(codes.Error, "Execution failed")
		span.SetAttributes(
			telemetry.ResultAttributes(false, 0, execTime.Milliseconds())...,
		)
		return Result{
			Error: fmt.Sprintf("executing task: %v", err),
		}
	}
	execSpan.SetStatus(codes.Ok, "Task executed successfully")
	execSpan.End()

	// Record success
	span.SetAttributes(
		telemetry.ResultAttributes(true, len(result.Output), execTime.Milliseconds())...,
	)
	span.SetStatus(codes.Ok, "Agent execution completed")

	return Result{
		Output: result.Output,
	}
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
