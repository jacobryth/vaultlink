package limit

// Ensure None level with zero count is valid (no-op limiter).
func NewNone() *Limiter {
	return &Limiter{level: None, count: 0}
}
