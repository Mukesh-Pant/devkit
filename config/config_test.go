package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	
	if m == nil {
		t.Fatal("NewManager() returned nil")
	}
	
	if m.v == nil {
		t.Fatal("Manager.v is nil")
	}
}

func TestManager_DefaultValues(t *testing.T) {
	m := NewManager()
	
	tests := []struct {
		name     string
		key      string
		expected interface{}
		getter   func(string) interface{}
	}{
		{
			name:     "default timeout",
			key:      "devkit.timeout",
			expected: 5 * time.Second,
			getter:   func(k string) interface{} { return m.GetDuration(k) },
		},
		{
			name:     "default output",
			key:      "devkit.output",
			expected: "table",
			getter:   func(k string) interface{} { return m.GetString(k) },
		},
		{
			name:     "default slack_webhook",
			key:      "devkit.slack_webhook",
			expected: "",
			getter:   func(k string) interface{} { return m.GetString(k) },
		},
		{
			name:     "default watch interval",
			key:      "watch.interval",
			expected: 30 * time.Second,
			getter:   func(k string) interface{} { return m.GetDuration(k) },
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.getter(tt.key)
			if got != tt.expected {
				t.Errorf("GetDefault(%s) = %v, want %v", tt.key, got, tt.expected)
			}
		})
	}
}

func TestManager_GetStringSlice(t *testing.T) {
	m := NewManager()
	
	// Test default empty slice
	urls := m.GetStringSlice("watch.urls")
	if len(urls) != 0 {
		t.Errorf("Default watch.urls should be empty, got %v", urls)
	}
}

func TestManager_IsSet(t *testing.T) {
	m := NewManager()
	
	// Default values should be set
	if !m.IsSet("devkit.timeout") {
		t.Error("devkit.timeout should be set by default")
	}
	
	// Non-existent key should not be set
	if m.IsSet("nonexistent.key") {
		t.Error("nonexistent.key should not be set")
	}
}

func TestManager_Load_NoConfigFile(t *testing.T) {
	// Create a temporary directory with no config file
	tmpDir := t.TempDir()
	
	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	m := NewManager()
	err = m.Load()
	
	// Should not return error when config file is not found
	if err != nil {
		t.Errorf("Load() should not error when config file not found, got: %v", err)
	}
	
	// Should still have default values
	if m.GetString("devkit.output") != "table" {
		t.Error("Should have default values when config file not found")
	}
}

func TestManager_Load_FromCurrentDirectory(t *testing.T) {
	// Create a temporary directory with a config file
	tmpDir := t.TempDir()
	
	configContent := `devkit:
  timeout: 10s
  output: json
  slack_webhook: https://hooks.slack.com/test
watch:
  interval: 60s
  urls:
    - https://example.com
    - https://example.org
`
	
	configPath := filepath.Join(tmpDir, ".devkit.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	m := NewManager()
	err = m.Load()
	
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	
	// Verify loaded values
	if got := m.GetDuration("devkit.timeout"); got != 10*time.Second {
		t.Errorf("devkit.timeout = %v, want 10s", got)
	}
	
	if got := m.GetString("devkit.output"); got != "json" {
		t.Errorf("devkit.output = %v, want json", got)
	}
	
	if got := m.GetString("devkit.slack_webhook"); got != "https://hooks.slack.com/test" {
		t.Errorf("devkit.slack_webhook = %v, want https://hooks.slack.com/test", got)
	}
	
	if got := m.GetDuration("watch.interval"); got != 60*time.Second {
		t.Errorf("watch.interval = %v, want 60s", got)
	}
	
	urls := m.GetStringSlice("watch.urls")
	if len(urls) != 2 {
		t.Errorf("watch.urls length = %d, want 2", len(urls))
	}
	if len(urls) > 0 && urls[0] != "https://example.com" {
		t.Errorf("watch.urls[0] = %v, want https://example.com", urls[0])
	}
}

func TestManager_LoadFrom_ValidFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	
	configContent := `devkit:
  timeout: 15s
  output: plain
`
	
	configPath := filepath.Join(tmpDir, "custom-config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	m := NewManager()
	err := m.LoadFrom(configPath)
	
	if err != nil {
		t.Fatalf("LoadFrom() error = %v", err)
	}
	
	// Verify loaded values
	if got := m.GetDuration("devkit.timeout"); got != 15*time.Second {
		t.Errorf("devkit.timeout = %v, want 15s", got)
	}
	
	if got := m.GetString("devkit.output"); got != "plain" {
		t.Errorf("devkit.output = %v, want plain", got)
	}
}

func TestManager_LoadFrom_NonExistentFile(t *testing.T) {
	m := NewManager()
	err := m.LoadFrom("/nonexistent/path/config.yaml")
	
	if err == nil {
		t.Error("LoadFrom() should return error for non-existent file")
	}
}

func TestManager_LoadFrom_InvalidYAML(t *testing.T) {
	// Create a temporary config file with invalid YAML
	tmpDir := t.TempDir()
	
	invalidContent := `devkit:
  timeout: 15s
  output: plain
    invalid indentation
`
	
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	m := NewManager()
	err := m.LoadFrom(configPath)
	
	if err == nil {
		t.Error("LoadFrom() should return error for invalid YAML")
	}
}

func TestManager_GetInt(t *testing.T) {
	tmpDir := t.TempDir()
	
	configContent := `devkit:
  max_retries: 3
`
	
	configPath := filepath.Join(tmpDir, ".devkit.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	m := NewManager()
	if err := m.LoadFrom(configPath); err != nil {
		t.Fatal(err)
	}
	
	if got := m.GetInt("devkit.max_retries"); got != 3 {
		t.Errorf("GetInt(devkit.max_retries) = %d, want 3", got)
	}
}

func TestManager_Load_FromHomeDirectory(t *testing.T) {
	// This test verifies that the home directory is added to the search path
	// We can't easily test actual loading from home without modifying the user's home
	// So we just verify the manager is created and Load doesn't panic
	
	m := NewManager()
	
	// Create a temp directory to use as working directory (no config file)
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// This should not panic and should return nil (config not found is okay)
	err = m.Load()
	if err != nil {
		t.Errorf("Load() should not error when searching home directory, got: %v", err)
	}
}
