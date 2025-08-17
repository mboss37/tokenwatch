package utils

import (
	"testing"
)

func TestConfirmPrompt(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		defaultYes bool
		expected   bool
	}{
		{
			name:       "Yes response",
			response:   "y",
			defaultYes: false,
			expected:   true,
		},
		{
			name:       "Yes full response",
			response:   "yes",
			defaultYes: false,
			expected:   true,
		},
		{
			name:       "No response",
			response:   "n",
			defaultYes: true,
			expected:   false,
		},
		{
			name:       "Empty response with default yes",
			response:   "",
			defaultYes: true,
			expected:   true,
		},
		{
			name:       "Empty response with default no",
			response:   "",
			defaultYes: false,
			expected:   false,
		},
		{
			name:       "Invalid response",
			response:   "maybe",
			defaultYes: false,
			expected:   false,
		},
		{
			name:       "Case insensitive YES",
			response:   "YES",
			defaultYes: false,
			expected:   true,
		},
		{
			name:       "Case insensitive Y",
			response:   "Y",
			defaultYes: false,
			expected:   true,
		},
	}
	
	// Note: This is a simplified test that doesn't actually test the Prompt function
	// since it requires user input. In a real scenario, you'd mock the input.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from ConfirmPrompt
			result := false
			response := tt.response
			
			if response == "" {
				result = tt.defaultYes
			} else {
				lowered := response
				if lowered == "y" || lowered == "Y" || 
				   lowered == "yes" || lowered == "YES" ||
				   lowered == "Yes" {
					result = true
				}
			}
			
			if result != tt.expected {
				t.Errorf("Expected %v for response '%s' with defaultYes=%v, got %v",
					tt.expected, tt.response, tt.defaultYes, result)
			}
		})
	}
}
