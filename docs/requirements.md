# Requirements Document

## Introduction
This document specifies the functional and related non-functional requirements for a terminal-based (CLI/CUI) calendar application for scheduling events. Upon launch, the application displays three calendar months side-by-side (previous, current, next). Users navigate the three-month window with keyboard controls, select specific days, view events for a selected date, and add new events. All event data is persisted to a text file.

The intended environment is a standard terminal. Interaction is entirely keyboard-driven.

## Requirements

### Requirement 1: Three-Month Calendar Display on Startup
- User Story: As a user, I want to see the previous, current, and next months side-by-side at startup so that I immediately understand the current date in context and can plan around it.
- Acceptance Criteria:
  - WHEN the application starts THEN the system SHALL render three calendar months side-by-side showing the previous, current, and next months relative to the system’s local date.
  - WHEN rendering a month THEN the system SHALL show day-of-week headers and align date numbers under the correct weekdays, correctly handling varying month lengths and leap years.
  - WHEN today’s date is within the three displayed months THEN the system SHALL visually distinguish today’s date from other dates.
  - WHEN the terminal size is at least 80x24 characters THEN the system SHALL render the three-month view without truncation.

### Requirement 2: Day Indicators in Calendar Grid
- User Story: As a user, I want day cells to display indicators so that I can quickly identify special states such as selected day, today, or days with events.
- Acceptance Criteria:
  - WHEN a day has one or more events THEN the system SHALL display an event indicator in that day’s cell (e.g., a dot or marker) using ASCII-safe characters.
  - WHEN a day is selected THEN the system SHALL visually highlight the selected day’s cell distinctly from other states.
  - WHEN a day is both today and selected THEN the system SHALL present a combined, non-ambiguous indication (e.g., inverted colors plus a marker), ensuring both states are perceivable.
  - WHEN indicators are rendered THEN the system SHALL avoid using symbols that conflict with date digits and SHALL maintain legibility in monochrome terminals.

### Requirement 3: Month Navigation with B/N Keys
- User Story: As a user, I want to navigate the three-month window backward or forward by month so that I can browse past and future months.
- Acceptance Criteria:
  - WHEN the user presses 'B' or 'b' THEN the system SHALL shift the three-month window backward by one month (e.g., May–Jun–Jul becomes Apr–May–Jun).
  - WHEN the user presses 'N' or 'n' THEN the system SHALL shift the three-month window forward by one month (e.g., May–Jun–Jul becomes Jun–Jul–Aug).
  - WHEN navigation crosses year boundaries THEN the system SHALL correctly adjust the year (e.g., Jan back to Dec of the previous year).
  - WHEN the window shifts THEN the system SHALL preserve the selected day number if it exists in the new context; otherwise, it SHALL select the last valid day of the corresponding month.

### Requirement 4: Day Navigation with H/J/K/L Keys
- User Story: As a user, I want to move the selection among days using H, J, K, L so that I can quickly choose a specific date without leaving the keyboard.
- Acceptance Criteria:
  - WHEN the user presses 'H' or 'h' THEN the system SHALL move the selection one day to the left; if it is the first day of a month, it SHALL move into the previous visible month if applicable.
  - WHEN the user presses 'L' or 'l' THEN the system SHALL move the selection one day to the right; if it is the last day of a month, it SHALL move into the next visible month if applicable.
  - WHEN the user presses 'K' or 'k' THEN the system SHALL move the selection one week up (−7 days), staying within the visible three-month range when possible.
  - WHEN the user presses 'J' or 'j' THEN the system SHALL move the selection one week down (+7 days), staying within the visible three-month range when possible.
  - WHEN a movement would go beyond the currently visible three-month window THEN the system SHALL keep the selection at the nearest boundary day; the user can use 'B'/'N' to shift the window further.

