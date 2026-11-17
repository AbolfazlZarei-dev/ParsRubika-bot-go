# `text` — Package text implements the Markdown/HTML → metadata conversion used

_Source: [`formatter.go`](../../ParsRubika/text/formatter.go)_

## Symbols

### `Formatter`

```go
type Formatter struct
```

Formatter is the public, stateless converter. Use New() to obtain a
zero-cost instance (it is just a struct with no fields).

### `New`

```go
func New() Formatter { return Formatter{} }
```

New returns a fresh Formatter. It is a zero-allocation constructor
that exists only so callers don't have to write `text.Formatter{}`
everywhere.
