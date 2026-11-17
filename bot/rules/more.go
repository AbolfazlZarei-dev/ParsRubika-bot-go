package rules

import (
	"context"
	"strings"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// Func adapts a plain function into a Rule.
type Func func(ctx context.Context, upd *models.Update) bool

// Check implements Rule.
func (f Func) Check(ctx context.Context, upd *models.Update) bool { return f != nil && f(ctx, upd) }

// Always matches every update.
type Always struct{}

// Check implements Rule.
func (Always) Check(context.Context, *models.Update) bool { return true }

// Never matches no update.
type Never struct{}

// Check implements Rule.
func (Never) Check(context.Context, *models.Update) bool { return false }

// HasText matches updates whose message carries non-empty text.
type HasText struct{}

// Check implements Rule.
func (HasText) Check(_ context.Context, u *models.Update) bool { return u.Text() != "" }

// TextInRule matches when the text equals any of the given strings.
type TextInRule struct{ set map[string]struct{} }

// TextIn builds a TextInRule.
func TextIn(values ...string) *TextInRule {
	r := &TextInRule{set: make(map[string]struct{}, len(values))}
	for _, v := range values {
		r.set[v] = struct{}{}
	}
	return r
}

// Check implements Rule.
func (r *TextInRule) Check(_ context.Context, u *models.Update) bool {
	_, ok := r.set[u.Text()]
	return ok
}

// TextContainsAny matches when the text contains any of the substrings.
type TextContainsAny struct {
	Substrings []string
	FoldCase   bool
}

// Check implements Rule.
func (r TextContainsAny) Check(_ context.Context, u *models.Update) bool {
	t := u.Text()
	if r.FoldCase {
		t = strings.ToLower(t)
	}
	for _, s := range r.Substrings {
		if r.FoldCase {
			s = strings.ToLower(s)
		}
		if s != "" && strings.Contains(t, s) {
			return true
		}
	}
	return false
}

// SenderIn matches updates produced by one of the given user IDs.
type SenderIn struct{ set map[string]struct{} }

// NewSenderIn builds a SenderIn rule.
func NewSenderIn(ids ...string) *SenderIn {
	r := &SenderIn{set: make(map[string]struct{}, len(ids))}
	for _, id := range ids {
		r.set[id] = struct{}{}
	}
	return r
}

// Check implements Rule.
func (r *SenderIn) Check(_ context.Context, u *models.Update) bool {
	_, ok := r.set[u.SenderID()]
	return ok
}

// ChatIn matches updates of one of the given chats.
type ChatIn struct{ set map[string]struct{} }

// NewChatIn builds a ChatIn rule.
func NewChatIn(ids ...string) *ChatIn {
	r := &ChatIn{set: make(map[string]struct{}, len(ids))}
	for _, id := range ids {
		r.set[id] = struct{}{}
	}
	return r
}

// Check implements Rule.
func (r *ChatIn) Check(_ context.Context, u *models.Update) bool {
	_, ok := r.set[u.ChatID]
	return ok
}

// MessageIDIs matches one exact message ID.
type MessageIDIs string

// Check implements Rule.
func (r MessageIDIs) Check(_ context.Context, u *models.Update) bool {
	return u.MessageID() == string(r)
}

// PrivateChat matches one-to-one chats only.
type PrivateChat struct{}

// Check implements Rule.
func (PrivateChat) Check(_ context.Context, u *models.Update) bool { return u.IsPrivateChat() }
