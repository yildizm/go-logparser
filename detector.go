package logparser

import (
	"encoding/json"
	"strings"
)

// detector handles format detection logic
type detector struct{}

// newDetector creates a new format detector
func newDetector() *detector {
	return &detector{}
}

// detectFormat attempts to detect log format from samples
func (d *detector) detectFormat(samples []string) Format {
	if len(samples) == 0 {
		return FormatText // Default to text
	}

	// Take up to 10 samples for detection
	sampleSize := 10
	if len(samples) < sampleSize {
		sampleSize = len(samples)
	}

	// Count successful detections for each format
	scores := make(map[Format]int)

	for _, sample := range samples[:sampleSize] {
		if d.isJSON(sample) {
			scores[FormatJSON]++
		}

		if d.isLogfmt(sample) {
			scores[FormatLogfmt]++
		}
		// Text format always matches as fallback
		scores[FormatText]++
	}

	// Find format with highest score (prefer JSON > logfmt > text)
	if scores[FormatJSON] > scores[FormatLogfmt] && scores[FormatJSON] > scores[FormatText]/2 {
		return FormatJSON
	}

	if scores[FormatLogfmt] > scores[FormatText]/2 {
		return FormatLogfmt
	}

	return FormatText
}

// isJSON checks if a line appears to be JSON
func (d *detector) isJSON(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
		return false
	}

	// Try to parse as JSON
	var obj map[string]interface{}

	return json.Unmarshal([]byte(line), &obj) == nil
}

// isLogfmt checks if a line appears to be logfmt
func (d *detector) isLogfmt(line string) bool {
	// Simple heuristic: contains key=value pattern
	return strings.Contains(line, "=") &&
		(strings.Contains(line, "level=") ||
			strings.Contains(line, "msg=") ||
			strings.Contains(line, "time=") ||
			strings.Contains(line, "timestamp="))
}
