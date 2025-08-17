# ADR-001: Platform Separation Architecture

## Status
Accepted

## Context
TokenWatch CLI needs to support multiple LLM platforms (OpenAI, Anthropic, Grok, Cursor) with different API specifications, authentication methods, rate limits, and response formats.

## Decision
Implement **complete platform separation** where each platform has its own dedicated files and implementations, with no mixing of platform-specific logic.

## Consequences

### Positive
- **Easy Maintenance**: Work on one platform without affecting others
- **Clear Debugging**: Issues isolated to specific platform implementations
- **Independent Development**: Multiple developers can work on different platforms simultaneously
- **Simple Testing**: Test each platform in isolation
- **Easy Extension**: Add new platforms without modifying existing code
- **No API Confusion**: Each platform handles its own API specifications
- **Clean Architecture**: Single responsibility principle for each file

### Negative
- **Code Duplication**: Some common patterns may be repeated across platforms
- **File Count**: More files in the project structure
- **Interface Complexity**: Need to design common interfaces that work for all platforms

### Neutral
- **Learning Curve**: New developers need to understand the separation principle
- **Documentation**: Need comprehensive documentation to explain the architecture

## Implementation Details

### File Structure
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

### Key Principles
1. **One Platform Per File**: Never mix platform logic in a single file
2. **No Cross-Platform Imports**: Platforms don't know about each other
3. **Common Interface**: All platforms implement the same Provider interface
4. **Factory Pattern**: Centralized provider creation without platform mixing
5. **Independent Testing**: Each platform can be tested separately

### Provider Interface
```go
type Provider interface {
    GetPlatform() string
    GetConsumptionSummary(period string) (*models.ConsumptionSummary, error)
    GetPricingSummary(period string) (*models.PricingSummary, error)
    IsAvailable() bool
}
```

## Alternatives Considered

### Alternative 1: Monolithic Provider
- **Description**: Single provider file handling all platforms
- **Rejected**: Would create massive, unmaintainable files with mixed API logic

### Alternative 2: Platform-Specific Packages
- **Description**: Separate packages for each platform
- **Rejected**: Over-engineering for current scope, adds unnecessary complexity

### Alternative 3: Configuration-Driven
- **Description**: Single provider with platform-specific configuration
- **Rejected**: Would still mix platform logic and create complex conditional code

## Related Decisions
- [ADR-002: Parallel Execution for Multi-Platform Commands](./002-parallel-execution.md)
- [ADR-003: Common Data Models Across Platforms](./003-common-data-models.md)

## References
- [Platform Separation Documentation](../platform-separation.md)
- [Component Architecture Documentation](../component-architecture.md)
- [Go Best Practices for Package Organization](https://golang.org/doc/effective_go.html#package-names)

---

**This decision establishes the foundation for TokenWatch CLI's maintainable and extensible architecture. Every developer working on this project must understand and follow these principles.**
