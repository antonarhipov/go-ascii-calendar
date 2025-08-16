package terminal

import (
	"testing"
	"time"

	"go-ascii-calendar/models"
)

func TestNavigationController(t *testing.T) {
	// Create test calendar and selection
	cal := models.NewCalendar()
	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	// Test initial state
	currentSelection := nc.GetCurrentSelection()
	if currentSelection.IsZero() {
		t.Error("Initial selection should not be zero time")
	}

	// Test that selection is within visible range initially
	if !nc.selection.IsWithinVisibleRange() {
		t.Error("Initial selection should be within visible range")
	}
}

func TestNavigateMonthBackward(t *testing.T) {
	cal := models.NewCalendar()
	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	// Get initial month
	initialMonth := cal.CurrentMonth
	initialSelectedDay := sel.SelectedDate.Day()

	// Navigate backward
	nc.NavigateMonthBackward()

	// Check that month moved backward
	expectedMonth := initialMonth.AddDate(0, -1, 0)
	if !cal.CurrentMonth.Equal(expectedMonth) {
		t.Errorf("Expected month %v, got %v", expectedMonth, cal.CurrentMonth)
	}

	// Check that selection day is preserved when possible
	if sel.SelectedDate.Day() != initialSelectedDay {
		// This is okay if the previous month doesn't have the same day (e.g., Jan 31 -> Feb 28)
		lastDayOfMonth := nc.getLastDayOfMonth(cal.CurrentMonth).Day()
		if sel.SelectedDate.Day() != lastDayOfMonth {
			t.Errorf("Expected selected day to be preserved (%d) or adjusted to last day (%d), got %d",
				initialSelectedDay, lastDayOfMonth, sel.SelectedDate.Day())
		}
	}

	// Test year boundary crossing
	cal.CurrentMonth = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	nc.NavigateMonthBackward()

	expectedYear := 2024
	expectedMonth2 := time.December
	if cal.CurrentMonth.Year() != expectedYear || cal.CurrentMonth.Month() != expectedMonth2 {
		t.Errorf("Expected December 2024, got %v %d", cal.CurrentMonth.Month(), cal.CurrentMonth.Year())
	}
}

func TestNavigateMonthForward(t *testing.T) {
	cal := models.NewCalendar()
	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	// Get initial month
	initialMonth := cal.CurrentMonth
	initialSelectedDay := sel.SelectedDate.Day()

	// Navigate forward
	nc.NavigateMonthForward()

	// Check that month moved forward
	expectedMonth := initialMonth.AddDate(0, 1, 0)
	if !cal.CurrentMonth.Equal(expectedMonth) {
		t.Errorf("Expected month %v, got %v", expectedMonth, cal.CurrentMonth)
	}

	// Check that selection day is preserved when possible
	if sel.SelectedDate.Day() != initialSelectedDay {
		// This is okay if the next month doesn't have the same day (e.g., Jan 31 -> Feb 28)
		lastDayOfMonth := nc.getLastDayOfMonth(cal.CurrentMonth).Day()
		if sel.SelectedDate.Day() != lastDayOfMonth {
			t.Errorf("Expected selected day to be preserved (%d) or adjusted to last day (%d), got %d",
				initialSelectedDay, lastDayOfMonth, sel.SelectedDate.Day())
		}
	}

	// Test year boundary crossing
	cal.CurrentMonth = time.Date(2025, time.December, 1, 0, 0, 0, 0, time.UTC)
	nc.NavigateMonthForward()

	expectedYear := 2026
	expectedMonth2 := time.January
	if cal.CurrentMonth.Year() != expectedYear || cal.CurrentMonth.Month() != expectedMonth2 {
		t.Errorf("Expected January 2026, got %v %d", cal.CurrentMonth.Month(), cal.CurrentMonth.Year())
	}
}

func TestNavigateDayLeft(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	// Set selection to August 15, 2025 (middle of the month)
	sel := models.NewSelection(cal)
	sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)
	nc := NewNavigationController(cal, sel)

	initialDate := sel.SelectedDate

	// Navigate left (one day backward)
	nc.NavigateDayLeft()

	expectedDate := initialDate.AddDate(0, 0, -1)
	if !sel.SelectedDate.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, sel.SelectedDate)
	}

	// Test month boundary - navigate from August 1st
	sel.SelectedDate = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)
	nc.NavigateDayLeft()

	// Should move to July 31st if July is visible
	if sel.SelectedDate.Month() == time.July && sel.SelectedDate.Day() == 31 {
		// Good - moved to previous month
	} else if sel.SelectedDate.Equal(time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)) {
		// Good - stayed at boundary because previous month not visible
	} else {
		t.Errorf("Unexpected date after left navigation from month boundary: %v", sel.SelectedDate)
	}
}

func TestNavigateDayRight(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	// Set selection to August 15, 2025 (middle of the month)
	sel := models.NewSelection(cal)
	sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)
	nc := NewNavigationController(cal, sel)

	initialDate := sel.SelectedDate

	// Navigate right (one day forward)
	nc.NavigateDayRight()

	expectedDate := initialDate.AddDate(0, 0, 1)
	if !sel.SelectedDate.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, sel.SelectedDate)
	}

	// Test month boundary - navigate from August 31st
	sel.SelectedDate = time.Date(2025, time.August, 31, 0, 0, 0, 0, time.UTC)
	nc.NavigateDayRight()

	// Should move to September 1st if September is visible
	if sel.SelectedDate.Month() == time.September && sel.SelectedDate.Day() == 1 {
		// Good - moved to next month
	} else if sel.SelectedDate.Equal(time.Date(2025, time.August, 31, 0, 0, 0, 0, time.UTC)) {
		// Good - stayed at boundary because next month not visible
	} else {
		t.Errorf("Unexpected date after right navigation from month boundary: %v", sel.SelectedDate)
	}
}

