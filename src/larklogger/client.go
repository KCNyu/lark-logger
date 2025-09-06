package larklogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LarkClient handles communication with Lark webhook
type LarkClient struct {
	webhookURL string
	httpClient *http.Client
	opts       *ClientOptions
}

// ClientOptions holds client configuration options
type ClientOptions struct {
	Timeout    time.Duration
	RetryCount int
	RetryDelay time.Duration
	UserAgent  string
	Headers    map[string]string
}

// ClientOption is a function that configures the client
type ClientOption func(*ClientOptions)

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.Timeout = timeout
	}
}

// WithRetry sets retry configuration
func WithRetry(count int, delay time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.RetryCount = count
		opts.RetryDelay = delay
	}
}

// WithUserAgent sets the user agent
func WithUserAgent(userAgent string) ClientOption {
	return func(opts *ClientOptions) {
		opts.UserAgent = userAgent
	}
}

// WithHeaders sets custom headers
func WithHeaders(headers map[string]string) ClientOption {
	return func(opts *ClientOptions) {
		opts.Headers = headers
	}
}

// NewLarkClient creates a new Lark client
func NewLarkClient(webhookURL string, opts ...ClientOption) *LarkClient {
	options := &ClientOptions{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: 1 * time.Second,
		UserAgent:  "larklogger-go/1.0.0",
		Headers:    make(map[string]string),
	}

	for _, opt := range opts {
		opt(options)
	}

	return &LarkClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: options.Timeout,
		},
		opts: options,
	}
}

// SendCard sends a card to the Lark webhook
func (c *LarkClient) SendCard(card *Card) error {
	jsonData, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to marshal card: %w", err)
	}

	return c.sendWithRetry(jsonData)
}

// SendText sends a simple text message to the Lark webhook
func (c *LarkClient) SendText(text string) error {
	payload := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": text,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return c.sendWithRetry(jsonData)
}

// sendWithRetry sends the request with retry logic
func (c *LarkClient) sendWithRetry(data []byte) error {
	var lastErr error

	for i := 0; i <= c.opts.RetryCount; i++ {
		if i > 0 {
			time.Sleep(c.opts.RetryDelay)
		}

		err := c.sendRequest(data)
		if err == nil {
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("failed to send message after %d retries: %w", c.opts.RetryCount, lastErr)
}

// sendRequest sends a single HTTP request
func (c *LarkClient) sendRequest(data []byte) error {
	req, err := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.opts.UserAgent)

	for key, value := range c.opts.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to check for errors
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if code, ok := response["code"].(float64); ok && code != 0 {
		msg := "unknown error"
		if message, ok := response["msg"].(string); ok {
			msg = message
		}
		return fmt.Errorf("lark API error (code: %.0f): %s", code, msg)
	}

	return nil
}
