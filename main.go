package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsf/termbox-go"
	"go-ascii-calendar/config"
	"go-ascii-calendar/events"
	"go-ascii-calendar/models"
	"go-ascii-calendar/terminal"
)

// AppState represents the current state of the application
type AppState int

const (
	StateCalendar               AppState = iota
	StateCalendarEventSelection          // New state for selecting events within calendar view
	StateCalendarEventAdd                // New state for adding events within calendar view
	StateCalendarEventEdit               // New state for editing events within calendar view
	StateSearch                          // New state for search functionality
	StateEventList
	StateAddEvent
)

// Application holds the main application components
type Application struct {
	config             *config.Config
	terminal           *terminal.Terminal
	renderer           *terminal.Renderer
	input              *terminal.InputHandler
	navigation         *terminal.NavigationController
	events             *events.Manager
	calendar           *models.Calendar
	selection          *models.Selection
	state              AppState
	selectedEventIndex int // Index of currently selected event in events view
	// Search-related fields
	searchQuery         string         // Current search query
	searchResults       []models.Event // Search results
	searchResultDates   []string       // Unique dates from search results for grouping
	selectedResultIndex int            // Index of currently selected search result
}

// NewApplication creates a new application instance with configuration
func NewApplication(cfg *config.Config) *Application {
	term := terminal.NewTerminal()
	eventManager := events.NewManagerWithConfig(cfg)
	cal := models.NewCalendar()
	sel := models.NewSelection(cal)

	return &Application{
		config:     cfg,
		terminal:   term,
		renderer:   terminal.NewRenderer(term, eventManager, cfg),
		input:      terminal.NewInputHandler(term),
		navigation: terminal.NewNavigationController(cal, sel),
		events:     eventManager,
		calendar:   cal,
		selection:  sel,
		state:      StateCalendar,
	}
}

// Initialize initializes the application
func (app *Application) Initialize() error {
	// Initialize terminal
	if err := app.terminal.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize terminal: %v", err)
	}

	// Check terminal size
	if !app.terminal.CheckSize() {
		app.terminal.Close()
		return fmt.Errorf("terminal too small - minimum 80x24 required")
	}

	// Load events from storage
	if err := app.events.LoadEvents(); err != nil {
		app.terminal.Close()
		return fmt.Errorf("failed to load events: %v", err)
	}

	return nil
}

// Run starts the main application loop
func (app *Application) Run() error {
	defer app.terminal.Close()

	// Initial render
	if err := app.renderCurrentView(); err != nil {
		return fmt.Errorf("initial render failed: %v", err)
	}

	// Main event loop
	for {
		// Wait for user input
		event := app.input.WaitForKey()
		action := app.input.ProcessKeyEvent(event)

		// Handle the action based on current state
		shouldExit := app.handleAction(action)
		if shouldExit {
			break
		}

		// Re-render the current view
		if err := app.renderCurrentView(); err != nil {
			app.showError(fmt.Sprintf("Render error: %v", err))
		}
	}

	return nil
}

// handleAction handles the given action based on current state
func (app *Application) handleAction(action terminal.KeyAction) bool {
	switch app.state {
	case StateCalendar:
		return app.handleCalendarAction(action)
	case StateCalendarEventSelection:
		return app.handleCalendarEventSelectionAction(action)
	case StateCalendarEventAdd:
		return app.handleCalendarEventAddAction(action)
	case StateCalendarEventEdit:
		return app.handleCalendarEventEditAction(action)
	case StateSearch:
		return app.handleSearchAction(action)
	case StateEventList:
		return app.handleEventListAction(action)
	case StateAddEvent:
		return app.handleAddEventAction(action)
	}
	return false
}

