# `matchers` — Text

_Source: [`text.go`](../../../ParsRubika/account/matchers/text.go)_

## Symbols

### `Text`

```go
type Text struct
```

Text matches the message text. When `pattern` is empty it matches
any text; when `regex` is true the pattern is treated as a Go regex.

### `NewText`

```go
func NewText(pattern string, regex bool, kind string) *Text
```

NewText constructs a Text matcher. `kind` controls which update type
it accepts (defaults to "MessageUpdates").

### `Kind`

```go
func (t *Text) Kind() string { return t.kind }
```

Kind implements Matcher.

### `Match`

```go
func (t *Text) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `Command`

```go
type Command struct
```

Command matches a bot command (e.g. `/start`, `!help`).

### `NewCommand`

```go
func NewCommand(cmds, prefixes []string, caseInsensitive, allowNoPrefix bool, kind string) *Command
```

NewCommand constructs a Command matcher.

### `Kind`

```go
func (c *Command) Kind() string { return c.kind }
```

Kind implements Matcher.

### `Match`

```go
func (c *Command) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `ObjectGUID`

```go
type ObjectGUID struct
```

ObjectGUID matches by chat GUID.

### `NewObjectGUID`

```go
func NewObjectGUID(kind string, guids ...string) *ObjectGUID
```

NewObjectGUID constructs an ObjectGUID matcher.

### `Kind`

```go
func (o *ObjectGUID) Kind() string { return o.kind }
```

Kind implements Matcher.

### `Match`

```go
func (o *ObjectGUID) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `AuthorGUID`

```go
type AuthorGUID struct
```

AuthorGUID matches by author GUID.

### `NewAuthorGUID`

```go
func NewAuthorGUID(kind string, guids ...string) *AuthorGUID
```

NewAuthorGUID constructs an AuthorGUID matcher.

### `Kind`

```go
func (a *AuthorGUID) Kind() string { return a.kind }
```

Kind implements Matcher.

### `Match`

```go
func (a *AuthorGUID) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsGroup`

```go
type IsGroup struct
```

IsGroup matches when the update came from a group.

### `NewIsGroup`

```go
func NewIsGroup(kind string) *IsGroup
```

NewIsGroup constructs an IsGroup matcher.

### `Kind`

```go
func (g *IsGroup) Kind() string { return g.kind }
```

Kind implements Matcher.

### `Match`

```go
func (g *IsGroup) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsPrivate`

```go
type IsPrivate struct
```

IsPrivate matches when the update came from a private chat.

### `NewIsPrivate`

```go
func NewIsPrivate(kind string) *IsPrivate
```

NewIsPrivate constructs an IsPrivate matcher.

### `Kind`

```go
func (p *IsPrivate) Kind() string { return p.kind }
```

Kind implements Matcher.

### `Match`

```go
func (p *IsPrivate) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsChannel`

```go
type IsChannel struct
```

IsChannel matches when the update came from a channel.

### `NewIsChannel`

```go
func NewIsChannel(kind string) *IsChannel
```

NewIsChannel constructs an IsChannel matcher.

### `Kind`

```go
func (c *IsChannel) Kind() string { return c.kind }
```

Kind implements Matcher.

### `Match`

```go
func (c *IsChannel) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsText`

```go
type IsText struct
```

IsText matches text messages (any non-empty text).

### `NewIsText`

```go
func NewIsText(kind string) *IsText
```

NewIsText constructs an IsText matcher.

### `Kind`

```go
func (t *IsText) Kind() string { return t.kind }
```

Kind implements Matcher.

### `Match`

```go
func (t *IsText) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsFile`

```go
type IsFile struct
```

IsFile matches messages with a file_inline of type "File".

### `NewIsFile`

```go
func NewIsFile(kind string) *IsFile
```

NewIsFile constructs an IsFile matcher.

### `Kind`

```go
func (f *IsFile) Kind() string { return f.kind }
```

Kind implements Matcher.

### `Match`

```go
func (f *IsFile) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsPhoto`

```go
type IsPhoto struct
```

IsPhoto matches messages with a file_inline of type "Image".

### `NewIsPhoto`

```go
func NewIsPhoto(kind string) *IsPhoto
```

NewIsPhoto constructs an IsPhoto matcher.

### `Kind`

```go
func (p *IsPhoto) Kind() string { return p.kind }
```

Kind implements Matcher.

### `Match`

```go
func (p *IsPhoto) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsSticker`

```go
type IsSticker struct
```

IsSticker matches sticker messages.

### `NewIsSticker`

```go
func NewIsSticker(kind string) *IsSticker
```

NewIsSticker constructs an IsSticker matcher.

### `Kind`

```go
func (s *IsSticker) Kind() string { return s.kind }
```

Kind implements Matcher.

### `Match`

```go
func (s *IsSticker) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsForward`

```go
type IsForward struct
```

IsForward matches forwarded messages.

### `NewIsForward`

```go
func NewIsForward(kind string) *IsForward
```

NewIsForward constructs an IsForward matcher.

### `Kind`

```go
func (f *IsForward) Kind() string { return f.kind }
```

Kind implements Matcher.

### `Match`

```go
func (f *IsForward) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsEdited`

```go
type IsEdited struct
```

IsEdited matches edited messages.

### `NewIsEdited`

```go
func NewIsEdited(kind string) *IsEdited
```

NewIsEdited constructs an IsEdited matcher.

### `Kind`

```go
func (e *IsEdited) Kind() string { return e.kind }
```

Kind implements Matcher.

### `Match`

```go
func (e *IsEdited) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.

### `IsMe`

```go
type IsMe struct
```

IsMe matches updates where the author is the current user. The
current user's GUID must be supplied at construction time.

### `NewIsMe`

```go
func NewIsMe(selfGUID, kind string) *IsMe
```

NewIsMe constructs an IsMe matcher.

### `Kind`

```go
func (m *IsMe) Kind() string { return m.kind }
```

Kind implements Matcher.

### `Match`

```go
func (m *IsMe) Match(_ context.Context, raw map[string]any) bool
```

Match implements Matcher.
