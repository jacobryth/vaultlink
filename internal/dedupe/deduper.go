package dedupe

import "strings"

// Strategy defines how duplicate keys are resolved.
type Strategy string

const (
	StrategyFirst Strategy = "first"
	StrategyLast  Strategy = "last"
	StrategyError Strategy = "error"
)

// Deduper removes or resolves duplicate keys in a secrets map.
type Deduper struct {
	strategy Strategy
}

// New creates a Deduper with the given strategy.
// Defaults to "first" if unknown.
func New(strategy string) (*Deduper, error) {
	s := Strategy(strings.ToLower(strategy))
	switch s {
	case StrategyFirst, StrategyLast, StrategyError:
		return &Deduper{strategy: s}, nil
	case "":
		return &Deduper{strategy: StrategyFirst}, nil
	default:
		return nil, fmt.Errorf("dedupe: unknown strategy %q", strategy)
	}
}

// Apply resolves duplicates in a slice of key-value pairs.
// Returns a deduplicated map and any error (e.g. for StrategyError).
func (d *Deduper) Apply(pairs []KV) (map[string]string, error) {
	if pairs == nil {
		return map[string]string{}, nil
	}
	seen := map[string]string{}
	for _, kv := range pairs {
		_, exists := seen[kv.Key]
		switch d.strategy {
		case StrategyFirst:
			if !exists {
				seen[kv.Key] = kv.Value
			}
		case StrategyLast:
			seen[kv.Key] = kv.Value
		case StrategyError:
			if exists {
				return nil, fmt.Errorf("dedupe: duplicate key %q", kv.Key)
			}
			seen[kv.Key] = kv.Value
		}
	}
	return seen, nil
}

// KV is a key-value pair that may contain duplicates.
type KV struct {
	Key   string
	Value string
}
