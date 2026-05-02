package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewChecker verifies the constructor creates a checker with proper defaults
func TestNewChecker(t *testing.T) {
	opts := Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	}

	checker := NewChecker(opts)

	if checker == nil {
		t.Fatal("NewChecker returned nil")
	}

	if checker.opts.Timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", checker.opts.Timeout)
	}

	if checker.opts.ExpectedStatus != 200 {
		t.Errorf("Expected status 200, got %d", checker.opts.ExpectedStatus)
	}

	if checker.opts.Client == nil {
		t.Error("Expected client to be initialized")
	}
}

// TestNewChecker_DefaultExpectedStatus verifies default expected status is 200
func TestNewChecker_DefaultExpectedStatus(t *testing.T) {
	opts := Options{
		Timeout: 5 * time.Second,
	}

	checker := NewChecker(opts)

	if checker.opts.ExpectedStatus != http.StatusOK {
		t.Errorf("Expected default status 200, got %d", checker.opts.ExpectedStatus)
	}
}

// TestCheck_Success verifies UP detection for 200 OK responses
func TestCheck_Success(t *testing.T) {
	// Create mock server that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	result := checker.Check(server.URL)

	if result.Status != StatusUp {
		t.Errorf("Expected status UP, got %s", result.Status)
	}

	if result.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", result.StatusCode)
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}

	if result.URL != server.URL {
		t.Errorf("Expected URL %s, got %s", server.URL, result.URL)
	}

	if result.ResponseTime <= 0 {
		t.Error("Expected positive response time")
	}
}

// TestCheck_InternalServerError verifies DOWN detection for 500 responses
func TestCheck_InternalServerError(t *testing.T) {
	// Create mock server that returns 500 Internal Server Error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	result := checker.Check(server.URL)

	if result.Status != StatusDown {
		t.Errorf("Expected status DOWN, got %s", result.Status)
	}

	if result.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", result.StatusCode)
	}

	if result.Error == nil {
		t.Error("Expected error for non-200 status code")
	}
}

// TestCheck_NotFound verifies DOWN detection for 404 responses
func TestCheck_NotFound(t *testing.T) {
	// Create mock server that returns 404 Not Found
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	result := checker.Check(server.URL)

	if result.Status != StatusDown {
		t.Errorf("Expected status DOWN, got %s", result.Status)
	}

	if result.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", result.StatusCode)
	}

	if result.Error == nil {
		t.Error("Expected error for 404 status code")
	}
}

// TestCheck_Timeout verifies timeout handling with slow server
func TestCheck_Timeout(t *testing.T) {
	// Create mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Set timeout shorter than server delay
	checker := NewChecker(Options{
		Timeout:        500 * time.Millisecond,
		ExpectedStatus: 200,
	})

	result := checker.Check(server.URL)

	if result.Status != StatusDown {
		t.Errorf("Expected status DOWN for timeout, got %s", result.Status)
	}

	if result.Error == nil {
		t.Error("Expected timeout error")
	}
}

// TestCheck_CustomExpectedStatus verifies custom expected status codes
func TestCheck_CustomExpectedStatus(t *testing.T) {
	// Create mock server that returns 201 Created
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 201,
	})

	result := checker.Check(server.URL)

	if result.Status != StatusUp {
		t.Errorf("Expected status UP for 201, got %s", result.Status)
	}

	if result.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", result.StatusCode)
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got %v", result.Error)
	}
}

// TestCheck_InvalidURL verifies handling of invalid URLs
func TestCheck_InvalidURL(t *testing.T) {
	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	result := checker.Check("not-a-valid-url")

	if result.Status != StatusDown {
		t.Errorf("Expected status DOWN for invalid URL, got %s", result.Status)
	}

	if result.Error == nil {
		t.Error("Expected error for invalid URL")
	}
}

// TestCheck_ConnectionRefused verifies handling of connection refused errors
func TestCheck_ConnectionRefused(t *testing.T) {
	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	// Use a URL that will refuse connection (port unlikely to be in use)
	result := checker.Check("http://localhost:59999")

	if result.Status != StatusDown {
		t.Errorf("Expected status DOWN for connection refused, got %s", result.Status)
	}

	if result.Error == nil {
		t.Error("Expected error for connection refused")
	}
}

// TestCheckMultiple_Success verifies concurrent checking of multiple URLs
func TestCheckMultiple_Success(t *testing.T) {
	// Create multiple mock servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server2.Close()

	server3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server3.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	urls := []string{server1.URL, server2.URL, server3.URL}
	results := checker.CheckMultiple(urls)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Verify all checks succeeded
	for _, result := range results {
		if result.Status != StatusUp {
			t.Errorf("Expected status UP for %s, got %s", result.URL, result.Status)
		}
		if result.StatusCode != 200 {
			t.Errorf("Expected status code 200 for %s, got %d", result.URL, result.StatusCode)
		}
	}
}

