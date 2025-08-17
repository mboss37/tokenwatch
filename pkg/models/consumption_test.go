package models

import (
	"testing"
	"time"
)

func TestNewConsumptionSummary(t *testing.T) {
	now := time.Now()
	endTime := now.Add(24 * time.Hour)
	
	summary := NewConsumptionSummary("openai", "gpt-4", "7d", now, endTime)
	
	if summary.Platform != "openai" {
		t.Errorf("Expected platform 'openai', got '%s'", summary.Platform)
	}
	if summary.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", summary.Model)
	}
	if summary.Period != "7d" {
		t.Errorf("Expected period '7d', got '%s'", summary.Period)
	}
	if summary.TotalInputTokens != 0 {
		t.Errorf("Expected initial input tokens to be 0, got %d", summary.TotalInputTokens)
	}
}

func TestConsumptionSummary_AddConsumption(t *testing.T) {
	now := time.Now()
	endTime := now.Add(24 * time.Hour)
	
	summary := NewConsumptionSummary("openai", "gpt-4", "7d", now, endTime)
	
	// Add first consumption
	consumption1 := &Consumption{
		Platform:     "openai",
		Model:        "gpt-4",
		InputTokens:  100,
		OutputTokens: 50,
		TotalTokens:  150,
		RequestCount: 1,
		StartTime:    now,
		EndTime:      now.Add(1 * time.Hour),
	}
	summary.AddConsumption(consumption1)
	
	if summary.TotalInputTokens != 100 {
		t.Errorf("Expected input tokens 100, got %d", summary.TotalInputTokens)
	}
	if summary.TotalOutputTokens != 50 {
		t.Errorf("Expected output tokens 50, got %d", summary.TotalOutputTokens)
	}
	if summary.TotalTokens != 150 {
		t.Errorf("Expected total tokens 150, got %d", summary.TotalTokens)
	}
	if summary.TotalRequests != 1 {
		t.Errorf("Expected requests 1, got %d", summary.TotalRequests)
	}
	
	// Add second consumption
	consumption2 := &Consumption{
		Platform:     "openai",
		Model:        "gpt-4",
		InputTokens:  200,
		OutputTokens: 100,
		TotalTokens:  300,
		RequestCount: 2,
		StartTime:    now.Add(1 * time.Hour),
		EndTime:      now.Add(2 * time.Hour),
	}
	summary.AddConsumption(consumption2)
	
	if summary.TotalInputTokens != 300 {
		t.Errorf("Expected input tokens 300, got %d", summary.TotalInputTokens)
	}
	if summary.TotalOutputTokens != 150 {
		t.Errorf("Expected output tokens 150, got %d", summary.TotalOutputTokens)
	}
	if summary.TotalTokens != 450 {
		t.Errorf("Expected total tokens 450, got %d", summary.TotalTokens)
	}
	if summary.TotalRequests != 3 {
		t.Errorf("Expected requests 3, got %d", summary.TotalRequests)
	}
}
