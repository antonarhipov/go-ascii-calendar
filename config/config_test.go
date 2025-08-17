package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}

	// Verify events file path is set
	if config.EventsFilePath == "" {
		t.Error("EventsFilePath should not be empty")
	}

	// Verify config file path is set
	if config.ConfigFilePath == "" {
		t.Error("ConfigFilePath should not be empty")
	}

	// Verify paths contain .ascii-calendar directory
	if !strings.Contains(config.EventsFilePath, ".ascii-calendar") {
		t.Error("EventsFilePath should contain .ascii-calendar directory")
	}

	if !strings.Contains(config.ConfigFilePath, ".ascii-calendar") {
		t.Error("ConfigFilePath should contain .ascii-calendar directory")
	}

	// Verify file extensions
	if !strings.HasSuffix(config.EventsFilePath, "events.json") {
		t.Error("EventsFilePath should end with events.json")
	}

	if !strings.HasSuffix(config.ConfigFilePath, "configuration.json") {
		t.Error("ConfigFilePath should end with configuration.json")
	}
}

func TestConfig_GetEventsFilePath(t *testing.T) {
	config := &Config{
		EventsFilePath: "/test/path/events.json",
	}

	result := config.GetEventsFilePath()
	expected := "/test/path/events.json"

	if result != expected {
		t.Errorf("GetEventsFilePath() = %s, want %s", result, expected)
	}
}

func TestConfig_GetConfigFilePath(t *testing.T) {
	config := &Config{
		ConfigFilePath: "/test/path/config.json",
	}

	result := config.GetConfigFilePath()
	expected := "/test/path/config.json"

	if result != expected {
		t.Errorf("GetConfigFilePath() = %s, want %s", result, expected)
	}
}

func TestConfig_SaveToFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test successful save
	config := &Config{
		EventsFilePath: filepath.Join(tempDir, "events.json"),
		ConfigFilePath: filepath.Join(tempDir, "config.json"),
	}

	err = config.SaveToFile()
	if err != nil {
		t.Errorf("SaveToFile() failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(config.ConfigFilePath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Verify file contents
	fileContent, err := os.ReadFile(config.ConfigFilePath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var savedConfig Config
	err = json.Unmarshal(fileContent, &savedConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if savedConfig.EventsFilePath != config.EventsFilePath {
		t.Errorf("Saved EventsFilePath = %s, want %s", savedConfig.EventsFilePath, config.EventsFilePath)
	}
}

func TestConfig_SaveToFile_InvalidPath(t *testing.T) {
	// Test save to invalid path (should fail)
	config := &Config{
		EventsFilePath: "/nonexistent/invalid/path/events.json",
		ConfigFilePath: "/nonexistent/invalid/path/config.json",
	}

	err := config.SaveToFile()
	if err == nil {
		t.Error("SaveToFile() should have failed with invalid path")
	}
}

func TestConfig_loadFromFile(t *testing.T) {
	// Create temporary directory and config file
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test_config.json")
	testConfig := &Config{
		EventsFilePath: "/test/events/path.json",
	}

	// Write test config to file
	file, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(testConfig)
	file.Close()
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test loading from file
	config := &Config{ConfigFilePath: configPath}
	err = config.loadFromFile()
	if err != nil {
		t.Errorf("loadFromFile() failed: %v", err)
	}

	if config.EventsFilePath != testConfig.EventsFilePath {
		t.Errorf("Loaded EventsFilePath = %s, want %s", config.EventsFilePath, testConfig.EventsFilePath)
	}
}

func TestConfig_loadFromFile_NonexistentFile(t *testing.T) {
	config := &Config{
		ConfigFilePath: "/nonexistent/file.json",
	}

	err := config.loadFromFile()
	if err == nil {
		t.Error("loadFromFile() should have failed with nonexistent file")
	}

	if !os.IsNotExist(err) {
		t.Errorf("loadFromFile() should return os.IsNotExist error, got: %v", err)
	}
}

func TestConfig_loadFromFile_InvalidJSON(t *testing.T) {
	// Create temporary file with invalid JSON
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(configPath, []byte("invalid json content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid JSON file: %v", err)
	}

	config := &Config{ConfigFilePath: configPath}
	err = config.loadFromFile()
	if err == nil {
		t.Error("loadFromFile() should have failed with invalid JSON")
	}
}

func TestConfig_ensureDirectoryExists(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test creating new directory
	newDir := filepath.Join(tempDir, "new", "nested", "directory")
	config := &Config{
		EventsFilePath: filepath.Join(newDir, "events.json"),
	}

	err = config.ensureDirectoryExists()
	if err != nil {
		t.Errorf("ensureDirectoryExists() failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}
}

func TestConfig_ensureDirectoryExists_ExistingDirectory(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with existing directory
	config := &Config{
		EventsFilePath: filepath.Join(tempDir, "events.json"),
	}

	err = config.ensureDirectoryExists()
	if err != nil {
		t.Errorf("ensureDirectoryExists() failed with existing directory: %v", err)
	}
}

func TestLoadConfig_WithDefaults(t *testing.T) {
	// Save original command line args
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// Clear command line args
	os.Args = []string{"test"}

	config, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() failed: %v", err)
	}

	if config == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// Verify default values are set
	if config.EventsFilePath == "" {
		t.Error("EventsFilePath should not be empty")
	}

	if config.ConfigFilePath == "" {
		t.Error("ConfigFilePath should not be empty")
	}
}

func TestLoadConfig_WithCommandLineArgs(t *testing.T) {
	// Skip this test for now due to flag redefinition issues
	// The LoadConfig function uses global flags which can't be easily reset in tests
	t.Skip("Skipping LoadConfig command line test due to global flag limitations")
}

func TestLoadConfig_WithConfigFile(t *testing.T) {
	// Skip this test for now due to flag redefinition issues
	// The LoadConfig function uses global flags which can't be easily reset in tests
	t.Skip("Skipping LoadConfig config file test due to global flag limitations")
}
