package events

import (
	"fmt"
	"sort"
	"time"

	"go-ascii-calendar/calendar"
	"go-ascii-calendar/config"
	"go-ascii-calendar/models"
	"go-ascii-calendar/storage"
)

// Manager handles event operations and integrates with storage
type Manager struct {
	events []models.Event
	config *config.Config
}

// NewManager creates a new event manager (legacy function)
func NewManager() *Manager {
	return &Manager{
		events: make([]models.Event, 0),
		config: nil,
	}
}

// NewManagerWithConfig creates a new event manager with configuration
func NewManagerWithConfig(cfg *config.Config) *Manager {
	return &Manager{
		events: make([]models.Event, 0),
		config: cfg,
	}
}

// LoadEvents loads all events from storage on application startup
func (m *Manager) LoadEvents() error {
	var events []models.Event
	var err error

	if m.config != nil {
		// Use configured path with automatic migration
		events, err = storage.LoadEventsWithConfig(m.config.GetEventsFilePath())
	} else {
		// Fallback to legacy text format
		events, err = storage.LoadEvents()
	}

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
	if m.config != nil {
		if err := storage.SaveEventWithConfig(event, m.config.GetEventsFilePath()); err != nil {
			return fmt.Errorf("failed to save event: %v", err)
		}
	} else {
		// Fallback to legacy format
		if err := storage.SaveEvent(event); err != nil {
			return fmt.Errorf("failed to save event: %v", err)
		}
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

// DeleteEvent deletes an event from both storage and memory
func (m *Manager) DeleteEvent(eventToDelete models.Event) error {
	// Delete from storage first
	if m.config != nil {
		if err := storage.DeleteEventWithConfig(eventToDelete, m.config.GetEventsFilePath()); err != nil {
			return fmt.Errorf("failed to delete event from storage: %v", err)
		}
	} else {
		// Fallback to legacy format
		if err := storage.DeleteEvent(eventToDelete); err != nil {
			return fmt.Errorf("failed to delete event from storage: %v", err)
		}
	}

	// Remove from in-memory collection
	var updatedEvents []models.Event
	found := false
	for _, event := range m.events {
		// Compare events by date, time, and description
		if event.Date.Equal(eventToDelete.Date) &&
			event.Time.Equal(eventToDelete.Time) &&
			event.Description == eventToDelete.Description {
			found = true
			continue // Skip this event (delete it)
		}
		updatedEvents = append(updatedEvents, event)
	}

	if !found {
		return fmt.Errorf("event not found in memory for deletion")
	}

	m.events = updatedEvents
	return nil
}

// EditEvent replaces an existing event with a new one in both storage and memory
func (m *Manager) EditEvent(oldEvent models.Event, date time.Time, timeStr, description string) error {
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

	// Create new event
	newEvent := models.Event{
		Date:        date,
		Time:        eventTime,
		Description: description,
	}

	// Validate the complete new event
	if err := storage.ValidateEvent(newEvent); err != nil {
		return fmt.Errorf("new event validation failed: %v", err)
	}

	// Update in storage first
	if m.config != nil {
		if err := storage.UpdateEventWithConfig(oldEvent, newEvent, m.config.GetEventsFilePath()); err != nil {
			return fmt.Errorf("failed to update event in storage: %v", err)
		}
	} else {
		// Fallback to legacy format
		if err := storage.UpdateEvent(oldEvent, newEvent); err != nil {
			return fmt.Errorf("failed to update event in storage: %v", err)
		}
	}

	// Update in-memory collection
	found := false
	for i, event := range m.events {
		// Compare events by date, time, and description
		if event.Date.Equal(oldEvent.Date) &&
			event.Time.Equal(oldEvent.Time) &&
			event.Description == oldEvent.Description {
			m.events[i] = newEvent
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("event not found in memory for update")
	}

	return nil
}
