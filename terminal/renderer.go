package terminal

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
	"go-ascii-calendar/calendar"
	"go-ascii-calendar/events"
	"go-ascii-calendar/models"
)

// Renderer handles calendar rendering operations
type Renderer struct {
	terminal     *Terminal
	eventManager *events.Manager
	monthWidth   int // Width of each month display
	monthSpacing int // Spacing between months
}

// NewRenderer creates a new calendar renderer
func NewRenderer(terminal *Terminal, eventManager *events.Manager) *Renderer {
	return &Renderer{
		terminal:     terminal,
		eventManager: eventManager,
		monthWidth:   24, // Width for each month (includes padding)
		monthSpacing: 2,  // Space between months
	}
}

// RenderCalendar renders the three-month calendar view
func (r *Renderer) RenderCalendar(cal *models.Calendar, selection *models.Selection) error {
	r.terminal.Clear()

	// Get terminal size
	width, height := r.terminal.GetSize()
	if width < 80 || height < 24 {
		r.terminal.PrintCentered(height/2, "Terminal too small! Minimum 80x24 required.",
			termbox.ColorRed, termbox.ColorDefault)
		return r.terminal.Flush()
	}

	// Calculate starting positions for three months
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2

	prevMonth := cal.GetPreviousMonth()
	currentMonth := cal.CurrentMonth
	nextMonth := cal.GetNextMonth()

	months := []time.Time{prevMonth, currentMonth, nextMonth}

	// Render each month
	for i, month := range months {
		x := startX + i*(r.monthWidth+r.monthSpacing)
		err := r.renderMonth(month, x, 2, selection)
		if err != nil {
			return err
		}
	}

	// Render events for selected date
	r.renderSelectedDateEvents(selection.SelectedDate)

	// Render key legend
	r.renderKeyLegend()

	return r.terminal.Flush()
}

// RenderCalendarWithEventSelection renders the calendar with event selection highlighting
func (r *Renderer) RenderCalendarWithEventSelection(cal *models.Calendar, selection *models.Selection, selectedEventIndex int) error {
	r.terminal.Clear()

	// Get terminal size
	width, height := r.terminal.GetSize()
	if width < 80 || height < 24 {
		r.terminal.PrintCentered(height/2, "Terminal too small! Minimum 80x24 required.",
			termbox.ColorRed, termbox.ColorDefault)
		return r.terminal.Flush()
	}

	// Calculate starting positions for three months
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2

	prevMonth := cal.GetPreviousMonth()
	currentMonth := cal.CurrentMonth
	nextMonth := cal.GetNextMonth()

	months := []time.Time{prevMonth, currentMonth, nextMonth}

	// Render each month
	for i, month := range months {
		x := startX + i*(r.monthWidth+r.monthSpacing)
		err := r.renderMonth(month, x, 2, selection)
		if err != nil {
			return err
		}
	}

	// Render events for selected date with selection highlighting
	r.renderSelectedDateEventsWithSelection(selection.SelectedDate, selectedEventIndex)

	// Render key legend for event selection mode
	r.renderEventSelectionKeyLegend()

	return r.terminal.Flush()
}

// RenderCalendarWithEventAdd renders the calendar with event add highlighting
func (r *Renderer) RenderCalendarWithEventAdd(cal *models.Calendar, selection *models.Selection) error {
	r.terminal.Clear()

	// Get terminal size
	width, height := r.terminal.GetSize()
	if width < 80 || height < 24 {
		r.terminal.PrintCentered(height/2, "Terminal too small! Minimum 80x24 required.",
			termbox.ColorRed, termbox.ColorDefault)
		return r.terminal.Flush()
	}

	// Calculate starting positions for three months
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2

	prevMonth := cal.GetPreviousMonth()
	currentMonth := cal.CurrentMonth
	nextMonth := cal.GetNextMonth()

	months := []time.Time{prevMonth, currentMonth, nextMonth}

	// Render each month
	for i, month := range months {
		x := startX + i*(r.monthWidth+r.monthSpacing)
		err := r.renderMonth(month, x, 2, selection)
		if err != nil {
			return err
		}
	}

	// Render events for selected date with add mode highlighting
	r.renderSelectedDateEventsWithAddMode(selection.SelectedDate)

	// Render key legend for event add mode
	r.renderEventAddKeyLegend()

	return r.terminal.Flush()
}

