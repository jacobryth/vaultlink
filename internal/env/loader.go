package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Loader reads an existing .env file into a key-value map.
type Loader struct {
	path string
}

// NewLoader creates a Loader for the given file path.
func NewLoader(path string) *Loader {
	return &Loader{path: path}
}

// Load reads the .env file and returns a map of key-value pairs.
// Lines beginning with '#' and blank lines are ignored.
// Inline comments (# ...) are stripped from values.
func (l *Loader) Load() (map[string]string, error) {
	f, err := os.Open(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("env loader: open %q: %w", l.path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("env loader: line %d: missing '=' in %q", lineNo, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip surrounding quotes if present.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		// Strip inline comments (unquoted).
		if ci := strings.Index(val, " #"); ci >= 0 {
			val = strings.TrimSpace(val[:ci])
		}

		result[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env loader: scan %q: %w", l.path, err)
	}

	return result, nil
}

// Keys returns all keys present in the .env file in the order they appear.
func (l *Loader) Keys() ([]string, error) {
	secrets, err := l.Load()
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	return keys, nil
}
