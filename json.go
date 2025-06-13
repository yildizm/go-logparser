package logparser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// parseJSON parses JSON formatted logs
func parseJSON(lines []string) ([]LogEntry, error) {
	entries := make([]LogEntry, 0, len(lines))

	for _, line := range lines {
		entry, err := parseJSONLine(line)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

// parseJSONLine parses a single JSON log line
func parseJSONLine(line string) (*LogEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, ErrEmptyLine
	}

	// Parse JSON
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	entry := &LogEntry{
		Fields: make(map[string]interface{}),
	}

	// Extract standard fields
	extractJSONTimestamp(raw, entry)
	extractJSONLevel(raw, entry)
	extractJSONMessage(raw, entry)

	// Remaining fields go to Fields map
	for k, v := range raw {
		entry.Fields[k] = v
	}

	return entry, nil
}

// extractJSONTimestamp extracts timestamp from various field names
func extractJSONTimestamp(raw map[string]interface{}, entry *LogEntry) {
	for _, key := range []string{"timestamp", "time", "@timestamp", "ts"} {
		if val, ok := raw[key]; ok {
			if t, err := parseTimestamp(val); err == nil {
				entry.Timestamp = t

				delete(raw, key)

				return
			}
		}
	}
	// Default to current time if no timestamp found
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
}

// extractJSONLevel extracts log level from various field names
func extractJSONLevel(raw map[string]interface{}, entry *LogEntry) {
	for _, key := range []string{"level", "severity", "log.level"} {
		if val, ok := raw[key]; ok {
			if s, ok := val.(string); ok {
				entry.Level = ParseLevel(s)

				delete(raw, key)

				return
			}
		}
	}
	// Default level
	entry.Level = LevelInfo
}

// extractJSONMessage extracts message from various field names
func extractJSONMessage(raw map[string]interface{}, entry *LogEntry) {
	for _, key := range []string{"message", "msg", "log"} {
		if val, ok := raw[key]; ok {
			if s, ok := val.(string); ok {
				entry.Message = s

				delete(raw, key)

				return
			}
		}
	}
}
