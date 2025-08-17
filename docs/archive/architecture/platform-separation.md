# Platform Separation - Why and How

## 🎯 **Why Platform Separation?**

This document explains the fundamental architectural decision to keep each LLM platform completely isolated in separate files. This is a **core principle** of TokenWatch CLI that must never be violated.

## 🚫 **The Problem We're Solving**

### **Without Platform Separation**
```
❌ BAD: Mixed platform logic in one file
├── openai_cmd.go
│   ├── OpenAI API calls
│   ├── Grok API calls          ← Mixed platforms!
│   ├── Anthropic API calls     ← Mixed platforms!
│   └── Cursor API calls        ← Mixed platforms!
```

**Problems:**
- **API Spec Confusion**: Mixing different API endpoints and formats
- **Authentication Chaos**: Different platforms need different auth methods
- **Rate Limiting Issues**: Each platform has different rate limits
- **Debugging Nightmare**: Hard to isolate which platform has issues
- **Maintenance Hell**: Changes to one platform affect others
- **Testing Complexity**: Can't test platforms independently

### **With Platform Separation**
```
✅ GOOD: Clean platform separation
├── openai_cmd.go (OpenAI only)
├── grok_cmd.go (Grok only)
├── anthropic_cmd.go (Anthropic only)
└── cursor_cmd.go (Cursor only)
```

**Benefits:**
- **Clear Boundaries**: Each file has one responsibility
- **Easy Debugging**: Issues isolated to specific platform
- **Independent Development**: Work on platforms separately
- **Simple Testing**: Test each platform in isolation
- **Easy Maintenance**: Changes don't affect other platforms

## 🏗️ **How Platform Separation Works**

### **1. File-Level Separation**

```
cmd/root/
├── all_cmd.go              ← Multi-platform orchestration ONLY
├── openai_cmd.go           ← OpenAI commands ONLY
├── grok_cmd.go             ← Grok commands ONLY
├── anthropic_cmd.go        ← Anthropic commands ONLY
└── cursor_cmd.go           ← Cursor commands ONLY

pkg/providers/
├── provider.go             ← Common interface (NO platform logic)
├── openai.go               ← OpenAI API implementation ONLY
├── grok.go                 ← Grok API implementation ONLY
├── anthropic.go            ← Anthropic API implementation ONLY
└── cursor.go               ← Cursor API implementation ONLY
```

### **2. Interface-Based Unification**

```go
// Common interface that all platforms implement
type Provider interface {
    GetPlatform() string
    GetConsumptionSummary(period string) (*models.ConsumptionSummary, error)
    GetPricingSummary(period string) (*models.PricingSummary, error)
    IsAvailable() bool
}

// Each platform implements this interface independently
type OpenAIProvider struct { /* OpenAI-specific fields */ }
type GrokProvider struct { /* Grok-specific fields */ }
type AnthropicProvider struct { /* Anthropic-specific fields */ }
type CursorProvider struct { /* Cursor-specific fields */ }
```

### **3. Factory Pattern for Provider Creation**

```go
// Centralized provider creation - NO platform mixing
func getProvider(platform string) providers.Provider {
    switch platform {
    case "openai":
        return providers.NewOpenAIProvider(apiKey, "")
    case "grok":
        return providers.NewGrokProvider(apiKey, "")
    case "anthropic":
        return providers.NewAnthropicProvider(apiKey, "")
    case "cursor":
        return providers.NewCursorProvider(apiKey, "")
    default:
        return nil
    }
}
```

## 🔒 **Separation Rules - NEVER VIOLATE THESE**

### **Rule 1: One Platform Per File**
```
✅ CORRECT: openai.go contains ONLY OpenAI logic
❌ WRONG: openai.go contains OpenAI + Grok logic
```

### **Rule 2: No Cross-Platform Imports**
```
✅ CORRECT: openai.go imports only openai-specific packages
❌ WRONG: openai.go imports grok packages
```

### **Rule 3: No Shared Platform Logic**
```
✅ CORRECT: Each platform handles its own API calls
❌ WRONG: Common function that calls multiple platform APIs
```

### **Rule 4: No Mixed Configuration**
```
✅ CORRECT: Each platform has its own config section
❌ WRONG: Shared config that mixes platform settings
```

## 📋 **What Each Platform File Contains**

### **Command Files (e.g., `openai_cmd.go`)**
```go
// ✅ CORRECT: OpenAI-specific command logic
var openaiCmd = &cobra.Command{
    Use:   "openai",
    Short: "Show OpenAI token consumption and costs",
    RunE: func(cmd *cobra.Command, args []string) error {
        // ONLY OpenAI logic here
        provider := getProvider("openai")
        // ... OpenAI-specific processing
    },
}
```

**Contains:**
- OpenAI command definition
- OpenAI-specific flag handling
- OpenAI data processing
- OpenAI output formatting

**Does NOT contain:**
- Grok API calls
- Anthropic data processing
- Cursor-specific logic

