# `enums` — Media type

_Source: [`media_type.go`](../../ParsRubika/enums/media_type.go)_

## Symbols

### `MediaType`

```go
type MediaType string
```

MediaType is the kind of inline attachment a message carries.

### `MediaTypeText`

```go
const MediaTypeText MediaType = "Text"
```

### `MediaTypeFile`

```go
const MediaTypeFile MediaType = "File"
```

### `MediaTypeFileInline`

```go
const MediaTypeFileInline MediaType = "FileInline"
```

### `MediaTypeSticker`

```go
const MediaTypeSticker MediaType = "Sticker"
```

### `MediaTypeCaption`

```go
const MediaTypeCaption MediaType = "FileInlineCaption"
```

### `MediaTypeEvent`

```go
const MediaTypeEvent MediaType = "Event"
```

### `String`

```go
func (m MediaType) String() string { return string(m) }
```

String implements fmt.Stringer.