func TestNavigateDayUp(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	// Set selection to August 15, 2025 (middle of the month)
	sel := models.NewSelection(cal)
	sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)
	nc := NewNavigationController(cal, sel)

	initialDate := sel.SelectedDate

	// Navigate up (one week backward)
	nc.NavigateDayUp()

	expectedDate := initialDate.AddDate(0, 0, -7)
	if nc.isDateInVisibleRange(expectedDate) {
		if !sel.SelectedDate.Equal(expectedDate) {
			t.Errorf("Expected date %v, got %v", expectedDate, sel.SelectedDate)
		}
	} else {
		// Should stay at current position if target is out of range
		if !sel.SelectedDate.Equal(initialDate) {
			t.Errorf("Expected to stay at %v when target out of range, got %v", initialDate, sel.SelectedDate)
		}
	}
}

func TestNavigateDayDown(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	// Set selection to August 15, 2025 (middle of the month)
	sel := models.NewSelection(cal)
	sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)
	nc := NewNavigationController(cal, sel)

	initialDate := sel.SelectedDate

	// Navigate down (one week forward)
	nc.NavigateDayDown()

	expectedDate := initialDate.AddDate(0, 0, 7)
	if nc.isDateInVisibleRange(expectedDate) {
		if !sel.SelectedDate.Equal(expectedDate) {
			t.Errorf("Expected date %v, got %v", expectedDate, sel.SelectedDate)
		}
	} else {
		// Should stay at current position if target is out of range
		if !sel.SelectedDate.Equal(initialDate) {
			t.Errorf("Expected to stay at %v when target out of range, got %v", initialDate, sel.SelectedDate)
		}
	}
}

func TestIsDateInVisibleRange(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025 (shows July, August, September)
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"July 1st (in range)", time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC), true},
		{"July 31st (in range)", time.Date(2025, time.July, 31, 0, 0, 0, 0, time.UTC), true},
		{"August 15th (in range)", time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC), true},
		{"September 30th (in range)", time.Date(2025, time.September, 30, 0, 0, 0, 0, time.UTC), true},
		{"June 30th (out of range)", time.Date(2025, time.June, 30, 0, 0, 0, 0, time.UTC), false},
		{"October 1st (out of range)", time.Date(2025, time.October, 1, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nc.isDateInVisibleRange(tt.date)
			if result != tt.expected {
				t.Errorf("isDateInVisibleRange(%v) = %v, want %v", tt.date, result, tt.expected)
			}
		})
	}
}

func TestGetDaysInMonth(t *testing.T) {
	cal := models.NewCalendar()
	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	tests := []struct {
		name     string
		date     time.Time
		expected int
	}{
		{"January 2025", time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC), 31},
		{"February 2025", time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC), 28},
		{"February 2024", time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC), 29},
		{"April 2025", time.Date(2025, time.April, 15, 0, 0, 0, 0, time.UTC), 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nc.getDaysInMonth(tt.date)
			if result != tt.expected {
				t.Errorf("getDaysInMonth(%v) = %d, want %d", tt.date, result, tt.expected)
			}
		})
	}
}

func TestSetSelection(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	// Test setting valid date within visible range
	validDate := time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)
	success := nc.SetSelection(validDate)
	if !success {
		t.Error("SetSelection should succeed for date within visible range")
	}
	if !nc.selection.SelectedDate.Equal(validDate) {
		t.Errorf("Expected selected date %v, got %v", validDate, nc.selection.SelectedDate)
	}

	// Test setting invalid date outside visible range
	invalidDate := time.Date(2025, time.December, 15, 0, 0, 0, 0, time.UTC)
	success = nc.SetSelection(invalidDate)
	if success {
		t.Error("SetSelection should fail for date outside visible range")
	}
	// Selection should remain unchanged
	if !nc.selection.SelectedDate.Equal(validDate) {
		t.Errorf("Selection should remain %v after failed SetSelection, got %v", validDate, nc.selection.SelectedDate)
	}
}

func TestIsSelectionInCurrentMonth(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	// Test selection in current month
	sel.SelectedDate = time.Date(2025, time.August, 15, 0, 0, 0, 0, time.UTC)
	if !nc.IsSelectionInCurrentMonth() {
		t.Error("IsSelectionInCurrentMonth should return true for selection in current month")
	}

	// Test selection in previous month
	sel.SelectedDate = time.Date(2025, time.July, 15, 0, 0, 0, 0, time.UTC)
	if nc.IsSelectionInCurrentMonth() {
		t.Error("IsSelectionInCurrentMonth should return false for selection in previous month")
	}

	// Test selection in next month
	sel.SelectedDate = time.Date(2025, time.September, 15, 0, 0, 0, 0, time.UTC)
	if nc.IsSelectionInCurrentMonth() {
		t.Error("IsSelectionInCurrentMonth should return false for selection in next month")
	}
}

func TestGetVisibleDateRange(t *testing.T) {
	cal := models.NewCalendar()
	// Set calendar to August 2025
	cal.CurrentMonth = time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC)

	sel := models.NewSelection(cal)
	nc := NewNavigationController(cal, sel)

	start, end := nc.GetVisibleDateRange()

	expectedStart := time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2025, time.September, 30, 0, 0, 0, 0, time.UTC)

	if !start.Equal(expectedStart) {
		t.Errorf("Expected start date %v, got %v", expectedStart, start)
	}

	if !end.Equal(expectedEnd) {
		t.Errorf("Expected end date %v, got %v", expectedEnd, end)
	}
}
