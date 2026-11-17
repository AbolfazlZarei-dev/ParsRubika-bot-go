# `tasks` — Package tasks provides a goroutine-pool task distributor: a bounded

_Source: [`tasks.go`](../../../ParsRubika/extras/tasks/tasks.go)_

## Symbols

### `Task`

```go
type Task func()
```

Task is a function executed by the pool.

### `Pool`

```go
type Pool struct
```

Pool is a bounded goroutine pool. It keeps `workers` goroutines alive
and distributes tasks across them via a channel.

Pool is safe for concurrent use. Submit never panics on a stopped
pool (it silently drops the task instead); use SubmitNonBlocking if
you need to detect that case.

### `NewPool`

```go
func NewPool(workers int) *Pool
```

NewPool constructs a pool with `workers` goroutines.

### `Submit`

```go
func (p *Pool) Submit(t Task)
```

Submit adds a task to the pool. It blocks when the channel is full.
Submitting to a stopped pool is a no-op (the task is dropped) — it
never panics on a closed channel. This is guaranteed by checking
p.stopped atomically BEFORE the select: if the pool is already
stopped we drop the task without touching the (possibly closed)
tasks channel.

### `SubmitNonBlocking`

```go
func (p *Pool) SubmitNonBlocking(t Task) bool
```

SubmitNonBlocking tries to add a task without blocking. Returns true
when the task was queued, false when the pool is stopping or full.

### `Stop`

```go
func (p *Pool) Stop()
```

Stop drains the pool and waits for all workers to finish. Calling
Stop more than once is a no-op. Stop sets the stopped flag FIRST so
that concurrent Submit calls see it and bail out before touching
the tasks channel; only then are stopCh and tasks closed.

### `Wait`

```go
func (p *Pool) Wait()
```

Wait blocks until all tasks submitted *so far* have completed. It
does NOT stop the pool — new tasks can be submitted after Wait
returns. To tear down the pool, call Stop.

### `Distributor`

```go
type Distributor struct
```

Distributor is a simple task distributor that limits concurrency to
`max` goroutines. Unlike Pool, it does not keep workers alive between
tasks — each Submit spawns a goroutine that exits when the task
finishes.

### `NewDistributor`

```go
func NewDistributor(max int) *Distributor
```

NewDistributor constructs a distributor with max concurrency.

### `Submit`

```go
func (d *Distributor) Submit(fn func())
```

Submit runs `fn` in a new goroutine, blocking when the semaphore is
full.

### `SubmitCtx`

```go
func (d *Distributor) SubmitCtx(ctx context.Context, fn func(context.Context))
```

SubmitCtx runs `fn` in a new goroutine, passing the supplied context.

### `Wait`

```go
func (d *Distributor) Wait() { d.wg.Wait() }
```

Wait blocks until all submitted tasks have finished.
