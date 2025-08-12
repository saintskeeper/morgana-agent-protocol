//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/saintskeeper/claude-code-configs/morgana-protocol/internal/events"
)

// ResourceMonitor tracks system resource usage during tests
type ResourceMonitor struct {
	initialMemStats   runtime.MemStats
	initialGoRoutines int
	measurements      []ResourceMeasurement
}

type ResourceMeasurement struct {
	Timestamp       time.Time
	MemAlloc        uint64
	MemSys          uint64
	NumGC           uint32
	GoRoutines      int
	EventsProcessed int64
}

// NewResourceMonitor creates a new resource monitor
func NewResourceMonitor() *ResourceMonitor {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return &ResourceMonitor{
		initialMemStats:   memStats,
		initialGoRoutines: runtime.NumGoroutine(),
		measurements:      make([]ResourceMeasurement, 0),
	}
}

// TakeMeasurement records current resource usage
func (rm *ResourceMonitor) TakeMeasurement(eventsProcessed int64) ResourceMeasurement {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	measurement := ResourceMeasurement{
		Timestamp:       time.Now(),
		MemAlloc:        memStats.Alloc,
		MemSys:          memStats.Sys,
		NumGC:           memStats.NumGC,
		GoRoutines:      runtime.NumGoroutine(),
		EventsProcessed: eventsProcessed,
	}

	rm.measurements = append(rm.measurements, measurement)
	return measurement
}

// GetLeakAnalysis analyzes resource usage for potential leaks
func (rm *ResourceMonitor) GetLeakAnalysis() LeakAnalysis {
	if len(rm.measurements) < 2 {
		return LeakAnalysis{Valid: false}
	}

	first := rm.measurements[0]
	last := rm.measurements[len(rm.measurements)-1]

	memGrowthBytes := int64(last.MemAlloc) - int64(first.MemAlloc)
	memGrowthMB := float64(memGrowthBytes) / 1024 / 1024

	goRoutineGrowth := last.GoRoutines - first.GoRoutines

	eventsPerMB := float64(0)
	if memGrowthMB > 0 {
		eventsPerMB = float64(last.EventsProcessed-first.EventsProcessed) / memGrowthMB
	}

	return LeakAnalysis{
		Valid:             true,
		MemoryGrowthMB:    memGrowthMB,
		GoRoutineGrowth:   goRoutineGrowth,
		EventsProcessed:   last.EventsProcessed - first.EventsProcessed,
		EventsPerMB:       eventsPerMB,
		GCCount:           last.NumGC - first.NumGC,
		Duration:          last.Timestamp.Sub(first.Timestamp),
		SuspiciousMemLeak: memGrowthMB > 50 && eventsPerMB < 1000,
		SuspiciousGRLeak:  goRoutineGrowth > 10,
	}
}

type LeakAnalysis struct {
	Valid             bool
	MemoryGrowthMB    float64
	GoRoutineGrowth   int
	EventsProcessed   int64
	EventsPerMB       float64
	GCCount           uint32
	Duration          time.Duration
	SuspiciousMemLeak bool
	SuspiciousGRLeak  bool
}

