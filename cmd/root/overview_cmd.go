package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"tokenwatch/internal/config"
	"tokenwatch/pkg/models"
	"tokenwatch/pkg/providers"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Show combined consumption and costs from all configured platforms",
	Long: `Display a comprehensive overview combining data from all configured platforms.
	
This command provides a unified dashboard showing:
• Combined usage statistics across all platforms
• Platform-by-platform breakdown
• Cost comparisons and analysis
• Time period analysis

Examples:
  tokenwatch all                    # Last 7 days
  tokenwatch all --period 30d      # Last 30 days
  tokenwatch all --period 90d      # Last 90 days
  tokenwatch all -w                 # Watch mode - refresh every 30s
  tokenwatch all -w -p 30d          # Watch mode with 30-day period`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get flags
		period, _ := cmd.Flags().GetString("period")
		watch, _ := cmd.Flags().GetBool("watch")

		// Validate period
		if period != "" && period != "7d" && period != "30d" && period != "90d" {
			return fmt.Errorf("invalid period: %s. Use 7d, 30d, or 90d", period)
		}
		if period == "" {
			period = "7d" // Default to 7 days
		}

		// Get available platforms
		platforms := getAvailablePlatforms()
		if len(platforms) == 0 {
			fmt.Println("❌ No platforms configured. Run 'tokenwatch setup' first.")
			return nil
		}

		// If watch mode, run in a loop
		if watch {
			for {
				// Clear screen
				fmt.Print("\033[H\033[2J")

				// Display data with cache bypassed for fresh data
				if err := displayAllPlatformsData(platforms, period, true); err != nil {
					fmt.Printf("❌ Error: %v\n", err)
				}

				// Show refresh info
				fmt.Printf("\n🔄 Refreshing every 30 seconds... (Press Ctrl+C to stop)\n")

				// Wait 30 seconds
				time.Sleep(30 * time.Second)
			}
		} else {
			// Single run
			return displayAllPlatformsData(platforms, period, false)
		}
	},
}

func init() {
	allCmd.Flags().StringP("period", "p", "7d", "Time period (7d, 30d, 90d)")
	allCmd.Flags().BoolP("watch", "w", false, "Watch mode - refresh every 30 seconds")
	RootCmd.AddCommand(allCmd)
}

// displayAllPlatformsData fetches and displays combined platform data
func displayAllPlatformsData(platforms []string, period string, bypassCache bool) error {
	fmt.Printf("🎯 TOKENWATCH ALL PLATFORMS - Last %s\n", period)
	fmt.Printf("⏰ Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("🔗 Period: %s\n\n", getPeriodDescription(period))

	fmt.Printf("🚀 Fetching data from %d configured platform(s)...\n\n", len(platforms))

	// Collect data from all platforms in parallel
	results := collectPlatformDataParallel(platforms, period, bypassCache) // Pass false for bypassCache

	if len(results) == 0 {
		fmt.Println("ℹ️  No data found for the specified period.")
		return nil
	}

	// Calculate combined totals
	var totalTokens int64
	var totalRequests int64
	var totalCost float64

	for _, result := range results {
		for _, consumption := range result.consumptions {
			totalTokens += consumption.TotalTokens
			totalRequests += consumption.RequestCount
		}
		for _, pricing := range result.pricings {
			totalCost += pricing.Amount
		}
	}

	// Display combined dashboard
	displayCombinedDashboard(period, totalTokens, totalRequests, totalCost)

	// Display detailed platform breakdown table
	displayAllPlatformsTable(results, period)

	return nil
}

// PlatformDataResult holds the data collected from a single platform
type PlatformDataResult struct {
	platform     string
	consumptions []*models.Consumption
	pricings     []*models.Pricing
	err          error
}

// collectPlatformDataParallel fetches data from all platforms concurrently
func collectPlatformDataParallel(platforms []string, period string, bypassCache bool) []PlatformDataResult {
	var wg sync.WaitGroup
	results := make([]PlatformDataResult, 0, len(platforms))
	resultChan := make(chan PlatformDataResult, len(platforms))

	// Start goroutines for each platform
	for _, platform := range platforms {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			provider := getProvider(p)
			if provider == nil || !provider.IsAvailable() {
				resultChan <- PlatformDataResult{
					platform: p,
					err:      fmt.Errorf("provider not available"),
				}
				return
			}

			// Get detailed consumption data
			startTime, endTime := providers.GetPeriodTimeRange(period)
			consumptions, err := provider.GetConsumption(startTime, endTime, bypassCache)
			if err != nil {
				resultChan <- PlatformDataResult{
					platform: p,
					err:      fmt.Errorf("consumption data not available - check API key"),
				}
				return
			}

			// Get detailed pricing data
			pricings, err := provider.GetPricing(startTime, endTime, bypassCache)
			if err != nil {
				resultChan <- PlatformDataResult{
					platform: p,
					err:      fmt.Errorf("pricing data not available - check API key"),
				}
				return
			}

			resultChan <- PlatformDataResult{
				platform:     p,
				consumptions: consumptions,
				pricings:     pricings,
				err:          nil,
			}
		}(platform)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		if result.err == nil {
			results = append(results, result)
		} else {
			fmt.Printf("⚠️  %s: %v\n", strings.Title(result.platform), result.err)
		}
	}

	return results
}

