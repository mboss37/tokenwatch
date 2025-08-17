package config

import (
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Initialize config
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}
	
	// Test getting non-existent key
	key := GetAPIKey("nonexistent")
	if key != "" {
		t.Errorf("Expected empty string for non-existent key, got '%s'", key)
	}
	
	// Test that GetAPIKey returns something for valid platforms
	// (could be from config file or env var)
	platforms := []string{"openai", "anthropic", "grok", "cursor"}
	for _, platform := range platforms {
		// Just verify the function doesn't panic
		_ = GetAPIKey(platform)
	}
}

func TestGetBool(t *testing.T) {
	// Initialize config
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}
	
	// Test default value
	debug := GetBool("settings.debug")
	if debug != false {
		t.Error("Expected debug to be false by default")
	}
}

func TestGetString(t *testing.T) {
	// Initialize config
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}
	
	// Test non-existent string
	value := GetString("nonexistent.key")
	if value != "" {
		t.Errorf("Expected empty string for non-existent key, got '%s'", value)
	}
}

func TestGetInt(t *testing.T) {
	// Initialize config
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}
	
	// Test non-existent int (should return 0)
	value := GetInt("nonexistent.number")
	if value != 0 {
		t.Errorf("Expected 0 for non-existent int key, got %d", value)
	}
}


