//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// CommandState represents the state of command polling
type CommandState string

const (
	CommandStateRunning CommandState = "running"
	CommandStatePaused  CommandState = "paused"
	CommandStateStopped CommandState = "stopped"
)

// CommandPoller simulates the command polling functionality
type CommandPoller struct {
	commandFile string
	state       CommandState
	stopCh      chan struct{}
	pauseCh     chan struct{}
	resumeCh    chan struct{}
	stateCh     chan CommandState
}

// NewCommandPoller creates a new command poller
func NewCommandPoller(commandFile string) *CommandPoller {
	return &CommandPoller{
		commandFile: commandFile,
		state:       CommandStateRunning,
		stopCh:      make(chan struct{}),
		pauseCh:     make(chan struct{}),
		resumeCh:    make(chan struct{}),
		stateCh:     make(chan CommandState, 10),
	}
}

// Start begins command polling
func (cp *CommandPoller) Start(ctx context.Context) error {
	ticker := time.NewTicker(100 * time.Millisecond) // Poll every 100ms
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cp.stopCh:
			cp.state = CommandStateStopped
			cp.stateCh <- cp.state
			return nil
		case <-ticker.C:
			if err := cp.pollCommand(); err != nil {
				return err
			}
		}

		// Handle pause/resume
		if cp.state == CommandStatePaused {
			select {
			case <-cp.resumeCh:
				cp.state = CommandStateRunning
				cp.stateCh <- cp.state
			case <-ctx.Done():
				return ctx.Err()
			case <-cp.stopCh:
				cp.state = CommandStateStopped
				cp.stateCh <- cp.state
				return nil
			}
		}
	}
}

// pollCommand checks for command file and processes commands
func (cp *CommandPoller) pollCommand() error {
	if _, err := os.Stat(cp.commandFile); os.IsNotExist(err) {
		return nil // No command file, continue polling
	}

	// Read command from file
	data, err := os.ReadFile(cp.commandFile)
	if err != nil {
		return err
	}

	command := string(data)
	switch command {
	case "pause\n", "pause":
		if cp.state == CommandStateRunning {
			cp.state = CommandStatePaused
			cp.stateCh <- cp.state
		}
	case "resume\n", "resume":
		if cp.state == CommandStatePaused {
			cp.resumeCh <- struct{}{}
		}
	case "stop\n", "stop":
		cp.stopCh <- struct{}{}
	}

	// Remove command file after processing
	os.Remove(cp.commandFile)
	return nil
}

// GetState returns the current state
func (cp *CommandPoller) GetState() CommandState {
	return cp.state
}

// Stop stops the command poller
func (cp *CommandPoller) Stop() {
	select {
	case cp.stopCh <- struct{}{}:
	default:
	}
}

// GetStateChanges returns a channel for state changes
func (cp *CommandPoller) GetStateChanges() <-chan CommandState {
	return cp.stateCh
}

