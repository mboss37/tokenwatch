package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var Config *viper.Viper

func Init() error {
	Config = viper.New()
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")

	// Default config path
	configDir := filepath.Join(os.Getenv("HOME"), ".tokenwatch")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create dir if not exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	// Bind env vars
	Config.AutomaticEnv()
	Config.SetEnvPrefix("TOKENWATCH")

	// Set defaults
	Config.SetDefault("settings.cache_duration", 300)
	Config.SetDefault("settings.request_timeout", 10)
	Config.SetDefault("settings.retry_attempts", 3)
	Config.SetDefault("settings.debug", false)
	Config.SetDefault("data_dir", configDir)
	Config.SetDefault("display.date_format", "2006-01-02 15:04:05")
	Config.SetDefault("display.colors", true)
	Config.SetDefault("display.table_style", "fancy")
	Config.SetDefault("display.show_progress", true)
	Config.SetDefault("output.default_format", "table")

	// Read config file if exists
	Config.SetConfigFile(configPath)
	err := Config.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// File not found; use defaults and env - no error
		} else if strings.Contains(err.Error(), "no such file or directory") {
			// File was deleted or doesn't exist - use defaults and env - no error
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

func GetAPIKey(platform string) string {
	key := Config.GetString("api_keys." + platform)
	if key == "" {
		// Try environment variables for each platform
		switch platform {
		case "openai":
			key = os.Getenv("OPENAI_API_KEY")
		case "anthropic":
			key = os.Getenv("ANTHROPIC_API_KEY")
		case "grok":
			key = os.Getenv("GROK_API_KEY")
		case "cursor":
			key = os.Getenv("CURSOR_API_KEY")
		}
	}
	return key
}

// GetCacheDuration retrieves the cache duration in seconds
func GetCacheDuration() int {
	// Default to 5 minutes if not set
	duration := Config.GetInt("settings.cache_duration")
	if duration <= 0 {
		return 300 // 5 minutes default
	}
	return duration
}

// GetString retrieves a string configuration value
func GetString(key string) string {
	return Config.GetString(key)
}

// GetInt retrieves an integer configuration value
func GetInt(key string) int {
	return Config.GetInt(key)
}

// GetBool retrieves a boolean configuration value
func GetBool(key string) bool {
	return Config.GetBool(key)
}

// GetFloat64 retrieves a float64 configuration value
func GetFloat64(key string) float64 {
	return Config.GetFloat64(key)
}

// Set sets a configuration value
func Set(key string, value interface{}) {
	Config.Set(key, value)
}

// WriteConfig writes the current configuration to file
func WriteConfig() error {
	return Config.WriteConfig()
}

// GetConfigFile returns the current config file path
func GetConfigFile() string {
	file := Config.ConfigFileUsed()
	if file == "" {
		return "none"
	}
	
	// Check if the file actually exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return "none"
	}
	
	return file
}

// IsSet checks if a key is set in configuration
func IsSet(key string) bool {
	return Config.IsSet(key)
}

// AllSettings returns all configuration settings
func AllSettings() map[string]interface{} {
	return Config.AllSettings()
}
