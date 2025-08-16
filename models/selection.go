package models

import (
	"time"
)

// Selection tracks the current selected date and position within the calendar
type Selection struct {
	SelectedDate time.Time // The currently selected date
	Calendar     *Calendar // Reference to the calendar for boundary checking
}

// NewSelection creates a new selection with today's date as the initial selection
func NewSelection(calendar *Calendar) *Selection {
	now := time.Now()
	selectedDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return &Selection{
		SelectedDate: selectedDate,
		Calendar:     calendar,
	}
}

// IsWithinVisibleRange checks if the selected date is within the three-month window
func (s *Selection) IsWithinVisibleRange() bool {
	prevMonth := s.Calendar.GetPreviousMonth()
	nextMonth := s.Calendar.GetNextMonth()

	// Check if selected date is within the range from first day of prev month to last day of next month
	startRange := time.Date(prevMonth.Year(), prevMonth.Month(), 1, 0, 0, 0, 0, prevMonth.Location())
	endRange := getLastDayOfMonth(nextMonth)

	return !s.SelectedDate.Before(startRange) && !s.SelectedDate.After(endRange)
}

// MoveLeft moves selection one day to the left (H key)
func (s *Selection) MoveLeft() {
	newDate := s.SelectedDate.AddDate(0, 0, -1)
	if s.isDateWithinBounds(newDate) {
		s.SelectedDate = newDate
	}
}

// MoveRight moves selection one day to the right (L key)
func (s *Selection) MoveRight() {
	newDate := s.SelectedDate.AddDate(0, 0, 1)
	if s.isDateWithinBounds(newDate) {
		s.SelectedDate = newDate
	}
}

// MoveUp moves selection one week up (K key)
func (s *Selection) MoveUp() {
	newDate := s.SelectedDate.AddDate(0, 0, -7)
	if s.isDateWithinBounds(newDate) {
		s.SelectedDate = newDate
	}
}

// MoveDown moves selection one week down (J key)
func (s *Selection) MoveDown() {
	newDate := s.SelectedDate.AddDate(0, 0, 7)
	if s.isDateWithinBounds(newDate) {
		s.SelectedDate = newDate
	}
}

// AdjustForMonthChange adjusts selection when the month window changes
// Preserves the selected day if it exists in the new context, otherwise selects last valid day
func (s *Selection) AdjustForMonthChange() {
	if !s.IsWithinVisibleRange() {
		// Try to preserve the same day number in the current month
		currentMonth := s.Calendar.CurrentMonth
		desiredDay := s.SelectedDate.Day()

		// Get the last day of the current month
		lastDayOfMonth := getLastDayOfMonth(currentMonth).Day()

		// Use the desired day or the last valid day of the month
		if desiredDay > lastDayOfMonth {
			desiredDay = lastDayOfMonth
		}

		s.SelectedDate = time.Date(currentMonth.Year(), currentMonth.Month(), desiredDay, 0, 0, 0, 0, currentMonth.Location())
	}
}

// isDateWithinBounds checks if a date is within the visible three-month range
func (s *Selection) isDateWithinBounds(date time.Time) bool {
	prevMonth := s.Calendar.GetPreviousMonth()
	nextMonth := s.Calendar.GetNextMonth()

	startRange := time.Date(prevMonth.Year(), prevMonth.Month(), 1, 0, 0, 0, 0, prevMonth.Location())
	endRange := getLastDayOfMonth(nextMonth)

	return !date.Before(startRange) && !date.After(endRange)
}

// getLastDayOfMonth returns the last day of the given month
func getLastDayOfMonth(month time.Time) time.Time {
	firstDayNextMonth := time.Date(month.Year(), month.Month()+1, 1, 0, 0, 0, 0, month.Location())
	return firstDayNextMonth.AddDate(0, 0, -1)
}
