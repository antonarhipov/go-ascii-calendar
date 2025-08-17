package models

import (
	"testing"
	"time"
)

func TestNewCalendar(t *testing.T) {
	calendar := NewCalendar()

	if calendar == nil {
		t.Fatal("NewCalendar() returned nil")
	}

	// Verify CurrentMonth is set to first day of current month
	now := time.Now()
	expectedMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	if !calendar.CurrentMonth.Equal(expectedMonth) {
		t.Errorf("CurrentMonth = %v, want %v", calendar.CurrentMonth, expectedMonth)
	}

	// Verify Events slice is initialized
	if calendar.Events == nil {
		t.Error("Events slice should be initialized")
	}

	if len(calendar.Events) != 0 {
		t.Errorf("Events slice should be empty initially, got %d events", len(calendar.Events))
	}

	// Verify it's set to the first day of the month (time components should be zero)
	if calendar.CurrentMonth.Day() != 1 {
		t.Errorf("CurrentMonth day = %d, want 1", calendar.CurrentMonth.Day())
	}

	if calendar.CurrentMonth.Hour() != 0 || calendar.CurrentMonth.Minute() != 0 || calendar.CurrentMonth.Second() != 0 {
		t.Error("CurrentMonth should have zero time components")
	}
}

func TestCalendar_GetPreviousMonth(t *testing.T) {
	calendar := NewCalendar()

	// Set a known month for testing
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	previousMonth := calendar.GetPreviousMonth()
	expected := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	if !previousMonth.Equal(expected) {
		t.Errorf("GetPreviousMonth() = %v, want %v", previousMonth, expected)
	}

	// Test year boundary crossing
	calendar.CurrentMonth = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousMonth = calendar.GetPreviousMonth()
	expected = time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)

	if !previousMonth.Equal(expected) {
		t.Errorf("GetPreviousMonth() across year boundary = %v, want %v", previousMonth, expected)
	}
}

func TestCalendar_GetNextMonth(t *testing.T) {
	calendar := NewCalendar()

	// Set a known month for testing
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	nextMonth := calendar.GetNextMonth()
	expected := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)

	if !nextMonth.Equal(expected) {
		t.Errorf("GetNextMonth() = %v, want %v", nextMonth, expected)
	}

	// Test year boundary crossing
	calendar.CurrentMonth = time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	nextMonth = calendar.GetNextMonth()
	expected = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if !nextMonth.Equal(expected) {
		t.Errorf("GetNextMonth() across year boundary = %v, want %v", nextMonth, expected)
	}
}

func TestCalendar_NavigateBackward(t *testing.T) {
	calendar := NewCalendar()

	// Set a known month for testing
	initialMonth := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	calendar.CurrentMonth = initialMonth

	calendar.NavigateBackward()
	expected := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	if !calendar.CurrentMonth.Equal(expected) {
		t.Errorf("NavigateBackward() CurrentMonth = %v, want %v", calendar.CurrentMonth, expected)
	}

	// Test multiple backward navigations
	calendar.NavigateBackward()
	expected = time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	if !calendar.CurrentMonth.Equal(expected) {
		t.Errorf("Second NavigateBackward() CurrentMonth = %v, want %v", calendar.CurrentMonth, expected)
	}

	// Test year boundary crossing
	calendar.CurrentMonth = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	calendar.NavigateBackward()
	expected = time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)

	if !calendar.CurrentMonth.Equal(expected) {
		t.Errorf("NavigateBackward() across year boundary = %v, want %v", calendar.CurrentMonth, expected)
	}
}

