package env

import (
	"fmt"
	"os"
	"strings"
)

// MergeMode controls how conflicts are resolved when merging env layers.
type MergeMode string

const (
	MergeModeOverwrite MergeMode = "overwrite"
	MergeModeKeep      MergeMode = "keep"
	MergeModeError     MergeMode = "error"
)

var validMergeModes = map[MergeMode]bool{
	MergeModeOverwrite: true,
	MergeModeKeep:      true,
	MergeModeError:     true,
}

// EnvMerger merges multiple env-file layers into a single key-value map.
type EnvMerger struct {
	mode   MergeMode
	loader *Loader
}

// NewEnvMerger returns an EnvMerger for the given merge mode.
func NewEnvMerger(mode MergeMode) (*EnvMerger, error) {
	if !validMergeModes[mode] {
		return nil, fmt.Errorf("env/merger: unknown merge mode %q", mode)
	}
	return &EnvMerger{mode: mode, loader: NewLoader()}, nil
}

// MergeFiles reads each path in order and merges them according to the mode.
func (m *EnvMerger) MergeFiles(paths []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("env/merger: reading %s: %w", p, err)
		}
		pairs, err := m.loader.LoadReader(strings.NewReader(string(data)))
		if err != nil {
			return nil, fmt.Errorf("env/merger: parsing %s: %w", p, err)
		}
		for k, v := range pairs {
			switch m.mode {
			case MergeModeOverwrite:
				result[k] = v
			case MergeModeKeep:
				if _, exists := result[k]; !exists {
					result[k] = v
				}
			case MergeModeError:
				if _, exists := result[k]; exists {
					return nil, fmt.Errorf("env/merger: duplicate key %q in %s", k, p)
				}
				result[k] = v
			}
		}
	}
	return result, nil
}
