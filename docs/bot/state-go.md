# `bot` — State

_Source: [`state.go`](../../ParsRubika/bot/state.go)_

## Symbols

### `StateManager`

```go
type StateManager struct
```

StateManager holds per-user state values with optional TTL.

### `NewStateManager`

```go
func NewStateManager(ttl time.Duration) *StateManager
```

NewStateManager creates a StateManager. If ttl > 0, state entries
expire after that duration of inactivity.

### `Set`

```go
func (sm *StateManager) Set(userID, state string)
```

Set sets the state for userID.

### `Get`

```go
func (sm *StateManager) Get(userID string) string
```

Get returns the current state for userID ("" if none).

### `Clear`

```go
func (sm *StateManager) Clear(userID string)
```

Clear removes the state for userID.

### `SetData`

```go
func (sm *StateManager) SetData(userID, key, value string)
```

SetData stores a key/value pair in the per-user state scratch area.
The state is created if it does not exist.

### `GetData`

```go
func (sm *StateManager) GetData(userID, key string) string
```

GetData retrieves a key from the per-user state scratch area.

### `HasState`

```go
func (sm *StateManager) HasState(userID string) bool
```

HasState returns true if the user has any state set (and not expired).

### `ListAll`

```go
func (sm *StateManager) ListAll() map[string]string
```

ListAll returns a snapshot of all active states (for debugging).

### `Close`

```go
func (sm *StateManager) Close() { sm.once.Do(func() { close(sm.done) }) }
```

Close stops the background expiry goroutine. It is safe to call more
than once.

### `StateRule`

```go
type StateRule struct
```

StateRule is a Rule that matches when the user is in (or not in) a
given state.

### `StateScope`

```go
type StateScope int
```

StateScope determines whether the state key is the sender_id (User)
or the chat_id (Chat) or a combination (Both).

### `StateScopeUser`

```go
const StateScopeUser StateScope = iota
```

### `StateScopeChat`

```go
const StateScopeChat
```

### `StateScopeBoth`

```go
const StateScopeBoth
```

### `NewStateRule`

```go
func NewStateRule(manager *StateManager, scope StateScope, invert bool, states ...string) *StateRule
```

NewStateRule creates a rule that matches when the user is in one of
the given states. Pass invert=true to match when NOT in any of them.

### `Check`

```go
func (s *StateRule) Check(ctx context.Context, upd *models.Update) bool
```

Check implements rules.Rule.

### `Conversation`

```go
type Conversation struct
```

Conversation is a multi-step flow tied to a single user. It tracks the
current step name and optional per-step data.

### `ConversationStep`

```go
type ConversationStep struct
```

ConversationStep is a single step in a conversation flow.

### `ConversationManager`

```go
type ConversationManager struct
```

ConversationManager owns all active conversations and their step
definitions.

### `NewConversationManager`

```go
func NewConversationManager(stateMgr *StateManager) *ConversationManager
```

NewConversationManager creates a new ConversationManager.

### `RegisterFlow`

```go
func (cm *ConversationManager) RegisterFlow(name string, steps []ConversationStep)
```

RegisterFlow registers a named conversation flow with its steps.

### `Start`

```go
func (cm *ConversationManager) Start(userID, flowName string) (string, error)
```

Start begins a conversation flow for userID. The first step's prompt
is returned.

### `Handle`

```go
func (cm *ConversationManager) Handle(userID, input string) (string, error)
```

Handle processes an input for an active conversation. Returns the next
step's prompt (or "" if the conversation is complete) and an error.

### `IsActive`

```go
func (cm *ConversationManager) IsActive(userID string) bool
```

IsActive returns true if the user has an active conversation.

### `Stop`

```go
func (cm *ConversationManager) Stop(userID string)
```

Stop forcibly ends a conversation.
