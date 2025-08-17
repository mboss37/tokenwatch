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
./tokenwatch usage -w -p 7d          # Watch mode with 7-day period
```

## 📖 What It Does

TokenWatch CLI provides comprehensive monitoring of your OpenAI API usage:

- **Token Consumption**: Input, output, and total tokens per model
- **Cost Analysis**: Real-time cost tracking and daily averages
- **Model Breakdown**: Detailed usage statistics for each model
- **Smart Recommendations**: Period-specific guidance for optimal monitoring
- **Real-time Updates**: Watch mode for continuous monitoring

## 🛠️ Installation

### Download Pre-built Binary (Recommended for Users)

#### **Quick Install**
1. **Go to [Releases](https://github.com/mboss37/tokenwatch/releases)**
2. **Download v0.1.0** for your platform:
   - **Linux (x64)**: `tokenwatch-linux-amd64`
   - **Linux (ARM64)**: `tokenwatch-linux-arm64`
   - **macOS (Intel)**: `tokenwatch-darwin-amd64`
   - **macOS (Apple Silicon)**: `tokenwatch-darwin-arm64`
   - **Windows**: `tokenwatch-windows-amd64.exe`

#### **Install Commands**
```bash
# Download and install (Linux/macOS)
wget https://github.com/mboss37/tokenwatch/releases/download/v0.1.0/tokenwatch-linux-amd64
chmod +x tokenwatch-linux-amd64
sudo mv tokenwatch-linux-amd64 /usr/local/bin/tokenwatch

# Test installation
tokenwatch --version
```

### Install via Go (If you have Go installed)

```bash
# Install the specific version
go install github.com/mboss37/tokenwatch/cmd/root@v0.1.0

# Add to PATH if needed
export PATH="$HOME/go/bin:$PATH"

# Test installation
tokenwatch --version
```

### Build from Source (For Developers)

```bash
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch
git checkout v0.1.0
go build -o tokenwatch ./cmd/root
chmod +x tokenwatch
./tokenwatch --help
```

### Using Makefile

```bash
make build      # Build only
make build-all  # Build for all platforms
make install    # Install locally
```

## 🔧 Setup

### Quick Start
```bash
# 1. Install (choose method above)
# 2. Setup your API key
tokenwatch setup

# 3. Start monitoring
tokenwatch usage                    # Last 7 days
tokenwatch usage --period 1d       # Last 24 hours
tokenwatch usage -w -p 1d          # Watch mode for real-time updates
```

### Requirements
- **OpenAI Admin API key** with `api.usage.read` scope
- **Organization-level access** (personal API keys won't work)
- **Go 1.21+** (only if building from source)

## 📚 Documentation

- **[📖 User Guide](docs/USER_GUIDE.md)** - Complete usage documentation with examples
- **[🛠️ Developer Guide](docs/DEVELOPER.md)** - Contributing and architecture details

## 🔍 Basic Usage

```bash
# View OpenAI usage (last 7 days)
./tokenwatch usage

# Specify time period
./tokenwatch usage --period 1d       # Last 24 hours
./tokenwatch usage --period 7d       # Last 7 days (default)
./tokenwatch usage --period 30d      # Last 30 days

# Watch mode (works with all periods)
./tokenwatch usage -w -p 1d          # Real-time updates
./tokenwatch usage -w -p 7d          # Watch 7-day period

# Debug mode
./tokenwatch usage --debug            # API request/response logging
```

## 🏗️ Architecture

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
│  │  consumption.go │  │   pricing.go    │                  │ │
│  │ (usage models)  │  │ (cost models)   │                  │ │
│  └─────────────────┘  └─────────────────┘                  │ │
└─────────────────────────────────────────────────────────────┘
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

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- Uses [Viper](https://github.com/spf13/viper) for configuration
- Terminal tables powered by [tablewriter](https://github.com/olekukonko/tablewriter)
- Colors provided by [fatih/color](https://github.com/fatih/color)

## 📞 Support

- **Documentation**: [User Guide](docs/USER_GUIDE.md) | [Developer Guide](docs/DEVELOPER.md)
- **Issues**: [GitHub Issues](https://github.com/mboss37/tokenwatch/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mboss37/tokenwatch/discussions)
- **Repository**: [https://github.com/mboss37/tokenwatch](https://github.com/mboss37/tokenwatch)

---

**Happy monitoring! 🚀**

*For detailed usage instructions, examples, and troubleshooting, see the [User Guide](docs/USER_GUIDE.md).*