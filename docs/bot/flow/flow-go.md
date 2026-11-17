# `flow` — Package flow provides a conversation flow manager for the bot

_Source: [`flow.go`](../../../ParsRubika/bot/flow/flow.go)_

## Symbols

### `Bot`

```go
type Bot interface
```

Bot is the minimal interface a flow client needs.

### `Step`

```go
type Step func(ctx context.Context, b Bot, chatID, userID, input string) (next string, err error)
```

Step is the signature of a single flow step.

### `Flow`

```go
type Flow struct
```

Flow is a named multi-step conversation.

### `Manager`

```go
type Manager struct
```

Manager tracks all active flows across all users.

### `NewManager`

```go
func NewManager() *Manager
```

NewManager constructs a fresh flow Manager.

### `Register`

```go
func (m *Manager) Register(name string) *Flow
```

Register creates a new flow under the given name and returns it
so the caller can add steps.

### `Add`

```go
func (f *Flow) Add(stepName string, step Step)
```

Add registers a step in the flow.

### `Start`

```go
func (m *Manager) Start(flowName, userID string)
```

Start begins the flow for a user at the "start" step.

### `IsActive`

```go
func (m *Manager) IsActive(flowName, userID string) bool
```

IsActive reports whether a flow is in progress for a user.

### `IsActiveAny`

```go
func (m *Manager) IsActiveAny(userID string) bool
```

IsActiveAny reports whether any flow is in progress for a user.

### `Stop`

```go
func (m *Manager) Stop(flowName, userID string)
```

Stop ends any active flow for a user.

### `StopAll`

```go
func (m *Manager) StopAll(userID string)
```

StopAll ends all flows for a user.

### `GetData`

```go
func (m *Manager) GetData(flowName, userID, key string) string
```

GetData retrieves a value from the user's flow scratch space.

### `SetData`

```go
func (m *Manager) SetData(flowName, userID, key, value string)
```

SetData stores a value in the user's flow scratch space.

### `Handle`

```go
func (m *Manager) Handle(ctx context.Context, b Bot, chatID, userID, input string) (bool, error)
```

Handle processes an incoming message for a user. If the user has
an active flow, the current step is invoked. Returns true when
the message was consumed by a flow, false otherwise.
