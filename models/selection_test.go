package models

import (
	"testing"
	"time"
)

func TestNewSelection(t *testing.T) {
	calendar := NewCalendar()
	selection := NewSelection(calendar)

	if selection == nil {
		t.Fatal("NewSelection() returned nil")
	}

	// Verify Calendar reference is set
	if selection.Calendar != calendar {
		t.Error("Selection should reference the provided calendar")
	}

	// Verify SelectedDate is set to today's date (normalized)
	now := time.Now()
	expectedDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if !selection.SelectedDate.Equal(expectedDate) {
		t.Errorf("SelectedDate = %v, want %v", selection.SelectedDate, expectedDate)
	}

	// Verify time components are zero
	if selection.SelectedDate.Hour() != 0 || selection.SelectedDate.Minute() != 0 || selection.SelectedDate.Second() != 0 {
		t.Error("SelectedDate should have zero time components")
	}
}

func TestSelection_IsWithinVisibleRange(t *testing.T) {
	calendar := NewCalendar()
	// Set calendar to August 2025 for predictable testing
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	selection := NewSelection(calendar)

	tests := []struct {
		name           string
		selectedDate   time.Time
		expectedResult bool
	}{
		{
			name:           "Date in current month",
			selectedDate:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name:           "Date in previous month",
			selectedDate:   time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name:           "Date in next month",
			selectedDate:   time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name:           "First day of previous month",
			selectedDate:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name:           "Last day of next month",
			selectedDate:   time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC),
			expectedResult: true,
		},
		{
			name:           "Date before visible range",
			selectedDate:   time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC),
			expectedResult: false,
		},
		{
			name:           "Date after visible range",
			selectedDate:   time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selection.SelectedDate = tt.selectedDate
			result := selection.IsWithinVisibleRange()
			if result != tt.expectedResult {
				t.Errorf("IsWithinVisibleRange() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestSelection_MoveLeft(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	selection := NewSelection(calendar)

	// Test normal left movement (within bounds)
	selection.SelectedDate = time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	initialDate := selection.SelectedDate

	selection.MoveLeft()
	expectedDate := initialDate.AddDate(0, 0, -1)

	if !selection.SelectedDate.Equal(expectedDate) {
		t.Errorf("MoveLeft() resulted in %v, want %v", selection.SelectedDate, expectedDate)
	}

	// Test boundary constraint - should not move beyond visible range
	selection.SelectedDate = time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC) // First day of prev month
	beforeMove := selection.SelectedDate

	selection.MoveLeft() // This should not change the date as it would go out of bounds

	if !selection.SelectedDate.Equal(beforeMove) {
		t.Errorf("MoveLeft() at boundary should not change date, got %v, want %v", selection.SelectedDate, beforeMove)
	}
}

func TestSelection_MoveRight(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	selection := NewSelection(calendar)

	// Test normal right movement (within bounds)
	selection.SelectedDate = time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	initialDate := selection.SelectedDate

	selection.MoveRight()
	expectedDate := initialDate.AddDate(0, 0, 1)

	if !selection.SelectedDate.Equal(expectedDate) {
		t.Errorf("MoveRight() resulted in %v, want %v", selection.SelectedDate, expectedDate)
	}

	// Test boundary constraint - should not move beyond visible range
	selection.SelectedDate = time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC) // Last day of next month
	beforeMove := selection.SelectedDate

	selection.MoveRight() // This should not change the date as it would go out of bounds

	if !selection.SelectedDate.Equal(beforeMove) {
		t.Errorf("MoveRight() at boundary should not change date, got %v, want %v", selection.SelectedDate, beforeMove)
	}
}

func TestSelection_MoveUp(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	selection := NewSelection(calendar)

	// Test normal up movement (within bounds) - moves 7 days back
	selection.SelectedDate = time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	initialDate := selection.SelectedDate

	selection.MoveUp()
	expectedDate := initialDate.AddDate(0, 0, -7)

	if !selection.SelectedDate.Equal(expectedDate) {
		t.Errorf("MoveUp() resulted in %v, want %v", selection.SelectedDate, expectedDate)
	}

	// Test boundary constraint - should not move beyond visible range
	selection.SelectedDate = time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC) // Early in prev month
	beforeMove := selection.SelectedDate

	selection.MoveUp() // This should not change the date as it would go out of bounds

	if !selection.SelectedDate.Equal(beforeMove) {
		t.Errorf("MoveUp() at boundary should not change date, got %v, want %v", selection.SelectedDate, beforeMove)
	}
}

func TestSelection_MoveDown(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	selection := NewSelection(calendar)

	// Test normal down movement (within bounds) - moves 7 days forward
	selection.SelectedDate = time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	initialDate := selection.SelectedDate

	selection.MoveDown()
	expectedDate := initialDate.AddDate(0, 0, 7)

	if !selection.SelectedDate.Equal(expectedDate) {
		t.Errorf("MoveDown() resulted in %v, want %v", selection.SelectedDate, expectedDate)
	}

	// Test boundary constraint - should not move beyond visible range
	selection.SelectedDate = time.Date(2025, 9, 25, 0, 0, 0, 0, time.UTC) // Late in next month
	beforeMove := selection.SelectedDate

	selection.MoveDown() // This should not change the date as it would go out of bounds

	if !selection.SelectedDate.Equal(beforeMove) {
		t.Errorf("MoveDown() at boundary should not change date, got %v, want %v", selection.SelectedDate, beforeMove)
	}
}

