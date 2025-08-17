package providers

import (
	"time"

	"tokenwatch/pkg/models"
)

// Provider defines the interface that all platform providers must implement
type Provider interface {
	// GetPlatform returns the platform name (e.g., "openai")
	GetPlatform() string

	// GetConsumption retrieves consumption data for a specific time period
	GetConsumption(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Consumption, error)

	// GetPricing retrieves pricing data for a specific time period
	GetPricing(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Pricing, error)

	// GetConsumptionSummary gets aggregated consumption data for common periods
	GetConsumptionSummary(period string) (*models.ConsumptionSummary, error)

	// GetPricingSummary gets aggregated pricing data for common periods
	GetPricingSummary(period string) (*models.PricingSummary, error)

	// IsAvailable checks if the provider is properly configured and available
	IsAvailable() bool
}

// Common periods that providers should support
const (
	Period7Days  = "7d"
	Period30Days = "30d"
	Period90Days = "90d"
	Period1Year  = "1y"
	PeriodAll    = "all"
)

// GetPeriodTimeRange returns start and end times for common periods
func GetPeriodTimeRange(period string) (time.Time, time.Time) {
	endTime := time.Now()
	var startTime time.Time

	switch period {
	case "1d":
		startTime = endTime.AddDate(0, 0, -1) // Last 24 hours
	case "7d":
		startTime = endTime.AddDate(0, 0, -7) // Last 7 days
	case "30d":
		startTime = endTime.AddDate(0, 0, -30) // Last 30 days
	case "90d":
		startTime = endTime.AddDate(0, 0, -90) // Last 90 days
	case "1y":
		startTime = endTime.AddDate(0, 0, -365) // Last 1 year
	case "all":
		startTime = endTime.AddDate(0, 0, -1825) // Last 5 years (maximum practical limit)
	default:
		// Default to 7 days if period is not recognized
		startTime = endTime.AddDate(0, 0, -7)
	}

	return startTime, endTime
}
