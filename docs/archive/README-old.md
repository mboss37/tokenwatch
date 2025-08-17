# TokenWatch CLI - Documentation

## 📚 **Documentation Overview**

This directory contains comprehensive documentation for TokenWatch CLI, designed to help developers understand the architecture, contribute to the project, and maintain consistency across all development work.

## 🏗️ **Architecture Documentation**

### **Core Architecture**
- **[README.md](architecture/README.md)** - High-level architecture overview and principles
- **[component-architecture.md](architecture/component-architecture.md)** - Detailed component relationships and interactions
- **[execution-flow.md](architecture/execution-flow.md)** - How data flows through the system for different commands
- **[platform-separation.md](architecture/platform-separation.md)** - Why and how platforms are separated (CRITICAL READING)

### **Architecture Decision Records (ADRs)**
- **[001-platform-separation.md](architecture/decisions/001-platform-separation.md)** - Platform separation architecture decision

## 🚀 **Quick Start for New Developers**

### **1. Start Here (Required Reading)**
1. **[Architecture Overview](architecture/README.md)** - Understand the big picture
2. **[Platform Separation](architecture/platform-separation.md)** - Learn the core principle
3. **[Component Architecture](architecture/component-architecture.md)** - See how everything fits together

### **2. Understand the Flows**
- **[Execution Flow](architecture/execution-flow.md)** - See how data moves through the system

### **3. Review Decisions**
- **[ADRs](architecture/decisions/)** - Understand why architectural choices were made

## 🔑 **Core Principles (Never Violate)**

### **1. 🚫 NO Platform Mixing**
- Each platform has its own files
- Never mix OpenAI API specs with Grok API specs
- Each provider knows only its own platform

### **2. ✅ Clean Separation**
- Command layer handles user interaction
- Provider layer implements platform-specific APIs
- Model layer provides common data structures
- Config layer manages centralized configuration

### **3. 🔄 Common Interface**
- All platforms implement the same Provider interface
- All platforms return data in the same format
- Commands work with any platform without knowing API details

### **4. 🚀 Easy Extension**
- Adding a new platform requires only 3 files
- No modification of existing platform code
- Clear extension points and patterns

## 📁 **Project Structure Reference**

```
tokenwatchcli/
├── cmd/root/                    # Command implementations
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
│   └── utils/                   # Utility functions
├── internal/                     # Internal packages
│   ├── config/                  # Configuration management
│   ├── api/                     # API client implementations
│   └── storage/                 # Local data storage
└── docs/                        # This documentation
```

## 🚨 **Common Mistakes to Avoid**

### **❌ Don't Do This**
1. **Mix platforms in one file**: `openai.go` should NOT contain Grok logic
2. **Cross-platform imports**: OpenAI shouldn't import Grok packages
3. **Shared platform logic**: Don't create functions that know about multiple platforms
4. **Mixed configuration**: Don't combine platform settings in one config struct

### **✅ Do This Instead**
1. **One platform per file**: Each file has one responsibility
2. **Use interfaces**: Common contracts, separate implementations
3. **Factory pattern**: Centralized provider creation
4. **Separate configs**: Each platform has its own configuration section

---

**This documentation is your guide to understanding and contributing to TokenWatch CLI. Read it thoroughly, follow the principles, and maintain the clean architecture that makes this project maintainable and extensible.**

**Remember: Platform separation is the foundation. Never violate it.**
