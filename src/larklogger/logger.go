package larklogger

import (
	"fmt"
	"strings"
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
	Service    string
	Env        string
	Hostname   string
	Title      string
	ShowConfig bool // Whether to show configuration section in logs
}

// LoggerOption is a function that configures the logger
type LoggerOption func(*LoggerConfig)

// NewLarkLogger creates a new LarkLogger instance
func NewLarkLogger(client *LarkClient, opts ...LoggerOption) Logger {
	config := &LoggerConfig{
		Service:    "default-service",
		Env:        "development",
		Hostname:   "localhost",
		Title:      "System Log",
		ShowConfig: false,
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

	// Add subtitle with message and level emoji
	subtitle := fmt.Sprintf("%s %s", getLogLevelEmoji(level), message)
	builder.AddSubtitle(subtitle)

	// Add timestamp
	builder.AddTimestamp()

	// Add divider
	builder.AddDivider()

	// Add custom fields if any
	if len(fields) > 0 {
		customFields := mapToKVItems(fields)
		builder.AddKVTable(customFields)
	}

	// Add configuration section only if ShowConfig is enabled
	if l.opts.ShowConfig {
		builder.AddDivider()

		// Add configuration section as 2x2 grid with emojis
		configData := map[string]string{
			"level":          "📊 Level",
			"level_value":    strings.ToUpper(string(level)),
			"service":        "🔧 Service",
			"service_value":  l.opts.Service,
			"env":            "🌍 Environment",
			"env_value":      l.opts.Env,
			"hostname":       "🖥️ Hostname",
			"hostname_value": l.opts.Hostname,
		}

		// Add config section as 2x2 grid
		builder.AddConfigGrid(configData)
	}

	return builder.Build()
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

func WithShowConfig(show bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.ShowConfig = show
	}
}
