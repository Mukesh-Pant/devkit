package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"devkit/internal/checker"
)

func TestTableFormatter_Format_CheckResults(t *testing.T) {
	tests := []struct {
		name     string
		results  []checker.CheckResult
		noColor  bool
		wantErr  bool
		contains []string
	}{
		{
			name: "successful check results",
			results: []checker.CheckResult{
				{
					URL:          "https://example.com",
					StatusCode:   200,
					ResponseTime: 123 * time.Millisecond,
					Status:       checker.StatusUp,
				},
				{
					URL:          "https://example.com/api",
					StatusCode:   500,
					ResponseTime: 5000 * time.Millisecond,
					Status:       checker.StatusDown,
				},
			},
			noColor: true,
			wantErr: false,
			contains: []string{
				"URL",
				"STATUS",
				"RESPONSE TIME",
				"CODE",
				"https://example.com",
				"UP",
				"123ms",
				"200",
				"https://example.com/api",
				"DOWN",
				"5000ms",
				"500",
			},
		},
		{
			name: "check result with error",
			results: []checker.CheckResult{
				{
					URL:          "https://invalid.example",
					StatusCode:   0,
					ResponseTime: 0,
					Status:       checker.StatusDown,
					Error:        nil,
				},
			},
			noColor: true,
			wantErr: false,
			contains: []string{
				"https://invalid.example",
				"DOWN",
				"-", // Code should be "-" when StatusCode is 0
			},
		},
		{
			name: "empty results",
			results: []checker.CheckResult{},
			noColor: true,
			wantErr: false,
			contains: []string{
				"URL",
				"STATUS",
				"RESPONSE TIME",
				"CODE",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{noColor: tt.noColor}
			var buf bytes.Buffer

			err := f.Format(tt.results, &buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("TableFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			for _, want := range tt.contains {
				if !strings.Contains(output, want) {
					t.Errorf("TableFormatter.Format() output missing expected string %q\nGot:\n%s", want, output)
				}
			}
		})
	}
}

func TestTableFormatter_Format_UnsupportedType(t *testing.T) {
	f := &TableFormatter{noColor: true}
	var buf bytes.Buffer

	err := f.Format("unsupported", &buf)

	if err == nil {
		t.Error("TableFormatter.Format() expected error for unsupported type, got nil")
	}

	if !strings.Contains(err.Error(), "unsupported data type") {
		t.Errorf("TableFormatter.Format() error = %v, want error containing 'unsupported data type'", err)
	}
}

func TestTableFormatter_FormatStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   checker.Status
		noColor  bool
		wantText string
	}{
		{
			name:     "UP status with color disabled",
			status:   checker.StatusUp,
			noColor:  true,
			wantText: "UP",
		},
		{
			name:     "DOWN status with color disabled",
			status:   checker.StatusDown,
			noColor:  true,
			wantText: "DOWN",
		},
		{
			name:     "UP status with color enabled",
			status:   checker.StatusUp,
			noColor:  false,
			wantText: "UP",
		},
		{
			name:     "DOWN status with color enabled",
			status:   checker.StatusDown,
			noColor:  false,
			wantText: "DOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &TableFormatter{noColor: tt.noColor}
			result := f.formatStatus(tt.status)

			if !strings.Contains(result, tt.wantText) {
				t.Errorf("TableFormatter.formatStatus() = %q, want to contain %q", result, tt.wantText)
			}

			// When noColor is true, result should be exactly the status text
			if tt.noColor && result != tt.wantText {
				t.Errorf("TableFormatter.formatStatus() with noColor = %q, want %q", result, tt.wantText)
			}
		})
	}
}
