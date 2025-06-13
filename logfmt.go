package logparser

import (
	"strings"
	"time"
)

// parseLogfmt parses logfmt formatted logs
func parseLogfmt(lines []string) ([]LogEntry, error) {
	entries := make([]LogEntry, 0, len(lines))

	for _, line := range lines {
		entry, err := parseLogfmtLine(line)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

// parseLogfmtLine parses a single logfmt line
func parseLogfmtLine(line string) (*LogEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, ErrEmptyLine
	}

	entry := &LogEntry{
		Fields: make(map[string]interface{}),
		Level:  "INFO", // Default level
	}

	// Parse key=value pairs
	pairs := parseLogfmtPairs(line)

	// Extract standard fields
	extractLogfmtTimestamp(pairs, entry)
	extractLogfmtLevel(pairs, entry)
	extractLogfmtMessage(pairs, entry)

	// Remaining pairs go to Fields
	for k, v := range pairs {
		entry.Fields[k] = v
	}

	// Default timestamp if not found
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	return entry, nil
}

// parseLogfmtPairs parses key=value pairs from a line
func parseLogfmtPairs(line string) map[string]interface{} {
	pairs := make(map[string]interface{})

	var key string

	var value strings.Builder

	inQuotes := false
	inKey := true

	for i := range len(line) {
		ch := line[i]

		switch {
		case ch == '=' && inKey && !inQuotes:
			inKey = false

		case ch == '"' && !inKey:
			if i > 0 && line[i-1] != '\\' {
				inQuotes = !inQuotes
			} else {
				value.WriteByte(ch)
			}

		case ch == ' ' && !inQuotes && !inKey:
			// End of value
			if key != "" {
				pairs[key] = value.String()
			}

			key = ""

			value.Reset()

			inKey = true

		case inKey:
			key += string(ch)

		default:
			value.WriteByte(ch)
		}
	}

	// Handle last pair
	if key != "" {
		pairs[key] = value.String()
	}

	return pairs
}

// extractLogfmtTimestamp extracts timestamp from logfmt pairs
func extractLogfmtTimestamp(pairs map[string]interface{}, entry *LogEntry) {
	for _, key := range []string{"timestamp", "time", "ts"} {
		if val, ok := pairs[key]; ok {
			if t, err := parseTimestamp(val); err == nil {
				entry.Timestamp = t

				delete(pairs, key)

				return
			}
		}
	}
}

// extractLogfmtLevel extracts log level from logfmt pairs
func extractLogfmtLevel(pairs map[string]interface{}, entry *LogEntry) {
	if val, ok := pairs["level"]; ok {
		if s, ok := val.(string); ok {
			entry.Level = ParseLevel(s)

			delete(pairs, "level")
		}
	}
}

// extractLogfmtMessage extracts message from logfmt pairs
func extractLogfmtMessage(pairs map[string]interface{}, entry *LogEntry) {
	for _, key := range []string{"msg", "message"} {
		if val, ok := pairs[key]; ok {
			if s, ok := val.(string); ok {
				entry.Message = s

				delete(pairs, key)

				return
			}
		}
	}
}