// getAvailablePlatforms returns a list of configured platforms
func getAvailablePlatforms() []string {
	var platforms []string

	// Check all supported platforms
	supportedPlatforms := []string{"openai", "anthropic", "grok", "cursor"}

	for _, platform := range supportedPlatforms {
		if config.GetAPIKey(platform) != "" {
			platforms = append(platforms, platform)
		}
	}

	return platforms
}

// getProvider creates a provider instance for the given platform
func getProvider(platform string) providers.Provider {
	switch platform {
	case "openai":
		apiKey := config.GetAPIKey("openai")
		if apiKey == "" {
			return nil
		}
		return providers.NewOpenAIProvider(apiKey, "")
	case "anthropic":
		// TODO: Implement when Anthropic provider is ready
		return nil
	case "grok":
		// TODO: Implement when Grok provider is ready
		return nil
	case "cursor":
		// TODO: Implement when Cursor provider is ready
		return nil
	default:
		return nil
	}
}

// getPeriodDescription returns a human-readable description of the period
func getPeriodDescription(period string) string {
	switch period {
	case "7d":
		return "Last 7 days"
	case "30d":
		return "Last 30 days"
	case "90d":
		return "Last 90 days"
	default:
		return "Custom period"
	}
}

// displayCombinedDashboard shows the main dashboard with combined metrics
func displayCombinedDashboard(period string, totalTokens, totalRequests int64, totalCost float64) {
	startTime, endTime := providers.GetPeriodTimeRange(period)
	days := int(endTime.Sub(startTime).Hours() / 24)

	fmt.Println("📊 COMBINED DASHBOARD METRICS")
	fmt.Println("─" + strings.Repeat("─", 50))

	// Token metrics
	fmt.Printf("🔤 Combined Token Usage:\n")
	fmt.Printf("   Total Tokens: %s\n", color.GreenString("%d", totalTokens))
	fmt.Printf("   Total Requests: %s\n", color.BlueString("%d", totalRequests))
	if totalRequests > 0 {
		fmt.Printf("   Average per Request: %s\n", color.YellowString("%.1f", float64(totalTokens)/float64(totalRequests)))
	}
	fmt.Printf("   Daily Average: %s\n", color.CyanString("%.1f", float64(totalTokens)/float64(days)))

	// Cost metrics
	fmt.Printf("\n💰 Combined Cost Analysis:\n")
	if totalCost > 0 {
		fmt.Printf("   Total Cost: %s\n", color.GreenString("$%.4f", totalCost))
		fmt.Printf("   Daily Average: %s\n", color.YellowString("$%.4f", totalCost/float64(days)))
		if totalTokens > 0 {
			fmt.Printf("   Cost per Token: %s\n", color.CyanString("$%.6f", totalCost/float64(totalTokens)))
		}
	} else {
		fmt.Printf("   Total Cost: %s\n", color.GreenString("$0.0000 (Free Tier)"))
		fmt.Printf("   Daily Average: %s\n", color.YellowString("$0.0000"))
	}

	fmt.Printf("\n📅 Time Period: %s to %s (%d days)\n",
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		days)
	fmt.Println()
}

