package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewNotifier(t *testing.T) {
	webhookURL := "https://hooks.slack.com/services/TEST/WEBHOOK/URL"
	notifier := NewNotifier(webhookURL)

	if notifier.webhookURL != webhookURL {
		t.Errorf("expected webhook URL %s, got %s", webhookURL, notifier.webhookURL)
	}

	if notifier.client == nil {
		t.Error("expected client to be initialized")
	}

	if notifier.client.Timeout != 10*time.Second {
		t.Errorf("expected timeout 10s, got %v", notifier.client.Timeout)
	}
}

func TestNewNotifierWithClient(t *testing.T) {
	webhookURL := "https://hooks.slack.com/services/TEST/WEBHOOK/URL"
	customClient := &http.Client{Timeout: 5 * time.Second}
	notifier := NewNotifierWithClient(webhookURL, customClient)

	if notifier.webhookURL != webhookURL {
		t.Errorf("expected webhook URL %s, got %s", webhookURL, notifier.webhookURL)
	}

	if notifier.client != customClient {
		t.Error("expected custom client to be used")
	}

	if notifier.client.Timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", notifier.client.Timeout)
	}
}

func TestSend_Success(t *testing.T) {
	var receivedPayload slackPayload
	var receivedContentType string

	// Create mock Slack server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture Content-Type header
		receivedContentType = r.Header.Get("Content-Type")

		// Decode the payload
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)

		// Return success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	// Create notifier with test server
	notifier := NewNotifierWithClient(server.URL, server.Client())

	// Send a message
	err := notifier.Send(Message{
		Title: "Test Alert",
		Text:  "This is a test notification",
		Color: ColorDanger,
	})

	// Verify no error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify Content-Type header
	if receivedContentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", receivedContentType)
	}

	// Verify payload structure
	if len(receivedPayload.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(receivedPayload.Attachments))
	}

	attachment := receivedPayload.Attachments[0]

	if attachment.Title != "Test Alert" {
		t.Errorf("expected title 'Test Alert', got '%s'", attachment.Title)
	}

	if attachment.Text != "This is a test notification" {
		t.Errorf("expected text 'This is a test notification', got '%s'", attachment.Text)
	}

	if attachment.Color != "danger" {
		t.Errorf("expected color 'danger', got '%s'", attachment.Color)
	}

	// Verify timestamp is recent (within last 5 seconds)
	now := time.Now().Unix()
	if attachment.Timestamp < now-5 || attachment.Timestamp > now+5 {
		t.Errorf("expected timestamp around %d, got %d", now, attachment.Timestamp)
	}
}

func TestSend_WithFields(t *testing.T) {
	var receivedPayload slackPayload

	// Create mock Slack server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewNotifierWithClient(server.URL, server.Client())

	// Send message with fields
	err := notifier.Send(Message{
		Title: "Alert",
		Text:  "Service down",
		Color: ColorWarning,
		Fields: []Field{
			{Title: "URL", Value: "https://example.com", Short: true},
			{Title: "Status", Value: "DOWN", Short: true},
			{Title: "Error", Value: "Connection timeout", Short: false},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify fields
	attachment := receivedPayload.Attachments[0]
	if len(attachment.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(attachment.Fields))
	}

	// Check first field
	if attachment.Fields[0].Title != "URL" {
		t.Errorf("expected field title 'URL', got '%s'", attachment.Fields[0].Title)
	}
	if attachment.Fields[0].Value != "https://example.com" {
		t.Errorf("expected field value 'https://example.com', got '%s'", attachment.Fields[0].Value)
	}
	if !attachment.Fields[0].Short {
		t.Error("expected field to be short")
	}

	// Check third field
	if attachment.Fields[2].Title != "Error" {
		t.Errorf("expected field title 'Error', got '%s'", attachment.Fields[2].Title)
	}
	if attachment.Fields[2].Short {
		t.Error("expected field to not be short")
	}
}

func TestSend_ColorConstants(t *testing.T) {
	tests := []struct {
		name          string
		color         Color
		expectedColor string
	}{
		{"good color", ColorGood, "good"},
		{"warning color", ColorWarning, "warning"},
		{"danger color", ColorDanger, "danger"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedPayload slackPayload

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				json.Unmarshal(body, &receivedPayload)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			notifier := NewNotifierWithClient(server.URL, server.Client())
			err := notifier.Send(Message{
				Text:  "Test",
				Color: tt.color,
			})

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if receivedPayload.Attachments[0].Color != tt.expectedColor {
				t.Errorf("expected color '%s', got '%s'", tt.expectedColor, receivedPayload.Attachments[0].Color)
			}
		})
	}
}

func TestSend_SlackAPIError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid_payload"))
	}))
	defer server.Close()

	notifier := NewNotifierWithClient(server.URL, server.Client())

	err := notifier.Send(Message{
		Text: "Test",
	})

	// Verify error is returned
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Verify error message contains status code and response
	if !strings.Contains(err.Error(), "slack API error") {
		t.Errorf("expected error to mention 'slack API error', got: %v", err)
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected error to mention status code 400, got: %v", err)
	}
	if !strings.Contains(err.Error(), "invalid_payload") {
		t.Errorf("expected error to include response body, got: %v", err)
	}
}

