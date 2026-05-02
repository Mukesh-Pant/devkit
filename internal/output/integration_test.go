package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"devkit/internal/checker"
)

// TestFormatters_Integration verifies that all three formatters
// (Table, JSON, Plain) can format the same data correctly
func TestFormatters_Integration(t *testing.T) {
	// Create sample check results
	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   200,
			ResponseTime: 123 * time.Millisecond,
			Status:       checker.StatusUp,
		},
		{
			URL:          "https://api.example.com",
			StatusCode:   500,
			ResponseTime: 5000 * time.Millisecond,
			Status:       checker.StatusDown,
		},
	}

	t.Run("Table formatter", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewFormatter(Options{Format: FormatTable, NoColor: true})

		err := formatter.Format(results, &buf)
		if err != nil {
			t.Fatalf("Table formatter failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "example.com") {
			t.Error("Table output missing URL")
		}
		if !strings.Contains(output, "UP") {
			t.Error("Table output missing status")
		}
	})

	t.Run("JSON formatter", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewFormatter(Options{Format: FormatJSON})

		err := formatter.Format(results, &buf)
		if err != nil {
			t.Fatalf("JSON formatter failed: %v", err)
		}

		// Verify valid JSON
		var jsonResults []checkResultJSON
		err = json.Unmarshal(buf.Bytes(), &jsonResults)
		if err != nil {
			t.Fatalf("Invalid JSON output: %v", err)
		}

		if len(jsonResults) != 2 {
			t.Errorf("Expected 2 results, got %d", len(jsonResults))
		}
	})

	t.Run("Plain formatter", func(t *testing.T) {
		var buf bytes.Buffer
		formatter := NewFormatter(Options{Format: FormatPlain})

		err := formatter.Format(results, &buf)
		if err != nil {
			t.Fatalf("Plain formatter failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "example.com") {
			t.Error("Plain output missing URL")
		}
		if !strings.Contains(output, "Status: UP") {
			t.Error("Plain output missing status")
		}
	})
}

// TestFormatters_ConsistentData verifies that all formatters
// produce output containing the same core data
func TestFormatters_ConsistentData(t *testing.T) {
	results := []checker.CheckResult{
		{
			URL:          "https://test.example.com",
			StatusCode:   404,
			ResponseTime: 250 * time.Millisecond,
			Status:       checker.StatusDown,
		},
	}

	// Table output
	var tableBuf bytes.Buffer
	tableFormatter := NewFormatter(Options{Format: FormatTable, NoColor: true})
	err := tableFormatter.Format(results, &tableBuf)
	if err != nil {
		t.Fatalf("Table formatter failed: %v", err)
	}

	// JSON output
	var jsonBuf bytes.Buffer
	jsonFormatter := NewFormatter(Options{Format: FormatJSON})
	err = jsonFormatter.Format(results, &jsonBuf)
	if err != nil {
		t.Fatalf("JSON formatter failed: %v", err)
	}

	// Plain output
	var plainBuf bytes.Buffer
	plainFormatter := NewFormatter(Options{Format: FormatPlain})
	err = plainFormatter.Format(results, &plainBuf)
	if err != nil {
		t.Fatalf("Plain formatter failed: %v", err)
	}

	// Verify all outputs contain the URL
	tableOutput := tableBuf.String()
	jsonOutput := jsonBuf.String()
	plainOutput := plainBuf.String()

	if !strings.Contains(tableOutput, "test.example.com") {
		t.Error("Table output missing URL")
	}
	if !strings.Contains(jsonOutput, "test.example.com") {
		t.Error("JSON output missing URL")
	}
	if !strings.Contains(plainOutput, "test.example.com") {
		t.Error("Plain output missing URL")
	}

	// Verify all outputs contain status code
	if !strings.Contains(tableOutput, "404") {
		t.Error("Table output missing status code")
	}
	if !strings.Contains(jsonOutput, "404") {
		t.Error("JSON output missing status code")
	}
	if !strings.Contains(plainOutput, "404") {
		t.Error("Plain output missing status code")
	}

	// Verify all outputs contain DOWN status
	if !strings.Contains(tableOutput, "DOWN") {
		t.Error("Table output missing DOWN status")
	}
	if !strings.Contains(jsonOutput, "DOWN") {
		t.Error("JSON output missing DOWN status")
	}
	if !strings.Contains(plainOutput, "DOWN") {
		t.Error("Plain output missing DOWN status")
	}
}
