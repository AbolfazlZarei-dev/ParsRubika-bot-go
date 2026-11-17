package bot

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/rules"
)

// HandlerFunc is the ergonomic handler signature: everything a handler
// needs is reachable from the *Ctx.
type HandlerFunc func(c *Ctx) error

// Ctx bundles an update with the bot and offers short helpers for the
// most common reactions (reply, edit, delete, forward, arguments, ...).
// It embeds the request context, so it can be passed anywhere a
// context.Context is expected.
type Ctx struct {
	context.Context
	Bot    *Bot
	Update *Update

	mu     sync.Mutex
	values map[string]any
	chat   *models.Chat
}

// NewCtx wraps an update. Handlers registered through Handle, Command,
// Text, Regex, Button and the other shortcuts receive one automatically.
func NewCtx(ctx context.Context, b *Bot, u *Update) *Ctx {
	return &Ctx{Context: ctx, Bot: b, Update: u}
}

// Message returns the incoming message, or nil.
func (c *Ctx) Message() *models.Message { return c.Update.Message() }

// ChatID returns the chat the update belongs to.
func (c *Ctx) ChatID() string { return c.Update.ChatID }

// SenderID returns the ID of the user that wrote the message.
func (c *Ctx) SenderID() string { return c.Update.SenderID() }

// MessageID returns the ID of the incoming message.
func (c *Ctx) MessageID() string { return c.Update.MessageID() }

// Text returns the text of the incoming message.
func (c *Ctx) Text() string { return c.Update.Text() }

// ButtonID returns the ID of the pressed inline button, if any.
func (c *Ctx) ButtonID() string { return c.Update.ButtonID() }

// IsPrivate reports whether the chat is a one-to-one chat.
func (c *Ctx) IsPrivate() bool { return c.Update.IsPrivateChat() }

// Command returns the command name of the message without its prefix
// ("/start now" -> "start"), or "" if the text is not a command.
func (c *Ctx) Command() string {
	name, _ := splitCommand(c.Text())
	return name
}

// Args returns the whitespace separated arguments following a command.
func (c *Ctx) Args() []string {
	_, rest := splitCommand(c.Text())
	return strings.Fields(rest)
}

// ArgsText returns everything after the command word, trimmed.
func (c *Ctx) ArgsText() string {
	_, rest := splitCommand(c.Text())
	return strings.TrimSpace(rest)
}

func splitCommand(text string) (name, rest string) {
	text = strings.TrimSpace(text)
	if text == "" || !strings.ContainsAny(text[:1], "/!.#") {
		return "", ""
	}
	word := text
	if i := strings.IndexAny(text, " \t\n"); i >= 0 {
		word, rest = text[:i], text[i+1:]
	}
	word = strings.TrimLeft(word, "/!.#")
	if i := strings.IndexByte(word, '@'); i >= 0 {
		word = word[:i]
	}
	return strings.ToLower(word), rest
}

// Send sends a message to the current chat.
func (c *Ctx) Send(text string, opts ...SendOptions) (*models.MessageID, error) {
	return c.Bot.SendMessage(c, c.ChatID(), text, opts...)
}

// Reply sends a message to the current chat as a reply to the incoming one.
func (c *Ctx) Reply(text string, opts ...SendOptions) (*models.MessageID, error) {
	o := SendOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.ReplyToMessageID == "" {
		o.ReplyToMessageID = c.MessageID()
	}
	return c.Bot.SendMessage(c, c.ChatID(), text, o)
}

// ReplyKeypad replies with an inline keypad attached.
func (c *Ctx) ReplyKeypad(text string, kp *models.Keypad) (*models.MessageID, error) {
	return c.Reply(text, SendOptions{InlineKeypad: kp})
}

// Edit replaces the text of a message of this chat (defaults to the
// incoming message when messageID is empty).
func (c *Ctx) Edit(messageID, text string, opts ...EditOptions) error {
	if messageID == "" {
		messageID = c.MessageID()
	}
	return c.Bot.EditMessageText(c, c.ChatID(), messageID, text, opts...)
}

// Delete removes the incoming message.
func (c *Ctx) Delete() error { return c.Bot.DeleteMessage(c, c.ChatID(), c.MessageID()) }

// DeleteMessage removes another message of the current chat.
func (c *Ctx) DeleteMessage(messageID string) error {
	return c.Bot.DeleteMessage(c, c.ChatID(), messageID)
}

// Forward forwards the incoming message to another chat.
func (c *Ctx) Forward(toChatID string) (*models.MessageID, error) {
	return c.Bot.ForwardMessage(c, c.ChatID(), c.MessageID(), toChatID, false)
}