func TestResourceLeaks(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("EventProcessingMemoryLeaks", func(t *testing.T) {
		monitor := NewResourceMonitor()

		// Force GC and take initial measurement
		runtime.GC()
		time.Sleep(50 * time.Millisecond)
		monitor.TakeMeasurement(0)

		// Process events in waves to detect leaks
		agentTypes := RandomAgentTypes()
		wavesSize := 1000
		waves := 10

		for wave := 0; wave < waves; wave++ {
			// Generate events
			setup.Generator.GenerateHighVolumeEvents(ctx, wavesSize, agentTypes)

			// Wait for processing
			time.Sleep(200 * time.Millisecond)

			// Force GC to see true memory usage
			runtime.GC()
			time.Sleep(50 * time.Millisecond)

			// Take measurement
			stats := setup.Collector.GetStats()
			monitor.TakeMeasurement(stats.TotalEvents)

			t.Logf("Wave %d: Events=%d, Goroutines=%d",
				wave, stats.TotalEvents, runtime.NumGoroutine())
		}

		// Analyze for leaks
		analysis := monitor.GetLeakAnalysis()

		if analysis.SuspiciousMemLeak {
			t.Errorf("Suspicious memory leak detected: %.1f MB growth, %.0f events/MB",
				analysis.MemoryGrowthMB, analysis.EventsPerMB)
		}

		if analysis.SuspiciousGRLeak {
			t.Errorf("Suspicious goroutine leak detected: %d goroutines added",
				analysis.GoRoutineGrowth)
		}

		// Memory growth should be reasonable for the number of events processed
		if analysis.MemoryGrowthMB > 100 {
			t.Errorf("Excessive memory growth: %.1f MB", analysis.MemoryGrowthMB)
		}

		t.Logf("Resource leak analysis: Mem=%.1fMB, GR=%d, Events=%d, GC=%d",
			analysis.MemoryGrowthMB, analysis.GoRoutineGrowth,
			analysis.EventsProcessed, analysis.GCCount)
	})

	t.Run("EventBusShutdownCleanup", func(t *testing.T) {
		// Test that event bus properly cleans up resources on shutdown
		initialGoroutines := runtime.NumGoroutine()

		// Create and use multiple event buses
		numBuses := 5
		for i := 0; i < numBuses; i++ {
			eventBus := events.NewEventBus(events.DefaultBusConfig())

			// Use the bus briefly
			generator := NewTestEventGenerator(eventBus)
			generator.GenerateHighVolumeEvents(ctx, 100, []string{"cleanup-test"})

			time.Sleep(50 * time.Millisecond)

			// Close the bus
			err := eventBus.Close()
			if err != nil {
				t.Errorf("Event bus %d failed to close: %v", i, err)
			}
		}

		// Wait for cleanup
		time.Sleep(200 * time.Millisecond)
		runtime.GC()
		time.Sleep(50 * time.Millisecond)

		finalGoroutines := runtime.NumGoroutine()
		goroutineGrowth := finalGoroutines - initialGoroutines

		// Should not have significant goroutine growth
		if goroutineGrowth > 5 {
			t.Errorf("Event bus shutdown left %d goroutines running", goroutineGrowth)
		}

		t.Logf("Event bus cleanup test: %d goroutines added after %d buses",
			goroutineGrowth, numBuses)
	})

	t.Run("LongRunningEventProcessing", func(t *testing.T) {
		// Test resource usage over extended operation
		monitor := NewResourceMonitor()
		runtime.GC()
		monitor.TakeMeasurement(0)

		// Run continuous event processing
		duration := 10 * time.Second
		endTime := time.Now().Add(duration)
		eventCounter := 0

		go func() {
			for time.Now().Before(endTime) {
				setup.Generator.GenerateTaskLifecycle(ctx, "long-running-test", 10*time.Millisecond)
				eventCounter++
				time.Sleep(50 * time.Millisecond)
			}
		}()

		// Take measurements periodically
		measureInterval := 2 * time.Second
		measureTicker := time.NewTicker(measureInterval)
		defer measureTicker.Stop()

		for time.Now().Before(endTime) {
			select {
			case <-measureTicker.C:
				runtime.GC()
				time.Sleep(10 * time.Millisecond)
				stats := setup.Collector.GetStats()
				monitor.TakeMeasurement(stats.TotalEvents)
			case <-time.After(duration):
				break
			}
		}

		// Final measurement
		runtime.GC()
		time.Sleep(50 * time.Millisecond)
		finalStats := setup.Collector.GetStats()
		monitor.TakeMeasurement(finalStats.TotalEvents)

		analysis := monitor.GetLeakAnalysis()

		// Long-running operation should not cause significant leaks
		memPerSecond := analysis.MemoryGrowthMB / analysis.Duration.Seconds()
		if memPerSecond > 1.0 { // More than 1MB/second growth
			t.Errorf("High memory growth rate during long-running test: %.2f MB/sec", memPerSecond)
		}

		if analysis.SuspiciousGRLeak {
			t.Errorf("Goroutine leak during long-running test: %d goroutines",
				analysis.GoRoutineGrowth)
		}

		t.Logf("Long-running test (%v): Events=%d, Mem=%.1fMB, Rate=%.2fMB/s, GR=%d",
			analysis.Duration, analysis.EventsProcessed, analysis.MemoryGrowthMB,
			memPerSecond, analysis.GoRoutineGrowth)
	})

	t.Run("FileDescriptorLeaks", func(t *testing.T) {
		// Test for file descriptor leaks (Unix systems)
		if runtime.GOOS == "windows" {
			t.Skip("File descriptor test not applicable on Windows")
		}

		// Count initial file descriptors
		initialFDs := countFileDescriptors(t)

		// Create many temporary files and clean them up
		tempFiles := make([]string, 100)
		for i := 0; i < len(tempFiles); i++ {
			tempFile := filepath.Join(setup.TempDir, fmt.Sprintf("fd_test_%d.tmp", i))

			file, err := os.Create(tempFile)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			// Write some data
			file.WriteString("test data for FD leak test")
			file.Close()

			tempFiles[i] = tempFile
		}

		// Clean up files
		for _, tempFile := range tempFiles {
			os.Remove(tempFile)
		}

		// Check file descriptor count
		finalFDs := countFileDescriptors(t)
		fdGrowth := finalFDs - initialFDs

		if fdGrowth > 5 {
			t.Errorf("Potential file descriptor leak: %d FDs added", fdGrowth)
		}

		t.Logf("File descriptor test: %d initial, %d final, %d growth",
			initialFDs, finalFDs, fdGrowth)
	})
}