func TestCalendar_NavigateForward(t *testing.T) {
	calendar := NewCalendar()

	// Set a known month for testing
	initialMonth := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	calendar.CurrentMonth = initialMonth

	calendar.NavigateForward()
	expected := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)

	if !calendar.CurrentMonth.Equal(expected) {
		t.Errorf("NavigateForward() CurrentMonth = %v, want %v", calendar.CurrentMonth, expected)
	}

	// Test multiple forward navigations
	calendar.NavigateForward()
	expected = time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)

	if !calendar.CurrentMonth.Equal(expected) {
		t.Errorf("Second NavigateForward() CurrentMonth = %v, want %v", calendar.CurrentMonth, expected)
	}

	// Test year boundary crossing
	calendar.CurrentMonth = time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	calendar.NavigateForward()
	expected = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	if !calendar.CurrentMonth.Equal(expected) {
		t.Errorf("NavigateForward() across year boundary = %v, want %v", calendar.CurrentMonth, expected)
	}
}

func TestCalendar_AddEvent(t *testing.T) {
	calendar := NewCalendar()

	// Test adding single event
	event1 := Event{
		Date:        time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Time:        time.Date(0, 1, 1, 10, 30, 0, 0, time.UTC),
		Description: "Test event 1",
	}

	calendar.AddEvent(event1)

	if len(calendar.Events) != 1 {
		t.Errorf("Event count = %d, want 1", len(calendar.Events))
	}

	if calendar.Events[0].Description != "Test event 1" {
		t.Errorf("Event description = %s, want 'Test event 1'", calendar.Events[0].Description)
	}

	// Test adding multiple events
	event2 := Event{
		Date:        time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		Time:        time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC),
		Description: "Test event 2",
	}

	calendar.AddEvent(event2)

	if len(calendar.Events) != 2 {
		t.Errorf("Event count = %d, want 2", len(calendar.Events))
	}

	// Verify both events are present
	descriptions := make(map[string]bool)
	for _, event := range calendar.Events {
		descriptions[event.Description] = true
	}

	if !descriptions["Test event 1"] {
		t.Error("Test event 1 not found")
	}

	if !descriptions["Test event 2"] {
		t.Error("Test event 2 not found")
	}
}

func TestCalendar_HasEventsForDate(t *testing.T) {
	calendar := NewCalendar()

	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	// Test empty calendar
	if calendar.HasEventsForDate(testDate) {
		t.Error("HasEventsForDate() should return false for empty calendar")
	}

	// Add events for testing
	events := []Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Event 1"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC), Description: "Event 2"},
	}

	for _, event := range events {
		calendar.AddEvent(event)
	}

	// Test date with events
	if !calendar.HasEventsForDate(testDate) {
		t.Error("HasEventsForDate() should return true for date with events")
	}

	// Test date without events
	noEventDate := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)
	if calendar.HasEventsForDate(noEventDate) {
		t.Error("HasEventsForDate() should return false for date without events")
	}

	// Test date normalization (time component should be ignored)
	testDateWithTime := time.Date(2025, 8, 15, 14, 30, 45, 0, time.UTC)
	if !calendar.HasEventsForDate(testDateWithTime) {
		t.Error("HasEventsForDate() should normalize dates and ignore time component")
	}
}

func TestCalendar_GetEventsForDate(t *testing.T) {
	calendar := NewCalendar()

	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	// Test empty calendar
	events := calendar.GetEventsForDate(testDate)
	if len(events) != 0 {
		t.Errorf("GetEventsForDate() on empty calendar = %d events, want 0", len(events))
	}

	// Add events for testing (unsorted to test sorting functionality)
	testEvents := []Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 14, 30, 0, 0, time.UTC), Description: "Afternoon event"},
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC), Description: "Morning event"},
		{Date: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC), Description: "Different day"},
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), Time: time.Date(0, 1, 1, 18, 45, 0, 0, time.UTC), Description: "Evening event"},
	}

	for _, event := range testEvents {
		calendar.AddEvent(event)
	}

	// Test getting events for specific date
	events = calendar.GetEventsForDate(testDate)
	if len(events) != 3 {
		t.Errorf("GetEventsForDate() = %d events, want 3", len(events))
	}

	// Verify events are sorted by time (ascending)
	expectedOrder := []string{"Morning event", "Afternoon event", "Evening event"}
	for i, event := range events {
		if event.Description != expectedOrder[i] {
			t.Errorf("Event %d description = %s, want %s", i, event.Description, expectedOrder[i])
		}
	}

	// Test date normalization
	testDateWithTime := time.Date(2025, 8, 15, 23, 59, 59, 0, time.UTC)
	events = calendar.GetEventsForDate(testDateWithTime)
	if len(events) != 3 {
		t.Errorf("GetEventsForDate() with time component should return 3 events, got %d", len(events))
	}
}

