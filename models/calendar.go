package models

import (
	"sort"
	"time"

	"go-ascii-calendar/calendar"
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
	targetDate := calendar.NormalizeDate(date)

	for _, event := range c.Events {
		eventDate := calendar.NormalizeDate(event.Date)
		if eventDate.Equal(targetDate) {
			events = append(events, event)
		}
	}

	// Sort events by time ascending using Go's built-in sort
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})

	return events
}

// HasEventsForDate checks if there are any events for a specific date
func (c *Calendar) HasEventsForDate(date time.Time) bool {
	targetDate := calendar.NormalizeDate(date)

	for _, event := range c.Events {
		eventDate := calendar.NormalizeDate(event.Date)
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
