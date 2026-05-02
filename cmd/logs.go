package cmd

import (
	"fmt"
	"os"

	"devkit/internal/logparser"
	"devkit/internal/output"

	"github.com/spf13/cobra"
)

var (
	// Flags for logs command
	logsFormat string
	logsTop    int
	logsOutput string
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <filepath>",
	Short: "Parse and summarize log files",
	Long: `Parse and summarize log files in various formats.

Supported formats:
  - combined: Apache/Nginx combined log format
  - json: JSON log entries (one per line)
  - plain: Plain text logs with pattern matching

The logs command analyzes log files and provides statistics including:
  - Total lines processed
  - Request counts by status code
  - Error rate (percentage of 4xx and 5xx responses)
  - Top N IP addresses by request count
  - Top N request paths by request count
  - Top N status codes by occurrence

Examples:
  # Parse Apache/Nginx combined log
  devkit logs /var/log/nginx/access.log --format combined

  # Parse JSON logs and show top 20 items
  devkit logs /var/log/app.json --format json --top 20

  # Parse plain text logs with JSON output
  devkit logs /var/log/app.log --format plain --output json`,
	Args: cobra.ExactArgs(1),
	RunE: runLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Add flags
	logsCmd.Flags().StringVar(&logsFormat, "format", "combined", "Log format: combined, json, or plain")
	logsCmd.Flags().IntVar(&logsTop, "top", 10, "Number of top items to display")
	logsCmd.Flags().StringVar(&logsOutput, "output", "table", "Output format: table, json, or plain")
}

func runLogs(cmd *cobra.Command, args []string) error {
	filepath := args[0]

	// Validate format
	var format logparser.LogFormat
	switch logsFormat {
	case "combined":
		format = logparser.FormatCombined
	case "json":
		format = logparser.FormatJSON
	case "plain":
		format = logparser.FormatPlain
	default:
		return fmt.Errorf("invalid format: %s (must be combined, json, or plain)", logsFormat)
	}

	// Validate output format
	var outputFormat output.Format
	switch logsOutput {
	case "table":
		outputFormat = output.FormatTable
	case "json":
		outputFormat = output.FormatJSON
	case "plain":
		outputFormat = output.FormatPlain
	default:
		return fmt.Errorf("invalid output format: %s (must be table, json, or plain)", logsOutput)
	}

	// Create parser
	parser := logparser.NewParser(logparser.Options{
		Format: format,
		TopN:   logsTop,
	})

	// Parse log file
	stats, err := parser.Parse(filepath)
	if err != nil {
		return fmt.Errorf("failed to parse log file: %w", err)
	}

	// Create formatter
	formatter := output.NewFormatter(output.Options{
		Format:  outputFormat,
		NoColor: false, // TODO: Auto-detect redirected output
	})

	// Format and display statistics
	err = formatter.Format(stats, os.Stdout)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	return nil
}
