package cmd

import (
	"fmt"
	"os"

	"devkit/internal/notifier"

	"github.com/spf13/cobra"
)

var (
	// Notify command flags
	notifyWebhook string
	notifyMessage string
	notifyTitle   string
	notifyColor   string
)

// notifyCmd represents the notify command
var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send a Slack notification",
	Long: `Send a formatted notification to a Slack channel via webhook.

The notify command sends a message to Slack using an incoming webhook URL.
The webhook URL can be provided via the --webhook flag or configured in
.devkit.yaml under devkit.slack_webhook.

The message is formatted as a Slack attachment with optional title, color,
and timestamp. Colors can be specified as "good" (green), "warning" (yellow),
"danger" (red), or as a hex color code (e.g., "#FF5733").

Examples:
  devkit notify --message "Deployment completed successfully" --color good
  devkit notify --webhook https://hooks.slack.com/... --message "Alert!" --title "Production Issue" --color danger
  devkit notify --message "Build finished" --title "CI/CD" --color warning`,
	RunE: runNotify,
}

func init() {
	rootCmd.AddCommand(notifyCmd)

	// Add flags
	notifyCmd.Flags().StringVar(&notifyWebhook, "webhook", "", "Slack webhook URL (required if not in config)")
	notifyCmd.Flags().StringVar(&notifyMessage, "message", "", "notification text (required)")
	notifyCmd.Flags().StringVar(&notifyTitle, "title", "", "notification title")
	notifyCmd.Flags().StringVar(&notifyColor, "color", "", "notification color: good, warning, danger, or hex color code")

	// Mark message as required
	notifyCmd.MarkFlagRequired("message")
}

func runNotify(cmd *cobra.Command, args []string) error {
	// Get webhook URL from flag or config
	webhookURL := notifyWebhook
	if webhookURL == "" {
		webhookURL = cfg.GetString("devkit.slack_webhook")
	}

	// Validate that webhook URL is provided
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is required: provide via --webhook flag or devkit.slack_webhook in config")
	}

	// Validate that message is provided (should be caught by MarkFlagRequired, but double-check)
	if notifyMessage == "" {
		return fmt.Errorf("message is required: provide via --message flag")
	}

	// Parse color
	var color notifier.Color
	switch notifyColor {
	case "good":
		color = notifier.ColorGood
	case "warning":
		color = notifier.ColorWarning
	case "danger":
		color = notifier.ColorDanger
	case "":
		// No color specified, leave empty
		color = ""
	default:
		// Assume it's a hex color code or custom color
		color = notifier.Color(notifyColor)
	}

	// Create notifier
	n := notifier.NewNotifier(webhookURL)

	// Build message
	msg := notifier.Message{
		Title: notifyTitle,
		Text:  notifyMessage,
		Color: color,
	}

	// Send notification
	err := n.Send(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to send notification: %v\n", err)
		os.Exit(1)
	}

	// Display confirmation message
	cmd.Println("✓ Notification sent successfully")

	return nil
}
