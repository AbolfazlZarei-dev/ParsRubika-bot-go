package rules

import (
	"context"
	"regexp"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// ===== Advanced filter rules ===================
//
// Granular filters for polls, locations, forwards, edits and metadata.

// PollQuestionContains matches polls whose question contains `keyword`
// (case-insensitive).
type PollQuestionContains struct{ Keyword string }

// Check implements Rule.
func (p PollQuestionContains) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Poll == nil {
		return false
	}
	return strings.Contains(
		strings.ToLower(upd.NewMessage.Poll.Question),
		strings.ToLower(p.Keyword),
	)
}

// PollOptionCount matches polls with between min and max options
// (inclusive).
type PollOptionCount struct{ Min, Max int }

// Check implements Rule.
func (p PollOptionCount) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Poll == nil {
		return false
	}
	n := len(upd.NewMessage.Poll.Options)
	if p.Min > 0 && n < p.Min {
		return false
	}
	if p.Max > 0 && n > p.Max {
		return false
	}
	return true
}

// LocationWithin matches location messages whose latitude/longitude fall
// within the given bounding box.
type LocationWithin struct {
	LatMin, LatMax, LongMin, LongMax float64
}

// Check implements Rule.
func (l LocationWithin) Check(_ context.Context, upd *models.Update) bool {
	var lat, long float64
	if upd.NewMessage != nil && upd.NewMessage.Location != nil {
		lat = parseFloat(upd.NewMessage.Location.Latitude)
		long = parseFloat(upd.NewMessage.Location.Longitude)
	} else if upd.NewMessage != nil && upd.NewMessage.LiveLocation != nil &&
		upd.NewMessage.LiveLocation.CurrentLocation != nil {
		lat = parseFloat(upd.NewMessage.LiveLocation.CurrentLocation.Latitude)
		long = parseFloat(upd.NewMessage.LiveLocation.CurrentLocation.Longitude)
	} else {
		return false
	}
	return lat >= l.LatMin && lat <= l.LatMax && long >= l.LongMin && long <= l.LongMax
}

// MetaLinksContain matches messages whose metadata contains a link URL
// containing `substring`.
type MetaLinksContain struct{ Substring string }

// Check implements Rule.
func (m MetaLinksContain) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Metadata == nil {
		return false
	}
	for _, part := range upd.NewMessage.Metadata.MetaDataParts {
		if strings.Contains(part.LinkURL, m.Substring) {
			return true
		}
	}
	return false
}

// MetaTypesInclude matches messages whose metadata includes at least one
// part of the given types (e.g. "Bold", "Italic", "Link").
type MetaTypesInclude struct{ Types []string }

// Check implements Rule.
func (m MetaTypesInclude) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.Metadata == nil {
		return false
	}
	for _, part := range upd.NewMessage.Metadata.MetaDataParts {
		for _, t := range m.Types {
			if part.Type == t {
				return true
			}
		}
	}
	return false
}

// IsReplyToUser matches messages that are a reply to a message from the
// given user ID.
type IsReplyToUser struct{ UserID string }

// Check implements Rule.
func (r IsReplyToUser) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.ReplyToMessageID == "" {
		return false
	}
	// We cannot check the reply author without fetching the replied
	// message; this rule is best used in conjunction with a middleware
	// that pre-populates ExternalData["reply_sender_id"].
	if upd.ExternalData != nil {
		if sid, ok := upd.ExternalData["reply_sender_id"].(string); ok {
			return sid == r.UserID
		}
	}
	return false
}

// IsForwardedFrom matches messages forwarded from the given user/channel.
type IsForwardedFrom struct{ FromID string }

// Check implements Rule.
func (f IsForwardedFrom) Check(_ context.Context, upd *models.Update) bool {
	if upd.NewMessage == nil || upd.NewMessage.ForwardedFrom == nil {
		return false
	}
	return upd.NewMessage.ForwardedFrom.FromChatID == f.FromID ||
		upd.NewMessage.ForwardedFrom.FromSenderID == f.FromID
}

// EditedTextContains matches edited messages whose new text contains
// `keyword` (case-insensitive).
type EditedTextContains struct{ Keyword string }

// Check implements Rule.
func (e EditedTextContains) Check(_ context.Context, upd *models.Update) bool {
	if upd.UpdatedMessage == nil {
		return false
	}
	return strings.Contains(
		strings.ToLower(upd.UpdatedMessage.Text),
		strings.ToLower(e.Keyword),
	)
}

// AuxDataContains matches updates whose aux_data contains the given
// key/value pair.
type AuxDataContains struct {
	Key   string
	Value string
}

// Check implements Rule.
func (a AuxDataContains) Check(_ context.Context, upd *models.Update) bool {
	var aux *models.AuxData
	if upd.NewMessage != nil {
		aux = upd.NewMessage.AuxData
	} else if upd.UpdatedMessage != nil {
		aux = upd.UpdatedMessage.AuxData
	}
	if aux == nil {
		return false
	}
	switch a.Key {
	case "start_id":
		return aux.StartID == a.Value
	case "button_id":
		return aux.ButtonID == a.Value
	}
	return false
}

// TextRegexCompiled matches text using a pre-compiled regex. This is
// more efficient than TextRegex when the same pattern is reused.
type TextRegexCompiled struct {
	re *regexp.Regexp
}

// NewTextRegexCompiled creates a TextRegexCompiled rule from a compiled
// regex.
func NewTextRegexCompiled(re *regexp.Regexp) *TextRegexCompiled {
	return &TextRegexCompiled{re: re}
}

// Check implements Rule.
func (t *TextRegexCompiled) Check(_ context.Context, upd *models.Update) bool {
	if t.re == nil {
		return false
	}
	text := upd.Text()
	if text == "" {
		return false
	}
	return t.re.MatchString(text)
}

// ChatTitleRegex matches when the chat title (from ExternalData) matches
// the given regex.
type ChatTitleRegex struct{ Pattern string }

// Check implements Rule.
func (c ChatTitleRegex) Check(_ context.Context, upd *models.Update) bool {
	if upd.ExternalData == nil {
		return false
	}
	title, ok := upd.ExternalData["chat_title"].(string)
	if !ok {
		return false
	}
	re, err := regexp.Compile(c.Pattern)
	if err != nil {
		return false
	}
	return re.MatchString(title)
}

// ChatHasUsername matches when the chat has a username (from ExternalData
// or from a pre-fetched Chat object).
type ChatHasUsername struct{}

// Check implements Rule.
func (ChatHasUsername) Check(_ context.Context, upd *models.Update) bool {
	if upd.ExternalData == nil {
		return false
	}
	username, _ := upd.ExternalData["chat_username"].(string)
	return username != ""
}

// parseFloat is a small helper that parses a string to float64. Returns
// 0 on error. Used by LocationWithin.
func parseFloat(s string) float64 {
	var f float64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			f = f*10 + float64(c-'0')
		} else if c == '.' {
			break
		}
	}
	// Handle decimal part (simplified)
	dotIdx := strings.Index(s, ".")
	if dotIdx >= 0 {
		var frac float64
		div := 10.0
		for _, c := range s[dotIdx+1:] {
			if c >= '0' && c <= '9' {
				frac += float64(c-'0') / div
				div *= 10
			} else {
				break
			}
		}
		f += frac
	}
	return f
}