// displayAllPlatformsTable shows a comprehensive table of all platforms with model breakdown
func displayAllPlatformsTable(results []PlatformDataResult, period string) {
	fmt.Println("📋 ALL PLATFORMS MODEL BREAKDOWN")

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Platform", "Model", "Input Tokens", "Output Tokens", "Total Tokens", "Requests", "Cost", "$/1K Tokens")

	// Collect data from all providers
	var rows [][]string
	var allModels []ModelStats

	for _, result := range results {
		// Create a map to aggregate data by model
		modelMap := make(map[string]*ModelStats)

		// Aggregate consumption data by model
		for _, consumption := range result.consumptions {
			if _, exists := modelMap[consumption.Model]; !exists {
				modelMap[consumption.Model] = &ModelStats{Model: consumption.Model}
			}
			stats := modelMap[consumption.Model]
			stats.InputTokens += consumption.InputTokens
			stats.OutputTokens += consumption.OutputTokens
			stats.TotalTokens += consumption.TotalTokens
			stats.Requests += consumption.RequestCount
		}

		// Add pricing data
		for _, pricing := range result.pricings {
			if stats, exists := modelMap[pricing.Model]; exists {
				stats.Cost += pricing.Amount
			}
		}

		// Convert to slice and sort by total tokens (descending)
		var models []ModelStats
		for _, stats := range modelMap {
			if stats.TotalTokens > 0 || stats.Cost > 0 {
				models = append(models, *stats)
				// Also collect for overall totals
				allModels = append(allModels, *stats)
			}
		}
		sort.Slice(models, func(i, j int) bool {
			return models[i].TotalTokens > models[j].TotalTokens
		})

		// Add rows for each model
		for _, model := range models {
			costPer1K := float64(0)
			if model.TotalTokens > 0 && model.Cost > 0 {
				costPer1K = (model.Cost / float64(model.TotalTokens)) * 1000
			}

			// Color-code the platform name for better distinction
			platformName := color.CyanString(strings.Title(result.platform))

			row := []string{
				platformName,
				color.YellowString(model.Model),
				color.GreenString("%d", model.InputTokens),
				color.BlueString("%d", model.OutputTokens),
				color.WhiteString("%d", model.TotalTokens),
				color.MagentaString("%d", model.Requests),
				color.CyanString("$%.4f", model.Cost),
				color.HiBlackString("$%.4f", costPer1K),
			}
			rows = append(rows, row)
		}
	}

	// Add separator line above total
	separatorRow := []string{
		"─",
		"─",
		"─",
		"─",
		"─",
		"─",
		"─",
		"─",
	}
	rows = append(rows, separatorRow)

	// Calculate overall totals across all platforms and models
	var totalInput, totalOutput, totalTokens, totalRequests int64
	var totalCost float64
	for _, model := range allModels {
		totalInput += model.InputTokens
		totalOutput += model.OutputTokens
		totalTokens += model.TotalTokens
		totalRequests += model.Requests
		totalCost += model.Cost
	}

	// Add summary row with overall totals
	costPer1K := float64(0)
	if totalTokens > 0 && totalCost > 0 {
		costPer1K = (totalCost / float64(totalTokens)) * 1000
	}

	summaryRow := []string{
		color.HiWhiteString("TOTAL"),
		color.HiWhiteString(""),
		color.HiGreenString("%d", totalInput),
		color.HiBlueString("%d", totalOutput),
		color.HiWhiteString("%d", totalTokens),
		color.HiMagentaString("%d", totalRequests),
		color.HiYellowString("$%.4f", totalCost),
		color.HiCyanString("$%.4f", costPer1K),
	}
	rows = append(rows, summaryRow)

	table.Bulk(rows)
	table.Render()
	fmt.Println()
}
