# `rubiai` — Sendmessage

_Source: [`sendmessage.go`](../../ParsRubika/rubiai/sendmessage.go)_

## Symbols

### `SendMessageOptions`

```go
type SendMessageOptions struct
```

SendMessageOptions configures an AI sendMessage call.

### `SendMessageResult`

```go
type SendMessageResult struct
```

SendMessageResult is the assembled outcome of an SSE-streamed
sendMessage call. The TextContent field holds the concatenated
assistant reply (incrementally delivered via updateMessage events).

### `SendMessage`

```go
func (r *RubiaAI) SendMessage(ctx context.Context, opts SendMessageOptions) (*SendMessageResult, error)
```

SendMessage sends a message to the AI service and returns the fully
assembled result. It blocks until the SSE stream terminates.

For incremental/token-by-token output use SendMessageStream, which
invokes a callback for each parsed event as it arrives.

### `SendMessageStream`

```go
func (r *RubiaAI) SendMessageStream(ctx context.Context, opts SendMessageOptions, onEvent func(event string, data map[string]any)) (*SendMessageResult, error)
```

SendMessageStream is like SendMessage but invokes `onEvent` for each
parsed SSE event as soon as it arrives (event name + parsed JSON data).
The callback is invoked from the goroutine reading the response body,
so it must not block on the same context. Pass nil to behave like
SendMessage. The returned result is the assembled outcome.