func TestSend_NetworkError(t *testing.T) {
	// Use invalid URL to trigger network error
	notifier := NewNotifier("http://invalid-host-that-does-not-exist-12345.com/webhook")

	err := notifier.Send(Message{
		Text: "Test",
	})

	// Verify error is returned
	if err == nil {
		t.Fatal("expected network error, got nil")
	}

	// Verify error message mentions network failure
	if !strings.Contains(err.Error(), "network failure") {
		t.Errorf("expected error to mention 'network failure', got: %v", err)
	}
}

func TestSendSimple(t *testing.T) {
	var receivedPayload slackPayload

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &receivedPayload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	notifier := NewNotifierWithClient(server.URL, server.Client())

	err := notifier.SendSimple("Simple notification text")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify payload
	if len(receivedPayload.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(receivedPayload.Attachments))
	}

	attachment := receivedPayload.Attachments[0]
	if attachment.Text != "Simple notification text" {
		t.Errorf("expected text 'Simple notification text', got '%s'", attachment.Text)
	}

	// Title should be empty for simple messages
	if attachment.Title != "" {
		t.Errorf("expected empty title, got '%s'", attachment.Title)
	}
}

func TestValidateWebhookURL(t *testing.T) {
	tests := []struct {
		name        string
		webhookURL  string
		expectError bool
		errorText   string
	}{
		{
			name:        "valid https slack URL",
			webhookURL:  "https://hooks.slack.com/services/T00/B00/XXXX",
			expectError: false,
		},
		{
			name:        "valid http URL (for testing)",
			webhookURL:  "http://localhost:8080/webhook",
			expectError: false,
		},
		{
			name:        "empty URL",
			webhookURL:  "",
			expectError: true,
			errorText:   "webhook URL is empty",
		},
		{
			name:        "invalid URL format",
			webhookURL:  "not-a-valid-url",
			expectError: true,
			errorText:   "must use http or https scheme",
		},
		{
			name:        "missing scheme",
			webhookURL:  "hooks.slack.com/services/T00/B00/XXXX",
			expectError: true,
			errorText:   "must use http or https scheme",
		},
		{
			name:        "invalid scheme",
			webhookURL:  "ftp://hooks.slack.com/services/T00/B00/XXXX",
			expectError: true,
			errorText:   "must use http or https scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notifier := NewNotifier(tt.webhookURL)
			err := notifier.validateWebhookURL()

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errorText) {
					t.Errorf("expected error to contain '%s', got: %v", tt.errorText, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestSend_InvalidWebhookURL(t *testing.T) {
	notifier := NewNotifier("")

	err := notifier.Send(Message{
		Text: "Test",
	})

	if err == nil {
		t.Fatal("expected error for empty webhook URL, got nil")
	}

	if !strings.Contains(err.Error(), "webhook URL is empty") {
		t.Errorf("expected error about empty URL, got: %v", err)
	}
}

func TestConvertFields(t *testing.T) {
	// Test with fields
	fields := []Field{
		{Title: "Field1", Value: "Value1", Short: true},
		{Title: "Field2", Value: "Value2", Short: false},
	}

	result := convertFields(fields)

	if len(result) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(result))
	}

	if result[0].Title != "Field1" || result[0].Value != "Value1" || !result[0].Short {
		t.Error("first field not converted correctly")
	}

	if result[1].Title != "Field2" || result[1].Value != "Value2" || result[1].Short {
		t.Error("second field not converted correctly")
	}

	// Test with empty fields
	emptyResult := convertFields([]Field{})
	if emptyResult != nil {
		t.Errorf("expected nil for empty fields, got %v", emptyResult)
	}

	// Test with nil fields
	nilResult := convertFields(nil)
	if nilResult != nil {
		t.Errorf("expected nil for nil fields, got %v", nilResult)
	}
}
