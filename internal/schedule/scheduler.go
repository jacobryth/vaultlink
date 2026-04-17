package schedule

import (
	"context"
	"fmt"
	"time"
)

// Scheduler runs a sync job on a fixed interval.
type Scheduler struct {
	interval time.Duration
	jobFn    func(ctx context.Context) error
	onError  func(err error)
}

// New creates a Scheduler with the given interval and job function.
func New(interval time.Duration, jobFn func(ctx context.Context) error, onError func(err error)) (*Scheduler, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("schedule: interval must be positive, got %s", interval)
	}
	if jobFn == nil {
		return nil, fmt.Errorf("schedule: jobFn must not be nil")
	}
	if onError == nil {
		onError = func(err error) {}
	}
	return &Scheduler{interval: interval, jobFn: jobFn, onError: onError}, nil
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.jobFn(ctx); err != nil {
				s.onError(err)
			}
		}
	}
}

// RunOnce executes the job function immediately once.
func (s *Scheduler) RunOnce(ctx context.Context) error {
	return s.jobFn(ctx)
}
