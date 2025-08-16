package main

import (
	"os"
	"testing"
	"time"

	"github.com/nsf/termbox-go"
	"go-ascii-calendar/events"
	"go-ascii-calendar/models"
	"go-ascii-calendar/storage"
	"go-ascii-calendar/terminal"
)

// Integration tests for the ASCII Calendar application
// These tests verify that all components work together correctly

func TestCompleteWorkflow(t *testing.T) {
	// Create temporary events file for testing
	tempFile := "test_integration_events.txt"
	defer os.Remove(tempFile)

	// Test complete workflow: create calendar, add events, navigate, view events
	t.Run("NavigationAndEventManagement", func(t *testing.T) {
		// Initialize components
		eventManager := events.NewManager()
		cal := models.NewCalendar()
		sel := models.NewSelection(cal)
		nc := terminal.NewNavigationController(cal, sel)

		// Set calendar to a known date for predictable testing
		cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
		sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)

		// Test navigation
		initialMonth := cal.CurrentMonth
		nc.NavigateMonthForward()
		if cal.CurrentMonth.Month() != initialMonth.AddDate(0, 1, 0).Month() {
			t.Errorf("Month navigation failed: expected %v, got %v",
				initialMonth.AddDate(0, 1, 0).Month(), cal.CurrentMonth.Month())
		}

		nc.NavigateMonthBackward()
		if cal.CurrentMonth.Month() != initialMonth.Month() {
			t.Errorf("Month navigation back failed: expected %v, got %v",
				initialMonth.Month(), cal.CurrentMonth.Month())
		}

		// Test day navigation
		initialDate := sel.SelectedDate
		nc.NavigateDayRight()
		expectedDate := initialDate.AddDate(0, 0, 1)
		if !sel.SelectedDate.Equal(expectedDate) {
			t.Errorf("Day navigation failed: expected %v, got %v", expectedDate, sel.SelectedDate)
		}

		// Test adding events
		testDate := time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC)
		err := eventManager.AddEvent(testDate, "10:30", "Integration test meeting")
		if err != nil {
			t.Fatalf("Failed to add event: %v", err)
		}

		err = eventManager.AddEvent(testDate, "14:00", "Another test event")
		if err != nil {
			t.Fatalf("Failed to add second event: %v", err)
		}

		// Verify events were added
		eventsForDate := eventManager.GetEventsForDate(testDate)
		if len(eventsForDate) != 2 {
			t.Errorf("Expected 2 events, got %d", len(eventsForDate))
		}

		// Verify events are sorted by time
		if eventsForDate[0].GetTimeString() != "10:30" || eventsForDate[1].GetTimeString() != "14:00" {
			t.Errorf("Events not sorted correctly: got %s, %s",
				eventsForDate[0].GetTimeString(), eventsForDate[1].GetTimeString())
		}

		// Test that events are detected for the date
		if !eventManager.HasEventsForDate(testDate) {
			t.Error("HasEventsForDate should return true for date with events")
		}

		// Test date without events
		emptyDate := time.Date(2025, time.August, 20, 0, 0, 0, 0, time.UTC)
		if eventManager.HasEventsForDate(emptyDate) {
			t.Error("HasEventsForDate should return false for date without events")
		}
	})
}

