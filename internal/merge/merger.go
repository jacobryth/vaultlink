package merge

import "fmt"

// Strategy defines how duplicate keys are resolved when merging secret maps.
type Strategy string

const (
	StrategyNone     Strategy = "none"
	StrategyOverwrite Strategy = "overwrite"
	StrategyKeepFirst Strategy = "keep-first"
)

var knownStrategies = map[Strategy]bool{
	StrategyNone:      true,
	StrategyOverwrite:  true,
	StrategyKeepFirst:  true,
}

// Merger combines multiple secret maps into one according to a strategy.
type Merger struct {
	strategy Strategy
}

// New returns a Merger for the given strategy, or an error if unrecognised.
func New(strategy Strategy) (*Merger, error) {
	if !knownStrategies[strategy] {
		return nil, fmt.Errorf("merge: unknown strategy %q", strategy)
	}
	return &Merger{strategy: strategy}, nil
}

// Apply merges layers in order. Later layers take precedence when
// StrategyOverwrite is used; earlier layers win under StrategyKeepFirst.
// StrategyNone returns the first layer unchanged.
func (m *Merger) Apply(layers ...map[string]string) map[string]string {
	if len(layers) == 0 {
		return map[string]string{}
	}
	if m.strategy == StrategyNone {
		if layers[0] == nil {
			return map[string]string{}
		}
		out := make(map[string]string, len(layers[0]))
		for k, v := range layers[0] {
			out[k] = v
		}
		return out
	}

	out := make(map[string]string)
	for _, layer := range layers {
		for k, v := range layer {
			_, exists := out[k]
			switch m.strategy {
			case StrategyOverwrite:
				out[k] = v
			case StrategyKeepFirst:
				if !exists {
					out[k] = v
				}
			}
		}
	}
	return out
}
