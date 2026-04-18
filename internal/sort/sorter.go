package sort

import (
	"fmt"
	"sort"
)

type Level string

const (
	None Level = "none"
	Asc  Level = "asc"
	Desc Level = "desc"
)

var knownLevels = map[Level]bool{
	None: true,
	Asc:  true,
	Desc: true,
}

type Sorter struct {
	level Level
}

func New(level string) (*Sorter, error) {
	l := Level(level)
	if !knownLevels[l] {
		return nil, fmt.Errorf("unknown sort level: %q (want: none, asc, desc)", level)
	}
	return &Sorter{level: l}, nil
}

func (s *Sorter) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return nil
	}
	if s.level == None {
		return secrets
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	if s.level == Asc {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	ordered := make(map[string]string, len(secrets))
	for _, k := range keys {
		ordered[k] = secrets[k]
	}
	return ordered
}
