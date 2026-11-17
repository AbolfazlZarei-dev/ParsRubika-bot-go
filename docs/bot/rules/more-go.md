# `rules` — More

_Source: [`more.go`](../../../ParsRubika/bot/rules/more.go)_

## Symbols

### `Func`

```go
type Func func(ctx context.Context, upd *models.Update) bool
```

Func adapts a plain function into a Rule.

### `Check`

```go
func (f Func) Check(ctx context.Context, upd *models.Update) bool { return f != nil && f(ctx, upd) }
```

Check implements Rule.

### `Always`

```go
type Always struct
```

Always matches every update.

### `Check`

```go
func (Always) Check(context.Context, *models.Update) bool { return true }
```

Check implements Rule.

### `Never`

```go
type Never struct
```

Never matches no update.

### `Check`

```go
func (Never) Check(context.Context, *models.Update) bool { return false }
```

Check implements Rule.

### `HasText`

```go
type HasText struct
```

HasText matches updates whose message carries non-empty text.

### `Check`

```go
func (HasText) Check(_ context.Context, u *models.Update) bool { return u.Text() != "" }
```

Check implements Rule.

### `TextInRule`

```go
type TextInRule struct
```

TextInRule matches when the text equals any of the given strings.

### `TextIn`

```go
func TextIn(values ...string) *TextInRule
```

TextIn builds a TextInRule.

### `Check`

```go
func (r *TextInRule) Check(_ context.Context, u *models.Update) bool
```

Check implements Rule.

### `TextContainsAny`

```go
type TextContainsAny struct
```

TextContainsAny matches when the text contains any of the substrings.

### `Check`

```go
func (r TextContainsAny) Check(_ context.Context, u *models.Update) bool
```

Check implements Rule.

### `SenderIn`

```go
type SenderIn struct
```

SenderIn matches updates produced by one of the given user IDs.

### `NewSenderIn`

```go
func NewSenderIn(ids ...string) *SenderIn
```

NewSenderIn builds a SenderIn rule.

### `Check`

```go
func (r *SenderIn) Check(_ context.Context, u *models.Update) bool
```

Check implements Rule.

### `ChatIn`

```go
type ChatIn struct
```

ChatIn matches updates of one of the given chats.

### `NewChatIn`

```go
func NewChatIn(ids ...string) *ChatIn
```

NewChatIn builds a ChatIn rule.

### `Check`

```go
func (r *ChatIn) Check(_ context.Context, u *models.Update) bool
```

Check implements Rule.

### `MessageIDIs`

```go
type MessageIDIs string
```

MessageIDIs matches one exact message ID.

### `Check`

```go
func (r MessageIDIs) Check(_ context.Context, u *models.Update) bool
```

Check implements Rule.

### `PrivateChat`

```go
type PrivateChat struct
```

PrivateChat matches one-to-one chats only.

### `Check`

```go
func (PrivateChat) Check(_ context.Context, u *models.Update) bool { return u.IsPrivateChat() }
```

Check implements Rule.