// handleCalendarAction handles actions when in calendar view
func (app *Application) handleCalendarAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		return app.confirmExit() // Exit application when Esc is pressed on main screen

	case terminal.ActionMonthPrev:
		app.navigation.NavigateMonthBackward()

	case terminal.ActionMonthNext:
		app.navigation.NavigateMonthForward()

	case terminal.ActionMoveLeft:
		app.navigation.NavigateDayLeft()

	case terminal.ActionMoveRight:
		app.navigation.NavigateDayRight()

	case terminal.ActionMoveUp:
		app.navigation.NavigateDayUp()

	case terminal.ActionMoveDown:
		app.navigation.NavigateDayDown()

	case terminal.ActionShowEvents:
		app.state = StateEventList
		app.selectedEventIndex = 0 // Initialize event selection

	case terminal.ActionAddEvent:
		// Directly start adding event from calendar view
		app.processAddEventFromCalendar()

	case terminal.ActionDeleteEvent:
		// Enter event selection mode in calendar view
		selectedDate := app.navigation.GetCurrentSelection()
		events := app.events.GetEventsForDate(selectedDate)
		if len(events) > 0 {
			app.state = StateCalendarEventSelection
			app.selectedEventIndex = 0
		} else {
			app.showError("No events to delete on this date")
		}

	case terminal.ActionEditEvent:
		// Enter event edit selection mode in calendar view
		selectedDate := app.navigation.GetCurrentSelection()
		events := app.events.GetEventsForDate(selectedDate)
		if len(events) > 0 {
			app.state = StateCalendarEventEdit
			app.selectedEventIndex = 0
		} else {
			app.showError("No events to edit on this date")
		}

	case terminal.ActionResetCurrent:
		app.navigation.ResetToCurrent()

	case terminal.ActionSearch:
		app.processSearch()
	}

	return false
}

// handleSearchAction handles actions when in search mode
func (app *Application) handleSearchAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		// Exit search mode and return to calendar
		app.state = StateCalendar
		app.searchQuery = ""
		app.searchResults = nil
		app.searchResultDates = nil
		app.selectedResultIndex = 0

	case terminal.ActionMoveUp:
		app.navigateSearchResultUp()

	case terminal.ActionMoveDown:
		app.navigateSearchResultDown()

	case terminal.ActionShowEvents:
		// Enter key - navigate to selected date and close search
		app.processSearchResultSelection()

	default:
		// For other keys, ignore them in search mode
		return false
	}

	return false
}

// handleCalendarEventSelectionAction handles actions when selecting events in calendar view
func (app *Application) handleCalendarEventSelectionAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		// Exit event selection mode and return to calendar navigation
		app.state = StateCalendar
		app.selectedEventIndex = 0

	case terminal.ActionMoveUp:
		app.navigateCalendarEventUp()

	case terminal.ActionMoveDown:
		app.navigateCalendarEventDown()

	case terminal.ActionShowEvents:
		// Enter key - confirm deletion of selected event
		app.processDeleteSelectedCalendarEvent()

	default:
		// For other keys, ignore them in event selection mode
		return false
	}

	return false
}

// handleCalendarEventAddAction handles actions when adding events in calendar view
func (app *Application) handleCalendarEventAddAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		// Exit event add mode and return to calendar navigation
		app.state = StateCalendar
		app.selectedEventIndex = 0

	case terminal.ActionShowEvents:
		// Enter key - start adding the event
		app.processAddEventFromCalendar()

	default:
		// For other keys, ignore them in event add mode
		return false
	}

	return false
}

// handleCalendarEventEditAction handles actions when editing events in calendar view
func (app *Application) handleCalendarEventEditAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		// Exit event edit mode and return to calendar navigation
		app.state = StateCalendar
		app.selectedEventIndex = 0

	case terminal.ActionMoveUp:
		app.navigateCalendarEventEditUp()

	case terminal.ActionMoveDown:
		app.navigateCalendarEventEditDown()

	case terminal.ActionShowEvents:
		// Enter key - confirm editing of selected event
		app.processEditSelectedCalendarEvent()

	default:
		// For other keys, ignore them in event edit mode
		return false
	}

	return false
}

// handleEventListAction handles actions when viewing events
func (app *Application) handleEventListAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		app.state = StateCalendar
		app.selectedEventIndex = 0 // Reset event selection

	case terminal.ActionMoveUp:
		app.navigateEventUp()

	case terminal.ActionMoveDown:
		app.navigateEventDown()

	case terminal.ActionAddEvent:
		app.processAddEventFromEventsList()

	case terminal.ActionDeleteEvent:
		app.processDeleteEventFromList()

	case terminal.ActionEditEvent:
		app.processEditEventFromList()
	}

	return false
}

