package events

import (
	"fmt"
	"sort"
	"time"

	"go-ascii-calendar/calendar"
	"go-ascii-calendar/models"
	"go-ascii-calendar/storage"
)

// Manager handles event operations and integrates with storage
type Manager struct {
	events []models.Event
}

// NewManager creates a new event manager
func NewManager() *Manager {
	return &Manager{
		events: make([]models.Event, 0),
	}
}

// LoadEvents loads all events from storage on application startup
func (m *Manager) LoadEvents() error {
	events, err := storage.LoadEvents()
	if err != nil {
		return fmt.Errorf("failed to load events: %v", err)
	}

	m.events = events
	return nil
}

// GetAllEvents returns all events loaded in memory
func (m *Manager) GetAllEvents() []models.Event {
	return m.events
}

// GetEventsForDate returns all events for a specific date, sorted by time ascending
func (m *Manager) GetEventsForDate(date time.Time) []models.Event {
	var dateEvents []models.Event
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for _, event := range m.events {
		eventDate := time.Date(event.Date.Year(), event.Date.Month(), event.Date.Day(), 0, 0, 0, 0, event.Date.Location())
		if eventDate.Equal(targetDate) {
			dateEvents = append(dateEvents, event)
		}
	}

	// Sort events by time ascending
	sort.Slice(dateEvents, func(i, j int) bool {
		return dateEvents[i].Time.Before(dateEvents[j].Time)
	})

	return dateEvents
}

// HasEventsForDate checks if there are any events for a specific date
func (m *Manager) HasEventsForDate(date time.Time) bool {
	targetDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for _, event := range m.events {
		eventDate := time.Date(event.Date.Year(), event.Date.Month(), event.Date.Day(), 0, 0, 0, 0, event.Date.Location())
		if eventDate.Equal(targetDate) {
			return true
		}
	}
	return false
}

// AddEvent adds a new event with validation and persistence
func (m *Manager) AddEvent(date time.Time, timeStr, description string) error {
	// Validate time string format
	if !calendar.ValidateTimeString(timeStr) {
		return fmt.Errorf("invalid time format '%s': expected HH:MM", timeStr)
	}

	// Validate description is not empty
	if len(description) == 0 {
		return fmt.Errorf("event description cannot be empty")
	}

	// Parse time
	eventTime, err := calendar.ParseTime(timeStr)
	if err != nil {
		return fmt.Errorf("failed to parse time '%s': %v", timeStr, err)
	}

	// Create event
	event := models.Event{
		Date:        date,
		Time:        eventTime,
		Description: description,
	}

	// Validate the complete event
	if err := storage.ValidateEvent(event); err != nil {
		return fmt.Errorf("event validation failed: %v", err)
	}

	// Save to storage
	if err := storage.SaveEvent(event); err != nil {
		return fmt.Errorf("failed to save event: %v", err)
	}

	// Add to in-memory collection
	m.events = append(m.events, event)

	return nil
}

// GetEventCount returns the total number of events
func (m *Manager) GetEventCount() int {
	return len(m.events)
}

// GetEventsForMonth returns all events for a specific month, sorted by date and time
func (m *Manager) GetEventsForMonth(month time.Time) []models.Event {
	var monthEvents []models.Event
	targetYear := month.Year()
	targetMonth := month.Month()

	for _, event := range m.events {
		if event.Date.Year() == targetYear && event.Date.Month() == targetMonth {
			monthEvents = append(monthEvents, event)
		}
	}

	// Sort events by date, then by time
	sort.Slice(monthEvents, func(i, j int) bool {
		if monthEvents[i].Date.Equal(monthEvents[j].Date) {
			return monthEvents[i].Time.Before(monthEvents[j].Time)
		}
		return monthEvents[i].Date.Before(monthEvents[j].Date)
	})

	return monthEvents
}

// GetEventsInDateRange returns all events within a date range, sorted by date and time
func (m *Manager) GetEventsInDateRange(startDate, endDate time.Time) []models.Event {
	var rangeEvents []models.Event

	for _, event := range m.events {
		eventDate := time.Date(event.Date.Year(), event.Date.Month(), event.Date.Day(), 0, 0, 0, 0, event.Date.Location())
		if !eventDate.Before(startDate) && !eventDate.After(endDate) {
			rangeEvents = append(rangeEvents, event)
		}
	}

	// Sort events by date, then by time
	sort.Slice(rangeEvents, func(i, j int) bool {
		if rangeEvents[i].Date.Equal(rangeEvents[j].Date) {
			return rangeEvents[i].Time.Before(rangeEvents[j].Time)
		}
		return rangeEvents[i].Date.Before(rangeEvents[j].Date)
	})

	return rangeEvents
}

// ReloadEvents reloads events from storage (useful for external file changes)
func (m *Manager) ReloadEvents() error {
	return m.LoadEvents()
}
