package output

import (
	"fmt"
	"io"

	"devkit/internal/checker"
	"devkit/internal/logparser"
)

// PlainFormatter formats data as plain text
type PlainFormatter struct{}

// Format implements the Formatter interface for plain text output
// Supports formatting of []checker.CheckResult and Statistics structs
func (f *PlainFormatter) Format(data interface{}, w io.Writer) error {
	switch v := data.(type) {
	case []checker.CheckResult:
		return f.formatCheckResults(v, w)
	case *logparser.Statistics:
		return f.formatStatistics(v, w)
	default:
		return fmt.Errorf("unsupported data type for plain text formatting: %T", data)
	}
}

// formatCheckResults formats health check results as plain text
func (f *PlainFormatter) formatCheckResults(results []checker.CheckResult, w io.Writer) error {
	for _, result := range results {
		// Format: URL: <url> | Status: <status> | Code: <code> | Time: <time>ms
		statusCode := fmt.Sprintf("%d", result.StatusCode)
		if result.StatusCode == 0 {
			statusCode = "-"
		}

		line := fmt.Sprintf("URL: %s | Status: %s | Code: %s | Time: %dms\n",
			result.URL,
			result.Status,
			statusCode,
			result.ResponseTime.Milliseconds(),
		)

		_, err := w.Write([]byte(line))
		if err != nil {
			return fmt.Errorf("failed to write plain text output: %w", err)
		}

		// If there's an error, print it on the next line
		if result.Error != nil {
			errorLine := fmt.Sprintf("  Error: %s\n", result.Error.Error())
			_, err := w.Write([]byte(errorLine))
			if err != nil {
				return fmt.Errorf("failed to write error line: %w", err)
			}
		}
	}

	return nil
}

// formatStatistics formats log statistics as plain text
func (f *PlainFormatter) formatStatistics(stats *logparser.Statistics, w io.Writer) error {
	// Print summary statistics
	fmt.Fprintf(w, "Total Lines: %d\n", stats.TotalLines)
	fmt.Fprintf(w, "Error Rate: %.2f%%\n\n", stats.ErrorRate)

	// Print top IPs
	if len(stats.TopIPs) > 0 {
		fmt.Fprintf(w, "Top IP Addresses:\n")
		for _, item := range stats.TopIPs {
			fmt.Fprintf(w, "  %s: %d\n", item.Value, item.Count)
		}
		fmt.Fprintln(w)
	}

	// Print top paths
	if len(stats.TopPaths) > 0 {
		fmt.Fprintf(w, "Top Request Paths:\n")
		for _, item := range stats.TopPaths {
			fmt.Fprintf(w, "  %s: %d\n", item.Value, item.Count)
		}
		fmt.Fprintln(w)
	}

	// Print top status codes
	if len(stats.TopStatusCodes) > 0 {
		fmt.Fprintf(w, "Top Status Codes:\n")
		for _, item := range stats.TopStatusCodes {
			fmt.Fprintf(w, "  %s: %d\n", item.Value, item.Count)
		}
	}

	return nil
}
