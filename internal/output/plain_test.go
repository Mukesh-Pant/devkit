package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"devkit/internal/checker"
)

// TestPlainFormatter_Format_CheckResults verifies that PlainFormatter correctly
// formats []checker.CheckResult as plain text
func TestPlainFormatter_Format_CheckResults(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   200,
			ResponseTime: 123 * time.Millisecond,
			Status:       checker.StatusUp,
			Error:        nil,
		},
		{
			URL:          "https://example.com/api",
			StatusCode:   500,
			ResponseTime: 5000 * time.Millisecond,
			Status:       checker.StatusDown,
			Error:        errors.New("unexpected status code: got 500, expected 200"),
		},
	}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()

	// Verify first result is in output
	if !strings.Contains(output, "URL: https://example.com") {
		t.Error("Expected first URL in output")
	}
	if !strings.Contains(output, "Status: UP") {
		t.Error("Expected 'Status: UP' in output")
	}
	if !strings.Contains(output, "Code: 200") {
		t.Error("Expected 'Code: 200' in output")
	}
	if !strings.Contains(output, "Time: 123ms") {
		t.Error("Expected 'Time: 123ms' in output")
	}

	// Verify second result is in output
	if !strings.Contains(output, "URL: https://example.com/api") {
		t.Error("Expected second URL in output")
	}
	if !strings.Contains(output, "Status: DOWN") {
		t.Error("Expected 'Status: DOWN' in output")
	}
	if !strings.Contains(output, "Code: 500") {
		t.Error("Expected 'Code: 500' in output")
	}
	if !strings.Contains(output, "Time: 5000ms") {
		t.Error("Expected 'Time: 5000ms' in output")
	}

	// Verify error message is in output
	if !strings.Contains(output, "Error: unexpected status code: got 500, expected 200") {
		t.Error("Expected error message in output")
	}
}

// TestPlainFormatter_Format_EmptyResults verifies that PlainFormatter handles
// empty result slices correctly
func TestPlainFormatter_Format_EmptyResults(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()

	if output != "" {
		t.Errorf("Expected empty output, got: %s", output)
	}
}

// TestPlainFormatter_Format_NoError verifies that PlainFormatter
// does not print error line when Error is nil
func TestPlainFormatter_Format_NoError(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   200,
			ResponseTime: 100 * time.Millisecond,
			Status:       checker.StatusUp,
			Error:        nil,
		},
	}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()

	// Verify no error line is present
	if strings.Contains(output, "Error:") {
		t.Error("Expected no error line in output when Error is nil")
	}

	// Verify output contains expected fields
	if !strings.Contains(output, "URL: https://example.com") {
		t.Error("Expected URL in output")
	}
}

// TestPlainFormatter_Format_ZeroStatusCode verifies that PlainFormatter
// displays "-" for zero status codes
func TestPlainFormatter_Format_ZeroStatusCode(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   0,
			ResponseTime: 0,
			Status:       checker.StatusDown,
			Error:        errors.New("connection refused"),
		},
	}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()

	// Verify status code is displayed as "-"
	if !strings.Contains(output, "Code: -") {
		t.Error("Expected 'Code: -' for zero status code")
	}
}

// TestPlainFormatter_Format_MultipleResults verifies that PlainFormatter
// formats multiple results with proper line breaks
func TestPlainFormatter_Format_MultipleResults(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{
		{
			URL:          "https://example1.com",
			StatusCode:   200,
			ResponseTime: 100 * time.Millisecond,
			Status:       checker.StatusUp,
		},
		{
			URL:          "https://example2.com",
			StatusCode:   200,
			ResponseTime: 150 * time.Millisecond,
			Status:       checker.StatusUp,
		},
		{
			URL:          "https://example3.com",
			StatusCode:   404,
			ResponseTime: 200 * time.Millisecond,
			Status:       checker.StatusDown,
		},
	}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Verify we have 3 lines (one per result)
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}

	// Verify each line contains expected URL
	if !strings.Contains(lines[0], "example1.com") {
		t.Error("Expected first line to contain example1.com")
	}
	if !strings.Contains(lines[1], "example2.com") {
		t.Error("Expected second line to contain example2.com")
	}
	if !strings.Contains(lines[2], "example3.com") {
		t.Error("Expected third line to contain example3.com")
	}
}

// TestPlainFormatter_Format_UnsupportedType verifies that PlainFormatter
// returns an error for unsupported data types
func TestPlainFormatter_Format_UnsupportedType(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	// Test with unsupported type
	unsupportedData := "string data"

	err := formatter.Format(unsupportedData, &buf)
	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}

	expectedError := "unsupported data type for plain text formatting"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedError, err.Error())
	}
}
