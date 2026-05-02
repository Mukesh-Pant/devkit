package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that root command can be executed
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute root command: %v", err)
	}
}

func TestConfigFlagWithValidFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")
	
	configContent := `devkit:
  timeout: 15s
  output: json
  slack_webhook: https://hooks.slack.com/test

watch:
  interval: 45s
  urls:
    - https://test.com
`
	
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	// Manually test the configuration loading
	cfgFile = configPath
	
	// Manually call the PreRunE logic
	err = rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify configuration was loaded
	if cfg == nil {
		t.Fatal("Configuration manager was not initialized")
	}
	
	// Verify configuration values
	timeout := cfg.GetDuration("devkit.timeout")
	if timeout.Seconds() != 15 {
		t.Errorf("Expected timeout 15s, got %v", timeout)
	}
	
	output := cfg.GetString("devkit.output")
	if output != "json" {
		t.Errorf("Expected output 'json', got '%s'", output)
	}
	
	webhook := cfg.GetString("devkit.slack_webhook")
	if webhook != "https://hooks.slack.com/test" {
		t.Errorf("Expected webhook 'https://hooks.slack.com/test', got '%s'", webhook)
	}
	
	interval := cfg.GetDuration("watch.interval")
	if interval.Seconds() != 45 {
		t.Errorf("Expected interval 45s, got %v", interval)
	}
	
	urls := cfg.GetStringSlice("watch.urls")
	if len(urls) != 1 || urls[0] != "https://test.com" {
		t.Errorf("Expected urls ['https://test.com'], got %v", urls)
	}
	
	// Reset
	cfgFile = ""
	cfg = nil
}

func TestConfigNotFound(t *testing.T) {
	// Set a non-existent config file
	cfgFile = "/nonexistent/config.yaml"
	
	// Manually call the PreRunE logic - should fail
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err == nil {
		t.Fatal("Expected error when config file not found, got nil")
	}
	
	// Reset
	cfgFile = ""
	cfg = nil
}

func TestDefaultConfig(t *testing.T) {
	// Reset config file flag
	cfgFile = ""
	
	// Manually call the PreRunE logic
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}
	
	// Verify configuration manager was initialized with defaults
	if cfg == nil {
		t.Fatal("Configuration manager was not initialized")
	}
	
	// Verify default values
	timeout := cfg.GetDuration("devkit.timeout")
	if timeout.Seconds() != 5 {
		t.Errorf("Expected default timeout 5s, got %v", timeout)
	}
	
	output := cfg.GetString("devkit.output")
	if output != "table" {
		t.Errorf("Expected default output 'table', got '%s'", output)
	}
	
	interval := cfg.GetDuration("watch.interval")
	if interval.Seconds() != 30 {
		t.Errorf("Expected default interval 30s, got %v", interval)
	}
	
	// Reset
	cfg = nil
}

func TestGetConfig(t *testing.T) {
	// Initialize config
	cfgFile = ""
	err := rootCmd.PersistentPreRunE(rootCmd, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}
	
	// Test GetConfig function
	retrievedCfg := GetConfig()
	if retrievedCfg == nil {
		t.Fatal("GetConfig returned nil")
	}
	
	if retrievedCfg != cfg {
		t.Error("GetConfig did not return the global config instance")
	}
	
	// Reset
	cfg = nil
}