func TestSelection_AdjustForMonthChange(t *testing.T) {
	calendar := NewCalendar()
	selection := NewSelection(calendar)

	tests := []struct {
		name           string
		currentMonth   time.Time
		selectedDate   time.Time
		expectedResult time.Time
		shouldAdjust   bool
	}{
		{
			name:           "Selection within range - no adjustment",
			currentMonth:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			selectedDate:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
			expectedResult: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
			shouldAdjust:   false,
		},
		{
			name:           "Selection out of range - adjust to current month",
			currentMonth:   time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			selectedDate:   time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC), // Way out of range
			expectedResult: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC), // Same day in current month
			shouldAdjust:   true,
		},
		{
			name:           "Day exists in new month",
			currentMonth:   time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			selectedDate:   time.Date(2025, 4, 15, 0, 0, 0, 0, time.UTC), // Out of range
			expectedResult: time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), // Same day preserved
			shouldAdjust:   true,
		},
		{
			name:           "Day doesn't exist in new month - adjust to last day",
			currentMonth:   time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),   // February has 28 days
			selectedDate:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC), // Out of range, day 31
			expectedResult: time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),  // Adjusted to last day of Feb
			shouldAdjust:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calendar.CurrentMonth = tt.currentMonth
			selection.SelectedDate = tt.selectedDate

			selection.AdjustForMonthChange()

			if !selection.SelectedDate.Equal(tt.expectedResult) {
				t.Errorf("AdjustForMonthChange() resulted in %v, want %v", selection.SelectedDate, tt.expectedResult)
			}
		})
	}
}

func TestSelection_BoundaryEdgeCases(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	selection := NewSelection(calendar)

	// Test movement from first day of visible range
	selection.SelectedDate = time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)

	// Should be able to move right
	selection.MoveRight()
	expected := time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("MoveRight() from first day = %v, want %v", selection.SelectedDate, expected)
	}

	// Reset and test movement from last day of visible range
	selection.SelectedDate = time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)

	// Should be able to move left
	selection.MoveLeft()
	expected = time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("MoveLeft() from last day = %v, want %v", selection.SelectedDate, expected)
	}
}

func TestSelection_CrossMonthMovement(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	selection := NewSelection(calendar)

	// Test movement across month boundaries within visible range
	// Move from July to August
	selection.SelectedDate = time.Date(2025, 7, 31, 0, 0, 0, 0, time.UTC)
	selection.MoveRight()
	expected := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)

	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("Cross-month MoveRight() = %v, want %v", selection.SelectedDate, expected)
	}

	// Move from August to September
	selection.SelectedDate = time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC)
	selection.MoveRight()
	expected = time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)

	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("Cross-month MoveRight() Aug->Sep = %v, want %v", selection.SelectedDate, expected)
	}
}

func TestSelection_WeekMovementEdgeCases(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	selection := NewSelection(calendar)

	// Test week movement across month boundaries
	selection.SelectedDate = time.Date(2025, 7, 28, 0, 0, 0, 0, time.UTC) // Monday in July
	selection.MoveDown()                                                  // Should go to Monday in August
	expected := time.Date(2025, 8, 4, 0, 0, 0, 0, time.UTC)

	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("MoveDown() across month boundary = %v, want %v", selection.SelectedDate, expected)
	}

	// Test week movement back
	selection.MoveUp() // Should go back to July
	expected = time.Date(2025, 7, 28, 0, 0, 0, 0, time.UTC)

	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("MoveUp() back across month boundary = %v, want %v", selection.SelectedDate, expected)
	}
}

func TestSelection_LeapYearHandling(t *testing.T) {
	calendar := NewCalendar()
	calendar.CurrentMonth = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC) // February 2024 (leap year)
	selection := NewSelection(calendar)

	// Test adjustment with leap year February
	selection.SelectedDate = time.Date(2023, 12, 29, 0, 0, 0, 0, time.UTC) // Out of range
	selection.AdjustForMonthChange()

	expected := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC) // Should adjust to Feb 29 (leap year)
	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("AdjustForMonthChange() in leap year = %v, want %v", selection.SelectedDate, expected)
	}

	// Test adjustment from Feb 30th (doesn't exist) to last day of February
	calendar.CurrentMonth = time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)   // February 2023 (non-leap year)
	selection.SelectedDate = time.Date(2022, 1, 30, 0, 0, 0, 0, time.UTC) // Out of range, day 30
	selection.AdjustForMonthChange()

	expected = time.Date(2023, 2, 28, 0, 0, 0, 0, time.UTC) // Should adjust to Feb 28 (non-leap year)
	if !selection.SelectedDate.Equal(expected) {
		t.Errorf("AdjustForMonthChange() in non-leap year = %v, want %v", selection.SelectedDate, expected)
	}
}