// handleAddEventAction handles actions when adding events
func (app *Application) handleAddEventAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return app.confirmExit() // Exit application with confirmation

	case terminal.ActionBack:
		app.state = StateCalendar
	}

	// For adding events, we handle the input differently
	app.processAddEvent()
	app.state = StateCalendar

	return false
}

// renderCurrentView renders the appropriate view based on current state
func (app *Application) renderCurrentView() error {
	switch app.state {
	case StateCalendar:
		return app.renderer.RenderCalendar(app.calendar, app.selection)

	case StateCalendarEventSelection:
		// Render calendar with event selection highlighting
		return app.renderer.RenderCalendarWithEventSelection(app.calendar, app.selection, app.selectedEventIndex)

	case StateCalendarEventAdd:
		// Render calendar with event add highlighting
		return app.renderer.RenderCalendarWithEventAdd(app.calendar, app.selection)

	case StateCalendarEventEdit:
		// Render calendar with event edit highlighting
		return app.renderer.RenderCalendarWithEventEdit(app.calendar, app.selection, app.selectedEventIndex)

	case StateSearch:
		// Render calendar with search results
		return app.renderer.RenderCalendarWithSearch(app.calendar, app.selection, app.searchQuery, app.searchResults, app.searchResultDates, app.selectedResultIndex)

	case StateEventList:
		selectedDate := app.navigation.GetCurrentSelection()
		eventList := app.events.GetEventsForDate(selectedDate)
		return app.renderer.RenderEventList(selectedDate, eventList, app.selectedEventIndex)

	case StateAddEvent:
		// This state is handled differently - we don't render here
		// but in processAddEvent()
		return nil
	}

	return nil
}

// processAddEvent handles the event addition workflow
func (app *Application) processAddEvent() {
	selectedDate := app.navigation.GetCurrentSelection()

	// Get time input with validation
	timeStr, ok := app.input.GetTimeInput("Enter time (HH:MM):", app.renderer)
	if !ok {
		return // User cancelled
	}

	// Get description input
	description, ok := app.input.GetTextInputWithPrompt("Enter description:", 100, app.renderer)
	if !ok {
		return // User cancelled
	}

	// Add the event
	err := app.events.AddEvent(selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error adding event: %v", err))
	} else {
		app.showMessage("Event added successfully!")
	}
}

// processDeleteEvent handles the event deletion workflow
func (app *Application) processDeleteEvent() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) == 0 {
		app.showError("No events to delete on this date")
		return
	}

	if len(events) == 1 {
		// Only one event, delete it directly after confirmation
		event := events[0]
		confirmMsg := fmt.Sprintf("Delete event: %s - %s? (Enter: confirm, Esc: cancel)", event.GetTimeString(), event.Description)

		if app.confirmAction(confirmMsg) {
			err := app.events.DeleteEvent(event)
			if err != nil {
				app.showError(fmt.Sprintf("Error deleting event: %v", err))
			} else {
				app.showMessage("Event deleted successfully!")
			}
		}
		return
	}

	// Multiple events - let user select which one to delete
	selectedEvent := app.selectEventFromList(events, "Select event to delete:")
	if selectedEvent != nil {
		confirmMsg := fmt.Sprintf("Delete event: %s - %s? (Enter: confirm, Esc: cancel)", selectedEvent.GetTimeString(), selectedEvent.Description)

		if app.confirmAction(confirmMsg) {
			err := app.events.DeleteEvent(*selectedEvent)
			if err != nil {
				app.showError(fmt.Sprintf("Error deleting event: %v", err))
			} else {
				app.showMessage("Event deleted successfully!")
			}
		}
	}
}

