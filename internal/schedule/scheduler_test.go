package schedule

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew_InvalidInterval(t *testing.T) {
	_, err := New(0, func(ctx context.Context) error { return nil }, nil)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_NilJobFn(t *testing.T) {
	_, err := New(time.Second, nil, nil)
	if err == nil {
		t.Fatal("expected error for nil jobFn")
	}
}

func TestRunOnce_Success(t *testing.T) {
	called := false
	s, err := New(time.Second, func(ctx context.Context) error {
		called = true
		return nil
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := s.RunOnce(context.Background()); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if !called {
		t.Fatal("expected jobFn to be called")
	}
}

func TestRun_TicksAndCancels(t *testing.T) {
	var count int64
	s, _ := New(20*time.Millisecond, func(ctx context.Context) error {
		atomic.AddInt64(&count, 1)
		return nil
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	got := atomic.LoadInt64(&count)
	if got < 2 {
		t.Fatalf("expected at least 2 ticks, got %d", got)
	}
}

func TestRun_OnErrorCalled(t *testing.T) {
	var errCount int64
	s, _ := New(20*time.Millisecond, func(ctx context.Context) error {
		return errors.New("boom")
	}, func(err error) {
		atomic.AddInt64(&errCount, 1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	if atomic.LoadInt64(&errCount) == 0 {
		t.Fatal("expected onError to be called")
	}
}
