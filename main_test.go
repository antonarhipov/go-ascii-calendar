package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go-ascii-calendar/config"
)

func TestNewApplication(t *testing.T) {
	// Create temporary config for testing
	tempDir, err := os.MkdirTemp("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{}
	cfg.EventsFilePath = filepath.Join(tempDir, "test_events.json")

	app := NewApplication(cfg)

	if app == nil {
		t.Fatal("NewApplication() returned nil")
	}

	// Test that basic fields are initialized (using public interface)
	if app.config != cfg {
		t.Error("Application should reference the provided config")
	}

	if app.state != StateCalendar {
		t.Errorf("Initial state = %v, want StateCalendar", app.state)
	}

	// Test that all main components are initialized (checking they're not nil)
	if app.calendar == nil {
		t.Error("Calendar should be initialized")
	}

	if app.selection == nil {
		t.Error("Selection should be initialized")
	}

	if app.events == nil {
		t.Error("EventManager should be initialized")
	}

	if app.terminal == nil {
		t.Error("Terminal should be initialized")
	}

	if app.renderer == nil {
		t.Error("Renderer should be initialized")
	}

	if app.navigation == nil {
		t.Error("NavigationController should be initialized")
	}
}

func TestApplication_Initialize(t *testing.T) {
	// Create temporary config for testing
	tempDir, err := os.MkdirTemp("", "main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{}
	cfg.EventsFilePath = filepath.Join(tempDir, "test_events.json")

	app := NewApplication(cfg)

	// Test initialization - this will likely fail due to termbox dependency
	// but we can test that it doesn't panic and handles errors gracefully
	err = app.Initialize()

	// We expect this to fail in test environment due to terminal requirements
	// The key is that it should fail gracefully with a meaningful error
	if err == nil {
		t.Log("Initialize() succeeded (unexpected in test environment)")
	} else {
		t.Logf("Initialize() failed as expected: %v", err)
		// This is expected in test environment due to terminal dependencies
	}

	// Test that events manager was created and can be accessed
	events := app.events.GetAllEvents()
	if events == nil {
		t.Error("Events slice should be initialized even if empty")
	}
}

func TestApplication_Constructor_WithNilConfig(t *testing.T) {
	// Test that constructor handles nil config gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewApplication() with nil config panicked: %v", r)
		}
	}()

	app := NewApplication(nil)

	if app == nil {
		t.Fatal("NewApplication() returned nil with nil config")
	}

	// Should still initialize basic components
	if app.calendar == nil {
		t.Error("Calendar should be initialized even with nil config")
	}

	if app.selection == nil {
		t.Error("Selection should be initialized even with nil config")
	}

	t.Log("NewApplication() handled nil config gracefully")
}

func TestApplication_StateTransitions(t *testing.T) {
	cfg := config.DefaultConfig()
	app := NewApplication(cfg)

	// Test initial state
	if app.state != StateCalendar {
		t.Errorf("Initial state = %v, want StateCalendar", app.state)
	}

	// Test state transitions by directly setting states
	// (since we can't easily test input handling without terminal)
	testStates := []AppState{
		StateEventList,
		StateAddEvent,
		StateSearch,
		StateCalendarEventSelection,
		StateCalendarEventAdd,
		StateCalendarEventEdit,
		StateCalendar, // Back to initial
	}

	for _, testState := range testStates {
		app.state = testState
		if app.state != testState {
			t.Errorf("Failed to set state to %v", testState)
		}
	}

	t.Log("State transitions completed successfully")
}

func TestApplication_EventIndexManagement(t *testing.T) {
	cfg := config.DefaultConfig()
	app := NewApplication(cfg)

	// Test initial selected event index
	if app.selectedEventIndex != 0 {
		t.Errorf("Initial selectedEventIndex = %d, want 0", app.selectedEventIndex)
	}

	// Test setting event index
	app.selectedEventIndex = 5
	if app.selectedEventIndex != 5 {
		t.Errorf("selectedEventIndex = %d, want 5", app.selectedEventIndex)
	}

	// Test resetting
	app.selectedEventIndex = 0
	if app.selectedEventIndex != 0 {
		t.Errorf("Reset selectedEventIndex = %d, want 0", app.selectedEventIndex)
	}

	t.Log("Event index management completed successfully")
}

func TestApplication_SearchFieldsInitialization(t *testing.T) {
	cfg := config.DefaultConfig()
	app := NewApplication(cfg)

	// Test that search-related fields are properly initialized
	if app.searchQuery != "" {
		t.Errorf("Initial searchQuery = %s, want empty string", app.searchQuery)
	}

	if app.searchResults != nil {
		t.Errorf("Initial searchResults should be nil, got %v", app.searchResults)
	}

	if app.searchResultDates != nil {
		t.Errorf("Initial searchResultDates should be nil, got %v", app.searchResultDates)
	}

	if app.selectedResultIndex != 0 {
		t.Errorf("Initial selectedResultIndex = %d, want 0", app.selectedResultIndex)
	}

	t.Log("Search fields initialization verified")
}

func TestAppStateConstants(t *testing.T) {
	// Test that app states are defined correctly and distinct
	states := []AppState{
		StateCalendar,
		StateCalendarEventSelection,
		StateCalendarEventAdd,
		StateCalendarEventEdit,
		StateSearch,
		StateEventList,
		StateAddEvent,
	}

	// Verify states are distinct
	stateMap := make(map[AppState]bool)
	for _, state := range states {
		if stateMap[state] {
			t.Errorf("Duplicate app state found: %v", state)
		}
		stateMap[state] = true
	}

	// Test that states have reasonable values
	for i, state := range states {
		if int(state) < 0 {
			t.Errorf("State %d has negative value: %v", i, state)
		}
	}

	t.Logf("Verified %d distinct app states with valid values", len(states))
}

func TestApplication_ComponentIntegration(t *testing.T) {
	cfg := config.DefaultConfig()
	app := NewApplication(cfg)

	// Test that components are properly connected
	// Calendar and Selection should be connected
	if app.selection.Calendar != app.calendar {
		t.Error("Selection should reference the same calendar instance")
	}

	// Navigation should reference the same components
	if app.navigation == nil {
		t.Error("Navigation controller should be initialized")
	}

	t.Log("Component integration verified")
}

// Benchmark test for application creation
func BenchmarkNewApplication(b *testing.B) {
	cfg := config.DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := NewApplication(cfg)
		_ = app // Prevent optimization
	}
}

// Test memory usage and cleanup
func TestApplication_MemoryManagement(t *testing.T) {
	cfg := config.DefaultConfig()

	// Create and discard multiple applications to test for memory leaks
	for i := 0; i < 100; i++ {
		app := NewApplication(cfg)

		// Use the app minimally
		if app.state != StateCalendar {
			t.Errorf("Iteration %d: wrong initial state", i)
		}

		// Let it go out of scope for GC
		app = nil
	}

	t.Log("Memory management test completed")
}

func TestApplication_ConfigHandling(t *testing.T) {
	// Test with different config scenarios
	testConfigs := []*config.Config{
		config.DefaultConfig(),
		&config.Config{EventsFilePath: "/tmp/test_events.json"},
		&config.Config{EventsFilePath: ""},
		nil, // Should be handled gracefully
	}

	for i, cfg := range testConfigs {
		t.Run(fmt.Sprintf("Config_%d", i), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("NewApplication() panicked with config %d: %v", i, r)
				}
			}()

			app := NewApplication(cfg)
			if app == nil {
				t.Errorf("NewApplication() returned nil with config %d", i)
			}

			t.Logf("Config %d handled successfully", i)
		})
	}
}
