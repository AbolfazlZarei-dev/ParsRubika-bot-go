# `models` — Message

_Source: [`message.go`](../../../ParsRubika/bot/models/message.go)_

## Symbols

### `ForwardedFrom`

```go
type ForwardedFrom struct
```

ForwardedFrom carries the forward metadata.

### `PaymentStatus`

```go
type PaymentStatus struct
```

PaymentStatus reflects the payment state attached to a message.

### `PollStatus`

```go
type PollStatus struct
```

PollStatus holds the live state of a poll.

### `Poll`

```go
type Poll struct
```

Poll is a poll attachment.

### `Location`

```go
type Location struct
```

Location is a (lat, lon) pair.

### `LiveLocation`

```go
type LiveLocation struct
```

LiveLocation carries a live-location update.

### `ContactMessage`

```go
type ContactMessage struct
```

ContactMessage carries a contact vCard payload.

### `Sticker`

```go
type Sticker struct
```

Sticker is the sticker metadata.

### `AuxData`

```go
type AuxData struct
```

AuxData carries the aux_data attached to a button click.

### `JoinChannelData`

```go
type JoinChannelData struct
```

JoinChannelData is the payload attached to a "join channel" button.

### `OpenChatData`

```go
type OpenChatData struct
```

OpenChatData carries the open-chat payload attached to a button.

### `Message`

```go
type Message struct
```

Message is the parsed form of a `new_message` / `updated_message` payload.

### `MessageID`

```go
type MessageID struct
```

MessageID is the response of every send* call.

### `MetaDataParts`

```go
type MetaDataParts struct
```

MetaDataParts is one slice of the metadata attached to a text message.

### `Metadata`

```go
type Metadata struct
```

Metadata is the structured-form metadata block.