// RenderCalendarWithEventEdit renders the calendar with event edit highlighting
func (r *Renderer) RenderCalendarWithEventEdit(cal *models.Calendar, selection *models.Selection, selectedEventIndex int) error {
	r.terminal.Clear()

	// Get terminal size
	width, height := r.terminal.GetSize()
	if width < 80 || height < 24 {
		r.terminal.PrintCentered(height/2, "Terminal too small! Minimum 80x24 required.",
			termbox.ColorRed, termbox.ColorDefault)
		return r.terminal.Flush()
	}

	// Calculate starting positions for three months
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2

	prevMonth := cal.GetPreviousMonth()
	currentMonth := cal.CurrentMonth
	nextMonth := cal.GetNextMonth()

	months := []time.Time{prevMonth, currentMonth, nextMonth}

	// Render each month
	for i, month := range months {
		x := startX + i*(r.monthWidth+r.monthSpacing)
		err := r.renderMonth(month, x, 2, selection)
		if err != nil {
			return err
		}
	}

	// Render events for selected date with edit mode highlighting
	r.renderSelectedDateEventsWithEditMode(selection.SelectedDate, selectedEventIndex)

	// Render key legend for event edit mode
	r.renderEventEditKeyLegend()

	return r.terminal.Flush()
}

// renderMonth renders a single month at the specified position
func (r *Renderer) renderMonth(month time.Time, x, y int, selection *models.Selection) error {
	fg, bg := r.terminal.GetDefaultColors()

	// Render month header (month name and year) with color
	monthHeader := fmt.Sprintf("%s %d", calendar.GetMonthName(month), month.Year())
	headerX := x + (r.monthWidth-len(monthHeader))/2

	var headerFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		// Use magenta for month headers in color terminals
		headerFg = termbox.ColorMagenta | termbox.AttrBold
	} else {
		headerFg = termbox.AttrBold
	}
	r.terminal.Print(headerX, y, monthHeader, headerFg, bg)

	// Render day-of-week headers with color
	dayHeaders := calendar.GetDayOfWeekHeaders()
	headerY := y + 2

	var dayHeaderFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		// Use cyan for day-of-week headers in color terminals
		dayHeaderFg = termbox.ColorCyan
	} else {
		dayHeaderFg = fg
	}

	for i, header := range dayHeaders {
		headerX := x + i*3 + 1
		r.terminal.Print(headerX, headerY, header, dayHeaderFg, bg)
	}

	// Render separator line
	separatorY := headerY + 1
	for i := 0; i < r.monthWidth-2; i++ {
		r.terminal.SetCell(x+1+i, separatorY, '-', fg, bg)
	}

	// Get calendar weeks for this month
	weeks := calendar.GetCalendarWeeks(month)

	// Render day grid
	startY := separatorY + 1
	for weekIndex, week := range weeks {
		weekY := startY + weekIndex
		for dayIndex, dayNum := range week {
			dayX := x + dayIndex*3 + 1

			if dayNum == 0 {
				// Empty cell
				r.terminal.Print(dayX, weekY, "  ", fg, bg)
			} else {
				// Create date for this day
				dayDate := time.Date(month.Year(), month.Month(), dayNum, 0, 0, 0, 0, month.Location())

				// Determine display attributes
				dayFg, dayBg, dayText := r.getDayAttributes(dayDate, selection)

				r.terminal.Print(dayX, weekY, dayText, dayFg, dayBg)
			}
		}
	}

	return nil
}

// getDayAttributes determines the display attributes for a day cell
func (r *Renderer) getDayAttributes(date time.Time, selection *models.Selection) (fg, bg termbox.Attribute, text string) {
	dayNum := date.Day()
	text = fmt.Sprintf("%2d", dayNum)

	// Check various states
	isToday := calendar.IsToday(date)
	isSelected := calendar.IsSameDate(date, selection.SelectedDate)
	hasEvents := r.eventManager.HasEventsForDate(date)

	// Default colors
	fg = termbox.ColorDefault
	bg = termbox.ColorDefault

	// Apply color themes based on state (with fallback for monochrome terminals)
	if r.terminal.IsColorSupported() {
		// Color terminal - use color themes
		if isSelected && isToday {
			// Selected + Today: bright cyan background with white text
			fg = termbox.ColorWhite | termbox.AttrBold
			bg = termbox.ColorCyan
		} else if isSelected {
			// Selected: blue background with white text
			fg = termbox.ColorWhite | termbox.AttrBold
			bg = termbox.ColorBlue
		} else if isToday {
			// Today: bright yellow text
			fg = termbox.ColorYellow | termbox.AttrBold
			bg = termbox.ColorDefault
		} else if hasEvents {
			// Days with events: green text
			fg = termbox.ColorGreen
			bg = termbox.ColorDefault
		}
	} else {
		// Monochrome terminal - use attribute-based styling
		if isSelected && isToday {
			// Selected + Today: reverse colors with bold
			fg = termbox.ColorDefault | termbox.AttrBold | termbox.AttrReverse
			bg = termbox.ColorDefault
		} else if isSelected {
			// Selected: reverse colors
			fg = termbox.ColorDefault | termbox.AttrReverse
			bg = termbox.ColorDefault
		} else if isToday {
			// Today: bold
			fg = termbox.ColorDefault | termbox.AttrBold
			bg = termbox.ColorDefault
		}
	}

	// Note: Event indication is now handled purely through color coding
	// No additional visual indicators (bullets, asterisks) are added

	return fg, bg, text
}

