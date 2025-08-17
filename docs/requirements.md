# Enhanced Requirements Document

## Introduction
This document specifies the functional and non-functional requirements for an advanced terminal-based (CLI/CUI) calendar application for comprehensive event management. Upon launch, the application displays three calendar months side-by-side (previous, current, next) with enhanced navigation, search capabilities, and sophisticated event management features. Users can navigate using both vim-style keys and arrow keys, perform inline event operations, search through events, and customize application behavior through configuration files.

The application supports JSON-based data persistence with automatic migration from legacy formats, configuration management through files and command-line options, real-time input validation, and enhanced visual feedback including color support with monochrome terminal compatibility.

The intended environment is a standard terminal with keyboard-driven interaction, supporting both basic monochrome terminals and modern color-enabled terminals.

## Enhanced Requirements

### Requirement 1: Three-Month Calendar Display with Enhanced Visual Indicators
- **User Story**: As a user, I want to see the previous, current, and next months side-by-side at startup with clear visual indicators for today, selected dates, and event-containing days so that I can immediately understand the temporal context and event distribution.
- **Acceptance Criteria**:
    - WHEN the application starts THEN the system SHALL render three calendar months side-by-side showing the previous, current, and next months relative to the system's local date.
    - WHEN rendering a month THEN the system SHALL show day-of-week headers with color coding (cyan in color terminals) and align date numbers under the correct weekdays.
    - WHEN today's date is within the three displayed months THEN the system SHALL visually distinguish it with bright yellow text in color terminals or bold formatting in monochrome terminals.
    - WHEN a day is selected THEN the system SHALL highlight it with blue background and white text in color terminals or reverse video in monochrome terminals.
    - WHEN a day contains events THEN the system SHALL indicate this with green text color in color terminals while maintaining readability in monochrome terminals.
    - WHEN a day is both today and selected THEN the system SHALL combine visual indicators with bright cyan background and white bold text.
    - WHEN the terminal size is at least 80x24 characters THEN the system SHALL render the three-month view without truncation.
    - WHEN events exist for the selected date THEN the system SHALL display up to 10 events in a dedicated section below the calendar, aligned with the calendar's left edge.

### Requirement 2: Dual Navigation System with Enhanced Movement
- **User Story**: As a user, I want flexible navigation options using both vim-style keys (H/J/K/L) and arrow keys, plus quick navigation shortcuts, so that I can efficiently move through dates regardless of my keyboard preference.
- **Acceptance Criteria**:
    - WHEN the user presses 'H', 'h', or Left Arrow THEN the system SHALL move the selection one day to the left, crossing month boundaries within the visible range.
    - WHEN the user presses 'L', 'l', or Right Arrow THEN the system SHALL move the selection one day to the right, crossing month boundaries within the visible range.
    - WHEN the user presses 'K', 'k', or Up Arrow THEN the system SHALL move the selection one week up (−7 days), staying within the visible three-month range.
    - WHEN the user presses 'J', 'j', or Down Arrow THEN the system SHALL move the selection one week down (+7 days), staying within the visible three-month range.
    - WHEN the user presses 'B' or 'b' THEN the system SHALL shift the three-month window backward by one month.
    - WHEN the user presses 'N' or 'n' THEN the system SHALL shift the three-month window forward by one month.
    - WHEN the user presses 'C' or 'c' THEN the system SHALL reset the calendar to the current month and select today's date.
    - WHEN navigation crosses year boundaries THEN the system SHALL correctly handle year transitions.

### Requirement 3: Advanced Event Management with Inline Operations
- **User Story**: As a user, I want to add, edit, and delete events directly from the calendar view with visual feedback and confirmation prompts so that I can efficiently manage my schedule without navigating through multiple screens.
- **Acceptance Criteria**:
    - WHEN the user presses 'A' or 'a' in calendar view THEN the system SHALL immediately start the inline event addition process with time input validation.
    - WHEN adding an event THEN the system SHALL highlight a "[New Event]" row below existing events and accept input directly on that line.
    - WHEN the user presses 'D' or 'd' in calendar view THEN the system SHALL enter event selection mode, highlighting events below the calendar with arrow key navigation.
    - WHEN the user presses 'E' or 'e' in calendar view THEN the system SHALL enter event edit mode, allowing selection and inline editing of existing events.
    - WHEN in event selection mode THEN the system SHALL highlight the selected event with yellow background and show navigation instructions.
    - WHEN confirming event deletion THEN the system SHALL use Enter for confirmation and Esc for cancellation.
    - WHEN editing events THEN the system SHALL pre-fill current values and allow inline modification with real-time validation.
    - WHEN operations complete THEN the system SHALL show non-blocking success messages and immediately update the calendar display.

