package logparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

// LogFormat represents the log file format
type LogFormat string

const (
	// FormatCombined represents Apache/Nginx combined log format
	FormatCombined LogFormat = "combined"
	// FormatJSON represents JSON log format
	FormatJSON LogFormat = "json"
	// FormatPlain represents plain text log format
	FormatPlain LogFormat = "plain"
)

// LogEntry represents a parsed log entry
type LogEntry struct {
	IP         string
	Timestamp  time.Time
	Method     string
	Path       string
	StatusCode int
	Size       int
	UserAgent  string
	Referrer   string
}

// Statistics represents aggregated log statistics
type Statistics struct {
	TotalLines     int            `json:"total_lines"`
	StatusCounts   map[int]int    `json:"status_counts"`
	IPCounts       map[string]int `json:"ip_counts"`
	PathCounts     map[string]int `json:"path_counts"`
	ErrorRate      float64        `json:"error_rate"`
	TopIPs         []TopItem      `json:"top_ips"`
	TopPaths       []TopItem      `json:"top_paths"`
	TopStatusCodes []TopItem      `json:"top_status_codes"`
}

// TopItem represents a top N item with count
type TopItem struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// Options configures log parsing behavior
type Options struct {
	Format LogFormat
	TopN   int
}

// Parser parses log files
type Parser struct {
	opts Options
}

// NewParser creates a new log parser with options
func NewParser(opts Options) *Parser {
	// Set default TopN if not specified
	if opts.TopN == 0 {
		opts.TopN = 10
	}

	return &Parser{
		opts: opts,
	}
}

// Parse parses a log file and returns statistics
func (p *Parser) Parse(filepath string) (*Statistics, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", filepath, err)
	}
	defer file.Close()

	// Parse from the file reader
	return p.ParseReader(file)
}

// ParseReader parses logs from an io.Reader
func (p *Parser) ParseReader(r io.Reader) (*Statistics, error) {
	// Initialize statistics
	stats := &Statistics{
		StatusCounts: make(map[int]int),
		IPCounts:     make(map[string]int),
		PathCounts:   make(map[string]int),
	}

	// Create scanner for line-by-line reading
	scanner := bufio.NewScanner(r)

	// Parse each line
	for scanner.Scan() {
		line := scanner.Text()
		stats.TotalLines++

		// Parse the line based on format
		var entry *LogEntry
		var err error

		switch p.opts.Format {
		case FormatCombined:
			entry, err = parseCombined(line)
		case FormatJSON:
			entry, err = parseJSON(line)
		case FormatPlain:
			entry, err = parsePlain(line)
		default:
			return nil, fmt.Errorf("unsupported log format: %s", p.opts.Format)
		}

		// Handle malformed lines gracefully (log warning, continue parsing)
		if err != nil {
			// Skip malformed lines silently
			continue
		}

		// Update statistics
		if entry.StatusCode > 0 {
			stats.StatusCounts[entry.StatusCode]++
		}
		if entry.IP != "" {
			stats.IPCounts[entry.IP]++
		}
		if entry.Path != "" {
			stats.PathCounts[entry.Path]++
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	// Calculate error rate as: (4xx + 5xx) / total * 100
	errorCount := 0
	for statusCode, count := range stats.StatusCounts {
		if statusCode >= 400 && statusCode < 600 {
			errorCount += count
		}
	}
	totalRequests := 0
	for _, count := range stats.StatusCounts {
		totalRequests += count
	}
	if totalRequests > 0 {
		stats.ErrorRate = float64(errorCount) / float64(totalRequests) * 100
	}

	// Extract top N items
	stats.TopIPs = extractTopN(stats.IPCounts, p.opts.TopN)
	stats.TopPaths = extractTopN(stats.PathCounts, p.opts.TopN)
	stats.TopStatusCodes = extractTopNInt(stats.StatusCounts, p.opts.TopN)

	return stats, nil
}

// extractTopN extracts top N items from a map by count
func extractTopN(counts map[string]int, n int) []TopItem {
	// Convert map to slice
	items := make([]TopItem, 0, len(counts))
	for value, count := range counts {
		items = append(items, TopItem{
			Value: value,
			Count: count,
		})
	}

	// Sort by count (descending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// Return top N
	if len(items) > n {
		return items[:n]
	}
	return items
}

// extractTopNInt extracts top N items from a map with int keys by count
func extractTopNInt(counts map[int]int, n int) []TopItem {
	// Convert map to slice
	items := make([]TopItem, 0, len(counts))
	for value, count := range counts {
		items = append(items, TopItem{
			Value: fmt.Sprintf("%d", value),
			Count: count,
		})
	}

	// Sort by count (descending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	// Return top N
	if len(items) > n {
		return items[:n]
	}
	return items
}
