package api

import (
	"errors"
	"fmt"
	"net/http"
)

// TodoistError represents an error from the Todoist API.
type TodoistError struct {
	Message    string
	StatusCode int
}

func (e *TodoistError) Error() string {
	return fmt.Sprintf("[Todoist] %s (HTTP %d)", e.Message, e.StatusCode)
}

// IsAuthError returns true if the error is an authentication error (401).
func (e *TodoistError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsRateLimitError returns true if the error is a rate limit error (429).
func (e *TodoistError) IsRateLimitError() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true if the error is a server error (5xx).
func (e *TodoistError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// IsNotFoundError returns true if the error is a not found error (404).
func (e *TodoistError) IsNotFoundError() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsRetryable returns true if the error is potentially retryable (429 or 5xx).
func (e *TodoistError) IsRetryable() bool {
	return e.IsRateLimitError() || e.IsServerError()
}

// Package-level helper functions for error checking via errors.As.

// IsTodoistError returns true if the error is a TodoistError.
func IsTodoistError(err error) bool {
	var apiErr *TodoistError
	return errors.As(err, &apiErr)
}

// IsAuthError returns true if the error is a Todoist authentication error.
func IsAuthError(err error) bool {
	var apiErr *TodoistError
	if errors.As(err, &apiErr) {
		return apiErr.IsAuthError()
	}
	return false
}

// IsRateLimitError returns true if the error is a Todoist rate limit error.
func IsRateLimitError(err error) bool {
	var apiErr *TodoistError
	if errors.As(err, &apiErr) {
		return apiErr.IsRateLimitError()
	}
	return false
}

// IsServerError returns true if the error is a Todoist server error (5xx).
func IsServerError(err error) bool {
	var apiErr *TodoistError
	if errors.As(err, &apiErr) {
		return apiErr.IsServerError()
	}
	return false
}

// IsNotFoundError returns true if the error is a Todoist not found error.
func IsNotFoundError(err error) bool {
	var apiErr *TodoistError
	if errors.As(err, &apiErr) {
		return apiErr.IsNotFoundError()
	}
	return false
}

// IsRetryable returns true if the error is potentially retryable.
func IsRetryable(err error) bool {
	var apiErr *TodoistError
	if errors.As(err, &apiErr) {
		return apiErr.IsRetryable()
	}
	return false
}

// NewAuthError creates a new authentication error.
func NewAuthError() *TodoistError {
	return &TodoistError{
		Message:    "Unauthorized: Invalid or missing access token",
		StatusCode: http.StatusUnauthorized,
	}
}

// NewRateLimitError creates a new rate limit error.
func NewRateLimitError(retryAfter int) *TodoistError {
	return &TodoistError{
		Message:    fmt.Sprintf("Rate limit exceeded. Retry after %d seconds", retryAfter),
		StatusCode: http.StatusTooManyRequests,
	}
}
