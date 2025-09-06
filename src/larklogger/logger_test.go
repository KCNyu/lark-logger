package larklogger

import (
	"os"
	"testing"
)

func TestNewLarkLogger(t *testing.T) {
	client := NewLarkClient("https://example.com/webhook")

	t.Run("with default options", func(t *testing.T) {
		logger := NewLarkLogger(client)
		larkLogger := logger.(*LarkLogger)

		if larkLogger.opts.Service != "default-service" {
			t.Errorf("Expected service 'default-service', got %s", larkLogger.opts.Service)
		}
		if larkLogger.opts.Env != "development" {
			t.Errorf("Expected env 'development', got %s", larkLogger.opts.Env)
		}
		if larkLogger.opts.Hostname != "localhost" {
			t.Errorf("Expected hostname 'localhost', got %s", larkLogger.opts.Hostname)
		}
	})

	t.Run("with custom options", func(t *testing.T) {
		logger := NewLarkLogger(client,
			WithService("test-service"),
			WithEnv("production"),
			WithHostname("server-01"),
			WithTitle("Custom Log Title"),
		)
		larkLogger := logger.(*LarkLogger)

		if larkLogger.opts.Service != "test-service" {
			t.Errorf("Expected service 'test-service', got %s", larkLogger.opts.Service)
		}
		if larkLogger.opts.Env != "production" {
			t.Errorf("Expected env 'production', got %s", larkLogger.opts.Env)
		}
		if larkLogger.opts.Hostname != "server-01" {
			t.Errorf("Expected hostname 'server-01', got %s", larkLogger.opts.Hostname)
		}
		if larkLogger.opts.Title != "Custom Log Title" {
			t.Errorf("Expected title 'Custom Log Title', got %s", larkLogger.opts.Title)
		}
	})
}

func TestLarkLoggerBuildLogCard(t *testing.T) {
	client := NewLarkClient("https://example.com/webhook")
	logger := NewLarkLogger(client, WithTitle("System Log"))

	t.Run("info level card", func(t *testing.T) {
		fields := map[string]interface{}{
			"port":    8080,
			"version": "1.0.0",
		}
		larkLogger := logger.(*LarkLogger)
		card := larkLogger.buildLogCard(LevelInfo, "Service started", fields)

		if card == nil {
			t.Error("Expected card to not be nil")
			return
		}
		if card.Card.Header.Title.Content == "" {
			t.Error("Expected header title content to not be empty")
			return
		}
		if card.Card.Header.Title.Content == "" {
			t.Error("Expected card header title to not be empty")
		}
		if !contains(card.Card.Header.Title.Content, "ℹ️") {
			t.Error("Expected title to contain ℹ️")
		}
		if !contains(card.Card.Header.Title.Content, "System Log") {
			t.Error("Expected title to contain System Log")
		}
		if card.Card.Header.Template != "blue" {
			t.Errorf("Expected template 'blue', got %s", card.Card.Header.Template)
		}

		if len(card.Card.Elements) < 2 {
			t.Errorf("Expected at least 2 elements, got %d", len(card.Card.Elements))
		}
	})

	t.Run("error level card", func(t *testing.T) {
		larkLogger := logger.(*LarkLogger)
		card := larkLogger.buildLogCard(LevelError, "Database error", nil)

		if card == nil {
			t.Error("Expected card to not be nil")
			return
		}
		if card.Card.Header.Title.Content == "" {
			t.Error("Expected header title content to not be empty")
			return
		}
		if !contains(card.Card.Header.Title.Content, "❌") {
			t.Error("Expected title to contain ❌")
		}
		if !contains(card.Card.Header.Title.Content, "System Log") {
			t.Error("Expected title to contain System Log")
		}
		if card.Card.Header.Template != "red" {
			t.Errorf("Expected template 'red', got %s", card.Card.Header.Template)
		}
	})

	t.Run("warn level card", func(t *testing.T) {
		larkLogger := logger.(*LarkLogger)
		card := larkLogger.buildLogCard(LevelWarn, "High memory usage", nil)

		if card == nil {
			t.Error("Expected card to not be nil")
			return
		}
		if card.Card.Header.Title.Content == "" {
			t.Error("Expected header title content to not be empty")
			return
		}
		if !contains(card.Card.Header.Title.Content, "⚠️") {
			t.Error("Expected title to contain ⚠️")
		}
		if !contains(card.Card.Header.Title.Content, "System Log") {
			t.Error("Expected title to contain System Log")
		}
		if card.Card.Header.Template != "yellow" {
			t.Errorf("Expected template 'yellow', got %s", card.Card.Header.Template)
		}
	})
}

func TestLoggerOptions(t *testing.T) {
	t.Run("WithService", func(t *testing.T) {
		config := &LoggerConfig{}
		WithService("my-service")(config)
		if config.Service != "my-service" {
			t.Errorf("Expected service 'my-service', got %s", config.Service)
		}
	})

	t.Run("WithEnv", func(t *testing.T) {
		config := &LoggerConfig{}
		WithEnv("production")(config)
		if config.Env != "production" {
			t.Errorf("Expected env 'production', got %s", config.Env)
		}
	})

	t.Run("WithHostname", func(t *testing.T) {
		config := &LoggerConfig{}
		WithHostname("server-01")(config)
		if config.Hostname != "server-01" {
			t.Errorf("Expected hostname 'server-01', got %s", config.Hostname)
		}
	})

	t.Run("WithTitle", func(t *testing.T) {
		config := &LoggerConfig{}
		WithTitle("My Custom Title")(config)
		if config.Title != "My Custom Title" {
			t.Errorf("Expected title 'My Custom Title', got %s", config.Title)
		}
	})
}

func TestLoggerInfof(t *testing.T) {
	// Skip tests that use mock webhook URLs
	if os.Getenv("LARK_TEST_MODE") == "true" {
		t.Skip("Skipping mock webhook tests")
	}
	
	client := NewLarkClient("https://test.webhook.url")
	logger := NewLarkLogger(client, WithTitle("Test Logger"))

	t.Run("infof with key-value pairs", func(t *testing.T) {
		// This will test the Infof method
		// Note: In a real test, you might want to mock the client
		logger.Infof("Service started", "port", 8080, "version", "1.0.0")
	})

	t.Run("warnf with key-value pairs", func(t *testing.T) {
		logger.Warnf("High memory usage", "usage", "85%", "threshold", "80%")
	})

	t.Run("errorf with key-value pairs", func(t *testing.T) {
		logger.Errorf("Database error", "error", "connection timeout", "retry_count", 3)
	})
}
