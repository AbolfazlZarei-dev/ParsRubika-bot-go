# `scheduler` — Package scheduler provides a lightweight, goroutine-safe scheduler

_Source: [`scheduler.go`](../../../ParsRubika/extras/scheduler/scheduler.go)_

## Symbols

### `Scheduler`

```go
type Scheduler struct
```

Scheduler is a goroutine-safe scheduler that runs registered
callbacks either once (after a delay) or repeatedly (every interval).
It uses a single goroutine per task — Go's runtime is efficient
enough that this scales to thousands of tasks.

### `Task`

```go
type Task struct
```

Task is a single scheduled callback.

### `New`

```go
func New(parent context.Context) *Scheduler
```

New returns a fresh Scheduler. Pass a parent context to be able to
stop all tasks at once via `Stop()`.

### `After`

```go
func (s *Scheduler) After(delay time.Duration, callback func(context.Context)) string
```

After schedules `callback` to run once after `delay`. Returns the
task ID, which can be passed to Cancel to prevent it from firing.

### `Every`

```go
func (s *Scheduler) Every(interval time.Duration, callback func(context.Context)) string
```

Every schedules `callback` to run repeatedly every `interval`. The
first run happens immediately after `interval` (NOT at t=0).

### `Cancel`

```go
func (s *Scheduler) Cancel(id string) bool
```

Cancel removes a scheduled task. Returns true if the task existed.

### `CancelAll`

```go
func (s *Scheduler) CancelAll()
```

CancelAll cancels every registered task.

### `Stop`

```go
func (s *Scheduler) Stop()
```

Stop cancels the parent context, which in turn cancels every task.
The Scheduler cannot be reused after Stop.

### `TaskIDs`

```go
func (s *Scheduler) TaskIDs() []string
```

TaskIDs returns the IDs of all currently registered tasks.

### `Count`

```go
func (s *Scheduler) Count() int
```

Count returns the number of currently registered tasks.