// renderSelectedDateEvents renders events for the selected date below the calendar
func (r *Renderer) renderSelectedDateEvents(selectedDate time.Time) {
	fg, bg := r.terminal.GetDefaultColors()

	// Calculate Y position for events section (after calendar, before key legend)
	// Calendar starts at Y=2, month header + day headers + separator + 6 weeks = ~10 lines per month
	eventsStartY := 13

	// Calculate left alignment position to match calendar's left edge
	width, _ := r.terminal.GetSize()
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2
	eventsLeftX := startX + 1 // Align with calendar's leftmost day column

	// Get events for the selected date
	events := r.eventManager.GetEventsForDate(selectedDate)

	// Render section header
	dateStr := calendar.FormatDate(selectedDate)
	headerText := fmt.Sprintf("Events for %s:", dateStr)

	var headerFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		headerFg = termbox.ColorYellow | termbox.AttrBold
	} else {
		headerFg = termbox.AttrBold
	}

	r.terminal.Print(eventsLeftX, eventsStartY, headerText, headerFg, bg)

	// Render events or "no events" message
	if len(events) == 0 {
		var noEventsFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			noEventsFg = termbox.ColorWhite
		} else {
			noEventsFg = fg
		}
		r.terminal.Print(eventsLeftX, eventsStartY+1, "No events scheduled", noEventsFg, bg)
	} else {
		// Show up to 10 events per date
		maxEvents := 10
		if len(events) > maxEvents {
			maxEvents = 10
		}

		for i := 0; i < maxEvents && i < len(events); i++ {
			event := events[i]
			timeStr := event.GetTimeString()
			description := event.Description

			var eventFg termbox.Attribute
			if r.terminal.IsColorSupported() {
				eventFg = termbox.ColorWhite
			} else {
				eventFg = fg
			}

			// Render event as single line
			eventY := eventsStartY + 1 + i
			eventText := fmt.Sprintf("%s - %s", timeStr, description)

			// Calculate available width from left position to right margin
			maxEventWidth := width - eventsLeftX - 4 // Leave some right margin
			if len(eventText) > maxEventWidth {
				eventText = eventText[:maxEventWidth-3] + "..."
			}

			r.terminal.Print(eventsLeftX, eventY, eventText, eventFg, bg)
		}

		// Show "and X more" if there are additional events
		if len(events) > maxEvents {
			moreText := fmt.Sprintf("... and %d more events", len(events)-maxEvents)
			var moreFg termbox.Attribute
			if r.terminal.IsColorSupported() {
				moreFg = termbox.ColorMagenta
			} else {
				moreFg = fg
			}
			r.terminal.Print(eventsLeftX, eventsStartY+1+maxEvents, moreText, moreFg, bg)
		}
	}
}

