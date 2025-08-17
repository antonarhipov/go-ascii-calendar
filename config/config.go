package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nsf/termbox-go"
)

// WeekStartDay represents the first day of the week
type WeekStartDay int

const (
	StartSunday WeekStartDay = iota // 0 = Sunday (default)
	StartMonday                     // 1 = Monday
)

// ColorTheme defines colors for all UI elements
type ColorTheme struct {
	// Month headers (e.g., "August 2025")
	MonthHeaderFg string `json:"month_header_fg"`
	MonthHeaderBg string `json:"month_header_bg"`

	// Day-of-week headers (e.g., "Su Mo Tu")
	DayHeaderFg string `json:"day_header_fg"`
	DayHeaderBg string `json:"day_header_bg"`

	// Calendar day cells
	RegularDayFg string `json:"regular_day_fg"`
	RegularDayBg string `json:"regular_day_bg"`

	// Today's date
	TodayFg string `json:"today_fg"`
	TodayBg string `json:"today_bg"`

	// Selected date
	SelectedFg string `json:"selected_fg"`
	SelectedBg string `json:"selected_bg"`

	// Selected date that is also today
	SelectedTodayFg string `json:"selected_today_fg"`
	SelectedTodayBg string `json:"selected_today_bg"`

	// Days with events
	EventDayFg string `json:"event_day_fg"`
	EventDayBg string `json:"event_day_bg"`

	// Event list section header
	EventHeaderFg string `json:"event_header_fg"`
	EventHeaderBg string `json:"event_header_bg"`

	// Event list items
	EventTextFg string `json:"event_text_fg"`
	EventTextBg string `json:"event_text_bg"`

	// Selected event in event list
	SelectedEventFg string `json:"selected_event_fg"`
	SelectedEventBg string `json:"selected_event_bg"`

	// "No events" message
	NoEventsFg string `json:"no_events_fg"`
	NoEventsBg string `json:"no_events_bg"`

	// "More events" indicator
	MoreEventsFg string `json:"more_events_fg"`
	MoreEventsBg string `json:"more_events_bg"`

	// Error messages
	ErrorFg string `json:"error_fg"`
	ErrorBg string `json:"error_bg"`

	// Success messages
	SuccessFg string `json:"success_fg"`
	SuccessBg string `json:"success_bg"`

	// Input prompts and fields
	InputFg string `json:"input_fg"`
	InputBg string `json:"input_bg"`

	// Search results
	SearchResultFg string `json:"search_result_fg"`
	SearchResultBg string `json:"search_result_bg"`

	// Key legend/instructions
	InstructionsFg string `json:"instructions_fg"`
	InstructionsBg string `json:"instructions_bg"`
}

