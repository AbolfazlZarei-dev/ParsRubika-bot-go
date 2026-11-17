# `models` — Keypad

_Source: [`keypad.go`](../../../ParsRubika/bot/models/keypad.go)_

## Symbols

### `ButtonSelectionItem`

```go
type ButtonSelectionItem struct
```

ButtonSelectionItem is one item shown inside a selection button.

### `ButtonSelection`

```go
type ButtonSelection struct
```

ButtonSelection describes the `button_selection` field of a button.

### `ButtonCalendar`

```go
type ButtonCalendar struct
```

ButtonCalendar is the calendar button.

### `ButtonNumberPicker`

```go
type ButtonNumberPicker struct
```

ButtonNumberPicker is the numeric picker button.

### `ButtonStringPicker`

```go
type ButtonStringPicker struct
```

ButtonStringPicker is the string-list picker button.

### `ButtonTextbox`

```go
type ButtonTextbox struct
```

ButtonTextbox is the textbox button.

### `ButtonLocation`

```go
type ButtonLocation struct
```

ButtonLocation is the location-picker button.

### `ButtonLink`

```go
type ButtonLink struct
```

ButtonLink wraps a hyperlink or join-channel / open-chat button.

### `NormalizeLinkURL`

```go
func (b *ButtonLink) NormalizeLinkURL()
```

NormalizeLinkURL converts the rubika.ir short links into the
deep-link form expected by the server.

### `Button`

```go
type Button struct
```

Button is a single key inside a keypad row.

### `KeypadRow`

```go
type KeypadRow struct
```

KeypadRow is a row of buttons.

### `Keypad`

```go
type Keypad struct
```

Keypad is the full keypad object attached to a message or chat.

### `MessageKeypadUpdate`

```go
type MessageKeypadUpdate struct
```

MessageKeypadUpdate is the payload for `editMessageKeypad`.
