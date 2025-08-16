package storage

import (
	"os"
	"testing"
	"time"

	"go-ascii-calendar/models"
)

func TestValidateEvent(t *testing.T) {
	tests := []struct {
		name      string
		event     models.Event
		expectErr bool
	}{
		{
			name: "Valid event",
			event: models.Event{
				Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, time.January, 1, 9, 30, 0, 0, time.UTC),
				Description: "Team meeting",
			},
			expectErr: false,
		},
		{
			name: "Empty description",
			event: models.Event{
				Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, time.January, 1, 9, 30, 0, 0, time.UTC),
				Description: "",
			},
			expectErr: true,
		},
		{
			name: "Whitespace-only description",
			event: models.Event{
				Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, time.January, 1, 9, 30, 0, 0, time.UTC),
				Description: "   ",
			},
			expectErr: true,
		},
		{
			name: "Valid event with spaces in description",
			event: models.Event{
				Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, time.January, 1, 14, 45, 0, 0, time.UTC),
				Description: "Client presentation with team",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEvent(tt.event)
			if tt.expectErr && err == nil {
				t.Errorf("ValidateEvent() expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("ValidateEvent() unexpected error: %v", err)
			}
		})
	}
}

func TestParseEventLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		expectErr    bool
		expectedDate string
		expectedTime string
		expectedDesc string
	}{
		{
			name:         "Valid event line",
			line:         "2025-08-16|09:30|Team meeting",
			expectErr:    false,
			expectedDate: "2025-08-16",
			expectedTime: "09:30",
			expectedDesc: "Team meeting",
		},
		{
			name:         "Valid event with spaces",
			line:         "2025-08-17|14:45|Client presentation with stakeholders",
			expectErr:    false,
			expectedDate: "2025-08-17",
			expectedTime: "14:45",
			expectedDesc: "Client presentation with stakeholders",
		},
		{
			name:      "Missing description",
			line:      "2025-08-16|09:30|",
			expectErr: true,
		},
		{
			name:      "Missing time",
			line:      "2025-08-16||Team meeting",
			expectErr: true,
		},
		{
			name:      "Missing date",
			line:      "|09:30|Team meeting",
			expectErr: true,
		},
		{
			name:      "Invalid date format",
			line:      "25-08-16|09:30|Team meeting",
			expectErr: true,
		},
		{
			name:      "Invalid time format",
			line:      "2025-08-16|25:30|Team meeting",
			expectErr: true,
		},
		{
			name:      "Wrong number of separators",
			line:      "2025-08-16|09:30",
			expectErr: true,
		},
		{
			name:         "Too many separators",
			line:         "2025-08-16|09:30|Team|meeting",
			expectErr:    false, // Should work - description can contain pipes
			expectedDate: "2025-08-16",
			expectedTime: "09:30",
			expectedDesc: "Team|meeting",
		},
		{
			name:      "Empty line",
			line:      "",
			expectErr: true,
		},
		{
			name:      "Invalid leap day",
			line:      "2025-02-29|09:30|Invalid date",
			expectErr: true,
		},
		{
			name:         "Valid leap day",
			line:         "2024-02-29|09:30|Leap day event",
			expectErr:    false,
			expectedDate: "2024-02-29",
			expectedTime: "09:30",
			expectedDesc: "Leap day event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := ParseEventLine(tt.line)

			if tt.expectErr && err == nil {
				t.Errorf("ParseEventLine(%s) expected error but got none", tt.line)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("ParseEventLine(%s) unexpected error: %v", tt.line, err)
			}

			if !tt.expectErr {
				if event.GetDateString() != tt.expectedDate {
					t.Errorf("ParseEventLine(%s) date = %s, want %s", tt.line, event.GetDateString(), tt.expectedDate)
				}
				if event.GetTimeString() != tt.expectedTime {
					t.Errorf("ParseEventLine(%s) time = %s, want %s", tt.line, event.GetTimeString(), tt.expectedTime)
				}
				if event.Description != tt.expectedDesc {
					t.Errorf("ParseEventLine(%s) description = %s, want %s", tt.line, event.Description, tt.expectedDesc)
				}
			}
		})
	}
}

func TestSaveAndLoadEvents(t *testing.T) {
	// Create a temporary test file
	tempFile := "test_events.txt"

	// Clean up any existing test file
	os.Remove(tempFile)
	defer os.Remove(tempFile)

	// Test saving an event
	testEvent := models.Event{
		Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
		Time:        time.Date(0, time.January, 1, 9, 30, 0, 0, time.UTC),
		Description: "Test event",
	}

	err := SaveEventToFile(testEvent, tempFile)
	if err != nil {
		t.Fatalf("SaveEventToFile() failed: %v", err)
	}

	// Test loading events
	events, err := LoadEventsFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadEventsFromFile() failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("LoadEventsFromFile() returned %d events, want 1", len(events))
	}

	loadedEvent := events[0]
	if loadedEvent.GetDateString() != testEvent.GetDateString() {
		t.Errorf("Loaded event date = %s, want %s", loadedEvent.GetDateString(), testEvent.GetDateString())
	}

	if loadedEvent.GetTimeString() != testEvent.GetTimeString() {
		t.Errorf("Loaded event time = %s, want %s", loadedEvent.GetTimeString(), testEvent.GetTimeString())
	}

	if loadedEvent.Description != testEvent.Description {
		t.Errorf("Loaded event description = %s, want %s", loadedEvent.Description, testEvent.Description)
	}
}

