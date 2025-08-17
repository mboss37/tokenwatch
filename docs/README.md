# TokenWatch CLI - User Guide 📚

A comprehensive guide to using TokenWatch CLI for monitoring OpenAI token usage and costs.

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
./tokenwatch usage                    # Last 7 days
./tokenwatch usage --period 1d       # Last 24 hours
./tokenwatch usage --period 30d      # Last 30 days
./tokenwatch usage -w -p 1d          # Watch mode for real-time monitoring
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
```

## Setup

```bash
# Interactive setup for OpenAI API key
./tokenwatch setup
```

This will:
- Prompt for your OpenAI Admin API key
- Validate the API key by making a test call
- Save the configuration to `~/.tokenwatch/config.yaml`

**Important**: You need an OpenAI Admin API key with `api.usage.read` scope for organization-level access.

## Basic Usage

### OpenAI Usage

```bash
# View OpenAI usage for last 7 days (default)
./tokenwatch usage

# Specify time period
./tokenwatch usage --period 1d       # Last 24 hours
./tokenwatch usage --period 7d       # Last 7 days (default)
./tokenwatch usage --period 30d      # Last 30 days

# Short flags
./tokenwatch usage -p 1d             # Same as --period 1d
```

**Available Time Periods:**
- `1d` - Last 24 hours (perfect for recent activity)
- `7d` - Last 7 days (ideal for historical data)
- `30d` - Last 30 days (may have limited data)

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

Watch mode provides real-time monitoring of your OpenAI usage with automatic refresh every 30 seconds:

```bash
# Watch OpenAI usage (1-day period only)
./tokenwatch usage -w -p 1d

# Stop watching: Press Ctrl+C
```

**Features:**
- **Auto-refresh**: Updates every 30 seconds
- **Screen clearing**: Clean display on each refresh
- **Fresh data**: Bypasses cache for real-time information
- **1-day only**: Watch mode is only available for recent activity
- **Easy exit**: Ctrl+C to stop

**Note**: Watch mode is only available for 1-day periods since longer periods don't need real-time updates.

## Debug Mode

Debug mode shows detailed API request/response information for troubleshooting and development:

```bash
# Enable debug mode for OpenAI
./tokenwatch usage --debug

# Debug mode with specific period
./tokenwatch usage --debug --period 30d
```

**Debug Output Includes:**
- **API Request Details**: URL, timestamps, parameters
- **Raw JSON Responses**: Complete OpenAI API responses
- **Pagination Flow**: Shows how data is fetched across multiple pages
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
export OPENAI_API_KEY="sk-..."

# Logging
export TOKENWATCH_LOG_LEVEL="debug"
```

### Config File Structure

```yaml
api_keys:
  openai: "sk-..."

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

### Debug Mode Output

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
- Try a shorter period: `./tokenwatch usage --period 7d`
- Verify you have recent API usage
- Check if your API key is valid

**"Permission denied"**
```bash
chmod +x tokenwatch
```

### Debug Mode for Troubleshooting

```bash
# Enable debug mode to see API details
./tokenwatch usage --debug

# Check configuration
./tokenwatch config check

# Verify API key
./tokenwatch setup
```

### Logging

```bash
# Enable debug logging
export TOKENWATCH_LOG_LEVEL=debug
./tokenwatch usage

# Or use debug mode for API details
./tokenwatch usage --debug
```

## Platform Support

| Platform | Status | Description |
|----------|--------|-------------|
| OpenAI | ✅ Ready | ChatGPT, GPT-4, DALL-E |

**Note**: Currently only OpenAI is supported as it's the only platform that provides comprehensive usage and costs APIs.

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

### Smart Recommendations

The CLI provides intelligent recommendations based on your selected time period:
- **1-day**: Perfect for recent activity monitoring
- **7-day**: Ideal for historical analysis and weekly patterns
- **30-day**: May have limited data due to API limitations

## Getting Help

```bash
# Command help
./tokenwatch --help
./tokenwatch usage --help

# Version info
./tokenwatch version

# Configuration check
./tokenwatch config check
```

For more information, see the [Developer Guide](DEVELOPER.md) or visit the [GitHub repository](https://github.com/mboss37/tokenwatch).