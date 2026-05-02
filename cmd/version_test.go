package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	// Set test version information
	SetVersionInfo("1.2.3", "abc123def", "2024-01-15T10:00:00Z")

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	
	// Create a new root command for testing to avoid state pollution
	testRootCmd := &cobra.Command{Use: "devkit"}
	testVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("devkit version %s\n", version)
			cmd.Printf("commit: %s\n", commit)
			cmd.Printf("built: %s\n", buildDate)
		},
	}
	testRootCmd.AddCommand(testVersionCmd)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs([]string{"version"})
	
	err := testRootCmd.Execute()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	output := buf.String()

	// Verify output contains version information
	if !strings.Contains(output, "devkit version 1.2.3") {
		t.Errorf("expected output to contain 'devkit version 1.2.3', got: %s", output)
	}

	if !strings.Contains(output, "commit: abc123def") {
		t.Errorf("expected output to contain 'commit: abc123def', got: %s", output)
	}

	if !strings.Contains(output, "built: 2024-01-15T10:00:00Z") {
		t.Errorf("expected output to contain 'built: 2024-01-15T10:00:00Z', got: %s", output)
	}
}

func TestVersionCommand_DefaultValues(t *testing.T) {
	// Set default version information (simulating no ldflags injection)
	SetVersionInfo("dev", "unknown", "unknown")

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	
	// Create a new root command for testing to avoid state pollution
	testRootCmd := &cobra.Command{Use: "devkit"}
	testVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("devkit version %s\n", version)
			cmd.Printf("commit: %s\n", commit)
			cmd.Printf("built: %s\n", buildDate)
		},
	}
	testRootCmd.AddCommand(testVersionCmd)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs([]string{"version"})
	
	err := testRootCmd.Execute()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	output := buf.String()

	// Verify output contains default values
	if !strings.Contains(output, "devkit version dev") {
		t.Errorf("expected output to contain 'devkit version dev', got: %s", output)
	}

	if !strings.Contains(output, "commit: unknown") {
		t.Errorf("expected output to contain 'commit: unknown', got: %s", output)
	}

	if !strings.Contains(output, "built: unknown") {
		t.Errorf("expected output to contain 'built: unknown', got: %s", output)
	}
}

func TestSetVersionInfo(t *testing.T) {
	// Test that SetVersionInfo correctly updates the package variables
	testVersion := "2.0.0"
	testCommit := "xyz789"
	testBuildDate := "2024-02-20T15:30:00Z"

	SetVersionInfo(testVersion, testCommit, testBuildDate)

	if version != testVersion {
		t.Errorf("expected version to be %s, got %s", testVersion, version)
	}

	if commit != testCommit {
		t.Errorf("expected commit to be %s, got %s", testCommit, commit)
	}

	if buildDate != testBuildDate {
		t.Errorf("expected buildDate to be %s, got %s", testBuildDate, buildDate)
	}
}
