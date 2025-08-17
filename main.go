package main

import (
	"fmt"
	"log"
	"time"

	"go-ascii-calendar/events"
	"go-ascii-calendar/models"
	"go-ascii-calendar/terminal"
)

// AppState represents the current state of the application
type AppState int

const (
	StateCalendar AppState = iota
	StateEventList
	StateAddEvent
)

// Application holds the main application components
type Application struct {
	terminal   *terminal.Terminal
	renderer   *terminal.Renderer
	input      *terminal.InputHandler
	navigation *terminal.NavigationController
	events     *events.Manager
	calendar   *models.Calendar
	selection  *models.Selection
	state      AppState
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	term := terminal.NewTerminal()
	eventManager := events.NewManager()
	cal := models.NewCalendar()
	sel := models.NewSelection(cal)

	return &Application{
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

	case terminal.ActionAddEvent:
		app.state = StateAddEvent

	case terminal.ActionResetCurrent:
		app.navigation.ResetToCurrent()
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

	case terminal.ActionAddEvent:
		app.state = StateAddEvent
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

	case StateEventList:
		selectedDate := app.navigation.GetCurrentSelection()
		eventList := app.events.GetEventsForDate(selectedDate)
		return app.renderer.RenderEventList(selectedDate, eventList)

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

// showError displays an error message
func (app *Application) showError(message string) {
	app.renderer.RenderMessage(message, true)
	app.terminal.Flush()

	// Wait for a short time to let user see the message
	time.Sleep(2 * time.Second)
}

// showMessage displays a success message
func (app *Application) showMessage(message string) {
	app.renderer.RenderMessage(message, false)
	app.terminal.Flush()

	// Wait for a short time to let user see the message
	time.Sleep(2 * time.Second)
}

func main() {
	app := NewApplication()

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	fmt.Println("ASCII Calendar - Goodbye!")
}
