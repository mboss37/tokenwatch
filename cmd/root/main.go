package main

import (
	"os"

	"tokenwatch/internal/config"
	"tokenwatch/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	Version   = "v0.1.1"
	BuildTime = "unknown"
)

var RootCmd = &cobra.Command{
	Use:     "tokenwatch",
	Version: Version,
	Short:   "TokenWatch: Track OpenAI token consumption and pricing",
	Long: `TokenWatch is a CLI tool for monitoring OpenAI token usage and costs.
Currently supports OpenAI organization-level usage and costs monitoring.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Initialize logger early
	logLevel := utils.InfoLevel

	// Try to load config to get debug setting
	if err := config.Init(); err == nil {
		if config.GetBool("settings.debug") {
			logLevel = utils.DebugLevel
		}
		// Check environment variable
		if envLevel := os.Getenv("TOKENWATCH_LOG_LEVEL"); envLevel != "" {
			logLevel = utils.ParseLogLevel(envLevel)
		}
	}

	// Initialize logger with color support for terminal
	utils.InitLogger(logLevel, true)
}

func main() {
	if err := Execute(); err != nil {
		utils.Error("Command execution failed", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
}