// renderSelectedDateEventsWithSelection renders events for the selected date with selection highlighting
func (r *Renderer) renderSelectedDateEventsWithSelection(selectedDate time.Time, selectedEventIndex int) {
	fg, bg := r.terminal.GetDefaultColors()

	// Calculate Y position for events section (after calendar, before key legend)
	// Calendar starts at Y=2, month header + day headers + separator + 6 weeks = ~10 lines per month
	eventsStartY := 13

	// Calculate left alignment position to match calendar's left edge
	width, _ := r.terminal.GetSize()
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2
	eventsLeftX := startX + 1 // Align with calendar's leftmost day column

	// Get events for the selected date
	events := r.eventManager.GetEventsForDate(selectedDate)

	// Render section header
	dateStr := calendar.FormatDate(selectedDate)
	headerText := fmt.Sprintf("Events for %s (Use ↑↓ to select, Enter to delete, Esc to cancel):", dateStr)

	var headerFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		headerFg = termbox.ColorYellow | termbox.AttrBold
	} else {
		headerFg = termbox.AttrBold
	}

	r.terminal.Print(eventsLeftX, eventsStartY, headerText, headerFg, bg)

	// Render events or "no events" message
	if len(events) == 0 {
		var noEventsFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			noEventsFg = termbox.ColorWhite
		} else {
			noEventsFg = fg
		}
		r.terminal.Print(eventsLeftX, eventsStartY+1, "No events scheduled", noEventsFg, bg)
	} else {
		// Show up to 10 events per date
		maxEvents := 10
		if len(events) > maxEvents {
			maxEvents = 10
		}

		for i := 0; i < maxEvents && i < len(events); i++ {
			event := events[i]
			timeStr := event.GetTimeString()
			description := event.Description

			// Check if this is the selected event
			isSelected := i == selectedEventIndex

			var eventFg, eventBg termbox.Attribute
			var prefix string

			if isSelected {
				// Selected event: use highlighting
				prefix = "> "
				if r.terminal.IsColorSupported() {
					eventFg = termbox.ColorBlack | termbox.AttrBold
					eventBg = termbox.ColorYellow // Yellow background for selection
				} else {
					eventFg = termbox.ColorDefault | termbox.AttrReverse | termbox.AttrBold
					eventBg = termbox.ColorDefault
				}
			} else {
				// Normal event colors
				prefix = "  "
				eventBg = bg
				if r.terminal.IsColorSupported() {
					eventFg = termbox.ColorWhite
				} else {
					eventFg = fg
				}
			}

			// Render event as single line with selection indicator
			eventY := eventsStartY + 1 + i
			eventText := fmt.Sprintf("%s%s - %s", prefix, timeStr, description)

			// Calculate available width from left position to right margin
			maxEventWidth := width - eventsLeftX - 4 // Leave some right margin
			if len(eventText) > maxEventWidth {
				eventText = eventText[:maxEventWidth-3] + "..."
			}

			r.terminal.Print(eventsLeftX, eventY, eventText, eventFg, eventBg)

			// Fill the rest of the line with the background color for selected events
			if isSelected {
				for x := eventsLeftX + len(eventText); x < width; x++ {
					r.terminal.SetCell(x, eventY, ' ', eventFg, eventBg)
				}
			}
		}

		// Show "and X more" if there are additional events
		if len(events) > maxEvents {
			moreText := fmt.Sprintf("... and %d more events", len(events)-maxEvents)
			var moreFg termbox.Attribute
			if r.terminal.IsColorSupported() {
				moreFg = termbox.ColorMagenta
			} else {
				moreFg = fg
			}
			r.terminal.Print(eventsLeftX, eventsStartY+1+maxEvents, moreText, moreFg, bg)
		}
	}
}

// renderSelectedDateEventsWithEditMode renders events for the selected date with edit mode highlighting
func (r *Renderer) renderSelectedDateEventsWithEditMode(selectedDate time.Time, selectedEventIndex int) {
	fg, bg := r.terminal.GetDefaultColors()

	// Calculate Y position for events section (after calendar, before key legend)
	// Calendar starts at Y=2, month header + day headers + separator + 6 weeks = ~10 lines per month
	eventsStartY := 13

	// Calculate left alignment position to match calendar's left edge
	width, _ := r.terminal.GetSize()
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2
	eventsLeftX := startX + 1 // Align with calendar's leftmost day column

	// Get events for the selected date
	events := r.eventManager.GetEventsForDate(selectedDate)

	// Render section header
	dateStr := calendar.FormatDate(selectedDate)
	headerText := fmt.Sprintf("Events for %s (Use ↑↓ to select, Enter to edit, Esc to cancel):", dateStr)

	var headerFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		headerFg = termbox.ColorYellow | termbox.AttrBold
	} else {
		headerFg = termbox.AttrBold
	}

	r.terminal.Print(eventsLeftX, eventsStartY, headerText, headerFg, bg)

	// Render events or "no events" message
	if len(events) == 0 {
		var noEventsFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			noEventsFg = termbox.ColorWhite
		} else {
			noEventsFg = fg
		}
		r.terminal.Print(eventsLeftX, eventsStartY+1, "No events scheduled", noEventsFg, bg)
	} else {
		// Show up to 10 events per date
		maxEvents := 10
		if len(events) > maxEvents {
			maxEvents = 10
		}

		for i := 0; i < maxEvents && i < len(events); i++ {
			event := events[i]
			timeStr := event.GetTimeString()
			description := event.Description

			// Check if this is the selected event
			isSelected := i == selectedEventIndex

			var eventFg, eventBg termbox.Attribute
			var prefix string

			if isSelected {
				// Selected event: use highlighting
				prefix = "> "
				if r.terminal.IsColorSupported() {
					eventFg = termbox.ColorBlack | termbox.AttrBold
					eventBg = termbox.ColorYellow // Yellow background for selection
				} else {
					eventFg = termbox.ColorDefault | termbox.AttrReverse | termbox.AttrBold
					eventBg = termbox.ColorDefault
				}
			} else {
				// Normal event colors
				prefix = "  "
				eventBg = bg
				if r.terminal.IsColorSupported() {
					eventFg = termbox.ColorWhite
				} else {
					eventFg = fg
				}
			}

			// Render event as single line with selection indicator
			eventY := eventsStartY + 1 + i
			eventText := fmt.Sprintf("%s%s - %s", prefix, timeStr, description)

			// Calculate available width from left position to right margin
			maxEventWidth := width - eventsLeftX - 4 // Leave some right margin
			if len(eventText) > maxEventWidth {
				eventText = eventText[:maxEventWidth-3] + "..."
			}

			r.terminal.Print(eventsLeftX, eventY, eventText, eventFg, eventBg)

			// Fill the rest of the line with the background color for selected events
			if isSelected {
				for x := eventsLeftX + len(eventText); x < width; x++ {
					r.terminal.SetCell(x, eventY, ' ', eventFg, eventBg)
				}
			}
		}

		// Show "and X more" if there are additional events
		if len(events) > maxEvents {
			moreText := fmt.Sprintf("... and %d more events", len(events)-maxEvents)
			var moreFg termbox.Attribute
			if r.terminal.IsColorSupported() {
				moreFg = termbox.ColorMagenta
			} else {
				moreFg = fg
			}
			r.terminal.Print(eventsLeftX, eventsStartY+1+maxEvents, moreText, moreFg, bg)
		}
	}
}

