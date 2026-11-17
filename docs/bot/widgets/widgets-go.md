# `widgets` — Package widgets provides a builder DSL for the inline / chat

_Source: [`widgets.go`](../../../ParsRubika/bot/widgets/widgets.go)_

## Symbols

### `KeypadBuilder`

```go
type KeypadBuilder struct
```

KeypadBuilder builds a chat keypad (the persistent keyboard shown
under the text input). Each call to Row() starts a new row; each
button constructor appends to the current row.

### `NewChat`

```go
func NewChat() *KeypadBuilder { return &KeypadBuilder{resizeKeyboard: true} }
```

NewChat returns a fresh KeypadBuilder for chat keypads.

### `Row`

```go
func (k *KeypadBuilder) Row() *KeypadBuilder
```

Row finalises the current row and starts a new one. Calling Row()
before any button is added starts the first row.

### `AddButton`

```go
func (k *KeypadBuilder) AddButton(b models.Button) *KeypadBuilder
```

AddButton appends a button to the current row. Calling AddButton
before Row() implicitly starts the first row.

### `Resize`

```go
func (k *KeypadBuilder) Resize(v bool) *KeypadBuilder
```

Resize sets the resize_keyboard flag. Defaults to true.

### `OneTime`

```go
func (k *KeypadBuilder) OneTime(v bool) *KeypadBuilder
```

OneTime sets the one_time_keyboard flag. Defaults to false.

### `Build`

```go
func (k *KeypadBuilder) Build() *models.Keypad
```

Build returns the assembled keypad.

### `InlineBuilder`

```go
type InlineBuilder struct
```

InlineBuilder builds an inline keypad (the buttons attached under
a specific message). Same API as KeypadBuilder but the resize /
one_time flags are ignored.

### `NewInline`

```go
func NewInline() *InlineBuilder { return &InlineBuilder{} }
```

NewInline returns a fresh InlineBuilder.

### `Row`

```go
func (k *InlineBuilder) Row() *InlineBuilder
```

Row finalises the current row and starts a new one.

### `AddButton`

```go
func (k *InlineBuilder) AddButton(b models.Button) *InlineBuilder
```

AddButton appends a button to the current row.

### `Build`

```go
func (k *InlineBuilder) Build() *models.Keypad
```

Build returns the assembled keypad.

### `Simple`

```go
func Simple(id, text string) models.Button
```

Simple creates a "Simple" button.

### `Link`

```go
func Link(id, text, url string) models.Button
```

Link creates a "Link" button with a URL.

### `JoinChannel`

```go
func JoinChannel(id, text, username string) models.Button
```

JoinChannel creates a "Link" button that opens a channel.

### `OpenChat`

```go
func OpenChat(id, text, objectGUID, objectType string) models.Button
```

OpenChat creates a "Link" button that opens a chat.

### `Payment`

```go
func Payment(id, text string, amount int64, description string) models.Button
```

Payment creates a "Payment" button.

### `Calendar`

```go
func Calendar(id, title, calendarType string, defaultValue, minYear, maxYear string) models.Button
```

Calendar creates a "Calendar" button.

### `NumberPicker`

```go
func NumberPicker(id, title, min, max, defaultValue string) models.Button
```

NumberPicker creates a "NumberPicker" button.

### `StringPicker`

```go
func StringPicker(id, title string, items []string, defaultValue string) models.Button
```

StringPicker creates a "StringPicker" button.

### `Location`

```go
func Location(id, title, locationType string, lat, lon float64) models.Button
```

Location creates a "Location" button.

### `Textbox`

```go
func Textbox(id, title, lineType, keypadType, placeholder, defaultValue string) models.Button
```

Textbox creates a "Textbox" button.

### `Selection`

```go
func Selection(id, title, selectionID string, items []models.ButtonSelectionItem, isMulti bool, columnsCount string) models.Button
```

Selection creates a "Selection" button with a list of items.

### `CameraImage`

```go
func CameraImage(id, text string) models.Button
```

CameraImage creates a "CameraImage" button.

### `CameraVideo`

```go
func CameraVideo(id, text string) models.Button
```

CameraVideo creates a "CameraVideo" button.

### `GalleryImage`

```go
func GalleryImage(id, text string) models.Button
```

GalleryImage creates a "GalleryImage" button.

### `GalleryVideo`

```go
func GalleryVideo(id, text string) models.Button
```

GalleryVideo creates a "GalleryVideo" button.

### `File`

```go
func File(id, text string) models.Button
```

File creates a "File" button.

### `Audio`

```go
func Audio(id, text string) models.Button
```

Audio creates an "Audio" button.

### `RecordAudio`

```go
func RecordAudio(id, text string) models.Button
```

RecordAudio creates a "RecordAudio" button.

### `MyPhoneNumber`

```go
func MyPhoneNumber(id, text string) models.Button
```

MyPhoneNumber creates a "MyPhoneNumber" button.

### `MyLocation`

```go
func MyLocation(id, text string) models.Button
```

MyLocation creates a "MyLocation" button.

### `AskMyPhoneNumber`

```go
func AskMyPhoneNumber(id, text string) models.Button
```

AskMyPhoneNumber creates an "AskMyPhoneNumber" button.

### `AskLocation`

```go
func AskLocation(id, text string) models.Button
```

AskLocation creates an "AskLocation" button.

### `Barcode`

```go
func Barcode(id, text string) models.Button
```

Barcode creates a "Barcode" button.