// processEditEvent handles the event editing workflow
func (app *Application) processEditEvent() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) == 0 {
		app.showError("No events to edit on this date")
		return
	}

	var eventToEdit *models.Event
	if len(events) == 1 {
		// Only one event, edit it directly
		eventToEdit = &events[0]
	} else {
		// Multiple events - let user select which one to edit
		eventToEdit = app.selectEventFromList(events, "Select event to edit:")
		if eventToEdit == nil {
			return // User cancelled selection
		}
	}

	// Get new time input with validation (default to current time)
	currentTime := eventToEdit.GetTimeString()
	prompt := fmt.Sprintf("Enter new time (current: %s):", currentTime)
	timeStr, ok := app.input.GetTimeInput(prompt, app.renderer)
	if !ok {
		return // User cancelled
	}

	// If user entered empty time, keep the current time
	if timeStr == "" {
		timeStr = currentTime
	}

	// Get new description input (default to current description)
	currentDesc := eventToEdit.Description
	prompt = fmt.Sprintf("Enter new description (current: %s):", currentDesc)
	description, ok := app.input.GetTextInputWithPrompt(prompt, 100, app.renderer)
	if !ok {
		return // User cancelled
	}

	// If user entered empty description, keep the current description
	if description == "" {
		description = currentDesc
	}

	// Update the event
	err := app.events.EditEvent(*eventToEdit, selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error editing event: %v", err))
	} else {
		app.showMessage("Event edited successfully!")
	}
}

// navigateEventUp moves selection up in the event list
func (app *Application) navigateEventUp() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) > 0 && app.selectedEventIndex > 0 {
		app.selectedEventIndex--
	}
}

// navigateEventDown moves selection down in the event list
func (app *Application) navigateEventDown() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) > 0 && app.selectedEventIndex < len(events)-1 {
		app.selectedEventIndex++
	}
}

// processDeleteEventFromList deletes the currently selected event from the events list
func (app *Application) processDeleteEventFromList() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) == 0 {
		app.showError("No events to delete on this date")
		return
	}

	if app.selectedEventIndex >= len(events) {
		app.selectedEventIndex = len(events) - 1
	}

	event := events[app.selectedEventIndex]
	confirmMsg := fmt.Sprintf("Delete event: %s - %s? (Enter: confirm, Esc: cancel)", event.GetTimeString(), event.Description)

	if app.confirmAction(confirmMsg) {
		err := app.events.DeleteEvent(event)
		if err != nil {
			app.showError(fmt.Sprintf("Error deleting event: %v", err))
		} else {
			app.showMessage("Event deleted successfully!")
			// Adjust selection if we deleted the last event
			if app.selectedEventIndex >= len(events)-1 && app.selectedEventIndex > 0 {
				app.selectedEventIndex--
			}
		}
	}
}

// processEditEventFromList edits the currently selected event from the events list using inline input
func (app *Application) processEditEventFromList() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) == 0 {
		app.showError("No events to edit on this date")
		return
	}

	if app.selectedEventIndex >= len(events) {
		app.selectedEventIndex = len(events) - 1
	}

	eventToEdit := events[app.selectedEventIndex]

	// Calculate coordinates for inline input on the selected event
	// Events view has title at Y=2, separator at Y=4, events start at Y=6
	startY := 6
	editEventY := startY + app.selectedEventIndex // Position of the selected event
	eventsLeftX := 2                              // Use left margin like the event list

	// Get new time input with current value as default using inline input with validation
	currentTime := eventToEdit.GetTimeString()
	timeStr, ok := app.input.GetInlineTimeInputWithDefault(eventsLeftX, editEventY, "Time:", currentTime, app.renderer)
	if !ok {
		return // User cancelled
	}

	// If user entered empty time, keep the current time
	if timeStr == "" {
		timeStr = currentTime
	}

	// Get new description input with current value as default using inline input
	currentDesc := eventToEdit.Description
	description, ok := app.input.GetInlineTextInputWithDefault(eventsLeftX, editEventY, "Description:", 100, currentDesc, app.renderer)
	if !ok {
		return // User cancelled
	}

	// If user entered empty description, keep the current description
	if description == "" {
		description = currentDesc
	}

	// Update the event
	err := app.events.EditEvent(eventToEdit, selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error editing event: %v", err))
	} else {
		app.showMessage("Event edited successfully!")
	}
}

