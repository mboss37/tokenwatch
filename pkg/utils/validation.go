package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ValidateOpenAIKey validates an OpenAI API key by making a test request
func ValidateOpenAIKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Basic format validation
	if !strings.HasPrefix(apiKey, "sk-") {
		return fmt.Errorf("invalid API key format: OpenAI keys should start with 'sk-'")
	}

	// Test the key with a minimal API call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.openai.com/v1/models", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API connection failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		Info("API key validated successfully")
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("invalid API key: authentication failed")
	case http.StatusForbidden:
		// This might happen if the key doesn't have the right permissions
		// For OpenAI organization keys, we need to check usage endpoint
		return validateOpenAIUsageAccess(apiKey)
	default:
		return fmt.Errorf("unexpected response: %d %s", resp.StatusCode, resp.Status)
	}
}

// validateOpenAIUsageAccess specifically checks if the key has usage API access
func validateOpenAIUsageAccess(apiKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to access the usage endpoint with a minimal date range
	endTime := time.Now()
	startTime := endTime.Add(-24 * time.Hour)

	url := fmt.Sprintf("https://api.openai.com/v1/organization/usage/completions?start_time=%d&end_time=%d",
		startTime.Unix(), endTime.Unix())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API connection failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		Info("API key has usage access")
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("invalid API key: authentication failed")
	case http.StatusForbidden:
		return fmt.Errorf("API key lacks required permissions: needs 'api.usage.read' scope for organization-level access")
	default:
		return fmt.Errorf("unexpected response from usage API: %d %s", resp.StatusCode, resp.Status)
	}
}

// ValidatePlatformKey validates an API key for the specified platform
func ValidatePlatformKey(platform, apiKey string) error {
	switch platform {
	case "openai":
		return ValidateOpenAIKey(apiKey)
	case "anthropic":
		// TODO: Implement Anthropic validation when provider is ready
		return validateGenericKey(platform, apiKey)
	case "grok":
		// TODO: Implement Grok validation when provider is ready
		return validateGenericKey(platform, apiKey)
	case "cursor":
		// TODO: Implement Cursor validation when provider is ready
		return validateGenericKey(platform, apiKey)
	default:
		return fmt.Errorf("unsupported platform: %s", platform)
	}
}

// validateGenericKey performs basic validation for platforms not yet implemented
func validateGenericKey(platform, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("API key seems too short")
	}

	// For now, just accept the key for platforms not yet implemented
	Warn(fmt.Sprintf("%s platform validation not yet implemented, accepting key", platform))
	return nil
}
