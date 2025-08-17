# TokenWatch CLI - User Guide 📚

A comprehensive guide to using TokenWatch CLI for monitoring AI token usage and costs.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Installation](#installation)
3. [Setup](#setup)
4. [Basic Usage](#basic-usage)
5. [Watch Mode](#watch-mode)
6. [Debug Mode](#debug-mode)
7. [Configuration](#configuration)
8. [Troubleshooting](#troubleshooting)

## Quick Start 🚀

```bash
# Install and setup
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch
go build -o tokenwatch ./cmd/root
./tokenwatch setup

# Basic usage
./tokenwatch openai                    # Last 7 days
./tokenwatch openai --period 30d      # Last 30 days
./tokenwatch openai -w                 # Watch mode - refresh every 30s
./tokenwatch all                       # All platforms
```

## Installation

### Build from Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch

# Build the binary
go build -o tokenwatch ./cmd/root

# Make it executable (Linux/macOS)
chmod +x tokenwatch

# Test it works
./tokenwatch --help
```

### Install to System PATH

**Linux/macOS:**
```bash
# Build and install to ~/go/bin
go build -o tokenwatch ./cmd/root
cp tokenwatch ~/go/bin/

# Add to PATH (add this to your ~/.bashrc, ~/.zshrc, or ~/.profile)
export PATH="$HOME/go/bin:$PATH"

# Or install system-wide (requires sudo)
sudo cp tokenwatch /usr/local/bin/
```

**Windows:**
```bash
# Build for Windows
go build -o tokenwatch.exe ./cmd/root

# Add the directory to your PATH environment variable
# Or run from the current directory: .\tokenwatch.exe --help
```

### Using Makefile

```bash
# Build only
make build

# Build for all platforms
make build-all

# Install locally (puts binary in ~/go/bin as 'tokenwatch')
make install

# Note: make install creates binary named 'tokenwatch', not 'root'
```

## Setup

```bash
# Interactive setup for API keys
./tokenwatch setup
```

This will:
- Prompt for your OpenAI API key
- Validate the API key by making a test call
- Save the configuration to `~/.tokenwatch/config.yaml`

## Basic Usage

### OpenAI Usage

```bash
# View OpenAI usage for last 7 days
./tokenwatch openai

# Specify time period
./tokenwatch openai --period 7d      # Last 7 days (default)
./tokenwatch openai --period 30d     # Last 30 days
./tokenwatch openai --period 90d     # Last 90 days

# Short flags
./tokenwatch openai -p 30d           # Same as --period 30d
```

### All Platforms

```bash
# View combined data from all configured platforms
./tokenwatch all

# With specific time period
./tokenwatch all --period 30d
```

### Configuration Management

```bash
# Check configuration status
./tokenwatch config check

# Reset configuration
./tokenwatch config reset

# View version
./tokenwatch version
```

## Watch Mode

Watch mode provides real-time monitoring of your AI usage with automatic refresh every 30 seconds:

```bash
# Watch OpenAI usage
./tokenwatch openai -w

# Watch all platforms
./tokenwatch all -w

# Watch with specific period
./tokenwatch openai -w -p 30d

# Stop watching: Press Ctrl+C
```

**Features:**
- **Auto-refresh**: Updates every 30 seconds
- **Screen clearing**: Clean display on each refresh
- **Fresh data**: Bypasses cache for real-time information
- **Easy exit**: Ctrl+C to stop

## Debug Mode

Debug mode shows detailed API request/response information for troubleshooting and development:

```bash
# Enable debug mode for OpenAI
./tokenwatch openai --debug

# Debug mode with watch
./tokenwatch openai -w --debug

# Debug mode for all platforms
./tokenwatch all --debug

# Debug with specific period
./tokenwatch openai --debug --period 30d
```

**Debug Output Includes:**
- **API Request Details**: URL, timestamps, parameters
- **Raw JSON Responses**: Complete OpenAI API responses
- **Request/Response Flow**: Full API call lifecycle

**Use Cases:**
- Troubleshooting API issues
- Verifying data freshness
- Development and testing
- Understanding API behavior

## Configuration

### Config File Location

Configuration is stored in `~/.tokenwatch/config.yaml`

### Environment Variables

```bash
# API Keys
export TOKENWATCH_OPENAI_API_KEY="sk-..."

# Logging
export TOKENWATCH_LOG_LEVEL="debug"
```

### Config File Structure

```yaml
api_keys:
  openai: "sk-..."
  anthropic: "your-key-here"
  grok: "your-key-here"
  cursor: "your-key-here"

settings:
  debug: false
```

## Example Output

### OpenAI Usage (Normal Mode)

```
🤖 OPENAI USAGE - Last 7d
⏰ Generated: 2025-08-17 13:30:13

📊 SUMMARY
──────────────────────────────────────────────────
📅 Period: 2025-08-10 to 2025-08-17 (7 days)
📈 Daily Averages: 721.7 tokens, 1.9 requests
💰 Daily Cost Average: $0.0051

📋 MODEL BREAKDOWN
┌────────────────────────┬──────────────┬───────────────┬──────────────┬──────────┬─────────┬───────────────┐
│         MODEL          │ INPUT TOKENS │ OUTPUT TOKENS │ TOTAL TOKENS │ REQUESTS │  COST   │ $/ 1 K TOKENS │
├────────────────────────┼──────────────┼───────────────┼──────────────┼──────────┼─────────┼───────────────┤
│ gpt-4o-2024-08-06      │ 1872         │ 3099          │ 4971         │ 12       │ $0.0357 │ $0.0072       │
│ gpt-4o-mini-2024-07-18 │ 10           │ 71            │ 81           │ 1        │ $0.0000 │ $0.0005       │
│ ─                      │ ─            │ ─             │ ─            │ ─        │ ─       │ ─             │
│ TOTAL                  │ 1882         │ 3170          │ 5052         │ 13       │ $0.0357 │ $0.0071       │
└────────────────────────┴──────────────┴───────────────┴──────────────┴──────────┴─────────┴───────────────┘
```

### Debug Mode Output

```
🔍 OPENAI USAGE API REQUEST:
   URL: https://api.openai.com/v1/organization/usage/completions?...
   Start Time: 2025-08-10 13:30:26 (1754825426)
   End Time: 2025-08-17 13:30:26 (1755430226)
   Bucket Width: 1d
   Group By: [model]

🔍 RAW OPENAI USAGE API RESPONSE:
{
  "data": [
    {
      "start_time": 1754825426,
      "end_time": 1754870400,
      "results": []
    },
    // ... more data
  ]
}
```

## Troubleshooting

### Common Issues

**"OpenAI not configured"**
```bash
./tokenwatch setup
```

**"API key lacks required permissions"**
- You need an Admin API key with `api.usage.read` scope
- Check your OpenAI organization settings

**"No data found"**
- Try a shorter period: `./tokenwatch openai --period 7d`
- Verify you have recent API usage
- Check if your API key is valid

**"Permission denied"**
```bash
chmod +x tokenwatch
```

### Debug Mode for Troubleshooting

```bash
# Enable debug mode to see API details
./tokenwatch openai --debug

# Check configuration
./tokenwatch config check

# Verify API key
./tokenwatch setup
```

### Logging

```bash
# Enable debug logging
export TOKENWATCH_LOG_LEVEL=debug
./tokenwatch openai

# Or use debug mode for API details
./tokenwatch openai --debug
```

## Platform Support

| Platform | Status | Description |
|----------|--------|-------------|
| OpenAI | ✅ Ready | ChatGPT, GPT-4, DALL-E |
| Anthropic | 🚧 Coming Soon | Claude models |
| Grok | 🚧 Coming Soon | xAI's Grok |
| Cursor | 🚧 Coming Soon | Cursor AI |

> **Note**: Anthropic, Grok, and Cursor currently don't provide usage APIs. Support will be added when they make these APIs available.

## Advanced Features

### Cache Management

- **Normal mode**: 5-minute cache for efficiency
- **Watch mode**: Cache bypassed for real-time data
- **Debug mode**: Shows cache behavior

### Rate Limiting

- **OpenAI**: 1 request/second with burst of 5
- **Automatic retries** with exponential backoff
- **Circuit breaker** to prevent cascading failures

### Data Freshness

- **Real-time data** in watch mode
- **Fresh API calls** every 30 seconds
- **Cache bypass** when needed

## Getting Help

```bash
# Command help
./tokenwatch --help
./tokenwatch openai --help
./tokenwatch all --help

# Version info
./tokenwatch version

# Configuration check
./tokenwatch config check
```

For more information, see the [Developer Guide](DEVELOPER.md) or visit the [GitHub repository](https://github.com/mboss37/tokenwatch).