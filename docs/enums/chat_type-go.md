# `enums` — Chat type

_Source: [`chat_type.go`](../../ParsRubika/enums/chat_type.go)_

## Symbols

### `ChatType`

```go
type ChatType string
```

ChatType enumerates the kinds of chat the server can return.

### `ChatTypeUser`

```go
const ChatTypeUser ChatType = "User"
```

### `ChatTypeBot`

```go
const ChatTypeBot ChatType = "Bot"
```

### `ChatTypeGroup`

```go
const ChatTypeGroup ChatType = "Group"
```

### `ChatTypeChannel`

```go
const ChatTypeChannel ChatType = "Channel"
```

### `ChatTypeService`

```go
const ChatTypeService ChatType = "Service"
```

### `String`

```go
func (c ChatType) String() string { return string(c) }
```

String implements fmt.Stringer.

### `ChatTypeFromGUID`

```go
func ChatTypeFromGUID(guid string) ChatType
```

FromGUID infers the chat type from the prefix of the supplied GUID:

- "c0..." → Channel
- "g0..." → Group
- anything else (including "u0...", "b0...") → User
