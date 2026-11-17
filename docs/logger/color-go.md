# `logger` — Color

_Source: [`color.go`](../../ParsRubika/logger/color.go)_

## Symbols

### `ColorMode`

```go
type ColorMode int
```

ColorMode controls when ANSI colours are emitted.

### `ColorAuto`

```go
const ColorAuto ColorMode = iota
```

ColorAuto enables colours when the output is a capable terminal.

### `ColorAlways`

```go
const ColorAlways
```

ColorAlways forces colours on.

### `ColorNever`

```go
const ColorNever
```

ColorNever disables colours.
