package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Manager wraps Viper configuration
type Manager struct {
	v *viper.Viper
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	v := viper.New()
	
	// Set configuration file name and type
	v.SetConfigName(".devkit")
	v.SetConfigType("yaml")
	
	// Set default values
	v.SetDefault("devkit.timeout", "5s")
	v.SetDefault("devkit.output", "table")
	v.SetDefault("devkit.slack_webhook", "")
	v.SetDefault("watch.interval", "30s")
	v.SetDefault("watch.urls", []string{})
	
	return &Manager{v: v}
}

// Load attempts to load .devkit.yaml from current dir or home dir
func (m *Manager) Load() error {
	// Add current directory to search path
	m.v.AddConfigPath(".")
	
	// Add home directory to search path
	home, err := os.UserHomeDir()
	if err == nil {
		m.v.AddConfigPath(home)
	}
	
	// Attempt to read configuration
	if err := m.v.ReadInConfig(); err != nil {
		// If config file not found, that's okay - we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}
		// For other errors, return them
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	return nil
}

// LoadFrom loads configuration from a specific file path
func (m *Manager) LoadFrom(path string) error {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve config path: %w", err)
	}
	
	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", absPath)
	}
	
	// Set the config file explicitly
	m.v.SetConfigFile(absPath)
	
	// Read the configuration
	if err := m.v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %s: %w", absPath, err)
	}
	
	return nil
}

// GetString retrieves a string configuration value
func (m *Manager) GetString(key string) string {
	return m.v.GetString(key)
}

// GetInt retrieves an integer configuration value
func (m *Manager) GetInt(key string) int {
	return m.v.GetInt(key)
}

// GetStringSlice retrieves a string slice configuration value
func (m *Manager) GetStringSlice(key string) []string {
	return m.v.GetStringSlice(key)
}

// GetDuration retrieves a duration configuration value
func (m *Manager) GetDuration(key string) time.Duration {
	return m.v.GetDuration(key)
}

// IsSet checks if a configuration key is set
func (m *Manager) IsSet(key string) bool {
	return m.v.IsSet(key)
}
