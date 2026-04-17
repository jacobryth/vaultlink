package rotate

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LogEntry is a JSON-serialisable rotation log entry.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Key       string `json:"key"`
	PrevHash  string `json:"prev_hash"`
	NewHash   string `json:"new_hash"`
}

// WriteLog appends rotation records to the given file path as JSON lines.
// If path is empty, entries are written to stdout.
func WriteLog(path string, records []Record) error {
	if len(records) == 0 {
		return nil
	}

	var f *os.File
	var err error
	if path == "" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("rotate log: open file: %w", err)
		}
		defer f.Close()
	}

	enc := json.NewEncoder(f)
	for _, rec := range records {
		entry := LogEntry{
			Timestamp: rec.RotatedAt.Format(time.RFC3339),
			Key:       rec.Key,
			PrevHash:  rec.PrevHash,
			NewHash:   rec.NewHash,
		}
		if err := enc.Encode(entry); err != nil {
			return fmt.Errorf("rotate log: encode entry: %w", err)
		}
	}
	return nil
}
