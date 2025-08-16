package terminal

import (
	"testing"

	"github.com/nsf/termbox-go"
)

func TestProcessKeyEvent(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	tests := []struct {
		name     string
		event    termbox.Event
		expected KeyAction
	}{
		// Character keys - lowercase
		{"q key", termbox.Event{Type: termbox.EventKey, Ch: 'q'}, ActionQuit},
		{"b key", termbox.Event{Type: termbox.EventKey, Ch: 'b'}, ActionMonthPrev},
		{"n key", termbox.Event{Type: termbox.EventKey, Ch: 'n'}, ActionMonthNext},
		{"h key", termbox.Event{Type: termbox.EventKey, Ch: 'h'}, ActionMoveLeft},
		{"l key", termbox.Event{Type: termbox.EventKey, Ch: 'l'}, ActionMoveRight},
		{"k key", termbox.Event{Type: termbox.EventKey, Ch: 'k'}, ActionMoveUp},
		{"j key", termbox.Event{Type: termbox.EventKey, Ch: 'j'}, ActionMoveDown},
		{"a key", termbox.Event{Type: termbox.EventKey, Ch: 'a'}, ActionAddEvent},

		// Character keys - uppercase
		{"Q key", termbox.Event{Type: termbox.EventKey, Ch: 'Q'}, ActionQuit},
		{"B key", termbox.Event{Type: termbox.EventKey, Ch: 'B'}, ActionMonthPrev},
		{"N key", termbox.Event{Type: termbox.EventKey, Ch: 'N'}, ActionMonthNext},
		{"H key", termbox.Event{Type: termbox.EventKey, Ch: 'H'}, ActionMoveLeft},
		{"L key", termbox.Event{Type: termbox.EventKey, Ch: 'L'}, ActionMoveRight},
		{"K key", termbox.Event{Type: termbox.EventKey, Ch: 'K'}, ActionMoveUp},
		{"J key", termbox.Event{Type: termbox.EventKey, Ch: 'J'}, ActionMoveDown},
		{"A key", termbox.Event{Type: termbox.EventKey, Ch: 'A'}, ActionAddEvent},

		// Special keys
		{"Escape key", termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc}, ActionBack},
		{"Enter key", termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEnter}, ActionShowEvents},
		{"Space key", termbox.Event{Type: termbox.EventKey, Key: termbox.KeySpace}, ActionNone},
		{"Ctrl+C", termbox.Event{Type: termbox.EventKey, Key: termbox.KeyCtrlC}, ActionQuit},

		// Invalid/unrecognized keys
		{"x key", termbox.Event{Type: termbox.EventKey, Ch: 'x'}, ActionNone},
		{"1 key", termbox.Event{Type: termbox.EventKey, Ch: '1'}, ActionNone},
		{"@ key", termbox.Event{Type: termbox.EventKey, Ch: '@'}, ActionNone},

		// Non-key events
		{"Mouse event", termbox.Event{Type: termbox.EventMouse}, ActionNone},
		{"Resize event", termbox.Event{Type: termbox.EventResize}, ActionNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ih.ProcessKeyEvent(tt.event)
			if result != tt.expected {
				t.Errorf("ProcessKeyEvent(%v) = %v, want %v", tt.event, result, tt.expected)
			}
		})
	}
}

func TestGetKeyDescription(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	tests := []struct {
		action      KeyAction
		description string
	}{
		{ActionQuit, "Quit application"},
		{ActionMonthPrev, "Previous month"},
		{ActionMonthNext, "Next month"},
		{ActionMoveLeft, "Move selection left"},
		{ActionMoveRight, "Move selection right"},
		{ActionMoveUp, "Move selection up"},
		{ActionMoveDown, "Move selection down"},
		{ActionShowEvents, "Show events for selected date"},
		{ActionAddEvent, "Add new event"},
		{ActionBack, "Back to previous view"},
		{ActionNone, "Unknown action"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := ih.GetKeyDescription(tt.action)
			if result != tt.description {
				t.Errorf("GetKeyDescription(%v) = %s, want %s", tt.action, result, tt.description)
			}
		})
	}
}

