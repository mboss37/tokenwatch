# TokenWatch CLI - Developer Guide 🛠️

A comprehensive guide for developers contributing to TokenWatch CLI.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Getting Started](#getting-started)
3. [Code Structure](#code-structure)
4. [Adding New Platforms](#adding-new-platforms)
5. [Debug Mode Implementation](#debug-mode-implementation)
6. [Testing](#testing)
7. [Contributing Guidelines](#contributing-guidelines)

## Architecture Overview

TokenWatch CLI follows a clean, layered architecture with clear separation of concerns:

### Core Principles

- **Platform Separation**: Each AI platform is completely isolated
- **Interface-Driven Design**: Common interfaces for all platforms
- **Layered Architecture**: Command → Provider → Model → Config layers
- **Resilience Patterns**: Retry logic, rate limiting, circuit breaker
- **Clean Extension Points**: Easy to add new platforms

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    Command Layer                            │
│  (openai_cmd.go, overview_cmd.go, setup_cmd.go, etc.)     │
├─────────────────────────────────────────────────────────────┤
│                   Provider Layer                            │
│  (Provider interface, OpenAIProvider, etc.)                │
├─────────────────────────────────────────────────────────────┤
│                    Model Layer                              │
│  (Consumption, Pricing, ConsumptionSummary, etc.)          │
├─────────────────────────────────────────────────────────────┤
│                   Config Layer                              │
│  (Viper-based configuration management)                    │
├─────────────────────────────────────────────────────────────┤
│                   Utility Layer                            │
│  (HTTP client, circuit breaker, logging, etc.)            │
└─────────────────────────────────────────────────────────────┘
```

## Getting Started

### Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For version control
- **Make**: For build automation (optional)

### Development Setup

```bash
# Clone the repository
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch

# Install dependencies
go mod download

# Build the binary
go build -o tokenwatch ./cmd/root

# Run tests
go test ./...

# Run the CLI
./tokenwatch --help
```

### Development Workflow

```bash
# Make changes to code
# Build and test
go build ./cmd/root

# Run tests
go test ./...

# Test the CLI
./tokenwatch openai --debug

# Commit changes
git add .
git commit -m "Description of changes"
git push origin master
```

## Code Structure

### Directory Layout

```
tokenwatch/
├── cmd/root/                 # Command implementations
│   ├── main.go              # Entry point
│   ├── openai_cmd.go        # OpenAI command
│   ├── overview_cmd.go      # All platforms command
│   ├── setup_cmd.go         # Setup command
│   ├── config_cmd.go        # Configuration management
│   └── version_cmd.go       # Version command
├── pkg/                     # Public packages
│   ├── models/              # Data models
│   ├── providers/           # Platform providers
│   └── utils/               # Utility functions
├── internal/                # Internal packages
│   └── config/              # Configuration management
├── docs/                    # Documentation
└── configs/                 # Configuration examples
```

### Key Components

#### Command Layer (`cmd/root/`)

Commands implement the Cobra CLI framework:

```go
var openaiCmd = &cobra.Command{
    Use:   "openai",
    Short: "Show OpenAI token consumption and costs",
    Long:  `Display comprehensive OpenAI usage...`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Command implementation
    },
}
```

**Key Features:**
- **Flag Management**: `--period`, `--watch`, `--debug`
- **Error Handling**: Structured error responses
- **User Experience**: Clear output and helpful messages

#### Provider Layer (`pkg/providers/`)

Providers implement the `Provider` interface:

```go
type Provider interface {
    GetPlatform() string
    GetConsumption(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Consumption, error)
    GetPricing(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Pricing, error)
    GetConsumptionSummary(period string) (*models.ConsumptionSummary, error)
    GetPricingSummary(period string) (*models.PricingSummary, error)
    IsAvailable() bool
}
```

**Key Features:**
- **Cache Management**: 5-minute TTL with bypass capability
- **Rate Limiting**: 1 req/sec with burst of 5
- **Circuit Breaker**: Prevents cascading failures
- **Debug Mode**: Optional API request/response logging

#### Model Layer (`pkg/models/`)

Data models for consistent representation:

```go
type Consumption struct {
    Platform      string
    Model         string
    InputTokens   int64
    OutputTokens  int64
    TotalTokens   int64
    RequestCount  int64
    StartTime     time.Time
    EndTime       time.Time
}

type Pricing struct {
    Platform string
    Model    string
    LineItem string
    Amount   float64
    Currency string
    StartTime time.Time
    EndTime   time.Time
}
```

## Adding New Platforms

### Step 1: Create Provider Implementation

Create a new file `pkg/providers/[platform].go`:

```go
package providers

import (
    "time"
    "tokenwatch/pkg/models"
)

type [Platform]Provider struct {
    // Platform-specific fields
}

func New[Platform]Provider(apiKey string) *[Platform]Provider {
    return &[Platform]Provider{
        // Initialize fields
    }
}

// Implement Provider interface methods
func (p *[Platform]Provider) GetPlatform() string {
    return "[platform]"
}

func (p *[Platform]Provider) GetConsumption(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Consumption, error) {
    // Implementation
}

func (p *[Platform]Provider) GetPricing(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Pricing, error) {
    // Implementation
}

// ... other methods
```

### Step 2: Update Provider Factory

Add to `cmd/root/overview_cmd.go`:

```go
func getProvider(platform string) providers.Provider {
    switch platform {
    case "[platform]":
        apiKey := config.GetAPIKey("[platform]")
        if apiKey == "" {
            return nil
        }
        return providers.New[Platform]Provider(apiKey)
    // ... other cases
    }
}
```

### Step 3: Add Configuration Support

Update `internal/config/config.go`:

```go
func GetAPIKey(platform string) string {
    switch platform {
    case "[platform]":
        return viper.GetString("api_keys.[platform]")
    // ... other cases
    }
}
```

### Step 4: Update Setup Command

Add to `cmd/root/setup_cmd.go`:

```go
func setup[Platform](cmd *cobra.Command, args []string) error {
    // Platform-specific setup logic
}
```

## Debug Mode Implementation

### Overview

Debug mode provides optional API request/response logging for troubleshooting and development:

```bash
# Enable debug mode
./tokenwatch openai --debug

# Debug with watch mode
./tokenwatch openai -w --debug

# Debug all platforms
./tokenwatch all --debug
```

### Implementation Details

#### Command Layer

```go
// Add debug flag
openaiCmd.Flags().BoolP("debug", "d", false, "Enable debug logging for API calls")

// Pass debug parameter to display function
func displayOpenAIData(provider *providers.OpenAIProvider, period string, bypassCache bool, debug bool) error {
    // Pass debug to provider calls
    consumptions, err := provider.GetConsumption(startTime, endTime, bypassCache, debug)
    pricings, err := provider.GetPricing(startTime, endTime, bypassCache, debug)
}
```

#### Provider Layer

```go
// Update interface to include debug parameter
type Provider interface {
    GetConsumption(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Consumption, error)
    GetPricing(startTime, endTime time.Time, bypassCache bool, debug bool) ([]*models.Pricing, error)
}

// Implement debug logging in provider methods
func (o *OpenAIProvider) GetUsage(startTime, endTime time.Time, bucketWidth string, groupBy []string, bypassCache bool, debug bool) (*OpenAIUsageResponse, error) {
    // Log request details when debug is enabled
    if debug {
        fmt.Printf("🔍 OPENAI USAGE API REQUEST:\n")
        fmt.Printf("   URL: %s\n", req.URL.String())
        // ... more details
    }
    
    // ... API call logic
    
    // Log raw response when debug is enabled
    if debug {
        rawJSON, _ := json.MarshalIndent(usageResp, "", "  ")
        fmt.Printf("🔍 RAW OPENAI USAGE API RESPONSE:\n%s\n\n", string(rawJSON))
    }
}
```

### Debug Output Features

- **API Request Details**: URL, timestamps, parameters
- **Raw JSON Responses**: Complete API responses
- **Cache Behavior**: Shows when cache is hit/bypassed
- **Request/Response Flow**: Full API call lifecycle

### Use Cases

- **Troubleshooting**: Debug API issues and errors
- **Development**: Verify API behavior during development
- **Data Verification**: Check raw data for accuracy
- **Performance Analysis**: Monitor API call patterns

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/models
go test ./pkg/providers
go test ./internal/config

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Test Structure

Tests follow Go conventions:

```go
// pkg/models/consumption_test.go
func TestConsumption_AddConsumption(t *testing.T) {
    // Test implementation
}

func TestConsumptionSummary_AddConsumption(t *testing.T) {
    // Test implementation
}
```

### Test Coverage

Current test coverage includes:
- ✅ **Models**: Consumption, Pricing, Summaries
- ✅ **Providers**: OpenAI provider methods
- ✅ **Config**: Configuration management
- ✅ **Utils**: Utility functions

## Contributing Guidelines

### Code Style

- **Go Format**: Use `gofmt` or `go fmt`
- **Naming**: Follow Go naming conventions
- **Comments**: Document exported functions and types
- **Error Handling**: Use structured errors with context

### Commit Messages

Follow conventional commit format:

```
type(scope): description

- Detailed change 1
- Detailed change 2

Fixes #123
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test additions/changes

### Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch
3. **Make** your changes
4. **Test** thoroughly
5. **Commit** with clear messages
6. **Push** to your fork
7. **Create** a pull request

### Review Process

- **Code Review**: All PRs require review
- **Testing**: Ensure tests pass
- **Documentation**: Update docs if needed
- **Architecture**: Follow established patterns

## Key Features Implemented

### ✅ **Core Infrastructure**
- Provider interface, common models, config management
- Platform separation architecture
- Command structure with Cobra

### ✅ **Production Features**
- Retry logic, rate limiting, circuit breaker
- Response caching with TTL
- Structured error handling
- API key validation

### ✅ **Watch Mode**
- Real-time monitoring with `-w` flag
- Auto-refresh every 30 seconds
- Cache bypass for fresh data
- Screen clearing for clean display

### ✅ **Debug Mode**
- Optional API request/response logging
- `--debug` flag for troubleshooting
- Clean normal mode output
- Development-friendly debugging

### ✅ **Enhanced UX**
- Beautiful terminal tables with colors
- Detailed model breakdowns
- Total rows for aggregation
- Clear error messages with suggestions

### ✅ **Resilience Patterns**
- Exponential backoff retry logic
- Rate limiting (1 req/sec, burst 5)
- Circuit breaker (5 failures, 1 min reset)
- Graceful error handling

## Future Enhancements

### Planned Features

- **Additional Platforms**: Anthropic, Grok, Cursor (when APIs available)
- **Advanced Analytics**: Usage trends, cost predictions
- **Export Options**: CSV, JSON, PDF reports
- **Web Dashboard**: Browser-based monitoring
- **Alerting**: Cost threshold notifications

### Architecture Improvements

- **Plugin System**: Dynamic platform loading
- **Metrics Collection**: Prometheus integration
- **Distributed Caching**: Redis support
- **API Gateway**: Centralized API management

## Getting Help

### Resources

- **User Guide**: [docs/README.md](README.md)
- **GitHub Issues**: [Report bugs or request features](https://github.com/mboss37/tokenwatch/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mboss37/tokenwatch/discussions)

### Contact

- **Repository**: [https://github.com/mboss37/tokenwatch](https://github.com/mboss37/tokenwatch)
- **Issues**: [https://github.com/mboss37/tokenwatch/issues](https://github.com/mboss37/tokenwatch/issues)

---

**Happy coding! 🚀**
