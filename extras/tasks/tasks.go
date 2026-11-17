// Package tasks provides a goroutine-pool task distributor: a bounded
// worker pool, useful for fan-out
// work where you want to cap concurrency.
package tasks

import (
	"context"
	"sync"
	"sync/atomic"
)

// Task is a function executed by the pool.
type Task func()

// Pool is a bounded goroutine pool. It keeps `workers` goroutines alive
// and distributes tasks across them via a channel.
//
// Pool is safe for concurrent use. Submit never panics on a stopped
// pool (it silently drops the task instead); use SubmitNonBlocking if
// you need to detect that case.
type Pool struct {
	tasks    chan Task
	wg       sync.WaitGroup // worker lifecycle
	inflight sync.WaitGroup // submitted-but-not-completed
	stopCh   chan struct{}
	stopped  atomic.Bool
}

// NewPool constructs a pool with `workers` goroutines.
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = 4
	}
	p := &Pool{
		tasks:  make(chan Task, workers*2),
		stopCh: make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	return p
}

// Submit adds a task to the pool. It blocks when the channel is full.
// Submitting to a stopped pool is a no-op (the task is dropped) — it
// never panics on a closed channel. This is guaranteed by checking
// p.stopped atomically BEFORE the select: if the pool is already
// stopped we drop the task without touching the (possibly closed)
// tasks channel.
func (p *Pool) Submit(t Task) {
	if p.stopped.Load() {
		return
	}
	p.inflight.Add(1)
	// Re-check stopped under the inflight Add so we don't race with
	// a concurrent Stop that closes p.tasks.
	if p.stopped.Load() {
		p.inflight.Done()
		return
	}
	select {
	case <-p.stopCh:
		p.inflight.Done()
	case p.tasks <- t:
	}
}

// SubmitNonBlocking tries to add a task without blocking. Returns true
// when the task was queued, false when the pool is stopping or full.
func (p *Pool) SubmitNonBlocking(t Task) bool {
	if p.stopped.Load() {
		return false
	}
	p.inflight.Add(1)
	if p.stopped.Load() {
		p.inflight.Done()
		return false
	}
	select {
	case <-p.stopCh:
		p.inflight.Done()
		return false
	case p.tasks <- t:
		return true
	default:
		p.inflight.Done()
		return false
	}
}

// Stop drains the pool and waits for all workers to finish. Calling
// Stop more than once is a no-op. Stop sets the stopped flag FIRST so
// that concurrent Submit calls see it and bail out before touching
// the tasks channel; only then are stopCh and tasks closed.
func (p *Pool) Stop() {
	if !p.stopped.CompareAndSwap(false, true) {
		return
	}
	close(p.stopCh) // reject new submits (and unblock waiting ones)
	close(p.tasks)  // tell workers to drain and exit
	p.wg.Wait()
}

// Wait blocks until all tasks submitted *so far* have completed. It
// does NOT stop the pool — new tasks can be submitted after Wait
// returns. To tear down the pool, call Stop.
func (p *Pool) Wait() {
	p.inflight.Wait()
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for t := range p.tasks {
		if t != nil {
			t()
		}
		p.inflight.Done()
	}
}

// Distributor is a simple task distributor that limits concurrency to
// `max` goroutines. Unlike Pool, it does not keep workers alive between
// tasks — each Submit spawns a goroutine that exits when the task
// finishes.
type Distributor struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// NewDistributor constructs a distributor with max concurrency.
func NewDistributor(max int) *Distributor {
	if max <= 0 {
		max = 4
	}
	return &Distributor{sem: make(chan struct{}, max)}
}

// Submit runs `fn` in a new goroutine, blocking when the semaphore is
// full.
func (d *Distributor) Submit(fn func()) {
	d.wg.Add(1)
	d.sem <- struct{}{}
	go func() {
		defer d.wg.Done()
		defer func() { <-d.sem }()
		if fn != nil {
			fn()
		}
	}()
}

// SubmitCtx runs `fn` in a new goroutine, passing the supplied context.
func (d *Distributor) SubmitCtx(ctx context.Context, fn func(context.Context)) {
	d.wg.Add(1)
	d.sem <- struct{}{}
	go func() {
		defer d.wg.Done()
		defer func() { <-d.sem }()
		if fn != nil {
			fn(ctx)
		}
	}()
}

// Wait blocks until all submitted tasks have finished.
func (d *Distributor) Wait() { d.wg.Wait() }
