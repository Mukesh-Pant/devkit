package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devkit/internal/checker"
	"devkit/internal/notifier"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Watch command flags
	watchInterval    time.Duration
	watchAlertWebhook string
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Continuously monitor URLs and send alerts on status changes",
	Long: `Continuously monitor URLs at regular intervals and send Slack alerts when status changes.

The watch command polls URLs from the configuration file at the specified interval
and displays real-time status updates. When a URL transitions from UP to DOWN,
an alert is sent via Slack webhook if configured.

The command runs until interrupted with Ctrl+C (SIGINT).

Examples:
  devkit watch
  devkit watch --interval 60s
  devkit watch --interval 30s --alert-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL

Configuration:
  URLs to monitor should be specified in .devkit.yaml:
    watch:
      urls:
        - https://example.com
        - https://api.example.com/health`,
	RunE: runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Add flags
	watchCmd.Flags().DurationVar(&watchInterval, "interval", 0, "polling interval (default from config or 30s)")
	watchCmd.Flags().StringVar(&watchAlertWebhook, "alert-webhook", "", "Slack webhook URL for alerts (default from config)")
}

func runWatch(cmd *cobra.Command, args []string) error {
	// Get URLs from configuration
	urls := cfg.GetStringSlice("watch.urls")
	if len(urls) == 0 {
		return fmt.Errorf("no URLs configured for monitoring. Add URLs to .devkit.yaml under watch.urls")
	}

	// Get interval from config if not specified via flag
	interval := watchInterval
	if interval == 0 {
		interval = cfg.GetDuration("watch.interval")
		if interval == 0 {
			interval = 30 * time.Second
		}
	}

	// Get alert webhook from config if not specified via flag
	alertWebhook := watchAlertWebhook
	if alertWebhook == "" {
		alertWebhook = cfg.GetString("devkit.slack_webhook")
	}

	// Get timeout from config
	timeout := cfg.GetDuration("devkit.timeout")
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	// Create checker with options
	chk := checker.NewChecker(checker.Options{
		Timeout:        timeout,
		ExpectedStatus: 200,
	})

	// Create notifier if webhook is configured
	var slackNotifier *notifier.Notifier
	if alertWebhook != "" {
		slackNotifier = notifier.NewNotifier(alertWebhook)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)

	// Create ticker for polling
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Track previous status for change detection
	previousStatus := make(map[string]checker.Status)

	// Display initial message
	fmt.Printf("Starting continuous monitoring of %d URL(s) every %s\n", len(urls), interval)
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Perform initial check immediately
	performCheck(chk, urls, previousStatus, slackNotifier)

	// Main monitoring loop
	for {
		select {
		case <-ticker.C:
			// Perform periodic check
			performCheck(chk, urls, previousStatus, slackNotifier)

		case <-sigChan:
			// Graceful shutdown
			fmt.Println("\n\nShutting down gracefully...")
			return nil
		}
	}
}

// performCheck executes health checks, detects status changes, sends alerts, and displays results
func performCheck(chk *checker.Checker, urls []string, previousStatus map[string]checker.Status, slackNotifier *notifier.Notifier) {
	// Clear terminal and move cursor to top
	clearTerminal()

	// Display timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Last check: %s\n\n", timestamp)

	// Perform health checks
	results := chk.CheckMultiple(urls)

	// Detect status changes and send alerts
	for _, result := range results {
		prevStatus, exists := previousStatus[result.URL]
		
		// Check if status changed from UP to DOWN
		if exists && prevStatus == checker.StatusUp && result.Status == checker.StatusDown {
			// Send alert if notifier is configured
			if slackNotifier != nil {
				sendAlert(slackNotifier, result)
			}
		}

		// Update previous status
		previousStatus[result.URL] = result.Status
	}

	// Display results in table format
	displayWatchResults(results)
}

// clearTerminal clears the terminal screen using ANSI escape codes
func clearTerminal() {
	// ANSI escape code to clear screen and move cursor to top-left
	fmt.Print("\033[2J\033[H")
}

// displayWatchResults displays health check results in a table format with colors
func displayWatchResults(results []checker.CheckResult) {
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
		
		// Display error if present
		if result.Error != nil {
			fmt.Printf("  Error: %s\n", color.RedString(result.Error.Error()))
		}
	}
}

// sendAlert sends a Slack notification when a URL goes down
func sendAlert(slackNotifier *notifier.Notifier, result checker.CheckResult) {
	// Build alert message
	msg := notifier.Message{
		Title: "🚨 URL Health Alert",
		Text:  fmt.Sprintf("URL is DOWN: %s", result.URL),
		Color: notifier.ColorDanger,
		Fields: []notifier.Field{
			{
				Title: "URL",
				Value: result.URL,
				Short: false,
			},
			{
				Title: "Status Code",
				Value: fmt.Sprintf("%d", result.StatusCode),
				Short: true,
			},
			{
				Title: "Response Time",
				Value: fmt.Sprintf("%dms", result.ResponseTime.Milliseconds()),
				Short: true,
			},
		},
	}

	// Add error field if present
	if result.Error != nil {
		msg.Fields = append(msg.Fields, notifier.Field{
			Title: "Error",
			Value: result.Error.Error(),
			Short: false,
		})
	}

	// Send notification (ignore errors to avoid disrupting monitoring)
	if err := slackNotifier.Send(msg); err != nil {
		fmt.Printf("Warning: Failed to send Slack alert: %s\n", err)
	}
}
