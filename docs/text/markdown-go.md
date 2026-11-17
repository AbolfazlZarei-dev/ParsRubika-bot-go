# `text` — Markdown

_Source: [`markdown.go`](../../ParsRubika/text/markdown.go)_

## Symbols

### `MarkupType`

```go
type MarkupType struct
```

MarkupType maps a markdown prefix (e.g. "**", "__", "```") to its
wire-format type name and the regex group index that holds the inner
content (1-indexed, like regexp submatch groups).

### `MarkupTypeSequence`

```go
var MarkupTypeSequence = []MarkupType{
```

MarkupTypeSequence lists the markup kinds in the order in which the
regex alternation tries to match them.

### `MentionTypes`

```go
var MentionTypes = map[byte]string{
```

MentionTypes maps a one-letter URL prefix (u/g/c/b) to its mention
type used on the wire.

### `Part`

```go
type Part struct
```

Part is the on-wire shape of one formatted region.

### `Link`

```go
type Link struct
```

Link is the structured link object attached to a Link-style Part.

### `HyperlinkData`

```go
type HyperlinkData struct
```

HyperlinkData holds the URL when the link is a plain hyperlink.

### `Parts`

```go
type Parts struct
```

Parts is the outer wrapper around the meta_data_parts slice.

### `Result`

```go
type Result struct
```

Result is the structured result of Format. It is what gets merged
into the `sendMessage` payload (only Text and Metadata fields).

### `Format`

```go
func (f Formatter) Format(text string) Result
```

Format converts a Markdown text into the wire-format Result.

Implementation note: the server measures offsets in UTF-16 code units,
while Go's `regexp.Regexp.FindAllStringSubmatchIndex` returns byte
indices and Go's string slicing is byte-based — which would corrupt
multi-byte UTF-8 sequences (e.g. Persian text). To produce exact
offsets the formatter:

1. Convert the input to `[]rune` so we can index by codepoint.
2. Build a byte→rune index lookup so the regex byte-indices can be
converted to rune indices.
3. Build the UTF-16 prefix array against the *original* text (so
the `from_index` / `length` values match what the server expects
for the original text, before any markup is stripped).

### `ToMarkdown`

```go
func (f Formatter) ToMarkdown(text string) string
```

ToMarkdown converts the most common HTML formatting tags to their
Markdown equivalents. If the input does not contain any HTML tag it
is returned unchanged (after trimming whitespace).

### `JSON`

```go
func (r Result) JSON() string
```

JSON returns the JSON representation of the result, useful when
merging it into a larger payload.
