# TokenWatch CLI - User Guide

Welcome to TokenWatch CLI! This guide will help you get started with monitoring your AI token consumption and costs across multiple platforms.

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Commands Overview](#commands-overview)
4. [Platform-Specific Commands](#platform-specific-commands)
5. [Watch Mode](#watch-mode)
6. [Configuration](#configuration)
7. [Troubleshooting](#troubleshooting)

## Installation

### Download Binary

Download the latest release for your platform:
- **macOS**: `tokenwatch-darwin-amd64` or `tokenwatch-darwin-arm64`
- **Linux**: `tokenwatch-linux-amd64`
- **Windows**: `tokenwatch-windows-amd64.exe`

Make it executable (macOS/Linux):
```bash
chmod +x tokenwatch
```

### Build from Source

```bash
git clone https://github.com/mboss37/tokenwatch.git
cd tokenwatch/tokenwatchcli
go build -o tokenwatch ./cmd/root
```

## Quick Start

1. **Set up your first platform:**
   ```bash
   tokenwatch setup
   ```

2. **View your usage:**
   ```bash
   tokenwatch openai    # View OpenAI usage
   tokenwatch all       # View all platforms
   ```

3. **Monitor in real-time:**
   ```bash
   tokenwatch openai -w  # Watch mode - refreshes every 30s
   ```

## Commands Overview

### Core Commands

- `tokenwatch setup` - Interactive setup for API keys
- `tokenwatch config` - Manage configuration
- `tokenwatch all` - View combined data from all platforms
- `tokenwatch openai` - View OpenAI usage and costs
- `tokenwatch version` - Display version information

### Global Flags

- `--period, -p` - Time period: `7d` (default), `30d`, or `90d`
- `--watch, -w` - Watch mode - auto-refresh every 30 seconds
- `--help, -h` - Show help for any command

## Platform-Specific Commands

### OpenAI

View your OpenAI token consumption and costs:

```bash
# Last 7 days (default)
tokenwatch openai

# Last 30 days
tokenwatch openai --period 30d

# Last 90 days
tokenwatch openai --period 90d

# Watch mode - real-time monitoring
tokenwatch openai -w
tokenwatch openai -w -p 30d
```

**Note:** Currently, only OpenAI provides usage APIs. Support for Anthropic, Grok, and Cursor will be added when they provide similar APIs.

## Watch Mode

Watch mode provides real-time monitoring of your AI usage with automatic refresh every 30 seconds:

```bash
# Watch OpenAI usage
tokenwatch openai -w

# Watch all platforms
tokenwatch all -w

# Watch with custom period
tokenwatch openai -w -p 30d
```

Features:
- Auto-refreshes every 30 seconds
- Clear screen between updates
- Press `Ctrl+C` to stop watching
- Works with all period options

## Configuration

### Setup Command

Interactive setup for all platforms:
```bash
tokenwatch setup
```

Follow the prompts to:
1. Select a platform
2. Enter your API key (masked input)
3. Validate the API key
4. Optionally configure additional platforms

### Config Management

```bash
# Check current configuration
tokenwatch config check

# Reset all configuration
tokenwatch config reset
```

### Manual Configuration

Configuration is stored in `~/.tokenwatch/config.yaml`:

```yaml
api_keys:
  openai: sk-your-api-key-here
  anthropic: your-anthropic-key
  grok: your-grok-key
  cursor: your-cursor-key

settings:
  debug: false
```

### Environment Variables

You can also use environment variables:
- `TOKENWATCH_OPENAI_API_KEY`
- `TOKENWATCH_ANTHROPIC_API_KEY`
- `TOKENWATCH_GROK_API_KEY`
- `TOKENWATCH_CURSOR_API_KEY`
- `TOKENWATCH_LOG_LEVEL` (debug, info, warn, error)

## Troubleshooting

### Common Issues

**"Platform not configured" error:**
- Run `tokenwatch setup` to configure the platform
- Or check `tokenwatch config check` to see what's missing

**"No data found" message:**
- Ensure you've made API calls during the selected period
- Try a shorter period like `--period 7d`
- Note: Cost data may not be available for periods longer than 30 days

**API Key Issues:**
- Verify your API key is correct (should start with 'sk-' for OpenAI)
- Check that your key has the necessary permissions
- For OpenAI: Ensure your organization has usage tracking enabled

### Debug Mode

Enable debug logging:
```bash
export TOKENWATCH_LOG_LEVEL=debug
tokenwatch openai
```

Or set in config:
```yaml
settings:
  debug: true
```

## Features

### Production-Ready Enhancements
- **Retry Logic**: Automatic retries for failed API calls
- **Rate Limiting**: Respects API rate limits (1 req/sec with burst of 5)
- **Circuit Breaker**: Prevents cascading failures
- **Caching**: 5-minute cache to reduce API calls
- **Structured Logging**: Better debugging and monitoring
- **API Key Validation**: Validates keys during setup

### Data Presentation
- **Detailed Model Breakdown**: Shows usage per model
- **Cost Analysis**: Per-token costs and daily averages
- **Total Rows**: Clear summaries in all tables
- **Color-Coded Output**: Easy to read terminal output

## Next Steps

- Check out the [Developer Guide](DEVELOPER.md) if you want to contribute
- Report issues on [GitHub](https://github.com/mboss37/tokenwatch/issues)
- Star the project if you find it useful!