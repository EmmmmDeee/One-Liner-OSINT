# One-Liner OSINT - Production Go System

A production-grade, comprehensive OSINT (Open-Source Intelligence) gathering system built in Go. This system provides powerful tools for security researchers, bug bounty hunters, and ethical hackers to automate intelligence gathering tasks.

## Features

### Core Capabilities

- **Email Intelligence**: Email validation, breach checking, domain analysis
- **Domain Reconnaissance**: DNS enumeration, subdomain discovery, mail server identification
- **IP Analysis**: Geolocation, reverse DNS, private/public IP classification
- **Image Metadata**: EXIF extraction, GPS coordinates, camera information
- **Concurrent Execution**: Worker pools with rate limiting and retry logic
- **Caching System**: File-based caching with TTL support
- **Multiple Output Formats**: Text, JSON, CSV, and table formats
- **Extensible Architecture**: Modular design for easy addition of new OSINT modules

### Production-Grade Features

- ✅ Robust error handling and recovery
- ✅ Comprehensive logging with multiple levels
- ✅ Rate limiting to respect API quotas
- ✅ Retry logic with exponential backoff
- ✅ Configuration management (YAML, environment variables)
- ✅ Concurrent task execution with worker pools
- ✅ Response caching for improved performance
- ✅ Colored terminal output
- ✅ Progress indicators
- ✅ Comprehensive test coverage
- ✅ Cross-platform support (Linux, macOS, Windows)

## Architecture

```
One-Liner-OSINT/
├── cmd/osint/              # CLI entry point
│   ├── main.go
│   └── commands/           # CLI commands
│       ├── root.go
│       ├── email.go
│       ├── domain.go
│       ├── ip.go
│       └── image.go
├── pkg/                    # Public packages
│   ├── osint/             # OSINT modules
│   │   ├── email.go
│   │   ├── domain.go
│   │   ├── ip.go
│   │   └── image.go
│   ├── config/            # Configuration management
│   ├── logger/            # Logging utilities
│   ├── output/            # Output formatting
│   └── utils/             # Utility functions
├── internal/              # Private packages
│   ├── engine/            # Core execution engine
│   ├── worker/            # Worker pool implementation
│   └── cache/             # Caching system
└── docs/                  # Documentation
```

## Installation

### Prerequisites

- Go 1.22 or higher
- Make (optional, for using Makefile)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/EmmmmDeee/One-Liner-OSINT.git
cd One-Liner-OSINT

# Download dependencies
go mod download

# Build the binary
make build

# Or build directly with go
go build -o bin/osint ./cmd/osint
```

### Install Binary

```bash
# Install to $GOPATH/bin
make install

# Or
go install ./cmd/osint
```

### Cross-Platform Builds

```bash
# Build for all platforms
make build-all

# Build for specific platform
make build-linux
make build-darwin
make build-windows
```

## Usage

### Basic Commands

```bash
# Get help
osint --help

# Analyze email
osint email user@example.com

# Analyze domain
osint domain example.com

# Analyze domain with subdomain enumeration
osint domain example.com --enumerate-subdomains

# Analyze IP address
osint ip 8.8.8.8

# Extract image metadata
osint image /path/to/image.jpg

# Or from URL
osint image https://example.com/image.jpg
```

### Output Formats

```bash
# JSON output
osint email user@example.com --output json

# CSV output
osint domain example.com --output csv

# Table output
osint ip 8.8.8.8 --output table
```

### Advanced Options

```bash
# Verbose logging
osint email user@example.com --verbose

# Custom worker count
osint domain example.com --workers 20

# Custom timeout
osint ip 8.8.8.8 --timeout 60

# Disable colors
osint email user@example.com --no-color

# Use config file
osint email user@example.com --config ~/.osint.yaml
```

## Configuration

Create a configuration file at `~/.osint.yaml`:

```yaml
# General settings
verbose: false
output: text
workers: 10
timeout: 30
user_agent: "OSINT-Tool/1.0"

