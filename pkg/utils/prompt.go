package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Prompt displays a prompt and reads user input
func Prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// PromptMasked displays a prompt and reads user input without echoing (for passwords)
func PromptMasked(label string) string {
	fmt.Print(label)

	// Get terminal file descriptor
	fd := int(os.Stdin.Fd())

	// Read password without echo
	bytePassword, err := term.ReadPassword(fd)
	if err != nil {
		// Fallback to regular input if terminal doesn't support password input
		fmt.Print(" (fallback mode): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		return strings.TrimSpace(input)
	}

	fmt.Println() // Add newline after masked input
	return string(bytePassword)
}

// ConfirmPrompt asks for yes/no confirmation
func ConfirmPrompt(label string, defaultYes bool) bool {
	defaultStr := "y/N"
	if defaultYes {
		defaultStr = "Y/n"
	}

	response := Prompt(fmt.Sprintf("%s (%s): ", label, defaultStr))
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}

// MaskAPIKey masks an API key for display
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}