func TestSaveMultipleEvents(t *testing.T) {
	// Create a temporary test file
	tempFile := "test_multiple_events.txt"

	// Clean up any existing test file
	os.Remove(tempFile)
	defer os.Remove(tempFile)

	// Test saving multiple events
	events := []models.Event{
		{
			Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
			Time:        time.Date(0, time.January, 1, 9, 30, 0, 0, time.UTC),
			Description: "First event",
		},
		{
			Date:        time.Date(2025, time.August, 16, 0, 0, 0, 0, time.UTC),
			Time:        time.Date(0, time.January, 1, 14, 45, 0, 0, time.UTC),
			Description: "Second event",
		},
		{
			Date:        time.Date(2025, time.August, 17, 0, 0, 0, 0, time.UTC),
			Time:        time.Date(0, time.January, 1, 10, 0, 0, 0, time.UTC),
			Description: "Third event",
		},
	}

	for _, event := range events {
		err := SaveEventToFile(event, tempFile)
		if err != nil {
			t.Fatalf("SaveEventToFile() failed: %v", err)
		}
	}

	// Load and verify all events
	loadedEvents, err := LoadEventsFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadEventsFromFile() failed: %v", err)
	}

	if len(loadedEvents) != len(events) {
		t.Fatalf("LoadEventsFromFile() returned %d events, want %d", len(loadedEvents), len(events))
	}

	for i, expected := range events {
		loaded := loadedEvents[i]
		if loaded.GetDateString() != expected.GetDateString() ||
			loaded.GetTimeString() != expected.GetTimeString() ||
			loaded.Description != expected.Description {
			t.Errorf("Event %d mismatch: got %s|%s|%s, want %s|%s|%s",
				i, loaded.GetDateString(), loaded.GetTimeString(), loaded.Description,
				expected.GetDateString(), expected.GetTimeString(), expected.Description)
		}
	}
}

func TestLoadEventsWithMalformedLines(t *testing.T) {
	// Create a test file with malformed lines
	tempFile := "test_malformed_events.txt"

	// Clean up any existing test file
	os.Remove(tempFile)
	defer os.Remove(tempFile)

	// Create test file with both valid and malformed lines
	content := `2025-08-16|09:30|Valid event
invalid line without pipes
2025-08-16|25:30|Invalid time
|09:30|Missing date
2025-08-16||Missing description
2025-08-17|14:45|Another valid event

2025-13-01|10:00|Invalid date`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Load events - should skip malformed lines and continue
	events, err := LoadEventsFromFile(tempFile)
	if err != nil {
		t.Fatalf("LoadEventsFromFile() failed: %v", err)
	}

	// Should have loaded only the 2 valid events
	if len(events) != 2 {
		t.Errorf("LoadEventsFromFile() returned %d events, want 2 (malformed lines should be skipped)", len(events))
	}

	// Verify the valid events were loaded correctly
	expectedEvents := []string{
		"Valid event",
		"Another valid event",
	}

	for i, expected := range expectedEvents {
		if i < len(events) && events[i].Description != expected {
			t.Errorf("Event %d description = %s, want %s", i, events[i].Description, expected)
		}
	}
}

func TestLoadEventsNonExistentFile(t *testing.T) {
	// Loading from non-existent file should return empty slice, not error
	events, err := LoadEventsFromFile("non_existent_file.txt")
	if err != nil {
		t.Errorf("LoadEventsFromFile() from non-existent file should not return error, got: %v", err)
	}

	if len(events) != 0 {
		t.Errorf("LoadEventsFromFile() from non-existent file should return empty slice, got %d events", len(events))
	}
}

func TestFileExistsAtPath(t *testing.T) {
	// Test with existing file
	tempFile := "test_exists.txt"
	err := os.WriteFile(tempFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tempFile)

	if !FileExistsAtPath(tempFile) {
		t.Error("FileExistsAtPath() should return true for existing file")
	}

	// Test with non-existent file
	if FileExistsAtPath("non_existent_file.txt") {
		t.Error("FileExistsAtPath() should return false for non-existent file")
	}
}

func TestCreateEventFileAtPath(t *testing.T) {
	tempFile := "test_create.txt"
	defer os.Remove(tempFile)

	// Ensure file doesn't exist
	os.Remove(tempFile)

	err := CreateEventFileAtPath(tempFile)
	if err != nil {
		t.Fatalf("CreateEventFileAtPath() failed: %v", err)
	}

	if !FileExistsAtPath(tempFile) {
		t.Error("CreateEventFileAtPath() should create the file")
	}

	// Test creating file that already exists (should not error)
	err = CreateEventFileAtPath(tempFile)
	if err != nil {
		t.Errorf("CreateEventFileAtPath() should not error when file already exists: %v", err)
	}
}
