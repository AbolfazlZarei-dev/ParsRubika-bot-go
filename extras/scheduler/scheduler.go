// Package scheduler provides a lightweight, goroutine-safe scheduler
// for one-shot and recurring tasks. Tasks can be cancelled, listed and
// replaced by name.
package scheduler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Scheduler is a goroutine-safe scheduler that runs registered
// callbacks either once (after a delay) or repeatedly (every interval).
// It uses a single goroutine per task — Go's runtime is efficient
// enough that this scales to thousands of tasks.
type Scheduler struct {
	mu      sync.Mutex
	tasks   map[string]*Task
	counter atomic.Int64
	cancel  context.CancelFunc
	ctx     context.Context
}

// Task is a single scheduled callback.
type Task struct {
	ID       string
	interval time.Duration
	once     bool
	callback func(context.Context)
	cancel   context.CancelFunc
}

// New returns a fresh Scheduler. Pass a parent context to be able to
// stop all tasks at once via `Stop()`.
func New(parent context.Context) *Scheduler {
	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)
	return &Scheduler{
		tasks:  map[string]*Task{},
		ctx:    ctx,
		cancel: cancel,
	}
}

// After schedules `callback` to run once after `delay`. Returns the
// task ID, which can be passed to Cancel to prevent it from firing.
func (s *Scheduler) After(delay time.Duration, callback func(context.Context)) string {
	id := s.nextID()
	t := &Task{ID: id, interval: delay, once: true, callback: callback}
	ctx, cancel := context.WithCancel(s.ctx)
	t.cancel = cancel
	s.mu.Lock()
	s.tasks[id] = t
	s.mu.Unlock()
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
			callback(ctx)
			s.mu.Lock()
			delete(s.tasks, id)
			s.mu.Unlock()
		}
	}()
	return id
}

// Every schedules `callback` to run repeatedly every `interval`. The
// first run happens immediately after `interval` (NOT at t=0).
func (s *Scheduler) Every(interval time.Duration, callback func(context.Context)) string {
	id := s.nextID()
	t := &Task{ID: id, interval: interval, once: false, callback: callback}
	ctx, cancel := context.WithCancel(s.ctx)
	t.cancel = cancel
	s.mu.Lock()
	s.tasks[id] = t
	s.mu.Unlock()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				callback(ctx)
			}
		}
	}()
	return id
}

// Cancel removes a scheduled task. Returns true if the task existed.
func (s *Scheduler) Cancel(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return false
	}
	t.cancel()
	delete(s.tasks, id)
	return true
}

// CancelAll cancels every registered task.
func (s *Scheduler) CancelAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tasks {
		t.cancel()
	}
	s.tasks = map[string]*Task{}
}

// Stop cancels the parent context, which in turn cancels every task.
// The Scheduler cannot be reused after Stop.
func (s *Scheduler) Stop() {
	s.cancel()
	s.CancelAll()
}

// TaskIDs returns the IDs of all currently registered tasks.
func (s *Scheduler) TaskIDs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.tasks))
	for id := range s.tasks {
		out = append(out, id)
	}
	return out
}

// Count returns the number of currently registered tasks.
func (s *Scheduler) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.tasks)
}

func (s *Scheduler) nextID() string {
	return fmt.Sprintf("task-%d", s.counter.Add(1))
}
