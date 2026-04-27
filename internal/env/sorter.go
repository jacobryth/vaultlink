package env

import (
	"fmt"
	"sort"
)

// SortMode controls how .env file keys are sorted.
type SortMode string

const (
	SortNone SortMode = "none"
	SortAsc  SortMode = "asc"
	SortDesc SortMode = "desc"
)

var validSortModes = map[SortMode]bool{
	SortNone: true,
	SortAsc:  true,
	SortDesc: true,
}

// EnvSorter orders key=value pairs in an env file.
type EnvSorter struct {
	mode SortMode
}

// NewEnvSorter creates an EnvSorter for the given mode.
func NewEnvSorter(mode string) (*EnvSorter, error) {
	m := SortMode(mode)
	if !validSortModes[m] {
		return nil, fmt.Errorf("env/sorter: unknown mode %q", mode)
	}
	return &EnvSorter{mode: m}, nil
}

// Apply sorts the provided key=value pairs according to the configured mode.
// Pairs with no '=' separator are passed through unchanged at their original
// position (comments, blanks).
func (s *EnvSorter) Apply(pairs []string) []string {
	if pairs == nil {
		return nil
	}
	if s.mode == SortNone {
		return pairs
	}

	type entry struct {
		key  string
		line string
	}

	var kvs []entry
	var passthrough []string

	for _, line := range pairs {
		idx := indexByte(line, '=')
		if idx < 0 {
			passthrough = append(passthrough, line)
			continue
		}
		kvs = append(kvs, entry{key: line[:idx], line: line})
	}

	sort.Slice(kvs, func(i, j int) bool {
		if s.mode == SortDesc {
			return kvs[i].key > kvs[j].key
		}
		return kvs[i].key < kvs[j].key
	})

	out := make([]string, 0, len(pairs))
	for _, e := range kvs {
		out = append(out, e.line)
	}
	out = append(out, passthrough...)
	return out
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