// processAddEventFromEventsList handles adding an event from the events view with inline input
func (app *Application) processAddEventFromEventsList() {
	selectedDate := app.navigation.GetCurrentSelection()

	// Calculate coordinates for inline input in events view
	// Events view has title at Y=2, separator at Y=4, events start at Y=6
	// We want to add the new event at the bottom of the existing events list
	events := app.events.GetEventsForDate(selectedDate)
	startY := 6
	addEventY := startY + len(events) // Position after existing events

	// Use left margin like the event list (X=2)
	eventsLeftX := 2

	// Get time input using inline input with validation
	timeStr, ok := app.input.GetInlineTimeInput(eventsLeftX, addEventY, "Time:", app.renderer)
	if !ok {
		// User cancelled
		return
	}

	// Get description input using inline input
	description, ok := app.input.GetInlineTextInput(eventsLeftX, addEventY, "Description:", 100, app.renderer)
	if !ok {
		// User cancelled
		return
	}

	// Add the event
	err := app.events.AddEvent(selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error adding event: %v", err))
	} else {
		app.showMessage("Event added successfully!")

		// After adding the event, select and highlight the newly added event
		// Get the updated events list
		updatedEvents := app.events.GetEventsForDate(selectedDate)

		// Find the newly added event (it should be the one with matching time and description)
		for i, event := range updatedEvents {
			if event.GetTimeString() == timeStr && event.Description == description {
				app.selectedEventIndex = i
				break
			}
		}

		// If we couldn't find it (unlikely), select the last event
		if app.selectedEventIndex >= len(updatedEvents) {
			app.selectedEventIndex = len(updatedEvents) - 1
		}
	}
}

// processAddEventFromCalendar handles adding an event from the calendar view with inline input
func (app *Application) processAddEventFromCalendar() {
	selectedDate := app.navigation.GetCurrentSelection()

	// Calculate coordinates for inline input (same as renderSelectedDateEventsWithAddMode)
	width, _ := app.terminal.GetSize()
	totalWidth := 3*24 + 2*2 // monthWidth=24, monthSpacing=2 (from renderer)
	startX := (width - totalWidth) / 2
	eventsLeftX := startX + 1
	eventsStartY := 13

	// Get existing events to calculate the Y position
	events := app.events.GetEventsForDate(selectedDate)
	maxExistingEvents := 9
	if len(events) > maxExistingEvents {
		maxExistingEvents = 9
	}
	addEventY := eventsStartY + 1 + maxExistingEvents

	// Get time input using inline input with validation
	timeStr, ok := app.input.GetInlineTimeInput(eventsLeftX, addEventY, "Time:", app.renderer)
	if !ok {
		// User cancelled, return to calendar
		app.state = StateCalendar
		app.selectedEventIndex = 0
		return
	}

	// Get description input using inline input
	description, ok := app.input.GetInlineTextInput(eventsLeftX, addEventY, "Description:", 100, app.renderer)
	if !ok {
		// User cancelled, return to calendar
		app.state = StateCalendar
		app.selectedEventIndex = 0
		return
	}

	// Add the event
	err := app.events.AddEvent(selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error adding event: %v", err))
	} else {
		app.showMessage("Event added successfully!")
	}

	// Return to calendar view
	app.state = StateCalendar
	app.selectedEventIndex = 0
}

// navigateCalendarEventUp moves selection up in the calendar events list
func (app *Application) navigateCalendarEventUp() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) > 0 && app.selectedEventIndex > 0 {
		app.selectedEventIndex--
	}
}

// navigateCalendarEventDown moves selection down in the calendar events list
func (app *Application) navigateCalendarEventDown() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) > 0 && app.selectedEventIndex < len(events)-1 {
		app.selectedEventIndex++
	}
}

// navigateCalendarEventEditUp moves selection up in the calendar events list for editing
func (app *Application) navigateCalendarEventEditUp() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) > 0 && app.selectedEventIndex > 0 {
		app.selectedEventIndex--
	}
}

