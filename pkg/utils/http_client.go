package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

// DefaultRetryConfig returns sensible defaults for retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		BackoffFactor:  2.0,
	}
}

// RateLimitedClient wraps an HTTP client with rate limiting and retry logic
type RateLimitedClient struct {
	client      *http.Client
	rateLimiter *rate.Limiter
	retryConfig RetryConfig
}

// NewRateLimitedClient creates a new rate-limited HTTP client
func NewRateLimitedClient(requestsPerSecond float64, burst int, timeout time.Duration) *RateLimitedClient {
	return &RateLimitedClient{
		client: &http.Client{
			Timeout: timeout,
		},
		rateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		retryConfig: DefaultRetryConfig(),
	}
}

// NewRateLimitedClientWithConfig creates a new rate-limited HTTP client with custom retry config
func NewRateLimitedClientWithConfig(requestsPerSecond float64, burst int, timeout time.Duration, retryConfig RetryConfig) *RateLimitedClient {
	return &RateLimitedClient{
		client: &http.Client{
			Timeout: timeout,
		},
		rateLimiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		retryConfig: retryConfig,
	}
}

// Do executes an HTTP request with rate limiting and retry logic
func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoWithContext(req.Context(), req)
}

// DoWithContext executes an HTTP request with rate limiting, retry logic, and context support
func (c *RateLimitedClient) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	backoff := c.retryConfig.InitialBackoff

	for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
		// Wait for rate limiter
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}

		// Clone the request for each attempt
		reqClone := req.Clone(ctx)

		// Execute request
		resp, err = c.client.Do(reqClone)

		// Check if we should retry
		if err == nil && resp.StatusCode < 500 {
			// Success or client error - don't retry
			return resp, nil
		}

		// Log retry attempt
		if attempt < c.retryConfig.MaxRetries {
			if err != nil {
				// Network error
				Debug("Request failed, retrying", map[string]interface{}{
					"attempt":     attempt + 1,
					"maxAttempts": c.retryConfig.MaxRetries + 1,
					"error":       err.Error(),
					"backoff":     backoff.String(),
					"url":         req.URL.String(),
				})
			} else {
				// Server error
				Debug("Request failed with server error, retrying", map[string]interface{}{
					"attempt":     attempt + 1,
					"maxAttempts": c.retryConfig.MaxRetries + 1,
					"status":      resp.StatusCode,
					"backoff":     backoff.String(),
					"url":         req.URL.String(),
				})
				resp.Body.Close()
			}

			// Wait before retry with exponential backoff
			select {
			case <-time.After(backoff):
				// Continue to next attempt
			case <-ctx.Done():
				return nil, ctx.Err()
			}

			// Increase backoff for next attempt
			backoff = time.Duration(float64(backoff) * c.retryConfig.BackoffFactor)
			if backoff > c.retryConfig.MaxBackoff {
				backoff = c.retryConfig.MaxBackoff
			}
		}
	}

	// All retries exhausted
	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", c.retryConfig.MaxRetries+1, err)
	}
	return resp, fmt.Errorf("request failed with status %d after %d attempts", resp.StatusCode, c.retryConfig.MaxRetries+1)
}

// IsRetryableError determines if an error or status code is retryable
func IsRetryableError(err error, statusCode int) bool {
	if err != nil {
		// Network errors are generally retryable
		return true
	}

	// Retry on server errors and rate limit errors
	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}

	return false
}
