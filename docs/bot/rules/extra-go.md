# `rules` — Extra

_Source: [`extra.go`](../../../ParsRubika/bot/rules/extra.go)_

## Symbols

### `TextStartsWith`

```go
type TextStartsWith struct
```

TextStartsWith matches messages whose text starts with `prefix`.

### `Check`

```go
func (t TextStartsWith) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextEndsWith`

```go
type TextEndsWith struct
```

TextEndsWith matches messages whose text ends with `suffix`.

### `Check`

```go
func (t TextEndsWith) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextContains`

```go
type TextContains struct
```

TextContains matches messages whose text contains `substring`
(case-insensitive by default).

### `Check`

```go
func (t TextContains) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextEquals`

```go
type TextEquals struct
```

TextEquals matches messages whose text exactly equals `value`.

### `Check`

```go
func (t TextEquals) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextNotEquals`

```go
type TextNotEquals struct
```

TextNotEquals matches messages whose text does NOT equal `value`.

### `Check`

```go
func (t TextNotEquals) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextLength`

```go
type TextLength struct
```

TextLength matches messages whose text length is within [min, max].

### `Check`

```go
func (t TextLength) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextWordCount`

```go
type TextWordCount struct
```

TextWordCount matches messages whose word count is within [min, max].

### `Check`

```go
func (t TextWordCount) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextIsUpper`

```go
type TextIsUpper struct
```

TextIsUpper matches messages whose text is all uppercase.

### `Check`

```go
func (TextIsUpper) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextIsLower`

```go
type TextIsLower struct
```

TextIsLower matches messages whose text is all lowercase.

### `Check`

```go
func (TextIsLower) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextIsDigit`

```go
type TextIsDigit struct
```

TextIsDigit matches messages whose text is all digits.

### `Check`

```go
func (TextIsDigit) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `TextRegex`

```go
type TextRegex struct
```

TextRegex matches messages whose text matches a regex.

### `NewTextRegex`

```go
func NewTextRegex(pattern string) *TextRegex
```

NewTextRegex constructs a TextRegex rule.

### `Check`

```go
func (t *TextRegex) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `FileSizeGT`

```go
type FileSizeGT struct
```

FileSizeGT matches messages whose file size is greater than `size` bytes.

### `Check`

```go
func (f FileSizeGT) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `FileSizeLT`

```go
type FileSizeLT struct
```

FileSizeLT matches messages whose file size is less than `size` bytes.

### `Check`

```go
func (f FileSizeLT) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `FileNameContains`

```go
type FileNameContains struct
```

FileNameContains matches messages whose file name contains `substring`.

### `Check`

```go
func (f FileNameContains) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `FileExtension`

```go
type FileExtension struct
```

FileExtension matches messages whose file name ends with `ext`.

### `Check`

```go
func (f FileExtension) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `FileIDIs`

```go
type FileIDIs struct
```

FileIDIs matches messages whose file_id equals `id`.

### `Check`

```go
func (f FileIDIs) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `StickerIDIs`

```go
type StickerIDIs struct
```

StickerIDIs matches messages whose sticker_id equals `id`.

### `Check`

```go
func (s StickerIDIs) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `StickerEmojiIs`

```go
type StickerEmojiIs struct
```

StickerEmojiIs matches messages whose sticker emoji equals `emoji`.

### `Check`

```go
func (s StickerEmojiIs) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ChatTitleContains`

```go
type ChatTitleContains struct
```

ChatTitleContains matches when the chat title contains `keyword`.
Note: this requires the chat info to be fetched first, which the
rule does lazily.

### `Check`

```go
func (c ChatTitleContains) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ChatTypeIs`

```go
type ChatTypeIs struct
```

ChatTypeIs matches when the update came from a chat of the given type.
`chatType` is "user", "private", "group", "channel", or "bot".

### `Check`

```go
func (c ChatTypeIs) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ForwardedFromUser`

```go
type ForwardedFromUser struct
```

ForwardedFromUser matches messages forwarded from a user.

### `Check`

```go
func (ForwardedFromUser) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ForwardedFromChannel`

```go
type ForwardedFromChannel struct
```

ForwardedFromChannel matches messages forwarded from a channel.

### `Check`

```go
func (ForwardedFromChannel) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ForwardedFromBot`

```go
type ForwardedFromBot struct
```

ForwardedFromBot matches messages forwarded from a bot.

### `Check`

```go
func (ForwardedFromBot) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `ForwardedNoLink`

```go
type ForwardedNoLink struct
```

ForwardedNoLink matches messages forwarded without a link to the
original sender.

### `Check`

```go
func (ForwardedNoLink) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsEventData`

```go
type IsEventData struct
```

IsEventData matches EventData updates (BotJoined, BotRemoved, etc.).

### `Check`

```go
func (e IsEventData) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `BotJoined`

```go
type BotJoined struct
```

BotJoined matches BotJoined events.

### `Check`

```go
func (BotJoined) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `BotRemoved`

```go
type BotRemoved struct
```

BotRemoved matches BotRemoved events.

### `Check`

```go
func (BotRemoved) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `BotPermissionsChanged`

```go
type BotPermissionsChanged struct
```

BotPermissionsChanged matches BotPermissionsChanged events.

### `Check`

```go
func (BotPermissionsChanged) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `MetadataType`

```go
type MetadataType struct
```

MetadataType matches messages that contain metadata of the given type(s).

### `Check`

```go
func (m MetadataType) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsBold`

```go
type IsBold struct
```

IsBold matches messages that contain bold formatting.

### `Check`

```go
func (IsBold) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsItalic`

```go
type IsItalic struct
```

IsItalic matches messages that contain italic formatting.

### `Check`

```go
func (IsItalic) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsStrike`

```go
type IsStrike struct
```

IsStrike matches messages that contain strike formatting.

### `Check`

```go
func (IsStrike) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsUnderline`

```go
type IsUnderline struct
```

IsUnderline matches messages that contain underline formatting.

### `Check`

```go
func (IsUnderline) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsQuote`

```go
type IsQuote struct
```

IsQuote matches messages that contain a quote block.

### `Check`

```go
func (IsQuote) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsSpoiler`

```go
type IsSpoiler struct
```

IsSpoiler matches messages that contain spoiler formatting.

### `Check`

```go
func (IsSpoiler) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsPre`

```go
type IsPre struct
```

IsPre matches messages that contain a pre-formatted code block.

### `Check`

```go
func (IsPre) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.

### `IsMono`

```go
type IsMono struct
```

IsMono matches messages that contain inline code formatting.

### `Check`

```go
func (IsMono) Check(_ context.Context, upd *models.Update) bool
```

Check implements Rule.
