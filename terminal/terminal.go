package terminal

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Terminal handles low-level terminal operations
type Terminal struct {
	width  int
	height int
}

// NewTerminal creates a new terminal handler
func NewTerminal() *Terminal {
	return &Terminal{}
}

// Initialize initializes the terminal for raw input mode
func (t *Terminal) Initialize() error {
	err := termbox.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize terminal: %v", err)
	}

	// Set input mode to ESC and Alt keys
	termbox.SetInputMode(termbox.InputEsc)

	// Update terminal dimensions
	t.updateSize()

	return nil
}

// Close cleans up and restores the terminal
func (t *Terminal) Close() {
	termbox.Close()
}

// Clear clears the entire screen
func (t *Terminal) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// Flush flushes all changes to the terminal
func (t *Terminal) Flush() error {
	return termbox.Flush()
}

// GetSize returns the current terminal dimensions
func (t *Terminal) GetSize() (width, height int) {
	return t.width, t.height
}

// updateSize updates the stored terminal dimensions
func (t *Terminal) updateSize() {
	t.width, t.height = termbox.Size()
}

// CheckSize checks if terminal is large enough (minimum 80x24)
func (t *Terminal) CheckSize() bool {
	t.updateSize()
	return t.width >= 80 && t.height >= 24
}

// SetCell sets a character at the specified position with colors
func (t *Terminal) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, ch, fg, bg)
}

// Print prints a string at the specified position with colors
func (t *Terminal) Print(x, y int, text string, fg, bg termbox.Attribute) {
	for i, ch := range text {
		if x+i < t.width {
			termbox.SetCell(x+i, y, ch, fg, bg)
		}
	}
}

// PrintCentered prints text centered horizontally at the specified y position
func (t *Terminal) PrintCentered(y int, text string, fg, bg termbox.Attribute) {
	x := (t.width - len(text)) / 2
	if x < 0 {
		x = 0
	}
	t.Print(x, y, text, fg, bg)
}

// PrintRight prints text right-aligned at the specified y position
func (t *Terminal) PrintRight(y int, text string, fg, bg termbox.Attribute) {
	x := t.width - len(text)
	if x < 0 {
		x = 0
	}
	t.Print(x, y, text, fg, bg)
}

// DrawBox draws a simple box using ASCII characters
func (t *Terminal) DrawBox(x, y, width, height int, fg, bg termbox.Attribute) {
	// Top and bottom borders
	for i := 0; i < width; i++ {
		if x+i < t.width {
			if y >= 0 && y < t.height {
				termbox.SetCell(x+i, y, '-', fg, bg)
			}
			if y+height-1 >= 0 && y+height-1 < t.height {
				termbox.SetCell(x+i, y+height-1, '-', fg, bg)
			}
		}
	}

	// Left and right borders
	for i := 0; i < height; i++ {
		if y+i >= 0 && y+i < t.height {
			if x >= 0 && x < t.width {
				termbox.SetCell(x, y+i, '|', fg, bg)
			}
			if x+width-1 >= 0 && x+width-1 < t.width {
				termbox.SetCell(x+width-1, y+i, '|', fg, bg)
			}
		}
	}

	// Corners
	if x >= 0 && x < t.width && y >= 0 && y < t.height {
		termbox.SetCell(x, y, '+', fg, bg) // Top-left
	}
	if x+width-1 >= 0 && x+width-1 < t.width && y >= 0 && y < t.height {
		termbox.SetCell(x+width-1, y, '+', fg, bg) // Top-right
	}
	if x >= 0 && x < t.width && y+height-1 >= 0 && y+height-1 < t.height {
		termbox.SetCell(x, y+height-1, '+', fg, bg) // Bottom-left
	}
	if x+width-1 >= 0 && x+width-1 < t.width && y+height-1 >= 0 && y+height-1 < t.height {
		termbox.SetCell(x+width-1, y+height-1, '+', fg, bg) // Bottom-right
	}
}

// FillRect fills a rectangle with the specified character and colors
func (t *Terminal) FillRect(x, y, width, height int, ch rune, fg, bg termbox.Attribute) {
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if x+col >= 0 && x+col < t.width && y+row >= 0 && y+row < t.height {
				termbox.SetCell(x+col, y+row, ch, fg, bg)
			}
		}
	}
}

// PollEvent waits for and returns the next keyboard event
func (t *Terminal) PollEvent() termbox.Event {
	return termbox.PollEvent()
}

// IsColorSupported checks if the terminal supports colors
func (t *Terminal) IsColorSupported() bool {
	// termbox-go handles color detection internally
	// For our purposes, we'll assume basic color support is available
	return true
}

// GetDefaultColors returns the default foreground and background colors
func (t *Terminal) GetDefaultColors() (fg, bg termbox.Attribute) {
	return termbox.ColorDefault, termbox.ColorDefault
}

// GetColors returns commonly used color combinations
func (t *Terminal) GetColors() map[string]termbox.Attribute {
	return map[string]termbox.Attribute{
		"default":   termbox.ColorDefault,
		"black":     termbox.ColorBlack,
		"red":       termbox.ColorRed,
		"green":     termbox.ColorGreen,
		"yellow":    termbox.ColorYellow,
		"blue":      termbox.ColorBlue,
		"magenta":   termbox.ColorMagenta,
		"cyan":      termbox.ColorCyan,
		"white":     termbox.ColorWhite,
		"bold":      termbox.AttrBold,
		"underline": termbox.AttrUnderline,
		"reverse":   termbox.AttrReverse,
	}
}
