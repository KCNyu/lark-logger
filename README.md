# lark-logger

[![CI](https://github.com/KCNyu/lark-logger/workflows/CI/badge.svg)](https://github.com/KCNyu/lark-logger/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/KCNyu/lark-logger)](https://goreportcard.com/report/github.com/KCNyu/lark-logger)
[![GoDoc](https://godoc.org/github.com/KCNyu/lark-logger?status.svg)](https://godoc.org/github.com/KCNyu/lark-logger)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

English | [中文文档](docs/README_zh.md)

An elegant Go SDK for sending structured log messages to Lark (Feishu) webhook bots with beautiful interactive card formatting.

## ✨ Features

- 🎨 **Beautiful Card Formatting** - Support for emojis, color templates, and Markdown formatting
- 📊 **Structured Logging** - Support for `map[string]interface{}` parameters with smart formatting
- 🚦 **Multi-level Logging** - Support for Info, Warn, and Error levels with unique styling and colors
- 🔄 **Retry Mechanism** - Built-in retry logic to ensure reliable message delivery
- ⚙️ **Type-safe Configuration** - Using functional options pattern for type-safe configuration
- 🎨 **Custom Titles** - Support for custom main titles and formatted subtitles with time and service info
- 📋 **Smart Field Display** - Clean KV table layout with proper alignment and readability
- 🔧 **Enhanced Layout** - Professional two-column layout with alternating row colors
- 📱 **Mobile Optimization** - Responsive design with proper padding and font sizing
- 🎨 **Visual Hierarchy** - Different font sizes and colors for headers and content
- 🚀 **Simplified API** - Support for Infof/Warnf/Errorf methods with key-value pairs
- 🔧 **Configuration Display** - Optional 2x2 grid layout for logger configuration (hidden by default)
- 🎯 **Environment Management** - Smart environment variable handling for local/CI/test modes
- 🧪 **Complete Test Coverage** - Includes unit tests and integration tests
- 🚀 **CI/CD Ready** - Includes GitHub Actions workflows and GoReleaser configuration

## 📦 Installation

```bash
go get github.com/KCNyu/lark-logger
```

## 🚀 Quick Start

### Environment Configuration

Before using lark-logger, you need to configure environment variables:

#### Option 1: Using .env files (Recommended)

1. Copy the example environment file:
```bash
cp env.example .env.local
```

2. Edit `.env.local` file with your actual webhook URL:
```bash
# For production use (sends real messages to Lark)
LARK_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url
LARK_TEST_MODE=false

# For testing/development (uses test webhook, no real messages sent)
# LARK_WEBHOOK_URL=https://test.webhook.url
# LARK_TEST_MODE=true
```

3. Load environment variables:
```bash
source scripts/load-env.sh
```

**Note**: `.env.local` files are automatically ignored by git to keep your webhook URLs secure.

#### Option 2: Direct export

```bash
# For production use (sends real messages to Lark)
export LARK_WEBHOOK_URL="https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"
export LARK_TEST_MODE="false"

# For testing/development (uses test webhook, no real messages sent)
export LARK_TEST_MODE="true"
export LARK_WEBHOOK_URL="https://test.webhook.url"
```

#### CI/Testing Configuration

For CI and testing environments, the following environment variables are automatically set:
- `LARK_TEST_MODE=true`
- `LARK_WEBHOOK_URL=https://test.webhook.url`

### Basic Usage

```go
package main

import (
    "github.com/KCNyu/lark-logger"
)

func main() {
    // Get webhook URL from environment variable
    webhookURL := larklogger.GetWebhookURL()

    // Create logger with webhook URL and options
    logger := larklogger.New(
        webhookURL,
        larklogger.WithEnv("production"),
        larklogger.WithTitle("System Monitor"),
    )

    // Send logs using traditional map format
    logger.Info("API Gateway initialized successfully", map[string]interface{}{
        "port":     8080,
        "version":  "2.1.0",
        "features": []string{"authentication", "rate_limiting", "metrics"},
    })

    // Send logs using simplified Infof/Warnf/Errorf format
    logger.Infof("Service health check", "status", "healthy", "response_time", "120ms", "uptime", "2h30m")
    logger.Warnf("Memory usage approaching threshold", "usage", "87%", "threshold", "85%", "recommendation", "consider horizontal scaling")
    logger.Errorf("Database connection pool exhausted", "error", "connection timeout after 30s", "retry_count", 3, "pool_size", 20)

    // Create logger with configuration display enabled
    configLogger := larklogger.New(
        webhookURL,
        larklogger.WithService("config-demo"),
        larklogger.WithEnv("production"),
        larklogger.WithHostname("server-01"),
        larklogger.WithTitle("Config Demo"),
        larklogger.WithShowConfig(true), // Enable configuration section
    )

    // This will show the configuration section with 2x2 grid layout
    configLogger.Info("Configuration section enabled", map[string]interface{}{
        "feature": "config_visibility",
        "status":  "enabled",
    })
}
```

### Send Simple Text Messages

```go
err := client.SendText("Hello from larklogger! 🚀")
if err != nil {
    log.Printf("Failed to send text message: %v", err)
}
```

### Custom Cards with Advanced Layout

```go
// Create a beautiful custom card with multiple layouts
card := larklogger.NewCardBuilder().
    SetHeader("🎉 System Status Report", "green").
    AddSection("**System is running smoothly**").
    AddKeyValueList("System Info", map[string]interface{}{
        "Uptime": "7 days, 12 hours",
        "Version": "v2.1.0",
        "Environment": "production",
    }).
    AddStatusBadge("success", "All services operational").
    AddMetricsGrid("Performance Metrics", map[string]interface{}{
        "CPU": "45%",
        "Memory": "67%", 
        "Disk": "23%",
        "Network": "12%",
    }).
    Build()

err := client.SendCard(card)
```

### Enhanced Logger with Custom Title

```go
// Create logger with custom title and enhanced formatting
logger := larklogger.NewLogger(client,
    larklogger.WithService("api-gateway"),
    larklogger.WithEnv("production"),
    larklogger.WithHostname("gateway-01"),
    larklogger.WithTitle("🌐 API Gateway Monitor"),
)

// Send log with rich structured data
logger.Info("Request processed successfully", map[string]interface{}{
    "method":     "POST",
    "path":       "/api/v1/users",
    "status":     200,
    "duration":   "45ms",
    "user_id":    "12345",
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "features":   []string{"auth", "rate-limit", "logging"},
})
```

This will create a beautiful enhanced card with:
- **Custom main title**: "🚀 System Monitor" with emoji and color template
- **Structured layout**: Clear sections with proper dividers and spacing
- **Enhanced KV table**: Two-column layout with alternating row colors
- **Smart field prioritization**: Important fields (like error codes) displayed prominently
- **Professional styling**: Proper padding, alignment, and mobile optimization
- **Visual hierarchy**: Different font sizes for headers and important values
- **Rich context**: Detailed system information with nested data structures

## 📋 API Documentation

### Logger Interface

```go
type Logger interface {
    Info(message string, fields map[string]interface{})
    Warn(message string, fields map[string]interface{})
    Error(message string, fields map[string]interface{})
}
```

### Client Configuration Options

```go
// Timeout setting
WithTimeout(timeout time.Duration) ClientOption

// Retry configuration
WithRetry(count int, delay time.Duration) ClientOption

// User agent
WithUserAgent(userAgent string) ClientOption

// Custom headers
WithHeaders(headers map[string]string) ClientOption
```

### Logger Configuration Options

```go
// Service name (optional, defaults to "larklogger")
WithService(service string) LoggerOption

// Environment identifier (optional, defaults to "development")
WithEnv(env string) LoggerOption

// Hostname (optional, defaults to "localhost")
WithHostname(hostname string) LoggerOption

// Custom title for log cards (optional, defaults to "System Log")
WithTitle(title string) LoggerOption

// Show configuration section in cards (optional, defaults to false)
WithShowConfig(show bool) LoggerOption
```

### Card Builder Methods

```go
// Set card header with title and color template
SetHeader(title, template string) *CardBuilder

// Add a section with text content
AddSection(text string) *CardBuilder

// Add a divider line
AddDivider() *CardBuilder

// Add key-value list with formatted display
AddKeyValueList(title string, items map[string]interface{}) *CardBuilder

// Add status badge with emoji and color
AddStatusBadge(status, message string) *CardBuilder

// Add metrics grid with table layout
AddMetricsGrid(title string, metrics map[string]interface{}) *CardBuilder

// Add two-column field layout
AddTwoColumnFields(fields []*CardField) *CardBuilder
```

### Enhanced Card Builder Methods

```go
// Create enhanced card builder
NewEnhancedCardBuilder() *EnhancedCardBuilder

// Set card header with title and color template
SetHeader(title, template string) *EnhancedCardBuilder

// Add subtitle with message
AddSubtitle(subtitle string) *EnhancedCardBuilder

// Add timestamp (right-aligned)
AddTimestamp() *EnhancedCardBuilder

// Add divider line
AddDivider() *EnhancedCardBuilder

// Add professional KV table with alternating colors
AddKVTable(kvList []KVItem) *EnhancedCardBuilder

// Add optional card link
AddCardLink(url string) *CardBuilder
```

### Simplified Logging Methods

```go
// Send logs with simplified key-value pairs (no map needed)
logger.Infof(title string, args ...interface{})
logger.Warnf(title string, args ...interface{})
logger.Errorf(title string, args ...interface{})

// Example usage:
logger.Infof("Service started", "port", 8080, "version", "1.0.0")
logger.Warnf("High memory", "usage", "85%", "threshold", "80%")
logger.Errorf("Database error", "error", "timeout", "retry_count", 3)
```

## 🎨 Log Level Styles

| Level | Emoji | Color Template | Usage |
|-------|-------|----------------|-------|
| Info  | ℹ️    | blue           | General information |
| Warn  | ⚠️    | orange         | Warning messages |
| Error | ❌    | red            | Error messages |

## 🎨 Available Color Templates

| Color | Usage |
|-------|-------|
| `blue` | Info messages, general notifications |
| `green` | Success messages, healthy status |
| `orange` | Warning messages, attention needed |
| `red` | Error messages, critical issues |
| `purple` | Debug messages, development info |
| `grey` | Default, neutral messages |

## 🧪 Testing

Run all tests:

```bash
go test -v ./...
```

Run tests with coverage report:

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📁 Project Structure

```
lark-logger/
├── .github/
│   └── workflows/          # GitHub Actions workflows
├── cmd/                    # Example application
├── docs/                   # Documentation
├── scripts/                # Utility scripts
├── src/
│   └── larklogger/         # Main package
│       ├── *.go           # Source files
│       └── *_test.go      # Test files
├── go.mod                  # Go module file
├── larklogger.go          # Main package exports
├── .golangci.yml          # Linter configuration
├── .goreleaser.yml         # GoReleaser configuration
├── Dockerfile              # Docker configuration
├── Makefile                # Build automation
├── env.example             # Environment variables example
└── README.md               # Project documentation
```

## 🔧 Development

### Prerequisites

- Go 1.19+
- Git

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/KCNyu/lark-logger.git
cd larklogger
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test -v ./src/larklogger/...
```

4. Run example:
```bash
go run ./cmd/main.go
```

5. Run with Makefile:
```bash
make run          # Run example with environment variables
make test         # Run tests
make lint         # Run linter
make ci           # Run full CI pipeline
```

### Code Quality

The project uses the following tools to ensure code quality:

- `golangci-lint` - Code static analysis
- `gosec` - Security scanning
- `go test` - Unit testing
- `go vet` - Code checking

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

If you encounter any issues or have suggestions, please:

- Submit an [Issue](https://github.com/KCNyu/lark-logger/issues)
- Send an email to [shengyu.li.evgeny@gmail.com](mailto:shengyu.li.evgeny@gmail.com)

## 🙏 Acknowledgments

Thanks to all developers who contributed to this project!

---

⭐ If this project helps you, please give it a star!