package cmd

import (
	"fmt"
	"os"
	"time"

	"devkit/internal/checker"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Check command flags
	checkTimeout        time.Duration
	checkExpectedStatus int
	checkOutput         string
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check <url1> <url2> ...",
	Short: "Check the health of one or more URLs",
	Long: `Perform HTTP health checks on one or more URLs concurrently.

The check command sends HTTP GET requests to the specified URLs and reports
their status, response time, and HTTP status code. URLs are checked concurrently
for improved performance.

Examples:
  devkit check https://example.com
  devkit check https://api.example.com https://example.com/health
  devkit check --timeout 10s --expected-status 200 https://example.com
  devkit check --output json https://example.com`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Add flags
	checkCmd.Flags().DurationVar(&checkTimeout, "timeout", 0, "timeout for health checks (default from config or 5s)")
	checkCmd.Flags().IntVar(&checkExpectedStatus, "expected-status", 200, "expected HTTP status code")
	checkCmd.Flags().StringVar(&checkOutput, "output", "", "output format: table, json, plain (default from config or table)")
}

func runCheck(cmd *cobra.Command, args []string) error {
	urls := args

	// Get timeout from config if not specified via flag
	timeout := checkTimeout
	if timeout == 0 {
		timeout = cfg.GetDuration("devkit.timeout")
		if timeout == 0 {
			timeout = 5 * time.Second
		}
	}

	// Get output format from config if not specified via flag
	outputFormat := checkOutput
	if outputFormat == "" {
		outputFormat = cfg.GetString("devkit.output")
		if outputFormat == "" {
			outputFormat = "table"
		}
	}

	// Create checker with options
	chk := checker.NewChecker(checker.Options{
		Timeout:        timeout,
		ExpectedStatus: checkExpectedStatus,
	})

	// Perform health checks with progress bar for batch operations
	var results []checker.CheckResult
	if len(urls) > 1 {
		results = chk.CheckWithProgress(urls)
	} else {
		// For single URL, don't show progress bar
		results = chk.CheckMultiple(urls)
	}

	// Display results based on output format
	switch outputFormat {
	case "json":
		displayResultsJSON(results)
	case "plain":
		displayResultsPlain(results)
	default:
		displayResultsTable(results)
	}

	// Exit with code 1 if any URL is DOWN
	for _, result := range results {
		if result.Status == checker.StatusDown {
			os.Exit(1)
		}
	}

	return nil
}

func displayResultsTable(results []checker.CheckResult) {
	// Print header
	fmt.Printf("%-50s %-10s %-15s %-10s\n", "URL", "STATUS", "RESPONSE TIME", "CODE")
	fmt.Println("--------------------------------------------------------------------------------------------")

	// Print results
	for _, result := range results {
		statusStr := string(result.Status)
		if result.Status == checker.StatusUp {
			statusStr = color.GreenString(statusStr)
		} else {
			statusStr = color.RedString(statusStr)
		}

		responseTime := fmt.Sprintf("%dms", result.ResponseTime.Milliseconds())
		statusCode := fmt.Sprintf("%d", result.StatusCode)
		if result.StatusCode == 0 {
			statusCode = "N/A"
		}

		fmt.Printf("%-50s %-20s %-15s %-10s\n", result.URL, statusStr, responseTime, statusCode)
	}
}

func displayResultsJSON(results []checker.CheckResult) {
	fmt.Println("[")
	for i, result := range results {
		fmt.Printf("  {\n")
		fmt.Printf("    \"url\": \"%s\",\n", result.URL)
		fmt.Printf("    \"status\": \"%s\",\n", result.Status)
		fmt.Printf("    \"status_code\": %d,\n", result.StatusCode)
		fmt.Printf("    \"response_time_ms\": %d", result.ResponseTime.Milliseconds())
		if result.Error != nil {
			fmt.Printf(",\n    \"error\": \"%s\"\n", result.Error.Error())
		} else {
			fmt.Println()
		}
		if i < len(results)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")
}

func displayResultsPlain(results []checker.CheckResult) {
	for _, result := range results {
		statusCode := fmt.Sprintf("%d", result.StatusCode)
		if result.StatusCode == 0 {
			statusCode = "N/A"
		}
		fmt.Printf("%s: %s (%s, %dms)\n", result.URL, result.Status, statusCode, result.ResponseTime.Milliseconds())
		if result.Error != nil {
			fmt.Printf("  Error: %s\n", result.Error.Error())
		}
	}
}
