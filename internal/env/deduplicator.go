package env

import "fmt"

// DedupeMode controls how duplicate keys in an env file are resolved.
type DedupeMode string

const (
	DedupeModeNone  DedupeMode = "none"
	DedupeModeFirst DedupeMode = "first"
	DedupeModeLast  DedupeMode = "last"
	DedupeModeError DedupeMode = "error"
)

var validDedupeModes = map[DedupeMode]bool{
	DedupeModeNone:  true,
	DedupeModeFirst: true,
	DedupeModeLast:  true,
	DedupeModeError: true,
}

// EnvDeduplicator removes or resolves duplicate keys from a slice of key-value pairs.
type EnvDeduplicator struct {
	mode DedupeMode
}

// NewEnvDeduplicator creates a new EnvDeduplicator with the given mode.
func NewEnvDeduplicator(mode DedupeMode) (*EnvDeduplicator, error) {
	if !validDedupeModes[mode] {
		return nil, fmt.Errorf("env/deduplicator: unknown mode %q", mode)
	}
	return &EnvDeduplicator{mode: mode}, nil
}

// Apply deduplicates the provided key-value pairs according to the configured mode.
// Input is a slice of [2]string{key, value} pairs to preserve ordering.
func (d *EnvDeduplicator) Apply(pairs [][2]string) ([][2]string, error) {
	if pairs == nil {
		return nil, nil
	}
	if d.mode == DedupeModeNone {
		return pairs, nil
	}

	seen := make(map[string]int) // key -> index in result
	result := make([][2]string, 0, len(pairs))

	for _, pair := range pairs {
		k := pair[0]
		if idx, exists := seen[k]; exists {
			switch d.mode {
			case DedupeModeError:
				return nil, fmt.Errorf("env/deduplicator: duplicate key %q", k)
			case DedupeModeFirst:
				// keep the existing entry, discard new
				_ = idx
				continue
			case DedupeModeLast:
				// overwrite the existing entry in-place
				result[idx] = pair
				continue
			}
		} else {
			seen[k] = len(result)
			result = append(result, pair)
		}
	}
	return result, nil
}
