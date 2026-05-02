package output

import "io"

// Format represents the output format type
type Format string

const (
	// FormatTable formats output as a table
	FormatTable Format = "table"
	// FormatJSON formats output as JSON
	FormatJSON Format = "json"
	// FormatPlain formats output as plain text
	FormatPlain Format = "plain"
)

// Formatter formats data for output
type Formatter interface {
	// Format formats data and writes to the writer
	Format(data interface{}, w io.Writer) error
}

// Options configures output formatting behavior
type Options struct {
	// Format specifies the output format (table, json, or plain)
	Format Format
	// NoColor disables colored output
	NoColor bool
	// NoProgress disables progress bars
	NoProgress bool
}

// NewFormatter creates a formatter based on the format type specified in options.
// Returns the appropriate formatter implementation (TableFormatter, JSONFormatter, or PlainFormatter).
func NewFormatter(opts Options) Formatter {
	switch opts.Format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatPlain:
		return &PlainFormatter{}
	case FormatTable:
		fallthrough
	default:
		return &TableFormatter{
			noColor: opts.NoColor,
		}
	}
}