### Requirement 4: Real-Time Input Validation and Enhanced User Experience
- **User Story**: As a user, I want real-time validation during time input and immediate feedback for all operations so that I can enter data efficiently without trial-and-error correction cycles.
- **Acceptance Criteria**:
    - WHEN entering time THEN the system SHALL validate input character-by-character, only allowing valid digits for each position.
    - WHEN entering the first hour digit THEN the system SHALL only accept '1' or '2' (no hours starting with 0, 3, etc.).
    - WHEN entering the second hour digit THEN the system SHALL validate based on the first digit (10-19 for '1', 20-23 for '2').
    - WHEN entering two hour digits THEN the system SHALL automatically insert a colon and show minute input format.
    - WHEN entering the first minute digit THEN the system SHALL only accept 0-5 (no minutes starting with 6-9).
    - WHEN entering the second minute digit THEN the system SHALL accept any digit 0-9 to complete the time.
    - WHEN time input is incomplete THEN the system SHALL show formatting hints with underscores (e.g., "14:3_").
    - WHEN input validation fails THEN the system SHALL prevent invalid characters without error messages, maintaining smooth user experience.

### Requirement 5: Event Search and Discovery System
- **User Story**: As a user, I want to search through my events by description and navigate to specific dates so that I can quickly find and access relevant events across my entire schedule.
- **Acceptance Criteria**:
    - WHEN the user presses 'F' or 'f' THEN the system SHALL prompt for a search query with live input feedback.
    - WHEN a search query is entered THEN the system SHALL perform case-insensitive matching against event descriptions.
    - WHEN search results exist THEN the system SHALL display them grouped by date under the calendar, sorted chronologically.
    - WHEN displaying search results THEN the system SHALL highlight date headers and format each result with time and description.
    - WHEN navigating search results THEN the system SHALL support up/down arrow key movement with visual selection highlighting.
    - WHEN the user presses Enter on a search result THEN the system SHALL navigate the calendar to that event's date and close the search.
    - WHEN the user presses Esc in search mode THEN the system SHALL cancel the search and return to normal calendar view.
    - WHEN no search results are found THEN the system SHALL display a clear "No events found" message.

### Requirement 6: Configuration Management and Persistence System
- **User Story**: As a user, I want flexible storage options and configuration management so that I can customize file locations and integrate the application into my preferred workflow and file organization system.
- **Acceptance Criteria**:
    - WHEN the application starts THEN the system SHALL create a `~/.ascii-calendar/` directory if it doesn't exist.
    - WHEN no custom paths are specified THEN the system SHALL use `~/.ascii-calendar/events.json` for event storage.
    - WHEN the user provides `-f <path>` THEN the system SHALL use the specified path for event storage.
    - WHEN the user provides `-c <path>` THEN the system SHALL load configuration from the specified JSON file.
    - WHEN a configuration file exists THEN the system SHALL load settings from `~/.ascii-calendar/configuration.json`.
    - WHEN an old `events.txt` file exists THEN the system SHALL automatically migrate it to JSON format and display confirmation.
    - WHEN saving events THEN the system SHALL use structured JSON format with proper date, time, and description fields.
    - WHEN file operations fail THEN the system SHALL display clear error messages without data corruption.

### Requirement 7: Enhanced Event List Management
- **User Story**: As a user, I want comprehensive event management capabilities within the event list view, including adding, editing, and deleting events with visual feedback and keyboard navigation.
- **Acceptance Criteria**:
    - WHEN viewing an event list THEN the system SHALL support arrow key navigation with visual highlighting of the selected event.
    - WHEN the user presses 'A' in event list view THEN the system SHALL add a new event with inline input at the bottom of the list.
    - WHEN the user presses 'D' in event list view THEN the system SHALL delete the currently selected event after confirmation.
    - WHEN the user presses 'E' in event list view THEN the system SHALL edit the selected event with inline input showing current values.
    - WHEN adding a new event from the list THEN the system SHALL automatically select the newly added event upon successful creation.
    - WHEN events are displayed THEN the system SHALL show time in green and description in white (color terminals) with proper formatting.
    - WHEN the list exceeds display capacity THEN the system SHALL show a "... and X more events" indicator.

### Requirement 8: Application State Management and Exit Confirmation
- **User Story**: As a user, I want confirmation before accidentally exiting the application and clear feedback about the current application mode so that I don't lose work due to unintended key presses.
- **Acceptance Criteria**:
    - WHEN the user attempts to quit the application THEN the system SHALL prompt "Exit ASCII Calendar? (Enter: confirm, Esc: cancel)".
    - WHEN quitting from any application state THEN the system SHALL require confirmation before actually exiting.
    - WHEN the user presses Esc from the main calendar view THEN the system SHALL prompt for exit confirmation.
    - WHEN in specialized modes (search, event selection, etc.) THEN the system SHALL show appropriate key legends for the current mode.
    - WHEN transitioning between states THEN the system SHALL provide clear visual feedback about the current mode and available actions.

