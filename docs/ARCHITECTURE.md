# Architecture Documentation

## System Overview

One-Liner OSINT is a production-grade Go system for Open-Source Intelligence gathering. The system is built with a modular, extensible architecture that prioritizes performance, reliability, and scalability.

## Design Principles

1. **Modularity**: Each OSINT capability is implemented as an independent module
2. **Concurrency**: Worker pools enable parallel execution of tasks
3. **Resilience**: Retry logic and error handling ensure robust operation
4. **Performance**: Caching and rate limiting optimize resource usage
5. **Extensibility**: Easy to add new modules and capabilities
6. **Production-Ready**: Comprehensive logging, configuration, and error handling

## Architecture Layers

### 1. CLI Layer (`cmd/osint`)
- Entry point for the application
- Command-line interface using Cobra
- Command handlers for each OSINT module

### 2. Command Layer (`cmd/osint/commands`)
- Individual command implementations
- Flag parsing and validation
- Integration with the engine layer

### 3. Engine Layer (`internal/engine`)
- Core execution engine
- Module registration and management
- Task orchestration
- Cache integration
- Result aggregation

### 4. Module Layer (`pkg/osint`)
- Individual OSINT modules (email, domain, IP, image)
- Module interface implementation
- External API integrations
- Data extraction and analysis

### 5. Infrastructure Layer
- **Worker Pool** (`internal/worker`): Concurrent task execution
- **Cache** (`internal/cache`): File-based result caching
- **Config** (`pkg/config`): Configuration management
- **Logger** (`pkg/logger`): Structured logging
- **Output** (`pkg/output`): Result formatting
- **Utils** (`pkg/utils`): HTTP client and utilities

## Component Details

### Engine (`internal/engine`)

The engine is the core orchestration layer that:
- Manages module lifecycle
- Coordinates concurrent execution
- Handles caching
- Aggregates results

```go
type Engine struct {
    config  *config.Config
    logger  *logger.Logger
    cache   *cache.Cache
    pool    *worker.Pool
    modules map[string]Module
}
```

Key methods:
- `RegisterModule()`: Add new OSINT modules
- `Execute()`: Run a single module
- `ExecuteMultiple()`: Run multiple modules concurrently
- `ExecuteAll()`: Run all registered modules

### Module Interface

All OSINT modules implement this interface:

```go
type Module interface {
    Name() string
    Description() string
    Execute(ctx context.Context, target string) ([]output.Result, error)
}
```

### Worker Pool (`internal/worker`)

Implements concurrent task execution with:
- Configurable worker count
- Rate limiting
- Retry logic with exponential backoff
- Context-based cancellation

```go
type Pool struct {
    workers    int
    tasks      chan Task
    results    chan error
    limiter    *rate.Limiter
    ctx        context.Context
    cancel     context.CancelFunc
    maxRetries int
    retryDelay time.Duration
}
```

### Cache (`internal/cache`)

File-based caching with:
- TTL-based expiration
- SHA-256 key hashing
- JSON serialization
- Automatic cleanup

```go
type Cache struct {
    dir string
    ttl time.Duration
}
```

### Configuration (`pkg/config`)

Hierarchical configuration from:
1. Default values
2. Configuration file (YAML)
3. Environment variables
4. Command-line flags

Configuration structure:
- General settings (workers, timeout, verbose)
- API keys for external services
- Rate limiting parameters
- Search settings
- Output preferences

### Output Formatter (`pkg/output`)

Supports multiple output formats:
- **Text**: Human-readable colored output
- **JSON**: Machine-readable structured data
- **CSV**: Spreadsheet-compatible format
- **Table**: Tabular display

```go
type Result struct {
    Timestamp time.Time
    Source    string
    Type      string
    Target    string
    Data      map[string]interface{}
    Error     string
    Metadata  map[string]string
}
```

## Data Flow

