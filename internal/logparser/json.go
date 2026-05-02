package logparser

import (
	"encoding/json"
	"fmt"
	"time"
)

// jsonLogEntry represents a JSON log entry structure
// This is flexible to handle various JSON log formats
type jsonLogEntry struct {
	IP         string    `json:"ip"`
	Timestamp  time.Time `json:"timestamp"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Status     int       `json:"status"` // Alternative field name
	Size       int       `json:"size"`
	UserAgent  string    `json:"user_agent"`
	Referrer   string    `json:"referrer"`
}

// parseJSON parses JSON log entries
func parseJSON(line string) (*LogEntry, error) {
	var jsonEntry jsonLogEntry

	// Unmarshal JSON
	err := json.Unmarshal([]byte(line), &jsonEntry)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Convert to LogEntry
	entry := &LogEntry{
		IP:        jsonEntry.IP,
		Timestamp: jsonEntry.Timestamp,
		Method:    jsonEntry.Method,
		Path:      jsonEntry.Path,
		Size:      jsonEntry.Size,
		UserAgent: jsonEntry.UserAgent,
		Referrer:  jsonEntry.Referrer,
	}

	// Handle both status_code and status field names
	if jsonEntry.StatusCode > 0 {
		entry.StatusCode = jsonEntry.StatusCode
	} else if jsonEntry.Status > 0 {
		entry.StatusCode = jsonEntry.Status
	}

	return entry, nil
}
