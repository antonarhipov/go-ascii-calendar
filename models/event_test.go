package models

import (
	"testing"
	"time"
)

func TestEvent_GetTimeString(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Morning time",
			time:     time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC),
			expected: "09:30",
		},
		{
			name:     "Afternoon time",
			time:     time.Date(0, 1, 1, 14, 45, 0, 0, time.UTC),
			expected: "14:45",
		},
		{
			name:     "Midnight",
			time:     time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "00:00",
		},
		{
			name:     "Late night",
			time:     time.Date(0, 1, 1, 23, 59, 0, 0, time.UTC),
			expected: "23:59",
		},
		{
			name:     "Noon",
			time:     time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
			expected: "12:00",
		},
		{
			name:     "Single digit hour and minute",
			time:     time.Date(0, 1, 1, 8, 5, 0, 0, time.UTC),
			expected: "08:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				Date:        time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
				Time:        tt.time,
				Description: "Test event",
			}

			result := event.GetTimeString()
			if result != tt.expected {
				t.Errorf("GetTimeString() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestEvent_GetDateString(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "Regular date",
			date:     time.Date(2025, 8, 15, 10, 30, 0, 0, time.UTC),
			expected: "2025-08-15",
		},
		{
			name:     "New Year's Day",
			date:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2025-01-01",
		},
		{
			name:     "December 31st",
			date:     time.Date(2024, 12, 31, 23, 59, 0, 0, time.UTC),
			expected: "2024-12-31",
		},
		{
			name:     "Leap year date",
			date:     time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			expected: "2024-02-29",
		},
		{
			name:     "Single digit month and day",
			date:     time.Date(2025, 3, 7, 15, 45, 0, 0, time.UTC),
			expected: "2025-03-07",
		},
		{
			name:     "Date with different timezone",
			date:     time.Date(2025, 6, 20, 8, 0, 0, 0, time.FixedZone("PST", -8*3600)),
			expected: "2025-06-20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &Event{
				Date:        tt.date,
				Time:        time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
				Description: "Test event",
			}

			result := event.GetDateString()
			if result != tt.expected {
				t.Errorf("GetDateString() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestEvent_String(t *testing.T) {
	tests := []struct {
		name        string
		event       Event
		expected    string
		description string
	}{
		{
			name: "Regular event",
			event: Event{
				Date:        time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC),
				Description: "Team meeting",
			},
			expected:    "2025-08-15|14:30|Team meeting",
			description: "Standard event formatting",
		},
		{
			name: "Morning event",
			event: Event{
				Date:        time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC),
				Description: "Christmas morning",
			},
			expected:    "2025-12-25|09:00|Christmas morning",
			description: "Morning event with single digit hour",
		},
		{
			name: "Midnight event",
			event: Event{
				Date:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
				Description: "New Year celebration",
			},
			expected:    "2026-01-01|00:00|New Year celebration",
			description: "Midnight event",
		},
		{
			name: "Event with special characters",
			event: Event{
				Date:        time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 18, 45, 0, 0, time.UTC),
				Description: "BBQ & fireworks: Independence Day!",
			},
			expected:    "2025-07-04|18:45|BBQ & fireworks: Independence Day!",
			description: "Event with special characters in description",
		},
		{
			name: "Empty description",
			event: Event{
				Date:        time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 12, 30, 0, 0, time.UTC),
				Description: "",
			},
			expected:    "2025-05-10|12:30|",
			description: "Event with empty description",
		},
		{
			name: "Late night event",
			event: Event{
				Date:        time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 23, 59, 0, 0, time.UTC),
				Description: "Halloween party ends",
			},
			expected:    "2025-10-31|23:59|Halloween party ends",
			description: "Late night event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.String()
			if result != tt.expected {
				t.Errorf("String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestEvent_StringConsistency(t *testing.T) {
	// Test that String() method is consistent with GetDateString() and GetTimeString()
	event := Event{
		Date:        time.Date(2025, 8, 15, 10, 20, 30, 0, time.UTC),
		Time:        time.Date(2023, 5, 3, 14, 45, 12, 0, time.UTC),
		Description: "Consistency test",
	}

	expected := event.GetDateString() + "|" + event.GetTimeString() + "|" + event.Description
	result := event.String()

	if result != expected {
		t.Errorf("String() = %s, want %s", result, expected)
		t.Errorf("String() should be consistent with GetDateString() + GetTimeString() + Description")
	}

	// Verify individual components
	if event.GetDateString() != "2025-08-15" {
		t.Errorf("GetDateString() = %s, want 2025-08-15", event.GetDateString())
	}

	if event.GetTimeString() != "14:45" {
		t.Errorf("GetTimeString() = %s, want 14:45", event.GetTimeString())
	}
}

func TestEvent_TimezoneHandling(t *testing.T) {
	// Test that formatting works correctly with different timezones
	est := time.FixedZone("EST", -5*3600)
	pst := time.FixedZone("PST", -8*3600)

	tests := []struct {
		name     string
		date     time.Time
		time     time.Time
		expected string
	}{
		{
			name:     "EST timezone",
			date:     time.Date(2025, 8, 15, 10, 0, 0, 0, est),
			time:     time.Date(0, 1, 1, 14, 30, 0, 0, est),
			expected: "2025-08-15|14:30|EST event",
		},
		{
			name:     "PST timezone",
			date:     time.Date(2025, 8, 15, 7, 0, 0, 0, pst),
			time:     time.Date(0, 1, 1, 9, 15, 0, 0, pst),
			expected: "2025-08-15|09:15|PST event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := Event{
				Date:        tt.date,
				Time:        tt.time,
				Description: tt.name[:3] + " event",
			}

			result := event.String()
			if result != tt.expected {
				t.Errorf("String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestEvent_EdgeCases(t *testing.T) {
	// Test edge cases and boundary conditions
	tests := []struct {
		name  string
		event Event
	}{
		{
			name: "Very long description",
			event: Event{
				Date:        time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				Description: "This is a very long event description that might be used to test how the Event model handles lengthy text content without any issues or truncation",
			},
		},
		{
			name: "Description with pipe characters",
			event: Event{
				Date:        time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
				Time:        time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC),
				Description: "Event with | pipe | characters",
			},
		},
		{
			name: "Zero time values",
			event: Event{
				Date:        time.Time{},
				Time:        time.Time{},
				Description: "Zero time event",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that methods don't panic and return some result
			dateStr := tt.event.GetDateString()
			timeStr := tt.event.GetTimeString()
			fullStr := tt.event.String()

			// Basic sanity checks
			if len(dateStr) == 0 {
				t.Error("GetDateString() returned empty string")
			}
			if len(timeStr) == 0 {
				t.Error("GetTimeString() returned empty string")
			}
			if len(fullStr) == 0 {
				t.Error("String() returned empty string")
			}

			// Verify String() contains the description
			if tt.event.Description != "" && !contains(fullStr, tt.event.Description) {
				t.Errorf("String() = %s should contain description %s", fullStr, tt.event.Description)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOf(s, substr) >= 0)
}

// Helper function to find index of substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
