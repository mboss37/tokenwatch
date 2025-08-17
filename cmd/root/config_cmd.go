package main

import (
	"fmt"
	"os"
	"strings"

	"tokenwatch/internal/config"
	"tokenwatch/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage TokenWatch configuration",
	Long: `Configuration management for TokenWatch.

This command allows you to:
• Check configuration status and API keys
• Reload configuration from file
• Reset settings to defaults`,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Verify configuration and API keys",
	Long:  `Check the current configuration status, validate settings, and verify API keys.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config - this will handle missing files gracefully
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println("🔍 CONFIGURATION STATUS")
		fmt.Println("─" + strings.Repeat("─", 50))

		// Check config file
		configFile := config.GetConfigFile()
		if configFile != "" && configFile != "none" {
			fmt.Printf("✅ Config file: %s\n", color.GreenString(configFile))
		} else {
			fmt.Printf("ℹ️  Config file: %s\n", color.CyanString("Using default configuration"))
		}

		// Check API keys
		fmt.Println("\n🔑 API KEYS:")
		platforms := []string{"openai"}
		hasKeys := false
		for _, platform := range platforms {
			key := config.GetAPIKey(platform)
			if key != "" {
				fmt.Printf("   ✅ %s: %s\n", strings.Title(platform), color.GreenString("Configured"))
				hasKeys = true
			} else {
				fmt.Printf("   ❌ %s: %s\n", strings.Title(platform), color.RedString("Not configured"))
			}
		}

		if !hasKeys {
			fmt.Printf("\n💡 Run 'tokenwatch setup' to configure your OpenAI API key\n")
		}

		fmt.Println("\n✅ Configuration check complete!")
		return nil
	},
}

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload configuration from file",
	Long:  `Reload configuration from the config file without restarting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to reload config: %w", err)
		}
		fmt.Println("✅ Configuration reloaded successfully")
		return nil
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset [key]",
	Short: "Reset configuration to defaults",
	Long: `Reset configuration values to defaults.

Examples:
  tokenwatch config reset                    # Reset all settings
  tokenwatch config reset settings.debug    # Reset specific setting`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(args) == 0 {
			// Reset all settings
			fmt.Println("⚠️  This will reset ALL configuration to defaults!")

			if !utils.ConfirmPrompt("Are you sure?", false) {
				fmt.Println("❌ Reset cancelled")
				return nil
			}

			// Remove config file to reset to defaults
			configFile := config.GetConfigFile()
			if configFile != "" && configFile != "none" {
				if err := os.Remove(configFile); err != nil {
					return fmt.Errorf("failed to remove config file: %w", err)
				}
				fmt.Println("✅ Configuration reset to defaults")
				fmt.Println("💡 Run 'tokenwatch config check' to verify the reset")
			} else {
				fmt.Println("ℹ️  Already using default configuration")
			}
		} else {
			// Reset specific key
			key := args[0]
			config.Set(key, nil) // Set to nil to remove

			if err := config.WriteConfig(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("✅ Reset %s to default\n", color.CyanString(key))
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(checkCmd)
	configCmd.AddCommand(reloadCmd)
	configCmd.AddCommand(resetCmd)
	RootCmd.AddCommand(configCmd)
}
