package rules

import (
	"context"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// Private matches when the update came from a private (user-to-bot)
// chat. User object GUIDs start with "u0".
type Private struct{}

// Check implements Rule.
func (Private) Check(_ context.Context, upd *models.Update) bool {
	return strings.HasPrefix(upd.ChatID, "u0")
}

// Group matches when the update came from a group.
type Group struct{}

// Check implements Rule.
func (Group) Check(_ context.Context, upd *models.Update) bool {
	return strings.HasPrefix(upd.ChatID, "g0")
}

// Channel matches when the update came from a channel.
type Channel struct{}

// Check implements Rule.
func (Channel) Check(_ context.Context, upd *models.Update) bool {
	return strings.HasPrefix(upd.ChatID, "c0")
}

// BotSender matches messages sent by a bot.
type BotSender struct{}

// Check implements Rule.
func (BotSender) Check(_ context.Context, upd *models.Update) bool {
	senderType := ""
	if upd.NewMessage != nil {
		senderType = upd.NewMessage.SenderType
	} else if upd.UpdatedMessage != nil {
		senderType = upd.UpdatedMessage.SenderType
	}
	return senderType == "Bot"
}

// ChatID matches by chat ID.
type ChatID struct{ IDs []string }

// NewChatID constructs a ChatID rule.
func NewChatID(ids ...string) *ChatID { return &ChatID{IDs: ids} }

// Check implements Rule.
func (c *ChatID) Check(_ context.Context, upd *models.Update) bool {
	for _, id := range c.IDs {
		if upd.ChatID == id {
			return true
		}
	}
	return false
}

// UpdateType matches when the update type is one of `types`.
type UpdateType struct{ Types []models.UpdateType }

// NewUpdateType constructs an UpdateType rule.
func NewUpdateType(types ...models.UpdateType) *UpdateType { return &UpdateType{Types: types} }

// Check implements Rule.
func (u *UpdateType) Check(_ context.Context, upd *models.Update) bool {
	for _, t := range u.Types {
		if upd.Type == t {
			return true
		}
	}
	return false
}

// Button matches a button click by `button_id`. Supports regex.
type Button struct {
	ButtonID string
	Regex    bool
	re       *regexp.Regexp
}

// NewButton constructs a Button rule.
func NewButton(buttonID string, regex bool) *Button {
	f := &Button{ButtonID: buttonID, Regex: regex}
	if regex {
		f.re, _ = regexp.Compile(buttonID)
	}
	return f
}

// Check implements Rule.
func (b *Button) Check(_ context.Context, upd *models.Update) bool {
	var aux *models.AuxData
	if upd.NewMessage != nil {
		aux = upd.NewMessage.AuxData
	} else if upd.UpdatedMessage != nil {
		aux = upd.UpdatedMessage.AuxData
	}
	if aux == nil {
		return false
	}
	bid := aux.ButtonID
	if bid == "" {
		return false
	}
	if b.Regex {
		return b.re != nil && b.re.MatchString(bid)
	}
	return bid == b.ButtonID
}

// Forward matches forwarded messages.
type Forward struct{}

// Check implements Rule.
func (Forward) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.ForwardedFrom != nil || upd.NewMessage.ForwardedNoLink != ""
	}
	if upd.UpdatedMessage != nil {
		return upd.UpdatedMessage.ForwardedFrom != nil || upd.UpdatedMessage.ForwardedNoLink != ""
	}
	return false
}

// File matches messages with a file attachment.
type File struct{}

// Check implements Rule.
func (File) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.File != nil
	}
	if upd.UpdatedMessage != nil {
		return upd.UpdatedMessage.File != nil
	}
	return false
}

// Sticker matches sticker messages.
type Sticker struct{}

// Check implements Rule.
func (Sticker) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.Sticker != nil
	}
	return false
}

// Location matches messages with a location attachment.
type Location struct{}

// Check implements Rule.
func (Location) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.Location != nil
	}
	if upd.UpdatedMessage != nil {
		return upd.UpdatedMessage.Location != nil
	}
	return false
}

// Poll matches messages with a poll attachment.
type Poll struct{}

// Check implements Rule.
func (Poll) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.Poll != nil
	}
	return false
}

// Contact matches messages with a contact attachment.
type Contact struct{}

// Check implements Rule.
func (Contact) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.ContactMessage != nil
	}
	return false
}

// Reply matches messages that reply to a given message ID.
type Reply struct{ ToMessageID string }

// NewReply constructs a Reply rule.
func NewReply(toMessageID string) *Reply { return &Reply{ToMessageID: toMessageID} }

// Check implements Rule.
func (r *Reply) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage != nil {
		return upd.NewMessage.ReplyToMessageID == r.ToMessageID
	}
	if upd.UpdatedMessage != nil {
		return upd.UpdatedMessage.ReplyToMessageID == r.ToMessageID
	}
	return false
}
