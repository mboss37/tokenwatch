package utils

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeConfig indicates a configuration error
	ErrorTypeConfig ErrorType = "CONFIG"
	// ErrorTypeAPI indicates an API error
	ErrorTypeAPI ErrorType = "API"
	// ErrorTypeAuth indicates an authentication error
	ErrorTypeAuth ErrorType = "AUTH"
	// ErrorTypeNetwork indicates a network error
	ErrorTypeNetwork ErrorType = "NETWORK"
	// ErrorTypeRateLimit indicates a rate limit error
	ErrorTypeRateLimit ErrorType = "RATE_LIMIT"
	// ErrorTypeValidation indicates a validation error
	ErrorTypeValidation ErrorType = "VALIDATION"
	// ErrorTypeInternal indicates an internal error
	ErrorTypeInternal ErrorType = "INTERNAL"
)

// StructuredError provides detailed error information with actionable suggestions
type StructuredError struct {
	Type        ErrorType
	Message     string
	Cause       error
	Suggestions []string
	Context     map[string]interface{}
}

// Error implements the error interface
func (e *StructuredError) Error() string {
	var parts []string

	parts = append(parts, fmt.Sprintf("[%s] %s", e.Type, e.Message))

	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("Cause: %v", e.Cause))
	}

	if len(e.Suggestions) > 0 {
		parts = append(parts, "\nSuggestions:")
		for i, suggestion := range e.Suggestions {
			parts = append(parts, fmt.Sprintf("  %d. %s", i+1, suggestion))
		}
	}

	return strings.Join(parts, "\n")
}

// Unwrap returns the underlying error
func (e *StructuredError) Unwrap() error {
	return e.Cause
}

// NewConfigError creates a configuration error
func NewConfigError(message string, cause error) *StructuredError {
	return &StructuredError{
		Type:    ErrorTypeConfig,
		Message: message,
		Cause:   cause,
		Suggestions: []string{
			"Run 'tokenwatch setup' to configure your API keys",
			"Check if your config file exists at ~/.tokenwatch/config.yaml",
			"Verify file permissions on the config directory",
		},
	}
}

// NewAPIError creates an API error with helpful suggestions
func NewAPIError(message string, statusCode int, cause error) *StructuredError {
	suggestions := []string{}

	switch statusCode {
	case 401:
		suggestions = append(suggestions,
			"Verify your API key is correct",
			"Run 'tokenwatch setup' to update your API key",
			"Check if your API key has been revoked or expired")
	case 403:
		suggestions = append(suggestions,
			"Ensure your API key has the required permissions",
			"For OpenAI, you need an Admin key with 'api.usage.read' scope",
			"Contact your organization administrator for proper access")
	case 429:
		suggestions = append(suggestions,
			"You've hit the rate limit - wait a moment and try again",
			"Consider upgrading your API plan for higher limits",
			"The request will be automatically retried")
	case 500, 502, 503, 504:
		suggestions = append(suggestions,
			"The API service is experiencing issues",
			"Wait a few moments and try again",
			"Check the platform's status page for outages")
	default:
		suggestions = append(suggestions,
			"Check your internet connection",
			"Verify the API service is available",
			"Try again in a few moments")
	}

	return &StructuredError{
		Type:        ErrorTypeAPI,
		Message:     message,
		Cause:       cause,
		Suggestions: suggestions,
		Context: map[string]interface{}{
			"status_code": statusCode,
		},
	}
}

// NewAuthError creates an authentication error
func NewAuthError(message string, platform string) *StructuredError {
	return &StructuredError{
		Type:    ErrorTypeAuth,
		Message: message,
		Suggestions: []string{
			fmt.Sprintf("Run 'tokenwatch setup' to configure your %s API key", platform),
			"Verify your API key is correct and hasn't been revoked",
			"Check if you're using the right type of API key (e.g., Admin key for OpenAI)",
		},
		Context: map[string]interface{}{
			"platform": platform,
		},
	}
}

// NewNetworkError creates a network error
func NewNetworkError(message string, cause error) *StructuredError {
	return &StructuredError{
		Type:    ErrorTypeNetwork,
		Message: message,
		Cause:   cause,
		Suggestions: []string{
			"Check your internet connection",
			"Verify you can reach the API endpoint",
			"Check if you're behind a firewall or proxy",
			"Try again in a few moments",
		},
	}
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(platform string, retryAfter string) *StructuredError {
	suggestions := []string{
		"Wait a moment for the rate limit to reset",
		"The request will be automatically retried",
	}

	if retryAfter != "" {
		suggestions = append(suggestions, fmt.Sprintf("Retry after: %s", retryAfter))
	}

	return &StructuredError{
		Type:        ErrorTypeRateLimit,
		Message:     fmt.Sprintf("Rate limit exceeded for %s", platform),
		Suggestions: suggestions,
		Context: map[string]interface{}{
			"platform":    platform,
			"retry_after": retryAfter,
		},
	}
}

// NewValidationError creates a validation error
func NewValidationError(field, message string) *StructuredError {
	return &StructuredError{
		Type:    ErrorTypeValidation,
		Message: fmt.Sprintf("Invalid %s: %s", field, message),
		Context: map[string]interface{}{
			"field": field,
		},
	}
}

// FormatError formats any error with additional context and suggestions
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	// Check if it's already a structured error
	if se, ok := err.(*StructuredError); ok {
		return se.Error()
	}

	// Try to determine error type from error message
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "config"):
		return NewConfigError(errStr, err).Error()
	case strings.Contains(errStr, "401") || strings.Contains(errStr, "unauthorized"):
		return NewAuthError(errStr, "unknown").Error()
	case strings.Contains(errStr, "network") || strings.Contains(errStr, "connection"):
		return NewNetworkError(errStr, err).Error()
	default:
		return err.Error()
	}
}
