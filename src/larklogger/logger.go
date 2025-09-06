package larklogger

import (
	"fmt"
	"strings"
	"time"
)

// LogLevel represents the log level
type LogLevel string

const (
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Logger interface defines the logging methods
type Logger interface {
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
	Infof(title string, args ...interface{})
	Warnf(title string, args ...interface{})
	Errorf(title string, args ...interface{})
}

// LarkLogger implements the Logger interface
type LarkLogger struct {
	client *LarkClient
	opts   *LoggerConfig
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Service  string
	Env      string
	Hostname string
	Title    string // Custom main title for log cards
}

// LoggerOption is a function that configures the logger
type LoggerOption func(*LoggerConfig)

// NewLarkLogger creates a new LarkLogger instance
func NewLarkLogger(client *LarkClient, opts ...LoggerOption) Logger {
	config := &LoggerConfig{
		Service:  "default-service",
		Env:      "development",
		Hostname: "localhost",
		Title:    "System Log",
	}

	for _, opt := range opts {
		opt(config)
	}

	return &LarkLogger{
		client: client,
		opts:   config,
	}
}

// Info logs an info level message
func (l *LarkLogger) Info(message string, fields map[string]interface{}) {
	l.log(LevelInfo, message, fields)
}

// Warn logs a warning level message
func (l *LarkLogger) Warn(message string, fields map[string]interface{}) {
	l.log(LevelWarn, message, fields)
}

// Error logs an error level message
func (l *LarkLogger) Error(message string, fields map[string]interface{}) {
	l.log(LevelError, message, fields)
}

// Infof logs an info level message with formatted title and key-value pairs
func (l *LarkLogger) Infof(title string, args ...interface{}) {
	fields := l.parseKeyValuePairs(args...)
	l.log(LevelInfo, title, fields)
}

// Warnf logs a warning level message with formatted title and key-value pairs
func (l *LarkLogger) Warnf(title string, args ...interface{}) {
	fields := l.parseKeyValuePairs(args...)
	l.log(LevelWarn, title, fields)
}

// Errorf logs an error level message with formatted title and key-value pairs
func (l *LarkLogger) Errorf(title string, args ...interface{}) {
	fields := l.parseKeyValuePairs(args...)
	l.log(LevelError, title, fields)
}

// parseKeyValuePairs parses alternating key-value pairs from args
func (l *LarkLogger) parseKeyValuePairs(args ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				key = fmt.Sprintf("field_%d", i/2)
			}
			fields[key] = args[i+1]
		} else {
			// Odd number of args, treat last one as a message
			key := fmt.Sprintf("extra_%d", i/2)
			fields[key] = args[i]
		}
	}

	return fields
}

// log sends a log message to Lark
func (l *LarkLogger) log(level LogLevel, message string, fields map[string]interface{}) {
	card := l.buildLogCard(level, message, fields)
	if err := l.client.SendCard(card); err != nil {
		// In a real implementation, you might want to fallback to console logging
		fmt.Printf("Failed to send log to Lark: %v\n", err)
	}
}

// buildLogCard builds a Lark card for the log message using enhanced design
func (l *LarkLogger) buildLogCard(level LogLevel, message string, fields map[string]interface{}) *Card {
	emoji := GetLogLevelEmoji(level)
	template := getVisualConfig(level)

	// Build main title with custom title and emoji
	mainTitle := fmt.Sprintf("%s %s", emoji, l.opts.Title)

	// Create enhanced card builder
	builder := NewCardBuilder().SetHeader(mainTitle, template)

	// Add subtitle with message only (no redundant default text)
	builder.AddSubtitle(message)

	// Add timestamp
	builder.AddTimestamp()

	// Add divider
	builder.AddDivider()

	// Prepare KV data
	kvData := map[string]interface{}{
		"Level": strings.ToUpper(string(level)),
		"Service": l.opts.Service,
		"Environment": l.opts.Env,
		"Hostname": l.opts.Hostname,
	}

	// Add custom fields to KV data
	for key, value := range fields {
		kvData[key] = value
	}

	// Convert to KV list
	kvList := mapToKVItems(kvData)

	// Add KV table
	builder.AddKVTable(kvList)

	// Add optional card link (you can customize this)
	// builder.AddCardLink("https://logs.example.com")

	return builder.Build()
}

// getStatusText returns status text based on log level
func (l *LarkLogger) getStatusText(level LogLevel) string {
	switch level {
	case LevelInfo:
		return "System running normally"
	case LevelWarn:
		return "Performance degradation detected"
	case LevelError:
		return "Critical error occurred"
	default:
		return "Unknown status"
	}
}

// getFieldEmoji returns appropriate emoji for field key
func (l *LarkLogger) getFieldEmoji(key string) string {
	keyLower := strings.ToLower(key)

	switch {
	case strings.Contains(keyLower, "error") || strings.Contains(keyLower, "exception"):
		return "❌"
	case strings.Contains(keyLower, "warning") || strings.Contains(keyLower, "warn"):
		return "⚠️"
	case strings.Contains(keyLower, "success") || strings.Contains(keyLower, "ok"):
		return "✅"
	case strings.Contains(keyLower, "time") || strings.Contains(keyLower, "duration"):
		return "⏱️"
	case strings.Contains(keyLower, "memory") || strings.Contains(keyLower, "ram"):
		return "💾"
	case strings.Contains(keyLower, "cpu") || strings.Contains(keyLower, "processor"):
		return "🖥️"
	case strings.Contains(keyLower, "network") || strings.Contains(keyLower, "connection"):
		return "🌐"
	case strings.Contains(keyLower, "database") || strings.Contains(keyLower, "db"):
		return "🗄️"
	case strings.Contains(keyLower, "user") || strings.Contains(keyLower, "client"):
		return "👤"
	case strings.Contains(keyLower, "request") || strings.Contains(keyLower, "api"):
		return "📡"
	case strings.Contains(keyLower, "file") || strings.Contains(keyLower, "path"):
		return "📁"
	case strings.Contains(keyLower, "port") || strings.Contains(keyLower, "address"):
		return "🔌"
	case strings.Contains(keyLower, "version") || strings.Contains(keyLower, "v"):
		return "🏷️"
	case strings.Contains(keyLower, "count") || strings.Contains(keyLower, "number"):
		return "🔢"
	case strings.Contains(keyLower, "size") || strings.Contains(keyLower, "length"):
		return "📏"
	case strings.Contains(keyLower, "status") || strings.Contains(keyLower, "state"):
		return "📊"
	default:
		return "📋"
	}
}

// formatFieldValue formats field value for display
func (l *LarkLogger) formatFieldValue(value interface{}) string {
	if value == nil {
		return "-"
	}

	switch v := value.(type) {
	case string:
		if len(v) > 100 {
			return v[:97] + "..."
		}
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Logger option functions
func WithService(service string) LoggerOption {
	return func(c *LoggerConfig) {
		c.Service = service
	}
}

func WithEnv(env string) LoggerOption {
	return func(c *LoggerConfig) {
		c.Env = env
	}
}

func WithHostname(hostname string) LoggerOption {
	return func(c *LoggerConfig) {
		c.Hostname = hostname
	}
}

func WithTitle(title string) LoggerOption {
	return func(c *LoggerConfig) {
		c.Title = title
	}
}
