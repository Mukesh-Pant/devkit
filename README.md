# devkit - DevOps CLI Toolkit

A production-ready command-line tool written in Go that provides essential DevOps utilities including health checking, log analysis, notifications, and continuous monitoring. The tool compiles to a single static binary with no external runtime dependencies.

## Features

- **Health Checking**: Concurrent HTTP health checks with configurable timeouts and expected status codes
- **Continuous Monitoring**: Real-time URL monitoring with Slack alerting on status changes
- **Log Analysis**: Parse and summarize Apache/Nginx combined logs, JSON logs, and plain text logs
- **Slack Notifications**: Send formatted notifications to Slack channels via webhooks
- **Flexible Output**: Support for table, JSON, and plain text output formats
- **Configuration Management**: YAML-based configuration with command-line flag overrides
- **Cross-Platform**: Compile to static binaries for Linux, macOS, and Windows

## Requirements

- **Go**: Version 1.21 or higher
- **Dependencies**: All dependencies are managed via Go modules (see `go.mod`)
  - `github.com/spf13/cobra` - CLI framework
  - `github.com/spf13/viper` - Configuration management
  - `github.com/olekukonko/tablewriter` - Table formatting
  - `github.com/fatih/color` - Colored terminal output
  - `github.com/schollz/progressbar/v3` - Progress bars

## Installation

### Manual Installation

