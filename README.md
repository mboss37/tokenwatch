# TokenWatch CLI 🔍

Track your AI token usage and costs across multiple platforms from your terminal.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Status](https://img.shields.io/badge/Status-Production-green)
![Platform](https://img.shields.io/badge/Platform-OpenAI-412991)

## Features ✨

- **📊 Real-time Usage Tracking** - Monitor your token consumption across models
- **💰 Cost Analysis** - See detailed costs with per-token pricing
- **📺 Watch Mode** - Live monitoring with auto-refresh every 30 seconds
- **🔍 Debug Mode** - Optional API request/response logging for troubleshooting
- **🚀 Fast & Lightweight** - Built with Go for speed and efficiency
- **🔒 Secure** - API keys stored locally, never transmitted
- **🎨 Beautiful Output** - Clean, colorful terminal tables with totals
- **⚡ Production Ready** - Retry logic, rate limiting, circuit breaker

## Quick Start 🚀

### Install

#### Option 1: Build from Source (Recommended)

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

#### Option 2: Install to System PATH

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

#### Option 3: Using Makefile

```bash
# Build only
make build

# Build for all platforms
make build-all

# Install locally (puts binary in ~/go/bin as 'tokenwatch')
make install

# Note: make install creates binary named 'tokenwatch', not 'root'
```

#### Verify Installation

After installation, verify it works:

```bash
# Check if tokenwatch is in PATH
which tokenwatch  # Linux/macOS
where tokenwatch  # Windows

# Test the command
tokenwatch --help
tokenwatch version
```

### Setup

```bash
# Interactive setup
./tokenwatch setup
```

### Use

```bash
# View OpenAI usage
./tokenwatch openai

# Watch mode - auto-refresh every 30s
./tokenwatch openai -w

# Debug mode - see raw API data
./tokenwatch openai --debug

# Last 30 days
./tokenwatch openai --period 30d

# All platforms
./tokenwatch all

# All platforms with debug
./tokenwatch all --debug
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

### All Platforms View

```
🎯 TOKENWATCH ALL PLATFORMS - Last 7d
⏰ Generated: 2025-08-17 12:22:21
🔗 Period: Last 7 days

🚀 Fetching data from 1 configured platform(s)...

📊 COMBINED DASHBOARD METRICS
───────────────────────────────────────────────────
🔤 Combined Token Usage:
   Total Tokens: 5052
   Total Requests: 13
   Average per Request: 388.6
   Daily Average: 721.7

💰 Combined Cost Analysis:
   Total Cost: $0.0357
   Daily Average: $0.0051
   Cost per Token: $0.000007

📅 Time Period: 2025-08-10 to 2025-08-17 (7 days)

📋 ALL PLATFORMS MODEL BREAKDOWN
┌──────────┬────────────────────────┬──────────────┬───────────────┬──────────────┬──────────┬─────────┬───────────────┐
│ PLATFORM │         MODEL          │ INPUT TOKENS │ OUTPUT TOKENS │ TOTAL TOKENS │ REQUESTS │  COST   │ $/ 1 K TOKENS │
├──────────┼────────────────────────┼──────────────┼───────────────┼──────────────┼──────────┼─────────┼───────────────┤
│ Openai   │ gpt-4o-2024-08-06      │ 1872         │ 3099          │ 4971         │ 12       │ $0.0357 │ $0.0072       │
│ Openai   │ gpt-4o-mini-2024-07-18 │ 10           │ 71            │ 81           │ 1        │ $0.0000 │ $0.0005       │
│ ─        │ ─                      │ ─            │ ─             │ ─            │ ─        │ ─       │ ─             │
│ TOTAL    │                        │ 1882         │ 3170          │ 5052         │ 13       │ $0.0357 │ $0.0071       │
└──────────┴────────────────────────┴──────────────┴───────────────┴──────────────┴──────────┴─────────┴───────────────┘
```

### Configuration Check

```
🔍 CONFIGURATION STATUS
───────────────────────────────────────────────────
✅ Config file: /home/mboss37/.tokenwatch/config.yaml

🔑 API KEYS:
   ✅ Openai: Configured
   ❌ Anthropic: Not configured
   ❌ Grok: Not configured
   ❌ Cursor: Not configured

✅ Configuration check complete!
```

### Watch Mode

```bash
./tokenwatch openai -w
# Auto-refreshes every 30 seconds
# Press Ctrl+C to stop

🔄 Refreshing every 30 seconds... (Press Ctrl+C to stop)
```

## Debug Mode 🔍

Debug mode provides optional API request/response logging for troubleshooting and development:

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

## Requirements

- **OpenAI**: Admin API key with `api.usage.read` scope
- **Go**: Version 1.21 or higher (for building from source)

## Platform Support

| Platform | Status | Description |
|----------|--------|-------------|
| OpenAI | ✅ Ready | ChatGPT, GPT-4, DALL-E |
| Anthropic | 🚧 Coming Soon | Claude models |
| Grok | 🚧 Coming Soon | xAI's Grok |
| Cursor | 🚧 Coming Soon | Cursor AI |

> **Note**: Anthropic, Grok, and Cursor currently don't provide usage APIs. Support will be added when they make these APIs available.

## Configuration

Config stored in `~/.tokenwatch/config.yaml`

```yaml
api_keys:
  openai: "sk-..."
  anthropic: "your-key-here"
  grok: "your-key-here"
  cursor: "your-key-here"

settings:
  debug: false
```

Environment variables also supported:
- `TOKENWATCH_OPENAI_API_KEY`
- `TOKENWATCH_LOG_LEVEL`

## Documentation 📚

- **[User Guide](docs/README.md)** - Complete usage documentation
- **[Developer Guide](docs/DEVELOPER.md)** - Contributing and architecture

## Key Features Under the Hood 🔧

- **Retry Logic** - Automatic retries with exponential backoff
- **Rate Limiting** - Respects API limits (1 req/sec with burst of 5)
- **Circuit Breaker** - Prevents cascading failures
- **Response Caching** - 5-minute cache to reduce API calls
- **Structured Logging** - Debug with `TOKENWATCH_LOG_LEVEL=debug`
- **Smart Errors** - Helpful suggestions when things go wrong
- **API Validation** - Validates API keys during setup
- **Debug Mode** - Optional API request/response logging with `--debug` flag

## Contributing 🤝

We welcome contributions! Please read our [Developer Guide](docs/DEVELOPER.md) first.

Key principle: **Platform Separation** - Each LLM platform has isolated implementations (dedicated files, no mixing of logic).

## Troubleshooting 🔧

### Installation Issues

**"command not found: tokenwatch"**
→ The binary isn't in your PATH. Try:
```bash
# Check where the binary is
ls -la ~/go/bin/  # If using go install
ls -la ./          # If built locally

# Add to PATH
export PATH="$HOME/go/bin:$PATH"  # For ~/go/bin
export PATH="$PWD:$PATH"          # For current directory
```

**"Permission denied"**
→ Make the binary executable:
```bash
chmod +x tokenwatch
```

**"go: command not found"**
→ Install Go first: https://golang.org/dl/

### Usage Issues

**"OpenAI not configured"**
→ Run `./tokenwatch setup`

**"API key lacks required permissions"**  
→ You need an Admin API key with `api.usage.read` scope

**"No data found"**
→ Try a shorter period: `./tokenwatch openai --period 7d`

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

## License

MIT License - see [LICENSE](LICENSE) file

---

Built with ❤️ for developers who care about their AI costs