func TestFileIOOperations(t *testing.T) {
	tempFile := "test_integration_io.txt"
	defer os.Remove(tempFile)

	t.Run("EventPersistenceWorkflow", func(t *testing.T) {
		// Test complete file I/O workflow
		testEvents := []struct {
			date        time.Time
			timeStr     string
			description string
		}{
			{time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC), "09:00", "Morning meeting"},
			{time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC), "14:30", "Afternoon session"},
			{time.Date(2025, time.August, 17, 0, 0, 0, 0, time.UTC), "10:15", "Client call"},
		}

		// Save events to file
		for _, te := range testEvents {
			event := models.Event{
				Date:        te.date,
				Time:        time.Date(0, time.January, 1, 9, 0, 0, 0, time.UTC), // Will be overwritten
				Description: te.description,
			}
			// Parse time properly
			parsedTime, err := time.Parse("15:04", te.timeStr)
			if err != nil {
				t.Fatalf("Failed to parse time %s: %v", te.timeStr, err)
			}
			event.Time = parsedTime

			err = storage.SaveEventToFile(event, tempFile)
			if err != nil {
				t.Fatalf("Failed to save event: %v", err)
			}
		}

		// Load events from file
		loadedEvents, err := storage.LoadEventsFromFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to load events: %v", err)
		}

		if len(loadedEvents) != len(testEvents) {
			t.Errorf("Expected %d events, loaded %d", len(testEvents), len(loadedEvents))
		}

		// Verify event contents
		for i, te := range testEvents {
			if i >= len(loadedEvents) {
				break
			}
			loaded := loadedEvents[i]

			if loaded.GetDateString() != te.date.Format("2006-01-02") {
				t.Errorf("Event %d date mismatch: expected %s, got %s",
					i, te.date.Format("2006-01-02"), loaded.GetDateString())
			}

			if loaded.GetTimeString() != te.timeStr {
				t.Errorf("Event %d time mismatch: expected %s, got %s",
					i, te.timeStr, loaded.GetTimeString())
			}

			if loaded.Description != te.description {
				t.Errorf("Event %d description mismatch: expected %s, got %s",
					i, te.description, loaded.Description)
			}
		}
	})

	t.Run("ErrorHandlingScenarios", func(t *testing.T) {
		// Test loading from non-existent file
		events, err := storage.LoadEventsFromFile("non_existent_file.txt")
		if err != nil {
			t.Errorf("Loading from non-existent file should not error: %v", err)
		}
		if len(events) != 0 {
			t.Errorf("Loading from non-existent file should return empty slice, got %d events", len(events))
		}

		// Test malformed event data handling
		malformedFile := "test_malformed_integration.txt"
		defer os.Remove(malformedFile)

		malformedContent := `2025-08-16|09:00|Valid event
invalid line format
2025-08-16|25:00|Invalid time
2025-13-01|10:00|Invalid date
2025-08-17|11:00|Another valid event`

		err = os.WriteFile(malformedFile, []byte(malformedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create malformed test file: %v", err)
		}

		events, err = storage.LoadEventsFromFile(malformedFile)
		if err != nil {
			t.Fatalf("Loading malformed file should not error: %v", err)
		}

		// Should load only the 2 valid events
		if len(events) != 2 {
			t.Errorf("Expected 2 valid events from malformed file, got %d", len(events))
		}

		expectedDescriptions := []string{"Valid event", "Another valid event"}
		for i, expected := range expectedDescriptions {
			if i < len(events) && events[i].Description != expected {
				t.Errorf("Event %d description: expected %s, got %s",
					i, expected, events[i].Description)
			}
		}
	})
}

func TestTerminalCompatibility(t *testing.T) {
	t.Run("InputProcessing", func(t *testing.T) {
		term := terminal.NewTerminal()
		ih := terminal.NewInputHandler(term)

		// Test key processing workflow
		testCases := []struct {
			description string
			key         rune
			expected    terminal.KeyAction
		}{
			{"Navigation up", 'k', terminal.ActionMoveUp},
			{"Navigation down", 'j', terminal.ActionMoveDown},
			{"Navigation left", 'h', terminal.ActionMoveLeft},
			{"Navigation right", 'l', terminal.ActionMoveRight},
			{"Previous month", 'b', terminal.ActionMonthPrev},
			{"Next month", 'n', terminal.ActionMonthNext},
			{"Add event", 'a', terminal.ActionAddEvent},
			{"Quit", 'q', terminal.ActionQuit},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				// Test both lowercase and uppercase using proper termbox.Event
				lowerEvent := termbox.Event{Type: termbox.EventKey, Ch: tc.key}
				upperEvent := termbox.Event{Type: termbox.EventKey, Ch: tc.key - 32} // Convert to uppercase

				lowerAction := ih.ProcessKeyEvent(lowerEvent)
				upperAction := ih.ProcessKeyEvent(upperEvent)

				if lowerAction != tc.expected {
					t.Errorf("Lowercase %c: expected %v, got %v", tc.key, tc.expected, lowerAction)
				}

				if upperAction != tc.expected {
					t.Errorf("Uppercase %c: expected %v, got %v", tc.key-32, tc.expected, upperAction)
				}
			})
		}
	})

	t.Run("ValidationAndEdgeCases", func(t *testing.T) {
		// Test input validation
		term := terminal.NewTerminal()
		ih := terminal.NewInputHandler(term)

		// Valid keys should be recognized
		validKeys := "bBnNhHjJkKlLaAqQ"
		for _, key := range validKeys {
			if !ih.IsValidKey(key) {
				t.Errorf("Key %c should be valid", key)
			}
		}

		// Invalid keys should not be recognized
		invalidKeys := "xyz123@#$"
		for _, key := range invalidKeys {
			if ih.IsValidKey(key) {
				t.Errorf("Key %c should be invalid", key)
			}
		}
	})
}

