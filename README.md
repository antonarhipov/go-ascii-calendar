# ASCII Calendar

A terminal-based calendar application for managing events with a clean, keyboard-driven interface. The application displays three months side-by-side (previous, current, next) and allows you to navigate, view events, and add new events using simple keyboard controls.

## Features

- **Three-month view**: See previous, current, and next months at once
- **Keyboard navigation**: Vim-style navigation keys (H/J/K/L) plus month switching (B/N)
- **Event management**: View and add events with time and description
- **Visual indicators**: See which days have events and today's date highlighted
- **Data persistence**: Events are saved to a text file and persist across sessions
- **Terminal compatibility**: Works in standard 80x24 monochrome terminals
- **ASCII-safe**: Uses only standard ASCII characters for maximum compatibility

## Installation

### Prerequisites

- Go 1.19 or later
- Terminal with at least 80x24 character display

### Building from Source

1. Clone or download the source code:
```bash
git clone <repository-url>
cd go-ascii-calendar
```

2. Build the application:
```bash
go build -o ascii-calendar
```

3. Run the application:
```bash
./ascii-calendar
```

### Alternative Build Method

You can also run directly without building:
```bash
go run main.go
```

## Usage

### Starting the Application

Launch the application from your terminal:
```bash
./ascii-calendar
```

The application will display three calendar months side-by-side, with today's date highlighted and any existing events indicated.

### Command Line Options

The application supports several command line options for configuration:

```bash
# Use default configuration
./ascii-calendar

# Specify custom events file location
./ascii-calendar -f /path/to/my/events.json

# Specify custom configuration file location
./ascii-calendar -c /path/to/my/config.json

# Use both custom configuration and events file
./ascii-calendar -c /path/to/config.json -f /path/to/events.json
```

**Available Options:**
- `-f <path>` - Path to events file (defaults to `~/.ascii-calendar/events.json`)
- `-c <path>` - Path to configuration file (defaults to `~/.ascii-calendar/configuration.json`)
- `-h` - Show help message with available options

### Key Bindings

#### Navigation
- **B** or **b** - Move backward one month (shifts the three-month window)
- **N** or **n** - Move forward one month (shifts the three-month window)
- **H** or **h** / **Left Arrow** - Move selection left (one day)
- **L** or **l** / **Right Arrow** - Move selection right (one day)
- **K** or **k** / **Up Arrow** - Move selection up (one week)
- **J** or **j** / **Down Arrow** - Move selection down (one week)
- **C** or **c** - Reset calendar to current month and select today's date

#### Event Management
- **Enter** - View events for the currently selected date
- **A** or **a** - Add a new event to the selected date (only available when viewing events)
- **Esc** - Exit application (from main calendar) / Back to previous view / Cancel current operation

#### Application Control
- **Q** or **q** - Quit the application
- **Ctrl+C** - Force quit the application

### Adding Events

1. Navigate to the desired date using the arrow keys
2. Press **Enter** to view events for that date
3. Press **A** to add a new event
4. Enter the time in HH:MM format (24-hour time, e.g., "14:30" for 2:30 PM)
5. Enter a description for the event
6. Press **Enter** to save, or **Esc** to cancel

### Visual Indicators

- **[Today]**: Current date is highlighted with square brackets
- **Selected**: Currently selected date has a different visual indication
- **Events**: Days with events show a dot (•) indicator
- **Combined**: Days can show multiple indicators (e.g., today + events)

## Event File Format

Events are stored in JSON format in `~/.ascii-calendar/events.json` by default. The JSON structure is:

```json
{
  "events": [
    {
      "date": "2025-08-16",
      "time": "09:00",
      "description": "Morning standup meeting"
    },
    {
      "date": "2025-08-16",
      "time": "14:30",
      "description": "Project review session"
    },
    {
      "date": "2025-08-17",
      "time": "10:00",
      "description": "Client presentation"
    }
  ]
}
```

### Configuration File

The application can be configured using a JSON configuration file at `~/.ascii-calendar/configuration.json`:

```json
{
  "events_file_path": "/path/to/custom/events.json"
}
```

### File Format Rules

- **Date**: YYYY-MM-DD format (ISO 8601)
- **Time**: HH:MM format in 24-hour time
- **Description**: Can contain spaces and most printable characters
- **Encoding**: UTF-8 JSON file
- **Location**: `~/.ascii-calendar/events.json` (configurable)

### Manual Editing

You can manually edit the JSON files with any text editor. The application will load changes on next startup. Invalid JSON will cause the application to report errors.

### Migration from Old Format

If you have an existing `events.txt` file in the old pipe-separated format, the application will automatically migrate it to the new JSON format on first startup. The migration will:

