package main

import (
	"fmt"
	"os"
	"path/filepath"

	"tokenwatch/internal/config"
	"tokenwatch/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runSetup() error {
	fmt.Println("🚀 Welcome to TokenWatch Setup!")
	fmt.Println("This will guide you through setting up your OpenAI API key for token usage monitoring.")
	fmt.Println()

	// Check if key already exists (try to load existing config first)
	var existingKey string
	if err := config.Init(); err == nil {
		existingKey = config.GetAPIKey("openai")
	}

	if existingKey != "" {
		fmt.Println("🔄 Updating existing OpenAI configuration...")
	} else {
		fmt.Println("🚀 Setting up OpenAI configuration...")
	}

	fmt.Println("\n⚠️  Important: OpenAI requires an Admin API key for organization-level access.")
	fmt.Println("   This key must have 'api.usage.read' scope to access usage and costs data.")
	fmt.Println("   Personal API keys won't work for these endpoints.")
	fmt.Println()

	// Initialize Viper and read existing config (if any)
	v := viper.New()

	// Set config path
	configDir := filepath.Join(os.Getenv("HOME"), ".tokenwatch")
	configPath := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set config file and try to read existing config
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		// If no existing config, set defaults
		v.SetDefault("settings.cache_duration", 300)
		v.SetDefault("settings.request_timeout", 10)
		v.SetDefault("settings.retry_attempts", 3)
		v.SetDefault("settings.debug", false)
		v.SetDefault("data_dir", configDir)
		v.SetDefault("display.date_format", "2006-01-02 15:04:05")
		v.SetDefault("display.colors", true)
		v.SetDefault("display.table_style", "fancy")
		v.SetDefault("display.show_progress", true)
		v.SetDefault("output.default_format", "table")
	}

	// Prompt for API key (masked input)
	keyPrompt := "Enter OpenAI Admin API Key (must have api.usage.read scope): "
	apiKey := utils.PromptMasked(keyPrompt)
	if apiKey == "" {
		return fmt.Errorf("API key is required for OpenAI setup")
	}

	// Validate the API key
	fmt.Printf("🔍 Validating OpenAI API key...\n")
	if err := utils.ValidatePlatformKey("openai", apiKey); err != nil {
		fmt.Printf("❌ API key validation failed: %v\n", err)

		// Ask if user wants to continue anyway
		if !utils.ConfirmPrompt("Do you want to save this key anyway?", false) {
			return fmt.Errorf("setup cancelled due to invalid API key")
		}
		fmt.Println("⚠️  Saving unvalidated API key. You may need to update it later.")
	} else {
		fmt.Println("✅ API key validated successfully!")
	}

	// Set the API key
	v.Set("api_keys.openai", apiKey)

	// Save config
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	fmt.Println()
	if existingKey != "" {
		fmt.Println("✅ OpenAI configuration updated successfully!")
	} else {
		fmt.Println("✅ OpenAI setup complete!")
	}
	fmt.Println()

	fmt.Println("🚀 You can now use these commands:")
	fmt.Println("   • tokenwatch usage              - View usage and costs")
	fmt.Println("   • tokenwatch usage --period 1d  - View last 24 hours")
	fmt.Println("   • tokenwatch usage --period 7d  - View last 7 days")
	fmt.Println("   • tokenwatch usage --period 30d - View last 30 days")
	fmt.Println("   • tokenwatch usage -w -p 1d     - Watch mode for real-time monitoring")
	fmt.Println("   • tokenwatch config check         - Verify your setup")

	fmt.Println("\n🎉 Setup complete! You can run 'tokenwatch setup' again anytime to update your API key.")
	return nil
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup for OpenAI API key",
	Long: `Interactive setup for OpenAI API key.
	
This command will guide you through setting up your OpenAI Admin API key for token usage monitoring.
Currently supports: OpenAI only.

Examples:
  tokenwatch setup    # Start interactive setup process

Note: OpenAI requires an Admin API key for organization-level access to usage and costs data.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetup()
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
