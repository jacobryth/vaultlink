package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Role      string    `json:"role"`
	Keys      []string  `json:"keys"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes audit entries to a file or stdout.
type Logger struct {
	f *os.File
}

// NewLogger creates a Logger. If path is empty, stdout is used.
func NewLogger(path string) (*Logger, error) {
	if path == "" {
		return &Logger{f: os.Stdout}, nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{f: f}, nil
}

// Log writes an audit entry as a JSON line.
func (l *Logger) Log(op, path, role string, keys []string, success bool, msg string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Path:      path,
		Role:      role,
		Keys:      keys,
		Success:   success,
		Message:   msg,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(l.f, string(b))
	return err
}

// Close closes the underlying file if it is not stdout.
func (l *Logger) Close() error {
	if l.f == os.Stdout {
		return nil
	}
	return l.f.Close()
}