// Predefined color themes
var (
	// DefaultTheme matches the current hardcoded colors in the renderer
	DefaultTheme = ColorTheme{
		MonthHeaderFg:   "magenta|bold",
		MonthHeaderBg:   "default",
		DayHeaderFg:     "cyan",
		DayHeaderBg:     "default",
		RegularDayFg:    "default",
		RegularDayBg:    "default",
		TodayFg:         "yellow|bold",
		TodayBg:         "default",
		SelectedFg:      "white|bold",
		SelectedBg:      "blue",
		SelectedTodayFg: "white|bold",
		SelectedTodayBg: "cyan",
		EventDayFg:      "green",
		EventDayBg:      "default",
		EventHeaderFg:   "yellow|bold",
		EventHeaderBg:   "default",
		EventTextFg:     "white",
		EventTextBg:     "default",
		SelectedEventFg: "black|bold",
		SelectedEventBg: "yellow",
		NoEventsFg:      "white",
		NoEventsBg:      "default",
		MoreEventsFg:    "magenta",
		MoreEventsBg:    "default",
		ErrorFg:         "red",
		ErrorBg:         "default",
		SuccessFg:       "green",
		SuccessBg:       "default",
		InputFg:         "black|bold",
		InputBg:         "yellow",
		SearchResultFg:  "white",
		SearchResultBg:  "default",
		InstructionsFg:  "cyan",
		InstructionsBg:  "default",
	}

	// DarkTheme provides better contrast for dark terminals
	DarkTheme = ColorTheme{
		MonthHeaderFg:   "bright_magenta|bold",
		MonthHeaderBg:   "default",
		DayHeaderFg:     "bright_cyan",
		DayHeaderBg:     "default",
		RegularDayFg:    "white",
		RegularDayBg:    "default",
		TodayFg:         "bright_yellow|bold",
		TodayBg:         "default",
		SelectedFg:      "black|bold",
		SelectedBg:      "bright_blue",
		SelectedTodayFg: "black|bold",
		SelectedTodayBg: "bright_cyan",
		EventDayFg:      "bright_green",
		EventDayBg:      "default",
		EventHeaderFg:   "bright_yellow|bold",
		EventHeaderBg:   "default",
		EventTextFg:     "bright_white",
		EventTextBg:     "default",
		SelectedEventFg: "black|bold",
		SelectedEventBg: "bright_yellow",
		NoEventsFg:      "bright_white",
		NoEventsBg:      "default",
		MoreEventsFg:    "bright_magenta",
		MoreEventsBg:    "default",
		ErrorFg:         "bright_red",
		ErrorBg:         "default",
		SuccessFg:       "bright_green",
		SuccessBg:       "default",
		InputFg:         "black|bold",
		InputBg:         "bright_yellow",
		SearchResultFg:  "bright_white",
		SearchResultBg:  "default",
		InstructionsFg:  "bright_cyan",
		InstructionsBg:  "default",
	}

	// LightTheme optimized for light backgrounds
	LightTheme = ColorTheme{
		MonthHeaderFg:   "blue|bold",
		MonthHeaderBg:   "default",
		DayHeaderFg:     "blue",
		DayHeaderBg:     "default",
		RegularDayFg:    "black",
		RegularDayBg:    "default",
		TodayFg:         "red|bold",
		TodayBg:         "default",
		SelectedFg:      "white|bold",
		SelectedBg:      "blue",
		SelectedTodayFg: "white|bold",
		SelectedTodayBg: "red",
		EventDayFg:      "green|bold",
		EventDayBg:      "default",
		EventHeaderFg:   "blue|bold",
		EventHeaderBg:   "default",
		EventTextFg:     "black",
		EventTextBg:     "default",
		SelectedEventFg: "white|bold",
		SelectedEventBg: "blue",
		NoEventsFg:      "black",
		NoEventsBg:      "default",
		MoreEventsFg:    "blue",
		MoreEventsBg:    "default",
		ErrorFg:         "red|bold",
		ErrorBg:         "default",
		SuccessFg:       "green|bold",
		SuccessBg:       "default",
		InputFg:         "black|bold",
		InputBg:         "white",
		SearchResultFg:  "black",
		SearchResultBg:  "default",
		InstructionsFg:  "blue",
		InstructionsBg:  "default",
	}
)

// ParseColor converts a color string like "magenta|bold" to termbox color attributes
func ParseColor(colorStr string) (termbox.Attribute, error) {
	if colorStr == "" || colorStr == "default" {
		return termbox.ColorDefault, nil
	}

	// Split color and attributes
	parts := strings.Split(colorStr, "|")
	colorName := strings.TrimSpace(parts[0])

	// Color name mapping
	colorMap := map[string]termbox.Attribute{
		"default":        termbox.ColorDefault,
		"black":          termbox.ColorBlack,
		"red":            termbox.ColorRed,
		"green":          termbox.ColorGreen,
		"yellow":         termbox.ColorYellow,
		"blue":           termbox.ColorBlue,
		"magenta":        termbox.ColorMagenta,
		"cyan":           termbox.ColorCyan,
		"white":          termbox.ColorWhite,
		"bright_black":   termbox.ColorBlack | termbox.AttrBold,
		"bright_red":     termbox.ColorRed | termbox.AttrBold,
		"bright_green":   termbox.ColorGreen | termbox.AttrBold,
		"bright_yellow":  termbox.ColorYellow | termbox.AttrBold,
		"bright_blue":    termbox.ColorBlue | termbox.AttrBold,
		"bright_magenta": termbox.ColorMagenta | termbox.AttrBold,
		"bright_cyan":    termbox.ColorCyan | termbox.AttrBold,
		"bright_white":   termbox.ColorWhite | termbox.AttrBold,
	}

	color, exists := colorMap[colorName]
	if !exists {
		return termbox.ColorDefault, fmt.Errorf("unknown color: %s", colorName)
	}

	// Apply additional attributes
	for i := 1; i < len(parts); i++ {
		attr := strings.TrimSpace(parts[i])
		switch attr {
		case "bold":
			color |= termbox.AttrBold
		case "underline":
			color |= termbox.AttrUnderline
		case "reverse":
			color |= termbox.AttrReverse
		default:
			return termbox.ColorDefault, fmt.Errorf("unknown attribute: %s", attr)
		}
	}

	return color, nil
}

