package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	EventsFilePath string `json:"events_file_path"`
	ConfigFilePath string `json:"-"` // Don't serialize this field
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
