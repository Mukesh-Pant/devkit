package logparser

import (
	"regexp"
	"strconv"
)

// Plain text log parsing uses pattern matching to extract basic information
var (
	// Pattern to match IP addresses
	ipPattern = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	
	// Pattern to match HTTP status codes
	statusPattern = regexp.MustCompile(`\b([2-5]\d{2})\b`)
	
	// Pattern to match common HTTP methods
	methodPattern = regexp.MustCompile(`\b(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\b`)
	
	// Pattern to match paths (starting with /)
	pathPattern = regexp.MustCompile(`\s(/[^\s"]*)\s`)
)

// parsePlain performs basic line counting and pattern matching for plain text logs
func parsePlain(line string) (*LogEntry, error) {
	entry := &LogEntry{}

	// Extract IP address
	if ipMatch := ipPattern.FindString(line); ipMatch != "" {
		entry.IP = ipMatch
	}

	// Extract status code
	if statusMatch := statusPattern.FindStringSubmatch(line); len(statusMatch) > 1 {
		if statusCode, err := strconv.Atoi(statusMatch[1]); err == nil {
			entry.StatusCode = statusCode
		}
	}

	// Extract HTTP method
	if methodMatch := methodPattern.FindString(line); methodMatch != "" {
		entry.Method = methodMatch
	}

	// Extract path
	if pathMatch := pathPattern.FindStringSubmatch(line); len(pathMatch) > 1 {
		entry.Path = pathMatch[1]
	}

	// For plain format, we always return a valid entry (even if mostly empty)
	// This allows line counting to work correctly
	return entry, nil
}
