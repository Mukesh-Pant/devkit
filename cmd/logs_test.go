package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogsCmd_ValidCombinedLog(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logContent := `192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
192.168.1.2 - - [01/Jan/2024:12:00:01 +0000] "POST /api/login HTTP/1.1" 401 567 "-" "curl/7.68.0"
192.168.1.1 - - [01/Jan/2024:12:00:02 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Test the command
	rootCmd.SetArgs([]string{"logs", logFile, "--format", "combined", "--output", "json"})
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("logs command failed: %v", err)
	}
}

func TestLogsCmd_ValidJSONLog(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.json")

	logContent := `{"ip":"192.168.1.1","timestamp":"2024-01-01T12:00:00Z","method":"GET","path":"/api/users","status_code":200,"size":1234}
{"ip":"192.168.1.2","timestamp":"2024-01-01T12:00:01Z","method":"POST","path":"/api/login","status_code":401,"size":567}`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Test the command
	rootCmd.SetArgs([]string{"logs", logFile, "--format", "json", "--output", "json"})
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("logs command failed: %v", err)
	}
}

func TestLogsCmd_ValidPlainLog(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logContent := `2024-01-01 12:00:00 INFO 192.168.1.1 GET /api/users 200
2024-01-01 12:00:01 ERROR 192.168.1.2 POST /api/login 401`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Test the command
	rootCmd.SetArgs([]string{"logs", logFile, "--format", "plain", "--output", "plain"})
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("logs command failed: %v", err)
	}
}

func TestLogsCmd_NonExistentFile(t *testing.T) {
	// Test with a non-existent file
	rootCmd.SetArgs([]string{"logs", "nonexistent.log", "--format", "combined"})
	err := rootCmd.Execute()

	if err == nil {
		t.Error("logs command should fail with non-existent file")
	}
}

func TestLogsCmd_InvalidFormat(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	err := os.WriteFile(logFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Test with invalid format
	rootCmd.SetArgs([]string{"logs", logFile, "--format", "invalid"})
	err = rootCmd.Execute()

	if err == nil {
		t.Error("logs command should fail with invalid format")
	}
}

func TestLogsCmd_InvalidOutputFormat(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	err := os.WriteFile(logFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Test with invalid output format
	rootCmd.SetArgs([]string{"logs", logFile, "--format", "combined", "--output", "invalid"})
	err = rootCmd.Execute()

	if err == nil {
		t.Error("logs command should fail with invalid output format")
	}
}

func TestLogsCmd_TopNFlag(t *testing.T) {
	// Create a temporary log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	logContent := `192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
192.168.1.2 - - [01/Jan/2024:12:00:01 +0000] "POST /api/login HTTP/1.1" 401 567 "-" "curl/7.68.0"
192.168.1.1 - - [01/Jan/2024:12:00:02 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"`

	err := os.WriteFile(logFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Test with custom top N
	rootCmd.SetArgs([]string{"logs", logFile, "--format", "combined", "--top", "5", "--output", "json"})
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("logs command failed: %v", err)
	}
}
