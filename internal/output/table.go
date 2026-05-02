package output

import (
	"fmt"
	"io"

	"devkit/internal/checker"
	"devkit/internal/logparser"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// TableFormatter formats data as tables
type TableFormatter struct {
	noColor bool
}

// Format implements the Formatter interface for table output
// Supports formatting of []checker.CheckResult with columns: URL, STATUS, RESPONSE TIME, CODE
func (f *TableFormatter) Format(data interface{}, w io.Writer) error {
	switch v := data.(type) {
	case []checker.CheckResult:
		return f.formatCheckResults(v, w)
	case *logparser.Statistics:
		return f.formatStatistics(v, w)
	default:
		return fmt.Errorf("unsupported data type for table formatting: %T", data)
	}
}

// formatCheckResults formats health check results as a table
func (f *TableFormatter) formatCheckResults(results []checker.CheckResult, w io.Writer) error {
	table := tablewriter.NewWriter(w)
	
	// Set header
	table.Header("URL", "STATUS", "RESPONSE TIME", "CODE")

	// Append rows
	for _, result := range results {
		// Format status with color
		statusStr := f.formatStatus(result.Status)

		// Format response time
		responseTimeStr := fmt.Sprintf("%dms", result.ResponseTime.Milliseconds())

		// Format status code
		codeStr := fmt.Sprintf("%d", result.StatusCode)
		if result.StatusCode == 0 {
			codeStr = "-"
		}

		err := table.Append(result.URL, statusStr, responseTimeStr, codeStr)
		if err != nil {
			return fmt.Errorf("failed to append row to table: %w", err)
		}
	}

	// Render the table
	err := table.Render()
	if err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}

	return nil
}

// formatStatus formats the status with color (green for UP, red for DOWN)
func (f *TableFormatter) formatStatus(status checker.Status) string {
	if f.noColor {
		return string(status)
	}

	switch status {
	case checker.StatusUp:
		return color.GreenString(string(status))
	case checker.StatusDown:
		return color.RedString(string(status))
	default:
		return string(status)
	}
}

// formatStatistics formats log statistics as tables
func (f *TableFormatter) formatStatistics(stats *logparser.Statistics, w io.Writer) error {
	// Print summary statistics
	fmt.Fprintf(w, "\n=== Log Statistics ===\n\n")
	fmt.Fprintf(w, "Total Lines: %d\n", stats.TotalLines)
	fmt.Fprintf(w, "Error Rate: %.2f%%\n\n", stats.ErrorRate)

	// Print top IPs
	if len(stats.TopIPs) > 0 {
		fmt.Fprintf(w, "Top IP Addresses:\n")
		table := tablewriter.NewWriter(w)
		table.Header("IP ADDRESS", "COUNT")
		for _, item := range stats.TopIPs {
			err := table.Append(item.Value, fmt.Sprintf("%d", item.Count))
			if err != nil {
				return fmt.Errorf("failed to append row to table: %w", err)
			}
		}
		err := table.Render()
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
		fmt.Fprintln(w)
	}

	// Print top paths
	if len(stats.TopPaths) > 0 {
		fmt.Fprintf(w, "Top Request Paths:\n")
		table := tablewriter.NewWriter(w)
		table.Header("PATH", "COUNT")
		for _, item := range stats.TopPaths {
			err := table.Append(item.Value, fmt.Sprintf("%d", item.Count))
			if err != nil {
				return fmt.Errorf("failed to append row to table: %w", err)
			}
		}
		err := table.Render()
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
		fmt.Fprintln(w)
	}

	// Print top status codes
	if len(stats.TopStatusCodes) > 0 {
		fmt.Fprintf(w, "Top Status Codes:\n")
		table := tablewriter.NewWriter(w)
		table.Header("STATUS CODE", "COUNT")
		for _, item := range stats.TopStatusCodes {
			err := table.Append(item.Value, fmt.Sprintf("%d", item.Count))
			if err != nil {
				return fmt.Errorf("failed to append row to table: %w", err)
			}
		}
		err := table.Render()
		if err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
	}

	return nil
}