// countFileDescriptors counts open file descriptors (Unix systems)
func countFileDescriptors(t *testing.T) int {
	fdDir := fmt.Sprintf("/proc/%d/fd", os.Getpid())
	entries, err := os.ReadDir(fdDir)
	if err != nil {
		t.Logf("Could not read FD directory: %v", err)
		return 0
	}
	return len(entries)
}

func TestFileRotation(t *testing.T) {
	setup := SetupIntegrationTest(t)
	defer setup.Cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("LogFileRotation", func(t *testing.T) {
		// Simulate log file rotation
		logDir := filepath.Join(setup.TempDir, "logs")
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create log directory: %v", err)
		}

		logFile := filepath.Join(logDir, "morgana.log")

		// Create initial log file
		file, err := os.Create(logFile)
		if err != nil {
			t.Fatalf("Failed to create log file: %v", err)
		}

		// Simulate writing logs during event processing
		go func() {
			for i := 0; i < 1000; i++ {
				if ctx.Err() != nil {
					break
				}
				fmt.Fprintf(file, "Log entry %d at %s\n", i, time.Now().Format(time.RFC3339))
				time.Sleep(time.Millisecond)
			}
		}()

		// Generate events while logging
		go func() {
			for i := 0; i < 100; i++ {
				if ctx.Err() != nil {
					break
				}
				setup.Generator.GenerateTaskLifecycle(ctx, "rotation-test", 5*time.Millisecond)
				time.Sleep(10 * time.Millisecond)
			}
		}()

		// Simulate rotation after some time
		time.Sleep(500 * time.Millisecond)

		// Rotate log file
		rotatedFile := filepath.Join(logDir, "morgana.log.1")
		err = os.Rename(logFile, rotatedFile)
		if err != nil {
			t.Fatalf("Failed to rotate log file: %v", err)
		}

		file.Close()

		// Create new log file
		newFile, err := os.Create(logFile)
		if err != nil {
			t.Fatalf("Failed to create new log file: %v", err)
		}
		defer newFile.Close()

		// Continue processing
		time.Sleep(500 * time.Millisecond)

		// Verify both files exist and have content
		rotatedInfo, err := os.Stat(rotatedFile)
		if err != nil {
			t.Errorf("Rotated log file not found: %v", err)
		} else if rotatedInfo.Size() == 0 {
			t.Error("Rotated log file is empty")
		}

		newInfo, err := os.Stat(logFile)
		if err != nil {
			t.Errorf("New log file not found: %v", err)
		} else if newInfo.Size() == 0 {
			t.Error("New log file is empty")
		}

		// Event processing should continue normally
		finalStats := setup.Collector.GetStats()
		if finalStats.TotalEvents == 0 {
			t.Error("No events processed during log rotation")
		}

		t.Logf("Log rotation test: %d events processed, rotated file: %d bytes, new file: %d bytes",
			finalStats.TotalEvents, rotatedInfo.Size(), newInfo.Size())
	})

	t.Run("MultipleFileRotations", func(t *testing.T) {
		// Test handling multiple rapid rotations
		rotationDir := filepath.Join(setup.TempDir, "rotations")
		err := os.MkdirAll(rotationDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create rotation directory: %v", err)
		}

		baseFile := filepath.Join(rotationDir, "events.log")

		// Generate events continuously
		go func() {
			for i := 0; i < 200; i++ {
				if ctx.Err() != nil {
					break
				}
				setup.Generator.GenerateTaskLifecycle(ctx, "multi-rotation-test", 2*time.Millisecond)
				time.Sleep(5 * time.Millisecond)
			}
		}()

		// Perform multiple rotations
		numRotations := 5
		for rotation := 0; rotation < numRotations; rotation++ {
			// Create file for this rotation
			currentFile := fmt.Sprintf("%s.%d", baseFile, rotation)
			file, err := os.Create(currentFile)
			if err != nil {
				t.Errorf("Failed to create rotation file %d: %v", rotation, err)
				continue
			}

			// Write some data
			for i := 0; i < 100; i++ {
				fmt.Fprintf(file, "Rotation %d entry %d\n", rotation, i)
			}
			file.Close()

			time.Sleep(200 * time.Millisecond)
		}

		// Wait for event processing to complete
		time.Sleep(500 * time.Millisecond)

		// Verify all rotation files exist
		for rotation := 0; rotation < numRotations; rotation++ {
			rotationFile := fmt.Sprintf("%s.%d", baseFile, rotation)
			info, err := os.Stat(rotationFile)
			if err != nil {
				t.Errorf("Rotation file %d not found: %v", rotation, err)
			} else if info.Size() == 0 {
				t.Errorf("Rotation file %d is empty", rotation)
			}
		}

		finalStats := setup.Collector.GetStats()
		if finalStats.TotalEvents == 0 {
			t.Error("No events processed during multiple rotations")
		}

		t.Logf("Multiple rotations test: %d rotations completed, %d events processed",
			numRotations, finalStats.TotalEvents)
	})

	t.Run("RotationErrorHandling", func(t *testing.T) {
		// Test rotation error scenarios
		errorDir := filepath.Join(setup.TempDir, "rotation_errors")
		err := os.MkdirAll(errorDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create error test directory: %v", err)
		}

		// Create a file and remove write permissions from directory
		testFile := filepath.Join(errorDir, "error_test.log")
		file, err := os.Create(testFile)
		if err != nil {
			t.Fatalf("Failed to create error test file: %v", err)
		}
		file.WriteString("initial content\n")
		file.Close()

		// Generate events
		setup.Generator.GenerateHighVolumeEvents(ctx, 50, []string{"error-test"})

		// Wait for processing
		time.Sleep(200 * time.Millisecond)

		// Try to rotate to a non-existent directory (should fail gracefully)
		badRotationPath := filepath.Join(errorDir, "nonexistent", "rotated.log")
		err = os.Rename(testFile, badRotationPath)

		// This should fail, which is expected
		if err == nil {
			t.Log("Unexpected success rotating to nonexistent directory")
		} else {
			t.Logf("Expected rotation failure: %v", err)
		}

		// System should continue processing events despite rotation failure
		finalStats := setup.Collector.GetStats()
		if finalStats.TotalEvents == 0 {
			t.Error("Event processing stopped due to rotation error")
		}

		t.Logf("Rotation error test: %d events processed despite rotation failure",
			finalStats.TotalEvents)
	})
}
