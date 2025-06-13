package logparser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Parser is the main interface for log parsing
type Parser interface {
	Parse(r io.Reader) ([]LogEntry, error)
	ParseString(s string) ([]LogEntry, error)
}

// parser implements the Parser interface
type parser struct {
	format   Format
	detector *detector
}

// New creates a parser with auto-detection
func New() Parser {
	return &parser{
		format:   FormatAuto,
		detector: newDetector(),
	}
}

// NewWithFormat creates a parser for specific format
func NewWithFormat(format Format) Parser {
	return &parser{
		format:   format,
		detector: newDetector(),
	}
}

// Parse parses logs from a reader
func (p *parser) Parse(r io.Reader) ([]LogEntry, error) {
	var lines []string

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, BufferSize), BufferSize) // 1MB buffer

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return p.parseLines(lines)
}

// ParseString parses a single log string
func (p *parser) ParseString(s string) ([]LogEntry, error) {
	lines := strings.Split(s, "\n")

	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return p.parseLines(cleanLines)
}

// parseLines parses an array of log lines
func (p *parser) parseLines(lines []string) ([]LogEntry, error) {
	if len(lines) == 0 {
		return []LogEntry{}, nil
	}

	format := p.format

	// Auto-detect format if needed
	if format == FormatAuto {
		format = p.detector.detectFormat(lines)
	}

	// Parse using the appropriate format parser
	switch format {
	case FormatAuto:
		// This should not happen as FormatAuto is handled above
		return nil, errors.New("auto-detection failed")
	case FormatJSON:
		return parseJSON(lines)
	case FormatLogfmt:
		return parseLogfmt(lines)
	case FormatText:
		return parseText(lines)
	default:
		return parseText(lines) // Default fallback
	}
}
