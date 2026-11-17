# `enums` — Parse mode

_Source: [`parse_mode.go`](../../ParsRubika/enums/parse_mode.go)_

## Symbols

### `ParseMode`

```go
type ParseMode string
```

ParseMode controls how text is formatted before being sent to the
server.

### `ParseModeNone`

```go
const ParseModeNone ParseMode = ""
```

ParseModeNone sends text as-is (no metadata generation).

### `ParseModeMarkdown`

```go
const ParseModeMarkdown ParseMode = "markdown"
```

ParseModeMarkdown runs the text through the Markdown formatter.

### `ParseModeHTML`

```go
const ParseModeHTML ParseMode = "html"
```

ParseModeHTML first converts HTML to Markdown via ToMarkdown and
then runs the Markdown formatter.

### `String`

```go
func (p ParseMode) String() string { return string(p) }
```

String implements fmt.Stringer.

### `Valid`

```go
func (p ParseMode) Valid() bool
```

Valid reports whether p is one of the supported parse modes.
