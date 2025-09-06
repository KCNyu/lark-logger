package larklogger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewLarkClient(t *testing.T) {
	t.Run("with default options", func(t *testing.T) {
		client := NewLarkClient("https://example.com/webhook")

		if client.webhookURL != "https://example.com/webhook" {
			t.Errorf("Expected webhook URL %s, got %s", "https://example.com/webhook", client.webhookURL)
		}

		if client.opts.Timeout != 30*time.Second {
			t.Errorf("Expected timeout %v, got %v", 30*time.Second, client.opts.Timeout)
		}

		if client.opts.RetryCount != 3 {
			t.Errorf("Expected retry count %d, got %d", 3, client.opts.RetryCount)
		}

		if client.opts.UserAgent != "larklogger-go/1.0.0" {
			t.Errorf("Expected user agent %s, got %s", "larklogger-go/1.0.0", client.opts.UserAgent)
		}
	})

	t.Run("with custom options", func(t *testing.T) {
		client := NewLarkClient("https://example.com/webhook",
			WithTimeout(10*time.Second),
			WithRetry(5, 2*time.Second),
			WithUserAgent("custom-agent/1.0"),
			WithHeaders(map[string]string{"X-Custom": "value"}),
		)

		if client.opts.Timeout != 10*time.Second {
			t.Errorf("Expected timeout %v, got %v", 10*time.Second, client.opts.Timeout)
		}

		if client.opts.RetryCount != 5 {
			t.Errorf("Expected retry count %d, got %d", 5, client.opts.RetryCount)
		}

		if client.opts.RetryDelay != 2*time.Second {
			t.Errorf("Expected retry delay %v, got %v", 2*time.Second, client.opts.RetryDelay)
		}

		if client.opts.UserAgent != "custom-agent/1.0" {
			t.Errorf("Expected user agent %s, got %s", "custom-agent/1.0", client.opts.UserAgent)
		}

		if client.opts.Headers["X-Custom"] != "value" {
			t.Errorf("Expected custom header value 'value', got %s", client.opts.Headers["X-Custom"])
		}
	})
}

func TestLarkClientSendText(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify user agent
		if r.Header.Get("User-Agent") != "larklogger-go/1.0.0" {
			t.Errorf("Expected User-Agent larklogger-go/1.0.0, got %s", r.Header.Get("User-Agent"))
		}

		// Parse request body
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		// Verify payload structure
		if payload["msg_type"] != "text" {
			t.Errorf("Expected msg_type 'text', got %v", payload["msg_type"])
		}

		content, ok := payload["content"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected content to be a map")
			return
		}

		if content["text"] != "Hello, Lark!" {
			t.Errorf("Expected text 'Hello, Lark!', got %v", content["text"])
		}

		// Send success response
		response := map[string]interface{}{
			"code": 0,
			"msg":  "success",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client and send message
	client := NewLarkClient(server.URL)
	err := client.SendText("Hello, Lark!")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLarkClientSendCard(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
			return
		}

		// Verify payload structure
		if payload["msg_type"] != "interactive" {
			t.Errorf("Expected msg_type 'interactive', got %v", payload["msg_type"])
		}

		card, ok := payload["card"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected card to be a map")
			return
		}

		// Verify card has required fields
		if card["elements"] == nil {
			t.Errorf("Expected card to have elements")
		}

		// Send success response
		response := map[string]interface{}{
			"code": 0,
			"msg":  "success",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client and send card
	client := NewLarkClient(server.URL)
	card := NewCardBuilder().
		SetHeader("Test Card", "blue").
		AddSection("Test content").
		Build()

	err := client.SendCard(card)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestLarkClientErrorHandling(t *testing.T) {
	t.Run("server returns error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"code": 1,
				"msg":  "invalid webhook",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewLarkClient(server.URL)
		err := client.SendText("test")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if !contains(err.Error(), "lark API error") {
			t.Errorf("Expected error to contain 'lark API error', got %v", err)
		}
	})

	t.Run("server returns non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Bad Request"))
		}))
		defer server.Close()

		client := NewLarkClient(server.URL)
		err := client.SendText("test")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if !contains(err.Error(), "status 400") {
			t.Errorf("Expected error to contain 'status 400', got %v", err)
		}
	})
}

func TestClientOptions(t *testing.T) {
	t.Run("WithTimeout", func(t *testing.T) {
		opts := &ClientOptions{}
		WithTimeout(5 * time.Second)(opts)
		if opts.Timeout != 5*time.Second {
			t.Errorf("Expected timeout %v, got %v", 5*time.Second, opts.Timeout)
		}
	})

	t.Run("WithRetry", func(t *testing.T) {
		opts := &ClientOptions{}
		WithRetry(2, 500*time.Millisecond)(opts)
		if opts.RetryCount != 2 {
			t.Errorf("Expected retry count %d, got %d", 2, opts.RetryCount)
		}
		if opts.RetryDelay != 500*time.Millisecond {
			t.Errorf("Expected retry delay %v, got %v", 500*time.Millisecond, opts.RetryDelay)
		}
	})

	t.Run("WithUserAgent", func(t *testing.T) {
		opts := &ClientOptions{}
		WithUserAgent("test-agent")(opts)
		if opts.UserAgent != "test-agent" {
			t.Errorf("Expected user agent 'test-agent', got %s", opts.UserAgent)
		}
	})

	t.Run("WithHeaders", func(t *testing.T) {
		opts := &ClientOptions{}
		headers := map[string]string{"X-Test": "value"}
		WithHeaders(headers)(opts)
		if opts.Headers["X-Test"] != "value" {
			t.Errorf("Expected header value 'value', got %s", opts.Headers["X-Test"])
		}
	})
}
