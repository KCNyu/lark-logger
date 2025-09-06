// Package larklogger provides a Go SDK for sending structured log messages
// to Lark (Feishu) webhook bots with beautiful interactive cards.
//
// Features:
// - Support for Info, Warn, and Error log levels
// - Beautiful card formatting with emojis and colors
// - Structured logging with map[string]interface{} fields
// - Configurable retry logic and timeouts
// - Type-safe configuration options
//
// Example usage:
//
//	client := larklogger.NewClient("https://open.feishu.cn/open-apis/bot/v2/hook/xxx")
//	logger := larklogger.NewLogger(client)
//
//	logger.Info("Service started successfully", map[string]interface{}{
//		"port": 8080,
//		"version": "1.0.0",
//	})
//
//	logger.Error("Database connection failed", map[string]interface{}{
//		"error": "connection timeout",
//		"retry_count": 3,
//	})
package larklogger

// Re-export all public types and functions from internal package
import (
	"time"

	"github.com/KCNyu/lark-logger/src/larklogger"
)

// Logger interface defines the logging methods
type Logger = larklogger.Logger

// LogLevel represents the log level
type LogLevel = larklogger.LogLevel

// Client represents the Lark webhook client
type Client = larklogger.LarkClient

// Card represents a Lark interactive card
type Card = larklogger.Card

// CardBuilder helps build Lark cards
type CardBuilder = larklogger.CardBuilder

// KVItem represents a prioritized key-value item
type KVItem = larklogger.KVItem

// Client options
type ClientOption = larklogger.ClientOption

// Logger options
type LoggerOption = larklogger.LoggerOption

// Log levels
const (
	LevelInfo  = larklogger.LevelInfo
	LevelWarn  = larklogger.LevelWarn
	LevelError = larklogger.LevelError
)

// NewClient creates a new Lark client
func NewClient(webhookURL string, opts ...ClientOption) *Client {
	return larklogger.NewLarkClient(webhookURL, opts...)
}

// NewLogger creates a new Logger instance
func NewLogger(client *Client, opts ...LoggerOption) Logger {
	return larklogger.NewLarkLogger(client, opts...)
}

// New creates a new Logger with webhook URL and options (simplified API)
func New(webhookURL string, opts ...interface{}) Logger {
	// Separate client and logger options
	var clientOpts []ClientOption
	var loggerOpts []LoggerOption

	for _, opt := range opts {
		switch o := opt.(type) {
		case ClientOption:
			clientOpts = append(clientOpts, o)
		case LoggerOption:
			loggerOpts = append(loggerOpts, o)
		}
	}

	client := NewClient(webhookURL, clientOpts...)
	return NewLogger(client, loggerOpts...)
}

// NewCardBuilder creates a new card builder
func NewCardBuilder() *CardBuilder {
	return larklogger.NewCardBuilder()
}

// NewCardField creates a new card field
func NewCardField(isShort bool, content string) *larklogger.CardField {
	return larklogger.NewCardField(isShort, content)
}

// FormatTimestamp formats timestamp for display
func FormatTimestamp(t time.Time) string {
	return larklogger.FormatTimestamp(t)
}

// Client options
func WithTimeout(timeout time.Duration) ClientOption {
	return larklogger.WithTimeout(timeout)
}

func WithRetry(count int, delay time.Duration) ClientOption {
	return larklogger.WithRetry(count, delay)
}

func WithUserAgent(userAgent string) ClientOption {
	return larklogger.WithUserAgent(userAgent)
}

func WithHeaders(headers map[string]string) ClientOption {
	return larklogger.WithHeaders(headers)
}

// Logger options
func WithService(service string) LoggerOption {
	return larklogger.WithService(service)
}

func WithEnv(env string) LoggerOption {
	return larklogger.WithEnv(env)
}

func WithHostname(hostname string) LoggerOption {
	return larklogger.WithHostname(hostname)
}

func WithTitle(title string) LoggerOption {
	return larklogger.WithTitle(title)
}

func WithShowConfig(show bool) LoggerOption {
	return larklogger.WithShowConfig(show)
}

// Environment configuration functions
func GetWebhookURL() string {
	return larklogger.GetWebhookURL()
}

func IsTestEnvironment() bool {
	return larklogger.IsTestEnvironment()
}
