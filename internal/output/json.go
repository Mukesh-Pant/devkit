package output

import (
	"encoding/json"
	"fmt"
	"io"

	"devkit/internal/checker"
)

// JSONFormatter formats data as JSON
type JSONFormatter struct{}

// Format implements the Formatter interface for JSON output
// Supports formatting of []checker.CheckResult and Statistics structs
// Produces valid JSON with proper escaping using encoding/json with MarshalIndent
func (f *JSONFormatter) Format(data interface{}, w io.Writer) error {
	// Convert data to JSON-serializable format if needed
	var jsonData []byte
	var err error

	switch v := data.(type) {
	case []checker.CheckResult:
		// Convert CheckResult to JSON-friendly format
		jsonResults := convertCheckResults(v)
		jsonData, err = json.MarshalIndent(jsonResults, "", "  ")
	default:
		// For other types, marshal directly
		jsonData, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Write JSON data to writer
	_, err = w.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON output: %w", err)
	}

	// Add trailing newline for better terminal output
	_, err = w.Write([]byte("\n"))
	if err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// checkResultJSON is a JSON-serializable representation of checker.CheckResult
// This struct ensures proper JSON field naming and handles error serialization
type checkResultJSON struct {
	URL          string `json:"url"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int64  `json:"response_time_ms"`
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
}

// convertCheckResults converts []checker.CheckResult to JSON-serializable format
func convertCheckResults(results []checker.CheckResult) []checkResultJSON {
	jsonResults := make([]checkResultJSON, len(results))
	for i, result := range results {
		jsonResults[i] = checkResultJSON{
			URL:          result.URL,
			StatusCode:   result.StatusCode,
			ResponseTime: result.ResponseTime.Milliseconds(),
			Status:       string(result.Status),
		}
		if result.Error != nil {
			jsonResults[i].Error = result.Error.Error()
		}
	}
	return jsonResults
}
