package main

import (
	"fmt"
	"log"

	larklogger "github.com/KCNyu/lark-logger"
)

func main() {
	// Get webhook URL from environment variable
	webhookURL := larklogger.GetWebhookURL()

	// Check if running in test mode
	if larklogger.IsTestEnvironment() {
		log.Println("⚠️  Running in TEST MODE - messages will not be sent to real Lark webhook")
		log.Printf("Test webhook URL: %s", webhookURL)
	} else {
		log.Println("🚀 Running in PRODUCTION MODE - messages will be sent to Lark")
	}

	// Create logger with webhook URL and options
	logger := larklogger.New(
		webhookURL,
		larklogger.WithEnv("production"),
		larklogger.WithTitle("System Monitor"),
	)

	// Send Info level log with rich context (traditional way)
	logger.Info("API Gateway initialized successfully", map[string]interface{}{
		"port":     8080,
		"version":  "2.1.0",
		"features": []string{"authentication", "rate_limiting", "metrics"},
		"config": map[string]interface{}{
			"database": "postgresql",
			"cache":    "redis",
			"ssl":      true,
		},
		"uptime": "0s",
	})

	// Create a logger with configuration section enabled
	configLogger := larklogger.New(
		webhookURL,
		larklogger.WithService("config-demo"),
		larklogger.WithEnv("production"),
		larklogger.WithHostname("server-01"),
		larklogger.WithTitle("Config Demo"),
		larklogger.WithShowConfig(true), // Enable configuration section
	)

	// This will show the configuration section with light blue background
	configLogger.Info("Configuration section enabled", map[string]interface{}{
		"feature": "config_visibility",
		"status":  "enabled",
	})

	// Create a logger with buttons for action items
	buttonLogger := larklogger.New(
		webhookURL,
		larklogger.WithService("action-demo"),
		larklogger.WithEnv("production"),
		larklogger.WithHostname("server-01"),
		larklogger.WithTitle("Action Required"),
		larklogger.WithShowConfig(true),
		larklogger.WithButtons([]larklogger.Button{
			{
				Text:  "View Logs",
				URL:   "https://logs.example.com/system/error-123",
				Style: larklogger.ButtonStylePrimary,
			},
			{
				Text:    "Restart Service",
				URL:     "https://ops.example.com/restart/service-abc",
				Style:   larklogger.ButtonStyleSecondary,
				Confirm: true,
			},
			{
				Text:    "Escalate Issue",
				URL:     "https://ops.example.com/escalate/incident-456",
				Style:   larklogger.ButtonStyleDanger,
				Confirm: true,
			},
		}),
	)

	// This will show buttons for user actions
	buttonLogger.Error("Critical system error detected", map[string]interface{}{
		"error_code": "SYS_001",
		"component":  "database",
		"severity":   "critical",
		"timestamp":  "2024-01-15T10:30:00Z",
		"details":    "Connection pool exhausted, unable to process requests",
	})

	// Send logs using new Infof/Warnf/Errorf format (simplified)
	logger.Infof("Service health check", "status", "healthy", "response_time", "120ms", "uptime", "2h30m")
	logger.Warnf("Memory usage approaching threshold", "usage", "87%", "threshold", "85%", "recommendation", "consider horizontal scaling")
	logger.Errorf("Database connection pool exhausted", "error", "connection timeout after 30s", "retry_count", 3, "pool_size", 20)

	// Send structured message using logger
	logger.Info("🚀 Lark Logger Demo - All systems operational!", map[string]interface{}{
		"demo":   true,
		"status": "operational",
	})

	// Send system health report using structured logging
	logger.Info("System Health Report", map[string]interface{}{
		"Memory":        "67%",
		"Disk Space":    "23%",
		"Network I/O":   "12%",
		"Response Time": "120ms",
		"Report Type":   "Health Check",
	})

	fmt.Println("🎊 All messages sent successfully! Check your Lark workspace.")
}
