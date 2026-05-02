// +build ignore

package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"devkit/internal/checker"
	"devkit/internal/output"
)

func main() {
	// Create sample check results
	results := []checker.CheckResult{
		{
			URL:          "https://example.com",
			StatusCode:   200,
			ResponseTime: 123 * time.Millisecond,
			Status:       checker.StatusUp,
		},
		{
			URL:          "https://api.example.com",
			StatusCode:   500,
			ResponseTime: 5000 * time.Millisecond,
			Status:       checker.StatusDown,
			Error:        errors.New("unexpected status code: got 500, expected 200"),
		},
		{
			URL:          "https://invalid.example.com",
			StatusCode:   0,
			ResponseTime: 0,
			Status:       checker.StatusDown,
			Error:        errors.New("failed to check https://invalid.example.com: dial tcp: connection refused"),
		},
	}

	// Demonstrate Table format
	fmt.Println("=== TABLE FORMAT ===")
	tableFormatter := output.NewFormatter(output.Options{
		Format:  output.FormatTable,
		NoColor: false,
	})
	tableFormatter.Format(results, os.Stdout)
	fmt.Println()

	// Demonstrate JSON format
	fmt.Println("=== JSON FORMAT ===")
	jsonFormatter := output.NewFormatter(output.Options{
		Format: output.FormatJSON,
	})
	jsonFormatter.Format(results, os.Stdout)
	fmt.Println()

	// Demonstrate Plain format
	fmt.Println("=== PLAIN FORMAT ===")
	plainFormatter := output.NewFormatter(output.Options{
		Format: output.FormatPlain,
	})
	plainFormatter.Format(results, os.Stdout)
}
