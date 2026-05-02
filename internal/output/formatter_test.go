package output

import (
	"bytes"
	"testing"
)

// TestNewFormatter_Table verifies that NewFormatter returns a TableFormatter
// when Format is set to FormatTable
func TestNewFormatter_Table(t *testing.T) {
	opts := Options{
		Format:  FormatTable,
		NoColor: false,
	}

	formatter := NewFormatter(opts)

	if _, ok := formatter.(*TableFormatter); !ok {
		t.Errorf("Expected TableFormatter, got %T", formatter)
	}
}

// TestNewFormatter_JSON verifies that NewFormatter returns a JSONFormatter
// when Format is set to FormatJSON
func TestNewFormatter_JSON(t *testing.T) {
	opts := Options{
		Format: FormatJSON,
	}

	formatter := NewFormatter(opts)

	if _, ok := formatter.(*JSONFormatter); !ok {
		t.Errorf("Expected JSONFormatter, got %T", formatter)
	}
}

// TestNewFormatter_Plain verifies that NewFormatter returns a PlainFormatter
// when Format is set to FormatPlain
func TestNewFormatter_Plain(t *testing.T) {
	opts := Options{
		Format: FormatPlain,
	}

	formatter := NewFormatter(opts)

	if _, ok := formatter.(*PlainFormatter); !ok {
		t.Errorf("Expected PlainFormatter, got %T", formatter)
	}
}

// TestNewFormatter_Default verifies that NewFormatter returns a TableFormatter
// when an invalid or empty Format is provided (default behavior)
func TestNewFormatter_Default(t *testing.T) {
	opts := Options{
		Format: Format("invalid"),
	}

	formatter := NewFormatter(opts)

	if _, ok := formatter.(*TableFormatter); !ok {
		t.Errorf("Expected TableFormatter as default, got %T", formatter)
	}
}

// TestNewFormatter_TableWithNoColor verifies that the NoColor option
// is properly passed to the TableFormatter
func TestNewFormatter_TableWithNoColor(t *testing.T) {
	opts := Options{
		Format:  FormatTable,
		NoColor: true,
	}

	formatter := NewFormatter(opts)

	tableFormatter, ok := formatter.(*TableFormatter)
	if !ok {
		t.Fatalf("Expected TableFormatter, got %T", formatter)
	}

	if !tableFormatter.noColor {
		t.Error("Expected noColor to be true, got false")
	}
}

// TestTableFormatter_Format verifies that TableFormatter implements
// the Formatter interface
func TestTableFormatter_Format(t *testing.T) {
	formatter := &TableFormatter{noColor: false}
	var buf bytes.Buffer

	// Test with nil data - should return error for unsupported type
	err := formatter.Format(nil, &buf)

	if err == nil {
		t.Error("Expected error for nil data, got nil")
	}
}

// TestJSONFormatter_Format verifies that JSONFormatter implements
// the Formatter interface and handles unsupported types
func TestJSONFormatter_Format(t *testing.T) {
	formatter := &JSONFormatter{}
	var buf bytes.Buffer

	// Test with unsupported type (channel cannot be marshaled to JSON)
	unsupportedData := make(chan int)
	err := formatter.Format(unsupportedData, &buf)

	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}
}

// TestPlainFormatter_Format verifies that PlainFormatter implements
// the Formatter interface and handles unsupported types
func TestPlainFormatter_Format(t *testing.T) {
	formatter := &PlainFormatter{}
	var buf bytes.Buffer

	// Test with unsupported type
	err := formatter.Format("string", &buf)

	if err == nil {
		t.Error("Expected error for unsupported type, got nil")
	}
}

// TestFormatConstants verifies that format constants have expected values
func TestFormatConstants(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected string
	}{
		{"Table format", FormatTable, "table"},
		{"JSON format", FormatJSON, "json"},
		{"Plain format", FormatPlain, "plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.format) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.format))
			}
		})
	}
}