// renderSelectedDateEventsWithAddMode renders events for the selected date with add mode highlighting
func (r *Renderer) renderSelectedDateEventsWithAddMode(selectedDate time.Time) {
	fg, bg := r.terminal.GetDefaultColors()

	// Calculate Y position for events section (after calendar, before key legend)
	// Calendar starts at Y=2, month header + day headers + separator + 6 weeks = ~10 lines per month
	eventsStartY := 13

	// Calculate left alignment position to match calendar's left edge
	width, _ := r.terminal.GetSize()
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2
	eventsLeftX := startX + 1 // Align with calendar's leftmost day column

	// Get events for the selected date
	events := r.eventManager.GetEventsForDate(selectedDate)

	// Render section header
	dateStr := calendar.FormatDate(selectedDate)
	headerText := fmt.Sprintf("Add new event for %s (Enter to add, Esc to cancel):", dateStr)

	var headerFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		headerFg = termbox.ColorYellow | termbox.AttrBold
	} else {
		headerFg = termbox.AttrBold
	}

	r.terminal.Print(eventsLeftX, eventsStartY, headerText, headerFg, bg)

	// First render existing events (up to 9 to leave room for the new event row)
	maxExistingEvents := 9
	if len(events) > maxExistingEvents {
		maxExistingEvents = 9
	}

	for i := 0; i < maxExistingEvents && i < len(events); i++ {
		event := events[i]
		timeStr := event.GetTimeString()
		description := event.Description

		var eventFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			eventFg = termbox.ColorWhite
		} else {
			eventFg = fg
		}

		// Render existing event as single line with normal formatting
		eventY := eventsStartY + 1 + i
		eventText := fmt.Sprintf("  %s - %s", timeStr, description)

		// Calculate available width from left position to right margin
		maxEventWidth := width - eventsLeftX - 4 // Leave some right margin
		if len(eventText) > maxEventWidth {
			eventText = eventText[:maxEventWidth-3] + "..."
		}

		r.terminal.Print(eventsLeftX, eventY, eventText, eventFg, bg)
	}

	// Now render the highlighted empty row for adding new event
	addEventY := eventsStartY + 1 + maxExistingEvents

	var addEventFg, addEventBg termbox.Attribute
	if r.terminal.IsColorSupported() {
		addEventFg = termbox.ColorBlack | termbox.AttrBold
		addEventBg = termbox.ColorYellow // Yellow background for selection
	} else {
		addEventFg = termbox.ColorDefault | termbox.AttrReverse | termbox.AttrBold
		addEventBg = termbox.ColorDefault
	}

	// Render the empty highlighted row for new event
	newEventText := "> [New Event]"
	r.terminal.Print(eventsLeftX, addEventY, newEventText, addEventFg, addEventBg)

	// Fill the rest of the line with the background color
	for x := eventsLeftX + len(newEventText); x < width; x++ {
		r.terminal.SetCell(x, addEventY, ' ', addEventFg, addEventBg)
	}

	// Show "and X more" if there are additional existing events
	if len(events) > maxExistingEvents {
		moreText := fmt.Sprintf("... and %d more existing events", len(events)-maxExistingEvents)
		var moreFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			moreFg = termbox.ColorMagenta
		} else {
			moreFg = fg
		}
		r.terminal.Print(eventsLeftX, addEventY+1, moreText, moreFg, bg)
	}
}

