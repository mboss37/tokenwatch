package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"tokenwatch/internal/config"
	"tokenwatch/pkg/providers"
	"tokenwatch/pkg/utils"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getProvider creates and returns a provider for the specified platform
func getProvider(platform string) providers.Provider {
	switch platform {
	case "openai":
		apiKey := config.GetAPIKey("openai")
		orgID := config.GetString("openai.organization_id")
		if apiKey == "" {
			return nil
		}
		return providers.NewOpenAIProvider(apiKey, orgID)
	default:
		return nil
	}
}

// ModelStats holds aggregated stats for a model
type ModelStats struct {
	Model        string
	InputTokens  int64
	OutputTokens int64
	TotalTokens  int64
	Requests     int64
	Cost         float64
}

// TotalStats holds the totals across all models
type TotalStats struct {
	TotalInput    int64
	TotalOutput   int64
	TotalTokens   int64
	TotalRequests int64
	TotalCost     float64
}

// calculateTotals computes the totals across all models
func calculateTotals(models []ModelStats) TotalStats {
	var totals TotalStats

	for _, m := range models {
		totals.TotalInput += m.InputTokens
		totals.TotalOutput += m.OutputTokens
		totals.TotalTokens += m.TotalTokens
		totals.TotalRequests += m.Requests
		totals.TotalCost += m.Cost
	}

	return totals
}

var (
	period string
	watch  bool
	debug  bool
)

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show OpenAI token consumption and costs",
	Long: `Display comprehensive OpenAI usage including tokens consumed and costs incurred.

This command shows:
• Token consumption by model (input/output)
• Associated costs
• Time period analysis
• Request statistics

Examples:
  tokenwatch usage                # Last 7 days
        tokenwatch usage --period 1d    # Last 24 hours
        tokenwatch usage --period 30d   # Last 30 days
        tokenwatch usage --period 90d   # Last 90 days
        tokenwatch usage --period 1y    # Last 1 year
        tokenwatch usage --period all   # Last 5 years (maximum)
        tokenwatch usage -w -p 1d       # Watch mode - refresh every 30s
        tokenwatch usage -w -p 7d       # Watch mode with 7-day period
        tokenwatch usage -w -p 90d      # Watch mode with 90-day period`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		if err := config.Init(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Check if OpenAI is configured
		apiKey := config.GetAPIKey("openai")
		if apiKey == "" {
			return utils.NewAuthError("OpenAI not configured", "openai")
		}

		// Get flags
		period, _ = cmd.Flags().GetString("period")
		watch, _ = cmd.Flags().GetBool("watch")
		debug, _ = cmd.Flags().GetBool("debug")

		// Validate period
		if period != "" && period != "1d" && period != "7d" && period != "30d" {
			return fmt.Errorf("invalid period: %s. Valid periods are: 1d, 7d, 30d", period)
		}
		if period == "" {
			period = "7d"
		}

		// Validate watch mode - only allow for 1d period
		//if watch && period != "1d" {
		// return fmt.Errorf("watch mode (-w) is only available for 1-day period (--period 1d). For longer periods, use regular mode")
		// }

		// Warn about 30d period limitations
		if period == "30d" {
			fmt.Println("⚠️  Note: 30-day period may take longer to load and may have limited data availability due to OpenAI API limitations.")
			fmt.Println("   Consider using --period 7d for more reliable results.")
			fmt.Println()
		}

		// Get provider
		provider := getProvider("openai")
		if provider == nil {
			return fmt.Errorf("OpenAI provider not available")
		}

		// Cast provider to OpenAIProvider for displayOpenAIData
		openaiProvider, ok := provider.(*providers.OpenAIProvider)
		if !ok {
			return fmt.Errorf("failed to get OpenAI provider")
		}

		// If watch mode, run in a loop
		if watch {
			for {
				// Clear screen
				fmt.Print("\033[H\033[2J")

				// Display data with cache bypassed for fresh data
				if err := displayOpenAIData(openaiProvider, period, true, debug); err != nil {
					fmt.Printf("❌ Error: %v\n", err)
				}

				// Show refresh info
				fmt.Printf("\n🔄 Refreshing every 30 seconds... (Press Ctrl+C to stop)\n")

				// Wait 30 seconds
				time.Sleep(30 * time.Second)
			}
		} else {
			// Single run
			return displayOpenAIData(openaiProvider, period, false, debug)
		}
	},
}

