package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of TokenWatch",
	Long:  `Print the version number and build information of TokenWatch CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TokenWatch %s\n", Version)
		fmt.Printf("Built: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", "go1.21+")
		fmt.Printf("Platform: OpenAI only\n")
		fmt.Printf("Status: First functional release\n")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