// renderEventSelectionKeyLegend renders the key bindings legend for event selection mode
func (r *Renderer) renderEventSelectionKeyLegend() {
	_, height := r.terminal.GetSize()
	legendY := height - 2

	fg, bg := r.terminal.GetDefaultColors()

	legend := "↑↓: select event  Enter: delete  Esc: cancel"
	r.terminal.PrintCentered(legendY, legend, fg, bg)
}

// renderEventAddKeyLegend renders the key bindings legend for event add mode
func (r *Renderer) renderEventAddKeyLegend() {
	_, height := r.terminal.GetSize()
	legendY := height - 2

	fg, bg := r.terminal.GetDefaultColors()

	legend := "Enter: add event  Esc: cancel"
	r.terminal.PrintCentered(legendY, legend, fg, bg)
}

// renderEventEditKeyLegend renders the key bindings legend for event edit mode
func (r *Renderer) renderEventEditKeyLegend() {
	_, height := r.terminal.GetSize()
	legendY := height - 2

	fg, bg := r.terminal.GetDefaultColors()

	legend := "↑↓: select event  Enter: edit  Esc: cancel"
	r.terminal.PrintCentered(legendY, legend, fg, bg)
}

// renderKeyLegend renders the key bindings legend at the bottom
func (r *Renderer) renderKeyLegend() {
	_, height := r.terminal.GetSize()
	legendY := height - 2

	fg, bg := r.terminal.GetDefaultColors()

	legend := "B/N: month  H/J/K/L: move  Enter: events  A: add  D: delete  E: edit  C: current  F: search  Q: quit"
	r.terminal.PrintCentered(legendY, legend, fg, bg)
}

// RenderEventList renders the event list for a selected date with selection highlighting
func (r *Renderer) RenderEventList(date time.Time, events []models.Event, selectedIndex int) error {
	r.terminal.Clear()

	width, height := r.terminal.GetSize()
	fg, bg := r.terminal.GetDefaultColors()

	// Title with color
	dateStr := calendar.FormatDate(date)
	title := fmt.Sprintf("Events for %s", dateStr)

	var titleFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		// Use yellow for the title in color terminals
		titleFg = termbox.ColorYellow | termbox.AttrBold
	} else {
		titleFg = termbox.AttrBold
	}
	r.terminal.PrintCentered(2, title, titleFg, bg)

	// Draw separator with color
	separatorY := 4
	var separatorFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		separatorFg = termbox.ColorCyan
	} else {
		separatorFg = fg
	}
	for i := 0; i < width; i++ {
		r.terminal.SetCell(i, separatorY, '-', separatorFg, bg)
	}

	startY := 6
	if len(events) == 0 {
		var noEventsFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			noEventsFg = termbox.ColorWhite
		} else {
			noEventsFg = fg
		}
		r.terminal.PrintCentered(startY, "No events scheduled for this date", noEventsFg, bg)
	} else {
		for i, event := range events {
			if startY+i >= height-4 {
				// Too many events to display
				moreText := fmt.Sprintf("... and %d more events", len(events)-i)
				var moreFg termbox.Attribute
				if r.terminal.IsColorSupported() {
					moreFg = termbox.ColorMagenta
				} else {
					moreFg = fg
				}
				r.terminal.PrintCentered(startY+i, moreText, moreFg, bg)
				break
			}

			// Check if this is the selected event
			isSelected := i == selectedIndex

			// Color the time and description differently
			timeStr := event.GetTimeString()
			description := event.Description

			var timeFg, descFg, eventBg termbox.Attribute
			if isSelected {
				// Selected event: use highlighting
				if r.terminal.IsColorSupported() {
					timeFg = termbox.ColorBlack | termbox.AttrBold
					descFg = termbox.ColorBlack
					eventBg = termbox.ColorYellow // Yellow background for selection
				} else {
					timeFg = termbox.ColorDefault | termbox.AttrReverse | termbox.AttrBold
					descFg = termbox.ColorDefault | termbox.AttrReverse
					eventBg = termbox.ColorDefault
				}
			} else {
				// Normal event colors
				eventBg = bg
				if r.terminal.IsColorSupported() {
					timeFg = termbox.ColorGreen | termbox.AttrBold // Green for time
					descFg = termbox.ColorWhite                    // White for description
				} else {
					timeFg = termbox.AttrBold
					descFg = fg
				}
			}

			// Add selection indicator
			var prefix string
			if isSelected {
				prefix = "> "
			} else {
				prefix = "  "
			}

			// Print prefix (selection indicator)
			r.terminal.Print(0, startY+i, prefix, timeFg, eventBg)

			// Print time
			r.terminal.Print(2, startY+i, timeStr, timeFg, eventBg)

			// Print separator
			separator := " - "
			r.terminal.Print(2+len(timeStr), startY+i, separator, timeFg, eventBg)

			// Print description (truncate if too long)
			descriptionText := description
			maxDescWidth := width - 4 - len(timeStr) - len(separator)
			if len(descriptionText) > maxDescWidth {
				descriptionText = descriptionText[:maxDescWidth-3] + "..."
			}
			r.terminal.Print(2+len(timeStr)+len(separator), startY+i, descriptionText, descFg, eventBg)

			// Fill the rest of the line with the background color for selected events
			if isSelected {
				lineLength := 2 + len(timeStr) + len(separator) + len(descriptionText)
				for x := lineLength; x < width; x++ {
					r.terminal.SetCell(x, startY+i, ' ', timeFg, eventBg)
				}
			}
		}
	}

	// Instructions with color
	instrY := height - 3
	var instrFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		instrFg = termbox.ColorCyan
	} else {
		instrFg = fg
	}
	r.terminal.PrintCentered(instrY, "J/K: navigate  A: add event  D: delete event  E: edit event  Esc: back to calendar", instrFg, bg)

	return r.terminal.Flush()
}

