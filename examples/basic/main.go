package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/yildizm/go-logparser"
)

func main() {
	// Example log data in different formats
	jsonLogs := `{"timestamp":"2024-01-02T15:04:05Z","level":"ERROR","message":"Database connection failed","service":"api"}
{"timestamp":"2024-01-02T15:04:06Z","level":"INFO","message":"Request processed","service":"api","duration":"123ms"}`

	logfmtLogs := `time=2024-01-02T15:04:05Z level=error msg="Connection timeout" service=worker duration=1.23
time=2024-01-02T15:04:06Z level=info msg="Task completed" service=worker task_id=abc123`

	textLogs := `2024-01-02 15:04:05 [ERROR] Failed to connect to database
2024-01-02 15:04:06 [INFO] Application started successfully`

	fmt.Println("=== go-logparser Example ===")
	fmt.Println()

	// Parse JSON logs with auto-detection
	fmt.Println("1. Auto-detection (JSON logs):")

	parser := logparser.New()

	entries, err := parser.ParseString(jsonLogs)
	if err != nil {
		log.Fatal(err)
	}

	printEntries(entries)

	// Parse logfmt logs with specific format
	fmt.Println("2. Specific format (logfmt):")

	logfmtParser := logparser.NewWithFormat(logparser.FormatLogfmt)

	entries, err = logfmtParser.ParseString(logfmtLogs)
	if err != nil {
		log.Fatal(err)
	}

	printEntries(entries)

	// Parse text logs with specific format
	fmt.Println("3. Specific format (text):")

	textParser := logparser.NewWithFormat(logparser.FormatText)

	entries, err = textParser.ParseString(textLogs)
	if err != nil {
		log.Fatal(err)
	}

	printEntries(entries)

	// Parse from a reader (using strings.Reader as example)
	fmt.Println("4. Parse from reader:")

	reader := strings.NewReader(jsonLogs)

	entries, err = parser.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}

	printEntries(entries)
}

func printEntries(entries []logparser.LogEntry) {
	for i, entry := range entries {
		fmt.Printf("  Entry %d:\n", i+1)
		fmt.Printf("    Timestamp: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("    Level: %s\n", entry.Level)
		fmt.Printf("    Message: %s\n", entry.Message)

		if len(entry.Fields) > 0 {
			fmt.Printf("    Fields: %+v\n", entry.Fields)
		}

		fmt.Println()
	}
}