# API Keys
api_keys:
  haveibeenpwned: "your-api-key"
  hunter_io: "your-api-key"
  shodan: "your-api-key"
  virustotal: "your-api-key"
  google_cse: "your-api-key"
  google_cse_cx: "your-cx-id"
  ipgeolocation: "your-api-key"
  github: "your-token"

# Rate limiting
rate_limit:
  requests_per_second: 10
  burst_size: 20
  retry_attempts: 3
  retry_delay: 2

# Search settings
search:
  max_results: 100
  max_depth: 3
  follow_redirects: true
  use_proxy: false

# Output settings
output:
  format: text
  pretty: true
  include_empty: false
  timestamp: true
  save_to_file: false
  output_dir: ./output

# Cache settings
cache:
  dir: .osint_cache
```

## API Keys

Some modules require API keys for full functionality:

- **Have I Been Pwned**: Breach checking - [Get API Key](https://haveibeenpwned.com/API/Key)
- **Hunter.io**: Email verification - [Get API Key](https://hunter.io/api)
- **Shodan**: Infrastructure scanning - [Get API Key](https://account.shodan.io/)
- **VirusTotal**: File/URL scanning - [Get API Key](https://www.virustotal.com/gui/join-us)
- **IPGeolocation**: IP geolocation - [Get API Key](https://ipgeolocation.io/)

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific package tests
go test ./pkg/osint/... -v
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

### Adding New Modules

1. Create a new file in `pkg/osint/`:

```go
package osint

import (
    "context"
    "github.com/EmmmmDeee/One-Liner-OSINT/pkg/config"
    "github.com/EmmmmDeee/One-Liner-OSINT/pkg/output"
)

type YourModule struct {
    config *config.Config
}

func NewYourModule(cfg *config.Config) *YourModule {
    return &YourModule{config: cfg}
}

func (m *YourModule) Name() string {
    return "your-module"
}

func (m *YourModule) Description() string {
    return "Description of your module"
}

func (m *YourModule) Execute(ctx context.Context, target string) ([]output.Result, error) {
    // Implementation
    return nil, nil
}
```

2. Register the module in the engine
3. Create a CLI command in `cmd/osint/commands/`

## Performance

- **Concurrent Execution**: Utilizes worker pools for parallel task execution
- **Rate Limiting**: Respects API rate limits with configurable requests per second
- **Caching**: Reduces redundant API calls with file-based caching
- **Efficient Memory Usage**: Streaming responses where possible
- **Timeout Management**: Prevents hanging operations

## Security Considerations

- API keys are never logged or included in output
- SSL/TLS verification is enabled by default
- User agents are configurable
- Rate limiting prevents abuse
- Input validation on all user-provided data

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details

## Disclaimer

This tool is intended for legal and ethical use only. Users are responsible for complying with applicable laws and regulations. The authors are not responsible for misuse or damage caused by this tool.

## Support

- **Issues**: [GitHub Issues](https://github.com/EmmmmDeee/One-Liner-OSINT/issues)
- **Documentation**: See `docs/` directory
- **Examples**: See `examples/` directory

## Roadmap

- [ ] Additional OSINT modules (social media, public records, etc.)
- [ ] Web UI interface
- [ ] API server mode
- [ ] Database storage for results
- [ ] Report generation (PDF, HTML)
- [ ] Plugin system for custom modules
- [ ] Docker support
- [ ] Kubernetes deployment manifests

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Logrus](https://github.com/sirupsen/logrus) - Logging
- [GoQuery](https://github.com/PuerkitoBio/goquery) - HTML parsing
- [GoExif](https://github.com/rwcarlsen/goexif) - EXIF extraction

## Authors

- Original OSINT documentation: YogSec
- Go implementation: EmmmmDeee

---

**Built with Go for maximum performance, reliability, and scalability.**
