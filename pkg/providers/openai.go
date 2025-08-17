package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tokenwatch/pkg/models"
	"tokenwatch/pkg/utils"
)

// OpenAIProvider handles API calls to OpenAI
type OpenAIProvider struct {
	client         *utils.RateLimitedClient
	circuitBreaker *utils.CircuitBreaker
	apiKey         string
	baseURL        string
	orgID          string
	cache          map[string]cacheItem
	cacheTTL       time.Duration
}

// cacheItem represents a cached API response
type cacheItem struct {
	data      interface{}
	expiresAt time.Time
}

// OpenAIUsageResponse represents the usage API response
type OpenAIUsageResponse struct {
	Data     []OpenAIUsageBucket `json:"data"`
	HasMore  bool                `json:"has_more"`
	NextPage string              `json:"next_page"`
}

type OpenAIUsageBucket struct {
	StartTime int64               `json:"start_time"`
	EndTime   int64               `json:"end_time"`
	Results   []OpenAIUsageResult `json:"results"`
}

type OpenAIUsageResult struct {
	Model            string `json:"model"`
	InputTokens      int64  `json:"input_tokens"`
	OutputTokens     int64  `json:"output_tokens"`
	NumModelRequests int64  `json:"num_model_requests"`
}

// OpenAICostResponse represents the costs API response
type OpenAICostResponse struct {
	Data     []OpenAICostBucket `json:"data"`
	HasMore  bool               `json:"has_more"`
	NextPage string             `json:"next_page"`
}

type OpenAICostBucket struct {
	StartTime int64              `json:"start_time"`
	EndTime   int64              `json:"end_time"`
	Results   []OpenAICostResult `json:"results"`
}

type OpenAICostResult struct {
	LineItem string           `json:"line_item"`
	Amount   OpenAICostAmount `json:"amount"`
}

type OpenAICostAmount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// NewOpenAIProvider creates a new OpenAI provider instance
func NewOpenAIProvider(apiKey, orgID string) *OpenAIProvider {
	// Default cache TTL of 5 minutes for normal operations
	cacheTTL := 5 * time.Minute

	// Rate limiting: OpenAI has various rate limits, using conservative defaults
	// 60 requests per minute = 1 request per second with burst of 5
	rateLimitedClient := utils.NewRateLimitedClient(1.0, 5, 30*time.Second)

	// Circuit breaker: Open after 5 consecutive failures, reset after 1 minute
	circuitBreaker := utils.NewCircuitBreaker(5, 1*time.Minute)

	return &OpenAIProvider{
		client:         rateLimitedClient,
		circuitBreaker: circuitBreaker,
		apiKey:         apiKey,
		baseURL:        "https://api.openai.com/v1",
		orgID:          orgID,
		cache:          make(map[string]cacheItem),
		cacheTTL:       cacheTTL,
	}
}

// GetPlatform returns the platform name
func (o *OpenAIProvider) GetPlatform() string {
	return "openai"
}

// IsAvailable checks if the provider is properly configured
func (o *OpenAIProvider) IsAvailable() bool {
	return o.apiKey != ""
}

// ClearCache clears all cached data
func (o *OpenAIProvider) ClearCache() {
	o.cache = make(map[string]cacheItem)
}

// GetConsumption retrieves consumption data and converts to common models
func (o *OpenAIProvider) GetConsumption(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Consumption, error) {
	usageResp, err := o.GetUsage(startTime, endTime, "1d", []string{"model"}, bypassCache, debug)
	if err != nil {
		return nil, err
	}

	var consumptions []*models.Consumption
	for _, bucket := range usageResp.Data {
		for _, result := range bucket.Results {
			consumption := models.NewConsumption(
				o.GetPlatform(),
				result.Model,
				result.InputTokens,
				result.OutputTokens,
				result.NumModelRequests,
				time.Unix(bucket.StartTime, 0),
				time.Unix(bucket.EndTime, 0),
			)
			consumptions = append(consumptions, consumption)
		}
	}

	return consumptions, nil
}

