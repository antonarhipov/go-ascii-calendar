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

	// Add event indicator
	if hasEvents {
		if dayNum < 10 {
			text = fmt.Sprintf(" %d•", dayNum) // Use bullet instead of asterisk
		} else {
			text = fmt.Sprintf("%d•", dayNum) // Use bullet instead of asterisk
		}
	}

	return fg, bg, text
}

// renderSelectedDateEvents renders events for the selected date below the calendar
func (r *Renderer) renderSelectedDateEvents(selectedDate time.Time) {
	fg, bg := r.terminal.GetDefaultColors()

	// Calculate Y position for events section (after calendar, before key legend)
	// Calendar starts at Y=2, month header + day headers + separator + 6 weeks = ~10 lines per month
	eventsStartY := 13

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

	r.terminal.PrintCentered(eventsStartY, headerText, headerFg, bg)

	// Render events or "no events" message
	if len(events) == 0 {
		var noEventsFg termbox.Attribute
		if r.terminal.IsColorSupported() {
			noEventsFg = termbox.ColorWhite
		} else {
			noEventsFg = fg
		}
		r.terminal.PrintCentered(eventsStartY+1, "No events scheduled", noEventsFg, bg)
	} else {
		// Show up to 3 events to keep it compact
		maxEvents := 3
		if len(events) > maxEvents {
			maxEvents = 3
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

			// Truncate if too long (leave space for margins)
			maxEventWidth := 76
			if len(eventText) > maxEventWidth {
				eventText = eventText[:maxEventWidth-3] + "..."
			}

			r.terminal.PrintCentered(eventY, eventText, eventFg, bg)
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
			r.terminal.PrintCentered(eventsStartY+1+maxEvents, moreText, moreFg, bg)
		}
	}
}

// renderKeyLegend renders the key bindings legend at the bottom
func (r *Renderer) renderKeyLegend() {
	_, height := r.terminal.GetSize()
	legendY := height - 2

	fg, bg := r.terminal.GetDefaultColors()

	legend := "B/N: month  H/J/K/L: move  Enter: events  A: add  C: current  Q: quit"
	r.terminal.PrintCentered(legendY, legend, fg, bg)
}

// RenderEventList renders the event list for a selected date
func (r *Renderer) RenderEventList(date time.Time, events []models.Event) error {
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

			// Color the time and description differently
			timeStr := event.GetTimeString()
			description := event.Description

			var timeFg, descFg termbox.Attribute
			if r.terminal.IsColorSupported() {
				timeFg = termbox.ColorGreen | termbox.AttrBold // Green for time
				descFg = termbox.ColorWhite                    // White for description
			} else {
				timeFg = termbox.AttrBold
				descFg = fg
			}

			// Print time
			r.terminal.Print(2, startY+i, timeStr, timeFg, bg)

			// Print separator
			separator := " - "
			r.terminal.Print(2+len(timeStr), startY+i, separator, fg, bg)

			// Print description (truncate if too long)
			descriptionText := description
			maxDescWidth := width - 4 - len(timeStr) - len(separator)
			if len(descriptionText) > maxDescWidth {
				descriptionText = descriptionText[:maxDescWidth-3] + "..."
			}
			r.terminal.Print(2+len(timeStr)+len(separator), startY+i, descriptionText, descFg, bg)
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
	r.terminal.PrintCentered(instrY, "A: add event  Esc: back to calendar", instrFg, bg)

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