### **Provider Files (e.g., `openai.go`)**
```go
// ✅ CORRECT: OpenAI API implementation only
type OpenAIProvider struct {
    apiKey string
    client *http.Client
}

func (p *OpenAIProvider) GetConsumptionSummary(period string) (*models.ConsumptionSummary, error) {
    // ONLY OpenAI API calls here
    resp, err := p.client.Get("https://api.openai.com/v1/usage")
    // ... OpenAI-specific response handling
}
```

**Contains:**
- OpenAI API endpoints
- OpenAI authentication
- OpenAI rate limiting
- OpenAI response parsing

**Does NOT contain:**
- Grok API URLs
- Anthropic response formats
- Cursor authentication methods

## 🔄 **How Multi-Platform Commands Work**

### **The `all` Command Orchestrates, Doesn't Mix**

```go
// ✅ CORRECT: Orchestration without mixing
func collectPlatformDataParallel(platforms []string, period string) []PlatformDataResult {
    var wg sync.WaitGroup
    results := make([]PlatformDataResult, 0, len(platforms))
    
    // Launch separate goroutine for each platform
    for _, platform := range platforms {
        wg.Add(1)
        go func(p string) {
            defer wg.Done()
            
            // Each platform runs independently
            provider := getProvider(p)  // Returns platform-specific provider
            // ... platform-specific API calls
        }(platform)
    }
    
    // Collect results without mixing platform logic
    // ... result collection
}
```

**Key Points:**
- **No platform mixing**: Each goroutine handles one platform
- **Independent execution**: Platforms don't interfere with each other
- **Clean aggregation**: Results combined without platform knowledge
- **Error isolation**: One platform's failure doesn't affect others

## 🚀 **Adding New Platforms - The Right Way**

### **Step 1: Create Platform-Specific Command**
```go
// cmd/root/anthropic_cmd.go
var anthropicCmd = &cobra.Command{
    Use:   "anthropic",
    Short: "Show Anthropic token consumption and costs",
    RunE: func(cmd *cobra.Command, args []string) error {
        // ONLY Anthropic logic here
        provider := getProvider("anthropic")
        // ... Anthropic-specific processing
    },
}
```

### **Step 2: Create Platform-Specific Provider**
```go
// pkg/providers/anthropic.go
type AnthropicProvider struct {
    apiKey string
    client *http.Client
}

func (p *AnthropicProvider) GetConsumptionSummary(period string) (*models.ConsumptionSummary, error) {
    // ONLY Anthropic API calls here
    resp, err := p.client.Get("https://api.anthropic.com/v1/usage")
    // ... Anthropic-specific response handling
}
```

### **Step 3: Update Factory Function**
```go
// In all_cmd.go - add to getProvider function
case "anthropic":
    return providers.NewAnthropicProvider(apiKey, "")
```

### **Step 4: Update Setup Command**
```go
// In setup_cmd.go - add to supported platforms
{"anthropic", "Anthropic", "Claude and other Anthropic models", true}
```

## 🧪 **Testing Platform Separation**

### **Unit Tests**
```go
// Test each platform independently
func TestOpenAIProvider(t *testing.T) {
    provider := NewOpenAIProvider("test-key", "")
    // Test ONLY OpenAI functionality
}

func TestGrokProvider(t *testing.T) {
    provider := NewGrokProvider("test-key", "")
    // Test ONLY Grok functionality
}
```

### **Integration Tests**
```go
// Test that platforms don't interfere
func TestPlatformIsolation(t *testing.T) {
    // Test that OpenAI commands don't affect Grok
    // Test that platform-specific configs are separate
    // Test that errors in one platform don't affect others
}
```

## 🚨 **Common Mistakes to Avoid**

### **Mistake 1: Shared Platform Logic**
```go
// ❌ WRONG: Don't do this
func getUsageData(platform string) (*models.ConsumptionSummary, error) {
    switch platform {
    case "openai":
        return callOpenAIAPI()
    case "grok":
        return callGrokAPI()
    }
}
```

**Why it's wrong:** This function knows about multiple platforms and their APIs.

### **Mistake 2: Cross-Platform Dependencies**
```go
// ❌ WRONG: Don't do this
// In openai.go
import "tokenwatch/pkg/providers/grok"  // Wrong! OpenAI shouldn't know about Grok
```

**Why it's wrong:** This creates tight coupling between platforms.

### **Mistake 3: Mixed Configuration**
```go
// ❌ WRONG: Don't do this
type Config struct {
    OpenAIAPIKey string
    GrokAPIKey   string  // Wrong! Mixing platform configs
    AnthropicAPIKey string
}
```

**Why it's wrong:** This creates a single point of failure and mixing.

## ✅ **Best Practices Summary**

1. **One platform per file** - Never mix platform logic
2. **Use interfaces** - Common contracts, separate implementations
3. **Factory pattern** - Centralized provider creation
4. **Independent testing** - Test each platform separately
5. **Clear boundaries** - Each file has one responsibility
6. **No cross-imports** - Platforms don't know about each other
7. **Separate configuration** - Each platform has its own config section

---

**Platform separation is the foundation of TokenWatch CLI's maintainability and extensibility. Every developer must understand and follow these principles to ensure the project remains clean and manageable.**