1. **User Input** → CLI parses command and flags
2. **Command Handler** → Validates input, creates engine
3. **Engine** → Checks cache, registers module
4. **Module** → Executes OSINT logic, calls APIs
5. **Results** → Cached and returned to engine
6. **Formatter** → Converts results to desired format
7. **Output** → Displayed to user

## Concurrency Model

### Worker Pool Pattern

```
┌─────────────┐
│   Tasks     │ ← Submit tasks
└──────┬──────┘
       │
   ┌───▼───┐
   │ Queue │
   └───┬───┘
       │
   ┌───▼────────────┐
   │  Rate Limiter  │
   └───┬────────────┘
       │
   ┌───▼────┬────┬────┐
   │ Worker │ ... │ Worker │
   └───┬────┴────┴────┘
       │
   ┌───▼────────┐
   │  Results   │
   └────────────┘
```

### Rate Limiting

Uses token bucket algorithm:
- Configurable requests per second
- Burst capacity
- Per-module rate limits

## Error Handling

### Error Strategy

1. **Validation Errors**: Return immediately to user
2. **Transient Errors**: Retry with exponential backoff
3. **API Errors**: Log and continue with degraded functionality
4. **Fatal Errors**: Graceful shutdown with cleanup

### Error Types

```go
type EngineError struct {
    message string
}
```

### Retry Logic

- Configurable retry attempts (default: 3)
- Exponential backoff delay (default: 2s)
- Context-based cancellation

## Caching Strategy

### Cache Key Generation

```
key = module_name + ":" + target
hash = SHA256(key)
filename = hash + ".json"
```

### Cache Lifecycle

1. Check cache on module execution
2. Return cached result if valid (not expired)
3. Execute module if cache miss
4. Store result in cache
5. Periodic cleanup of expired entries

### Cache TTL

- Default: 1 hour
- Configurable per deployment
- Different TTL for different module types

## Module Development

### Creating a New Module

1. Create file in `pkg/osint/`:

```go
package osint

type MyModule struct {
    config *config.Config
    logger *logger.Logger
    client *utils.HTTPClient
}

func NewMyModule(cfg *config.Config) *MyModule {
    return &MyModule{
        config: cfg,
        logger: logger.NewLogger(),
        client: utils.NewHTTPClient(cfg.Timeout, cfg.UserAgent),
    }
}

func (m *MyModule) Name() string {
    return "mymodule"
}

func (m *MyModule) Description() string {
    return "My custom OSINT module"
}

func (m *MyModule) Execute(ctx context.Context, target string) ([]output.Result, error) {
    // Implementation
    return []output.Result{}, nil
}
```

2. Create command in `cmd/osint/commands/`:

```go
package commands

var mymoduleCmd = &cobra.Command{
    Use:   "mymodule [target]",
    Short: "Description",
    Args:  cobra.ExactArgs(1),
    RunE:  runMyModule,
}

func init() {
    rootCmd.AddCommand(mymoduleCmd)
}

func runMyModule(cmd *cobra.Command, args []string) error {
    // Command implementation
}
```

3. Register module in engine

## API Integration

### HTTP Client

Wrapper around `net/http` with:
- Timeout configuration
- Custom user agent
- Header management
- Proxy support
- TLS configuration

### External APIs

Supported APIs:
- Have I Been Pwned (breach checking)
- Hunter.io (email verification)
- Shodan (infrastructure scanning)
- IPGeolocation (IP data)
- VirusTotal (file/URL scanning)

### API Key Management

- Stored in configuration file
- Never logged or displayed
- Loaded at startup
- Used per-request

## Performance Optimization

### Optimizations Implemented

1. **Connection Pooling**: Reuse HTTP connections
2. **Worker Pools**: Parallel task execution
3. **Rate Limiting**: Prevent API throttling
4. **Caching**: Avoid redundant requests
5. **Streaming**: Process large responses incrementally
6. **Timeout Management**: Prevent hanging operations

