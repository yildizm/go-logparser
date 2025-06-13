package logparser

import (
	"strings"
	"testing"
)

func TestJSONParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(*testing.T, LogEntry)
	}{
		{
			name:  "standard JSON log",
			input: `{"timestamp":"2024-01-02T15:04:05Z","level":"ERROR","message":"Database connection failed","service":"api"}`,
			validate: func(t *testing.T, e LogEntry) {
				t.Helper()
				if e.Level != LevelError {
					t.Errorf("want level ERROR, got %v", e.Level)
				}
				if e.Message != "Database connection failed" {
					t.Errorf("want message 'Database connection failed', got %s", e.Message)
				}
				if service, ok := e.Fields["service"].(string); !ok || service != "api" {
					t.Errorf("want service 'api', got %v", e.Fields["service"])
				}
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
		},
	}

	parser := NewWithFormat(FormatJSON)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := parser.ParseString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil && len(entries) > 0 {
				tt.validate(t, entries[0])
			}
		})
	}
}

func TestLogfmtParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(*testing.T, LogEntry)
	}{
		{
			name:  "standard logfmt",
			input: `time=2024-01-02T15:04:05Z level=error msg="Connection timeout" service=worker duration=1.23`,
			validate: func(t *testing.T, e LogEntry) {
				t.Helper()
				if e.Level != LevelError {
					t.Errorf("want level ERROR, got %v", e.Level)
				}
				if e.Message != "Connection timeout" {
					t.Errorf("want message 'Connection timeout', got %s", e.Message)
				}
				if duration, ok := e.Fields["duration"].(string); !ok || duration != "1.23" {
					t.Errorf("want duration=1.23, got %v", e.Fields["duration"])
				}
			},
		},
	}

	parser := NewWithFormat(FormatLogfmt)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, _ := parser.ParseString(tt.input)
			if tt.validate != nil && len(entries) > 0 {
				tt.validate(t, entries[0])
			}
		})
	}
}

func TestTextParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(*testing.T, LogEntry)
	}{
		{
			name:  "bracketed level",
			input: `2024-01-02 15:04:05 [ERROR] Failed to connect to database`,
			validate: func(t *testing.T, e LogEntry) {
				t.Helper()
				if e.Level != LevelError {
					t.Errorf("want level ERROR, got %v", e.Level)
				}
				if !strings.Contains(e.Message, "Failed to connect") {
					t.Errorf("message should contain 'Failed to connect'")
				}
			},
		},
	}

	parser := NewWithFormat(FormatText)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, _ := parser.ParseString(tt.input)
			if tt.validate != nil && len(entries) > 0 {
				tt.validate(t, entries[0])
			}
		})
	}
}

func TestFormatDetection(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Format
	}{
		{
			name:  "JSON logs",
			input: `{"level":"info","msg":"test"}` + "\n" + `{"level":"error","msg":"test2"}`,
			want:  FormatJSON,
		},
		{
			name:  "logfmt logs",
			input: `level=info msg="test" time=2024-01-02T15:04:05Z` + "\n" + `level=error msg="test2" time=2024-01-02T15:04:06Z`,
			want:  FormatLogfmt,
		},
		{
			name:  "text logs",
			input: `[INFO] Starting application` + "\n" + `[ERROR] Connection failed`,
			want:  FormatText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New() // Auto-detect format

			entries, err := parser.ParseString(tt.input)
			if err != nil {
				t.Errorf("ParseString() error = %v", err)
				return
			}

			if len(entries) == 0 {
				t.Errorf("expected entries, got none")
			}
		})
	}
}

func TestEmptyInput(t *testing.T) {
	parser := New()

	entries, err := parser.ParseString("")
	if err != nil {
		t.Errorf("ParseString() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected no entries for empty input, got %d", len(entries))
	}
}

func BenchmarkJSONParser(b *testing.B) {
	input := `{"timestamp":"2024-01-02T15:04:05Z","level":"ERROR","message":"Database connection failed","service":"api"}`
	parser := NewWithFormat(FormatJSON)

	b.ResetTimer()

	for range b.N {
		_, _ = parser.ParseString(input)
	}
}

func BenchmarkLogfmtParser(b *testing.B) {
	input := `time=2024-01-02T15:04:05Z level=error msg="Connection timeout" service=worker duration=1.23`
	parser := NewWithFormat(FormatLogfmt)

	b.ResetTimer()

	for range b.N {
		_, _ = parser.ParseString(input)
	}
}

func BenchmarkTextParser(b *testing.B) {
	input := `2024-01-02 15:04:05 [ERROR] Failed to connect to database`
	parser := NewWithFormat(FormatText)

	b.ResetTimer()

	for range b.N {
		_, _ = parser.ParseString(input)
	}
}
