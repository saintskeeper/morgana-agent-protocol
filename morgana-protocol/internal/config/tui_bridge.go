package config

import (
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
)

// ToTUIConfig converts the configuration TUIConfig to the internal tui.TUIConfig
// This maintains compatibility between the config system and TUI implementation
func (c *Config) ToTUIConfig() tui.TUIConfig {
	return tui.TUIConfig{
		RefreshRate:     c.TUI.Performance.RefreshRate,
		EventBufferSize: c.TUI.Events.BufferSize,
		MaxLogLines:     c.TUI.Performance.MaxLogLines,
		Theme:           c.toTUITheme(),
		ShowDebugInfo:   c.TUI.Visual.ShowDebugInfo,
		ShowTimeStamps:  c.TUI.Visual.ShowTimeStamps,
		CompactMode:     c.TUI.Visual.CompactMode,
		EnableFiltering: c.TUI.Features.EnableFiltering,
		EnableSearch:    c.TUI.Features.EnableSearch,
		EnableExport:    c.TUI.Features.EnableExport,
	}
}

// toTUITheme converts config theme to TUI theme
func (c *Config) toTUITheme() tui.Theme {
	theme := tui.Theme{
		Primary:     c.TUI.Visual.Theme.Primary,
		Secondary:   c.TUI.Visual.Theme.Secondary,
		Success:     c.TUI.Visual.Theme.Success,
		Warning:     c.TUI.Visual.Theme.Warning,
		Error:       c.TUI.Visual.Theme.Error,
		Info:        c.TUI.Visual.Theme.Info,
		Background:  c.TUI.Visual.Theme.Background,
		Foreground:  c.TUI.Visual.Theme.Foreground,
		Border:      c.TUI.Visual.Theme.Border,
		Highlight:   c.TUI.Visual.Theme.Highlight,
		Muted:       c.TUI.Visual.Theme.Muted,
		AgentColors: make(map[string]string),
	}

	// Copy agent colors
	for agent, color := range c.TUI.Visual.Theme.AgentColors {
		theme.AgentColors[agent] = color
	}

	// Set progress bar colors
	theme.ProgressBar.Complete = c.TUI.Visual.Theme.ProgressBar.Complete
	theme.ProgressBar.Incomplete = c.TUI.Visual.Theme.ProgressBar.Incomplete
	theme.ProgressBar.Background = c.TUI.Visual.Theme.ProgressBar.Background

	return theme
}

// CreateTUIConfigPresets creates common TUI configuration presets
func CreateTUIConfigPresets() map[string]TUIConfig {
	presets := make(map[string]TUIConfig)

	// Performance preset - optimized for high performance
	presets["performance"] = TUIConfig{
		Enabled: true,
		Performance: TUIPerformanceConfig{
			RefreshRate:        33 * time.Millisecond, // 30fps for lower CPU usage
			MaxLogLines:        5000,                  // Reduced memory usage
			OptimizedRendering: true,
			TargetFPS:          30,
		},
		Visual: TUIVisualConfig{
			Theme:          createDarkTheme(),
			ShowDebugInfo:  false,
			ShowTimeStamps: false, // Disable for less rendering work
			CompactMode:    true,  // More compact display
			MinWidth:       60,
			MinHeight:      20,
		},
		Features: TUIFeatureConfig{
			EnableFiltering:         false, // Disable expensive filtering
			EnableSearch:            false, // Disable search overhead
			EnableExport:            false,
			EnableKeyboardShortcuts: true,
			EnableMouse:             false, // Disable mouse for performance
			EnableAutoScroll:        true,
		},
		Events: TUIEventConfig{
			BufferSize:        5000, // Large buffer
			EnableBatching:    true,
			BatchSize:         100,                    // Larger batches
			BatchTimeout:      200 * time.Millisecond, // Longer timeout
			ProcessingTimeout: 10 * time.Second,
		},
	}

	// Development preset - all features enabled
	presets["development"] = TUIConfig{
		Enabled: true,
		Performance: TUIPerformanceConfig{
			RefreshRate:        16 * time.Millisecond, // 60fps
			MaxLogLines:        20000,                 // Large history
			OptimizedRendering: false,                 // All features enabled
			TargetFPS:          60,
		},
		Visual: TUIVisualConfig{
			Theme:          createDarkTheme(),
			ShowDebugInfo:  true, // Show performance info
			ShowTimeStamps: true,
			CompactMode:    false,
			MinWidth:       100, // Wider for debug info
			MinHeight:      30,  // Taller for debug info
		},
		Features: TUIFeatureConfig{
			EnableFiltering:         true, // Enable all features
			EnableSearch:            true,
			EnableExport:            true, // Enable export for debugging
			EnableKeyboardShortcuts: true,
			EnableMouse:             true,
			EnableAutoScroll:        true,
		},
		Events: TUIEventConfig{
			BufferSize:        2000,
			EnableBatching:    true,
			BatchSize:         25,                    // Smaller batches for responsiveness
			BatchTimeout:      50 * time.Millisecond, // Quick processing
			ProcessingTimeout: 5 * time.Second,
		},
	}

	// Minimal preset - basic functionality only
	presets["minimal"] = TUIConfig{
		Enabled: true,
		Performance: TUIPerformanceConfig{
			RefreshRate:        50 * time.Millisecond, // 20fps
			MaxLogLines:        1000,                  // Minimal memory usage
			OptimizedRendering: true,
			TargetFPS:          20,
		},
		Visual: TUIVisualConfig{
			Theme:          createLightTheme(), // Light theme for minimal contrast
			ShowDebugInfo:  false,
			ShowTimeStamps: false,
			CompactMode:    true,
			MinWidth:       40,
			MinHeight:      12,
		},
		Features: TUIFeatureConfig{
			EnableFiltering:         false,
			EnableSearch:            false,
			EnableExport:            false,
			EnableKeyboardShortcuts: false,
			EnableMouse:             false,
			EnableAutoScroll:        false,
		},
		Events: TUIEventConfig{
			BufferSize:        500,
			EnableBatching:    false, // No batching for simplicity
			BatchSize:         10,
			BatchTimeout:      1 * time.Second,
			ProcessingTimeout: 3 * time.Second,
		},
	}

	return presets
}