1. Read events from the old `events.txt` file
2. Convert them to JSON format 
3. Save them to `~/.ascii-calendar/events.json`
4. Display a confirmation message

The old `events.txt` file will remain unchanged for backup purposes.

## Troubleshooting

### Common Issues

#### Application Won't Start

**Problem**: Application exits immediately or shows initialization errors.

**Solutions**:
- Ensure your terminal supports at least 80x24 characters
- Check that you have Go 1.19 or later installed
- Verify the executable has proper permissions
- Try running in a different terminal application

#### Display Issues

**Problem**: Calendar appears garbled or improperly formatted.

**Solutions**:
- Resize terminal to at least 80 characters wide and 24 lines tall
- Use a monospace font in your terminal
- Check terminal encoding is set to UTF-8
- Try a different terminal application (some terminals have rendering issues)

#### Navigation Not Working

**Problem**: Key presses don't move the selection or trigger actions.

**Solutions**:
- Ensure you're using the correct keys (H/J/K/L for day navigation, B/N for months)
- Try both uppercase and lowercase versions of the keys
- Check that your terminal isn't intercepting the key combinations
- Some SSH connections may interfere with certain key combinations

#### Events Not Saving

**Problem**: Events disappear after restarting the application.

**Solutions**:
- Check that the application has write permissions to `~/.ascii-calendar/` directory
- Verify the `~/.ascii-calendar/events.json` file exists and isn't read-only
- Look for error messages during event creation
- Check available disk space in your home directory
- Try running with `-f /tmp/test-events.json` to test with a different location

#### Configuration Issues

**Problem**: Application doesn't use custom configuration or file paths.

**Solutions**:
- Verify the configuration file is valid JSON format
- Check that command line arguments are properly specified (`-c` for config, `-f` for events file)
- Ensure the `~/.ascii-calendar/` directory exists and is writable
- Try using absolute paths instead of relative paths
- Run `./ascii-calendar -h` to see available options

#### JSON File Corruption

**Problem**: Error messages about malformed JSON or events not loading properly.

**Solutions**:
- Open `~/.ascii-calendar/events.json` in a text editor and validate the JSON format
- Use an online JSON validator to check for syntax errors
- Ensure the file follows the correct structure with `"events"` array
- Check that all dates are in `YYYY-MM-DD` format and times in `HH:MM` format
- Restore from backup or delete the file to start fresh (application will recreate it)

#### Migration Issues

**Problem**: Old `events.txt` file not migrating to new JSON format.

**Solutions**:
- Ensure the old `events.txt` file is in the same directory as the application
- Check that the `~/.ascii-calendar/` directory can be created (permissions)
- Look for migration messages in the terminal output
- Manually specify the events file location with `-f` if needed
- Verify the old file format follows `YYYY-MM-DD|HH:MM|description` pattern

### Performance Issues

#### Slow Rendering

**Problem**: Calendar updates slowly or flickers.

**Solutions**:
- Use a modern terminal application
- Reduce terminal font size if using very large fonts
- Close other resource-intensive applications
- Try running on a different system to isolate hardware issues

### Terminal Compatibility

#### Supported Terminals

The application has been tested on:
- **Linux**: gnome-terminal, konsole, xterm, tmux, screen
- **macOS**: Terminal.app, iTerm2, tmux
- **Windows**: Windows Terminal, WSL terminals

#### Known Issues

- Some very old terminal applications may have rendering issues
- Windows Command Prompt (cmd.exe) may have limited compatibility
- SSH connections with high latency may cause input delays
- Terminal multiplexers (tmux/screen) may require specific configuration

### Getting Help

If you encounter issues not covered here:

1. Check that your system meets the minimum requirements
2. Try reproducing the issue in a different terminal
3. Look for error messages in the terminal output
4. Verify the `events.txt` file format if event-related issues occur
5. Consider reinstalling or rebuilding the application

### Error Messages

- **"Terminal too small"**: Resize terminal to at least 80x24
- **"Failed to initialize terminal"**: Terminal compatibility issue
- **"Permission denied"**: Check file/directory permissions
- **"Invalid time format"**: Use HH:MM format for event times
- **"Empty description not allowed"**: Event descriptions cannot be empty

## System Requirements

- **Operating System**: Linux, macOS, Windows (with proper terminal)
- **Go Version**: 1.19 or later (for building from source)
- **Terminal**: Minimum 80x24 characters, monospace font recommended
- **Memory**: Minimal (< 10MB typical usage)
- **Storage**: Minimal (events.txt typically < 1KB per 100 events)

## License

This project is open source. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please ensure any changes maintain terminal compatibility and follow the existing code structure.