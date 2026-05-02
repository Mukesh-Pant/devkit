package output_test

import (
	"os"
	"time"

	"devkit/internal/checker"
	"devkit/internal/output"
)

// ExampleTableFormatter demonstrates the table formatter with check results
func ExampleTableFormatter() {
	// Create sample check results
	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   200,
			ResponseTime: 123 * time.Millisecond,
			Status:       checker.StatusUp,
		},
		{
			URL:          "https://example.com/api",
			StatusCode:   500,
			ResponseTime: 5000 * time.Millisecond,
			Status:       checker.StatusDown,
		},
	}

	// Create formatter with colors disabled for consistent output
	formatter := output.NewFormatter(output.Options{
		Format:  output.FormatTable,
		NoColor: true,
	})

	// Format and output to stdout
	_ = formatter.Format(results, os.Stdout)

	// Output will be a table with URL, STATUS, RESPONSE TIME, and CODE columns
}