### Requirement 9: Color Support with Accessibility Compatibility
- **User Story**: As a user, I want enhanced visual appeal through color support while maintaining full functionality on monochrome terminals so that the application works well across different terminal environments.
- **Acceptance Criteria**:
    - WHEN the terminal supports color THEN the system SHALL use color themes for enhanced visual appeal (magenta headers, cyan day labels, green event indicators, yellow highlighting).
    - WHEN the terminal is monochrome THEN the system SHALL fall back to formatting using bold, reverse video, and spacing without losing information.
    - WHEN displaying month headers THEN the system SHALL use magenta bold text in color terminals or just bold in monochrome terminals.
    - WHEN showing day-of-week headers THEN the system SHALL use cyan text in color terminals or normal text in monochrome terminals.
    - WHEN indicating selected items THEN the system SHALL use yellow backgrounds in color terminals or reverse video in monochrome terminals.
    - WHEN all visual indicators are rendered THEN the system SHALL ensure no information is lost when color is unavailable.

### Requirement 10: Advanced Terminal Compatibility and Error Handling
- **User Story**: As a user, I want the application to work reliably across different terminal environments with graceful error handling and informative messages so that I can use it consistently regardless of my system configuration.
- **Acceptance Criteria**:
    - WHEN the terminal size is smaller than 80x24 THEN the system SHALL display a clear error message and exit gracefully.
    - WHEN terminal capabilities are detected THEN the system SHALL automatically adjust color usage and formatting options.
    - WHEN file system errors occur THEN the system SHALL display descriptive error messages without terminating unexpectedly.
    - WHEN JSON parsing fails THEN the system SHALL report the error and suggest corrective actions without data loss.
    - WHEN configuration files are malformed THEN the system SHALL fall back to default settings and continue operation.
    - WHEN migration from old format is needed THEN the system SHALL preserve the original file and create the new format safely.

### Requirement 11: Keyboard Input Processing and User Interface Consistency
- **User Story**: As a user, I want consistent keyboard controls across all application modes with clear feedback and comprehensive key legends so that I can efficiently navigate and operate the application.
- **Acceptance Criteria**:
    - WHEN in calendar mode THEN the system SHALL display the key legend: "B/N: month H/J/K/L: move Enter: events A: add D: delete E: edit C: current F: search Q: quit".
    - WHEN in search mode THEN the system SHALL display: "↑↓: navigate results Enter: go to date Esc: back to calendar F: search".
    - WHEN in event selection modes THEN the system SHALL display appropriate legends for the current operation (delete, edit, add).
    - WHEN the user presses any supported key THEN the system SHALL accept both uppercase and lowercase variants.
    - WHEN the user presses unsupported keys THEN the system SHALL ignore them without error messages or visual disturbance.
    - WHEN providing input prompts THEN the system SHALL show clear instructions and formatting requirements.

### Requirement 12: Data Persistence and Migration System
- **User Story**: As a user, I want reliable data storage with automatic format migration and backup preservation so that my event data remains safe and accessible across application updates.
- **Acceptance Criteria**:
    - WHEN the application stores events THEN the system SHALL use structured JSON format with "events" array containing date, time, and description fields.
    - WHEN loading events THEN the system SHALL parse dates in local timezone to ensure consistency with display and navigation.
    - WHEN an old events.txt file is detected THEN the system SHALL automatically convert it to JSON format and preserve the original file.
    - WHEN migration occurs THEN the system SHALL display "Successfully migrated X events from events.txt to events.json".
    - WHEN saving events THEN the system SHALL write formatted JSON with proper indentation for human readability.
    - WHEN concurrent file access issues occur THEN the system SHALL handle them gracefully without data corruption.

## Technical Specifications

### Compatibility Requirements
- Terminal size minimum: 80 characters wide × 24 lines tall
- Color support: Optional, with full functionality in monochrome mode
- Platform support: Linux, macOS, Windows (with appropriate terminal)
- Go version: 1.19 or later for building from source

### Security Considerations
- File permissions: Configuration directory created with 0755 permissions
- Data privacy: All data stored locally, no network communication
- Input validation: All user input validated before processing
- Configuration safety: Malformed configs fall back to safe defaults

## Conclusion
This enhanced requirements document reflects the evolution of the ASCII Calendar application from a basic terminal calendar to a comprehensive event management system. The application now provides professional-grade functionality while maintaining the simplicity and efficiency of terminal-based operation, supporting both novice and advanced users across diverse terminal environments.