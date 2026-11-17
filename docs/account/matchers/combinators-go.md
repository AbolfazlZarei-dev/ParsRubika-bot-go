# `matchers` — Combinators

_Source: [`combinators.go`](../../../ParsRubika/account/matchers/combinators.go)_

## Symbols

### `And`

```go
type And struct
```

And matches when ALL sub-matchers match (short-circuit).

### `NewAnd`

```go
func NewAnd(kind string, ms ...Matcher) *And
```

NewAnd constructs an And matcher.

### `Kind`

```go
func (a *And) Kind() string { return a.kind }
```

Kind implements Matcher.

### `Match`

```go
func (a *And) Match(ctx context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `Or`

```go
type Or struct
```

Or matches when ANY sub-matcher matches.

### `NewOr`

```go
func NewOr(kind string, ms ...Matcher) *Or
```

NewOr constructs an Or matcher.

### `Kind`

```go
func (o *Or) Kind() string { return o.kind }
```

Kind implements Matcher.

### `Match`

```go
func (o *Or) Match(ctx context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `Not`

```go
type Not struct
```

Not negates the wrapped matcher.

### `NewNot`

```go
func NewNot(kind string, m Matcher) *Not
```

NewNot constructs a Not matcher.

### `Kind`

```go
func (n *Not) Kind() string { return n.kind }
```

Kind implements Matcher.

### `Match`

```go
func (n *Not) Match(ctx context.Context, raw map[string]any) bool
```

Match implements Matcher.
