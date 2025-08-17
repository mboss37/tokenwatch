package models

import (
	"time"
)

// Consumption represents token usage data from any platform
type Consumption struct {
	Platform     string    `json:"platform"`
	Model        string    `json:"model"`
	InputTokens  int64     `json:"input_tokens"`
	OutputTokens int64     `json:"output_tokens"`
	TotalTokens  int64     `json:"total_tokens"`
	RequestCount int64     `json:"request_count"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Timestamp    time.Time `json:"timestamp"`
}

// ConsumptionSummary represents aggregated consumption data
type ConsumptionSummary struct {
	Platform          string    `json:"platform"`
	Model             string    `json:"model"`
	TotalInputTokens  int64     `json:"total_input_tokens"`
	TotalOutputTokens int64     `json:"total_output_tokens"`
	TotalTokens       int64     `json:"total_tokens"`
	TotalRequests     int64     `json:"total_requests"`
	Period            string    `json:"period"` // e.g., "7d", "30d"
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
}

// NewConsumption creates a new Consumption instance
func NewConsumption(platform, model string, inputTokens, outputTokens, requestCount int64, startTime, endTime time.Time) *Consumption {
	return &Consumption{
		Platform:     platform,
		Model:        model,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalTokens:  inputTokens + outputTokens,
		RequestCount: requestCount,
		StartTime:    startTime,
		EndTime:      endTime,
		Timestamp:    time.Now(),
	}
}

// NewConsumptionSummary creates a new ConsumptionSummary instance
func NewConsumptionSummary(platform, model string, period string, startTime, endTime time.Time) *ConsumptionSummary {
	return &ConsumptionSummary{
		Platform:  platform,
		Model:     model,
		Period:    period,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// AddConsumption adds consumption data to the summary
func (cs *ConsumptionSummary) AddConsumption(consumption *Consumption) {
	cs.TotalInputTokens += consumption.InputTokens
	cs.TotalOutputTokens += consumption.OutputTokens
	cs.TotalTokens += consumption.TotalTokens
	cs.TotalRequests += consumption.RequestCount
}