### Requirement 5: View Events for Selected Date (Enter Key)
- User Story: As a user, I want to press Enter on a selected date to view its events so that I can see what is scheduled that day.
- Acceptance Criteria:
  - WHEN the user presses Enter on a selected date THEN the system SHALL display the list of events for that date within the UI (e.g., in a pane or panel) while retaining the three-month calendar context on screen.
  - WHEN there are no events for the selected date THEN the system SHALL display a message indicating that no events are scheduled.
  - WHEN events are listed THEN the system SHALL present them sorted by start time ascending.
  - WHEN the events list is displayed THEN the selection for the date SHALL remain unchanged in the calendar.

### Requirement 6: Add New Event (A Key)
- User Story: As a user, I want to add a new event to the selected date by pressing 'A' so that I can schedule items directly from the calendar view.
- Acceptance Criteria:
  - WHEN the user presses 'A' or 'a' while the events list for a selected date is visible THEN the system SHALL prompt for the event time and description in sequence.
  - WHEN prompted for time THEN the system SHALL accept input in HH:MM 24-hour format, using the selected date as the event date.
  - WHEN the time input is invalid (e.g., not HH:MM, out of range) THEN the system SHALL inform the user of the error and re-prompt for a valid time without losing context.
  - WHEN prompted for description THEN the system SHALL accept a non-empty text line as the event description; an empty description SHALL be rejected with an error prompt.
  - WHEN a valid time and description are provided THEN the system SHALL persist the event to the storage file and immediately refresh the UI to show the new event and update the day indicator.

### Requirement 7: Data Persistence to Text File
- User Story: As a user, I want my events to be saved to a text file so that they persist across application restarts.
- Acceptance Criteria:
  - WHEN the application starts THEN the system SHALL load existing events from a text file named "events.txt" located in the application’s working directory; if the file does not exist, it SHALL start with an empty dataset.
  - WHEN saving an event THEN the system SHALL append or write the event in a line-oriented format: YYYY-MM-DD|HH:MM|description, preserving spaces in the description.
  - WHEN multiple events exist for the same date THEN the system SHALL load and display them all, sorted by time when listing.
  - WHEN the storage file contains malformed lines THEN the system SHALL skip those lines, continue processing, and log or show a minimal warning message without terminating.
  - WHEN file I/O errors occur (e.g., permission denied) THEN the system SHALL report the error to the user and SHALL not corrupt the existing file contents.

### Requirement 8: Keyboard Input Handling and Hints
- User Story: As a user, I want consistent keyboard controls and basic on-screen hints so that I can quickly learn and use the application.
- Acceptance Criteria:
  - WHEN the user presses any of the supported keys THEN the system SHALL accept both uppercase and lowercase variants for B, N, H, J, K, L, A.
  - WHEN the user presses an unrecognized key THEN the system SHALL ignore it and MAY present a brief status message indicating the key is not bound.
  - WHEN the application is in the calendar view THEN the system SHALL display a concise key legend (e.g., B/N: month, H/J/K/L: move, Enter: events, A: add) within the available space.

### Requirement 9: Visual and Terminal Compatibility
- User Story: As a user, I want the calendar to be readable in a typical terminal so that I can use it on various systems.
- Acceptance Criteria:
  - WHEN rendered in a standard 80x24, monochrome terminal THEN the system SHALL keep the calendar and indicators legible using ASCII-friendly characters and spacing.
  - WHEN color support is available THEN the system MAY use color to enhance readability; however, all information SHALL remain understandable without color.

### Requirement 10: Date, Time, and Localization Assumptions
- User Story: As a user, I want the calendar to respect my system’s date and time so that displayed information aligns with my expectations.
- Acceptance Criteria:
  - WHEN determining the current date and time THEN the system SHALL use the local system clock and time zone.
  - WHEN displaying day-of-week headers THEN the system SHALL default to Sunday-first ordering unless locale information indicates otherwise.
  - WHEN calculating month boundaries (including leap years) THEN the system SHALL use correct Gregorian calendar rules.