// RenderMessage renders a status message at the bottom
func (r *Renderer) RenderMessage(message string, isError bool) {
	_, height := r.terminal.GetSize()
	messageY := height - 1

	var fg termbox.Attribute
	if isError {
		fg = termbox.ColorRed
	} else {
		fg = termbox.ColorGreen
	}

	// Clear the line first
	r.terminal.FillRect(0, messageY, 80, 1, ' ', termbox.ColorDefault, termbox.ColorDefault)

	// Display message (truncate if too long)
	if len(message) > 78 {
		message = message[:75] + "..."
	}

	r.terminal.PrintCentered(messageY, message, fg, termbox.ColorDefault)
}

// RenderInputPrompt renders an input prompt for adding events
func (r *Renderer) RenderInputPrompt(prompt, input string) error {
	_, height := r.terminal.GetSize()
	promptY := height - 4
	inputY := height - 3

	fg, bg := r.terminal.GetDefaultColors()

	// Clear the input area
	r.terminal.FillRect(0, promptY, 80, 3, ' ', fg, bg)

	// Display prompt
	r.terminal.PrintCentered(promptY, prompt, fg, bg)

	// Display input with cursor
	inputText := input + "_"
	r.terminal.PrintCentered(inputY, inputText, fg, bg)

	return r.terminal.Flush()
}

// RenderInlineInput renders input directly on the highlighted event line
func (r *Renderer) RenderInlineInput(x, y int, prompt, input string) error {
	width, _ := r.terminal.GetSize()

	// Use highlighting colors similar to event selection
	var inputFg, inputBg termbox.Attribute
	if r.terminal.IsColorSupported() {
		inputFg = termbox.ColorBlack | termbox.AttrBold
		inputBg = termbox.ColorYellow // Yellow background for input
	} else {
		inputFg = termbox.ColorDefault | termbox.AttrReverse | termbox.AttrBold
		inputBg = termbox.ColorDefault
	}

	// Clear the entire line first
	for i := x; i < width; i++ {
		r.terminal.SetCell(i, y, ' ', inputFg, inputBg)
	}

	// Create the display text with cursor
	displayText := fmt.Sprintf("> %s %s_", prompt, input)

	// Truncate if too long
	maxWidth := width - x - 2
	if len(displayText) > maxWidth {
		displayText = displayText[:maxWidth-3] + "..."
	}

	// Display the input line
	r.terminal.Print(x, y, displayText, inputFg, inputBg)

	return r.terminal.Flush()
}

// RenderCalendarWithSearch renders the calendar with search results
func (r *Renderer) RenderCalendarWithSearch(cal *models.Calendar, selection *models.Selection, query string, results []models.Event, resultDates []string, selectedIndex int) error {
	r.terminal.Clear()

	// Get terminal size
	width, height := r.terminal.GetSize()
	if width < 80 || height < 24 {
		r.terminal.PrintCentered(height/2, "Terminal too small! Minimum 80x24 required.",
			termbox.ColorRed, termbox.ColorDefault)
		return r.terminal.Flush()
	}

	// Calculate starting positions for three months
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2

	prevMonth := cal.GetPreviousMonth()
	currentMonth := cal.CurrentMonth
	nextMonth := cal.GetNextMonth()

	months := []time.Time{prevMonth, currentMonth, nextMonth}

	// Render each month
	for i, month := range months {
		x := startX + i*(r.monthWidth+r.monthSpacing)
		err := r.renderMonth(month, x, 2, selection)
		if err != nil {
			return err
		}
	}

	// Render search results under the calendar
	r.renderSearchResults(query, results, selectedIndex)

	// Render search key legend
	r.renderSearchKeyLegend()

	return r.terminal.Flush()
}