// navigateCalendarEventEditDown moves selection down in the calendar events list for editing
func (app *Application) navigateCalendarEventEditDown() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) > 0 && app.selectedEventIndex < len(events)-1 {
		app.selectedEventIndex++
	}
}

// processDeleteSelectedCalendarEvent deletes the currently selected event in calendar view
func (app *Application) processDeleteSelectedCalendarEvent() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) == 0 {
		// No events to delete, exit selection mode
		app.state = StateCalendar
		app.selectedEventIndex = 0
		return
	}

	if app.selectedEventIndex >= len(events) {
		app.selectedEventIndex = len(events) - 1
	}

	event := events[app.selectedEventIndex]
	confirmMsg := fmt.Sprintf("Delete event: %s - %s? (Enter: confirm, Esc: cancel)", event.GetTimeString(), event.Description)

	if app.confirmAction(confirmMsg) {
		err := app.events.DeleteEvent(event)
		if err != nil {
			app.showError(fmt.Sprintf("Error deleting event: %v", err))
		} else {
			app.showMessage("Event deleted successfully!")
			// Adjust selection if we deleted the last event
			if app.selectedEventIndex >= len(events)-1 && app.selectedEventIndex > 0 {
				app.selectedEventIndex--
			}
		}
	}

	// Return to calendar navigation after deletion attempt
	app.state = StateCalendar
	app.selectedEventIndex = 0
}

// processEditSelectedCalendarEvent edits the currently selected event in calendar view
func (app *Application) processEditSelectedCalendarEvent() {
	selectedDate := app.navigation.GetCurrentSelection()
	events := app.events.GetEventsForDate(selectedDate)

	if len(events) == 0 {
		// No events to edit, exit edit mode
		app.state = StateCalendar
		app.selectedEventIndex = 0
		return
	}

	if app.selectedEventIndex >= len(events) {
		app.selectedEventIndex = len(events) - 1
	}

	eventToEdit := events[app.selectedEventIndex]

	// Calculate coordinates for inline input (same as add mode)
	width, _ := app.terminal.GetSize()
	totalWidth := 3*24 + 2*2 // monthWidth=24, monthSpacing=2 (from renderer)
	startX := (width - totalWidth) / 2
	eventsLeftX := startX + 1
	eventsStartY := 13

	// Calculate Y position for the selected event
	editEventY := eventsStartY + 1 + app.selectedEventIndex

	// Get new time input with current value as default using validation
	currentTime := eventToEdit.GetTimeString()
	timeStr, ok := app.input.GetInlineTimeInputWithDefault(eventsLeftX, editEventY, "Time:", currentTime, app.renderer)
	if !ok {
		// User cancelled, return to calendar
		app.state = StateCalendar
		app.selectedEventIndex = 0
		return
	}

	// If user entered empty time, keep the current time
	if timeStr == "" {
		timeStr = currentTime
	}

	// Get new description input with current value as default
	currentDesc := eventToEdit.Description
	description, ok := app.input.GetInlineTextInputWithDefault(eventsLeftX, editEventY, "Description:", 100, currentDesc, app.renderer)
	if !ok {
		// User cancelled, return to calendar
		app.state = StateCalendar
		app.selectedEventIndex = 0
		return
	}

	// If user entered empty description, keep the current description
	if description == "" {
		description = currentDesc
	}

	// Update the event
	err := app.events.EditEvent(eventToEdit, selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error editing event: %v", err))
	} else {
		app.showMessage("Event edited successfully!")
	}

	// Return to calendar view
	app.state = StateCalendar
	app.selectedEventIndex = 0
}

// showError displays an error message
func (app *Application) showError(message string) {
	app.renderer.RenderMessage(message, true)
	app.terminal.Flush()
}

// showMessage displays a success message
func (app *Application) showMessage(message string) {
	app.renderer.RenderMessage(message, false)
	app.terminal.Flush()
}

// confirmAction prompts the user for confirmation (Enter/Esc)
func (app *Application) confirmAction(message string) bool {
	// Display the confirmation message
	app.renderer.RenderMessage(message, false)
	app.terminal.Flush()

	// Wait for user input
	event := app.input.WaitForKey()

	// Check for Enter key (confirm) or Esc key (cancel)
	if event.Key == termbox.KeyEnter {
		return true
	}

	return false // Any other key (including Esc) cancels
}

