package larklogger

import (
	"os"
	"strings"
)

// EnvConfig holds environment-based configuration
type EnvConfig struct {
	WebhookURL string
	IsTestMode bool
}

// GetConfig returns configuration based on environment variables
func GetConfig() *EnvConfig {
	webhookURL := os.Getenv("LARK_WEBHOOK_URL")
	isTestMode := strings.ToLower(os.Getenv("LARK_TEST_MODE")) == "true"

	// If no webhook URL is provided, use a test URL
	if webhookURL == "" {
		webhookURL = "https://test.webhook.url"
		isTestMode = true
	}

	return &EnvConfig{
		WebhookURL: webhookURL,
		IsTestMode: isTestMode,
	}
}

// GetWebhookURL returns the webhook URL from environment or default test URL
func GetWebhookURL() string {
	config := GetConfig()
	return config.WebhookURL
}

// IsTestEnvironment returns true if running in test mode
func IsTestEnvironment() bool {
	config := GetConfig()
	return config.IsTestMode
}
