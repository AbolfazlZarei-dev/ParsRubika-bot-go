# `rules` — Text

_Source: [`text.go`](../../../ParsRubika/bot/rules/text.go)_

## Symbols

### `Text`

```go
type Text struct
```

Text matches the message text either by exact string or by regex.
When `Pattern` is empty it matches any text message.

### `NewText`

```go
func NewText(pattern string, regex bool) *Text
```

NewText constructs a Text rule.

### `Check`

```go
func (f *Text) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `Command`

```go
type Command struct
```

Command matches a bot command (e.g. `/start`, `!help`).

### `NewCommand`

```go
func NewCommand(cmds, prefixes []string, caseInsensitive, allowNoPrefix bool) *Command
```

NewCommand constructs a Command rule.

### `Check`

```go
func (c *Command) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.
