package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	FilePath string
	Overwrite bool
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string, overwrite bool) *Writer {
	return &Writer{
		FilePath: filePath,
		Overwrite: overwrite,
	}
}

// Write writes the provided secrets map to the .env file.
// If Overwrite is false and the file exists, it merges without overwriting existing keys.
func (w *Writer) Write(secrets map[string]string) error {
	existing := map[string]string{}

	if !w.Overwrite {
		var err error
		existing, err = readExisting(w.FilePath)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("reading existing env file: %w", err)
		}
	}

	merged := make(map[string]string)
	for k, v := range secrets {
		merged[k] = v
	}
	for k, v := range existing {
		if _, ok := merged[k]; !ok {
			merged[k] = v
		}
	}

	f, err := os.Create(w.FilePath)
	if err != nil {
		return fmt.Errorf("creating env file: %w", err)
	}
	defer f.Close()

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		_, err := fmt.Fprintf(f, "%s=%s\n", k, merged[k])
		if err != nil {
			return fmt.Errorf("writing key %s: %w", k, err)
		}
	}
	return nil
}

func readExisting(filePath string) (map[string]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, nil
}
