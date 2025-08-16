package models

import (
	"time"
)

// Calendar manages the three-month view state (previous, current, next months)
type Calendar struct {
	CurrentMonth time.Time // The middle month of the three-month view
	Events       []Event   // All events loaded from storage
}

// NewCalendar creates a new calendar with the current month as the middle month
func NewCalendar() *Calendar {
	now := time.Now()
	// Set to first day of current month for consistent calculations
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	return &Calendar{
		CurrentMonth: currentMonth,
		Events:       make([]Event, 0),
	}
}

// GetPreviousMonth returns the month before the current month
func (c *Calendar) GetPreviousMonth() time.Time {
	return c.CurrentMonth.AddDate(0, -1, 0)
}

// GetNextMonth returns the month after the current month
func (c *Calendar) GetNextMonth() time.Time {
	return c.CurrentMonth.AddDate(0, 1, 0)
}

// NavigateBackward shifts the three-month window backward by one month
func (c *Calendar) NavigateBackward() {
	c.CurrentMonth = c.CurrentMonth.AddDate(0, -1, 0)
}

// NavigateForward shifts the three-month window forward by one month
func (c *Calendar) NavigateForward() {
	c.CurrentMonth = c.CurrentMonth.AddDate(0, 1, 0)
}

// GetEventsForDate returns all events for a specific date, sorted by time
func (c *Calendar) GetEventsForDate(date time.Time) []Event {
	var events []Event
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for _, event := range c.Events {
		eventDate := time.Date(event.Date.Year(), event.Date.Month(), event.Date.Day(), 0, 0, 0, 0, event.Date.Location())
		if eventDate.Equal(targetDate) {
			events = append(events, event)
		}
	}

	// Sort events by time (bubble sort for simplicity)
	for i := 0; i < len(events)-1; i++ {
		for j := 0; j < len(events)-i-1; j++ {
			if events[j].Time.After(events[j+1].Time) {
				events[j], events[j+1] = events[j+1], events[j]
			}
		}
	}

	return events
}

// HasEventsForDate checks if there are any events for a specific date
func (c *Calendar) HasEventsForDate(date time.Time) bool {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for _, event := range c.Events {
		eventDate := time.Date(event.Date.Year(), event.Date.Month(), event.Date.Day(), 0, 0, 0, 0, event.Date.Location())
		if eventDate.Equal(targetDate) {
			return true
		}
	}
	return false
}

// AddEvent adds a new event to the calendar
func (c *Calendar) AddEvent(event Event) {
	c.Events = append(c.Events, event)
}
