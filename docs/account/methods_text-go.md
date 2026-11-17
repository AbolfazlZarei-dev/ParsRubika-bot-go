# `account` — Methods text

_Source: [`methods_text.go`](../../ParsRubika/account/methods_text.go)_

## Symbols

### `Bold`

```go
func Bold(text string) string { return "**" + text + "**" }
```

Bold wraps text in **bold** markers.

### `Italic`

```go
func Italic(text string) string { return "__" + text + "__" }
```

Italic wraps text in __italic__ markers.

### `Underline`

```go
func Underline(text string) string { return "--" + text + "--" }
```

Underline wraps text in --underline-- markers.

### `Strike`

```go
func Strike(text string) string { return "~~" + text + "~~" }
```

Strike wraps text in ~~strike~~ markers.

### `Spoiler`

```go
func Spoiler(text string) string { return "||" + text + "||" }
```

Spoiler wraps text in ||spoiler|| markers.

### `Code`

```go
func Code(text string) string { return "`" + text + "`" }
```

Code wraps text in `code` markers (inline monospace).

### `Pre`

```go
func Pre(text string) string { return "```" + text + "```" }
```

Pre wraps text in ```pre``` markers (block monospace).

### `Hyperlink`

```go
func Hyperlink(text, url string) string { return "[" + text + "](" + url + ")" }
```

Hyperlink creates a [text](url) markdown link.

### `Mention`

```go
func Mention(text, guid string) string
```

Mention creates a mention link for a GUID. The URL prefix is inferred
from the GUID: u0→User, g0→Group, c0→Channel, b0→Bot.

### `Quote`

```go
func Quote(text string) string { return "> " + text }
```

Quote creates a >quote block.
