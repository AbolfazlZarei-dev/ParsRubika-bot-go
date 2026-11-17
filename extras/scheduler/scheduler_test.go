package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestAfter(t *testing.T) {
	s := New(context.Background())
	var called atomic.Int32
	s.After(50*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	time.Sleep(150 * time.Millisecond)
	if called.Load() != 1 {
		t.Errorf("After callback called %d times, want 1", called.Load())
	}
}

func TestEvery(t *testing.T) {
	s := New(context.Background())
	var called atomic.Int32
	id := s.Every(20*time.Millisecond, func(ctx context.Context) {
		called.Add(1)
	})
	time.Sleep(110 * time.Millisecond)
	if called.Load() < 3 {
		t.Errorf("Every callback called %d times, want >= 3", called.Load())
	}
	if !s.Cancel(id) {
		t.Error("Cancel should return true")
	}
	if s.Cancel(id) {
		t.Error("Cancel after removal should return false")
	}
}

func TestCancelAll(t *testing.T) {
	s := New(context.Background())
	var called atomic.Int32
	for i := 0; i < 5; i++ {
		s.After(100*time.Millisecond, func(ctx context.Context) {
			called.Add(1)
		})
	}
	if s.Count() != 5 {
		t.Errorf("Count = %d, want 5", s.Count())
	}
	s.CancelAll()
	if s.Count() != 0 {
		t.Errorf("Count after CancelAll = %d, want 0", s.Count())
	}
	time.Sleep(200 * time.Millisecond)
	if called.Load() != 0 {
		t.Errorf("called = %d, want 0 (all cancelled)", called.Load())
	}
}

func TestStop(t *testing.T) {
	s := New(context.Background())
	s.After(50*time.Millisecond, func(ctx context.Context) {})
	s.Stop()
	if s.Count() != 0 {
		t.Errorf("Count after Stop = %d, want 0", s.Count())
	}
}
