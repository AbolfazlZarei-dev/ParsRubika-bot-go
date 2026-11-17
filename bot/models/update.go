package models

// Update is the high-level update dispatched by the bot's polling /
// webhook loop.
type Update struct {
	Type             UpdateType     `json:"type"`
	ChatID           string         `json:"chat_id"`
	RemovedMessageID string         `json:"removed_message_id,omitempty"`
	NewMessage       *Message       `json:"new_message,omitempty"`
	UpdatedMessage   *Message       `json:"updated_message,omitempty"`
	UpdatedPayment   *PaymentStatus `json:"updated_payment,omitempty"`
	UpdateTime       string         `json:"update_time,omitempty"`
	ExternalData     map[string]any `json:"-"`
}

// MessageID convenience accessor.
func (u *Update) MessageID() string {
	if u.NewMessage != nil {
		return u.NewMessage.MessageID
	}
	if u.UpdatedMessage != nil {
		return u.UpdatedMessage.MessageID
	}
	return u.RemovedMessageID
}

// Text returns the message text (if any).
func (u *Update) Text() string {
	if u.NewMessage != nil {
		return u.NewMessage.Text
	}
	if u.UpdatedMessage != nil {
		return u.UpdatedMessage.Text
	}
	return ""
}

// InlineMessage is the InlineMessage update type. It is dispatched by
// the server when the bot is invoked in inline-mode.
type InlineMessage struct {
	SenderID     string         `json:"sender_id"`
	Text         string         `json:"text,omitempty"`
	File         *File          `json:"file,omitempty"`
	Location     *Location      `json:"location,omitempty"`
	AuxData      *AuxData       `json:"aux_data,omitempty"`
	MessageID    string         `json:"message_id"`
	ChatID       string         `json:"chat_id"`
	ExternalData map[string]any `json:"-"`
}

// FindKey is the convenience helper used by the rules to inspect nested
// fields by name (for example "text", "message_id" or "chat_id").
func (u *Update) FindKey(key string) any {
	switch key {
	case "text":
		return u.Text()
	case "message_id":
		return u.MessageID()
	case "chat_id":
		return u.ChatID
	case "type":
		return string(u.Type)
	case "sender_id":
		if u.NewMessage != nil {
			return u.NewMessage.SenderID
		}
		if u.UpdatedMessage != nil {
			return u.UpdatedMessage.SenderID
		}
	case "sender_type":
		if u.NewMessage != nil {
			return u.NewMessage.SenderType
		}
		if u.UpdatedMessage != nil {
			return u.UpdatedMessage.SenderType
		}
	case "aux_data":
		if u.NewMessage != nil {
			return u.NewMessage.AuxData
		}
		if u.UpdatedMessage != nil {
			return u.UpdatedMessage.AuxData
		}
	case "file":
		if u.NewMessage != nil {
			return u.NewMessage.File
		}
		if u.UpdatedMessage != nil {
			return u.UpdatedMessage.File
		}
	case "sticker":
		if u.NewMessage != nil {
			return u.NewMessage.Sticker
		}
	case "forwarded_no_link":
		if u.NewMessage != nil {
			return u.NewMessage.ForwardedNoLink
		}
	}
	return nil
}

// FindKey is the inline-message counterpart of Update.FindKey.
func (i *InlineMessage) FindKey(key string) any {
	switch key {
	case "text":
		return i.Text
	case "message_id":
		return i.MessageID
	case "chat_id":
		return i.ChatID
	case "sender_id":
		return i.SenderID
	case "file":
		return i.File
	case "location":
		return i.Location
	case "aux_data":
		return i.AuxData
	}
	return nil
}

// Message returns the message carried by the update (new or edited),
// or nil for updates that carry none.
func (u *Update) Message() *Message {
	if u.NewMessage != nil {
		return u.NewMessage
	}
	return u.UpdatedMessage
}

// SenderID returns the ID of the user that produced the update.
func (u *Update) SenderID() string {
	if m := u.Message(); m != nil {
		return m.SenderID
	}
	return ""
}

// ButtonID returns the ID of the inline button that triggered the update.
func (u *Update) ButtonID() string {
	if m := u.Message(); m != nil && m.AuxData != nil {
		return m.AuxData.ButtonID
	}
	return ""
}

// IsPrivateChat reports whether the update comes from a one-to-one
// (user-to-bot) chat. In Rubika, user object GUIDs start with "u0";
// "b0" is the bot GUID prefix, not a private-chat marker. Groups use
// "g0" and channels use "c0".
func (u *Update) IsPrivateChat() bool { return len(u.ChatID) > 1 && u.ChatID[:2] == "u0" }
