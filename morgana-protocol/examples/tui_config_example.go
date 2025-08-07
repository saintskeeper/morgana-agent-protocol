package main

import (
	"fmt"
	"log"
	"os"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/config"
	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/tui"
)

// Example demonstrating comprehensive TUI configuration usage
func main() {
	if len(os.Args) > 1 && os.Args[1] == "help" {
		showUsage()
		return
	}

	// Load configuration from YAML file
	cfg, err := config.LoadFile("morgana.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("=== Morgana Protocol TUI Configuration Demo ===")

	// Display current configuration
	showCurrentConfig(cfg)

	// Demonstrate configuration presets
	demonstratePresets()

	// Show environment variable usage
	showEnvironmentVariables()

	// Demonstrate runtime configuration updates
	demonstrateRuntimeUpdates(cfg)

	// Show TUI integration example
	demonstrateTUIIntegration(cfg)
}

func showCurrentConfig(cfg *config.Config) {
	fmt.Println("\n--- Current TUI Configuration ---")
	fmt.Printf("TUI Enabled: %v\n", cfg.IsTUIEnabled())
	fmt.Printf("Refresh Rate: %v (%d FPS)\n", cfg.GetTUIRefreshRate(), cfg.TUI.Performance.TargetFPS)
	fmt.Printf("Theme: %s\n", cfg.TUI.Visual.Theme.Name)
	fmt.Printf("Primary Color: %s\n", cfg.TUI.Visual.Theme.Primary)
	fmt.Printf("Compact Mode: %v\n", cfg.TUI.Visual.CompactMode)
	fmt.Printf("Debug Info: %v\n", cfg.TUI.Visual.ShowDebugInfo)
	fmt.Printf("Max Log Lines: %d\n", cfg.TUI.Performance.MaxLogLines)
	fmt.Printf("Event Buffer Size: %d\n", cfg.TUI.Events.BufferSize)

	fmt.Println("\nAgent Colors:")
	for agent, color := range cfg.TUI.Visual.Theme.AgentColors {
		fmt.Printf("  %s: %s\n", agent, color)
	}

	fmt.Printf("\nFeature Flags:\n")
	fmt.Printf("  Filtering: %v\n", cfg.TUI.Features.EnableFiltering)
	fmt.Printf("  Search: %v\n", cfg.TUI.Features.EnableSearch)
	fmt.Printf("  Export: %v\n", cfg.TUI.Features.EnableExport)
	fmt.Printf("  Mouse: %v\n", cfg.TUI.Features.EnableMouse)
	fmt.Printf("  Auto-scroll: %v\n", cfg.TUI.Features.EnableAutoScroll)
}

func demonstratePresets() {
	fmt.Println("\n--- Configuration Presets ---")
	presets := config.CreateTUIConfigPresets()

	for name, preset := range presets {
		fmt.Printf("\n%s Preset:\n", name)
		fmt.Printf("  Refresh Rate: %v\n", preset.Performance.RefreshRate)
		fmt.Printf("  Max Log Lines: %d\n", preset.Performance.MaxLogLines)
		fmt.Printf("  Show Debug: %v\n", preset.Visual.ShowDebugInfo)
		fmt.Printf("  Filtering: %v\n", preset.Features.EnableFiltering)
		fmt.Printf("  Export: %v\n", preset.Features.EnableExport)
		fmt.Printf("  Buffer Size: %d\n", preset.Events.BufferSize)
	}

	fmt.Println("\nTo use a preset:")
	fmt.Println("  cfg, err := config.LoadConfigWithPreset(\"morgana.yaml\", \"performance\")")
}

func showEnvironmentVariables() {
	fmt.Println("\n--- Environment Variable Overrides ---")
	fmt.Println("Available environment variables:")

	envVars := []struct {
		name        string
		description string
		example     string
	}{
		{"MORGANA_TUI_ENABLED", "Enable/disable TUI", "true or false"},
		{"MORGANA_TUI_REFRESH_RATE", "Target refresh rate", "16ms (60fps) or 33ms (30fps)"},
		{"MORGANA_TUI_MAX_LOG_LINES", "Maximum log lines in memory", "10000"},
		{"MORGANA_TUI_BUFFER_SIZE", "Event buffer size", "1000"},
		{"MORGANA_TUI_THEME", "Color theme", "dark, light, or custom"},
		{"MORGANA_TUI_COMPACT_MODE", "Use compact display", "true or false"},
		{"MORGANA_TUI_SHOW_DEBUG", "Show debug information", "true or false"},
	}

	for _, env := range envVars {
		fmt.Printf("  %s\n    %s\n    Example: %s\n\n", env.name, env.description, env.example)
	}

	fmt.Println("Usage:")
	fmt.Println("  MORGANA_TUI_THEME=light MORGANA_TUI_COMPACT_MODE=true ./morgana-protocol")
}

func demonstrateRuntimeUpdates(cfg *config.Config) {
	fmt.Println("\n--- Runtime Configuration Updates ---")

	// Create a copy for testing
	originalRefreshRate := cfg.TUI.Performance.RefreshRate
	fmt.Printf("Original refresh rate: %v\n", originalRefreshRate)

	// Create an updated configuration
	updatedTUI := cfg.TUI
	updatedTUI.Performance.RefreshRate = 33 * 1000000 // 33ms
	updatedTUI.Visual.CompactMode = !updatedTUI.Visual.CompactMode

	// Apply the update
	err := cfg.UpdateTUIConfig(updatedTUI)
	if err != nil {
		fmt.Printf("Failed to update config: %v\n", err)
		return
	}

	fmt.Printf("Updated refresh rate: %v\n", cfg.TUI.Performance.RefreshRate)
	fmt.Printf("Updated compact mode: %v\n", cfg.TUI.Visual.CompactMode)

	// Restore original
	cfg.TUI.Performance.RefreshRate = originalRefreshRate
	fmt.Println("Configuration restored")
}

func demonstrateTUIIntegration(cfg *config.Config) {
	fmt.Println("\n--- TUI Integration Example ---")
	fmt.Println("Converting config format for TUI usage...")

	// Convert to TUI format
	tuiConfig := cfg.ToTUIConfig()
	fmt.Printf("Converted refresh rate: %v\n", tuiConfig.RefreshRate)
	fmt.Printf("Converted buffer size: %d\n", tuiConfig.EventBufferSize)
	fmt.Printf("Converted theme primary: %s\n", tuiConfig.Theme.Primary)

	// Show how to use with TUI system
	fmt.Println("\nTUI System Integration:")
	fmt.Println("  ctx := context.Background()")
	fmt.Println("  eventBus := events.NewEventBus()")
	fmt.Println("  tui := tui.New(ctx, eventBus, tuiConfig)")
	fmt.Println("  err := tui.Start()")

	// Example validation check
	fmt.Println("\nConfiguration validation:")
	if err := cfg.Validate(); err != nil {
		fmt.Printf("  Validation failed: %v\n", err)
	} else {
		fmt.Printf("  ✓ Configuration is valid\n")
	}

	// Terminal compatibility check
	if tui.IsTerminalSupported() {
		fmt.Println("  ✓ Terminal supports TUI features")
	} else {
		fmt.Println("  ⚠ Terminal may not fully support TUI features")
	}
}

func showUsage() {
	fmt.Println("Morgana Protocol TUI Configuration Example")
	fmt.Println()
	fmt.Println("This example demonstrates:")
	fmt.Println("• Loading TUI configuration from YAML")
	fmt.Println("• Environment variable overrides")
	fmt.Println("• Configuration presets")
	fmt.Println("• Runtime configuration updates")
	fmt.Println("• Integration with TUI system")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run examples/tui_config_example.go")
	fmt.Println("  go run examples/tui_config_example.go help")
	fmt.Println()
	fmt.Println("Environment variable examples:")
	fmt.Println("  MORGANA_TUI_THEME=light go run examples/tui_config_example.go")
	fmt.Println("  MORGANA_TUI_COMPACT_MODE=true go run examples/tui_config_example.go")
}