1. **Download the binary** for your platform from the [releases page](https://github.com/yourusername/devkit/releases):

   ```bash
   # Linux (amd64)
   curl -L https://github.com/yourusername/devkit/releases/download/v1.0.0/devkit-linux-amd64 -o devkit
   
   # macOS (amd64)
   curl -L https://github.com/yourusername/devkit/releases/download/v1.0.0/devkit-darwin-amd64 -o devkit
   
   # macOS (arm64 - Apple Silicon)
   curl -L https://github.com/yourusername/devkit/releases/download/v1.0.0/devkit-darwin-arm64 -o devkit
   
   # Windows (amd64)
   curl -L https://github.com/yourusername/devkit/releases/download/v1.0.0/devkit-windows-amd64.exe -o devkit.exe
   ```

2. **Make the binary executable** (Linux/macOS):

   ```bash
   chmod +x devkit
   ```

3. **Move to a directory in your PATH**:

   ```bash
   # Linux/macOS
   sudo mv devkit /usr/local/bin/
   
   # Or to a user directory
   mv devkit ~/bin/
   ```

4. **Verify the installation**:

   ```bash
   devkit version
   ```

### Build from Source

1. **Clone the repository**:

   ```bash
   git clone https://github.com/yourusername/devkit.git
   cd devkit
   ```

2. **Build the binary**:

   ```bash
   make build
   ```

   The binary will be created in the `dist/` directory.

3. **Install the binary**:

   ```bash
   sudo cp dist/devkit /usr/local/bin/
   ```

## Usage

### Global Flags

- `--config <path>` - Specify a custom configuration file path (default: `./.devkit.yaml` or `~/.devkit.yaml`)

### Commands

#### `devkit check` - Health Check URLs

Perform HTTP health checks on one or more URLs concurrently.

**Usage:**
```bash
devkit check <url1> <url2> ...
```

**Flags:**
- `--timeout <duration>` - Timeout for health checks (default: 5s)
- `--expected-status <code>` - Expected HTTP status code (default: 200)
- `--output <format>` - Output format: table, json, plain (default: table)

**Examples:**
```bash
# Check a single URL
devkit check https://example.com

# Check multiple URLs
devkit check https://api.example.com https://example.com/health

# Check with custom timeout and expected status
devkit check --timeout 10s --expected-status 200 https://example.com

# Output as JSON
devkit check --output json https://example.com
```

**Exit Codes:**
- `0` - All URLs are UP
- `1` - One or more URLs are DOWN

---

#### `devkit watch` - Continuous Monitoring

Continuously monitor URLs at regular intervals and send Slack alerts when status changes from UP to DOWN.

**Usage:**
```bash
devkit watch
```

**Flags:**
- `--interval <duration>` - Polling interval (default: 30s)
- `--alert-webhook <url>` - Slack webhook URL for alerts (default from config)

**Examples:**
```bash
# Monitor URLs from config with default interval
devkit watch

# Monitor with custom interval
devkit watch --interval 60s

# Monitor with Slack alerts
devkit watch --interval 30s --alert-webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

**Configuration:**

URLs to monitor must be specified in `.devkit.yaml`:

```yaml
watch:
  urls:
    - https://example.com
    - https://api.example.com/health
```

**Stopping:**

Press `Ctrl+C` to gracefully stop monitoring.

---

#### `devkit logs` - Parse and Analyze Log Files

Parse and summarize log files in various formats.

**Usage:**
```bash
devkit logs <filepath>
```

**Flags:**
- `--format <format>` - Log format: combined, json, plain (default: combined)
- `--top <n>` - Number of top items to display (default: 10)
- `--output <format>` - Output format: table, json, plain (default: table)

**Supported Formats:**
- `combined` - Apache/Nginx combined log format
- `json` - JSON log entries (one per line)
- `plain` - Plain text logs with basic pattern matching

**Statistics Provided:**
- Total lines processed
- Request counts by status code
- Error rate (percentage of 4xx and 5xx responses)
- Top N IP addresses by request count
- Top N request paths by request count
- Top N status codes by occurrence

**Examples:**
```bash
# Parse Apache/Nginx combined log
devkit logs /var/log/nginx/access.log --format combined

# Parse JSON logs and show top 20 items
devkit logs /var/log/app.json --format json --top 20

# Parse plain text logs with JSON output
devkit logs /var/log/app.log --format plain --output json
```

---

#### `devkit notify` - Send Slack Notifications

Send a formatted notification to a Slack channel via webhook.

**Usage:**
```bash
devkit notify --message <text>
```

**Flags:**
- `--webhook <url>` - Slack webhook URL (required if not in config)
- `--message <text>` - Notification text (required)
- `--title <text>` - Notification title (optional)
- `--color <color>` - Notification color: good, warning, danger, or hex code (optional)

**Colors:**
- `good` - Green (success)
- `warning` - Yellow (warning)
- `danger` - Red (error/alert)
- Hex code - Custom color (e.g., `#FF5733`)

**Examples:**
```bash
# Simple notification
devkit notify --message "Deployment completed successfully" --color good

# Notification with title and custom color
devkit notify --webhook https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
  --message "Production issue detected" \
  --title "🚨 Alert" \
  --color danger

# Build notification
devkit notify --message "Build finished" --title "CI/CD" --color warning
```

---

#### `devkit version` - Version Information

Display the version number, git commit hash, and build date.

**Usage:**
```bash
devkit version
```

**Example Output:**
```
devkit version 1.0.0
commit: a1b2c3d
built: 2024-01-15T10:30:00Z
```

---

## Configuration

devkit uses a YAML configuration file named `.devkit.yaml` for default settings. The file is searched in the following order:

1. Current directory: `./.devkit.yaml`
2. Home directory: `~/.devkit.yaml`
3. Custom path via `--config` flag

Command-line flags always override configuration file values.

### Configuration File Example

See `.devkit.yaml.example` for a complete example. Here's a basic configuration:

```yaml
# Global devkit settings
devkit:
  # Default timeout for HTTP health checks
  timeout: "5s"
  
  # Default output format for all commands
  output: "table"
  
  # Default Slack webhook URL for notifications
  slack_webhook: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

# Watch command settings
watch:
  # Polling interval for continuous monitoring
  interval: "30s"
  
  # List of URLs to monitor continuously
  urls:
    - "https://example.com"
    - "https://api.example.com/health"
    - "https://example.com/status"
```

### Configuration Keys

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| `devkit.timeout` | duration | Default timeout for health checks | `5s` |
| `devkit.output` | string | Default output format (table, json, plain) | `table` |
| `devkit.slack_webhook` | string | Default Slack webhook URL | `""` |
| `watch.interval` | duration | Polling interval for watch command | `30s` |
| `watch.urls` | []string | List of URLs to monitor | `[]` |

### Duration Format

Duration values support these units:
- `ns` - nanoseconds
- `us` or `µs` - microseconds
- `ms` - milliseconds
- `s` - seconds
- `m` - minutes
- `h` - hours

Examples: `5s`, `30s`, `1m`, `500ms`, `1h30m`

---

## Build Instructions

### Prerequisites

- Go 1.21 or higher
- Make (optional, but recommended)

### Build Targets

The project includes a Makefile with the following targets:

#### `make build`

Build the binary for the current platform:

```bash
make build
```

The binary will be created at `dist/devkit`.

#### `make test`

Run all unit tests:

```bash
make test
```

#### `make test-verbose`

Run tests with verbose output:

```bash
make test-verbose
```

#### `make clean`

Remove build artifacts:

```bash
make clean
```

#### `make release`

Cross-compile for all supported platforms:

```bash
make release
```

This creates binaries for:
- Linux (amd64): `dist/devkit-linux-amd64`
- macOS (amd64): `dist/devkit-darwin-amd64`
- macOS (arm64): `dist/devkit-darwin-arm64`
- Windows (amd64): `dist/devkit-windows-amd64.exe`

### Version Injection

Version information is injected at build time via ldflags:

```bash
make build VERSION=1.0.0
```

The Makefile automatically injects:
- `VERSION` - Version string (default: `dev`)
- `COMMIT` - Git commit hash (auto-detected)
- `BUILD_DATE` - Build timestamp (auto-generated)

### Manual Build

If you prefer not to use Make:

```bash
# Build for current platform
go build -ldflags "-X main.Version=1.0.0 -X main.Commit=$(git rev-parse --short HEAD) -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o devkit .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o devkit-linux-amd64 .

# Cross-compile for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o devkit-darwin-amd64 .

# Cross-compile for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o devkit-darwin-arm64 .

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o devkit-windows-amd64.exe .
```

---

## Project Structure

```
devkit/
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command and global flags
│   ├── check.go                  # Health check command
│   ├── watch.go                  # Continuous monitoring command
│   ├── logs.go                   # Log parsing command
│   ├── notify.go                 # Slack notification command
│   └── version.go                # Version information command
├── internal/                     # Internal packages
│   ├── checker/                  # Health checking logic
│   │   ├── checker.go            # Health check implementation
│   │   └── checker_test.go       # Unit tests
│   ├── logparser/                # Log parsing logic
│   │   ├── parser.go             # Log parser interface
│   │   ├── combined.go           # Apache/Nginx combined format
│   │   ├── json.go               # JSON log format
│   │   ├── plain.go              # Plain text format
│   │   └── parser_test.go        # Unit tests
│   ├── notifier/                 # Slack notification logic
│   │   ├── slack.go              # Slack webhook client
│   │   └── slack_test.go         # Unit tests
│   └── output/                   # Output formatting
│       ├── formatter.go          # Output formatter interface
│       ├── table.go              # Table format implementation
│       ├── json.go               # JSON format implementation
│       └── plain.go              # Plain text implementation
├── config/                       # Configuration management
│   ├── config.go                 # Viper configuration wrapper
│   └── config_test.go            # Unit tests
├── main.go                       # Application entry point
├── Makefile                      # Build automation
├── .devkit.yaml.example          # Example configuration file
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
└── README.md                     # This file
```

### Package Responsibilities

- **cmd**: CLI command definitions using Cobra framework
- **internal/checker**: HTTP health checking with concurrent execution
- **internal/logparser**: Log file parsing for multiple formats
- **internal/notifier**: Slack webhook integration for notifications
- **internal/output**: Output formatting (table, JSON, plain text)
- **config**: Configuration file management using Viper

---

## Examples

### Example 1: Production Monitoring with Slack Alerts

**Configuration** (`.devkit.yaml`):
```yaml
devkit:
  timeout: "10s"
  slack_webhook: "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

watch:
  interval: "60s"
  urls:
    - "https://api.production.com/health"
    - "https://www.production.com"
    - "https://admin.production.com/status"
```

**Command:**
```bash
devkit watch
```

This will monitor all configured URLs every 60 seconds and send Slack alerts when any URL goes down.

---

### Example 2: CI/CD Health Check

**Script** (`check-services.sh`):
```bash
#!/bin/bash

# Check critical services
devkit check \
  --timeout 15s \
  --output json \
  https://api.example.com/health \
  https://db.example.com/health \
  https://cache.example.com/ping

# Exit code will be 1 if any service is down
if [ $? -ne 0 ]; then
  devkit notify \
    --message "Service health check failed in CI/CD pipeline" \
    --title "🚨 CI/CD Alert" \
    --color danger
  exit 1
fi
```

---

### Example 3: Log Analysis for Nginx

**Command:**
```bash
# Analyze today's access log
devkit logs /var/log/nginx/access.log --format combined --top 20

# Export statistics as JSON for further processing
devkit logs /var/log/nginx/access.log --format combined --output json > stats.json
```

---

### Example 4: Development Environment Monitoring

**Configuration** (`.devkit.yaml`):
```yaml
devkit:
  timeout: "3s"
  output: "table"

watch:
  interval: "10s"
  urls:
    - "http://localhost:3000"
    - "http://localhost:8080/health"
    - "http://localhost:5432"
```

**Command:**
```bash
devkit watch
```

This provides fast feedback during development by checking local services every 10 seconds.

---

## Error Handling

devkit provides clear, actionable error messages:

```bash
# Invalid URL
$ devkit check invalid-url
Error: failed to check invalid-url: unsupported protocol scheme ""

# File not found
$ devkit logs /nonexistent/file.log
Error: failed to parse log file: open /nonexistent/file.log: no such file or directory

# Missing webhook URL
$ devkit notify --message "Test"
Error: webhook URL is required: provide via --webhook flag or devkit.slack_webhook in config
```

### Exit Codes

- `0` - Success
- `1` - Error (invalid arguments, configuration errors, file errors, health check failures)

---

## Security Considerations

### Webhook URLs

- **Never commit** `.devkit.yaml` files containing webhook URLs to version control
- Use environment-specific configuration files or command-line flags for sensitive values
- Consider using environment variables for webhook URLs in CI/CD pipelines

### Network Security

- All HTTP requests use proper TLS certificate validation
- Timeouts prevent hanging connections
- Concurrent checks are limited to prevent resource exhaustion

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for new functionality
4. Ensure all tests pass (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/yourusername/devkit).
