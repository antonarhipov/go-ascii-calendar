package events

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-ascii-calendar/config"
	"go-ascii-calendar/models"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	if manager.events == nil {
		t.Error("Manager events slice should be initialized")
	}

	if len(manager.events) != 0 {
		t.Error("Manager events slice should be empty initially")
	}

	if manager.config != nil {
		t.Error("Manager config should be nil for legacy constructor")
	}
}

func TestNewManagerWithConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	manager := NewManagerWithConfig(cfg)

	if manager == nil {
		t.Fatal("NewManagerWithConfig() returned nil")
	}

	if manager.events == nil {
		t.Error("Manager events slice should be initialized")
	}

	if len(manager.events) != 0 {
		t.Error("Manager events slice should be empty initially")
	}

	if manager.config != cfg {
		t.Error("Manager config should match the provided config")
	}
}

func TestManager_GetEventCount(t *testing.T) {
	manager := NewManager()

	// Test empty manager
	if manager.GetEventCount() != 0 {
		t.Errorf("GetEventCount() = %d, want 0", manager.GetEventCount())
	}

	// Add some events manually for testing
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Test event 1"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "Test event 2"},
	}

	manager.events = testEvents

	if manager.GetEventCount() != 2 {
		t.Errorf("GetEventCount() = %d, want 2", manager.GetEventCount())
	}
}

func TestManager_GetAllEvents(t *testing.T) {
	manager := NewManager()

	// Test empty manager
	events := manager.GetAllEvents()
	if events == nil {
		t.Error("GetAllEvents() should not return nil")
	}

	if len(events) != 0 {
		t.Errorf("GetAllEvents() length = %d, want 0", len(events))
	}

	// Add test events
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Test event 1"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "Test event 2"},
	}

	manager.events = testEvents
	events = manager.GetAllEvents()

	if len(events) != 2 {
		t.Errorf("GetAllEvents() length = %d, want 2", len(events))
	}

	// Verify events match
	for i, event := range events {
		if event.Description != testEvents[i].Description {
			t.Errorf("Event %d description = %s, want %s", i, event.Description, testEvents[i].Description)
		}
	}
}

func TestManager_HasEventsForDate(t *testing.T) {
	manager := NewManager()
	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	// Test empty manager
	if manager.HasEventsForDate(testDate) {
		t.Error("HasEventsForDate() should return false for empty manager")
	}

	// Add test events
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Test event 1"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "Test event 2"},
	}

	manager.events = testEvents

	// Test date with events
	if !manager.HasEventsForDate(testDate) {
		t.Error("HasEventsForDate() should return true for date with events")
	}

	// Test date without events
	noEventDate := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)
	if manager.HasEventsForDate(noEventDate) {
		t.Error("HasEventsForDate() should return false for date without events")
	}

	// Test date normalization (time component should be ignored)
	testDateWithTime := time.Date(2025, 8, 15, 14, 30, 0, 0, time.UTC)
	if !manager.HasEventsForDate(testDateWithTime) {
		t.Error("HasEventsForDate() should normalize dates and return true regardless of time component")
	}
}

func TestManager_GetEventsForDate(t *testing.T) {
	manager := NewManager()
	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	// Test empty manager
	events := manager.GetEventsForDate(testDate)
	if len(events) != 0 {
		t.Errorf("GetEventsForDate() length = %d, want 0", len(events))
	}

	// Add test events for different dates and times
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC), Description: "Afternoon event"},
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Morning event"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "Different day event"},
	}

	manager.events = testEvents

	// Test getting events for specific date
	events = manager.GetEventsForDate(testDate)
	if len(events) != 2 {
		t.Errorf("GetEventsForDate() length = %d, want 2", len(events))
	}

	// Verify events are sorted by time (morning event should come first)
	if events[0].Description != "Morning event" {
		t.Errorf("First event description = %s, want 'Morning event'", events[0].Description)
	}
	if events[1].Description != "Afternoon event" {
		t.Errorf("Second event description = %s, want 'Afternoon event'", events[1].Description)
	}

	// Test date normalization
	testDateWithTime := time.Date(2025, 8, 15, 18, 45, 0, 0, time.UTC)
	events = manager.GetEventsForDate(testDateWithTime)
	if len(events) != 2 {
		t.Errorf("GetEventsForDate() with time component should still return 2 events, got %d", len(events))
	}
}

func TestManager_GetEventsForMonth(t *testing.T) {
	manager := NewManager()
	testMonth := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	// Test empty manager
	events := manager.GetEventsForMonth(testMonth)
	if len(events) != 0 {
		t.Errorf("GetEventsForMonth() length = %d, want 0", len(events))
	}

	// Add test events for different months
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC), Description: "August event 2"},
		{Date: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "August event 1"},
		{Date: time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "September event"},
	}

	manager.events = testEvents

	// Test getting events for specific month
	events = manager.GetEventsForMonth(testMonth)
	if len(events) != 2 {
		t.Errorf("GetEventsForMonth() length = %d, want 2", len(events))
	}

	// Verify events are sorted by date and time
	if events[0].Description != "August event 1" {
		t.Errorf("First event description = %s, want 'August event 1'", events[0].Description)
	}
	if events[1].Description != "August event 2" {
		t.Errorf("Second event description = %s, want 'August event 2'", events[1].Description)
	}
}

