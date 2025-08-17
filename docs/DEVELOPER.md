# TokenWatch Developer Guide 🛠️

Developer documentation for TokenWatch CLI - a tool for monitoring OpenAI token usage and costs.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Project Structure](#project-structure)
3. [Development Setup](#development-setup)
4. [Adding New Features](#adding-new-features)
5. [Testing](#testing)
6. [Building and Deployment](#building-and-deployment)

## Architecture Overview

TokenWatch follows a clean, layered architecture designed for simplicity and maintainability:

```
┌─────────────────────────────────────────────────────────────┐
│                    Command Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   usage_cmd.go  │  │   setup_cmd.go  │  │ config_cmd  │ │
│  │  (OpenAI usage) │  │ (API key setup) │  │ (settings)  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Provider Layer                            │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                openai.go                               │ │
│  │           (OpenAI API implementation)                  │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                    Model Layer                              │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  consumption.go │  │   pricing.go    │                  │
│  │ (usage models)  │  │ (cost models)   │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────────┐
│                   Utility Layer                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │  http_client.go │  │ circuit_breaker │  │   logger.go │ │
│  │ (rate limiting) │  │ (fault tolerance)│  │ (logging)  │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Single Responsibility**: Each component has one clear purpose
2. **Dependency Injection**: Providers are injected into commands
3. **Interface Segregation**: Clean interfaces for each layer
4. **Error Handling**: Consistent error handling throughout
5. **Configuration**: Centralized configuration management

## Project Structure

```
tokenwatch/
├── cmd/root/                 # Command implementations
│   ├── main.go              # Application entry point
│   ├── usage_cmd.go         # OpenAI usage command
│   ├── setup_cmd.go         # API key setup
│   ├── config_cmd.go        # Configuration management
│   └── version_cmd.go       # Version information
├── pkg/                     # Reusable packages
│   ├── models/              # Data models
│   │   ├── consumption.go   # Usage data structures
│   │   └── pricing.go       # Cost data structures
│   ├── providers/           # Platform implementations
│   │   ├── provider.go      # Common interface
│   │   └── openai.go        # OpenAI API implementation
│   └── utils/               # Utility functions
│       ├── http_client.go   # Rate-limited HTTP client
│       ├── circuit_breaker.go # Fault tolerance
│       ├── logger.go        # Structured logging
│       ├── prompt.go        # User input handling
│       ├── validation.go    # Input validation
│       └── errors.go        # Error definitions
├── internal/                 # Internal packages
│   └── config/              # Configuration management
│       └── config.go        # Viper-based config
├── docs/                    # Documentation
│   ├── README.md            # User guide
│   └── DEVELOPER.md         # This file
├── Makefile                 # Build and development tasks
├── go.mod                   # Go module definition
└── README.md                # Project overview
```

## Development Setup

### Prerequisites

- Go 1.21+
- OpenAI Admin API key with `api.usage.read` scope

### Local Development

```bash
# Clone the repository
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch

# Install dependencies
go mod tidy

# Build the application
go build -o tokenwatch ./cmd/root

# Run tests
go test ./...

# Install locally
go install ./cmd/root
```

### Development Commands

```bash
# Build
make build

# Build for all platforms
make build-all

# Install locally
make install

# Format code
make fmt

# Run linter
make lint

# Run tests
make test
```

## Adding New Features

### Adding a New Command

1. **Create the command file** in `cmd/root/`
2. **Implement the command logic** following existing patterns
3. **Register the command** in the main command structure
4. **Add tests** for the new functionality

Example:
```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Description of new command",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation here
        return nil
    },
}

func init() {
    RootCmd.AddCommand(newCmd)
}
```

### Adding a New Provider

1. **Implement the Provider interface** in `pkg/providers/`
2. **Add configuration support** for the new platform
3. **Update setup command** to handle the new platform
4. **Add validation** for the new API key format

### Adding New Models

1. **Define the data structure** in `pkg/models/`
2. **Add conversion functions** from API responses
3. **Update display logic** to handle the new data

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/utils

# Run tests with verbose output
go test -v ./...
```

### Test Structure

- **Unit tests** for individual functions
- **Integration tests** for command workflows
- **Mock external APIs** for reliable testing

## Building and Deployment

### Build Targets

```bash
# Single platform
make build

# All platforms
make build-all

# Install to system
make install
```

### Supported Platforms

- **Linux**: amd64, arm64
- **macOS**: amd64, arm64  
- **Windows**: amd64

### Release Process

1. **Update version** in `main.go`
2. **Build all platforms** with `make build-all`
3. **Create GitHub release** with built binaries
4. **Tag the release** in git

## Configuration

### Configuration File

Located at `~/.tokenwatch/config.yaml`:

```yaml
api_keys:
  openai: "sk-..."

settings:
  debug: false
  cache_duration: 300
  request_timeout: 10
  retry_attempts: 3
```

### Environment Variables

```bash
# API Keys
export OPENAI_API_KEY="sk-..."

# Logging
export TOKENWATCH_LOG_LEVEL="debug"
```

## Error Handling

### Error Types

- **ValidationError**: Invalid input or configuration
- **APIError**: External API communication issues
- **ConfigError**: Configuration problems
- **InternalError**: Unexpected internal issues

### Error Recovery

- **Automatic retries** with exponential backoff
- **Circuit breaker** to prevent cascading failures
- **Graceful degradation** when possible
- **Clear error messages** for users

## Logging

### Log Levels

- **INFO**: General application flow
- **DEBUG**: Detailed debugging information
- **WARN**: Warning conditions
- **ERROR**: Error conditions

### Structured Logging

All logs include structured data for better debugging:

```go
utils.Info("API request completed", map[string]interface{}{
    "platform": "openai",
    "duration": "1.2s",
    "status": "success",
})
```

## Performance Considerations

### Caching

- **5-minute TTL** for normal operations
- **Cache bypass** in watch mode
- **Smart cache invalidation** based on usage patterns

### Rate Limiting

- **1 request/second** with burst of 5
- **Automatic backoff** on rate limit errors
- **Circuit breaker** for fault tolerance

### Memory Management

- **Efficient data structures** for large datasets
- **Streaming responses** where possible
- **Memory cleanup** after operations

## Security

### API Key Management

- **Secure storage** in user's home directory
- **Environment variable support** for CI/CD
- **No hardcoded keys** in source code
- **Key validation** before use

### Data Privacy

- **Local processing** of sensitive data
- **No external logging** of API responses
- **Configurable debug output** for development

## Contributing

### Code Style

- **Go fmt** for formatting
- **golangci-lint** for linting
- **Consistent naming** conventions
- **Documentation** for public APIs

### Pull Request Process

1. **Fork the repository**
2. **Create a feature branch**
3. **Implement your changes**
4. **Add tests** for new functionality
5. **Update documentation**
6. **Submit a pull request**

### Testing Requirements

- **All new code** must have tests
- **Existing tests** must pass
- **Coverage** should not decrease
- **Integration tests** for new commands

## Troubleshooting

### Common Development Issues

**Build failures**
```bash
go mod tidy
go clean -cache
```

**Test failures**
```bash
go test -v ./...
go vet ./...
```

**Linting issues**
```bash
golangci-lint run
make fmt
```

### Debug Mode

Enable debug logging for development:

```bash
export TOKENWATCH_LOG_LEVEL=debug
./tokenwatch usage --debug
```

## Future Enhancements

### Planned Features

- **Additional platforms** when APIs become available
- **Export functionality** (CSV, JSON)
- **Web dashboard** for visualization
- **Alerting system** for cost thresholds

### Architecture Evolution

- **Plugin system** for platform support
- **Database backend** for historical data
- **API server** for remote access
- **Metrics collection** for monitoring

## Getting Help

- **GitHub Issues**: Report bugs and request features
- **GitHub Discussions**: Ask questions and share ideas
- **Documentation**: Check this guide and [User Guide](USER_GUIDE.md)
- **Code**: Review the source code for examples

For more information, visit the [GitHub repository](https://github.com/mboss37/tokenwatch).