// GetPricing retrieves pricing data and converts to common models
func (o *OpenAIProvider) GetPricing(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Pricing, error) {
	costResp, err := o.GetCosts(startTime, endTime, []string{"line_item"}, bypassCache, debug)
	if err != nil {
		return nil, err
	}

	var pricings []*models.Pricing
	for _, bucket := range costResp.Data {
		for _, result := range bucket.Results {
			// Extract model from line item (e.g., "gpt-4o-input" -> "gpt-4o")
			model := o.extractModelFromLineItem(result.LineItem)

			pricing := models.NewPricing(
				o.GetPlatform(),
				model,
				result.LineItem,
				result.Amount.Value,
				result.Amount.Currency,
				time.Unix(bucket.StartTime, 0),
				time.Unix(bucket.EndTime, 0),
			)
			pricings = append(pricings, pricing)
		}
	}

	return pricings, nil
}

// GetConsumptionSummary gets aggregated consumption data for common periods
func (o *OpenAIProvider) GetConsumptionSummary(period string) (*models.ConsumptionSummary, error) {
	startTime, endTime := GetPeriodTimeRange(period)

	consumptions, err := o.GetConsumption(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}

	// Group by model
	modelSummaries := make(map[string]*models.ConsumptionSummary)
	for _, consumption := range consumptions {
		summary, exists := modelSummaries[consumption.Model]
		if !exists {
			summary = models.NewConsumptionSummary(o.GetPlatform(), consumption.Model, period, startTime, endTime)
			modelSummaries[consumption.Model] = summary
		}
		summary.AddConsumption(consumption)
	}

	// Return first summary for now (we can enhance this later)
	for _, summary := range modelSummaries {
		return summary, nil
	}

	// Return empty summary if no data
	return models.NewConsumptionSummary(o.GetPlatform(), "", period, startTime, endTime), nil
}

// GetPricingSummary gets aggregated pricing data for common periods
func (o *OpenAIProvider) GetPricingSummary(period string) (*models.PricingSummary, error) {
	startTime, endTime := GetPeriodTimeRange(period)

	pricings, err := o.GetPricing(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}

	// Group by model
	modelSummaries := make(map[string]*models.PricingSummary)
	for _, pricing := range pricings {
		summary, exists := modelSummaries[pricing.Model]
		if !exists {
			summary = models.NewPricingSummary(o.GetPlatform(), pricing.Model, period, startTime, endTime)
			modelSummaries[pricing.Model] = summary
		}
		summary.AddPricing(pricing)
	}

	// Return first summary for now (we can enhance this later)
	for _, summary := range modelSummaries {
		return summary, nil
	}

	// Return empty summary if no data
	return models.NewPricingSummary(o.GetPlatform(), "", period, startTime, endTime), nil
}

// extractModelFromLineItem extracts the model name from OpenAI's line item format
func (o *OpenAIProvider) extractModelFromLineItem(lineItem string) string {
	// OpenAI line items format: "model, type" e.g., "gpt-4o-2024-08-06, input"
	// We want to extract just the model name
	parts := strings.Split(lineItem, ", ")
	if len(parts) > 0 {
		return parts[0]
	}
	return lineItem
}

// getCacheKey generates a cache key for the given request parameters
func (o *OpenAIProvider) getCacheKey(endpoint string, params map[string]string) string {
	key := endpoint
	for k, v := range params {
		key += fmt.Sprintf(":%s=%s", k, v)
	}
	return key
}

// getFromCache attempts to retrieve data from cache
func (o *OpenAIProvider) getFromCache(key string, result interface{}) bool {
	item, found := o.cache[key]
	if !found {
		return false
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		delete(o.cache, key)
		return false
	}

	// Copy cached data to result
	switch data := item.data.(type) {
	case *OpenAIUsageResponse:
		if resp, ok := result.(**OpenAIUsageResponse); ok {
			*resp = data
			return true
		}
	case *OpenAICostResponse:
		if resp, ok := result.(**OpenAICostResponse); ok {
			*resp = data
			return true
		}
	}

	return false
}

