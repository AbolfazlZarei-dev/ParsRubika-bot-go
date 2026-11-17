package models

import "strings"

// ButtonSelectionItem is one item shown inside a selection button.
type ButtonSelectionItem struct {
	Text     string `json:"text"`
	ImageURL string `json:"image_url,omitempty"`
	Type     string `json:"type,omitempty"`
}

// ButtonSelection describes the `button_selection` field of a button.
type ButtonSelection struct {
	SelectionID      string                `json:"selection_id,omitempty"`
	SearchType       string                `json:"search_type,omitempty"`
	GetType          string                `json:"get_type,omitempty"`
	Items            []ButtonSelectionItem `json:"items,omitempty"`
	IsMultiSelection *bool                 `json:"is_multi_selection,omitempty"`
	ColumnsCount     string                `json:"columns_count,omitempty"`
	Title            string                `json:"title,omitempty"`
}

// ButtonCalendar is the calendar button.
type ButtonCalendar struct {
	DefaultValue string `json:"default_value,omitempty"`
	Type         string `json:"type,omitempty"`
	MinYear      string `json:"min_year,omitempty"`
	MaxYear      string `json:"max_year,omitempty"`
	Title        string `json:"title,omitempty"`
}

// ButtonNumberPicker is the numeric picker button.
type ButtonNumberPicker struct {
	MinValue     string `json:"min_value,omitempty"`
	MaxValue     string `json:"max_value,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Title        string `json:"title,omitempty"`
}

// ButtonStringPicker is the string-list picker button.
type ButtonStringPicker struct {
	Items        []string `json:"items,omitempty"`
	DefaultValue string   `json:"default_value,omitempty"`
	Title        string   `json:"title,omitempty"`
}

// ButtonTextbox is the textbox button.
type ButtonTextbox struct {
	TypeLine     string `json:"type_line,omitempty"`
	TypeKeypad   string `json:"type_keypad,omitempty"`
	PlaceHolder  string `json:"place_holder,omitempty"`
	Title        string `json:"title,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
}

// ButtonLocation is the location-picker button.
type ButtonLocation struct {
	DefaultPointerLocation *Location `json:"default_pointer_location,omitempty"`
	DefaultMapLocation     *Location `json:"default_map_location,omitempty"`
	Type                   string    `json:"type,omitempty"`
	Title                  string    `json:"title,omitempty"`
	LocationImageURL       string    `json:"location_image_url,omitempty"`
}

// ButtonLink wraps a hyperlink or join-channel / open-chat button.
type ButtonLink struct {
	Type            string           `json:"type,omitempty"`
	LinkURL         string           `json:"link_url,omitempty"`
	JoinChannelData *JoinChannelData `json:"joinchannel_data,omitempty"`
	OpenChatData    *OpenChatData    `json:"open_chat_data,omitempty"`
}

// NormalizeLinkURL converts the rubika.ir short links into the
// deep-link form expected by the server.
func (b *ButtonLink) NormalizeLinkURL() {
	if b.LinkURL == "" {
		return
	}
	prefixes := map[string]string{
		"https://rubika.ir/joing/": "rubika://g.rubika.ir/",
		"https://rubika.ir/joinc/": "rubika://c.rubika.ir/",
		"https://rubika.ir/post/":  "rubika://p.rubika.ir/",
	}
	for prefix, deep := range prefixes {
		if strings.HasPrefix(b.LinkURL, prefix) {
			b.LinkURL = deep + b.LinkURL[len(prefix):]
			break
		}
	}
}

// Button is a single key inside a keypad row.
type Button struct {
	ID                 string              `json:"id,omitempty"`
	Type               ButtonType          `json:"type,omitempty"`
	ButtonText         string              `json:"button_text,omitempty"`
	ButtonSelection    *ButtonSelection    `json:"button_selection,omitempty"`
	ButtonCalendar     *ButtonCalendar     `json:"button_calendar,omitempty"`
	ButtonNumberPicker *ButtonNumberPicker `json:"button_number_picker,omitempty"`
	ButtonStringPicker *ButtonStringPicker `json:"button_string_picker,omitempty"`
	ButtonLocation     *ButtonLocation     `json:"button_location,omitempty"`
	ButtonTextbox      *ButtonTextbox      `json:"button_textbox,omitempty"`
	ButtonLink         *ButtonLink         `json:"button_link,omitempty"`
}

// KeypadRow is a row of buttons.
type KeypadRow struct {
	Buttons []Button `json:"buttons,omitempty"`
}

// Keypad is the full keypad object attached to a message or chat.
type Keypad struct {
	Rows           []KeypadRow `json:"rows,omitempty"`
	ResizeKeyboard *bool       `json:"resize_keyboard,omitempty"`
	OnTimeKeyboard *bool       `json:"on_time_keyboard,omitempty"`
}

// MessageKeypadUpdate is the payload for `editMessageKeypad`.
type MessageKeypadUpdate struct {
	MessageID    string  `json:"message_id"`
	InlineKeypad *Keypad `json:"inline_keypad,omitempty"`
}
