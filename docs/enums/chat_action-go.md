# `enums` — Chat action

_Source: [`chat_action.go`](../../ParsRubika/enums/chat_action.go)_

## Symbols

### `ChatAction`

```go
type ChatAction string
```

ChatAction is the activity type reported during typing / uploading /
recording.

### `ChatActionTyping`

```go
const ChatActionTyping ChatAction = "Typing"
```

### `ChatActionUploading`

```go
const ChatActionUploading ChatAction = "Uploading"
```

### `ChatActionRecording`

```go
const ChatActionRecording ChatAction = "Recording"
```

### `String`

```go
func (c ChatAction) String() string { return string(c) }
```

String implements fmt.Stringer.
