# `forms` — Package forms is a ready-made plugin that provides a tiny

_Source: [`forms.go`](../../../ParsRubika/extras/forms/forms.go)_

## Symbols

### `Step`

```go
type Step func(ctx context.Context, b *bot.Bot, upd *models.Update, answer string) (next string, err error)
```

Step is a single step of a form. The callback receives the user's
answer and returns the next step name (or "" to finish).

### `Form`

```go
type Form struct
```

Form is a multi-step conversation scoped to a single chat. Steps
are registered by name; the first step is named "start".

### `New`

```go
func New() *Form
```

New constructs a fresh Form.

### `Add`

```go
func (f *Form) Add(name string, step Step)
```

Add registers a step under the supplied name.

### `Start`

```go
func (f *Form) Start(chatID string)
```

Start begins the form for a chat at the "start" step.

### `IsActive`

```go
func (f *Form) IsActive(chatID string) bool
```

IsActive reports whether a form is in progress for a chat.

### `Stop`

```go
func (f *Form) Stop(chatID string)
```

Stop ends the form for a chat.

### `Handle`

```go
func (f *Form) Handle(ctx context.Context, b *bot.Bot, upd *models.Update) error
```

Handle is the callback you register on the bot for any message
that arrives while a form is active.
