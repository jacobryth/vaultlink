package env

import (
	"fmt"
	"os"
	"strings"
)

// PatchMode controls how the patcher handles existing keys.
type PatchMode string

const (
	PatchModeUpsert PatchMode = "upsert" // add new, update existing
	PatchModeAddOnly PatchMode = "add"   // only add missing keys
	PatchModeRemove  PatchMode = "remove" // remove specified keys
)

var validPatchModes = map[PatchMode]bool{
	PatchModeUpsert:  true,
	PatchModeAddOnly: true,
	PatchModeRemove:  true,
}

// Patcher applies targeted key-level patches to an existing .env file.
type Patcher struct {
	mode PatchMode
}

// NewPatcher creates a Patcher with the given mode.
func NewPatcher(mode PatchMode) (*Patcher, error) {
	if !validPatchModes[mode] {
		return nil, fmt.Errorf("env/patcher: unknown mode %q", mode)
	}
	return &Patcher{mode: mode}, nil
}

// Patch reads the file at path, applies the patch map according to the mode,
// and writes the result back to the same file.
func (p *Patcher) Patch(path string, patch map[string]string) error {
	existing, err := readEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("env/patcher: read %s: %w", path, err)
	}
	if existing == nil {
		existing = map[string]string{}
	}

	switch p.mode {
	case PatchModeUpsert:
		for k, v := range patch {
			existing[k] = v
		}
	case PatchModeAddOnly:
		for k, v := range patch {
			if _, ok := existing[k]; !ok {
				existing[k] = v
			}
		}
	case PatchModeRemove:
		for k := range patch {
			delete(existing, k)
		}
	}

	return writeEnvFile(path, existing)
}

// readEnvFile parses a simple KEY=VALUE .env file into a map.
func readEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
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

// writeEnvFile serialises a map back to KEY=VALUE lines.
func writeEnvFile(path string, m map[string]string) error {
	var sb strings.Builder
	for k, v := range m {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(v)
		sb.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(sb.String()), 0o600)
}
