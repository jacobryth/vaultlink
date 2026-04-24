package filter

// Chain applies multiple secret maps through a sequence of filter functions,
// returning only the key-value pairs that pass all filters.

// FilterFunc is a function that takes a map of secrets and returns a filtered map.
type FilterFunc func(secrets map[string]string) map[string]string

// Chain holds an ordered list of FilterFuncs to apply in sequence.
type Chain struct {
	steps []FilterFunc
}

// NewChain creates a Chain from the provided FilterFuncs.
func NewChain(steps ...FilterFunc) *Chain {
	return &Chain{steps: steps}
}

// Apply runs the secrets map through each step in order.
// If secrets is nil, an empty map is returned immediately.
// Each step receives the output of the previous step, so filters are cumulative.
// Processing stops early if the result becomes empty.
func (c *Chain) Apply(secrets map[string]string) map[string]string {
	if secrets == nil {
		return map[string]string{}
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		result[k] = v
	}
	for _, step := range c.steps {
		result = step(result)
		if len(result) == 0 {
			return result
		}
	}
	return result
}

// Append returns a new Chain with the provided FilterFuncs added after the
// existing steps. The original Chain is not modified.
func (c *Chain) Append(steps ...FilterFunc) *Chain {
	newSteps := make([]FilterFunc, len(c.steps)+len(steps))
	copy(newSteps, c.steps)
	copy(newSteps[len(c.steps):], steps)
	return &Chain{steps: newSteps}
}

// Len returns the number of steps in the chain.
func (c *Chain) Len() int {
	return len(c.steps)
}
