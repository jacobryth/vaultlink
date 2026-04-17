package export

import (
	"encoding/json"
	"fmt"
	"os"\n	"path/filepath"
	"strings"
)

type Format string

const (
	FormatEnv  Format = "env"
	FormatJSON Format = "json"
)

type Exporter struct {
	format Format
}

func New(format string) (*Exporter, error) {
	f := Format(strings.ToLower(format))
	if f != FormatEnv && f != FormatJSON {
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
	return &Exporter{format: f}, nil
}

func (e *Exporter) Write(secrets map[string]string, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("export: create dirs: %w", err)
	}
	switch e.format {
	case FormatJSON:
		return e.writeJSON(secrets, dest)
	default:
		return e.writeEnv(secrets, dest)
	}
}

func (e *Exporter) writeEnv(secrets map[string]string, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("export: open file: %w", err)
	}
	defer f.Close()
	for k, v := range secrets {
		if _, err := fmt.Fprintf(f, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(secrets map[string]string, dest string) error {
	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return fmt.Errorf("export: marshal json: %w", err)
	}
	return os.WriteFile(dest, data, 0600)
}
