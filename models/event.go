package models

import (
	"time"
)

// Event represents a calendar event with date, time, and description
type Event struct {
	Date        time.Time // The date of the event (YYYY-MM-DD)
	Time        time.Time // The time of the event (HH:MM) - date part will be ignored
	Description string    // The event description
}

// GetTimeString returns the time in HH:MM format
func (e *Event) GetTimeString() string {
	return e.Time.Format("15:04")
}

// GetDateString returns the date in YYYY-MM-DD format
func (e *Event) GetDateString() string {
	return e.Date.Format("2006-01-02")
}

// String returns the event in the storage format: YYYY-MM-DD|HH:MM|description
func (e *Event) String() string {
	return e.GetDateString() + "|" + e.GetTimeString() + "|" + e.Description
}
