# Component Architecture - Detailed View

## 🏗️ **Component Relationships**

This document provides a detailed view of how all components in TokenWatch CLI interact with each other.

## 📊 **Component Diagram**

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              USER INTERFACE                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ tokenwatch all  │  │ tokenwatch openai│  │ tokenwatch grok │  │ tokenwatch  │ │
│  │ (multi-platform)│  │ (OpenAI only)   │  │ (Grok only)     │  │   setup     │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              COMMAND LAYER                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   all_cmd.go    │  │  openai_cmd.go  │  │   grok_cmd.go   │  │ setup_cmd.go│ │
│  │                 │  │                 │  │                 │  │             │ │
│  │ • Parallel exec │  │ • OpenAI logic  │  │ • Grok logic    │  │ • Platform  │ │
│  │ • Aggregation  │  │ • Period flags  │  │ • Period flags  │  │   selection │ │
│  │ • Multi-platform│  │ • Data display  │  │ • Data display  │  │ • API key   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  │   input     │ │
└─────────────────────────────────────────────────────────────────┘  └─────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              PROVIDER LAYER                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  openai.go      │  │   grok.go       │  │  anthropic.go   │  │  cursor.go  │ │
│  │                 │  │                 │  │                 │  │             │ │
│  │ • OpenAI API    │  │ • Grok API      │  │ • Anthropic API │  │ • Cursor API│ │
│  │ • Usage calls   │  │ • Usage calls   │  │ • Usage calls   │  │ • Usage calls│ │
│  │ • Cost calls    │  │ • Cost calls    │  │ • Cost calls    │  │ • Cost calls │ │
│  │ • Rate limiting │  │ • Rate limiting │  │ • Rate limiting │  │ • Rate limit │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
│           │                      │                      │                      │ │
│           └──────────────────────┼──────────────────────┼──────────────────────┘ │
│                                  ▼                      ▼                      │ │
│  ┌─────────────────────────────────────────────────────────────────────────────┐ │
│  │                           PROVIDER INTERFACE                               │ │
│  │  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │  │ GetPlatform() string                                                   │ │
│  │  │ GetConsumption(start, end time.Time) ([]*Consumption, error)          │ │
│  │  │ GetPricing(start, end time.Time) ([]*Pricing, error)                  │ │
│  │  │ GetConsumptionSummary(period string) (*ConsumptionSummary, error)     │ │
│  │  │ GetPricingSummary(period string) (*PricingSummary, error)             │ │
│  │  │ IsAvailable() bool                                                     │ │
│  │  └─────────────────────────────────────────────────────────────────────────┘ │
│  └─────────────────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              MODEL LAYER                                       │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ consumption.go  │  │   pricing.go    │  │   common.go     │  │   types.go  │ │
│  │                 │  │                 │  │                 │  │             │ │
│  │ • Consumption   │  │ • Pricing       │  │ • Common        │  │ • Platform  │ │
│  │ • Usage data    │  │ • Cost data     │  │ • Time periods  │  │   constants │ │
│  │ • Token counts  │  │ • Model rates   │  │ • Period calc   │  │ • API URLs  │ │
│  │ • Request info  │  │ • Total costs   │  │ • Date utils    │  │ • Endpoints │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              CONFIG LAYER                                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   config.go     │  │  viper.yaml     │  │  env_vars       │  │  defaults   │ │
│  │                 │  │                 │  │                 │  │             │ │
│  │ • Config init   │  │ • User settings │  │ • API keys      │  │ • Cache     │ │
│  │ • File loading  │  │ • Display opts  │  │ • Timeouts      │  │ • Timeouts  │ │
│  │ • API key mgmt  │  │ • Output format │  │ • Debug flags   │  │ • Retries   │ │
│  │ • Validation    │  │ • Color themes  │  │ • Log levels    │  │ • Formats   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              EXTERNAL APIs                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │ OpenAI API      │  │ Grok API        │  │ Anthropic API   │  │ Cursor API  │ │
│  │                 │  │                 │  │                 │  │             │ │
│  │ • Usage         │  │ • Usage         │  │ • Usage         │  │ • Usage     │ │
│  │ • Costs         │  │ • Costs         │  │ • Costs         │  │ • Costs     │ │
│  │ • Models        │  │ • Models        │  │ • Models        │  │ • Models    │ │
│  │ • Rate limits   │  │ • Rate limits   │  │ • Rate limits   │  │ • Rate limit│ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## 🔗 **Component Dependencies**