// renderSearchResults renders search results grouped by date under the calendar
func (r *Renderer) renderSearchResults(query string, results []models.Event, selectedIndex int) {
	fg, bg := r.terminal.GetDefaultColors()

	// Calculate Y position for search results section
	searchStartY := 13

	// Calculate left alignment position to match calendar's left edge
	width, height := r.terminal.GetSize()
	totalWidth := 3*r.monthWidth + 2*r.monthSpacing
	startX := (width - totalWidth) / 2
	searchLeftX := startX + 1

	// Render section header
	headerText := fmt.Sprintf("Search results for \"%s\":", query)
	if query == "" {
		headerText = "Search results:"
	}

	var headerFg termbox.Attribute
	if r.terminal.IsColorSupported() {
		headerFg = termbox.ColorYellow | termbox.AttrBold
	} else {
		headerFg = termbox.AttrBold
	}

	r.terminal.Print(searchLeftX, searchStartY, headerText, headerFg, bg)

	// Render search results
	if len(results) == 0 {
		var noResultsFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			noResultsFg = termbox.ColorWhite
		} else {
			noResultsFg = fg
		}
		r.terminal.Print(searchLeftX, searchStartY+1, "No events found matching your search", noResultsFg, bg)
	} else {
		// Group results by date and render
		currentY := searchStartY + 1
		currentDate := ""

		for i, event := range results {
			if currentY >= height-4 {
				// Too many results to display
				moreText := fmt.Sprintf("... and %d more results", len(results)-i)
				var moreFg termbox.Attribute
				if r.terminal.IsColorSupported() {
					moreFg = termbox.ColorMagenta
				} else {
					moreFg = fg
				}
				r.terminal.Print(searchLeftX, currentY, moreText, moreFg, bg)
				break
			}

			eventDateStr := event.Date.Format("2006-01-02")

			// Show date header if this is a new date
			if eventDateStr != currentDate {
				currentDate = eventDateStr
				if currentY > searchStartY+1 {
					currentY++ // Add space between date groups
				}

				// Format date header
				dateHeader := event.Date.Format("Monday, January 2, 2006")
				var dateFg termbox.Attribute
				if r.terminal.IsColorSupported() {
					dateFg = termbox.ColorCyan | termbox.AttrBold
				} else {
					dateFg = termbox.AttrBold
				}
				r.terminal.Print(searchLeftX, currentY, dateHeader, dateFg, bg)
				currentY++
			}

			// Check if this is the selected result
			isSelected := i == selectedIndex

			var eventFg, eventBg termbox.Attribute
			var prefix string

			if isSelected {
				// Selected result: use highlighting
				prefix = "  > "
				if r.terminal.IsColorSupported() {
					eventFg = termbox.ColorBlack | termbox.AttrBold
					eventBg = termbox.ColorYellow // Yellow background for selection
				} else {
					eventFg = termbox.ColorDefault | termbox.AttrReverse | termbox.AttrBold
					eventBg = termbox.ColorDefault
				}
			} else {
				// Normal result colors
				prefix = "    "
				eventBg = bg
				if r.terminal.IsColorSupported() {
					eventFg = termbox.ColorWhite
				} else {
					eventFg = fg
				}
			}

			// Render event as single line
			timeStr := event.GetTimeString()
			description := event.Description
			eventText := fmt.Sprintf("%s%s - %s", prefix, timeStr, description)

			// Calculate available width from left position to right margin
			maxEventWidth := width - searchLeftX - 4 // Leave some right margin
			if len(eventText) > maxEventWidth {
				eventText = eventText[:maxEventWidth-3] + "..."
			}

			r.terminal.Print(searchLeftX, currentY, eventText, eventFg, eventBg)

			// Fill the rest of the line with the background color for selected results
			if isSelected {
				for x := searchLeftX + len(eventText); x < width; x++ {
					r.terminal.SetCell(x, currentY, ' ', eventFg, eventBg)
				}
			}

			currentY++
		}
	}
}

// renderSearchKeyLegend renders the key bindings legend for search mode
func (r *Renderer) renderSearchKeyLegend() {
	_, height := r.terminal.GetSize()
	legendY := height - 2

	fg, bg := r.terminal.GetDefaultColors()

	legend := "↑↓: navigate results  Enter: go to date  Esc: back to calendar  F: search"
	r.terminal.PrintCentered(legendY, legend, fg, bg)
}
