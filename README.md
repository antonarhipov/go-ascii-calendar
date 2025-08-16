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

### Key Bindings

#### Navigation
- **B** or **b** - Move backward one month (shifts the three-month window)
- **N** or **n** - Move forward one month (shifts the three-month window)
- **H** or **h** - Move selection left (one day)
- **L** or **l** - Move selection right (one day)
- **K** or **k** - Move selection up (one week)
- **J** or **j** - Move selection down (one week)

#### Event Management
- **Enter** - View events for the currently selected date
- **A** or **a** - Add a new event to the selected date (only available when viewing events)
- **Esc** - Back to previous view / Cancel current operation

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

Events are stored in a plain text file called `events.txt` in the same directory as the application. The format is:

```
YYYY-MM-DD|HH:MM|description
```

### Examples

```
2025-08-16|09:00|Morning standup meeting
2025-08-16|14:30|Project review session
2025-08-17|10:00|Client presentation
2025-08-18|15:00|Team lunch
```

### File Format Rules

- **Date**: YYYY-MM-DD format (ISO 8601)
- **Time**: HH:MM format in 24-hour time
- **Separator**: Pipe character (|) separates fields
- **Description**: Can contain spaces and most printable characters
- **Encoding**: UTF-8 text file
- **Line endings**: Unix-style (LF) preferred, but CR+LF also supported

### Manual Editing

You can manually edit the `events.txt` file with any text editor. The application will load changes on next startup. Invalid lines are skipped with a warning message.

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
- Check that the application has write permissions in its directory
- Verify the `events.txt` file exists and isn't read-only
- Look for error messages during event creation
- Check available disk space

#### Events File Corruption

**Problem**: Error messages about malformed lines or events not loading properly.

**Solutions**:
- Open `events.txt` in a text editor and check the format
- Ensure each line follows the `YYYY-MM-DD|HH:MM|description` format
- Remove or fix any lines that don't match the expected format
- Check for invisible characters or encoding issues

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