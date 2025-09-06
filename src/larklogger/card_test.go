package larklogger

import (
	"testing"
	"time"
)

func TestCardBuilder(t *testing.T) {
	t.Run("build basic card", func(t *testing.T) {
		card := NewCardBuilder().
			SetHeader("Test Title", "blue").
			AddSection("Test content").
			Build()

		if card == nil {
			t.Error("Expected card to not be nil")
		}
		if card.Card.Config.WideScreenMode != true {
			t.Error("Expected WideScreenMode to be true")
		}
		if card.Card.Config.EnableForward != true {
			t.Error("Expected EnableForward to be true")
		}
		if card.Card.Config.UpdateMulti != true {
			t.Error("Expected UpdateMulti to be true")
		}

		if card.Card.Header.Title.Content != "Test Title" {
			t.Errorf("Expected title 'Test Title', got %s", card.Card.Header.Title.Content)
		}
		if card.Card.Header.Template != "blue" {
			t.Errorf("Expected template 'blue', got %s", card.Card.Header.Template)
		}

		if len(card.Card.Elements) < 1 {
			t.Errorf("Expected at least 1 element, got %d", len(card.Card.Elements))
		}
	})

	t.Run("build card with key value list", func(t *testing.T) {
		kv := map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}

		card := NewCardBuilder().
			SetHeader("KV Test", "green").
			AddKeyValueList("Test KV", kv).
			Build()

		if card == nil {
			t.Error("Expected card to not be nil")
		}
		if len(card.Card.Elements) < 2 {
			t.Errorf("Expected at least 2 elements, got %d", len(card.Card.Elements))
		}
	})

	t.Run("build card with status badge", func(t *testing.T) {
		card := NewCardBuilder().
			SetHeader("Status Test", "red").
			AddStatusBadge("success", "All systems operational").
			Build()

		if card == nil {
			t.Error("Expected card to not be nil")
		}
		if len(card.Card.Elements) < 1 {
			t.Errorf("Expected at least 1 element, got %d", len(card.Card.Elements))
		}
	})

	t.Run("build card with metrics grid", func(t *testing.T) {
		metrics := map[string]interface{}{
			"CPU":    "45%",
			"Memory": "67%",
			"Disk":   "23%",
		}

		card := NewCardBuilder().
			SetHeader("Metrics Test", "yellow").
			AddMetricsGrid("System Metrics", metrics).
			Build()

		if card == nil {
			t.Error("Expected card to not be nil")
		}
		if len(card.Card.Elements) < 1 {
			t.Errorf("Expected at least 1 element, got %d", len(card.Card.Elements))
		}
	})

	t.Run("build enhanced card with KV table", func(t *testing.T) {
		kvData := map[string]interface{}{
			"error": "connection timeout",
			"port":  8080,
			"host":  "localhost",
		}

		kvList := mapToKVItems(kvData)
		card := NewCardBuilder().
			SetHeader("Enhanced Test", "blue").
			AddSubtitle("Test subtitle").
			AddTimestamp().
			AddDivider().
			AddKVTable(kvList).
			Build()

		if card == nil {
			t.Error("Expected card to not be nil")
		}
		if len(card.Card.Elements) < 4 {
			t.Errorf("Expected at least 4 elements, got %d", len(card.Card.Elements))
		}
	})
}

func TestCardToJSON(t *testing.T) {
	card := NewCardBuilder().
		SetHeader("JSON Test", "blue").
		AddSection("Test content").
		Build()

	jsonStr, err := card.ToJSON()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if jsonStr == "" {
		t.Error("Expected non-empty JSON string")
	}

	// Check if JSON contains expected fields
	if !contains(jsonStr, "msg_type") {
		t.Error("Expected JSON to contain 'msg_type'")
	}
	if !contains(jsonStr, "card") {
		t.Error("Expected JSON to contain 'card'")
	}
	if !contains(jsonStr, "JSON Test") {
		t.Error("Expected JSON to contain 'JSON Test'")
	}
}

func TestNewCardField(t *testing.T) {
	field := NewCardField(true, "Test content")

	if field == nil {
		t.Error("Expected field to not be nil")
	}
	if !field.IsShort {
		t.Error("Expected IsShort to be true")
	}
	if field.Text == nil {
		t.Error("Expected Text to not be nil")
	}
	if field.Text.Content != "Test content" {
		t.Errorf("Expected content 'Test content', got %s", field.Text.Content)
	}
}

func TestFormatTimestamp(t *testing.T) {
	now := time.Now()
	formatted := FormatTimestamp(now)

	if formatted == "" {
		t.Error("Expected non-empty formatted timestamp")
	}

	// Check if it's in the expected format
	expectedFormat := "2006-01-02 15:04:05"
	_, err := time.Parse(expectedFormat, formatted)
	if err != nil {
		t.Errorf("Expected timestamp to be in format %s, got %s", expectedFormat, formatted)
	}
}

func TestGetLogLevelEmoji(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelInfo, "ℹ️"},
		{LevelWarn, "⚠️"},
		{LevelError, "❌"},
		{"unknown", "📋"},
	}

	for _, test := range tests {
		result := GetLogLevelEmoji(test.level)
		if result != test.expected {
			t.Errorf("Expected emoji %s for level %s, got %s", test.expected, test.level, result)
		}
	}
}
