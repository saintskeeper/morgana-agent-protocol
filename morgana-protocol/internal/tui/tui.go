package tui

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// TUI represents the main TUI interface
type TUI struct {
	model   *Model
	program *tea.Program
	config  TUIConfig
}

// New creates a new TUI instance
func New(ctx context.Context, eventBus events.EventBus, config ...TUIConfig) *TUI {
	// Use default config if none provided
	var cfg TUIConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultTUIConfig()
	}

	// Create the model
	model := NewModel(ctx, eventBus, cfg)

	// Create the bubbletea program with optimized settings for performance
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	}

	program := tea.NewProgram(model, opts...)

	tui := &TUI{
		model:   model,
		program: program,
		config:  cfg,
	}

	// Set the program in the model for the event bridge
	if err := model.SetProgram(program); err != nil {
		// Log error but don't fail - TUI can still work without events
		fmt.Fprintf(os.Stderr, "Warning: Failed to connect event bridge: %v\n", err)
	}

	return tui
}

// Start begins the TUI execution
func (t *TUI) Start() error {
	// Run the bubbletea program
	_, err := t.program.Run()
	return err
}

// StartAsync starts the TUI in a separate goroutine
func (t *TUI) StartAsync() error {
	go func() {
		if _, err := t.program.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		}
	}()
	return nil
}

// Stop gracefully stops the TUI
func (t *TUI) Stop() error {
	if t.program != nil {
		t.program.Quit()
	}
	return nil
}

// Kill forcefully terminates the TUI
func (t *TUI) Kill() error {
	if t.program != nil {
		t.program.Kill()
	}
	return nil
}

// Send allows sending custom messages to the TUI
func (t *TUI) Send(msg tea.Msg) {
	if t.program != nil {
		t.program.Send(msg)
	}
}

// IsRunning returns whether the TUI is currently active
func (t *TUI) IsRunning() bool {
	// In the current bubbletea API, there's no direct way to check if running
	// This would need to be tracked manually or via program state
	return t.program != nil
}

// GetStats returns performance statistics
func (t *TUI) GetStats() TUIStats {
	if t.model == nil || t.model.bridge == nil {
		return TUIStats{}
	}

	bridgeStats := t.model.bridge.GetStats()

	return TUIStats{
		RenderCount:     t.model.renderCounter,
		EventsProcessed: bridgeStats.EventsProcessed,
		FPS:             bridgeStats.FPS,
		IsRunning:       bridgeStats.IsRunning,
		Uptime:          t.model.state.StatusInfo.Uptime,
		MemoryMB:        t.model.state.StatusInfo.SystemLoad.MemoryMB,
	}
}

// TUIStats contains statistics about TUI performance
type TUIStats struct {
	RenderCount     int64         `json:"render_count"`
	EventsProcessed int64         `json:"events_processed"`
	FPS             float64       `json:"fps"`
	IsRunning       bool          `json:"is_running"`
	Uptime          time.Duration `json:"uptime"`
	MemoryMB        float64       `json:"memory_mb"`
}

// TUIManager manages multiple TUI instances or provides utilities
type TUIManager struct {
	instances map[string]*TUI
}

// NewTUIManager creates a new TUI manager
func NewTUIManager() *TUIManager {
	return &TUIManager{
		instances: make(map[string]*TUI),
	}
}

// Create creates a new TUI instance with the given ID
func (tm *TUIManager) Create(id string, ctx context.Context, eventBus events.EventBus, config ...TUIConfig) (*TUI, error) {
	if _, exists := tm.instances[id]; exists {
		return nil, fmt.Errorf("TUI instance with ID %s already exists", id)
	}

	tui := New(ctx, eventBus, config...)
	tm.instances[id] = tui

	return tui, nil
}

// Get retrieves a TUI instance by ID
func (tm *TUIManager) Get(id string) (*TUI, bool) {
	tui, exists := tm.instances[id]
	return tui, exists
}

// Remove stops and removes a TUI instance
func (tm *TUIManager) Remove(id string) error {
	if tui, exists := tm.instances[id]; exists {
		if err := tui.Stop(); err != nil {
			return fmt.Errorf("failed to stop TUI %s: %w", id, err)
		}
		delete(tm.instances, id)
	}
	return nil
}

// StopAll stops all TUI instances
func (tm *TUIManager) StopAll() error {
	for id, tui := range tm.instances {
		if err := tui.Stop(); err != nil {
			return fmt.Errorf("failed to stop TUI %s: %w", id, err)
		}
	}
	tm.instances = make(map[string]*TUI)
	return nil
}