func TestCommandPolling(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("BasicCommandPolling", func(t *testing.T) {
		commandFile := filepath.Join(setup.TempDir, "commands.txt")
		poller := NewCommandPoller(commandFile)

		// Start poller in background
		go poller.Start(ctx)

		// Wait for initial state
		time.Sleep(100 * time.Millisecond)

		if poller.GetState() != CommandStateRunning {
			t.Errorf("Expected initial state %s, got %s", CommandStateRunning, poller.GetState())
		}

		// Test pause command
		err := os.WriteFile(commandFile, []byte("pause"), 0644)
		if err != nil {
			t.Fatalf("Failed to write pause command: %v", err)
		}

		// Wait for state change
		select {
		case state := <-poller.GetStateChanges():
			if state != CommandStatePaused {
				t.Errorf("Expected state %s, got %s", CommandStatePaused, state)
			}
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for pause state change")
		}

		// Test resume command
		err = os.WriteFile(commandFile, []byte("resume"), 0644)
		if err != nil {
			t.Fatalf("Failed to write resume command: %v", err)
		}

		// Wait for state change
		select {
		case state := <-poller.GetStateChanges():
			if state != CommandStateRunning {
				t.Errorf("Expected state %s, got %s", CommandStateRunning, state)
			}
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for resume state change")
		}

		// Test stop command
		err = os.WriteFile(commandFile, []byte("stop"), 0644)
		if err != nil {
			t.Fatalf("Failed to write stop command: %v", err)
		}

		// Wait for state change
		select {
		case state := <-poller.GetStateChanges():
			if state != CommandStateStopped {
				t.Errorf("Expected state %s, got %s", CommandStateStopped, state)
			}
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for stop state change")
		}

		t.Log("Successfully tested basic command polling")
	})

	t.Run("CommandPollingWithEvents", func(t *testing.T) {
		// Test command polling while processing events
		commandFile := filepath.Join(setup.TempDir, "event_commands.txt")
		poller := NewCommandPoller(commandFile)

		go poller.Start(ctx)
		time.Sleep(50 * time.Millisecond)

		// Start generating events
		go func() {
			for i := 0; i < 100; i++ {
				if poller.GetState() == CommandStateStopped {
					break
				}
				setup.Generator.GenerateTaskLifecycle(ctx, "polling-test", 10*time.Millisecond)
				time.Sleep(20 * time.Millisecond)
			}
		}()

		// Let events generate for a bit
		time.Sleep(200 * time.Millisecond)
		initialEvents := setup.Collector.GetStats().TotalEvents

		// Pause event processing
		err := os.WriteFile(commandFile, []byte("pause"), 0644)
		if err != nil {
			t.Fatalf("Failed to write pause command: %v", err)
		}

		// Wait for pause
		select {
		case state := <-poller.GetStateChanges():
			if state != CommandStatePaused {
				t.Errorf("Expected paused state, got %s", state)
			}
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for pause")
		}

		// Let some time pass while paused
		time.Sleep(200 * time.Millisecond)
		pausedEvents := setup.Collector.GetStats().TotalEvents

		// Resume event processing
		err = os.WriteFile(commandFile, []byte("resume"), 0644)
		if err != nil {
			t.Fatalf("Failed to write resume command: %v", err)
		}

		// Wait for resume
		select {
		case state := <-poller.GetStateChanges():
			if state != CommandStateRunning {
				t.Errorf("Expected running state, got %s", state)
			}
		case <-time.After(time.Second):
			t.Fatal("Timeout waiting for resume")
		}

		// Let events process after resume
		time.Sleep(200 * time.Millisecond)
		finalEvents := setup.Collector.GetStats().TotalEvents

		// Stop everything
		poller.Stop()

		t.Logf("Event processing during commands: Initial=%d, Paused=%d, Final=%d",
			initialEvents, pausedEvents, finalEvents)

		// Events should have been processed initially and after resume
		if finalEvents <= initialEvents {
			t.Error("No events processed after resume")
		}
	})

	t.Run("PollingPerformanceImpact", func(t *testing.T) {
		// Test that command polling doesn't significantly impact performance
		commandFile := filepath.Join(setup.TempDir, "perf_commands.txt")

		// First, measure performance without polling
		start := time.Now()
		setup.Generator.GenerateHighVolumeEvents(ctx, 500, RandomAgentTypes())
		time.Sleep(time.Second)
		nopollingStats := setup.Collector.GetStats()
		nopollingDuration := time.Since(start)

		// Reset collector
		setup.Collector.Reset()

		// Now measure with polling
		poller := NewCommandPoller(commandFile)
		go poller.Start(ctx)

		start = time.Now()
		setup.Generator.GenerateHighVolumeEvents(ctx, 500, RandomAgentTypes())
		time.Sleep(time.Second)
		pollingStats := setup.Collector.GetStats()
		pollingDuration := time.Since(start)

		poller.Stop()

		// Performance should be comparable (within 20% degradation)
		performanceRatio := float64(pollingDuration) / float64(nopollingDuration)
		if performanceRatio > 1.2 {
			t.Errorf("Command polling caused significant performance degradation: %.2fx slower",
				performanceRatio)
		}

		eventRatio := float64(pollingStats.TotalEvents) / float64(nopollingStats.TotalEvents)
		if eventRatio < 0.8 {
			t.Errorf("Command polling caused significant event loss: %.1f%% of events processed",
				eventRatio*100)
		}

		t.Logf("Polling performance impact: Duration %.2fx, Events %.1f%%",
			performanceRatio, eventRatio*100)
	})

	t.Run("CommandFileCleanup", func(t *testing.T) {
		// Test that command files are properly cleaned up
		commandFile := filepath.Join(setup.TempDir, "cleanup_commands.txt")
		poller := NewCommandPoller(commandFile)

		go poller.Start(ctx)
		time.Sleep(50 * time.Millisecond)

		// Write multiple commands and verify they're cleaned up
		commands := []string{"pause", "resume", "pause", "resume"}

		for _, cmd := range commands {
			err := os.WriteFile(commandFile, []byte(cmd), 0644)
			if err != nil {
				t.Fatalf("Failed to write command %s: %v", cmd, err)
			}

			// Wait for processing
			time.Sleep(150 * time.Millisecond)

			// Verify file is cleaned up
			if _, err := os.Stat(commandFile); !os.IsNotExist(err) {
				t.Errorf("Command file not cleaned up after processing command: %s", cmd)
			}
		}

		poller.Stop()
		t.Log("Successfully verified command file cleanup")
	})

	t.Run("InvalidCommands", func(t *testing.T) {
		// Test handling of invalid commands
		commandFile := filepath.Join(setup.TempDir, "invalid_commands.txt")
		poller := NewCommandPoller(commandFile)

		go poller.Start(ctx)
		time.Sleep(50 * time.Millisecond)

		initialState := poller.GetState()

		// Write invalid commands
		invalidCommands := []string{"invalid", "badcommand", "123", "pause_typo"}

		for _, cmd := range invalidCommands {
			err := os.WriteFile(commandFile, []byte(cmd), 0644)
			if err != nil {
				t.Fatalf("Failed to write invalid command %s: %v", cmd, err)
			}

			time.Sleep(100 * time.Millisecond)

			// State should remain unchanged for invalid commands
			if poller.GetState() != initialState {
				t.Errorf("Invalid command %s changed state from %s to %s",
					cmd, initialState, poller.GetState())
			}
		}

		poller.Stop()
		t.Log("Successfully handled invalid commands")
	})

	t.Run("ConcurrentCommandsAndEvents", func(t *testing.T) {
		// Test concurrent command processing and event handling
		commandFile := filepath.Join(setup.TempDir, "concurrent_commands.txt")
		poller := NewCommandPoller(commandFile)

		go poller.Start(ctx)
		time.Sleep(50 * time.Millisecond)

		// Generate events continuously
		go func() {
			for i := 0; i < 50; i++ {
				if poller.GetState() == CommandStateStopped {
					break
				}
				setup.Generator.GenerateTaskLifecycle(ctx, "concurrent-test", 5*time.Millisecond)
				time.Sleep(10 * time.Millisecond)
			}
		}()

		// Issue commands concurrently
		go func() {
			commands := []string{"pause", "resume", "pause", "resume"}
			for _, cmd := range commands {
				if poller.GetState() == CommandStateStopped {
					break
				}
				os.WriteFile(commandFile, []byte(cmd), 0644)
				time.Sleep(100 * time.Millisecond)
			}
		}()

		// Monitor state changes
		stateChanges := 0
		timeout := time.After(2 * time.Second)

	stateLoop:
		for {
			select {
			case <-poller.GetStateChanges():
				stateChanges++
			case <-timeout:
				break stateLoop
			}
		}

		poller.Stop()

		if stateChanges < 4 {
			t.Errorf("Expected at least 4 state changes, got %d", stateChanges)
		}

		finalStats := setup.Collector.GetStats()
		if finalStats.TotalEvents == 0 {
			t.Error("No events processed during concurrent operations")
		}

		t.Logf("Concurrent test: %d state changes, %d events processed",
			stateChanges, finalStats.TotalEvents)
	})
}
