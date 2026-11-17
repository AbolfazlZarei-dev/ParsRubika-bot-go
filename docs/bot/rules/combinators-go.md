# `rules` — Package rules provides the composable rules used by the bot client

_Source: [`combinators.go`](../../../ParsRubika/bot/rules/combinators.go)_

## Symbols

### `Rule`

```go
type Rule interface
```

Rule is the contract implemented by every bot rule.

### `Check`

```go
func (a And) Check(ctx context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Check`

```go
func (o Or) Check(ctx context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Not`

```go
type Not struct
```

Not negates the wrapped rule.

### `Check`

```go
func (n Not) Check(ctx context.Context, upd *models.Update) bool { return !n.Inner.Check(ctx, upd) }
```

Check implements Rule.
