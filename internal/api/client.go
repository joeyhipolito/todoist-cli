// Package api provides the Todoist REST API v2 client.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// BaseURL is the Todoist REST API v2 base URL.
	BaseURL = "https://api.todoist.com/rest/v2"

	// MaxRetries is the maximum number of retry attempts.
	MaxRetries = 3

	// InitialBackoff is the initial backoff duration.
	InitialBackoff = 1 * time.Second
)

// Client is the Todoist API client.
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Todoist API client.
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, fmt.Errorf("access token is required")
	}

	return &Client{
		token:   token,
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// request performs an HTTP request with retry logic and rate limit handling.
// Body bytes are captured up front so the request can be retried safely.
func (c *Client) request(method, endpoint string, body io.Reader) ([]byte, int, error) {
	// Read body once so we can replay it on retries.
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	var lastErr error
	backoff := InitialBackoff

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff: 1s, 2s, 4s
		}

		url := c.baseURL + endpoint

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(method, url, reqBody)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Via-Todoist/0.1")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue // Retry on network errors
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		// Handle rate limiting (429) — retry after backoff
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = NewRateLimitError(60)
			time.Sleep(60 * time.Second)
			continue
		}

		// Handle non-2xx status codes
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			// Fail-fast on 401 — no point retrying with a bad token
			if resp.StatusCode == http.StatusUnauthorized {
				return nil, resp.StatusCode, NewAuthError()
			}

			apiErr := &TodoistError{
				StatusCode: resp.StatusCode,
			}

			// Todoist error responses may be a plain string or {"error": "..."}
			if len(respBody) > 0 {
				apiErr.Message = parseErrorBody(respBody)
			} else {
				apiErr.Message = http.StatusText(resp.StatusCode)
			}

			// Retry server errors (5xx) with exponential backoff
			if apiErr.IsServerError() {
				lastErr = apiErr
				continue
			}

			// Fail-fast on other 4xx (400, 403, 404, etc.)
			return nil, resp.StatusCode, apiErr
		}

		return respBody, resp.StatusCode, nil
	}

	if lastErr != nil {
		return nil, 0, fmt.Errorf("request failed after %d retries: %w", MaxRetries, lastErr)
	}
	return nil, 0, fmt.Errorf("request failed after %d retries", MaxRetries)
}

// parseErrorBody extracts a message from a Todoist error response.
// The API may return a JSON string, {"error": "..."}, or plain text.
func parseErrorBody(body []byte) string {
	// Try {"error": "message"} format
	var structured struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &structured); err == nil && structured.Error != "" {
		return structured.Error
	}

	// Try bare JSON string (e.g., "\"error message\"")
	var plain string
	if err := json.Unmarshal(body, &plain); err == nil && plain != "" {
		return plain
	}

	// Fall back to raw body text
	return string(body)
}
