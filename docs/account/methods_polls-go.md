# `account` — Methods polls

_Source: [`methods_polls.go`](../../ParsRubika/account/methods_polls.go)_

## Symbols

### `CreatePoll`

```go
func (a *Account) CreatePoll(ctx context.Context, objectGUID, question string, options []string, pollType string, isAnonymous, allowMulti bool, correctIdx *int, explanation, replyToMessageID string) (*events.Event, error)
```

CreatePoll creates and sends a poll message.

### `VotePoll`

```go
func (a *Account) VotePoll(ctx context.Context, pollID string, selectionIndex int) (*events.Event, error)
```

VotePoll votes on a single poll option. The API expects a single
`selection_index: int`, not a list.

### `GetPollStatus`

```go
func (a *Account) GetPollStatus(ctx context.Context, pollID string) (*events.Event, error)
```

GetPollStatus fetches the status of a poll. The API only needs
`{poll_id}` — no `object_guid`.

### `GetPollOptionVoters`

```go
func (a *Account) GetPollOptionVoters(ctx context.Context, pollID string, optionIndex int, startID string) (*events.Event, error)
```

GetPollOptionVoters returns the voters of a poll option, paginated
by `startID` (the API uses `start_id`, not `cursor`).

### `SetPollAction`

```go
func (a *Account) SetPollAction(ctx context.Context, pollID, action string) (*events.Event, error)
```

SetPollAction stops or retracts a poll.

### `StopPoll`

```go
func (a *Account) StopPoll(ctx context.Context, pollID string) (*events.Event, error)
```

StopPoll stops a poll (action = "Stop").

### `RetractPoll`

```go
func (a *Account) RetractPoll(ctx context.Context, pollID string) (*events.Event, error)
```

RetractPoll retracts a poll (action = "Retract").