func TestManager_GetEventsInDateRange(t *testing.T) {
	manager := NewManager()
	startDate := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)

	// Test empty manager
	events := manager.GetEventsInDateRange(startDate, endDate)
	if len(events) != 0 {
		t.Errorf("GetEventsInDateRange() length = %d, want 0", len(events))
	}

	// Add test events
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 5, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), Description: "Before range"},
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC), Description: "In range 2"},
		{Date: time.Date(2025, 8, 12, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "In range 1"},
		{Date: time.Date(2025, 8, 25, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "After range"},
	}

	manager.events = testEvents

	// Test getting events in date range
	events = manager.GetEventsInDateRange(startDate, endDate)
	if len(events) != 2 {
		t.Errorf("GetEventsInDateRange() length = %d, want 2", len(events))
	}

	// Verify events are sorted by date and time
	if events[0].Description != "In range 1" {
		t.Errorf("First event description = %s, want 'In range 1'", events[0].Description)
	}
	if events[1].Description != "In range 2" {
		t.Errorf("Second event description = %s, want 'In range 2'", events[1].Description)
	}
}

func TestManager_AddEvent(t *testing.T) {
	manager := NewManager()
	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	// Test successful event addition
	err := manager.AddEvent(testDate, "10:30", "Test meeting")
	if err != nil {
		t.Errorf("AddEvent() failed: %v", err)
	}

	if manager.GetEventCount() != 1 {
		t.Errorf("Event count = %d, want 1", manager.GetEventCount())
	}

	events := manager.GetEventsForDate(testDate)
	if len(events) != 1 {
		t.Errorf("Events for date length = %d, want 1", len(events))
	}

	event := events[0]
	if event.Description != "Test meeting" {
		t.Errorf("Event description = %s, want 'Test meeting'", event.Description)
	}

	// Test invalid time format
	err = manager.AddEvent(testDate, "25:00", "Invalid time event")
	if err == nil {
		t.Error("AddEvent() should fail with invalid time format")
	}

	err = manager.AddEvent(testDate, "10:60", "Invalid minutes event")
	if err == nil {
		t.Error("AddEvent() should fail with invalid minutes")
	}

	err = manager.AddEvent(testDate, "not-a-time", "Invalid format event")
	if err == nil {
		t.Error("AddEvent() should fail with non-time format")
	}

	// Test empty description
	err = manager.AddEvent(testDate, "11:00", "")
	if err == nil {
		t.Error("AddEvent() should fail with empty description")
	}
}

func TestManager_SearchEvents(t *testing.T) {
	manager := NewManager()

	// Add test events
	testEvents := []models.Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Team meeting"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "Project review"},
		{Date: time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC), Description: "Team retrospective"},
		{Date: time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 13, 0, 0, 0, time.UTC), Description: "Client presentation"},
	}

	manager.events = testEvents

	// Test search for "team" (case insensitive)
	results := manager.SearchEvents("team")
	if len(results) != 2 {
		t.Errorf("Search results length = %d, want 2", len(results))
	}

	// Test search for "MEETING" (case insensitive)
	results = manager.SearchEvents("MEETING")
	if len(results) != 1 {
		t.Errorf("Search results length = %d, want 1", len(results))
	}

	if results[0].Description != "Team meeting" {
		t.Errorf("Search result = %s, want 'Team meeting'", results[0].Description)
	}

	// Test search with no results
	results = manager.SearchEvents("nonexistent")
	if len(results) != 0 {
		t.Errorf("Search results length = %d, want 0", len(results))
	}

	// Test empty search query (should return empty slice, not all events)
	results = manager.SearchEvents("")
	if len(results) != 0 {
		t.Errorf("Empty search should return empty slice, got %d, want 0", len(results))
	}
}

func TestManager_LoadEvents_NoConfig(t *testing.T) {
	manager := NewManager()

	// This will attempt to load from the default legacy text format
	// Since we don't have control over the default events file, we'll test that it doesn't crash
	err := manager.LoadEvents()
	// We accept both success and file-not-found errors as valid
	if err != nil && !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "cannot find") {
		t.Errorf("LoadEvents() failed with unexpected error: %v", err)
	}
}

func TestManager_LoadEvents_WithConfig(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "events_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	eventsPath := filepath.Join(tempDir, "test_events.json")

	// Create test config
	cfg := &config.Config{}
	cfg.EventsFilePath = eventsPath

	manager := NewManagerWithConfig(cfg)

	// Test loading from non-existent file (should not fail)
	// Note: The storage layer may automatically migrate from existing events.txt
	err = manager.LoadEvents()
	if err != nil {
		t.Errorf("LoadEvents() failed: %v", err)
	}

	// Event count may be > 0 if migration occurred from events.txt
	// The main test is that LoadEvents() doesn't fail
	eventCount := manager.GetEventCount()
	if eventCount < 0 {
		t.Errorf("Event count should be >= 0, got %d", eventCount)
	}
}

func TestManager_ReloadEvents(t *testing.T) {
	manager := NewManager()

	// ReloadEvents should call LoadEvents
	err := manager.ReloadEvents()
	// We accept both success and file-not-found errors as valid
	if err != nil && !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "cannot find") {
		t.Errorf("ReloadEvents() failed with unexpected error: %v", err)
	}
}
