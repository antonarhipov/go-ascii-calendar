package terminal

import (
	"strings"

	"github.com/nsf/termbox-go"
)

// InputHandler handles keyboard input processing
type InputHandler struct {
	terminal *Terminal
}

// NewInputHandler creates a new input handler
func NewInputHandler(terminal *Terminal) *InputHandler {
	return &InputHandler{
		terminal: terminal,
	}
}

// KeyAction represents different types of actions that can be triggered by keys
type KeyAction int

const (
	ActionNone KeyAction = iota
	ActionQuit
	ActionMonthPrev
	ActionMonthNext
	ActionMoveLeft
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionShowEvents
	ActionAddEvent
	ActionDeleteEvent
	ActionEditEvent
	ActionBack
	ActionResetCurrent
)

// ProcessKeyEvent processes a keyboard event and returns the corresponding action
func (ih *InputHandler) ProcessKeyEvent(event termbox.Event) KeyAction {
	if event.Type != termbox.EventKey {
		return ActionNone
	}

	// Handle special keys first
	switch event.Key {
	case termbox.KeyEsc:
		return ActionBack
	case termbox.KeyEnter:
		return ActionShowEvents
	case termbox.KeySpace:
		return ActionNone // Ignore space
	case termbox.KeyCtrlC:
		return ActionQuit
	case termbox.KeyArrowLeft:
		return ActionMoveLeft
	case termbox.KeyArrowRight:
		return ActionMoveRight
	case termbox.KeyArrowUp:
		return ActionMoveUp
	case termbox.KeyArrowDown:
		return ActionMoveDown
	}

	// Handle character keys (convert to lowercase for consistent processing)
	ch := event.Ch
	if ch == 0 {
		return ActionNone
	}

	// Convert to lowercase for case-insensitive processing
	lowerCh := strings.ToLower(string(ch))[0]

	switch lowerCh {
	case 'q':
		return ActionQuit
	case 'b':
		return ActionMonthPrev
	case 'n':
		return ActionMonthNext
	case 'h':
		return ActionMoveLeft
	case 'l':
		return ActionMoveRight
	case 'k':
		return ActionMoveUp
	case 'j':
		return ActionMoveDown
	case 'a':
		return ActionAddEvent
	case 'd':
		return ActionDeleteEvent
	case 'e':
		return ActionEditEvent
	case 'c':
		return ActionResetCurrent
	default:
		// Unrecognized key - could show a brief message
		return ActionNone
	}
}

// GetKeyDescription returns a human-readable description of the key action
func (ih *InputHandler) GetKeyDescription(action KeyAction) string {
	switch action {
	case ActionQuit:
		return "Quit application"
	case ActionMonthPrev:
		return "Previous month"
	case ActionMonthNext:
		return "Next month"
	case ActionMoveLeft:
		return "Move selection left"
	case ActionMoveRight:
		return "Move selection right"
	case ActionMoveUp:
		return "Move selection up"
	case ActionMoveDown:
		return "Move selection down"
	case ActionShowEvents:
		return "Show events for selected date"
	case ActionAddEvent:
		return "Add new event"
	case ActionDeleteEvent:
		return "Delete event"
	case ActionEditEvent:
		return "Edit event"
	case ActionBack:
		return "Back to previous view"
	case ActionResetCurrent:
		return "Reset to current month/day"
	default:
		return "Unknown action"
	}
}

// WaitForKey waits for a key press and returns the event
func (ih *InputHandler) WaitForKey() termbox.Event {
	return ih.terminal.PollEvent()
}

// GetTextInput handles text input for prompts (like adding events)
func (ih *InputHandler) GetTextInput(maxLength int) (string, bool) {
	var input strings.Builder

	for {
		event := ih.terminal.PollEvent()

		if event.Type != termbox.EventKey {
			continue
		}

		switch event.Key {
		case termbox.KeyEsc:
			return "", false // User cancelled

		case termbox.KeyEnter:
			result := input.String()
			return result, true // User confirmed

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if input.Len() > 0 {
				// Remove last character
				str := input.String()
				input.Reset()
				if len(str) > 0 {
					input.WriteString(str[:len(str)-1])
				}
			}

		case termbox.KeySpace:
			if input.Len() < maxLength {
				input.WriteRune(' ')
			}

		default:
			// Handle printable characters
			if event.Ch != 0 && input.Len() < maxLength {
				// Allow printable ASCII characters
				if event.Ch >= 32 && event.Ch <= 126 {
					input.WriteRune(event.Ch)
				}
			}
		}

		// Return current input for live display (but continue loop)
		// This is used by the caller to update the display in real-time
	}
}

// GetTextInputWithPrompt handles text input with live prompt updating
func (ih *InputHandler) GetTextInputWithPrompt(prompt string, maxLength int, renderer *Renderer) (string, bool) {
	var input strings.Builder

	for {
		// Update display with current input
		renderer.RenderInputPrompt(prompt, input.String())

		event := ih.terminal.PollEvent()

		if event.Type != termbox.EventKey {
			continue
		}

		switch event.Key {
		case termbox.KeyEsc:
			return "", false // User cancelled

		case termbox.KeyEnter:
			result := strings.TrimSpace(input.String())
			return result, true // User confirmed

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if input.Len() > 0 {
				// Remove last character
				str := input.String()
				input.Reset()
				if len(str) > 0 {
					input.WriteString(str[:len(str)-1])
				}
			}

		case termbox.KeySpace:
			if input.Len() < maxLength {
				input.WriteRune(' ')
			}

		default:
			// Handle printable characters
			if event.Ch != 0 && input.Len() < maxLength {
				// Allow printable ASCII characters
				if event.Ch >= 32 && event.Ch <= 126 {
					input.WriteRune(event.Ch)
				}
			}
		}
	}
}

// IsValidKey checks if a character is a valid key for the application
func (ih *InputHandler) IsValidKey(ch rune) bool {
	validKeys := "bBnNhHjJkKlLaAqQ"
	return strings.ContainsRune(validKeys, ch)
}

// GetKeyMappings returns all valid key mappings for display
func (ih *InputHandler) GetKeyMappings() map[string]string {
	return map[string]string{
		"B/b":   "Previous month",
		"N/n":   "Next month",
		"H/h":   "Move left",
		"L/l":   "Move right",
		"K/k":   "Move up",
		"J/j":   "Move down",
		"Enter": "Show events",
		"A/a":   "Add event",
		"Q/q":   "Quit",
		"Esc":   "Back",
	}
}
