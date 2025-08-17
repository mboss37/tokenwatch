package providers

import (
	"time"

	"tokenwatch/pkg/models"
)

// Provider defines the interface that all platform providers must implement
type Provider interface {
	// GetPlatform returns the platform name (e.g., "openai", "grok")
	GetPlatform() string

	// GetConsumption retrieves consumption data for a specific time period
	GetConsumption(startTime, endTime time.Time, bypassCache bool) ([]*models.Consumption, error)

	// GetPricing retrieves pricing data for a specific time period
	GetPricing(startTime, endTime time.Time, bypassCache bool) ([]*models.Pricing, error)

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
)

// GetPeriodTimeRange returns start and end times for common periods
func GetPeriodTimeRange(period string) (time.Time, time.Time) {
	endTime := time.Now()
	var startTime time.Time

	switch period {
	case Period7Days:
		startTime = endTime.AddDate(0, 0, -7)
	case Period30Days:
		startTime = endTime.AddDate(0, 0, -30)
	case Period90Days:
		startTime = endTime.AddDate(0, 0, -90)
	default:
		// Default to 7 days
		startTime = endTime.AddDate(0, 0, -7)
	}

	return startTime, endTime
}