// TestCheckMultiple_MixedResults verifies handling of mixed UP/DOWN results
func TestCheckMultiple_MixedResults(t *testing.T) {
	// Create servers with different responses
	serverUp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer serverUp.Close()

	serverDown := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverDown.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	urls := []string{serverUp.URL, serverDown.URL}
	results := checker.CheckMultiple(urls)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Count UP and DOWN statuses
	upCount := 0
	downCount := 0
	for _, result := range results {
		if result.Status == StatusUp {
			upCount++
		} else if result.Status == StatusDown {
			downCount++
		}
	}

	if upCount != 1 {
		t.Errorf("Expected 1 UP status, got %d", upCount)
	}

	if downCount != 1 {
		t.Errorf("Expected 1 DOWN status, got %d", downCount)
	}
}

// TestCheckMultiple_EmptyList verifies handling of empty URL list
func TestCheckMultiple_EmptyList(t *testing.T) {
	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	results := checker.CheckMultiple([]string{})

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty list, got %d", len(results))
	}
}

// TestCheckMultiple_Concurrency verifies that checks run concurrently
func TestCheckMultiple_Concurrency(t *testing.T) {
	// Create servers that delay response
	delay := 500 * time.Millisecond
	createDelayedServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(delay)
			w.WriteHeader(http.StatusOK)
		}))
	}

	server1 := createDelayedServer()
	defer server1.Close()

	server2 := createDelayedServer()
	defer server2.Close()

	server3 := createDelayedServer()
	defer server3.Close()

	checker := NewChecker(Options{
		Timeout:        2 * time.Second,
		ExpectedStatus: 200,
	})

	urls := []string{server1.URL, server2.URL, server3.URL}

	start := time.Now()
	results := checker.CheckMultiple(urls)
	elapsed := time.Since(start)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// If checks were sequential, it would take 3 * delay
	// If concurrent, it should take approximately delay
	// Allow some overhead for goroutine scheduling
	maxExpectedTime := delay + 300*time.Millisecond
	if elapsed > maxExpectedTime {
		t.Errorf("Checks appear to be sequential. Expected ~%v, took %v", delay, elapsed)
	}
}

// TestCheck_ResponseTimeMeasurement verifies response time is measured correctly
func TestCheck_ResponseTimeMeasurement(t *testing.T) {
	delay := 100 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	result := checker.Check(server.URL)

	// Response time should be at least the delay
	if result.ResponseTime < delay {
		t.Errorf("Expected response time >= %v, got %v", delay, result.ResponseTime)
	}

	// Response time should be reasonable (not way too long)
	maxExpectedTime := delay + 200*time.Millisecond
	if result.ResponseTime > maxExpectedTime {
		t.Errorf("Response time too long. Expected ~%v, got %v", delay, result.ResponseTime)
	}
}

// TestCheckWithProgress_Success verifies progress bar integration with concurrent checking
func TestCheckWithProgress_Success(t *testing.T) {
	// Create multiple mock servers
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server2.Close()

	server3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server3.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	urls := []string{server1.URL, server2.URL, server3.URL}
	results := checker.CheckWithProgress(urls)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Verify all checks succeeded
	for _, result := range results {
		if result.Status != StatusUp {
			t.Errorf("Expected status UP for %s, got %s", result.URL, result.Status)
		}
		if result.StatusCode != 200 {
			t.Errorf("Expected status code 200 for %s, got %d", result.URL, result.StatusCode)
		}
	}
}

// TestCheckWithProgress_MixedResults verifies progress bar with mixed UP/DOWN results
func TestCheckWithProgress_MixedResults(t *testing.T) {
	// Create servers with different responses
	serverUp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer serverUp.Close()

	serverDown := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer serverDown.Close()

	checker := NewChecker(Options{
		Timeout:        5 * time.Second,
		ExpectedStatus: 200,
	})

	urls := []string{serverUp.URL, serverDown.URL}
	results := checker.CheckWithProgress(urls)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Count UP and DOWN statuses
	upCount := 0
	downCount := 0
	for _, result := range results {
		if result.Status == StatusUp {
			upCount++
		} else if result.Status == StatusDown {
			downCount++
		}
	}

	if upCount != 1 {
		t.Errorf("Expected 1 UP status, got %d", upCount)
	}

	if downCount != 1 {
		t.Errorf("Expected 1 DOWN status, got %d", downCount)
	}
}

// TestCheckWithProgress_Concurrency verifies that checks with progress bar run concurrently
func TestCheckWithProgress_Concurrency(t *testing.T) {
	// Create servers that delay response
	delay := 500 * time.Millisecond
	createDelayedServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(delay)
			w.WriteHeader(http.StatusOK)
		}))
	}

	server1 := createDelayedServer()
	defer server1.Close()

	server2 := createDelayedServer()
	defer server2.Close()

	server3 := createDelayedServer()
	defer server3.Close()

	checker := NewChecker(Options{
		Timeout:        2 * time.Second,
		ExpectedStatus: 200,
	})

	urls := []string{server1.URL, server2.URL, server3.URL}

	start := time.Now()
	results := checker.CheckWithProgress(urls)
	elapsed := time.Since(start)

	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// If checks were sequential, it would take 3 * delay
	// If concurrent, it should take approximately delay
	// Allow some overhead for goroutine scheduling and progress bar rendering
	maxExpectedTime := delay + 300*time.Millisecond
	if elapsed > maxExpectedTime {
		t.Errorf("Checks appear to be sequential. Expected ~%v, took %v", delay, elapsed)
	}
}
