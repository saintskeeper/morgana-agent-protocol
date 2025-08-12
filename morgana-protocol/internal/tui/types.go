package tui

import (
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// TUIConfig contains configuration for the TUI
type TUIConfig struct {
	// Performance settings
	RefreshRate     time.Duration // Target refresh rate (default: 16ms for 60fps)
	EventBufferSize int           // Buffer size for event-to-tea conversion
	MaxLogLines     int           // Maximum log lines to keep in memory

	// Visual settings
	Theme          Theme // Color theme configuration
	ShowDebugInfo  bool  // Show debug information panel
	ShowTimeStamps bool  // Show timestamps in logs
	CompactMode    bool  // Use compact display mode

	// Feature flags
	EnableFiltering bool // Enable log filtering capabilities
	EnableSearch    bool // Enable log search functionality
	EnableExport    bool // Enable export functionality
}

// DefaultTUIConfig returns a default configuration optimized for performance
func DefaultTUIConfig() TUIConfig {
	return TUIConfig{
		RefreshRate:     16 * time.Millisecond, // 60fps
		EventBufferSize: 1000,
		MaxLogLines:     10000,
		Theme:           DefaultTheme(),
		ShowDebugInfo:   false,
		ShowTimeStamps:  true,
		CompactMode:     false,
		EnableFiltering: true,
		EnableSearch:    true,
		EnableExport:    false,
	}
}

// Theme contains color scheme and styling configuration
type Theme struct {
	// Base colors
	Primary   string // Primary accent color
	Secondary string // Secondary accent color
	Success   string // Success/completed color
	Warning   string // Warning color
	Error     string // Error color
	Info      string // Info color

	// UI element colors
	Background string // Background color
	Foreground string // Default text color
	Border     string // Border color
	Highlight  string // Selection/highlight color
	Muted      string // Muted/secondary text color

	// Component-specific colors
	ProgressBar struct {
		Complete   string
		Incomplete string
		Background string
	}

	// Status colors by agent type
	AgentColors map[string]string
}

// DefaultTheme returns a professional dark theme with accessibility considerations
func DefaultTheme() Theme {
	theme := Theme{
		Primary:   "#7C3AED", // Violet
		Secondary: "#06B6D4", // Cyan
		Success:   "#10B981", // Emerald
		Warning:   "#F59E0B", // Amber
		Error:     "#EF4444", // Red
		Info:      "#3B82F6", // Blue

		Background: "#0F172A", // Slate 900
		Foreground: "#F1F5F9", // Slate 100
		Border:     "#334155", // Slate 700
		Highlight:  "#1E293B", // Slate 800
		Muted:      "#64748B", // Slate 500

		AgentColors: map[string]string{
			"code-implementer":  "#10B981", // Green
			"sprint-planner":    "#3B82F6", // Blue
			"test-specialist":   "#F59E0B", // Amber
			"validation-expert": "#EF4444", // Red
			"default":           "#64748B", // Slate
		},
	}

	theme.ProgressBar.Complete = theme.Success
	theme.ProgressBar.Incomplete = theme.Muted
	theme.ProgressBar.Background = theme.Border

	return theme
}

// EventMessage wraps events to be sent as tea.Msg
type EventMessage struct {
	Event     events.Event
	Timestamp time.Time
}

// TickMessage for periodic updates (fps control)
type TickMessage time.Time

// ResizeMessage for terminal resize events
type ResizeMessage struct {
	Width  int
	Height int
}

// KeyMessage wraps keyboard input
type KeyMessage struct {
	Key  string
	Rune rune
	Alt  bool
	Ctrl bool
}

// FilterUpdateMessage for log filtering
type FilterUpdateMessage struct {
	Filter string
	Active bool
}

// SearchUpdateMessage for search functionality
type SearchUpdateMessage struct {
	Query  string
	Active bool
	Index  int
}

// ExportMessage for data export operations
type ExportMessage struct {
	Type     string // "logs", "events", "status"
	Filepath string
	Success  bool
	Error    string
}

// ModelState represents the current state of different components
type ModelState struct {
	// Layout state
	Width   int
	Height  int
	Focused ComponentType

	// Data state
	TaskStates map[string]*TaskState
	LogEntries []*LogEntry
	StatusInfo *StatusInfo

	// UI state
	CurrentFilter string
	SearchQuery   string
	SearchIndex   int
	ShowHelp      bool

	// Performance tracking
	LastRender  time.Time
	RenderCount int64
	EventCount  int64
}

// ComponentType represents different UI components
type ComponentType int

const (
	ComponentDashboard ComponentType = iota
	ComponentLogs
	ComponentStatus
	ComponentStatistics
	ComponentHelp
)

// TaskState represents the current state of a task
type TaskState struct {
	ID         string
	AgentType  string
	Status     TaskStatus
	Progress   float64
	Stage      string
	Message    string
	StartTime  time.Time
	Duration   time.Duration
	Error      string
	Model      string
	RetryCount int
	Output     string
}

// TaskStatus represents task execution status
type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusRunning
	StatusCompleted
	StatusFailed
)

// LogEntry represents a log entry in the TUI
type LogEntry struct {
	ID        string
	Timestamp time.Time
	Level     LogLevel
	TaskID    string
	AgentType string
	Stage     string
	Message   string
	Event     events.Event
}

// LogLevel represents log severity levels
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// StatusInfo represents overall system status
type StatusInfo struct {
	ActiveTasks    int
	CompletedTasks int
	FailedTasks    int
	TotalEvents    int64
	Uptime         time.Duration
	SystemLoad     SystemLoad
}

// SystemLoad represents system performance metrics
type SystemLoad struct {
	CPUPercent    float64
	MemoryMB      float64
	EventsPerSec  float64
	RenderFPS     float64
	EventQueueLen int
}
