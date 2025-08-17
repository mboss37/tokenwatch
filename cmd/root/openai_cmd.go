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

var openaiCmd = &cobra.Command{
	Use:   "openai",
	Short: "Show OpenAI token consumption and costs",
	Long: `Display comprehensive OpenAI usage including tokens consumed and costs incurred.

This command shows:
• Token consumption by model (input/output)
• Associated costs
• Time period analysis
• Request statistics

Examples:
  tokenwatch openai                # Last 7 days
  tokenwatch openai --period 30d   # Last 30 days
  tokenwatch openai --period 90d   # Last 90 days
  tokenwatch openai -w             # Watch mode - refresh every 30s
  tokenwatch openai -w -p 30d      # Watch mode with 30-day period`,
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
		period, _ := cmd.Flags().GetString("period")
		watch, _ := cmd.Flags().GetBool("watch")

		if period != "" && period != "7d" && period != "30d" && period != "90d" {
			return fmt.Errorf("invalid period: %s. Use 7d, 30d, or 90d", period)
		}
		if period == "" {
			period = "7d" // Default
		}

		// Create provider
		provider := providers.NewOpenAIProvider(apiKey, "")

		// If watch mode, run in a loop
		if watch {
			for {
				// Clear screen
				fmt.Print("\033[H\033[2J")
				
				// Display data with cache bypassed for fresh data
				if err := displayOpenAIData(provider, period, true); err != nil {
					fmt.Printf("❌ Error: %v\n", err)
				}
				
				// Show refresh info
				fmt.Printf("\n🔄 Refreshing every 30 seconds... (Press Ctrl+C to stop)\n")
				
				// Wait 30 seconds
				time.Sleep(30 * time.Second)
			}
		} else {
			// Single run
			return displayOpenAIData(provider, period, false)
		}
	},
}

func init() {
	openaiCmd.Flags().StringP("period", "p", "7d", "Time period (7d, 30d, 90d)")
	openaiCmd.Flags().BoolP("watch", "w", false, "Watch mode - refresh every 30 seconds")
	RootCmd.AddCommand(openaiCmd)
}

// displayOpenAIData fetches and displays OpenAI usage data
func displayOpenAIData(provider *providers.OpenAIProvider, period string, bypassCache bool) error {
	// Display header
	fmt.Printf("🤖 OPENAI USAGE - Last %s\n", period)
	fmt.Printf("⏰ Generated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Get time range
	startTime, endTime := providers.GetPeriodTimeRange(period)

	// Fetch consumption data
	consumptions, err := provider.GetConsumption(startTime, endTime, bypassCache)
	if err != nil {
		return fmt.Errorf("failed to get consumption data: %w", err)
	}

	// Fetch pricing data
	pricings, err := provider.GetPricing(startTime, endTime, bypassCache)
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
		if period == "30d" || period == "90d" {
			fmt.Printf("💰 Cost Data: %s\n", color.YellowString("Not available for this period"))
		} else {
			fmt.Printf("💰 Cost Data: %s\n", color.GreenString("Free Tier"))
		}
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
