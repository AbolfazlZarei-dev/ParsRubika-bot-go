package bot

import (
	"sync"
	"time"
)

// dispatcher runs tasks concurrently while keeping tasks that share a
// key (the chat ID) strictly ordered. At most `limit` tasks run at the
// same time across all keys.
type dispatcher struct {
	sem    chan struct{}
	mu     sync.Mutex
	queues map[string]*keyQueue
	wg     sync.WaitGroup
}

type keyQueue struct {
	tasks   []func()
	running bool
}

func newDispatcher(limit int) *dispatcher {
	if limit < 1 {
		limit = 1
	}
	return &dispatcher{sem: make(chan struct{}, limit), queues: map[string]*keyQueue{}}
}

// submit schedules fn behind every earlier task of the same key.
func (d *dispatcher) submit(key string, fn func()) {
	d.mu.Lock()
	q := d.queues[key]
	if q == nil {
		q = &keyQueue{}
		d.queues[key] = q
	}
	q.tasks = append(q.tasks, fn)
	if q.running {
		d.mu.Unlock()
		return
	}
	q.running = true
	d.wg.Add(1)
	d.mu.Unlock()
	go d.drain(key, q)
}

func (d *dispatcher) drain(key string, q *keyQueue) {
	defer d.wg.Done()
	for {
		d.mu.Lock()
		if len(q.tasks) == 0 {
			q.running = false
			delete(d.queues, key)
			d.mu.Unlock()
			return
		}
		fn := q.tasks[0]
		q.tasks[0] = nil
		q.tasks = q.tasks[1:]
		d.mu.Unlock()

		d.sem <- struct{}{}
		fn()
		<-d.sem
	}
}

// wait blocks until every queued task finished or the timeout elapsed.
// It reports whether the queue drained in time.
func (d *dispatcher) wait(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() { d.wg.Wait(); close(done) }()
	if timeout <= 0 {
		<-done
		return true
	}
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// pending returns the number of keys with queued or running work.
func (d *dispatcher) pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.queues)
}

// seenSet is a fixed-capacity set that forgets its oldest entries first.
type seenSet struct {
	mu   sync.Mutex
	m    map[string]struct{}
	ring []string
	pos  int
}

func newSeenSet(capacity int) *seenSet {
	if capacity < 16 {
		capacity = 16
	}
	return &seenSet{m: make(map[string]struct{}, capacity), ring: make([]string, capacity)}
}

// add inserts key and reports whether it was new.
func (s *seenSet) add(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.m[key]; ok {
		return false
	}
	if old := s.ring[s.pos]; old != "" {
		delete(s.m, old)
	}
	s.ring[s.pos] = key
	s.m[key] = struct{}{}
	s.pos = (s.pos + 1) % len(s.ring)
	return true
}

// has reports whether key is currently remembered.
func (s *seenSet) has(key string) bool {
	s.mu.Lock()
	_, ok := s.m[key]
	s.mu.Unlock()
	return ok
}
