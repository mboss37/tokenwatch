package utils

import (
	"testing"
)

// TestValidatePlatformKeyFormat tests basic format validation
// Note: This doesn't test actual API validation which requires real API keys
func TestValidatePlatformKeyFormat(t *testing.T) {
	tests := []struct {
		name        string
		platform    string
		apiKey      string
		expectError bool
	}{
		{
			name:        "Empty API key",
			platform:    "openai",
			apiKey:      "",
			expectError: true,
		},
		{
			name:        "Unsupported platform",
			platform:    "unknown",
			apiKey:      "some-key",
			expectError: true,
		},
		{
			name:        "Generic platform with short key",
			platform:    "grok",
			apiKey:      "123",
			expectError: true,
		},
		{
			name:        "Generic platform with valid length key",
			platform:    "cursor",
			apiKey:      "1234567890abc",
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePlatformKey(tt.platform, tt.apiKey)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestValidateOpenAIKeyFormat tests OpenAI key format validation
func TestValidateOpenAIKeyFormat(t *testing.T) {
	// Test format validation without making actual API calls
	// For real API validation, use integration tests with valid keys
	
	tests := []struct {
		name        string
		apiKey      string
		shouldFail  bool
		reason      string
	}{
		{
			name:       "Empty key",
			apiKey:     "",
			shouldFail: true,
			reason:     "empty key should fail",
		},
		{
			name:       "Wrong prefix",
			apiKey:     "pk-1234567890",
			shouldFail: true,
			reason:     "non 'sk-' prefix should fail",
		},
		{
			name:       "Valid format",
			apiKey:     "sk-proj-1234567890abcdef",
			shouldFail: false,
			reason:     "valid format (actual validation would require real API)",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: ValidateOpenAIKey now makes real API calls
			// so we can only test format validation for invalid keys
			if tt.shouldFail {
				err := ValidateOpenAIKey(tt.apiKey)
				if err == nil {
					t.Errorf("Expected error for %s", tt.reason)
				}
			}
			// Skip testing valid format as it would make real API calls
		})
	}
}
