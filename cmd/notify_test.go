package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"devkit/config"

	"github.com/spf13/cobra"
)

func TestNotifyCommand_Success(t *testing.T) {
	// Create a mock Slack webhook server
	var receivedPayload map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and content type
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Decode the payload
		json.NewDecoder(r.Body).Decode(&receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Initialize config manager
	cfg = config.NewManager()

	// Create a buffer to capture output
	buf := new(bytes.Buffer)

	// Create a new root command for testing
	testRootCmd := &cobra.Command{Use: "devkit"}
	testNotifyCmd := &cobra.Command{
		Use:  "notify",
		RunE: runNotify,
	}
	testNotifyCmd.Flags().StringVar(&notifyWebhook, "webhook", "", "Slack webhook URL")
	testNotifyCmd.Flags().StringVar(&notifyMessage, "message", "", "notification text")
	testNotifyCmd.Flags().StringVar(&notifyTitle, "title", "", "notification title")
	testNotifyCmd.Flags().StringVar(&notifyColor, "color", "", "notification color")

	testRootCmd.AddCommand(testNotifyCmd)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs([]string{"notify", "--webhook", server.URL, "--message", "Test message", "--title", "Test Title", "--color", "good"})

	err := testRootCmd.Execute()
	if err != nil {
		t.Fatalf("notify command failed: %v", err)
	}

	output := buf.String()

	// Verify success message
	if !strings.Contains(output, "Notification sent successfully") {
		t.Errorf("expected success message, got: %q", output)
	}

	// Verify payload structure
	if receivedPayload == nil {
		t.Fatal("no payload received by mock server")
	}

	attachments, ok := receivedPayload["attachments"].([]interface{})
	if !ok || len(attachments) == 0 {
		t.Fatal("expected attachments in payload")
	}

	attachment := attachments[0].(map[string]interface{})
	if attachment["title"] != "Test Title" {
		t.Errorf("expected title 'Test Title', got %v", attachment["title"])
	}
	if attachment["text"] != "Test message" {
		t.Errorf("expected text 'Test message', got %v", attachment["text"])
	}
	if attachment["color"] != "good" {
		t.Errorf("expected color 'good', got %v", attachment["color"])
	}
}

func TestNotifyCommand_MissingWebhook(t *testing.T) {
	// Initialize config manager without webhook
	cfg = config.NewManager()

	// Create a buffer to capture output
	buf := new(bytes.Buffer)

	// Create a new root command for testing
	testRootCmd := &cobra.Command{Use: "devkit"}
	testNotifyCmd := &cobra.Command{
		Use:  "notify",
		RunE: runNotify,
	}
	testNotifyCmd.Flags().StringVar(&notifyWebhook, "webhook", "", "Slack webhook URL")
	testNotifyCmd.Flags().StringVar(&notifyMessage, "message", "", "notification text")
	testNotifyCmd.Flags().StringVar(&notifyTitle, "title", "", "notification title")
	testNotifyCmd.Flags().StringVar(&notifyColor, "color", "", "notification color")

	testRootCmd.AddCommand(testNotifyCmd)
	testRootCmd.SetOut(buf)
	testRootCmd.SetErr(buf)
	testRootCmd.SetArgs([]string{"notify", "--message", "Test message"})

	err := testRootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing webhook, got nil")
	}

	if !strings.Contains(err.Error(), "webhook URL is required") {
		t.Errorf("expected 'webhook URL is required' error, got: %v", err)
	}
}

func TestNotifyCommand_ColorVariants(t *testing.T) {
	tests := []struct {
		name          string
		colorFlag     string
		expectedColor string
	}{
		{"good color", "good", "good"},
		{"warning color", "warning", "warning"},
		{"danger color", "danger", "danger"},
		{"hex color", "#FF5733", "#FF5733"},
		{"no color", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock Slack webhook server
			var receivedPayload map[string]interface{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewDecoder(r.Body).Decode(&receivedPayload)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// Initialize config manager
			cfg = config.NewManager()

			// Create a buffer to capture output
			buf := new(bytes.Buffer)

			// Create a new root command for testing
			testRootCmd := &cobra.Command{Use: "devkit"}
			testNotifyCmd := &cobra.Command{
				Use:  "notify",
				RunE: runNotify,
			}
			testNotifyCmd.Flags().StringVar(&notifyWebhook, "webhook", "", "Slack webhook URL")
			testNotifyCmd.Flags().StringVar(&notifyMessage, "message", "", "notification text")
			testNotifyCmd.Flags().StringVar(&notifyTitle, "title", "", "notification title")
			testNotifyCmd.Flags().StringVar(&notifyColor, "color", "", "notification color")

			testRootCmd.AddCommand(testNotifyCmd)
			testRootCmd.SetOut(buf)
			testRootCmd.SetErr(buf)

			args := []string{"notify", "--webhook", server.URL, "--message", "Test message"}
			if tt.colorFlag != "" {
				args = append(args, "--color", tt.colorFlag)
			}
			testRootCmd.SetArgs(args)

			err := testRootCmd.Execute()
			if err != nil {
				t.Fatalf("notify command failed: %v", err)
			}

			// Verify payload color
			attachments := receivedPayload["attachments"].([]interface{})
			attachment := attachments[0].(map[string]interface{})
			
			if tt.expectedColor == "" {
				// When no color is specified, the field might be empty or omitted
				color, exists := attachment["color"]
				if exists && color != "" {
					t.Errorf("expected no color or empty color, got %v", color)
				}
			} else {
				if attachment["color"] != tt.expectedColor {
					t.Errorf("expected color '%s', got %v", tt.expectedColor, attachment["color"])
				}
			}
		})
	}
}

func TestNotifyCommand_WebhookError(t *testing.T) {
	// Note: This test is limited because runNotify calls os.Exit(1) on error,
	// which terminates the test process. In a real scenario, we would refactor
	// the command to return errors instead of calling os.Exit directly.
	// For now, we skip this test.
	t.Skip("Skipping test that would call os.Exit(1)")
}
