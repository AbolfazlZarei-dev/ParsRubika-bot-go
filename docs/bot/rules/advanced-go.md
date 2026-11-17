# `rules` — Advanced

_Source: [`advanced.go`](../../../ParsRubika/bot/rules/advanced.go)_

## Symbols

### `PollQuestionContains`

```go
type PollQuestionContains struct
```

PollQuestionContains matches polls whose question contains `keyword`
(case-insensitive).

### `Check`

```go
func (p PollQuestionContains) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `PollOptionCount`

```go
type PollOptionCount struct
```

PollOptionCount matches polls with between min and max options
(inclusive).

### `Check`

```go
func (p PollOptionCount) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `LocationWithin`

```go
type LocationWithin struct
```

LocationWithin matches location messages whose latitude/longitude fall
within the given bounding box.

### `Check`

```go
func (l LocationWithin) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `MetaLinksContain`

```go
type MetaLinksContain struct
```

MetaLinksContain matches messages whose metadata contains a link URL
containing `substring`.

### `Check`

```go
func (m MetaLinksContain) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `MetaTypesInclude`

```go
type MetaTypesInclude struct
```

MetaTypesInclude matches messages whose metadata includes at least one
part of the given types (e.g. "Bold", "Italic", "Link").

### `Check`

```go
func (m MetaTypesInclude) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsReplyToUser`

```go
type IsReplyToUser struct
```

IsReplyToUser matches messages that are a reply to a message from the
given user ID.

### `Check`

```go
func (r IsReplyToUser) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsForwardedFrom`

```go
type IsForwardedFrom struct
```

IsForwardedFrom matches messages forwarded from the given user/channel.

### `Check`

```go
func (f IsForwardedFrom) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `EditedTextContains`

```go
type EditedTextContains struct
```

EditedTextContains matches edited messages whose new text contains
`keyword` (case-insensitive).

### `Check`

```go
func (e EditedTextContains) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `AuxDataContains`

```go
type AuxDataContains struct
```

AuxDataContains matches updates whose aux_data contains the given
key/value pair.

### `Check`

```go
func (a AuxDataContains) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextRegexCompiled`

```go
type TextRegexCompiled struct
```

TextRegexCompiled matches text using a pre-compiled regex. This is
more efficient than TextRegex when the same pattern is reused.

### `NewTextRegexCompiled`

```go
func NewTextRegexCompiled(re *regexp.Regexp) *TextRegexCompiled
```

NewTextRegexCompiled creates a TextRegexCompiled rule from a compiled
regex.

### `Check`

```go
func (t *TextRegexCompiled) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ChatTitleRegex`

```go
type ChatTitleRegex struct
```

ChatTitleRegex matches when the chat title (from ExternalData) matches
the given regex.

### `Check`

```go
func (c ChatTitleRegex) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ChatHasUsername`

```go
type ChatHasUsername struct
```

ChatHasUsername matches when the chat has a username (from ExternalData
or from a pre-fetched Chat object).

### `Check`

```go
func (ChatHasUsername) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.