// Helper functions for creating theme presets

func createDarkTheme() TUIThemeConfig {
	return TUIThemeConfig{
		Name:       "dark",
		Primary:    "#7C3AED", // Violet
		Secondary:  "#06B6D4", // Cyan
		Success:    "#10B981", // Emerald
		Warning:    "#F59E0B", // Amber
		Error:      "#EF4444", // Red
		Info:       "#3B82F6", // Blue
		Background: "#0F172A", // Slate 900
		Foreground: "#F1F5F9", // Slate 100
		Border:     "#334155", // Slate 700
		Highlight:  "#1E293B", // Slate 800
		Muted:      "#64748B", // Slate 500
		ProgressBar: TUIProgressBarTheme{
			Complete:   "#10B981",
			Incomplete: "#64748B",
			Background: "#334155",
		},
		AgentColors: map[string]string{
			"code-implementer":  "#10B981", // Green
			"sprint-planner":    "#3B82F6", // Blue
			"test-specialist":   "#F59E0B", // Amber
			"validation-expert": "#EF4444", // Red
			"default":           "#64748B", // Slate
		},
	}
}

func createLightTheme() TUIThemeConfig {
	return TUIThemeConfig{
		Name:       "light",
		Primary:    "#8B5CF6", // Violet 400
		Secondary:  "#06B6D4", // Cyan 500
		Success:    "#059669", // Emerald 600
		Warning:    "#D97706", // Amber 600
		Error:      "#DC2626", // Red 600
		Info:       "#2563EB", // Blue 600
		Background: "#F8FAFC", // Slate 50
		Foreground: "#0F172A", // Slate 900
		Border:     "#CBD5E1", // Slate 300
		Highlight:  "#E2E8F0", // Slate 200
		Muted:      "#64748B", // Slate 500
		ProgressBar: TUIProgressBarTheme{
			Complete:   "#059669",
			Incomplete: "#CBD5E1",
			Background: "#E2E8F0",
		},
		AgentColors: map[string]string{
			"code-implementer":  "#059669", // Green 600
			"sprint-planner":    "#2563EB", // Blue 600
			"test-specialist":   "#D97706", // Amber 600
			"validation-expert": "#DC2626", // Red 600
			"default":           "#64748B", // Slate 500
		},
	}
}

// LoadConfigWithPreset loads configuration and applies a preset if specified
func LoadConfigWithPreset(path string, preset string) (*Config, error) {
	// Load base configuration
	cfg, err := LoadFile(path)
	if err != nil {
		return nil, err
	}

	// Apply preset if specified
	if preset != "" {
		presets := CreateTUIConfigPresets()
		if presetConfig, exists := presets[preset]; exists {
			cfg.TUI = presetConfig
		}
	}

	return cfg, nil
}
