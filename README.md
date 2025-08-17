# TokenWatch CLI 🚀

A simple and focused CLI tool for monitoring OpenAI token usage and costs with real-time capabilities.

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-OpenAI%20Only-orange.svg)](https://openai.com)

## ✨ Features

- **🤖 OpenAI Monitoring**: Track token consumption and costs across all models
- **⏰ Real-time Watch Mode**: Monitor usage with 30-second auto-refresh
- **📊 Smart Periods**: 1-day (recent), 7-day (historical), 30-day (extended)
- **🔍 Debug Mode**: Detailed API request/response logging for troubleshooting
- **💾 Intelligent Caching**: 5-minute cache with watch mode bypass
- **🔄 Resilience**: Rate limiting, circuit breaker, and automatic retries
- **🎨 Beautiful Output**: Colorful terminal tables with model breakdowns

## 🚀 Quick Start

```bash
# Clone and build
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch
go build -o tokenwatch ./cmd/root

# Setup your OpenAI API key
./tokenwatch setup

# Start monitoring
./tokenwatch usage                    # Last 7 days
./tokenwatch usage --period 1d       # Last 24 hours
./tokenwatch usage -w -p 1d          # Watch mode for real-time updates
```

## 📖 What It Does

TokenWatch CLI provides comprehensive monitoring of your OpenAI API usage:

- **Token Consumption**: Input, output, and total tokens per model
- **Cost Analysis**: Real-time cost tracking and daily averages
- **Model Breakdown**: Detailed usage statistics for each model
- **Smart Recommendations**: Period-specific guidance for optimal monitoring
- **Real-time Updates**: Watch mode for continuous monitoring

## 🛠️ Installation

### Build from Source (Recommended)

```bash
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch
go build -o tokenwatch ./cmd/root
chmod +x tokenwatch
./tokenwatch --help
```

### Install to System PATH

```bash
# Build and install
go build -o tokenwatch ./cmd/root
cp tokenwatch ~/go/bin/

# Add to PATH (add to ~/.bashrc, ~/.zshrc, or ~/.profile)
export PATH="$HOME/go/bin:$PATH"
```

### Using Makefile

```bash
make build      # Build only
make build-all  # Build for all platforms
make install    # Install locally
```

## 🔧 Setup

```bash
./tokenwatch setup
```

**Requirements:**
- OpenAI Admin API key with `api.usage.read` scope
- Organization-level access (personal API keys won't work)

The setup will:
- Prompt for your OpenAI Admin API key
- Validate the key with a test API call
- Save configuration to `~/.tokenwatch/config.yaml`

## 📚 Usage

### Basic Commands

```bash
# View OpenAI usage (last 7 days)
./tokenwatch usage

# Specify time period
./tokenwatch usage --period 1d       # Last 24 hours
./tokenwatch usage --period 7d       # Last 7 days (default)
./tokenwatch usage --period 30d      # Last 30 days

# Short flags
./tokenwatch usage -p 1d             # Same as --period 1d
```

### Watch Mode

Real-time monitoring with automatic refresh every 30 seconds:

```bash
# Watch mode (1-day period only)
./tokenwatch usage -w -p 1d

# Stop watching: Press Ctrl+C
```

**Features:**
- Auto-refresh every 30 seconds
- Cache bypass for fresh data
- Screen clearing for clean display
- Only available for 1-day periods (most logical)

### Debug Mode

Detailed API logging for troubleshooting:

```bash
# Enable debug mode
./tokenwatch usage --debug

# Debug with specific period
./tokenwatch usage --debug --period 30d
```

**Debug Output:**
- API request details (URL, timestamps, parameters)
- Raw JSON responses from OpenAI
- Pagination flow across multiple API calls
- Request/response lifecycle

### Configuration Management

```bash
# Check configuration status
./tokenwatch config check

# Reset configuration to defaults
./tokenwatch config reset

# View version
./tokenwatch version
```

## 📊 Example Output

### Normal Mode

```
🤖 OPENAI USAGE - Last 7d
⏰ Generated: 2025-08-17 13:30:13

📊 SUMMARY
──────────────────────────────────────────────────
📅 Period: 2025-08-10 to 2025-08-17 (7 days)
📈 Daily Averages: 721.7 tokens, 1.9 requests
💰 Daily Cost Average: $0.0051

💡 SMART RECOMMENDATIONS
──────────────────────────────────────────────────
📊 7-day period is ideal for:
   • Weekly usage patterns
   • Historical cost analysis
   • Model performance comparison
   • Budget planning

🔄 For recent activity, try: --period 1d

📋 MODEL BREAKDOWN
┌────────────────────────┬──────────────┬───────────────┬──────────────┬──────────┬─────────┬───────────────┐
│         MODEL          │ INPUT TOKENS │ OUTPUT TOKENS │ TOTAL TOKENS │ REQUESTS │  COST   │ $/ 1 K TOKENS │
├────────────────────────┼──────────────┼───────────────┬──────────────┬──────────┬─────────┬───────────────┤
│ gpt-4o-2024-08-06      │ 1872         │ 3099          │ 4971         │ 12       │ $0.0357 │ $0.0072       │
│ gpt-4o-mini-2024-07-18 │ 10           │ 71            │ 81           │ 1        │ $0.0000 │ $0.0005       │
│ ─                      │ ─            │ ─             │ ─            │ ─        │ ─       │ ─             │
│ TOTAL                  │ 1882         │ 3170          │ 5052         │ 13       │ $0.0357 │ $0.0071       │
└────────────────────────┴──────────────┴───────────────┴──────────────┴──────────┴─────────┴───────────────┘
```

### Debug Mode

```
🔍 OPENAI USAGE API REQUEST (Page 1, Token: ):
   URL: https://api.openai.com/v1/organization/usage/completions?...
   Start Time: 2025-08-10 13:30:26 (1754825426)
   End Time: 2025-08-17 13:30:26 (1755430226)
   Bucket Width: 1d
   Group By: [model]

🔍 RAW OPENAI USAGE API RESPONSE (Page 1, Token: ):
   Has More: true
   Next Page: page_AAAAAGijH7QR2l2hAAAAAGihG4A=
   Data Buckets: 7
   Total Results: 2
{
  "data": [
    {
      "start_time": 1754825426,
      "end_time": 1754870400,
      "results": []
    }
    // ... more data
  ]
}
```

## 🔍 Troubleshooting

### Common Issues

**"OpenAI not configured"**
```bash
./tokenwatch setup
```

**"API key lacks required permissions"**
- You need an Admin API key with `api.usage.read` scope
- Check your OpenAI organization settings

**"No data found"**
- Try a shorter period: `./tokenwatch usage --period 7d`
- Verify you have recent API usage
- Check if your API key is valid

### Debug Mode for Troubleshooting

```bash
# Enable debug mode to see API details
./tokenwatch usage --debug

# Check configuration
./tokenwatch config check

# Verify API key
./tokenwatch setup
```

## 🏗️ Architecture

TokenWatch follows a clean, layered architecture:

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
│  │  consumption.go │  │   pricing.go    │                  │ │
│  │ (usage models)  │  │ (cost models)   │                  │ │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Single Responsibility**: Each component has one clear purpose
2. **Dependency Injection**: Providers are injected into commands
3. **Interface Segregation**: Clean interfaces for each layer
4. **Error Handling**: Consistent error handling throughout
5. **Configuration**: Centralized configuration management

## 🚀 Advanced Features

### Cache Management

- **Normal mode**: 5-minute cache for efficiency
- **Watch mode**: Cache bypassed for real-time data
- **Debug mode**: Shows cache behavior

### Rate Limiting

- **OpenAI**: 1 request/second with burst of 5
- **Automatic retries** with exponential backoff
- **Circuit breaker** to prevent cascading failures

### Smart Recommendations

The CLI provides intelligent recommendations based on your selected time period:
- **1-day**: Perfect for recent activity monitoring
- **7-day**: Ideal for historical analysis and weekly patterns
- **30-day**: May have limited data due to API limitations

## 📁 Project Structure

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
│   ├── providers/           # Platform implementations
│   └── utils/               # Utility functions
├── internal/                 # Internal packages
│   └── config/              # Configuration management
├── docs/                    # Documentation
├── Makefile                 # Build and development tasks
└── README.md                # This file
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) and [Developer Guide](docs/DEVELOPER.md) for details.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch
go mod tidy

# Build and test
go build -o tokenwatch ./cmd/root
go test ./...

# Install locally
go install ./cmd/root
```

### Development Commands

```bash
make build      # Build the application
make test       # Run tests
make fmt        # Format code
make lint       # Run linter
make install    # Install locally
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [Viper](https://github.com/spf13/viper) for configuration
- Terminal tables powered by [tablewriter](https://github.com/olekukonko/tablewriter)
- Colors provided by [fatih/color](https://github.com/fatih/color)

## 📞 Support

- **Documentation**: [User Guide](docs/README.md) | [Developer Guide](docs/DEVELOPER.md)
- **Issues**: [GitHub Issues](https://github.com/mboss37/tokenwatch/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mboss37/tokenwatch/discussions)
- **Repository**: [https://github.com/mboss37/tokenwatch](https://github.com/mboss37/tokenwatch)

---

**Happy monitoring! 🚀**