// List returns all TUI instance IDs
func (tm *TUIManager) List() []string {
	var ids []string
	for id := range tm.instances {
		ids = append(ids, id)
	}
	return ids
}

// Utility Functions

// ValidateConfig validates a TUI configuration
func ValidateConfig(config TUIConfig) error {
	if config.RefreshRate <= 0 {
		return fmt.Errorf("refresh rate must be positive, got %v", config.RefreshRate)
	}

	if config.RefreshRate < time.Millisecond {
		return fmt.Errorf("refresh rate too high, minimum is 1ms, got %v", config.RefreshRate)
	}

	if config.EventBufferSize <= 0 {
		return fmt.Errorf("event buffer size must be positive, got %d", config.EventBufferSize)
	}

	if config.MaxLogLines <= 0 {
		return fmt.Errorf("max log lines must be positive, got %d", config.MaxLogLines)
	}

	return nil
}

// CreateOptimizedConfig creates a TUI config optimized for performance
func CreateOptimizedConfig() TUIConfig {
	config := DefaultTUIConfig()

	// Optimize for performance
	config.RefreshRate = 16 * time.Millisecond // 60 FPS
	config.EventBufferSize = 2000              // Larger buffer for high event volume
	config.MaxLogLines = 5000                  // Reasonable memory usage
	config.ShowDebugInfo = false               // Disable debug overhead

	return config
}

// CreateHighPerformanceConfig creates a TUI config for maximum performance
func CreateHighPerformanceConfig() TUIConfig {
	config := DefaultTUIConfig()

	// Maximum performance settings
	config.RefreshRate = 33 * time.Millisecond // 30 FPS for lower CPU usage
	config.EventBufferSize = 5000              // Very large buffer
	config.MaxLogLines = 2000                  // Reduced memory usage
	config.ShowDebugInfo = false               // No debug overhead
	config.ShowTimeStamps = false              // Less rendering work
	config.EnableFiltering = false             // Disable expensive filtering
	config.EnableSearch = false                // Disable search overhead

	return config
}

// CreateDevelopmentConfig creates a TUI config optimized for development
func CreateDevelopmentConfig() TUIConfig {
	config := DefaultTUIConfig()

	// Development-friendly settings
	config.ShowDebugInfo = true   // Show performance info
	config.EnableFiltering = true // Enable all features
	config.EnableSearch = true    // Enable all features
	config.EnableExport = true    // Enable export for debugging

	return config
}

// Terminal Detection and Setup

// IsTerminalSupported checks if the current terminal supports TUI features
func IsTerminalSupported() bool {
	// Check if we're in a TTY
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return false
	}

	// Check terminal capabilities
	term := os.Getenv("TERM")
	if term == "" {
		return false
	}

	// Basic terminal support check
	unsupportedTerms := []string{"dumb", "unknown", "linux"}
	for _, unsupported := range unsupportedTerms {
		if term == unsupported {
			return false
		}
	}

	return true
}

// GetTerminalSize returns the current terminal dimensions
func GetTerminalSize() (int, int, error) {
	// This would use a terminal detection library
	// For now, return default values
	return 80, 24, nil
}

// SetupTerminal configures the terminal for optimal TUI display
func SetupTerminal() error {
	// This would configure terminal settings like:
	// - Raw mode for immediate key input
	// - Disable cursor
	// - Clear screen
	// - Set up color support

	// For now, this is handled by bubbletea automatically
	return nil
}

// RestoreTerminal restores the terminal to its original state
func RestoreTerminal() error {
	// This would restore original terminal settings
	// For now, this is handled by bubbletea automatically
	return nil
}

// Integration helpers for common use cases

// RunWithEventBus is a convenience function to run TUI with an event bus
func RunWithEventBus(ctx context.Context, eventBus events.EventBus, config ...TUIConfig) error {
	if !IsTerminalSupported() {
		return fmt.Errorf("terminal does not support TUI mode")
	}

	tui := New(ctx, eventBus, config...)
	return tui.Start()
}

// RunAsync is a convenience function to run TUI asynchronously
func RunAsync(ctx context.Context, eventBus events.EventBus, config ...TUIConfig) (*TUI, error) {
	if !IsTerminalSupported() {
		return nil, fmt.Errorf("terminal does not support TUI mode")
	}

	tui := New(ctx, eventBus, config...)
	if err := tui.StartAsync(); err != nil {
		return nil, err
	}

	return tui, nil
}
