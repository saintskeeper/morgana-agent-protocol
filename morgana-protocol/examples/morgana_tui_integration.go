package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/adapter"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/config"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/orchestrator"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/prompt"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/telemetry"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/pkg/task"
)

func main() {
	var (
		configFile     = flag.String("config", "", "Path to config file")
		enableTUI      = flag.Bool("tui", false, "Enable TUI mode")
		tuiMode        = flag.String("tui-mode", "dev", "TUI mode: dev, optimized, or high-performance")
		agentDir       = flag.String("agent-dir", os.ExpandEnv("$HOME/.claude/agents"), "Agent prompt directory")
		otelExporter   = flag.String("otel-exporter", "stdout", "OpenTelemetry exporter: stdout, otlp, none")
		otelEndpoint   = flag.String("otel-endpoint", "localhost:4317", "OpenTelemetry collector endpoint")
		maxConcurrency = flag.Int("max-concurrency", 5, "Maximum concurrent tasks")
	)

	flag.Parse()

	if *enableTUI && !tui.IsTerminalSupported() {
		log.Fatalf("Terminal does not support TUI mode")
	}

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Load configuration
	var cfg *config.Config
	if *configFile != "" {
		loadedCfg, err := config.LoadFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		cfg = loadedCfg
	} else {
		cfg = config.DefaultConfig()
	}

	// Apply CLI overrides
	if *agentDir != os.ExpandEnv("$HOME/.claude/agents") {
		cfg.Agents.PromptDir = *agentDir
	}
	if *maxConcurrency != 5 {
		cfg.Execution.MaxConcurrency = *maxConcurrency
	}
	if *otelExporter != "stdout" {
		cfg.Telemetry.Exporter = *otelExporter
	}
	if *otelEndpoint != "localhost:4317" {
		cfg.Telemetry.OTLPEndpoint = *otelEndpoint
	}

	// Initialize telemetry
	telCfg := telemetry.Config{
		ServiceName:    cfg.Telemetry.ServiceName,
		ServiceVersion: "dev",
		Environment:    cfg.Telemetry.Environment,
		ExporterType:   cfg.Telemetry.Exporter,
		OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
		Debug:          os.Getenv("MORGANA_DEBUG") == "true",
	}

	telProvider, err := telemetry.NewProvider(ctx, telCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize telemetry: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := telProvider.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to shutdown telemetry: %v\n", err)
		}
	}()

	// Create event bus for TUI integration
	eventBusConfig := events.DefaultBusConfig()
	eventBusConfig.Debug = os.Getenv("MORGANA_DEBUG") == "true"
	eventBus := events.NewEventBus(eventBusConfig)
	defer eventBus.Close()

	// Initialize components with event bus integration
	promptLoader := prompt.NewPromptLoader(cfg.Agents.PromptDir)
	taskClient := task.NewClientWithConfig(
		cfg.TaskClient.BridgePath,
		cfg.TaskClient.PythonPath,
		cfg.TaskClient.MockMode,
	)

	adptr := adapter.New(promptLoader, taskClient, telProvider.Tracer)
	adptr.SetTimeouts(cfg.Agents.DefaultTimeout, cfg.Agents.Timeouts)
	adptr.SetEventBus(eventBus) // Connect adapter to event bus

	orch := orchestrator.New(adptr, cfg.Execution.MaxConcurrency, telProvider.Tracer)
	orch.SetEventBus(eventBus) // Connect orchestrator to event bus

	// Start TUI if enabled
	var tuiInstance *tui.TUI
	if *enableTUI {
		var tuiConfig tui.TUIConfig

		switch *tuiMode {
		case "dev":
			tuiConfig = tui.CreateDevelopmentConfig()
		case "optimized":
			tuiConfig = tui.CreateOptimizedConfig()
		case "high-performance":
			tuiConfig = tui.CreateHighPerformanceConfig()
		default:
			tuiConfig = tui.DefaultTUIConfig()
		}

		// Validate TUI config
		if err := tui.ValidateConfig(tuiConfig); err != nil {
			log.Fatalf("Invalid TUI config: %v", err)
		}

		// Start TUI asynchronously
		tuiInstance, err = tui.RunAsync(ctx, eventBus, tuiConfig)
		if err != nil {
			log.Fatalf("Failed to start TUI: %v", err)
		}

		log.Println("TUI started. Press 'q' in TUI or Ctrl+C to quit.")
	}

	// Parse input tasks
	tasks, err := parseInput()
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	// If no tasks provided, run in interactive mode
	if len(tasks) == 0 {
		log.Println("No tasks provided. Running in interactive mode...")
		if !*enableTUI {
			log.Println("Use --tui flag to enable interactive TUI mode")
			return
		}

		// Create demo tasks for TUI demonstration
		tasks = createDemoTasks()
	}

	// Start root span
	ctx, rootSpan := telProvider.Tracer.Start(ctx, "morgana.execute")
	defer rootSpan.End()

	// Execute tasks
	var results []adapter.Result
	isParallel := cfg.Execution.DefaultMode == "parallel"

	if isParallel && len(tasks) > 1 {
		results = orch.RunParallel(ctx, tasks)
	} else {
		results = orch.RunSequential(ctx, tasks)
	}

	// If not in TUI mode, output results
	if !*enableTUI {
		if err := outputResults(results); err != nil {
			log.Fatalf("Failed to output results: %v", err)
		}
	} else {
		// In TUI mode, wait for user to quit
		log.Println("Tasks completed. TUI is still running...")
		<-ctx.Done()

		// Stop TUI gracefully
		if tuiInstance != nil {
			if err := tuiInstance.Stop(); err != nil {
				log.Printf("Error stopping TUI: %v", err)
			}
		}

		// Display final stats
		if tuiInstance != nil {
			stats := tuiInstance.GetStats()
			fmt.Printf("\nTUI Session Summary:\n")
			fmt.Printf("  Renders: %d\n", stats.RenderCount)
			fmt.Printf("  Events Processed: %d\n", stats.EventsProcessed)
			fmt.Printf("  Average FPS: %.1f\n", stats.FPS)
			fmt.Printf("  Memory Usage: %.1f MB\n", stats.MemoryMB)
			fmt.Printf("  Session Duration: %v\n", stats.Uptime.Round(time.Second))
		}
	}
}

