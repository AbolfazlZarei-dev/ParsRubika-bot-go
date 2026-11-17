# `enums` — Access

_Source: [`access.go`](../../ParsRubika/enums/access.go)_

## Symbols

### `Access`

```go
type Access string
```

Access is the access-right flag used in chat permissions.

### `AccessRead`

```go
const AccessRead Access = "Read"
```

### `AccessWrite`

```go
const AccessWrite Access = "Write"
```

### `String`

```go
func (a Access) String() string { return string(a) }
```

String implements fmt.Stringer.
