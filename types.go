package logparser

import (
	"errors"
	"strings"
	"time"
)

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

// Static errors
var (
	ErrEmptyLine = errors.New("empty line")
)

// Log level constants
const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"
)

// Buffer and pattern constants
const (
	BufferSize      = 1024 * 1024 // 1MB buffer
	LevelIndex      = 2
	MessageIndex    = 3
	MessageIndexAlt = 2 // Alternative message index for some patterns
)

// String returns the string representation of the format
func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatLogfmt:
		return "logfmt"
	case FormatText:
		return "text"
	case FormatAuto:
		return "auto"
	default:
		return "unknown"
	}
}

// ParseLevel parses string to standard level
func ParseLevel(s string) string {
	switch strings.ToUpper(s) {
	case "DEBUG", "DBG":
		return "DEBUG"
	case LevelInfo, "INF":
		return "INFO"
	case "WARN", "WARNING", "WRN":
		return "WARN"
	case "ERROR", "ERR":
		return "ERROR"
	case "FATAL", "FTL":
		return "FATAL"
	default:
		return "INFO"
	}
}

// parseTimestamp attempts to parse various timestamp formats
func parseTimestamp(val interface{}) (time.Time, error) {
	switch v := val.(type) {
	case string:
		// Try common formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05.000Z",
			"2006-01-02 15:04:05",
			"Jan 02 15:04:05",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, nil
			}
		}

		return time.Time{}, &ParseError{Type: "timestamp", Value: v, Err: "unknown time format"}
	case float64:
		// Unix timestamp
		return time.Unix(int64(v), 0), nil
	default:
		return time.Time{}, &ParseError{Type: "timestamp", Value: val, Err: "unsupported timestamp type"}
	}
}

// ParseError represents a parsing error
type ParseError struct {
	Type  string
	Value interface{}
	Err   string
}

func (e *ParseError) Error() string {
	return e.Err
}
