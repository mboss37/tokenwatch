# TokenWatch Developer Guide

## Architecture Overview

TokenWatch follows a clean, modular architecture with strict platform separation.

```
┌─────────────────────────────────────────────────────┐
│                 CLI Commands                        │
│  (tokenwatch openai, tokenwatch all, etc.)         │
└─────────────────────────┬───────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────┐
│              Provider Interface                     │
│  (Common contract for all platforms)               │
└─────────────────────────┬───────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────┐
│           Platform Providers                        │
│  (OpenAI, Anthropic*, Grok*, Cursor*)              │
└─────────────────────────┬───────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────┐
│              External APIs                          │
│  (OpenAI API, Anthropic API, etc.)                 │
└─────────────────────────────────────────────────────┘
```

## Core Principles

### 1. Platform Separation (Most Important!)

Each platform is **completely isolated**:

```
✅ GOOD:
cmd/root/
├── openai_cmd.go      # Only OpenAI logic
├── anthropic_cmd.go   # Only Anthropic logic
└── grok_cmd.go        # Only Grok logic

❌ BAD:
cmd/root/
└── platforms_cmd.go   # Mixed platform logic - DON'T DO THIS!
```

**Why?** 
- Easy debugging (issues isolated to one file)
- Safe parallel development
- Simple testing
- No platform interference

### 2. Common Interface

All platforms implement:
```go
type Provider interface {
    GetPlatform() string
    GetConsumptionSummary(period string) (*models.ConsumptionSummary, error)
    GetPricingSummary(period string) (*models.PricingSummary, error)
    IsAvailable() bool
}
```

### 3. Error Handling

Use structured errors with helpful suggestions:
```go
// Good
return utils.NewAPIError("OpenAI request failed", 401, err)
// This automatically adds suggestions like "Check your API key"

// Bad
return fmt.Errorf("request failed")
```

## Development Workflow

### Running Locally

```bash
# Build
make build

# Run
./build/tokenwatch openai

# Debug mode
export TOKENWATCH_LOG_LEVEL=debug
./build/tokenwatch openai
```

### Adding a New Platform (Step-by-Step)

#### 1. Create the Provider
`pkg/providers/newplatform.go`:
```go
package providers

import (
    "tokenwatch/pkg/models"
    "tokenwatch/pkg/utils"
)

type NewPlatformProvider struct {
    client         *utils.RateLimitedClient
    circuitBreaker *utils.CircuitBreaker
    apiKey         string
}

func NewNewPlatformProvider(apiKey string) *NewPlatformProvider {
    return &NewPlatformProvider{
        client:         utils.NewRateLimitedClient(1.0, 5, 30*time.Second),
        circuitBreaker: utils.NewCircuitBreaker(5, 1*time.Minute),
        apiKey:         apiKey,
    }
}

func (p *NewPlatformProvider) GetPlatform() string {
    return "newplatform"
}

func (p *NewPlatformProvider) GetConsumptionSummary(period string) (*models.ConsumptionSummary, error) {
    // Implement API call
    // Use circuit breaker for resilience
    // Return data in common format
}

// ... implement other interface methods
```

#### 2. Create the Command
`cmd/root/newplatform_cmd.go`:
```go
package main

import (
    "tokenwatch/internal/config"
    "tokenwatch/pkg/providers"
    "tokenwatch/pkg/utils"
    "github.com/spf13/cobra"
)

var newplatformCmd = &cobra.Command{
    Use:   "newplatform",
    Short: "Show NewPlatform token consumption and costs",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Check configuration
        apiKey := config.GetAPIKey("newplatform")
        if apiKey == "" {
            return utils.NewAuthError("NewPlatform not configured", "newplatform")
        }
        
        // Create provider
        provider := providers.NewNewPlatformProvider(apiKey)
        
        // Get data and display
        // ... implementation
    },
}

func init() {
    newplatformCmd.Flags().StringP("period", "p", "7d", "Time period")
    RootCmd.AddCommand(newplatformCmd)
}
```

#### 3. Update Setup Command
In `setup_cmd.go`, add to platforms list:
```go
{"newplatform", "NewPlatform", "Description", true}
```

#### 4. Update Provider Factory
In `overview_cmd.go`, add to `getProvider()`:
```go
case "newplatform":
    return providers.NewNewPlatformProvider(apiKey)
```

### Testing Your Changes

```bash
# Build
go build ./cmd/root

# Test setup
./tokenwatch setup
# Select your new platform

# Test command
./tokenwatch newplatform

# Test in 'all' command
./tokenwatch all
```

## Code Style Guide

### Imports
Group imports in order:
1. Standard library
2. External packages
3. Internal packages

```go
import (
    "fmt"
    "time"
    
    "github.com/spf13/cobra"
    
    "tokenwatch/pkg/models"
    "tokenwatch/pkg/utils"
)
```

### Error Handling
- Always use structured errors for user-facing errors
- Add context with `fmt.Errorf("context: %w", err)`
- Provide actionable suggestions

### Logging
- Use structured logging
- Debug level for detailed info
- Info level for important events
- Error level for failures

```go
utils.Debug("API call started", map[string]interface{}{
    "platform": "openai",
    "endpoint": "usage",
})
```

## Common Patterns

### Rate Limiting
Already built-in via `utils.RateLimitedClient`

### Retry Logic
Automatic with exponential backoff

### Circuit Breaker
Prevents cascading failures:
```go
err := provider.circuitBreaker.Call(func() error {
    // API call here
})
```

### Caching
Implement in provider if needed:
```go
if cached := getFromCache(key); cached != nil {
    return cached
}
// Make API call
saveToCache(key, result)
```

## Debugging

### Enable Debug Logging
```bash
export TOKENWATCH_LOG_LEVEL=debug
```

### Common Issues

**"circuit breaker is open"**
- Too many failures, wait 1 minute
- Check API status
- Verify credentials

**Rate limit errors**
- Automatic retry with backoff
- Check API plan limits

**Network errors**
- Check internet connection
- Verify API endpoint
- Check firewall/proxy

## Project Maintenance

### Dependencies
```bash
# Update dependencies
go get -u ./...
go mod tidy
```

### Building
```bash
# Local build
make build

# Cross-platform
make build-all
```

### Features Implemented

- ✅ **Core Infrastructure**: Provider interface, common models, config management
- ✅ **Production Features**: Retry logic, rate limiting, circuit breaker, caching
- ✅ **Watch Mode**: Real-time monitoring with `-w` flag
- ✅ **Enhanced UX**: API key validation, masked input, colored output
- ✅ **Error Handling**: Structured errors with actionable suggestions
- ✅ **Logging**: Structured logging with debug support

### Future Improvements
- Add comprehensive unit tests
- Implement remaining platforms (when APIs available)
- Create CI/CD pipeline
- Add data export features

---

Remember: **Platform separation is key!** When in doubt, keep platforms isolated.
