# ASCII Calendar Configuration Guide

This document describes the configuration options available for the ASCII Calendar application and how to customize its appearance and behavior.

## Configuration File Location

The ASCII Calendar application uses a JSON configuration file located at:
- **Default location**: `~/.ascii-calendar/configuration.json`
- **Custom location**: Specify with `-c <config-file>` command line option

If the configuration file doesn't exist, the application will use default values and create the file automatically when settings are saved.

## Command Line Options

- `-c <config-file>`: Specify custom configuration file path
- `-f <events-file>`: Override events file path (takes precedence over config file setting)

## Configuration Structure

### Basic Configuration

```json
{
  "events_file_path": "~/.ascii-calendar/events.json",
  "week_start_day": 0,
  "ui_theme": {
    "month_header_fg": "magenta|bold",
    "day_header_fg": "cyan"
  }
}
```

### Configuration Fields

#### `events_file_path` (string)
Path to the JSON file where events are stored.
- Supports absolute and relative paths
- Supports `~` for home directory expansion
- **Default**: `~/.ascii-calendar/events.json`

#### `week_start_day` (integer)
Determines which day of the week appears first in the calendar.
- `0`: Sunday first (default)
- `1`: Monday first

#### `ui_theme` (object)
Complete color theme configuration for all UI elements. See [Color Theme Configuration](#color-theme-configuration) below.

## Color Theme Configuration

The `ui_theme` object allows customization of colors for all visual elements in the application.

### Color Syntax

Colors are specified as strings with the following format:
- **Basic colors**: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`, `default`
- **Bright colors**: `bright_black`, `bright_red`, `bright_green`, `bright_yellow`, `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`
- **Attributes**: `bold`, `underline`, `reverse`
- **Combinations**: Use `|` to combine (e.g., `red|bold`, `cyan|underline`)

### Theme Color Fields

Each UI element has separate foreground (`_fg`) and background (`_bg`) color settings:

#### Calendar Elements
- `month_header_fg/bg`: Month/year headers (e.g., "August 2025")
- `day_header_fg/bg`: Day-of-week headers (e.g., "Su Mo Tu We Th Fr Sa")
- `regular_day_fg/bg`: Regular calendar day numbers
- `today_fg/bg`: Today's date highlighting
- `selected_fg/bg`: Currently selected date
- `selected_today_fg/bg`: When selected date is also today
- `event_day_fg/bg`: Days that have events

#### Event Display Elements
- `event_header_fg/bg`: Event list section headers
- `event_text_fg/bg`: Event list item text
- `selected_event_fg/bg`: Selected event in event lists
- `no_events_fg/bg`: "No events" messages
- `more_events_fg/bg`: "... and X more events" indicators

#### Interface Elements
- `error_fg/bg`: Error messages
- `success_fg/bg`: Success messages
- `input_fg/bg`: Input prompts and fields
- `search_result_fg/bg`: Search result items
- `instructions_fg/bg`: Key legends and instruction text

## Predefined Themes

The application includes three predefined themes:

### Default Theme
Moderate colors suitable for most terminals:
```json
{
  "month_header_fg": "magenta|bold",
  "day_header_fg": "cyan",
  "today_fg": "yellow|bold",
  "selected_fg": "white|bold",
  "selected_bg": "blue",
  "event_day_fg": "green"
}
```

### Dark Theme
High-contrast theme optimized for dark terminal backgrounds:
```json
{
  "month_header_fg": "bright_magenta|bold",
  "day_header_fg": "bright_cyan",
  "regular_day_fg": "white",
  "today_fg": "bright_yellow|bold",
  "selected_fg": "black|bold",
  "selected_bg": "bright_blue",
  "event_day_fg": "bright_green"
}
```

### Light Theme
Optimized for light terminal backgrounds with darker text:
```json
{
  "month_header_fg": "blue|bold",
  "day_header_fg": "blue",
  "regular_day_fg": "black",
  "today_fg": "red|bold",
  "selected_fg": "white|bold",
  "selected_bg": "blue",
  "event_day_fg": "green|bold"
}
```

## Complete Configuration Example

```json
{
  "events_file_path": "~/.ascii-calendar/events.json",
  "week_start_day": 0,
  "ui_theme": {
    "month_header_fg": "magenta|bold",
    "month_header_bg": "default",
    "day_header_fg": "cyan",
    "day_header_bg": "default",
    "regular_day_fg": "default",
    "regular_day_bg": "default",
    "today_fg": "yellow|bold",
    "today_bg": "default",
    "selected_fg": "white|bold",
    "selected_bg": "blue",
    "selected_today_fg": "white|bold",
    "selected_today_bg": "cyan",
    "event_day_fg": "green",
    "event_day_bg": "default",
    "event_header_fg": "yellow|bold",
    "event_header_bg": "default",
    "event_text_fg": "white",
    "event_text_bg": "default",
    "selected_event_fg": "black|bold",
    "selected_event_bg": "yellow",
    "no_events_fg": "white",
    "no_events_bg": "default",
    "more_events_fg": "magenta",
    "more_events_bg": "default",
    "error_fg": "red",
    "error_bg": "default",
    "success_fg": "green",
    "success_bg": "default",
    "input_fg": "black|bold",
    "input_bg": "yellow",
    "search_result_fg": "white",
    "search_result_bg": "default",
    "instructions_fg": "cyan",
    "instructions_bg": "default"
  }
}
```

## Terminal Compatibility

### Color Support
- **Color terminals**: Full theme support with all specified colors
- **Monochrome terminals**: Automatic fallback to text attributes (bold, underline, reverse)
- **Invalid colors**: Fall back to defaults with error logging

### Terminal Types
- Works with any terminal that supports basic ANSI color codes
- Optimized for common terminals: xterm, gnome-terminal, iTerm2, Terminal.app, etc.
- Graceful degradation on limited terminals

## Configuration Management

### Creating Configuration
1. Run the application - it will create default configuration automatically
2. Or manually create `~/.ascii-calendar/configuration.json` with desired settings
3. Or copy and modify the example configuration from this documentation

### Editing Configuration
1. Edit the JSON file directly with any text editor
2. Restart the application to apply changes
3. Invalid configurations will fall back to defaults with error messages

### Validation
- Color values are validated at startup
- Invalid colors are logged and fall back to defaults
- Malformed JSON files are reported with helpful error messages
- Missing fields use default values

## Troubleshooting

### Common Issues
1. **Colors not appearing**: Check if your terminal supports color
2. **Configuration ignored**: Verify JSON syntax is valid
3. **File not found errors**: Ensure the directory `~/.ascii-calendar` exists and is writable
4. **Invalid colors**: Check color names match the supported list above

### Debug Tips
- Start the application from command line to see error messages
- Use `-c` option to specify a test configuration file
- Validate JSON syntax with online JSON validators
- Test with minimal configuration first, then add customizations

## Examples

### Monday-First Calendar
```json
{
  "week_start_day": 1
}
```

### Custom Events Location
```json
{
  "events_file_path": "/path/to/my/calendar-events.json"
}
```

### Minimal Dark Theme
```json
{
  "ui_theme": {
    "regular_day_fg": "white",
    "today_fg": "bright_yellow|bold",
    "selected_bg": "bright_blue",
    "event_day_fg": "bright_green"
  }
}
```