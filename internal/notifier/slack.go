package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Color represents Slack attachment colors
type Color string

const (
	ColorGood    Color = "good"    // Green
	ColorWarning Color = "warning" // Yellow
	ColorDanger  Color = "danger"  // Red
)

// Message represents a Slack notification message
type Message struct {
	Title  string
	Text   string
	Color  Color
	Fields []Field // Optional additional fields
}

// Field represents a Slack attachment field
type Field struct {
	Title string
	Value string
	Short bool
}

// slackPayload represents the JSON payload sent to Slack
type slackPayload struct {
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Title     string  `json:"title,omitempty"`
	Text      string  `json:"text"`
	Color     string  `json:"color,omitempty"`
	Timestamp int64   `json:"ts"`
	Fields    []field `json:"fields,omitempty"`
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Notifier sends Slack notifications
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// NewNotifier creates a new Slack notifier
func NewNotifier(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewNotifierWithClient creates a notifier with a custom HTTP client (for testing)
func NewNotifierWithClient(webhookURL string, client *http.Client) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     client,
	}
}

// Send sends a notification to Slack
func (n *Notifier) Send(msg Message) error {
	// Validate webhook URL format before sending
	if err := n.validateWebhookURL(); err != nil {
		return err
	}

	// Build the payload
	payload := slackPayload{
		Attachments: []attachment{
			{
				Title:     msg.Title,
				Text:      msg.Text,
				Color:     string(msg.Color),
				Timestamp: time.Now().Unix(),
				Fields:    convertFields(msg.Fields),
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	// Create HTTP POST request
	req, err := http.NewRequest("POST", n.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("network failure sending notification: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Parse Slack API error response
		var errorMsg bytes.Buffer
		errorMsg.ReadFrom(resp.Body)
		return fmt.Errorf("slack API error (status %d): %s", resp.StatusCode, errorMsg.String())
	}

	return nil
}

// SendSimple sends a simple text notification
func (n *Notifier) SendSimple(text string) error {
	return n.Send(Message{
		Text: text,
	})
}

// validateWebhookURL validates the webhook URL format
func (n *Notifier) validateWebhookURL() error {
	if n.webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	// Parse the URL
	parsedURL, err := url.Parse(n.webhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook URL format: %w", err)
	}

	// Validate scheme
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return fmt.Errorf("webhook URL must use http or https scheme")
	}

	// Validate host
	if parsedURL.Host == "" {
		return fmt.Errorf("webhook URL must have a valid host")
	}

	// Check if it looks like a Slack webhook (optional but helpful)
	if !strings.Contains(parsedURL.Host, "slack.com") && !strings.Contains(parsedURL.Host, "hooks.slack.com") {
		// This is just a warning - we'll still allow it for testing purposes
		// In production, you might want to be stricter
	}

	return nil
}

// convertFields converts Message.Fields to attachment fields
func convertFields(fields []Field) []field {
	if len(fields) == 0 {
		return nil
	}

	result := make([]field, len(fields))
	for i, f := range fields {
		result[i] = field{
			Title: f.Title,
			Value: f.Value,
			Short: f.Short,
		}
	}
	return result
}