// confirmExit prompts the user to confirm application exit
func (app *Application) confirmExit() bool {
	return app.confirmAction("Exit ASCII Calendar? (Enter: confirm, Esc: cancel)")
}

// selectEventFromList allows the user to select an event from a list
func (app *Application) selectEventFromList(events []models.Event, title string) *models.Event {
	if len(events) == 0 {
		return nil
	}

	// Simple implementation: show list with numbers and get user input
	app.terminal.Clear()

	fg, bg := app.terminal.GetDefaultColors()

	// Display title
	app.terminal.PrintCentered(2, title, termbox.AttrBold, bg)

	// Display events with numbers
	startY := 4
	for i, event := range events {
		eventText := fmt.Sprintf("%d. %s - %s", i+1, event.GetTimeString(), event.Description)
		// Truncate if too long
		if len(eventText) > 70 {
			eventText = eventText[:67] + "..."
		}
		app.terminal.PrintCentered(startY+i, eventText, fg, bg)
	}

	// Instructions
	instrY := startY + len(events) + 2
	app.terminal.PrintCentered(instrY, "Enter event number (1-9) or Esc to cancel:", fg, bg)

	app.terminal.Flush()

	// Wait for user selection
	for {
		event := app.input.WaitForKey()

		if event.Type != termbox.EventKey {
			continue
		}

		// Handle Esc key
		if event.Key == termbox.KeyEsc {
			return nil
		}

		// Handle number keys
		if event.Ch >= '1' && event.Ch <= '9' {
			choice := int(event.Ch - '1') // Convert to 0-based index
			if choice < len(events) {
				return &events[choice]
			}
		}
	}
}

// processSearch handles the search functionality workflow
func (app *Application) processSearch() {
	// Get search query input
	query, ok := app.input.GetTextInputWithPrompt("Enter search query:", 100, app.renderer)
	if !ok {
		return // User cancelled
	}

	// Perform search
	app.searchQuery = query
	app.searchResults = app.events.SearchEvents(query)
	app.selectedResultIndex = 0

	// Build unique dates list for grouping
	app.searchResultDates = make([]string, 0)
	datesSeen := make(map[string]bool)
	for _, event := range app.searchResults {
		dateStr := event.Date.Format("2006-01-02")
		if !datesSeen[dateStr] {
			app.searchResultDates = append(app.searchResultDates, dateStr)
			datesSeen[dateStr] = true
		}
	}

	// Switch to search mode
	app.state = StateSearch
}

// navigateSearchResultUp moves selection up in the search results
func (app *Application) navigateSearchResultUp() {
	if len(app.searchResults) > 0 && app.selectedResultIndex > 0 {
		app.selectedResultIndex--
	}
}

// navigateSearchResultDown moves selection down in the search results
func (app *Application) navigateSearchResultDown() {
	if len(app.searchResults) > 0 && app.selectedResultIndex < len(app.searchResults)-1 {
		app.selectedResultIndex++
	}
}

// processSearchResultSelection handles Enter key in search mode
func (app *Application) processSearchResultSelection() {
	if len(app.searchResults) == 0 {
		return
	}

	// Get the selected search result
	selectedEvent := app.searchResults[app.selectedResultIndex]

	// Navigate calendar to the event's date
	app.navigation.SetSelection(selectedEvent.Date)

	// Update calendar to show the month containing this date
	eventMonth := selectedEvent.Date
	app.calendar.CurrentMonth = time.Date(eventMonth.Year(), eventMonth.Month(), 1, 0, 0, 0, 0, eventMonth.Location())

	// Clear search state and return to calendar
	app.searchQuery = ""
	app.searchResults = nil
	app.searchResultDates = nil
	app.selectedResultIndex = 0
	app.state = StateCalendar
}

func main() {
	// Load configuration from command line and config file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create application with configuration
	app := NewApplication(cfg)

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	fmt.Println("ASCII Calendar - Goodbye!")
}
