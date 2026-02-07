// Package api provides the Todoist REST API v2 client.
package api

import (
	"encoding/json"
	"errors"
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
		return nil, errors.New("access token is required")
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
func (c *Client) request(method, endpoint string, body io.Reader) ([]byte, int, error) {
	var lastErr error
	backoff := InitialBackoff

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(backoff)
			backoff *= 2
		}

		url := c.baseURL + endpoint
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Via-Todoist/0.1")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response: %w", err)
			continue
		}

		// Handle rate limiting (429)
		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = NewRateLimitError(60)
			time.Sleep(60 * time.Second)
			continue
		}

		// Handle non-2xx status codes
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			apiErr := &TodoistError{
				StatusCode: resp.StatusCode,
			}

			// Try to parse error detail from response body
			if len(respBody) > 0 {
				var errMsg string
				if err := json.Unmarshal(respBody, &errMsg); err == nil {
					apiErr.Message = errMsg
				} else {
					apiErr.Message = string(respBody)
				}
			} else {
				apiErr.Message = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
			}

			if resp.StatusCode == http.StatusUnauthorized {
				return nil, resp.StatusCode, NewAuthError()
			}

			if apiErr.IsServerError() {
				lastErr = apiErr
				continue
			}

			return nil, resp.StatusCode, apiErr
		}

		return respBody, resp.StatusCode, nil
	}

	if lastErr != nil {
		return nil, 0, fmt.Errorf("request failed after %d retries: %w", MaxRetries, lastErr)
	}
	return nil, 0, fmt.Errorf("request failed after %d retries", MaxRetries)
}