### Benchmarking

Performance metrics:
- Worker count: 10 (default, configurable)
- Rate limit: 10 req/s (default, configurable)
- Timeout: 30s (default, configurable)
- Cache TTL: 1 hour

## Security Considerations

### Security Features

1. **SSL/TLS Verification**: Enabled by default
2. **Input Validation**: All user input validated
3. **API Key Protection**: Never logged or exposed
4. **Rate Limiting**: Prevents abuse
5. **Timeout Protection**: Prevents DoS
6. **Error Handling**: No sensitive data in errors

### Threat Model

Protected against:
- Command injection
- Path traversal
- DoS attacks
- API key leakage
- Man-in-the-middle attacks

## Testing Strategy

### Test Types

1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Module interaction testing
3. **End-to-End Tests**: Full workflow testing

### Test Coverage

Current coverage:
- Cache: 27.5%
- OSINT modules: 5.0%
- Utils: 22.9%

Target coverage: 80%

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Deployment

### Build Options

```bash
# Development build
go build -o bin/osint ./cmd/osint

# Production build (optimized)
go build -ldflags "-s -w" -o bin/osint ./cmd/osint

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o bin/osint-linux
GOOS=darwin GOARCH=amd64 go build -o bin/osint-darwin
GOOS=windows GOARCH=amd64 go build -o bin/osint-windows.exe
```

### Docker Support (Future)

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags "-s -w" -o osint ./cmd/osint

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/osint /usr/local/bin/
ENTRYPOINT ["osint"]
```

### Configuration Management

1. Default configuration in code
2. System-wide config: `/etc/osint/config.yaml`
3. User config: `~/.osint.yaml`
4. Project config: `./.osint.yaml`
5. Environment variables: `OSINT_*`
6. Command-line flags (highest priority)

## Monitoring and Logging

### Log Levels

- **DEBUG**: Detailed execution information
- **INFO**: General informational messages
- **WARN**: Warning messages (non-fatal)
- **ERROR**: Error messages (recoverable)
- **FATAL**: Fatal errors (unrecoverable)

### Log Format

```
[LEVEL][timestamp] message key=value key=value
```

### Metrics (Future)

Planned metrics:
- Request count by module
- Request duration by module
- Cache hit/miss ratio
- Error rate by type
- API usage by service

## Future Enhancements

### Roadmap

1. **Additional Modules**
   - Social media intelligence (Twitter, LinkedIn, Facebook)
   - Public records search
   - Cryptocurrency analysis
   - Dark web monitoring

2. **Web Interface**
   - React-based UI
   - Real-time results
   - Visualization
   - Report generation

3. **API Server Mode**
   - REST API
   - GraphQL endpoint
   - WebSocket support
   - Rate limiting per client

4. **Database Storage**
   - PostgreSQL integration
   - Result persistence
   - Historical analysis
   - Query interface

5. **Machine Learning**
   - Pattern recognition
   - Anomaly detection
   - Relationship mapping
   - Predictive analysis

6. **Plugin System**
   - Dynamic module loading
   - Custom plugins
   - Third-party extensions
   - Plugin marketplace

## Contributing

### Development Workflow

1. Fork repository
2. Create feature branch
3. Write tests
4. Implement feature
5. Run tests and linters
6. Submit pull request

### Code Standards

- Follow Go conventions
- Write comprehensive tests
- Document public APIs
- Use meaningful variable names
- Keep functions focused and small

### Review Process

All contributions must:
- Pass all tests
- Meet coverage requirements
- Follow code standards
- Include documentation
- Be reviewed by maintainer

## License

MIT License - see LICENSE file for details

## Support

- GitHub Issues: Bug reports and feature requests
- Documentation: See docs/ directory
- Examples: See examples/ directory

---

**Last Updated**: 2026-03-15
**Version**: 1.0.0
**Authors**: EmmmmDeee
