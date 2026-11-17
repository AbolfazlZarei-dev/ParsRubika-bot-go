# `analytics` — Package analytics is a plugin that counts updates by type and

_Source: [`analytics.go`](../../../../ParsRubika/plugins/builtin/analytics/analytics.go)_

## Symbols

### `Tracker`

```go
type Tracker interface
```

Tracker is the published service.

### `Plugin`

```go
type Plugin struct
```

Plugin implements plugins.Plugin and Tracker.

### `New`

```go
func New() *Plugin
```

New creates the plugin.

### `Setup`

```go
func (p *Plugin) Setup(_ context.Context, h *plugins.Host) error
```

Setup implements plugins.Plugin.

### `Inc`

```go
func (p *Plugin) Inc(name string, d float64) { p.mu.Lock(); p.counters[name] += d; p.mu.Unlock() }
```

Inc implements Tracker.

### `SetGauge`

```go
func (p *Plugin) SetGauge(name string, v float64) { p.mu.Lock(); p.gauges[name] = v; p.mu.Unlock() }
```

SetGauge implements Tracker.

### `Counter`

```go
func (p *Plugin) Counter(name string) float64
```

Counter implements Tracker.

### `Gauge`

```go
func (p *Plugin) Gauge(name string) float64
```

Gauge implements Tracker.

### `Snapshot`

```go
func (p *Plugin) Snapshot() (map[string]float64, map[string]float64)
```

Snapshot implements Tracker (copies of both maps).

### `Names`

```go
func (p *Plugin) Names() []string
```

Names returns the sorted counter names.
