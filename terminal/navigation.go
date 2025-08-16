package terminal

import (
	"time"

	"go-ascii-calendar/models"
)

// NavigationController handles navigation logic for the calendar
type NavigationController struct {
	calendar  *models.Calendar
	selection *models.Selection
}

// NewNavigationController creates a new navigation controller
func NewNavigationController(calendar *models.Calendar, selection *models.Selection) *NavigationController {
	return &NavigationController{
		calendar:  calendar,
		selection: selection,
	}
}

// NavigateMonthBackward shifts the three-month window backward by one month (B key)
func (nc *NavigationController) NavigateMonthBackward() {
	// Store the current selected day number for preservation
	selectedDay := nc.selection.SelectedDate.Day()

	// Shift the calendar window backward
	nc.calendar.NavigateBackward()

	// Adjust selection to preserve the day number if possible
	nc.adjustSelectionForMonthChange(selectedDay)
}

// NavigateMonthForward shifts the three-month window forward by one month (N key)
func (nc *NavigationController) NavigateMonthForward() {
	// Store the current selected day number for preservation
	selectedDay := nc.selection.SelectedDate.Day()

	// Shift the calendar window forward
	nc.calendar.NavigateForward()

	// Adjust selection to preserve the day number if possible
	nc.adjustSelectionForMonthChange(selectedDay)
}

// NavigateDayLeft moves selection one day to the left (H key)
func (nc *NavigationController) NavigateDayLeft() {
	newDate := nc.selection.SelectedDate.AddDate(0, 0, -1)

	// Check if the new date is within the visible three-month range
	if nc.isDateInVisibleRange(newDate) {
		nc.selection.SelectedDate = newDate
	} else {
		// Move to the previous month if we're at the beginning of a month
		if nc.selection.SelectedDate.Day() == 1 {
			// Try to move to the last day of the previous month if it's visible
			prevMonth := nc.selection.SelectedDate.AddDate(0, -1, 0)
			lastDayOfPrevMonth := nc.getLastDayOfMonth(prevMonth)

			if nc.isDateInVisibleRange(lastDayOfPrevMonth) {
				nc.selection.SelectedDate = lastDayOfPrevMonth
			}
		}
	}
}

// NavigateDayRight moves selection one day to the right (L key)
func (nc *NavigationController) NavigateDayRight() {
	newDate := nc.selection.SelectedDate.AddDate(0, 0, 1)

	// Check if the new date is within the visible three-month range
	if nc.isDateInVisibleRange(newDate) {
		nc.selection.SelectedDate = newDate
	} else {
		// Move to the next month if we're at the end of a month
		daysInCurrentMonth := nc.getDaysInMonth(nc.selection.SelectedDate)
		if nc.selection.SelectedDate.Day() == daysInCurrentMonth {
			// Try to move to the first day of the next month if it's visible
			nextMonth := nc.selection.SelectedDate.AddDate(0, 1, 0)
			firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

			if nc.isDateInVisibleRange(firstDayOfNextMonth) {
				nc.selection.SelectedDate = firstDayOfNextMonth
			}
		}
	}
}

// NavigateDayUp moves selection one week up (K key)
func (nc *NavigationController) NavigateDayUp() {
	newDate := nc.selection.SelectedDate.AddDate(0, 0, -7)

	// Check if the new date is within the visible three-month range
	if nc.isDateInVisibleRange(newDate) {
		nc.selection.SelectedDate = newDate
	}
	// If not in range, keep the current selection (boundary constraint)
}

// NavigateDayDown moves selection one week down (J key)
func (nc *NavigationController) NavigateDayDown() {
	newDate := nc.selection.SelectedDate.AddDate(0, 0, 7)

	// Check if the new date is within the visible three-month range
	if nc.isDateInVisibleRange(newDate) {
		nc.selection.SelectedDate = newDate
	}
	// If not in range, keep the current selection (boundary constraint)
}

// adjustSelectionForMonthChange adjusts selection when the month window changes
// Preserves the selected day if it exists in the new context, otherwise selects last valid day
func (nc *NavigationController) adjustSelectionForMonthChange(desiredDay int) {
	// If the selection is no longer in the visible range, adjust it
	if !nc.selection.IsWithinVisibleRange() {
		// Try to preserve the same day number in the current month
		currentMonth := nc.calendar.CurrentMonth

		// Get the last day of the current month
		daysInMonth := nc.getDaysInMonth(currentMonth)

		// Use the desired day or the last valid day of the month
		actualDay := desiredDay
		if desiredDay > daysInMonth {
			actualDay = daysInMonth
		}

		nc.selection.SelectedDate = time.Date(currentMonth.Year(), currentMonth.Month(), actualDay, 0, 0, 0, 0, currentMonth.Location())
	}
}

// isDateInVisibleRange checks if a date is within the visible three-month range
func (nc *NavigationController) isDateInVisibleRange(date time.Time) bool {
	prevMonth := nc.calendar.GetPreviousMonth()
	nextMonth := nc.calendar.GetNextMonth()

	// Calculate the start and end of the visible range
	startRange := time.Date(prevMonth.Year(), prevMonth.Month(), 1, 0, 0, 0, 0, prevMonth.Location())
	endRange := nc.getLastDayOfMonth(nextMonth)

	return !date.Before(startRange) && !date.After(endRange)
}

// getDaysInMonth returns the number of days in the given month
func (nc *NavigationController) getDaysInMonth(date time.Time) int {
	// Get the first day of the next month, then subtract one day
	firstOfNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastOfThisMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfThisMonth.Day()
}

// getLastDayOfMonth returns the last day of the given month
func (nc *NavigationController) getLastDayOfMonth(date time.Time) time.Time {
	firstDayNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	return firstDayNextMonth.AddDate(0, 0, -1)
}

// GetCurrentSelection returns the currently selected date
func (nc *NavigationController) GetCurrentSelection() time.Time {
	return nc.selection.SelectedDate
}

// SetSelection sets the selected date (with validation)
func (nc *NavigationController) SetSelection(date time.Time) bool {
	if nc.isDateInVisibleRange(date) {
		nc.selection.SelectedDate = date
		return true
	}
	return false
}

// IsSelectionInCurrentMonth checks if the selection is in the middle (current) month
func (nc *NavigationController) IsSelectionInCurrentMonth() bool {
	selectedMonth := nc.selection.SelectedDate.Month()
	selectedYear := nc.selection.SelectedDate.Year()
	currentMonth := nc.calendar.CurrentMonth.Month()
	currentYear := nc.calendar.CurrentMonth.Year()

	return selectedMonth == currentMonth && selectedYear == currentYear
}

// GetVisibleDateRange returns the start and end dates of the visible range
func (nc *NavigationController) GetVisibleDateRange() (start, end time.Time) {
	prevMonth := nc.calendar.GetPreviousMonth()
	nextMonth := nc.calendar.GetNextMonth()

	start = time.Date(prevMonth.Year(), prevMonth.Month(), 1, 0, 0, 0, 0, prevMonth.Location())
	end = nc.getLastDayOfMonth(nextMonth)

	return start, end
}
