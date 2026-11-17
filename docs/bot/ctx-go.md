# `bot` — Ctx

_Source: [`ctx.go`](../../ParsRubika/bot/ctx.go)_

## Symbols

### `HandlerFunc`

```go
type HandlerFunc func(c *Ctx) error
```

HandlerFunc is the ergonomic handler signature: everything a handler
needs is reachable from the *Ctx.

### `Ctx`

```go
type Ctx struct
```

Ctx bundles an update with the bot and offers short helpers for the
most common reactions (reply, edit, delete, forward, arguments, ...).
It embeds the request context, so it can be passed anywhere a
context.Context is expected.

### `NewCtx`

```go
func NewCtx(ctx context.Context, b *Bot, u *Update) *Ctx
```

NewCtx wraps an update. Handlers registered through Handle, Command,
Text, Regex, Button and the other shortcuts receive one automatically.

### `Message`

```go
func (c *Ctx) Message() *models.Message { return c.Update.Message() }
```

Message returns the incoming message, or nil.

### `ChatID`

```go
func (c *Ctx) ChatID() string { return c.Update.ChatID }
```

ChatID returns the chat the update belongs to.

### `SenderID`

```go
func (c *Ctx) SenderID() string { return c.Update.SenderID() }
```

SenderID returns the ID of the user that wrote the message.

### `MessageID`

```go
func (c *Ctx) MessageID() string { return c.Update.MessageID() }
```

MessageID returns the ID of the incoming message.

### `Text`

```go
func (c *Ctx) Text() string { return c.Update.Text() }
```

Text returns the text of the incoming message.

### `ButtonID`

```go
func (c *Ctx) ButtonID() string { return c.Update.ButtonID() }
```

ButtonID returns the ID of the pressed inline button, if any.

### `IsPrivate`

```go
func (c *Ctx) IsPrivate() bool { return c.Update.IsPrivateChat() }
```

IsPrivate reports whether the chat is a one-to-one chat.

### `Command`

```go
func (c *Ctx) Command() string
```

Command returns the command name of the message without its prefix
("/start now" -> "start"), or "" if the text is not a command.

### `Args`

```go
func (c *Ctx) Args() []string
```

Args returns the whitespace separated arguments following a command.

### `ArgsText`

```go
func (c *Ctx) ArgsText() string
```

ArgsText returns everything after the command word, trimmed.

### `Send`

```go
func (c *Ctx) Send(text string, opts ...SendOptions) (*models.MessageID, error)
```

Send sends a message to the current chat.

### `Reply`

```go
func (c *Ctx) Reply(text string, opts ...SendOptions) (*models.MessageID, error)
```

Reply sends a message to the current chat as a reply to the incoming one.

### `ReplyKeypad`

```go
func (c *Ctx) ReplyKeypad(text string, kp *models.Keypad) (*models.MessageID, error)
```

ReplyKeypad replies with an inline keypad attached.

### `Edit`

```go
func (c *Ctx) Edit(messageID, text string, opts ...EditOptions) error
```

Edit replaces the text of a message of this chat (defaults to the
incoming message when messageID is empty).

### `Delete`

```go
func (c *Ctx) Delete() error { return c.Bot.DeleteMessage(c, c.ChatID(), c.MessageID()) }
```

Delete removes the incoming message.

### `DeleteMessage`

```go
func (c *Ctx) DeleteMessage(messageID string) error
```

DeleteMessage removes another message of the current chat.

### `Forward`

```go
func (c *Ctx) Forward(toChatID string) (*models.MessageID, error)
```

Forward forwards the incoming message to another chat.

### `CopyTo`

```go
func (c *Ctx) CopyTo(toChatID string, opts ...SendOptions) (*models.MessageID, error)
```

CopyTo re-sends the content of the incoming message to another chat
without a "forwarded from" header.

### `Chat`

```go
func (c *Ctx) Chat() (*models.Chat, error)
```

Chat returns (and caches) the chat object of the current chat.

### `Ban`

```go
func (c *Ctx) Ban() error { return c.Bot.BanChatMember(c, c.ChatID(), c.SenderID()) }
```

Ban removes the sender of the message from the chat.

### `Set`

```go
func (c *Ctx) Set(key string, v any)
```

Set stores a value for the lifetime of this update (handy to pass
data from middleware to handlers).

### `Get`

```go
func (c *Ctx) Get(key string) (any, bool)
```

Get returns a value stored with Set.

### `GetString`

```go
func (c *Ctx) GetString(key string) string
```

GetString returns a string value stored with Set ("" if absent).

### `CtxFrom`

```go
func CtxFrom(ctx context.Context, b *Bot, u *Update) *Ctx
```

CtxFrom returns the *Ctx created for this update by Handle, or builds
a fresh one. Middlewares can use it to Set values for the handler.

### `Handle`

```go
func (b *Bot) Handle(rs []Rule, fn HandlerFunc) int { return b.HandlePriority(0, rs, fn) }
```

Handle registers an ergonomic handler for updates matching ALL rules.

### `HandlePriority`

```go
func (b *Bot) HandlePriority(prio int, rs []Rule, fn HandlerFunc) int
```

HandlePriority is Handle with an explicit priority (see OnPriority).

### `Command`

```go
func (b *Bot) Command(name string, fn HandlerFunc, aliases ...string) int
```

Command handles "/name", "!name" (and any alias) case-insensitively.

### `Text`

```go
func (b *Bot) Text(fn HandlerFunc, exact ...string) int
```

Text handles messages whose text equals one of the given strings.

### `Regex`

```go
func (b *Bot) Regex(pattern string, fn HandlerFunc) int
```

Regex handles messages whose text matches the regular expression.

### `Button`

```go
func (b *Bot) Button(id string, fn HandlerFunc) int
```

Button handles presses of the inline button with the given ID.

### `Message`

```go
func (b *Bot) Message(fn HandlerFunc) int
```

Message handles every message that carries text.

### `Fallback`

```go
func (b *Bot) Fallback(fn HandlerFunc) int { return b.HandlePriority(PriorityFallback, nil, fn) }
```

Fallback registers a catch-all handler. It always runs after every
other handler, whatever the registration order.

### `CommandRule`

```go
func CommandRule(name string, aliases ...string) Rule
```

CommandRule builds the rule used by Command: "/name" or "!name"
(and aliases), case-insensitive.

### `RegexRule`

```go
func RegexRule(pattern string) Rule
```

RegexRule builds a text regular-expression rule; it panics on an
invalid pattern, like regexp.MustCompile.
