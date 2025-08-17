package models

import (
	"math"
	"testing"
	"time"
)

func TestNewPricingSummary(t *testing.T) {
	now := time.Now()
	endTime := now.Add(24 * time.Hour)
	
	summary := NewPricingSummary("openai", "gpt-4", "7d", now, endTime)
	
	if summary.Platform != "openai" {
		t.Errorf("Expected platform 'openai', got '%s'", summary.Platform)
	}
	if summary.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", summary.Model)
	}
	if summary.Period != "7d" {
		t.Errorf("Expected period '7d', got '%s'", summary.Period)
	}
	// Currency is set when pricing is added, not in constructor
	if summary.Currency != "" {
		t.Errorf("Expected empty currency initially, got '%s'", summary.Currency)
	}
	if summary.TotalCost != 0 {
		t.Errorf("Expected initial cost to be 0, got %f", summary.TotalCost)
	}
}

func TestPricingSummary_AddPricing(t *testing.T) {
	now := time.Now()
	endTime := now.Add(24 * time.Hour)
	
	summary := NewPricingSummary("openai", "gpt-4", "7d", now, endTime)
	
	// Add first pricing
	pricing1 := &Pricing{
		Platform:  "openai",
		Model:     "gpt-4",
		LineItem:  "gpt-4, input",
		Amount:    0.05,
		Currency:  "USD",
		StartTime: now,
		EndTime:   now.Add(1 * time.Hour),
	}
	summary.AddPricing(pricing1)
	
	if summary.TotalCost != 0.05 {
		t.Errorf("Expected total cost 0.05, got %f", summary.TotalCost)
	}
	if len(summary.LineItems) != 1 {
		t.Errorf("Expected 1 line item, got %d", len(summary.LineItems))
	}
	
	// Add second pricing
	pricing2 := &Pricing{
		Platform:  "openai",
		Model:     "gpt-4",
		LineItem:  "gpt-4, output",
		Amount:    0.10,
		Currency:  "USD",
		StartTime: now.Add(1 * time.Hour),
		EndTime:   now.Add(2 * time.Hour),
	}
	summary.AddPricing(pricing2)
	
	expectedCost := 0.15
	if math.Abs(summary.TotalCost - expectedCost) > 0.0001 {
		t.Errorf("Expected total cost %.4f, got %.4f", expectedCost, summary.TotalCost)
	}
	if len(summary.LineItems) != 2 {
		t.Errorf("Expected 2 line items, got %d", len(summary.LineItems))
	}
	
	// Verify line items
	expectedItems := map[string]bool{
		"gpt-4, input":  false,
		"gpt-4, output": false,
	}
	
	for _, item := range summary.LineItems {
		if _, exists := expectedItems[item.LineItem]; exists {
			expectedItems[item.LineItem] = true
		} else {
			t.Errorf("Unexpected line item: %s", item.LineItem)
		}
	}
	
	for item, found := range expectedItems {
		if !found {
			t.Errorf("Expected line item not found: %s", item)
		}
	}
}
