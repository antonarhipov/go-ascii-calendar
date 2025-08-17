package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-ascii-calendar/calendar"
	"go-ascii-calendar/models"
)

const EventsFileName = "events.txt"

// LoadEvents loads all events from the events.txt file
func LoadEvents() ([]models.Event, error) {
	return LoadEventsFromFile(EventsFileName)
}

// LoadEventsFromFile loads events from a specified file (for testing)
func LoadEventsFromFile(filename string) ([]models.Event, error) {
	var events []models.Event

	file, err := os.Open(filename)
	if err != nil {
		// If file doesn't exist, return empty slice (not an error)
		if os.IsNotExist(err) {
			return events, nil
		}
		return nil, fmt.Errorf("failed to open events file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		event, err := ParseEventLine(line)
		if err != nil {
			// Log warning but continue processing other lines
			fmt.Printf("Warning: Skipping malformed line %d: %s (error: %v)\n", lineNum, line, err)
			continue
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading events file: %v", err)
	}

	return events, nil
}

// SaveEvent appends a new event to the events.txt file
func SaveEvent(event models.Event) error {
	return SaveEventToFile(event, EventsFileName)
}

// SaveEventToFile appends a new event to a specified file (for testing)
func SaveEventToFile(event models.Event, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open events file for writing: %v", err)
	}
	defer file.Close()

	// Write event in the format: YYYY-MM-DD|HH:MM|description
	eventLine := event.String()
	_, err = file.WriteString(eventLine + "\n")
	if err != nil {
		return fmt.Errorf("failed to write event to file: %v", err)
	}

	return nil
}

// ParseEventLine parses a single line from the events file
// Expected format: YYYY-MM-DD|HH:MM|description
func ParseEventLine(line string) (models.Event, error) {
	parts := strings.SplitN(line, "|", 3)
	if len(parts) != 3 {
		return models.Event{}, fmt.Errorf("invalid format: expected YYYY-MM-DD|HH:MM|description")
	}

	dateStr := strings.TrimSpace(parts[0])
	timeStr := strings.TrimSpace(parts[1])
	description := strings.TrimSpace(parts[2])

	// Validate that description is not empty
	if description == "" {
		return models.Event{}, fmt.Errorf("description cannot be empty")
	}

	// Parse date
	eventDate, err := calendar.ParseDate(dateStr)
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid date format '%s': %v", dateStr, err)
	}

	// Validate and parse time
	if !calendar.ValidateTimeString(timeStr) {
		return models.Event{}, fmt.Errorf("invalid time format '%s': expected HH:MM", timeStr)
	}

	eventTime, err := calendar.ParseTime(timeStr)
	if err != nil {
		return models.Event{}, fmt.Errorf("failed to parse time '%s': %v", timeStr, err)
	}

	return models.Event{
		Date:        eventDate,
		Time:        eventTime,
		Description: description,
	}, nil
}

// ValidateEvent validates an event before saving
func ValidateEvent(event models.Event) error {
	// Check that description is not empty
	if strings.TrimSpace(event.Description) == "" {
		return fmt.Errorf("event description cannot be empty")
	}

	// Validate time format by checking if it can be formatted properly
	timeStr := event.GetTimeString()
	if !calendar.ValidateTimeString(timeStr) {
		return fmt.Errorf("invalid time format: %s", timeStr)
	}

	return nil
}

// FileExists checks if the events file exists
func FileExists() bool {
	return FileExistsAtPath(EventsFileName)
}

// FileExistsAtPath checks if a file exists at the specified path (for testing)
func FileExistsAtPath(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// CreateEventFile creates an empty events file if it doesn't exist
func CreateEventFile() error {
	return CreateEventFileAtPath(EventsFileName)
}

// CreateEventFileAtPath creates an empty events file at the specified path (for testing)
func CreateEventFileAtPath(filename string) error {
	if FileExistsAtPath(filename) {
		return nil // File already exists
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create events file: %v", err)
	}
	defer file.Close()

	return nil
}

// DeleteEvent removes an event from the storage file
func DeleteEvent(eventToDelete models.Event) error {
	return DeleteEventFromFile(eventToDelete, EventsFileName)
}

// DeleteEventFromFile removes an event from a specified file (for testing)
func DeleteEventFromFile(eventToDelete models.Event, filename string) error {
	// Load all events
	events, err := LoadEventsFromFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load events for deletion: %v", err)
	}

	// Find and remove the matching event
	var updatedEvents []models.Event
	found := false
	for _, event := range events {
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
		return fmt.Errorf("event not found for deletion")
	}

	// Rewrite the entire file with the updated events
	return SaveAllEventsToFile(updatedEvents, filename)
}

// UpdateEvent replaces an existing event with a new one
func UpdateEvent(oldEvent, newEvent models.Event) error {
	return UpdateEventInFile(oldEvent, newEvent, EventsFileName)
}

// UpdateEventInFile replaces an existing event with a new one in a specified file (for testing)
func UpdateEventInFile(oldEvent, newEvent models.Event, filename string) error {
	// Load all events
	events, err := LoadEventsFromFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load events for update: %v", err)
	}

	// Find and replace the matching event
	found := false
	for i, event := range events {
		// Compare events by date, time, and description
		if event.Date.Equal(oldEvent.Date) &&
			event.Time.Equal(oldEvent.Time) &&
			event.Description == oldEvent.Description {
			events[i] = newEvent
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("event not found for update")
	}

	// Validate the new event
	if err := ValidateEvent(newEvent); err != nil {
		return fmt.Errorf("new event validation failed: %v", err)
	}

	// Rewrite the entire file with the updated events
	return SaveAllEventsToFile(events, filename)
}

// SaveAllEventsToFile writes all events to a file, replacing the existing content
func SaveAllEventsToFile(events []models.Event, filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open events file for writing: %v", err)
	}
	defer file.Close()

	for _, event := range events {
		eventLine := event.String()
		_, err = file.WriteString(eventLine + "\n")
		if err != nil {
			return fmt.Errorf("failed to write event to file: %v", err)
		}
	}

	return nil
}
