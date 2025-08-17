package models

import (
	"time"
)

// Pricing represents cost data from any platform
type Pricing struct {
	Platform  string    `json:"platform"`
	Model     string    `json:"model"`
	LineItem  string    `json:"line_item"` // e.g., "gpt-4o-input", "gpt-4o-output"
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Timestamp time.Time `json:"timestamp"`
}

// PricingSummary represents aggregated pricing data
type PricingSummary struct {
	Platform  string    `json:"platform"`
	Model     string    `json:"model"`
	TotalCost float64   `json:"total_cost"`
	Currency  string    `json:"currency"`
	Period    string    `json:"period"` // e.g., "7d", "30d"
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	LineItems []Pricing `json:"line_items"` // Breakdown by line item
}

// NewPricing creates a new Pricing instance
func NewPricing(platform, model, lineItem string, amount float64, currency string, startTime, endTime time.Time) *Pricing {
	return &Pricing{
		Platform:  platform,
		Model:     model,
		LineItem:  lineItem,
		Amount:    amount,
		Currency:  currency,
		StartTime: startTime,
		EndTime:   endTime,
		Timestamp: time.Now(),
	}
}

// NewPricingSummary creates a new PricingSummary instance
func NewPricingSummary(platform, model, period string, startTime, endTime time.Time) *PricingSummary {
	return &PricingSummary{
		Platform:  platform,
		Model:     model,
		Period:    period,
		StartTime: startTime,
		EndTime:   endTime,
		LineItems: make([]Pricing, 0),
	}
}

// AddPricing adds pricing data to the summary
func (ps *PricingSummary) AddPricing(pricing *Pricing) {
	ps.TotalCost += pricing.Amount
	if ps.Currency == "" {
		ps.Currency = pricing.Currency
	}
	ps.LineItems = append(ps.LineItems, *pricing)
}
