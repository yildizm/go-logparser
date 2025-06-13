package logparser

import (
	"regexp"
	"strings"
	"time"
)

// textPattern represents a text log pattern
type textPattern struct {
	regex    *regexp.Regexp
	tsFormat string
	tsIndex  int
	lvlIndex int
	msgIndex int
}

// parseText parses plain text logs with common patterns
func parseText(lines []string) ([]LogEntry, error) {
	patterns := initTextPatterns()
	entries := make([]LogEntry, 0, len(lines))

	for _, line := range lines {
		entry, err := parseTextLine(line, patterns)
		if err != nil {
			return nil, err
		}

		entries = append(entries, *entry)
	}

	return entries, nil
}

// parseTextLine parses a single text log line
func parseTextLine(line string, patterns []*textPattern) (*LogEntry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, ErrEmptyLine
	}

	entry := &LogEntry{
		Message: line,   // Default to full line
		Level:   "INFO", // Default level
		Fields:  make(map[string]interface{}),
	}

	// Try each pattern
	for _, pattern := range patterns {
		matches := pattern.regex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		// Extract timestamp
		if pattern.tsIndex > 0 && pattern.tsIndex < len(matches) && pattern.tsFormat != "" {
			if t, err := time.Parse(pattern.tsFormat, matches[pattern.tsIndex]); err == nil {
				entry.Timestamp = t
			}
		}

		// Extract level
		if pattern.lvlIndex > 0 && pattern.lvlIndex < len(matches) {
			entry.Level = ParseLevel(matches[pattern.lvlIndex])
		}

		// Extract message
		if pattern.msgIndex > 0 && pattern.msgIndex < len(matches) {
			entry.Message = matches[pattern.msgIndex]
		}

		break // Use first matching pattern
	}

	// If no timestamp found, use current time
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	return entry, nil
}

// initTextPatterns initializes common log patterns
func initTextPatterns() []*textPattern {
	patterns := []struct {
		pattern  string
		tsFormat string
		tsIndex  int
		lvlIndex int
		msgIndex int
	}{
		// Syslog format: Jan 02 15:04:05 hostname process[pid]: message
		{
			pattern:  `^(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+\S+\s+\S+:\s+\[?(\w+)\]?\s+(.*)$`,
			tsFormat: "Jan 02 15:04:05",
			tsIndex:  1,
			lvlIndex: LevelIndex,
			msgIndex: MessageIndex,
		},
		// Common format: 2006-01-02 15:04:05 [LEVEL] message
		{
			pattern:  `^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2})\s+\[(\w+)\]\s+(.*)$`,
			tsFormat: "2006-01-02 15:04:05",
			tsIndex:  1,
			lvlIndex: LevelIndex,
			msgIndex: MessageIndex,
		},
		// ISO format: 2006-01-02T15:04:05.000Z [LEVEL] message
		{
			pattern:  `^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z?)\s+\[?(\w+)\]?\s+(.*)$`,
			tsFormat: time.RFC3339,
			tsIndex:  1,
			lvlIndex: LevelIndex,
			msgIndex: MessageIndex,
		},
		// Simple format: [LEVEL] message
		{
			pattern:  `^\[(\w+)\]\s+(.*)$`,
			tsFormat: "",
			tsIndex:  0,
			lvlIndex: 1,
			msgIndex: MessageIndexAlt,
		},
	}

	textPatterns := make([]*textPattern, 0, len(patterns))

	for _, pt := range patterns {
		re, err := regexp.Compile(pt.pattern)
		if err != nil {
			continue
		}

		textPatterns = append(textPatterns, &textPattern{
			regex:    re,
			tsFormat: pt.tsFormat,
			tsIndex:  pt.tsIndex,
			lvlIndex: pt.lvlIndex,
			msgIndex: pt.msgIndex,
		})
	}

	return textPatterns
}
