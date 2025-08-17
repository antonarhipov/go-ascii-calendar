package main

import (
	"fmt"
	"log"

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
		renderer:   terminal.NewRenderer(term, eventManager),
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
		return true // Exit application

	case terminal.ActionBack:
		return true // Exit application when Esc is pressed on main screen

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
		// Enter event add mode in calendar view
		app.state = StateCalendarEventAdd
		app.selectedEventIndex = 0

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
		app.processEditEvent()

	case terminal.ActionResetCurrent:
		app.navigation.ResetToCurrent()
	}

	return false
}

// handleCalendarEventSelectionAction handles actions when selecting events in calendar view
func (app *Application) handleCalendarEventSelectionAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return true // Exit application

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
		return true // Exit application

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

// handleEventListAction handles actions when viewing events
func (app *Application) handleEventListAction(action terminal.KeyAction) bool {
	switch action {
	case terminal.ActionQuit:
		return true // Exit application

	case terminal.ActionBack:
		app.state = StateCalendar
		app.selectedEventIndex = 0 // Reset event selection

	case terminal.ActionMoveUp:
		app.navigateEventUp()

	case terminal.ActionMoveDown:
		app.navigateEventDown()

	case terminal.ActionAddEvent:
		app.processAddEvent()

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
		return true // Exit application

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

	// Get time input
	timeStr, ok := app.input.GetTextInputWithPrompt("Enter time (HH:MM):", 5, app.renderer)
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

	// Get new time input (default to current time)
	currentTime := eventToEdit.GetTimeString()
	prompt := fmt.Sprintf("Enter new time (current: %s):", currentTime)
	timeStr, ok := app.input.GetTextInputWithPrompt(prompt, 5, app.renderer)
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

// processEditEventFromList edits the currently selected event from the events list
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

	// Get new time input (default to current time)
	currentTime := eventToEdit.GetTimeString()
	prompt := fmt.Sprintf("Enter new time (current: %s):", currentTime)
	timeStr, ok := app.input.GetTextInputWithPrompt(prompt, 5, app.renderer)
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
	err := app.events.EditEvent(eventToEdit, selectedDate, timeStr, description)
	if err != nil {
		app.showError(fmt.Sprintf("Error editing event: %v", err))
	} else {
		app.showMessage("Event edited successfully!")
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

	// Get time input using inline input
	timeStr, ok := app.input.GetInlineTextInput(eventsLeftX, addEventY, "Time:", 5, app.renderer)
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
