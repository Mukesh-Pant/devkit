package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"devkit/internal/checker"
)

// TestJSONFormatter_Format_CheckResults verifies that JSONFormatter correctly
// formats []checker.CheckResult as valid JSON
func TestJSONFormatter_Format_CheckResults(t *testing.T) {
	formatter := &JSONFormatter{}
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

	// Verify output is valid JSON
	var jsonResults []checkResultJSON
	err = json.Unmarshal(buf.Bytes(), &jsonResults)
	if err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %s", err, buf.String())
	}

	// Verify first result
	if jsonResults[0].URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got '%s'", jsonResults[0].URL)
	}
	if jsonResults[0].StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", jsonResults[0].StatusCode)
	}
	if jsonResults[0].ResponseTime != 123 {
		t.Errorf("Expected response time 123ms, got %d", jsonResults[0].ResponseTime)
	}
	if jsonResults[0].Status != "UP" {
		t.Errorf("Expected status 'UP', got '%s'", jsonResults[0].Status)
	}
	if jsonResults[0].Error != "" {
		t.Errorf("Expected no error, got '%s'", jsonResults[0].Error)
	}

	// Verify second result with error
	if jsonResults[1].URL != "https://example.com/api" {
		t.Errorf("Expected URL 'https://example.com/api', got '%s'", jsonResults[1].URL)
	}
	if jsonResults[1].StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", jsonResults[1].StatusCode)
	}
	if jsonResults[1].Status != "DOWN" {
		t.Errorf("Expected status 'DOWN', got '%s'", jsonResults[1].Status)
	}
	if jsonResults[1].Error != "unexpected status code: got 500, expected 200" {
		t.Errorf("Expected error message, got '%s'", jsonResults[1].Error)
	}
}

// TestJSONFormatter_Format_EmptyResults verifies that JSONFormatter handles
// empty result slices correctly
func TestJSONFormatter_Format_EmptyResults(t *testing.T) {
	formatter := &JSONFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify output is valid JSON array
	var jsonResults []checkResultJSON
	err = json.Unmarshal(buf.Bytes(), &jsonResults)
	if err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	if len(jsonResults) != 0 {
		t.Errorf("Expected empty array, got %d elements", len(jsonResults))
	}
}

// TestJSONFormatter_Format_ProperEscaping verifies that JSONFormatter
// properly escapes special characters in URLs and error messages
func TestJSONFormatter_Format_ProperEscaping(t *testing.T) {
	formatter := &JSONFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{
		{
			URL:          "https://example.com/path?query=\"value\"&key=<script>",
			StatusCode:   0,
			ResponseTime: 0,
			Status:       checker.StatusDown,
			Error:        errors.New("error with \"quotes\" and <tags>"),
		},
	}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify output is valid JSON (proper escaping)
	var jsonResults []checkResultJSON
	err = json.Unmarshal(buf.Bytes(), &jsonResults)
	if err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify special characters are preserved after unmarshaling
	if jsonResults[0].URL != "https://example.com/path?query=\"value\"&key=<script>" {
		t.Errorf("URL not properly escaped/unescaped: %s", jsonResults[0].URL)
	}
	if jsonResults[0].Error != "error with \"quotes\" and <tags>" {
		t.Errorf("Error message not properly escaped/unescaped: %s", jsonResults[0].Error)
	}
}

// TestJSONFormatter_Format_Indentation verifies that JSONFormatter
// produces pretty-printed JSON with proper indentation
func TestJSONFormatter_Format_Indentation(t *testing.T) {
	formatter := &JSONFormatter{}
	var buf bytes.Buffer

	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   200,
			ResponseTime: 100 * time.Millisecond,
			Status:       checker.StatusUp,
		},
	}

	err := formatter.Format(results, &buf)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()

	// Verify output contains indentation (2 spaces)
	if !bytes.Contains([]byte(output), []byte("  \"url\"")) {
		t.Error("Expected indented JSON output with 2 spaces")
	}

	// Verify output ends with newline
	if output[len(output)-1] != '\n' {
		t.Error("Expected output to end with newline")
	}
}

// TestJSONFormatter_Format_UnsupportedType verifies that JSONFormatter
// handles unsupported data types gracefully
func TestJSONFormatter_Format_UnsupportedType(t *testing.T) {
	formatter := &JSONFormatter{}
	var buf bytes.Buffer

	// Test with a type that has unexported fields or channels (cannot be marshaled)
	unsupportedData := make(chan int)

	err := formatter.Format(unsupportedData, &buf)
	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}
}
