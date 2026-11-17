package bot

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// ===== Built-in middleware factories ===========
//
// Reusable Middleware factories for common needs: timing, panic
// recovery, per-user rate limiting, access control, logging, etc.

// TimerMiddleware logs the time taken to process each update.
func TimerMiddleware(logFn func(format string, args ...any)) Middleware {
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		start := time.Now()
		err := next()
		elapsed := time.Since(start)
		if logFn != nil {
			chatID := ""
			if upd != nil {
				chatID = upd.ChatID
			}
			logFn("bot: update from %s took %s (err=%v)", chatID, elapsed, err)
		}
		return err
	}
}

// RecovererMiddleware catches panics in downstream handlers/middleware
// and logs the stack trace. Without this, a panic would crash the whole
// process.
func RecovererMiddleware(logFn func(format string, args ...any)) Middleware {
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) (err error) {
		defer func() {
			if r := recover(); r != nil {
				if logFn != nil {
					logFn("bot: PANIC recovered: %v\n%s", r, debug.Stack())
				}
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		return next()
	}
}

// LoggingMiddleware logs every update at info level.
func LoggingMiddleware(logFn func(format string, args ...any)) Middleware {
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		if logFn != nil && upd != nil {
			logFn("bot: update type=%s chat=%s", upd.Type, upd.ChatID)
		}
		return next()
	}
}

// AllowedUsersMiddleware only allows updates from the given user IDs.
// Other updates are silently dropped (next() is not called).
func AllowedUsersMiddleware(allowedUserIDs []string, rejectMessage string) Middleware {
	allowed := make(map[string]struct{}, len(allowedUserIDs))
	for _, id := range allowedUserIDs {
		allowed[id] = struct{}{}
	}
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		senderID := senderOf(upd)
		if _, ok := allowed[senderID]; !ok {
			if rejectMessage != "" && upd != nil && upd.ChatID != "" {
				_, _ = b.SendMessage(ctx, upd.ChatID, rejectMessage)
			}
			return nil
		}
		return next()
	}
}

// BlockedUsersMiddleware drops updates from the given user IDs.
func BlockedUsersMiddleware(blockedUserIDs []string) Middleware {
	blocked := make(map[string]struct{}, len(blockedUserIDs))
	for _, id := range blockedUserIDs {
		blocked[id] = struct{}{}
	}
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		senderID := senderOf(upd)
		if _, ok := blocked[senderID]; ok {
			return nil
		}
		return next()
	}
}

// PerUserRateLimitMiddleware limits how many updates per user per
// window are processed. Excess updates are silently dropped.
func PerUserRateLimitMiddleware(maxUpdates int, window time.Duration) Middleware {
	type counter struct {
		count       int
		windowStart time.Time
	}
	var counters = make(map[string]*counter)
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		senderID := senderOf(upd)
		if senderID == "" {
			return next()
		}
		c, ok := counters[senderID]
		if !ok || time.Since(c.windowStart) > window {
			c = &counter{count: 0, windowStart: time.Now()}
			counters[senderID] = c
		}
		c.count++
		if c.count > maxUpdates {
			return nil // drop
		}
		return next()
	}
}

// ChatTypeMiddleware only allows updates from the given chat types
// ("private", "group", "channel").
func ChatTypeMiddleware(chatTypes ...string) Middleware {
	allowed := make(map[string]struct{}, len(chatTypes))
	for _, t := range chatTypes {
		allowed[t] = struct{}{}
	}
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		if upd == nil {
			return next()
		}
		ct := chatTypeOf(upd)
		if _, ok := allowed[ct]; !ok {
			return nil
		}
		return next()
	}
}

// chatTypeOf returns "private", "group", or "channel" for an update.
// User (private) chats use the "u0" prefix; "b0" is a bot GUID, not a
// private-chat marker. Groups use "g0" and channels use "c0".
func chatTypeOf(upd *Update) string {
	if upd == nil || upd.ChatID == "" {
		return ""
	}
	switch {
	case len(upd.ChatID) > 2 && upd.ChatID[:2] == "u0":
		return "private"
	case len(upd.ChatID) > 2 && upd.ChatID[:2] == "g0":
		return "group"
	case len(upd.ChatID) > 2 && upd.ChatID[:2] == "c0":
		return "channel"
	}
	return ""
}

// CommandMiddleware extracts the command from the message text (if it
// starts with a prefix like / or !) and stores it in ExternalData for
// downstream rules to use.
func CommandMiddleware(prefixes ...string) Middleware {
	if len(prefixes) == 0 {
		prefixes = []string{"/", "!"}
	}
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		if upd == nil || upd.NewMessage == nil || upd.NewMessage.Text == "" {
			return next()
		}
		text := upd.NewMessage.Text
		for _, p := range prefixes {
			if len(text) > len(p) && text[:len(p)] == p {
				// Extract command (up to first space)
				rest := text[len(p):]
				cmd := rest
				if idx := indexByte(rest, ' '); idx >= 0 {
					cmd = rest[:idx]
				}
				if upd.ExternalData == nil {
					upd.ExternalData = make(map[string]any)
				}
				upd.ExternalData["command"] = cmd
				upd.ExternalData["command_prefix"] = p
				break
			}
		}
		return next()
	}
}

// indexByte returns the index of the first occurrence of b in s, or -1.
func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}

// UserStatsMiddleware tracks per-user message counts in a MetricsTracker.
func UserStatsMiddleware(tracker *MetricsTracker) Middleware {
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		if tracker != nil && upd != nil {
			senderID := senderOf(upd)
			if senderID != "" {
				tracker.IncrementCounter("user:"+senderID+":messages", 1)
				tracker.IncrementCounter("total:messages", 1)
			}
		}
		return next()
	}
}

// ChainMiddlewares combines multiple middlewares into a single middleware
// that runs them in order. This is useful when you want to apply a group
// of middlewares as a unit.
func ChainMiddlewares(mws ...Middleware) Middleware {
	return func(ctx context.Context, b *Bot, upd *Update, next func() error) error {
		// Build the chain from the inside out.
		chain := next
		for i := len(mws) - 1; i >= 0; i-- {
			chain = func(current Middleware, nxt func() error) func() error {
				return func() error {
					return current(ctx, b, upd, nxt)
				}
			}(mws[i], chain)
		}
		return chain()
	}
}

// EnsureUpdate is a helper that returns the *models.Update from a bot
// update, or nil if the update is nil. This is useful in middleware
// that needs to access update fields.
func EnsureUpdate(upd *Update) *models.Update {
	if upd == nil {
		return nil
	}
	return &models.Update{
		Type:           upd.Type,
		ChatID:         upd.ChatID,
		NewMessage:     upd.NewMessage,
		UpdatedMessage: upd.UpdatedMessage,
		UpdateTime:     upd.UpdateTime,
		ExternalData:   upd.ExternalData,
	}
}