### **Command Layer Dependencies**
```
all_cmd.go
├── internal/config (for API keys)
├── pkg/providers (for platform providers)
├── pkg/models (for data structures)
└── sync (for parallel execution)

openai_cmd.go
├── internal/config (for API keys)
├── pkg/providers/openai (for OpenAI provider)
└── pkg/models (for data structures)

grok_cmd.go
├── internal/config (for API keys)
├── pkg/providers/grok (for Grok provider)
└── pkg/models (for data structures)
```

### **Provider Layer Dependencies**
```
openai.go
├── pkg/models (for data structures)
├── net/http (for API calls)
├── encoding/json (for JSON parsing)
└── time (for time handling)

grok.go
├── pkg/models (for data structures)
├── net/http (for API calls)
├── encoding/json (for JSON parsing)
└── time (for time handling)
```

### **Model Layer Dependencies**
```
consumption.go
├── time (for time handling)
└── encoding/json (for JSON tags)

pricing.go
├── time (for time handling)
└── encoding/json (for JSON tags)
```

## 🔄 **Data Flow Between Components**

### **1. Single Platform Command Flow**
```
User → openai_cmd.go → openai.go → OpenAI API → openai.go → openai_cmd.go → User
```

### **2. Multi-Platform Command Flow**
```
User → all_cmd.go → [openai.go, grok.go, anthropic.go] → [APIs] → [providers] → all_cmd.go → User
```

### **3. Configuration Flow**
```
User → setup_cmd.go → config.go → config.yaml → config.go → providers
```

## 📋 **Interface Contracts**

### **Provider Interface**
```go
type Provider interface {
    GetPlatform() string
    GetConsumption(startTime, endTime time.Time) ([]*models.Consumption, error)
    GetPricing(startTime, endTime time.Time) ([]*models.Pricing, error)
    GetConsumptionSummary(period string) (*models.ConsumptionSummary, error)
    GetPricingSummary(period string) (*models.PricingSummary, error)
    IsAvailable() bool
}
```

### **Data Model Contracts**
```go
type ConsumptionSummary struct {
    Platform      string
    Model         string
    TotalTokens   int64
    TotalRequests int64
    StartTime     time.Time
    EndTime       time.Time
}

type PricingSummary struct {
    Platform   string
    Model      string
    TotalCost  float64
    StartTime  time.Time
    EndTime    time.Time
}
```

## 🎯 **Key Design Patterns**

### **1. Strategy Pattern**
- **Each platform implements the same Provider interface**
- **Commands can work with any provider without knowing implementation details**
- **Easy to swap or add new platform implementations**

### **2. Factory Pattern**
- **`getProvider()` function creates appropriate provider instances**
- **Centralized provider creation logic**
- **Easy to add new platform types**

### **3. Observer Pattern**
- **Parallel execution with result collection via channels**
- **Non-blocking API calls with result aggregation**
- **Graceful handling of platform failures**

### **4. Template Method Pattern**
- **Common command structure across all platforms**
- **Platform-specific logic isolated in provider implementations**
- **Consistent user experience across all commands**

## 🔒 **Encapsulation & Boundaries**

### **Command Layer**
- **Public**: Command definitions and help text
- **Private**: Internal command logic and data processing
- **Boundary**: User input validation and output formatting

### **Provider Layer**
- **Public**: Provider interface implementation
- **Private**: API call logic, authentication, rate limiting
- **Boundary**: Data transformation between API format and common models

### **Model Layer**
- **Public**: Data structures and validation methods
- **Private**: Internal data processing logic
- **Boundary**: Data serialization/deserialization

### **Config Layer**
- **Public**: Configuration access methods
- **Private**: File I/O, validation, defaults
- **Boundary**: Configuration file format and environment variables

## 🚀 **Extension Points**

### **Adding New Platforms**
1. **Create new provider file**: `pkg/providers/{platform}.go`
2. **Implement Provider interface**: All required methods
3. **Create command file**: `cmd/root/{platform}_cmd.go`
4. **Update provider factory**: Add to `getProvider()` function
5. **Update setup command**: Add to supported platforms list

### **Adding New Commands**
1. **Create command file**: `cmd/root/{command}_cmd.go`
2. **Define command structure**: Use cobra.Command
3. **Implement command logic**: Use existing providers and models
4. **Add to root command**: Update main.go

### **Adding New Data Models**
1. **Create model file**: `pkg/models/{model}.go`
2. **Define data structure**: Use Go structs with JSON tags
3. **Add validation**: Implement validation methods
4. **Update providers**: Use new models in provider implementations

---

**This component architecture ensures clean separation, easy maintenance, and simple extension while maintaining consistent interfaces and data flow throughout the system.**