func TestCalendar_GetEventsForDate_SortingEdgeCases(t *testing.T) {
	calendar := NewCalendar()

	testDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	// Add events with same time to test stable sorting
	testEvents := []Event{
		{Date: testDate, Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Event A at 10:00"},
		{Date: testDate, Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC), Description: "Event B at 10:00"},
		{Date: testDate, Time: time.Date(0, 1, 1, 9, 30, 0, 0, time.UTC), Description: "Event at 09:30"},
	}

	for _, event := range testEvents {
		calendar.AddEvent(event)
	}

	events := calendar.GetEventsForDate(testDate)
	if len(events) != 3 {
		t.Errorf("GetEventsForDate() = %d events, want 3", len(events))
	}

	// First event should be the earliest time
	if events[0].Description != "Event at 09:30" {
		t.Errorf("First event = %s, want 'Event at 09:30'", events[0].Description)
	}

	// Remaining events should be the 10:00 events (order may vary due to stable sort)
	if events[1].Time.Hour() != 10 || events[2].Time.Hour() != 10 {
		t.Error("Events 1 and 2 should both be at 10:00")
	}
}

func TestCalendar_NavigationEdgeCases(t *testing.T) {
	calendar := NewCalendar()

	// Test navigation with leap year
	calendar.CurrentMonth = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // February 2024 (leap year)

	nextMonth := calendar.GetNextMonth()
	expectedNext := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	if !nextMonth.Equal(expectedNext) {
		t.Errorf("GetNextMonth() from leap year February = %v, want %v", nextMonth, expectedNext)
	}

	prevMonth := calendar.GetPreviousMonth()
	expectedPrev := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if !prevMonth.Equal(expectedPrev) {
		t.Errorf("GetPreviousMonth() from leap year February = %v, want %v", prevMonth, expectedPrev)
	}

	// Test navigation doesn't affect original CurrentMonth until Navigate* methods are called
	originalMonth := calendar.CurrentMonth
	calendar.GetNextMonth()
	calendar.GetPreviousMonth()

	if !calendar.CurrentMonth.Equal(originalMonth) {
		t.Error("Get methods should not modify CurrentMonth")
	}
}

func TestCalendar_TimezoneBehavior(t *testing.T) {
	calendar := NewCalendar()

	// Test with events in the same timezone
	utc := time.UTC

	// Add events with same logical date but different times
	events := []Event{
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, utc), Time: time.Date(0, 1, 1, 14, 0, 0, 0, utc), Description: "UTC event 1"},
		{Date: time.Date(2025, 8, 15, 0, 0, 0, 0, utc), Time: time.Date(0, 1, 1, 11, 0, 0, 0, utc), Description: "UTC event 2"},
	}

	for _, event := range events {
		calendar.AddEvent(event)
	}

	// Both events should be found for the same date
	testDate := time.Date(2025, 8, 15, 12, 0, 0, 0, utc)
	foundEvents := calendar.GetEventsForDate(testDate)

	if len(foundEvents) != 2 {
		t.Errorf("GetEventsForDate() with same date events = %d, want 2", len(foundEvents))
	}

	// Verify HasEventsForDate works
	if !calendar.HasEventsForDate(testDate) {
		t.Error("HasEventsForDate() should find events for the same date")
	}

	// Test that different dates don't match
	differentDate := time.Date(2025, 8, 16, 12, 0, 0, 0, utc)
	foundEventsDifferentDate := calendar.GetEventsForDate(differentDate)

	if len(foundEventsDifferentDate) != 0 {
		t.Errorf("GetEventsForDate() with different date = %d, want 0", len(foundEventsDifferentDate))
	}
}