func init() {
	usageCmd.Flags().StringP("period", "p", "7d", "Time period: 1d (recent activity), 7d (historical data), 30d, 90d, 1y, all")
	usageCmd.Flags().BoolP("watch", "w", false, "Watch mode - refresh every 30 seconds")
	usageCmd.Flags().BoolP("debug", "d", false, "Enable debug logging for API calls")
	RootCmd.AddCommand(usageCmd)
}

// displayOpenAIData fetches and displays OpenAI usage data
func displayOpenAIData(provider *providers.OpenAIProvider, period string, bypassCache bool, debug bool) error {
	// Display header
	fmt.Printf("🤖 OPENAI USAGE - Last %s\n", period)
	fmt.Printf("⏰ Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Get time range
	startTime, endTime := providers.GetPeriodTimeRange(period)

	// Fetch consumption data
	consumptions, err := provider.GetConsumption(startTime, endTime, bypassCache, debug)
	if err != nil {
		return fmt.Errorf("failed to get consumption data: %w", err)
	}

	// Fetch pricing data
	pricings, err := provider.GetPricing(startTime, endTime, bypassCache, debug)
	if err != nil {
		// Don't fail if pricing data is unavailable - just log a warning
		fmt.Printf("⚠️  Warning: Could not fetch pricing data: %v\n", err)
		fmt.Println("   This is normal for longer time periods or when costs are not yet available.")
		fmt.Println()
	}

	// Create pricing map for easy lookup
	pricingMap := make(map[string]float64)
	for _, p := range pricings {
		pricingMap[p.Model] = pricingMap[p.Model] + p.Amount
	}

	// Aggregate data by model
	modelMap := make(map[string]*ModelStats)

	for _, c := range consumptions {
		if _, exists := modelMap[c.Model]; !exists {
			modelMap[c.Model] = &ModelStats{Model: c.Model}
		}

		stats := modelMap[c.Model]
		stats.InputTokens += c.InputTokens
		stats.OutputTokens += c.OutputTokens
		stats.TotalTokens += c.TotalTokens
		stats.Requests += c.RequestCount
	}

	// Add pricing data
	for model, cost := range pricingMap {
		if stats, exists := modelMap[model]; exists {
			stats.Cost = cost
		} else {
			// Create entry for models with costs but no usage (shouldn't happen normally)
			modelMap[model] = &ModelStats{
				Model: model,
				Cost:  cost,
			}
		}
	}

	// Convert to slice and sort by total tokens (descending)
	var models []ModelStats
	for _, stats := range modelMap {
		// Include models that have either tokens or costs
		if stats.TotalTokens > 0 || stats.Cost > 0 {
			models = append(models, *stats)
		}
	}
	sort.Slice(models, func(i, j int) bool {
		// Sort by total tokens first, then by cost
		if models[i].TotalTokens != models[j].TotalTokens {
			return models[i].TotalTokens > models[j].TotalTokens
		}
		return models[i].Cost > models[j].Cost
	})

	// Check if we have any data to display
	if len(models) == 0 {
		fmt.Println("ℹ️  No consumption or cost data found for the specified period.")
		fmt.Println("   This could mean:")
		fmt.Println("   • No API calls were made during this period")
		fmt.Println("   • The data is not yet available from OpenAI")
		fmt.Println("   • You're checking a period before your account was created")
		fmt.Println()
		fmt.Println("💡 Try using a shorter period like '--period 7d' to see recent data.")
		fmt.Println("   OpenAI only returns data for periods with actual activity.")
		return nil
	}

	// Calculate totals and display
	totals := calculateTotals(models)

	// Display summary
	displayOpenAISummary(period, totals.TotalTokens, totals.TotalRequests, totals.TotalCost, startTime, endTime)

	// Display smart recommendations
	displaySmartRecommendations(period)

	// Display table
	displayOpenAITable(models, totals)

	return nil
}

// displayOpenAISummary shows overall statistics
func displayOpenAISummary(period string, totalTokens, totalRequests int64, totalCost float64, startTime, endTime time.Time) {
	days := int(endTime.Sub(startTime).Hours() / 24)

	fmt.Println("📊 SUMMARY")
	fmt.Println("─" + color.HiBlackString("─────────────────────────────────────────────────"))

	fmt.Printf("📅 Period: %s to %s (%d days)\n",
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
		days)

	fmt.Printf("📈 Daily Averages: %s tokens, %s requests\n",
		color.CyanString("%.1f", float64(totalTokens)/float64(days)),
		color.CyanString("%.1f", float64(totalRequests)/float64(days)))

	if totalCost > 0 {
		fmt.Printf("💰 Daily Cost Average: %s\n",
			color.YellowString("$%.4f", totalCost/float64(days)))
	} else {
		if period == "30d" {
			fmt.Printf("💰 Cost Data: %s\n", color.YellowString("Not available for this period"))
		} else {
			fmt.Printf("💰 Cost Data: %s\n", color.GreenString("Free Tier"))
		}
	}

	fmt.Println()
}

// displaySmartRecommendations provides smart recommendations based on the selected time period
func displaySmartRecommendations(period string) {
	fmt.Println("💡 SMART RECOMMENDATIONS")
	fmt.Println("─" + color.HiBlackString("─────────────────────────────────────────────────"))

	switch period {
	case "1d":
		fmt.Println("📊 1-day period is perfect for:")
		fmt.Println("   • Recent activity monitoring")
		fmt.Println("   • Real-time usage tracking")
		fmt.Println("   • Immediate cost calculations")
		fmt.Println("   • Debugging current API calls")
		fmt.Println()
		fmt.Println("🔄 For historical analysis, try: --period 7d")

	case "7d":
		fmt.Println("📊 7-day period is ideal for:")
		fmt.Println("   • Weekly usage patterns")
		fmt.Println("   • Historical cost analysis")
		fmt.Println("   • Model performance comparison")
		fmt.Println("   • Budget planning")
		fmt.Println()
		fmt.Println("🔄 For recent activity, try: --period 1d")

	case "30d":
		fmt.Println("📊 30-day period may have limited data:")
		fmt.Println("   • OpenAI API data availability varies")
		fmt.Println("   • Some periods may return empty results")
		fmt.Println("   • Consider using 7d for reliable data")
		fmt.Println()
		fmt.Println("🔄 For best results, try: --period 7d")

	default:
		fmt.Println("📊 For optimal results:")
		fmt.Println("   • Use --period 1d for recent activity")
		fmt.Println("   • Use --period 7d for historical data")
		fmt.Println("   • 30d period may have data limitations")
	}

	fmt.Println()
}

// displayOpenAITable shows detailed model breakdown
func displayOpenAITable(models []ModelStats, totals TotalStats) {
	fmt.Println("📋 MODEL BREAKDOWN")

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Model", "Input Tokens", "Output Tokens", "Total Tokens", "Requests", "Cost", "$/1K Tokens")

	// Add rows
	var rows [][]string

	for _, m := range models {
		costPer1K := float64(0)
		if m.TotalTokens > 0 && m.Cost > 0 {
			costPer1K = (m.Cost / float64(m.TotalTokens)) * 1000
		}

		row := []string{
			color.YellowString(m.Model),
			color.GreenString("%d", m.InputTokens),
			color.BlueString("%d", m.OutputTokens),
			color.WhiteString("%d", m.TotalTokens),
			color.MagentaString("%d", m.Requests),
			color.CyanString("$%.4f", m.Cost),
			color.HiBlackString("$%.4f", costPer1K),
		}
		rows = append(rows, row)
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
	}
	rows = append(rows, separatorRow)

	// Add summary row using pre-calculated totals
	costPer1K := float64(0)
	if totals.TotalTokens > 0 && totals.TotalCost > 0 {
		costPer1K = (totals.TotalCost / float64(totals.TotalTokens)) * 1000
	}

	summaryRow := []string{
		color.HiWhiteString("TOTAL"),
		color.HiGreenString("%d", totals.TotalInput),
		color.HiBlueString("%d", totals.TotalOutput),
		color.HiWhiteString("%d", totals.TotalTokens),
		color.HiMagentaString("%d", totals.TotalRequests),
		color.HiYellowString("$%.4f", totals.TotalCost),
		color.HiCyanString("$%.4f", costPer1K),
	}
	rows = append(rows, summaryRow)

	table.Bulk(rows)
	table.Render()
}