// parseInput parses command line arguments for tasks
func parseInput() ([]adapter.Task, error) {
	args := flag.Args()
	if len(args) == 0 {
		return nil, nil
	}

	// Simple task parsing - in real implementation this would be more sophisticated
	var tasks []adapter.Task
	if len(args) >= 2 && args[0] == "--agent" {
		task := adapter.Task{
			AgentType: args[1],
			Prompt:    "Interactive task execution",
			Options:   make(map[string]interface{}),
		}
		if len(args) >= 4 && args[2] == "--prompt" {
			task.Prompt = args[3]
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// createDemoTasks creates sample tasks for TUI demonstration
func createDemoTasks() []adapter.Task {
	return []adapter.Task{
		{
			AgentType: "code-implementer",
			Prompt:    "Implement a simple HTTP server with health check endpoint",
			Options:   map[string]interface{}{"complexity": "medium"},
		},
		{
			AgentType: "test-specialist",
			Prompt:    "Generate comprehensive tests for the HTTP server implementation",
			Options:   map[string]interface{}{"coverage": "high"},
		},
		{
			AgentType: "validation-expert",
			Prompt:    "Review and validate the HTTP server code for security issues",
			Options:   map[string]interface{}{"focus": "security"},
		},
	}
}

// outputResults outputs task results in JSON format
func outputResults(results []adapter.Result) error {
	// In a real implementation, this would format results properly
	fmt.Printf("Completed %d tasks\n", len(results))
	for i, result := range results {
		success := result.Error == ""
		fmt.Printf("Task %d: Success: %v, Output Length: %d\n", i+1, success, len(result.Output))
		if result.Error != "" {
			fmt.Printf("  Error: %s\n", result.Error)
		}
	}
	return nil
}
