# TokenWatch CLI - Architecture Overview

## 🎯 **Project Vision**

TokenWatch CLI is a fast and elegant command-line tool for tracking token consumption and pricing across multiple LLM platforms. The architecture is designed for **clean separation**, **easy maintenance**, and **simple extension** to new platforms.

## 🏗️ **High-Level Architecture**

```
┌─────────────────────────────────────────────────────────────────┐
│                    USER COMMAND                                │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ tokenwatch all  │  │ tokenwatch openai│  │ tokenwatch grok │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    COMMAND LAYER                               │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   all_cmd.go    │  │  openai_cmd.go  │  │   grok_cmd.go   │ │
│  │ (parallel exec) │  │ (OpenAI logic)  │  │ (Grok logic)    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    PROVIDER LAYER                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │  openai.go      │  │   grok.go       │  │  anthropic.go   │ │
│  │ (OpenAI API)    │  │ (Grok API)      │  │ (Anthropic API) │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    EXTERNAL APIs                               │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │ OpenAI API      │  │ Grok API        │  │ Anthropic API   │ │
│  │ (usage/costs)   │  │ (usage/costs)   │  │ (usage/costs)   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## 🔑 **Core Principles**

### **1. 🚫 NO Platform Mixing**
- **Each platform has its own file** with its own API logic
- **Never mix OpenAI API specs with Grok API specs**
- **Each provider knows only its own platform's requirements**

### **2. ✅ Clean Separation**
- **Command layer**: Handles user interaction and command logic
- **Provider layer**: Implements platform-specific API calls
- **Model layer**: Common data structures across all platforms
- **Config layer**: Centralized configuration management

### **3. 🔄 Common Interface**
- **All platforms implement the same Provider interface**
- **All platforms return data in the same format**
- **Commands can work with any platform without knowing API details**

### **4. 🚀 Easy Extension**
- **Adding a new platform requires only 3 files**:
  - `cmd/root/{platform}_cmd.go` (command logic)
  - `pkg/providers/{platform}.go` (API implementation)
  - Update `getProvider()` function in `all_cmd.go`

## 📁 **Project Structure**

```
tokenwatchcli/
├── cmd/root/                    # Command implementations
│   ├── main.go                  # Entry point
│   ├── all_cmd.go              # Multi-platform aggregation
│   ├── openai_cmd.go           # OpenAI-specific commands
│   ├── grok_cmd.go             # Grok-specific commands (future)
│   ├── anthropic_cmd.go        # Anthropic-specific commands (future)
│   ├── cursor_cmd.go            # Cursor-specific commands (future)
│   ├── config_cmd.go            # Configuration management
│   └── setup_cmd.go             # Interactive setup
├── pkg/                         # Public packages
│   ├── providers/               # Platform implementations
│   │   ├── provider.go          # Common interface
│   │   ├── openai.go            # OpenAI API implementation
│   │   ├── grok.go              # Grok API implementation (future)
│   │   ├── anthropic.go         # Anthropic API implementation (future)
│   │   └── cursor.go            # Cursor API implementation (future)
│   ├── models/                  # Common data structures
│   │   ├── consumption.go       # Consumption data models
│   │   └── pricing.go           # Pricing data models
│   └── utils/                   # Utility functions
├── internal/                     # Internal packages
│   ├── config/                  # Configuration management
│   ├── api/                     # API client implementations
│   └── storage/                 # Local data storage
├── configs/                      # Configuration templates
├── docs/                        # Documentation
│   ├── architecture/            # Architecture documentation
│   ├── api/                     # API documentation
│   └── user-guide/              # User documentation
└── go.mod                       # Go module file
```

## 🎯 **Key Design Decisions**

### **Why Parallel Execution in `all` Command?**
- **Performance**: Fetch data from all platforms simultaneously
- **User Experience**: Faster response times
- **Scalability**: Easy to add more platforms without performance degradation

### **Why Platform Separation?**
- **Maintainability**: Work on one platform without affecting others
- **Debugging**: Isolate issues to specific platform implementations
- **Team Development**: Multiple developers can work on different platforms
- **Testing**: Test each platform independently

### **Why Common Interface?**
- **Code Reuse**: Commands work with any platform
- **Consistency**: Same data format across all platforms
- **Extensibility**: Easy to add new platforms

## 🚀 **Development Workflow**

### **Adding a New Platform**
1. **Create command file**: `cmd/root/{platform}_cmd.go`
2. **Create provider**: `pkg/providers/{platform}.go`
3. **Implement Provider interface**: `GetConsumptionSummary()`, `GetPricingSummary()`
4. **Update setup command**: Add platform to supported platforms list
5. **Test**: Verify platform-specific commands work
6. **Verify integration**: Ensure `tokenwatch all` includes new platform

### **Modifying Existing Platform**
1. **Work only in the platform-specific files**
2. **Never modify other platform files**
3. **Test platform-specific commands**
4. **Verify integration still works**

## 📚 **Documentation Structure**

- **`README.md`** (this file): High-level architecture overview
- **`component-architecture.md`**: Detailed component relationships
- **`execution-flow.md`**: How data flows through the system
- **`platform-separation.md`**: Why and how platforms are separated
- **`decisions/`**: Architecture Decision Records (ADRs)

## 🔒 **Security & Configuration**

- **API keys stored locally** in `~/.tokenwatch/config.yaml`
- **Never committed to git** (protected by .gitignore)
- **Environment variable support** for CI/CD environments
- **Secure input handling** with masked prompts

## 📈 **Performance Considerations**

- **Parallel API calls** for multi-platform commands
- **Caching support** for frequently accessed data
- **Configurable timeouts** for API calls
- **Graceful degradation** when platforms are unavailable

---

**This architecture ensures that TokenWatch CLI is maintainable, extensible, and follows Go best practices. Every developer working on this project should understand and follow these principles.**