// ValidateColorTheme validates that all colors in a theme are parseable
func ValidateColorTheme(theme *ColorTheme) error {
	colorFields := []string{
		theme.MonthHeaderFg, theme.MonthHeaderBg,
		theme.DayHeaderFg, theme.DayHeaderBg,
		theme.RegularDayFg, theme.RegularDayBg,
		theme.TodayFg, theme.TodayBg,
		theme.SelectedFg, theme.SelectedBg,
		theme.SelectedTodayFg, theme.SelectedTodayBg,
		theme.EventDayFg, theme.EventDayBg,
		theme.EventHeaderFg, theme.EventHeaderBg,
		theme.EventTextFg, theme.EventTextBg,
		theme.SelectedEventFg, theme.SelectedEventBg,
		theme.NoEventsFg, theme.NoEventsBg,
		theme.MoreEventsFg, theme.MoreEventsBg,
		theme.ErrorFg, theme.ErrorBg,
		theme.SuccessFg, theme.SuccessBg,
		theme.InputFg, theme.InputBg,
		theme.SearchResultFg, theme.SearchResultBg,
		theme.InstructionsFg, theme.InstructionsBg,
	}

	for _, colorStr := range colorFields {
		if _, err := ParseColor(colorStr); err != nil {
			return fmt.Errorf("invalid color '%s': %v", colorStr, err)
		}
	}

	return nil
}

// GetThemeByName returns a predefined theme by name
func GetThemeByName(name string) (ColorTheme, error) {
	switch strings.ToLower(name) {
	case "default":
		return DefaultTheme, nil
	case "dark":
		return DarkTheme, nil
	case "light":
		return LightTheme, nil
	default:
		return DefaultTheme, fmt.Errorf("unknown theme: %s", name)
	}
}

// Config holds the application configuration
type Config struct {
	EventsFilePath string       `json:"events_file_path"`
	ConfigFilePath string       `json:"-"` // Don't serialize this field
	WeekStartDay   WeekStartDay `json:"week_start_day"`
	UITheme        ColorTheme   `json:"ui_theme"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not accessible
		homeDir = "."
	}

	configDir := filepath.Join(homeDir, ".ascii-calendar")

	return &Config{
		EventsFilePath: filepath.Join(configDir, "events.json"),
		ConfigFilePath: filepath.Join(configDir, "configuration.json"),
		WeekStartDay:   StartSunday, // Default to Sunday-first
		UITheme:        DefaultTheme,
	}
}

// LoadConfig loads configuration from command line arguments and configuration file
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Parse command line arguments
	var configFileFlag string
	var eventsFileFlag string

	flag.StringVar(&configFileFlag, "c", "", "Path to configuration file")
	flag.StringVar(&eventsFileFlag, "f", "", "Path to events file")
	flag.Parse()

	// Use command line config file path if provided
	if configFileFlag != "" {
		config.ConfigFilePath = configFileFlag
	}

	// Try to load configuration file
	if err := config.loadFromFile(); err != nil {
		// If configuration file doesn't exist, that's okay - use defaults
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load configuration file: %v", err)
		}
	}

	// Command line arguments override configuration file
	if eventsFileFlag != "" {
		config.EventsFilePath = eventsFileFlag
	}

	// Ensure the directory exists
	if err := config.ensureDirectoryExists(); err != nil {
		return nil, fmt.Errorf("failed to create configuration directory: %v", err)
	}

	return config, nil
}

// loadFromFile loads configuration from the configuration file
func (c *Config) loadFromFile() error {
	file, err := os.Open(c.ConfigFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(c)
}

// SaveToFile saves the current configuration to the configuration file
func (c *Config) SaveToFile() error {
	// Ensure directory exists
	if err := c.ensureDirectoryExists(); err != nil {
		return fmt.Errorf("failed to create configuration directory: %v", err)
	}

	file, err := os.Create(c.ConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to create configuration file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	return encoder.Encode(c)
}

// ensureDirectoryExists creates the configuration directory if it doesn't exist
func (c *Config) ensureDirectoryExists() error {
	// Get directory from events file path (since that's where we store everything)
	dir := filepath.Dir(c.EventsFilePath)

	// Create directory with appropriate permissions
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	return nil
}

// GetEventsFilePath returns the full path to the events file
func (c *Config) GetEventsFilePath() string {
	return c.EventsFilePath
}

// GetConfigFilePath returns the full path to the configuration file
func (c *Config) GetConfigFilePath() string {
	return c.ConfigFilePath
}
