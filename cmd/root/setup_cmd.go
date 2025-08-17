package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tokenwatch/internal/config"
	"tokenwatch/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runSetup() error {
	fmt.Println("🚀 Welcome to TokenWatch Setup!")
	fmt.Println("This will guide you through setting up API keys for supported platforms.")
	fmt.Println()

	// Define supported platforms
	platforms := []struct {
		id          string
		name        string
		description string
		implemented bool
	}{
		{"openai", "OpenAI", "ChatGPT, GPT-4, and other OpenAI models", true},
		{"anthropic", "Anthropic", "Claude and other Anthropic models", false},
		{"grok", "Grok", "xAI's Grok model", false},
		{"cursor", "Cursor", "Cursor AI models", false},
	}

	// Display platform options
	fmt.Println("Available platforms:")
	for i, platform := range platforms {
		status := "✅"
		if !platform.implemented {
			status = "🚧"
		}
		fmt.Printf("  %d. %s %s - %s\n", i+1, status, platform.name, platform.description)
	}
	fmt.Println()

	// Get platform selection
	var selectedPlatform string
	for {
		choice := utils.Prompt("Select a platform (1-4) or 'q' to quit: ")
		if strings.ToLower(choice) == "q" {
			fmt.Println("Setup cancelled.")
			return nil
		}

		switch choice {
		case "1":
			selectedPlatform = "openai"
			break
		case "2":
			selectedPlatform = "anthropic"
			break
		case "3":
			selectedPlatform = "grok"
			break
		case "4":
			selectedPlatform = "cursor"
			break
		default:
			fmt.Println("❌ Invalid choice. Please select 1-4 or 'q' to quit.")
			continue
		}
		break
	}

	platform := strings.ToLower(selectedPlatform)

	// Check if key already exists (try to load existing config first)
	var existingKey string
	if err := config.Init(); err == nil {
		existingKey = config.GetAPIKey(platform)
	}

	if existingKey != "" {
		fmt.Printf("🔄 Updating existing %s configuration...\n", strings.Title(platform))
	} else {
		fmt.Printf("🚀 Setting up %s configuration...\n", strings.Title(platform))
	}

	// Platform-specific instructions
	switch platform {
	case "openai":
		fmt.Println("\n⚠️  Important: OpenAI requires an Admin API key for organization-level access.")
		fmt.Println("   This key must have 'api.usage.read' scope to access usage and costs data.")
		fmt.Println("   Personal API keys won't work for these endpoints.")
	case "anthropic":
		fmt.Println("\nℹ️  Anthropic API key setup (placeholder implementation)")
		fmt.Println("   Full functionality coming soon!")
	case "grok":
		fmt.Println("\nℹ️  Grok API key setup (placeholder implementation)")
		fmt.Println("   Full functionality coming soon!")
	case "cursor":
		fmt.Println("\nℹ️  Cursor API key setup (placeholder implementation)")
		fmt.Println("   Full functionality coming soon!")
	}

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
	var keyPrompt string
	switch platform {
	case "openai":
		keyPrompt = fmt.Sprintf("Enter %s Admin API Key (must have api.usage.read scope): ", strings.Title(platform))
	default:
		keyPrompt = fmt.Sprintf("Enter %s API Key: ", strings.Title(platform))
	}

	apiKey := utils.PromptMasked(keyPrompt)
	if apiKey == "" {
		return fmt.Errorf("API key is required for %s setup", platform)
	}

	// Validate the API key
	fmt.Printf("🔍 Validating %s API key...\n", strings.Title(platform))
	if err := utils.ValidatePlatformKey(platform, apiKey); err != nil {
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
	v.Set(fmt.Sprintf("api_keys.%s", platform), apiKey)

	// Save config
	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	fmt.Println()
	if existingKey != "" {
		fmt.Printf("✅ %s configuration updated successfully!\n", strings.Title(platform))
	} else {
		fmt.Printf("✅ %s setup complete!\n", strings.Title(platform))
	}
	fmt.Println()

	// Platform-specific commands and status
	switch platform {
	case "openai":
		fmt.Println("🚀 You can now use these commands:")
		fmt.Println("   • tokenwatch openai              - View usage and costs")
		fmt.Println("   • tokenwatch openai --period 1d  - View last 24 hours")
		fmt.Println("   • tokenwatch openai --period 7d  - View last 7 days")
		fmt.Println("   • tokenwatch openai --period 30d - View last 30 days")
		fmt.Println("   • tokenwatch config check        - Verify your setup")
	default:
		fmt.Printf("🚧 %s is currently in development.\n", strings.Title(platform))
		fmt.Println("   Your API key has been saved and will be used when full support is implemented.")
		fmt.Println("   • tokenwatch config check        - Verify your setup")
	}

	// Ask if user wants to set up another platform
	fmt.Println()
	another := utils.Prompt("Would you like to set up another platform? (y/n): ")
	if strings.ToLower(another) == "y" || strings.ToLower(another) == "yes" {
		fmt.Println()
		// Recursively call setup for another platform
		return runSetup()
	}

	fmt.Println("\n🎉 Setup complete! You can run 'tokenwatch setup' again anytime to add more platforms.")
	return nil
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup for platform API keys",
	Long: `Interactive setup for platform API keys.
	
This command will guide you through setting up API keys for supported platforms.
Currently supports: OpenAI, Anthropic, Grok, and Cursor.

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
