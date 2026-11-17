# `models` — Bot

_Source: [`bot.go`](../../../ParsRubika/bot/models/bot.go)_

## Symbols

### `Bot`

```go
type Bot struct
```

Bot is the bot identity returned by `getMe`.

### `BotCommand`

```go
type BotCommand struct
```

BotCommand is a single entry in the bot's command list.

### `UpdateType`

```go
type UpdateType string
```

UpdateType enumerates the update types returned by `getUpdates`.

### `UpdateTypeNewMessage`

```go
const UpdateTypeNewMessage UpdateType = "NewMessage"
```

### `UpdateTypeUpdatedMessage`

```go
const UpdateTypeUpdatedMessage UpdateType = "UpdatedMessage"
```

### `UpdateTypeRemovedMessage`

```go
const UpdateTypeRemovedMessage UpdateType = "RemovedMessage"
```

### `UpdateTypeStartedBot`

```go
const UpdateTypeStartedBot UpdateType = "StartedBot"
```

### `UpdateTypeStoppedBot`

```go
const UpdateTypeStoppedBot UpdateType = "StoppedBot"
```

### `UpdateTypeUpdatedPayment`

```go
const UpdateTypeUpdatedPayment UpdateType = "UpdatedPayment"
```

### `UpdateTypeInlineMessage`

```go
const UpdateTypeInlineMessage UpdateType = "InlineMessage"
```

### `ChatKeypadType`

```go
type ChatKeypadType string
```

ChatKeypadType selects how a chat keypad is applied (None, New, Remove).

### `ChatKeypadNone`

```go
const ChatKeypadNone ChatKeypadType = "None"
```

### `ChatKeypadNew`

```go
const ChatKeypadNew ChatKeypadType = "New"
```

### `ChatKeypadRemove`

```go
const ChatKeypadRemove ChatKeypadType = "Remove"
```

### `FileType`

```go
type FileType string
```

FileType enumerates the supported file upload types.

### `FileTypeFile`

```go
const FileTypeFile FileType = "File"
```

### `FileTypeImage`

```go
const FileTypeImage FileType = "Image"
```

### `FileTypeVoice`

```go
const FileTypeVoice FileType = "Voice"
```

### `FileTypeMusic`

```go
const FileTypeMusic FileType = "Music"
```

### `FileTypeGif`

```go
const FileTypeGif FileType = "Gif"
```

### `FileTypeVideo`

```go
const FileTypeVideo FileType = "Video"
```

### `ButtonType`

```go
type ButtonType string
```

ButtonType enumerates the supported button types.

### `ButtonTypeSimple`

```go
const ButtonTypeSimple ButtonType = "Simple"
```

### `ButtonTypeSelection`

```go
const ButtonTypeSelection ButtonType = "Selection"
```

### `ButtonTypeCalendar`

```go
const ButtonTypeCalendar ButtonType = "Calendar"
```

### `ButtonTypeNumberPicker`

```go
const ButtonTypeNumberPicker ButtonType = "NumberPicker"
```

### `ButtonTypeStringPicker`

```go
const ButtonTypeStringPicker ButtonType = "StringPicker"
```

### `ButtonTypeLocation`

```go
const ButtonTypeLocation ButtonType = "Location"
```

### `ButtonTypePayment`

```go
const ButtonTypePayment ButtonType = "Payment"
```

### `ButtonTypeCameraImage`

```go
const ButtonTypeCameraImage ButtonType = "CameraImage"
```

### `ButtonTypeCameraVideo`

```go
const ButtonTypeCameraVideo ButtonType = "CameraVideo"
```

### `ButtonTypeGalleryImage`

```go
const ButtonTypeGalleryImage ButtonType = "GalleryImage"
```

### `ButtonTypeGalleryVideo`

```go
const ButtonTypeGalleryVideo ButtonType = "GalleryVideo"
```

### `ButtonTypeFile`

```go
const ButtonTypeFile ButtonType = "File"
```

### `ButtonTypeAudio`

```go
const ButtonTypeAudio ButtonType = "Audio"
```

### `ButtonTypeRecordAudio`

```go
const ButtonTypeRecordAudio ButtonType = "RecordAudio"
```

### `ButtonTypeMyPhone`

```go
const ButtonTypeMyPhone ButtonType = "MyPhoneNumber"
```

### `ButtonTypeMyLocation`

```go
const ButtonTypeMyLocation ButtonType = "MyLocation"
```

### `ButtonTypeTextbox`

```go
const ButtonTypeTextbox ButtonType = "Textbox"
```

### `ButtonTypeLink`

```go
const ButtonTypeLink ButtonType = "Link"
```

### `ButtonTypeAskMyPhone`

```go
const ButtonTypeAskMyPhone ButtonType = "AskMyPhoneNumber"
```

### `ButtonTypeAskLocation`

```go
const ButtonTypeAskLocation ButtonType = "AskLocation"
```

### `ButtonTypeBarcode`

```go
const ButtonTypeBarcode ButtonType = "Barcode"
```

### `UpdateEndpointType`

```go
type UpdateEndpointType string
```

UpdateEndpointType names the webhook endpoint kinds accepted by
`updateBotEndpoints`.

### `EndpointReceiveUpdate`

```go
const EndpointReceiveUpdate UpdateEndpointType = "ReceiveUpdate"
```

### `EndpointReceiveInlineMessage`

```go
const EndpointReceiveInlineMessage UpdateEndpointType = "ReceiveInlineMessage"
```

### `EndpointReceiveQuery`

```go
const EndpointReceiveQuery UpdateEndpointType = "ReceiveQuery"
```

### `EndpointGetSelectionItem`

```go
const EndpointGetSelectionItem UpdateEndpointType = "GetSelectionItem"
```

### `EndpointSearchSelectionItems`

```go
const EndpointSearchSelectionItems UpdateEndpointType = "SearchSelectionItems"
```
