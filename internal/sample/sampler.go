package sample

import (
	"fmt"
	"math/rand"
)

// Level controls the sampling strategy.
type Level string

const (
	LevelNone   Level = "none"
	LevelRandom Level = "random"
	LevelNth    Level = "nth"
)

var validLevels = map[Level]bool{
	LevelNone:   true,
	LevelRandom: true,
	LevelNth:    true,
}

// Sampler reduces a secret map to a subset based on the chosen strategy.
type Sampler struct {
	level Level
	n     int
}

// New creates a Sampler. For LevelRandom and LevelNth, n must be >= 1.
func New(level Level, n int) (*Sampler, error) {
	if !validLevels[level] {
		return nil, fmt.Errorf("sample: unknown level %q", level)
	}
	if level != LevelNone && n < 1 {
		return nil, fmt.Errorf("sample: n must be >= 1 for level %q", level)
	}
	return &Sampler{level: level, n: n}, nil
}

// Apply returns a sampled subset of secrets according to the configured strategy.
func (s *Sampler) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	switch s.level {
	case LevelRandom:
		return s.applyRandom(secrets)
	case LevelNth:
		return s.applyNth(secrets)
	default:
		return secrets
	}
}

func (s *Sampler) applyRandom(secrets map[string]string) map[string]string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	limit := s.n
	if limit > len(keys) {
		limit = len(keys)
	}
	out := make(map[string]string, limit)
	for _, k := range keys[:limit] {
		out[k] = secrets[k]
	}
	return out
}

func (s *Sampler) applyNth(secrets map[string]string) map[string]string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	out := make(map[string]string)
	for i, k := range keys {
		if i%s.n == 0 {
			out[k] = secrets[k]
		}
	}
	return out
}
