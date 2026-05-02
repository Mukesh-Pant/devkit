package logparser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Combined log format regex
// Format: IP - - [timestamp] "METHOD path HTTP/1.1" status size "referrer" "user-agent"
// Example: 192.168.1.1 - - [01/Jan/2024:12:00:00 +0000] "GET /api/users HTTP/1.1" 200 1234 "-" "Mozilla/5.0"
var combinedLogRegex = regexp.MustCompile(
	`^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) \S+" (\d+) (\d+|-) "([^"]*)" "([^"]*)"`,
)

// parseCombined parses Apache/Nginx combined log format
func parseCombined(line string) (*LogEntry, error) {
	matches := combinedLogRegex.FindStringSubmatch(line)
	if matches == nil || len(matches) < 9 {
		return nil, fmt.Errorf("invalid combined log format")
	}

	entry := &LogEntry{
		IP:        matches[1],
		Method:    matches[3],
		Path:      matches[4],
		Referrer:  matches[7],
		UserAgent: matches[8],
	}

	// Parse timestamp
	// Format: 01/Jan/2024:12:00:00 +0000
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		// If timestamp parsing fails, use zero time but continue
		entry.Timestamp = time.Time{}
	} else {
		entry.Timestamp = timestamp
	}

	// Parse status code
	statusCode, err := strconv.Atoi(matches[5])
	if err != nil {
		return nil, fmt.Errorf("invalid status code: %s", matches[5])
	}
	entry.StatusCode = statusCode

	// Parse size (may be "-" for no content)
	if matches[6] != "-" {
		size, err := strconv.Atoi(matches[6])
		if err != nil {
			// If size parsing fails, set to 0 but continue
			entry.Size = 0
		} else {
			entry.Size = size
		}
	}

	return entry, nil
}