func TestRequirementsValidation(t *testing.T) {
	t.Run("ThreeMonthCalendarRequirement", func(t *testing.T) {
		cal := models.NewCalendar()
		cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

		// Verify three-month view calculation
		prevMonth := cal.GetPreviousMonth()
		currMonth := cal.CurrentMonth
		nextMonth := cal.GetNextMonth()

		if prevMonth.Month() != time.July || prevMonth.Year() != 2025 {
			t.Errorf("Previous month should be July 2025, got %v %d", prevMonth.Month(), prevMonth.Year())
		}

		if currMonth.Month() != time.August || currMonth.Year() != 2025 {
			t.Errorf("Current month should be August 2025, got %v %d", currMonth.Month(), currMonth.Year())
		}

		if nextMonth.Month() != time.September || nextMonth.Year() != 2025 {
			t.Errorf("Next month should be September 2025, got %v %d", nextMonth.Month(), nextMonth.Year())
		}
	})

	t.Run("EventPersistenceRequirement", func(t *testing.T) {
		tempFile := "test_requirements_events.txt"
		defer os.Remove(tempFile)

		// Test the exact format specified: YYYY-MM-DD|HH:MM|description
		testEvent := models.Event{
			Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
			Time:        time.Date(0, time.January, 1, 14, 30, 0, 0, time.UTC),
			Description: "Test event with spaces and | pipes",
		}

		err := storage.SaveEventToFile(testEvent, tempFile)
		if err != nil {
			t.Fatalf("Failed to save event: %v", err)
		}

		// Read file content and verify format
		content, err := os.ReadFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		expectedFormat := "2025-08-16|14:30|Test event with spaces and | pipes\n"
		if string(content) != expectedFormat {
			t.Errorf("File format mismatch:\nExpected: %s\nGot: %s", expectedFormat, string(content))
		}

		// Test loading preserves format
		events, err := storage.LoadEventsFromFile(tempFile)
		if err != nil {
			t.Fatalf("Failed to load events: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(events))
		}

		loaded := events[0]
		if loaded.String() != "2025-08-16|14:30|Test event with spaces and | pipes" {
			t.Errorf("Event string format mismatch: %s", loaded.String())
		}
	})

	t.Run("NavigationRequirement", func(t *testing.T) {
		cal := models.NewCalendar()
		sel := models.NewSelection(cal)
		nc := terminal.NewNavigationController(cal, sel)

		// Set to August 2025 for predictable testing
		cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
		sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)

		// Test B/N keys (month navigation)
		nc.NavigateMonthBackward() // Should go to July
		if cal.CurrentMonth.Month() != time.July {
			t.Errorf("B key navigation failed: expected July, got %v", cal.CurrentMonth.Month())
		}

		nc.NavigateMonthForward() // Should go back to August
		nc.NavigateMonthForward() // Should go to September
		if cal.CurrentMonth.Month() != time.September {
			t.Errorf("N key navigation failed: expected September, got %v", cal.CurrentMonth.Month())
		}

		// Reset for day navigation testing
		nc.NavigateMonthBackward() // Back to August
		sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)

		// Test H/J/K/L keys (day navigation)
		initialDate := sel.SelectedDate

		nc.NavigateDayLeft() // H key
		if !sel.SelectedDate.Equal(initialDate.AddDate(0, 0, -1)) {
			t.Errorf("H key navigation failed: expected %v, got %v",
				initialDate.AddDate(0, 0, -1), sel.SelectedDate)
		}

		nc.NavigateDayRight() // L key (back to original)
		nc.NavigateDayRight() // L key (one day forward)
		if !sel.SelectedDate.Equal(initialDate.AddDate(0, 0, 1)) {
			t.Errorf("L key navigation failed: expected %v, got %v",
				initialDate.AddDate(0, 0, 1), sel.SelectedDate)
		}

		// Reset for week navigation
		sel.SelectedDate = initialDate
		nc.NavigateDayUp() // K key
		expectedUp := initialDate.AddDate(0, 0, -7)
		if nc.GetCurrentSelection().Day() != expectedUp.Day() {
			t.Errorf("K key navigation failed: expected day %d, got day %d",
				expectedUp.Day(), nc.GetCurrentSelection().Day())
		}

		sel.SelectedDate = initialDate
		nc.NavigateDayDown() // J key
		expectedDown := initialDate.AddDate(0, 0, 7)
		if nc.GetCurrentSelection().Day() != expectedDown.Day() {
			t.Errorf("J key navigation failed: expected day %d, got day %d",
				expectedDown.Day(), nc.GetCurrentSelection().Day())
		}
	})
}
