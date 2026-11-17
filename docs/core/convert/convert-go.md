# `convert` — Package convert provides type-conversion helpers for the loose

_Source: [`convert.go`](../../../ParsRubika/core/convert/convert.go)_

## Symbols

### `ToString`

```go
func ToString(v any) string
```

ToString converts any scalar to its string form.

### `ToInt`

```go
func ToInt(v any) int
```

ToInt converts any numeric scalar to int.

### `ToInt64`

```go
func ToInt64(v any) int64
```

ToInt64 converts any numeric scalar to int64.

### `ToBool`

```go
func ToBool(v any) bool
```

ToBool converts any scalar to bool.

### `ToJSON`

```go
func ToJSON(v any) string
```

ToJSON marshals `v` to a JSON string.

### `FromJSON`

```go
func FromJSON(s string, out any) error
```

FromJSON unmarshals `s` into `out`.

### `MapToJSON`

```go
func MapToJSON(m map[string]any) string { return ToJSON(m) }
```

MapToJSON converts a map to a JSON string.

### `JSONToMap`

```go
func JSONToMap(s string) (map[string]any, error)
```

JSONToMap parses a JSON string into a map[string]any.
