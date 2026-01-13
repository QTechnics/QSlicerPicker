package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDirName  = ".qslicerpicker"
	ConfigFileName = "config.json"
)

type Config struct {
	Language      string          `json:"language"`
	Slicers       []SlicerConfig  `json:"slicers"`
	CustomSlicers []CustomSlicer  `json:"custom_slicers"`
}

type SlicerConfig struct {
	ID         string   `json:"id"`
	Enabled    bool     `json:"enabled"`
	Order      int      `json:"order"`
	CustomPath string   `json:"custom_path,omitempty"`
	Arguments  []string `json:"arguments,omitempty"`
	WorkingDir string   `json:"working_dir,omitempty"`
}

type CustomSlicer struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Arguments   []string `json:"arguments,omitempty"`
	WorkingDir  string   `json:"working_dir,omitempty"`
	Enabled     bool     `json:"enabled"`
	Order       int      `json:"order"`
}

var (
	configInstance *Config
	configPath     string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}

	configDir := filepath.Join(homeDir, ConfigDirName)
	configPath = filepath.Join(configDir, ConfigFileName)

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create config directory: %v", err))
	}
}

// GetConfig returns the current configuration, loading it if necessary
func GetConfig() *Config {
	if configInstance == nil {
		configInstance = LoadConfig()
	}
	return configInstance
}

// LoadConfig loads the configuration from file, or returns default if file doesn't exist
func LoadConfig() *Config {
	config := &Config{
		Language:      "en",
		Slicers:       []SlicerConfig{},
		CustomSlicers: []CustomSlicer{},
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config
		}
		// If there's an error reading, return default config
		return config
	}

	if err := json.Unmarshal(data, config); err != nil {
		// If there's an error parsing, return default config
		return config
	}

	return config
}

// SaveConfig saves the current configuration to file
func SaveConfig() error {
	config := GetConfig()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	return configPath
}