// CopyTo re-sends the content of the incoming message to another chat
// without a "forwarded from" header.
func (c *Ctx) CopyTo(toChatID string, opts ...SendOptions) (*models.MessageID, error) {
	m := c.Message()
	if m == nil {
		return nil, ErrNoMessage
	}
	return c.Bot.CopyMessage(c, m, toChatID, opts...)
}

// Chat returns (and caches) the chat object of the current chat.
func (c *Ctx) Chat() (*models.Chat, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.chat != nil {
		return c.chat, nil
	}
	ch, err := c.Bot.GetChat(c, c.ChatID())
	if err == nil {
		c.chat = ch
	}
	return ch, err
}

// Ban removes the sender of the message from the chat.
func (c *Ctx) Ban() error { return c.Bot.BanChatMember(c, c.ChatID(), c.SenderID()) }

// Set stores a value for the lifetime of this update (handy to pass
// data from middleware to handlers).
func (c *Ctx) Set(key string, v any) {
	c.mu.Lock()
	if c.values == nil {
		c.values = map[string]any{}
	}
	c.values[key] = v
	c.mu.Unlock()
}

// Get returns a value stored with Set.
func (c *Ctx) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.values[key]
	return v, ok
}

// GetString returns a string value stored with Set ("" if absent).
func (c *Ctx) GetString(key string) string {
	v, _ := c.Get(key)
	s, _ := v.(string)
	return s
}

// ===== Registration shortcuts ==========================================

// ctxKey carries the *Ctx through the middleware chain.
type ctxKey struct{}

// CtxFrom returns the *Ctx created for this update by Handle, or builds
// a fresh one. Middlewares can use it to Set values for the handler.
func CtxFrom(ctx context.Context, b *Bot, u *Update) *Ctx {
	if c, ok := ctx.Value(ctxKey{}).(*Ctx); ok && c.Update == u {
		return c
	}
	return NewCtx(ctx, b, u)
}

// Handle registers an ergonomic handler for updates matching ALL rules.
func (b *Bot) Handle(rs []Rule, fn HandlerFunc) int { return b.HandlePriority(0, rs, fn) }

// HandlePriority is Handle with an explicit priority (see OnPriority).
func (b *Bot) HandlePriority(prio int, rs []Rule, fn HandlerFunc) int {
	return b.OnPriority(prio, rs, func(ctx context.Context, bb *Bot, u *Update) error {
		return fn(CtxFrom(ctx, bb, u))
	})
}

// Command handles "/name", "!name" (and any alias) case-insensitively.
func (b *Bot) Command(name string, fn HandlerFunc, aliases ...string) int {
	names := append([]string{name}, aliases...)
	return b.Handle([]Rule{rules.NewCommand(names, []string{"/", "!"}, true, false)}, fn)
}

// Text handles messages whose text equals one of the given strings.
func (b *Bot) Text(fn HandlerFunc, exact ...string) int {
	return b.Handle([]Rule{rules.TextIn(exact...)}, fn)
}

// Regex handles messages whose text matches the regular expression.
func (b *Bot) Regex(pattern string, fn HandlerFunc) int {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("bot: invalid regular expression " + strconv.Quote(pattern) + ": " + err.Error())
	}
	return b.Handle([]Rule{rules.NewTextRegexCompiled(re)}, fn)
}

// Button handles presses of the inline button with the given ID.
func (b *Bot) Button(id string, fn HandlerFunc) int {
	return b.Handle([]Rule{rules.NewButton(id, false)}, fn)
}

// Message handles every message that carries text.
func (b *Bot) Message(fn HandlerFunc) int {
	return b.Handle([]Rule{rules.HasText{}}, fn)
}

// Fallback registers a catch-all handler. It always runs after every
// other handler, whatever the registration order.
func (b *Bot) Fallback(fn HandlerFunc) int { return b.HandlePriority(PriorityFallback, nil, fn) }

// CommandRule builds the rule used by Command: "/name" or "!name"
// (and aliases), case-insensitive.
func CommandRule(name string, aliases ...string) Rule {
	return rules.NewCommand(append([]string{name}, aliases...), []string{"/", "!"}, true, false)
}

// RegexRule builds a text regular-expression rule; it panics on an
// invalid pattern, like regexp.MustCompile.
func RegexRule(pattern string) Rule {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic("bot: invalid regular expression " + strconv.Quote(pattern) + ": " + err.Error())
	}
	return rules.NewTextRegexCompiled(re)
}
