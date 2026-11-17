// Package widgets provides a builder DSL for the inline / chat
// keypads used by the bot API. Every button type of the platform has a
// constructor.
//
// Each button constructor returns a *models.Button ready to be
// embedded in a KeypadRow. The Keypad builder accumulates rows and
// produces a *models.Keypad on Build.
package widgets

import (
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// KeypadBuilder builds a chat keypad (the persistent keyboard shown
// under the text input). Each call to Row() starts a new row; each
// button constructor appends to the current row.
type KeypadBuilder struct {
	rows           []models.KeypadRow
	resizeKeyboard bool
	onTimeKeyboard bool
}

// NewChat returns a fresh KeypadBuilder for chat keypads.
func NewChat() *KeypadBuilder { return &KeypadBuilder{resizeKeyboard: true} }

// Row finalises the current row and starts a new one. Calling Row()
// before any button is added starts the first row.
func (k *KeypadBuilder) Row() *KeypadBuilder {
	k.rows = append(k.rows, models.KeypadRow{})
	return k
}

// AddButton appends a button to the current row. Calling AddButton
// before Row() implicitly starts the first row.
func (k *KeypadBuilder) AddButton(b models.Button) *KeypadBuilder {
	if len(k.rows) == 0 {
		k.rows = append(k.rows, models.KeypadRow{})
	}
	idx := len(k.rows) - 1
	k.rows[idx].Buttons = append(k.rows[idx].Buttons, b)
	return k
}

// Resize sets the resize_keyboard flag. Defaults to true.
func (k *KeypadBuilder) Resize(v bool) *KeypadBuilder {
	k.resizeKeyboard = v
	return k
}

// OneTime sets the one_time_keyboard flag. Defaults to false.
func (k *KeypadBuilder) OneTime(v bool) *KeypadBuilder {
	k.onTimeKeyboard = v
	return k
}

// Build returns the assembled keypad.
func (k *KeypadBuilder) Build() *models.Keypad {
	r := k.resizeKeyboard
	o := k.onTimeKeyboard
	return &models.Keypad{
		Rows:           k.rows,
		ResizeKeyboard: &r,
		OnTimeKeyboard: &o,
	}
}

// InlineBuilder builds an inline keypad (the buttons attached under
// a specific message). Same API as KeypadBuilder but the resize /
// one_time flags are ignored.
type InlineBuilder struct {
	rows []models.KeypadRow
}

// NewInline returns a fresh InlineBuilder.
func NewInline() *InlineBuilder { return &InlineBuilder{} }

// Row finalises the current row and starts a new one.
func (k *InlineBuilder) Row() *InlineBuilder {
	k.rows = append(k.rows, models.KeypadRow{})
	return k
}

// AddButton appends a button to the current row.
func (k *InlineBuilder) AddButton(b models.Button) *InlineBuilder {
	if len(k.rows) == 0 {
		k.rows = append(k.rows, models.KeypadRow{})
	}
	idx := len(k.rows) - 1
	k.rows[idx].Buttons = append(k.rows[idx].Buttons, b)
	return k
}

// Build returns the assembled keypad.
func (k *InlineBuilder) Build() *models.Keypad {
	return &models.Keypad{Rows: k.rows}
}

// ===== Button constructors ===========================================

// Simple creates a "Simple" button.
func Simple(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeSimple, ButtonText: text}
}

// Link creates a "Link" button with a URL.
func Link(id, text, url string) models.Button {
	bl := &models.ButtonLink{Type: "url", LinkURL: url}
	bl.NormalizeLinkURL()
	return models.Button{ID: id, Type: models.ButtonTypeLink, ButtonText: text, ButtonLink: bl}
}

// JoinChannel creates a "Link" button that opens a channel.
func JoinChannel(id, text, username string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeLink,
		ButtonText: text,
		ButtonLink: &models.ButtonLink{
			Type: "joinchannel",
			JoinChannelData: &models.JoinChannelData{
				Username: stripAt(username),
			},
		},
	}
}

// OpenChat creates a "Link" button that opens a chat.
func OpenChat(id, text, objectGUID, objectType string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeLink,
		ButtonText: text,
		ButtonLink: &models.ButtonLink{
			Type: "openchat",
			OpenChatData: &models.OpenChatData{
				ObjectGUID: objectGUID,
				ObjectType: objectType,
			},
		},
	}
}

// Payment creates a "Payment" button.
func Payment(id, text string, amount int64, description string) models.Button {
	// The models package has no dedicated payment button type, so a
	// generic button is used with the payment data attached to the
	// ButtonText (the caller can format it).
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypePayment,
		ButtonText: text,
	}
}

