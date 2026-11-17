# `cron` — Package cron is a tiny cron-style scheduler built on top of the

_Source: [`cron.go`](../../../ParsRubika/extras/cron/cron.go)_

## Symbols

### `Job`

```go
type Job struct
```

Job is a single cron entry: a schedule + a callback.

### `Cronner`

```go
type Cronner struct
```

Cronner wraps a scheduler and a list of parsed jobs. Use New to
construct one, Add to register jobs, then Start to begin firing.

### `New`

```go
func New(parent context.Context) *Cronner
```

New constructs a fresh Cronner bound to the supplied context.

### `Add`

```go
func (c *Cronner) Add(spec string, cb func(context.Context)) error
```

Add parses `spec` and registers `cb` for it. Supported spec forms:

- "@yearly" / "@annually" — once a year, Jan 1 00:00
- "@monthly" — once a month, day 1 00:00
- "@weekly" — once a week, Sunday 00:00
- "@daily" / "@midnight" — once a day, 00:00
- "@hourly" — once an hour, at minute 0
- 5-field cron: "M H DOM MON DOW"
Each field can be a number, "*", "*/N" (every N units), or a
comma-separated list of numbers.

### `Start`

```go
func (c *Cronner) Start()
```

Start begins firing all registered jobs. Returns immediately.

### `Stop`

```go
func (c *Cronner) Stop()
```

Stop cancels every scheduled job.

### `Running`

```go
func (c *Cronner) Running() bool
```

Running reports whether the cronner is currently active.

### `Next`

```go
func (j *parsedJob) Next(from time.Time) time.Time
```

Next returns the next time after `from` that matches the schedule.
It uses the standard cron field-by-field forward-scanning algorithm
(the same one used by Vixie cron and robfig/cron) instead of a
minute-by-minute brute force, so even `@yearly` resolves in O(1)
iterations rather than ~5 million.

Returns the zero time when no match exists in the next 5 years.
