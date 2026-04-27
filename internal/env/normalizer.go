package env

import (
	"fmt"
	"regexp"
	"strings"
)

// NormalizeLevel controls how env key normalization is applied.
type NormalizeLevel string

const (
	NormalizeLevelNone    NormalizeLevel = "none"
	NormalizeLevelKeys    NormalizeLevel = "keys"
	NormalizeLevelValues  NormalizeLevel = "values"
	NormalizeLevelBoth    NormalizeLevel = "both"
)

var validNormalizeLevels = map[NormalizeLevel]bool{
	NormalizeLevelNone:   true,
	NormalizeLevelKeys:   true,
	NormalizeLevelValues: true,
	NormalizeLevelBoth:   true,
}

var invalidKeyChars = regexp.MustCompile(`[^A-Z0-9_]`)

// Normalizer sanitizes env key names and/or values into a canonical form.
type Normalizer struct {
	level NormalizeLevel
}

// NewNormalizer creates a Normalizer for the given level.
func NewNormalizer(level NormalizeLevel) (*Normalizer, error) {
	if !validNormalizeLevels[level] {
		return nil, fmt.Errorf("normalizer: unknown level %q", level)
	}
	return &Normalizer{level: level}, nil
}

// Apply returns a new map with keys and/or values normalized according to the level.
func (n *Normalizer) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if n.level == NormalizeLevelNone {
		out := make(map[string]string, len(secrets))
		for k, v := range secrets {
			out[k] = v
		}
		return out
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		nk := k
		nv := v
		if n.level == NormalizeLevelKeys || n.level == NormalizeLevelBoth {
			nk = normalizeKey(k)
		}
		if n.level == NormalizeLevelValues || n.level == NormalizeLevelBoth {
			nv = normalizeValue(v)
		}
		out[nk] = nv
	}
	return out
}

// normalizeKey uppercases the key and replaces invalid characters with underscores.
func normalizeKey(k string) string {
	upper := strings.ToUpper(k)
	return invalidKeyChars.ReplaceAllString(upper, "_")
}

// normalizeValue trims surrounding whitespace from the value.
func normalizeValue(v string) string {
	return strings.TrimSpace(v)
}