// saveToCache stores data in cache
func (o *OpenAIProvider) saveToCache(key string, data interface{}) {
	o.cache[key] = cacheItem{
		data:      data,
		expiresAt: time.Now().Add(o.cacheTTL),
	}
}

// countTotalResults counts the total number of results across all buckets
func countTotalResults(buckets []OpenAIUsageBucket) int {
	total := 0
	for _, bucket := range buckets {
		total += len(bucket.Results)
	}
	return total
}

// countTotalCostResults counts the total number of results across all buckets for costs
func countTotalCostResults(buckets []OpenAICostBucket) int {
	total := 0
	for _, bucket := range buckets {
		total += len(bucket.Results)
	}
	return total
}

// GetUsage retrieves token usage data from OpenAI (internal method)
func (o *OpenAIProvider) GetUsage(startTime, endTime time.Time, bucketWidth string, groupBy []string, bypassCache bool, debug bool) (*OpenAIUsageResponse, error) {
	// Create cache key
	params := map[string]string{
		"start_time":   fmt.Sprintf("%d", startTime.Unix()),
		"end_time":     fmt.Sprintf("%d", endTime.Unix()),
		"bucket_width": bucketWidth,
	}
	for i, group := range groupBy {
		params[fmt.Sprintf("group_by_%d", i)] = group
	}
	cacheKey := o.getCacheKey("usage", params)

	// Try to get from cache (unless bypassing)
	if !bypassCache {
		var result *OpenAIUsageResponse
		if o.getFromCache(cacheKey, &result) {
			return result, nil
		}
	}

	// Not in cache or bypassing cache, make API request with pagination
	var allData []OpenAIUsageBucket
	var nextPage string
	maxPages := 50 // Safety limit to prevent infinite loops
	pageCount := 0
	seenPages := make(map[string]bool) // Track seen pages to detect loops

	for pageCount < maxPages {
		pageCount++

		// Build URL with query parameters
		url := fmt.Sprintf("%s/organization/usage/completions", o.baseURL)

		// Add timeout for longer periods to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create HTTP request with context
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add query parameters
		q := req.URL.Query()
		q.Add("start_time", fmt.Sprintf("%d", startTime.Unix()))
		if !endTime.IsZero() {
			q.Add("end_time", fmt.Sprintf("%d", endTime.Unix()))
		}
		if bucketWidth != "" {
			q.Add("bucket_width", bucketWidth)
		}
		if len(groupBy) > 0 {
			for _, group := range groupBy {
				q.Add("group_by", group)
			}
		}
		// Add next_page token if we have one
		if nextPage != "" {
			q.Add("page", nextPage)
		}
		req.URL.RawQuery = q.Encode()

		// Log request details for debugging (only when debug is enabled)
		if debug {
			fmt.Printf("🔍 OPENAI USAGE API REQUEST (Page %d, Token: %s):\n", pageCount, nextPage)
			fmt.Printf("   URL: %s\n", req.URL.String())
			fmt.Printf("   Start Time: %s (%d)\n", startTime.Format("2006-01-02 15:04:05"), startTime.Unix())
			fmt.Printf("   End Time: %s (%d)\n", endTime.Format("2006-01-02 15:04:05"), endTime.Unix())
			fmt.Printf("   Bucket Width: %s\n", bucketWidth)
			fmt.Printf("   Group By: %v\n", groupBy)
			if nextPage != "" {
				fmt.Printf("   Next Page: %s\n", nextPage)
			}
			fmt.Println()
		}

		// Add headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))
		if o.orgID != "" {
			req.Header.Set("OpenAI-Organization", o.orgID)
		}

		// Make request with circuit breaker
		var resp *http.Response
		err = o.circuitBreaker.Call(func() error {
			var reqErr error
			resp, reqErr = o.client.Do(req)
			if reqErr != nil {
				return fmt.Errorf("failed to make request: %w", reqErr)
			}

			if resp.StatusCode != http.StatusOK {
				defer resp.Body.Close()
				return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Parse response
		var usageResp OpenAIUsageResponse
		if err := json.NewDecoder(resp.Body).Decode(&usageResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		// Log raw response for debugging
		rawJSON, _ := json.MarshalIndent(usageResp, "", "  ")
		if debug {
			fmt.Printf("🔍 RAW OPENAI USAGE API RESPONSE (Page %d, Token: %s):\n", pageCount, nextPage)
			fmt.Printf("   Has More: %v\n", usageResp.HasMore)
			if usageResp.HasMore {
				fmt.Printf("   Next Page: %s\n", usageResp.NextPage)
			}
			fmt.Printf("   Data Buckets: %d\n", len(usageResp.Data))
			fmt.Printf("   Total Results: %d\n", countTotalResults(usageResp.Data))
			fmt.Printf("%s\n\n", string(rawJSON))
		}

		// Append results to allData
		allData = append(allData, usageResp.Data...)

		// Check if there's a next page
		if !usageResp.HasMore {
			if debug {
				fmt.Printf("🔍 PAGINATION COMPLETE: Fetched %d pages, %d total buckets\n",
					pageCount, len(allData))
			}
			break
		}

		// Check for pagination loops
		if usageResp.NextPage == "" {
			if debug {
				fmt.Printf("⚠️  WARNING: API returned has_more=true but no next_page token\n")
			}
			break
		}

		// Check if we've seen this page token before (loop detection)
		if seenPages[usageResp.NextPage] {
			if debug {
				fmt.Printf("⚠️  WARNING: Detected pagination loop at page %d, stopping\n", pageCount)
			}
			break
		}
		seenPages[usageResp.NextPage] = true

		nextPage = usageResp.NextPage
		if debug {
			fmt.Printf("🔍 FETCHING NEXT PAGE: %s\n\n", nextPage)
		}
	}

	// Safety check - if we hit max pages, log a warning
	if pageCount >= maxPages {
		if debug {
			fmt.Printf("⚠️  WARNING: Hit maximum page limit (%d), stopping pagination\n", maxPages)
		}
	}

	// Cache the result
	o.saveToCache(cacheKey, &OpenAIUsageResponse{Data: allData})

	return &OpenAIUsageResponse{Data: allData}, nil
}

// GetCosts retrieves cost data from OpenAI (internal method)
func (o *OpenAIProvider) GetCosts(startTime, endTime time.Time, groupBy []string, bypassCache bool, debug bool) (*OpenAICostResponse, error) {
	// Create cache key
	params := map[string]string{
		"start_time":   fmt.Sprintf("%d", startTime.Unix()),
		"end_time":     fmt.Sprintf("%d", endTime.Unix()),
		"bucket_width": "1d", // Costs API only supports daily buckets
	}
	for i, group := range groupBy {
		params[fmt.Sprintf("group_by_%d", i)] = group
	}
	cacheKey := o.getCacheKey("costs", params)

	// Try to get from cache (unless bypassing)
	if !bypassCache {
		var result *OpenAICostResponse
		if o.getFromCache(cacheKey, &result) {
			return result, nil
		}
	}

	// Not in cache or bypassing cache, make API request with pagination
	var allData []OpenAICostBucket
	var nextPage string
	maxPages := 50 // Safety limit to prevent infinite loops
	pageCount := 0
	seenPages := make(map[string]bool) // Track seen pages to detect loops

	for pageCount < maxPages {
		pageCount++

		// Build URL with query parameters
		url := fmt.Sprintf("%s/organization/costs", o.baseURL)

		// Add timeout for longer periods to prevent hanging
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Create HTTP request with context
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add query parameters
		q := req.URL.Query()
		q.Add("start_time", fmt.Sprintf("%d", startTime.Unix()))
		if !endTime.IsZero() {
			q.Add("end_time", fmt.Sprintf("%d", endTime.Unix()))
		}
		q.Add("bucket_width", "1d") // Costs API only supports daily buckets
		if len(groupBy) > 0 {
			for _, group := range groupBy {
				q.Add("group_by", group)
			}
		}
		// Add next_page token if we have one
		if nextPage != "" {
			q.Add("page", nextPage)
		}
		req.URL.RawQuery = q.Encode()

		// Log request details for debugging (only when debug is enabled)
		if debug {
			fmt.Printf("🔍 OPENAI COSTS API REQUEST (Page %d, Token: %s):\n", pageCount, nextPage)
			fmt.Printf("   URL: %s\n", req.URL.String())
			fmt.Printf("   Start Time: %s (%d)\n", startTime.Format("2006-01-02 15:04:05"), startTime.Unix())
			fmt.Printf("   End Time: %s (%d)\n", endTime.Format("2006-01-02 15:04:05"), endTime.Unix())
			fmt.Printf("   Group By: %v\n", groupBy)
			if nextPage != "" {
				fmt.Printf("   Next Page: %s\n", nextPage)
			}
			fmt.Println()
		}

		// Add headers
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))
		if o.orgID != "" {
			req.Header.Set("OpenAI-Organization", o.orgID)
		}

		// Make request with circuit breaker
		var resp *http.Response
		err = o.circuitBreaker.Call(func() error {
			var reqErr error
			resp, reqErr = o.client.Do(req)
			if reqErr != nil {
				return fmt.Errorf("failed to make request: %w", reqErr)
			}

			if resp.StatusCode != http.StatusOK {
				defer resp.Body.Close()
				return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Parse response
		var costResp OpenAICostResponse
		if err := json.NewDecoder(resp.Body).Decode(&costResp); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		// Log raw response for debugging
		rawJSON, _ := json.MarshalIndent(costResp, "", "  ")
		if debug {
			fmt.Printf("🔍 RAW OPENAI COSTS API RESPONSE (Page %d, Token: %s):\n", pageCount, nextPage)
			fmt.Printf("   Has More: %v\n", costResp.HasMore)
			if costResp.HasMore {
				fmt.Printf("   Next Page: %s\n", costResp.NextPage)
			}
			fmt.Printf("   Data Buckets: %d\n", len(costResp.Data))
			fmt.Printf("   Total Results: %d\n", countTotalCostResults(costResp.Data))
			fmt.Printf("%s\n\n", string(rawJSON))
		}

		// Append results to allData
		allData = append(allData, costResp.Data...)

		// Check if there's a next page
		if !costResp.HasMore {
			if debug {
				fmt.Printf("🔍 PAGINATION COMPLETE: Fetched %d pages, %d total buckets\n",
					pageCount, len(allData))
			}
			break
		}

		// Check for pagination loops
		if costResp.NextPage == "" {
			if debug {
				fmt.Printf("⚠️  WARNING: API returned has_more=true but no next_page token\n")
			}
			break
		}

		// Check if we've seen this page token before (loop detection)
		if seenPages[costResp.NextPage] {
			if debug {
				fmt.Printf("⚠️  WARNING: Detected pagination loop at page %d, stopping\n", pageCount)
			}
			break
		}
		seenPages[costResp.NextPage] = true

		nextPage = costResp.NextPage
		if debug {
			fmt.Printf("🔍 FETCHING NEXT PAGE: %s\n\n", nextPage)
		}
	}

	// Safety check - if we hit max pages, log a warning
	if pageCount >= maxPages {
		if debug {
			fmt.Printf("⚠️  WARNING: Hit maximum page limit (%d), stopping pagination\n", maxPages)
		}
	}

	// Cache the result
	o.saveToCache(cacheKey, &OpenAICostResponse{Data: allData})

	return &OpenAICostResponse{Data: allData}, nil
}

// Legacy methods for backward compatibility
func (o *OpenAIProvider) GetLast7DaysUsage() (*OpenAIUsageResponse, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -7)

	return o.GetUsage(startTime, endTime, "1d", []string{"model"}, false, false)
}

func (o *OpenAIProvider) GetLast30DaysCosts() (*OpenAICostResponse, error) {
	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -30)

	return o.GetCosts(startTime, endTime, []string{"line_item"}, false, false)
}
