# `events` — Package events defines the dynamic Event wrapper used across the

_Source: [`event.go`](../../ParsRubika/events/event.go)_

## Symbols

### `Event`

```go
type Event struct
```

Event wraps the raw JSON object returned by the server.

### `New`

```go
func New(m map[string]any) *Event
```

New wraps an existing map. The map is NOT copied.

Optimisation: External is left nil and only allocated on first
SetExternal call. The vast majority of events never use the scratch
space, so this saves one map allocation per event (~120 B/op on the
benchmark). Benchmarked: BenchmarkFindKey/BenchmarkStr stay at 0
allocs/op — they don't touch External — while New itself drops one
allocation.

### `JSON`

```go
func (e *Event) JSON(indent string) string
```

JSON pretty-prints the raw payload.

### `String`

```go
func (e *Event) String() string { return e.JSON(" ") }
```

String implements fmt.Stringer.

### `Get`

```go
func (e *Event) Get(key string) (any, bool)
```

Get returns the value stored at `key` in the top level.

### `Set`

```go
func (e *Event) Set(key string, value any)
```

Set writes a value into the top level.

### `SetExternal`

```go
func (e *Event) SetExternal(key string, value any)
```

SetExternal attaches an extra value to the event's scratch space.

### `GetExternal`

```go
func (e *Event) GetExternal(key string) (any, bool)
```

GetExternal retrieves a value from the scratch space.

### `Find`

```go
func (e *Event) Find(keys ...string) any
```

Find recursively searches the event for the given keys. Returns nil
if none of them are present at any level.

### `Str`

```go
func (e *Event) Str(keys ...string) string
```

Str finds a string under any of `keys` (returns "" if absent).

### `AnyStr`

```go
func (e *Event) AnyStr(keys ...string) string
```

AnyStr finds a value under any of `keys` and returns it as a string.
Unlike Str, this also coerces numeric values (int, int64, float64,
json.Number) to their decimal string form. Use this when the
server may return a field as either a string or a number (for
example `file_id` from `requestSendFile`).

### `Int`

```go
func (e *Event) Int(keys ...string) int
```

Int finds a numeric value under any of `keys` (returns 0 if absent).

In addition to the obvious int / int64 / float64 / json.Number
types, this also tries to parse STRING values as decimal integers.
This is needed because some server fields (notably `dc_id` returned
by `requestSendFile`) come back as JSON strings like `"480"` instead
of bare numbers — without the string case we would silently return
0 and then the server would reject the subsequent request with
INVALID_INPUT.

### `Int64`

```go
func (e *Event) Int64(keys ...string) int64
```

Int64 finds a numeric value as int64.

Like `Int`, this also accepts STRING values that parse as decimal
integers — see the comment on `Int` for why this matters.

### `Bool`

```go
func (e *Event) Bool(keys ...string) bool
```

Bool finds a boolean under any of `keys`.

### `Map`

```go
func (e *Event) Map(keys ...string) map[string]any
```

Map finds a sub-dict under any of `keys`.

### `EventOf`

```go
func (e *Event) EventOf(keys ...string) *Event
```

EventOf finds a *Event under any of `keys`.

### `Slice`

```go
func (e *Event) Slice(keys ...string) []any
```

Slice finds a list under any of `keys`.

### `Message`

```go
func (e *Event) Message() *Event { return e.EventOf("message") }
```

Message returns the `message` sub-object as an Event (or nil).

### `Text`

```go
func (e *Event) Text() string
```

Text returns the message text (if any).

### `Type`

```go
func (e *Event) Type() string { return e.Str("type", "author_type") }
```

Type returns the update type / author type.

### `IsMe`

```go
func (e *Event) IsMe(selfGUID string) bool { return e.AuthorGUID() == selfGUID }
```

IsMe reports whether the message author is the current user.

### `IsEdited`

```go
func (e *Event) IsEdited() bool
```

IsEdited reports whether the message has the `is_edited` flag.

### `IsForward`

```go
func (e *Event) IsForward() bool
```

IsForward reports whether the message has a `forwarded_from` block.

### `IsEvent`

```go
func (e *Event) IsEvent() bool
```

IsEvent reports whether the message type is "Event".

### `IsFileInline`

```go
func (e *Event) IsFileInline() bool
```

IsFileInline reports whether the message contains a file_inline block.

### `FileInline`

```go
func (e *Event) FileInline() *Event { return e.EventOf("file_inline") }
```

FileInline returns the file_inline block (if any).

### `IsGroup`

```go
func (e *Event) IsGroup() bool { return e.Type() == "Group" }
```

IsGroup / IsChannel / IsPrivate reflect the update's chat type.

### `IsChannel`

```go
func (e *Event) IsChannel() bool { return e.Type() == "Channel" }
```

### `IsPrivate`

```go
func (e *Event) IsPrivate() bool { return e.Type() == "User" }
```

### `ObjectGUID`

```go
func (e *Event) ObjectGUID() string
```

ObjectGUID returns the chat GUID, whichever name the server used.

### `AuthorGUID`

```go
func (e *Event) AuthorGUID() string
```

AuthorGUID returns the message author's GUID.

### `MessageID`

```go
func (e *Event) MessageID() string
```

MessageID returns the message id (or pinned_message_id if that's set).

### `ReplyMessageID`

```go
func (e *Event) ReplyMessageID() string
```

ReplyMessageID returns the replied-to message id (if any).

### `Sticker`

```go
func (e *Event) Sticker() *Event { return e.EventOf("sticker") }
```

Sticker returns the sticker sub-object (if any).

### `GUIDType`

```go
func GUIDType(guid string) string
```

GUIDType returns "Channel", "Group", or "User" depending on the prefix
of the supplied GUID (c0*, g0*, anything else).
