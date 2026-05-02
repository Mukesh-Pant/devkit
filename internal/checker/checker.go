package checker

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Common errors
var (
	ErrTimeout = errors.New("request timeout")
)

// Status represents the health status
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// CheckResult represents the result of a single health check
type CheckResult struct {
	URL          string        `json:"url"`
	StatusCode   int           `json:"status_code"`
	ResponseTime time.Duration `json:"response_time_ms"`
	Status       Status        `json:"status"`
	Error        error         `json:"error,omitempty"`
}

// Options configures health check behavior
type Options struct {
	Timeout        time.Duration
	ExpectedStatus int
	Client         *http.Client // Allow custom client for testing
}

// Checker performs health checks
type Checker struct {
	opts Options
}

// NewChecker creates a new health checker with options
func NewChecker(opts Options) *Checker {
	// If no custom client provided, create one with the specified timeout
	if opts.Client == nil {
		opts.Client = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	// Set default expected status if not specified
	if opts.ExpectedStatus == 0 {
		opts.ExpectedStatus = http.StatusOK
	}

	return &Checker{
		opts: opts,
	}
}

// Check performs a health check on a single URL
func (c *Checker) Check(url string) CheckResult {
	result := CheckResult{
		URL:    url,
		Status: StatusDown, // Default to DOWN
	}

	// Measure response time
	start := time.Now()

	// Perform HTTP GET request
	resp, err := c.opts.Client.Get(url)
	responseTime := time.Since(start)
	result.ResponseTime = responseTime

	// Handle network errors, timeouts, and invalid URLs
	if err != nil {
		result.Error = fmt.Errorf("failed to check %s: %w", url, err)
		return result
	}
	defer resp.Body.Close()

	// Record status code
	result.StatusCode = resp.StatusCode

	// Validate status code against expected status
	if resp.StatusCode == c.opts.ExpectedStatus {
		result.Status = StatusUp
	} else {
		result.Status = StatusDown
		result.Error = fmt.Errorf("unexpected status code: got %d, expected %d", resp.StatusCode, c.opts.ExpectedStatus)
	}

	return result
}

// CheckMultiple performs health checks on multiple URLs concurrently
func (c *Checker) CheckMultiple(urls []string) []CheckResult {
	// Create a channel to collect results
	resultsChan := make(chan CheckResult, len(urls))

	// Use WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Launch a goroutine for each URL
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := c.Check(u)
			resultsChan <- result
		}(url)
	}

	// Wait for all checks to complete and close the channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results from the channel
	results := make([]CheckResult, 0, len(urls))
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// CheckWithProgress performs health checks on multiple URLs concurrently with a progress bar
func (c *Checker) CheckWithProgress(urls []string) []CheckResult {
	// Create progress bar
	bar := progressbar.NewOptions(len(urls),
		progressbar.OptionSetDescription("Checking URLs"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	// Create a channel to collect results
	resultsChan := make(chan CheckResult, len(urls))

	// Use WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	// Launch a goroutine for each URL
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			result := c.Check(u)
			resultsChan <- result
			bar.Add(1) // Update progress bar
		}(url)
	}

	// Wait for all checks to complete and close the channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results from the channel
	results := make([]CheckResult, 0, len(urls))
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}
