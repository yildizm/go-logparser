# go-logparser

[![Go Reference](https://pkg.go.dev/badge/github.com/yildizm/go-logparser.svg)](https://pkg.go.dev/github.com/yildizm/go-logparser)
[![CI](https://github.com/yildizm/go-logparser/workflows/CI/badge.svg)](https://github.com/yildizm/go-logparser/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/yildizm/go-logparser)](https://goreportcard.com/report/github.com/yildizm/go-logparser)

A fast, standalone Go library for parsing multi-format logs (JSON, logfmt, plain text) with automatic format detection.

## Features

- **Multi-format support**: JSON, logfmt, and plain text logs
- **Automatic format detection**: Intelligently detects log format from samples
- **High performance**: Optimized for processing large log files
- **Simple API**: Easy-to-use interface with sensible defaults
- **Zero dependencies**: Pure Go implementation
- **Comprehensive testing**: Thoroughly tested with benchmarks

## Installation

```bash
go get github.com/yildizm/go-logparser
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/yildizm/go-logparser"
)

func main() {
    // Create parser with auto-detection
    parser := logparser.New()
    
    // Parse log string
    logs := `{"timestamp":"2024-01-02T15:04:05Z","level":"ERROR","message":"Database connection failed"}
level=info msg="Request processed" duration=123ms`
    
    entries, err := parser.ParseString(logs)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, entry := range entries {
        fmt.Printf("%s [%s] %s\n", 
            entry.Timestamp.Format("15:04:05"), 
            entry.Level, 
            entry.Message)
    }
}
```

## API Reference

### Types

Core data structures used throughout the library for parsing and representing log entries.

```go
// Parser is the main interface for log parsing
type Parser interface {
    Parse(r io.Reader) ([]LogEntry, error)
    ParseString(s string) ([]LogEntry, error)
}

// LogEntry represents a parsed log entry
type LogEntry struct {
    Timestamp time.Time              `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Format represents log format types
type Format int

const (
    FormatAuto Format = iota
    FormatJSON
    FormatLogfmt
    FormatText
)
```

### Functions

Main entry points for creating parsers with different configuration options.

```go
// New creates a parser with auto-detection
func New() Parser

// NewWithFormat creates a parser for specific format
func NewWithFormat(format Format) Parser
```

## Supported Log Formats

### JSON Logs
Structured logs in JSON format with key-value pairs for easy machine parsing.
```json
{"timestamp":"2024-01-02T15:04:05Z","level":"ERROR","message":"Database connection failed","service":"api"}
```

### Logfmt Logs
Key-value structured logs popular in cloud-native applications for human-readable output.
```
time=2024-01-02T15:04:05Z level=error msg="Connection timeout" service=worker duration=1.23
```

### Plain Text Logs
Traditional unstructured log formats with various timestamp and message patterns.
```
2024-01-02 15:04:05 [ERROR] Failed to connect to database
[INFO] Application started successfully
Jan 02 15:04:05 hostname process[pid]: System event occurred
```

## Examples

### Auto-Detection
Automatically detect log format from the first few lines of input.
```go
parser := logparser.New()
entries, err := parser.ParseString(logs)
```

### Specific Format
Create parsers optimized for known log formats to improve performance.
```go
// JSON parser
jsonParser := logparser.NewWithFormat(logparser.FormatJSON)
entries, err := jsonParser.ParseString(jsonLogs)

// Logfmt parser
logfmtParser := logparser.NewWithFormat(logparser.FormatLogfmt)
entries, err := logfmtParser.ParseString(logfmtLogs)

// Text parser
textParser := logparser.NewWithFormat(logparser.FormatText)
entries, err := textParser.ParseString(textLogs)
```

### Parse from Reader
Stream parse logs directly from files or other io.Reader sources.
```go
file, err := os.Open("app.log")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

parser := logparser.New()
entries, err := parser.Parse(file)
```

## Field Extraction

The library automatically extracts common fields from log entries:

### Standard Fields
Commonly used log fields that are automatically extracted and mapped to the LogEntry struct.
- **Timestamp**: `timestamp`, `time`, `@timestamp`, `ts`
- **Level**: `level`, `severity`, `log.level`
- **Message**: `message`, `msg`, `log`

### Additional Fields
Custom fields not mapped to standard fields are preserved for application-specific processing.
All other fields are preserved in the `Fields` map with their original types.

## Performance

Benchmarks on a modern machine:

```
BenchmarkJSONParser-8     1000000   1043 ns/op   512 B/op   12 allocs/op
BenchmarkLogfmtParser-8    800000   1387 ns/op   648 B/op   15 allocs/op
BenchmarkTextParser-8      600000   2156 ns/op   712 B/op   18 allocs/op
```

## Error Handling

The library is designed to be resilient:
- Invalid lines are skipped but processing continues
- Malformed timestamps default to current time
- Unknown levels default to "INFO"
- Parse errors are reported but don't stop processing

## Testing

Comprehensive test suite covering all parsers, edge cases, and performance benchmarks.

Run the test suite:

```bash
cd go-logparser
go test -v
```

Run benchmarks:

```bash
go test -bench=.
```

## Examples

See the [examples/](examples/) directory for complete working examples:

- [Basic Usage](examples/basic/main.go) - Demonstrates all parser formats

Run the basic example:

```bash
cd examples/basic
go run main.go
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Changelog

### v1.0.0
- Initial release
- Support for JSON, logfmt, and plain text formats
- Automatic format detection
- Comprehensive test suite
- Performance benchmarks