package env

import (
	"fmt"
	"strings"
)

// SplitMode controls how the splitter partitions secrets.
type SplitMode string

const (
	SplitNone   SplitMode = "none"
	SplitPrefix SplitMode = "prefix"
	SplitAlpha  SplitMode = "alpha"
)

var splitModes = map[SplitMode]bool{
	SplitNone:   true,
	SplitPrefix: true,
	SplitAlpha:  true,
}

// Splitter partitions a secrets map into named buckets.
type Splitter struct {
	mode     SplitMode
	delim    string
}

// New returns a Splitter for the given mode. For SplitPrefix, delim is the
// separator used to extract the prefix (e.g. "_").
func NewSplitter(mode SplitMode, delim string) (*Splitter, error) {
	if !splitModes[mode] {
		return nil, fmt.Errorf("env/splitter: unknown mode %q", mode)
	}
	if mode == SplitPrefix && delim == "" {
		return nil, fmt.Errorf("env/splitter: prefix mode requires a non-empty delimiter")
	}
	return &Splitter{mode: mode, delim: delim}, nil
}

// Split partitions secrets into buckets. Returns a map of bucket name →
// key/value pairs. With SplitNone all secrets land in bucket "default".
func (s *Splitter) Split(secrets map[string]string) map[string]map[string]string {
	result := make(map[string]map[string]string)

	if secrets == nil {
		return result
	}

	for k, v := range secrets {
		bucket := s.bucket(k)
		if result[bucket] == nil {
			result[bucket] = make(map[string]string)
		}
		result[bucket][k] = v
	}
	return result
}

func (s *Splitter) bucket(key string) string {
	switch s.mode {
	case SplitPrefix:
		if idx := strings.Index(key, s.delim); idx > 0 {
			return strings.ToLower(key[:idx])
		}
		return "default"
	case SplitAlpha:
		if len(key) > 0 {
			return strings.ToLower(string(key[0]))
		}
		return "default"
	default:
		return "default"
	}
}