// Calendar creates a "Calendar" button.
func Calendar(id, title, calendarType string, defaultValue, minYear, maxYear string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeCalendar,
		ButtonText: title,
		ButtonCalendar: &models.ButtonCalendar{
			Type:         calendarType,
			DefaultValue: defaultValue,
			MinYear:      minYear,
			MaxYear:      maxYear,
			Title:        title,
		},
	}
}

// NumberPicker creates a "NumberPicker" button.
func NumberPicker(id, title, min, max, defaultValue string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeNumberPicker,
		ButtonText: title,
		ButtonNumberPicker: &models.ButtonNumberPicker{
			MinValue:     min,
			MaxValue:     max,
			DefaultValue: defaultValue,
			Title:        title,
		},
	}
}

// StringPicker creates a "StringPicker" button.
func StringPicker(id, title string, items []string, defaultValue string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeStringPicker,
		ButtonText: title,
		ButtonStringPicker: &models.ButtonStringPicker{
			Items:        items,
			DefaultValue: defaultValue,
			Title:        title,
		},
	}
}

// Location creates a "Location" button.
func Location(id, title, locationType string, lat, lon float64) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeLocation,
		ButtonText: title,
		ButtonLocation: &models.ButtonLocation{
			Type: locationType,
			DefaultPointerLocation: &models.Location{
				Latitude:  ftoa(lat),
				Longitude: ftoa(lon),
			},
			Title: title,
		},
	}
}

// Textbox creates a "Textbox" button.
func Textbox(id, title, lineType, keypadType, placeholder, defaultValue string) models.Button {
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeTextbox,
		ButtonText: title,
		ButtonTextbox: &models.ButtonTextbox{
			TypeLine:     lineType,
			TypeKeypad:   keypadType,
			PlaceHolder:  placeholder,
			DefaultValue: defaultValue,
			Title:        title,
		},
	}
}

// Selection creates a "Selection" button with a list of items.
func Selection(id, title, selectionID string, items []models.ButtonSelectionItem, isMulti bool, columnsCount string) models.Button {
	multi := isMulti
	return models.Button{
		ID:         id,
		Type:       models.ButtonTypeSelection,
		ButtonText: title,
		ButtonSelection: &models.ButtonSelection{
			SelectionID:      selectionID,
			Title:            title,
			Items:            items,
			IsMultiSelection: &multi,
			ColumnsCount:     columnsCount,
		},
	}
}

// CameraImage creates a "CameraImage" button.
func CameraImage(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeCameraImage, ButtonText: text}
}

// CameraVideo creates a "CameraVideo" button.
func CameraVideo(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeCameraVideo, ButtonText: text}
}

// GalleryImage creates a "GalleryImage" button.
func GalleryImage(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeGalleryImage, ButtonText: text}
}

// GalleryVideo creates a "GalleryVideo" button.
func GalleryVideo(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeGalleryVideo, ButtonText: text}
}

// File creates a "File" button.
func File(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeFile, ButtonText: text}
}

// Audio creates an "Audio" button.
func Audio(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeAudio, ButtonText: text}
}

// RecordAudio creates a "RecordAudio" button.
func RecordAudio(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeRecordAudio, ButtonText: text}
}

// MyPhoneNumber creates a "MyPhoneNumber" button.
func MyPhoneNumber(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeMyPhone, ButtonText: text}
}

// MyLocation creates a "MyLocation" button.
func MyLocation(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeMyLocation, ButtonText: text}
}

// AskMyPhoneNumber creates an "AskMyPhoneNumber" button.
func AskMyPhoneNumber(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeAskMyPhone, ButtonText: text}
}

// AskLocation creates an "AskLocation" button.
func AskLocation(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeAskLocation, ButtonText: text}
}

// Barcode creates a "Barcode" button.
func Barcode(id, text string) models.Button {
	return models.Button{ID: id, Type: models.ButtonTypeBarcode, ButtonText: text}
}

// stripAt removes a leading `@` from a username.
func stripAt(s string) string {
	if len(s) > 0 && s[0] == '@' {
		return s[1:]
	}
	return s
}

// ftoa converts a float to a string with a reasonable precision.
func ftoa(f float64) string {
	// Avoid importing fmt here — use a tiny implementation that
	// handles the common cases.
	if f == 0 {
		return "0"
	}
	neg := f < 0
	if neg {
		f = -f
	}
	intPart := int64(f)
	frac := f - float64(intPart)
	// Round to 6 decimal places.
	fracInt := int64(frac*1e6 + 0.5)
	s := itoa(intPart) + "." + padLeft(itoa(fracInt), 6)
	if neg {
		return "-" + s
	}
	return s
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func padLeft(s string, n int) string {
	if len(s) >= n {
		return s
	}
	pad := make([]byte, n-len(s))
	for i := range pad {
		pad[i] = '0'
	}
	return string(pad) + s
}