func TestIsValidKey(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	tests := []struct {
		name     string
		key      rune
		expected bool
	}{
		// Valid keys - lowercase
		{"b", 'b', true},
		{"n", 'n', true},
		{"h", 'h', true},
		{"j", 'j', true},
		{"k", 'k', true},
		{"l", 'l', true},
		{"a", 'a', true},
		{"q", 'q', true},

		// Valid keys - uppercase
		{"B", 'B', true},
		{"N", 'N', true},
		{"H", 'H', true},
		{"J", 'J', true},
		{"K", 'K', true},
		{"L", 'L', true},
		{"A", 'A', true},
		{"Q", 'Q', true},

		// Invalid keys
		{"x", 'x', false},
		{"1", '1', false},
		{"@", '@', false},
		{" ", ' ', false},
		{"z", 'z', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ih.IsValidKey(tt.key)
			if result != tt.expected {
				t.Errorf("IsValidKey(%c) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

func TestGetKeyMappings(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	mappings := ih.GetKeyMappings()

	expectedMappings := map[string]string{
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

	if len(mappings) != len(expectedMappings) {
		t.Errorf("Expected %d key mappings, got %d", len(expectedMappings), len(mappings))
	}

	for key, expectedDesc := range expectedMappings {
		if actualDesc, exists := mappings[key]; !exists {
			t.Errorf("Missing key mapping for %s", key)
		} else if actualDesc != expectedDesc {
			t.Errorf("Key mapping for %s: got %s, want %s", key, actualDesc, expectedDesc)
		}
	}
}

// Test key action constants
func TestKeyActionConstants(t *testing.T) {
	// Verify that all key actions have unique values
	actions := []KeyAction{
		ActionNone,
		ActionQuit,
		ActionMonthPrev,
		ActionMonthNext,
		ActionMoveLeft,
		ActionMoveRight,
		ActionMoveUp,
		ActionMoveDown,
		ActionShowEvents,
		ActionAddEvent,
		ActionBack,
	}

	// Check that each action has a unique value
	seen := make(map[KeyAction]bool)
	for _, action := range actions {
		if seen[action] {
			t.Errorf("Duplicate action value: %v", action)
		}
		seen[action] = true
	}

	// Check that ActionNone is 0 (the default/zero value)
	if ActionNone != 0 {
		t.Errorf("ActionNone should be 0, got %v", ActionNone)
	}
}

// Test case-insensitive key processing
func TestCaseInsensitiveProcessing(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	// Test that uppercase and lowercase versions of keys produce the same action
	testKeys := []rune{'q', 'b', 'n', 'h', 'j', 'k', 'l', 'a'}

	for _, key := range testKeys {
		lowerEvent := termbox.Event{Type: termbox.EventKey, Ch: key}
		upperEvent := termbox.Event{Type: termbox.EventKey, Ch: key - 32} // Convert to uppercase

		lowerAction := ih.ProcessKeyEvent(lowerEvent)
		upperAction := ih.ProcessKeyEvent(upperEvent)

		if lowerAction != upperAction {
			t.Errorf("Case insensitive processing failed for key %c: lowercase=%v, uppercase=%v",
				key, lowerAction, upperAction)
		}

		// Both should not be ActionNone (unless it's an unrecognized key)
		if lowerAction == ActionNone {
			t.Errorf("Key %c should be recognized, got ActionNone", key)
		}
	}
}

// Test special key handling
func TestSpecialKeyHandling(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	// Test events with Key field set but Ch field as 0
	specialKeys := []struct {
		key      termbox.Key
		expected KeyAction
	}{
		{termbox.KeyEsc, ActionBack},
		{termbox.KeyEnter, ActionShowEvents},
		{termbox.KeySpace, ActionNone},
		{termbox.KeyCtrlC, ActionQuit},
	}

	for _, tt := range specialKeys {
		event := termbox.Event{Type: termbox.EventKey, Key: tt.key, Ch: 0}
		result := ih.ProcessKeyEvent(event)
		if result != tt.expected {
			t.Errorf("Special key %v: got %v, want %v", tt.key, result, tt.expected)
		}
	}
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	// Test event with both Key and Ch set (Key field takes priority in the implementation)
	event := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc, Ch: 'q'}
	result := ih.ProcessKeyEvent(event)
	// Should process Escape key, not 'q' character (Key field has priority)
	if result != ActionBack {
		t.Errorf("Event with both Key and Ch: got %v, want %v", result, ActionBack)
	}

	// Test event with Ch = 0 (should check Key field)
	event2 := termbox.Event{Type: termbox.EventKey, Key: termbox.KeyEsc, Ch: 0}
	result2 := ih.ProcessKeyEvent(event2)
	if result2 != ActionBack {
		t.Errorf("Event with Ch=0: got %v, want %v", result2, ActionBack)
	}

	// Test invalid event type
	event3 := termbox.Event{Type: termbox.EventError}
	result3 := ih.ProcessKeyEvent(event3)
	if result3 != ActionNone {
		t.Errorf("Invalid event type: got %v, want %v", result3, ActionNone)
	}
}

// Test that key descriptions are comprehensive
func TestKeyDescriptionsComprehensive(t *testing.T) {
	terminal := NewTerminal()
	ih := NewInputHandler(terminal)

	// Test all defined actions have descriptions
	actions := []KeyAction{
		ActionNone,
		ActionQuit,
		ActionMonthPrev,
		ActionMonthNext,
		ActionMoveLeft,
		ActionMoveRight,
		ActionMoveUp,
		ActionMoveDown,
		ActionShowEvents,
		ActionAddEvent,
		ActionBack,
	}

	for _, action := range actions {
		description := ih.GetKeyDescription(action)
		if description == "" {
			t.Errorf("Action %v has empty description", action)
		}
	}

	// Test that undefined action returns default description
	undefinedAction := KeyAction(999)
	description := ih.GetKeyDescription(undefinedAction)
	if description != "Unknown action" {
		t.Errorf("Undefined action should return 'Unknown action', got '%s'", description)
	}
}
