package bot

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// submit queues an update on the per-chat executor so that updates of
// one chat are handled in order while different chats run in parallel.
func (b *Bot) submit(ctx context.Context, u *models.Update) {
	if u == nil {
		return
	}
	b.stats.Received.Add(1)
	key := u.ChatID
	b.disp.submit(key, func() {
		if err := b.processUpdate(ctx, u); err != nil {
			b.errorf("update %s: %v", u.Type, err)
		}
	})
}

// Dispatch feeds an update to the bot synchronously, exactly as the
// polling loop would. It is handy for tests and custom transports.
func (b *Bot) Dispatch(ctx context.Context, u *models.Update) error {
	b.stats.Received.Add(1)
	return b.processUpdate(ctx, u)
}

// processUpdate de-duplicates the update, runs the middleware chain and
// finally the first matching handler. Panics never escape.
func (b *Bot) processUpdate(ctx context.Context, u *models.Update) (err error) {
	defer func() {
		if r := recover(); r != nil {
			b.stats.Panics.Add(1)
			ec := &ErrorContext{Ctx: ctx, Bot: b, Update: u, Panic: r, Stack: debug.Stack()}
			b.reportError(ec)
			err = nil
		}
	}()

	if key := updateKey(u); key != "" && !b.seen.add(key) {
		b.stats.Duplicate.Add(1)
		return nil
	}

	// One Ctx per update, shared by every middleware and the handler.
	c := NewCtx(ctx, b, u)
	ctx = context.WithValue(ctx, ctxKey{}, c)
	c.Context = ctx

	b.mu.RLock()
	mws := make([]Middleware, len(b.middlewares))
	for i, m := range b.middlewares {
		mws[i] = m.fn
	}
	handlers := append([]handlerEntry(nil), b.handlers...)
	b.mu.RUnlock()

	// The handler dispatch is the terminal element of the chain, so a
	// middleware that does not call next() really does stop the update
	// and one that wraps next() observes the handler's result and panics.
	chain := func() error { return b.runHandlers(ctx, u, handlers) }
	for i := len(mws) - 1; i >= 0; i-- {
		mw, next := mws[i], chain
		chain = func() error { return mw(ctx, b, u, next) }
	}
	if err := chain(); err != nil {
		b.stats.Failed.Add(1)
		b.reportError(&ErrorContext{Ctx: ctx, Bot: b, Update: u, Err: err})
	}
	return nil
}

// runHandlers calls the first handler whose rules all match. A handler
// returning ErrContinue hands the update to the next matching handler.
func (b *Bot) runHandlers(ctx context.Context, u *models.Update, handlers []handlerEntry) error {
	for _, h := range handlers {
		if !matchAll(ctx, h.rules, u) {
			continue
		}
		err := h.callback(ctx, b, u)
		if errors.Is(err, ErrContinue) {
			continue
		}
		b.stats.Handled.Add(1)
		return err
	}
	return nil
}

func matchAll(ctx context.Context, rs []Rule, u *models.Update) bool {
	for _, r := range rs {
		if !r.Check(ctx, u) {
			return false
		}
	}
	return true
}

// reportError forwards a failure to the configured error handler, or
// logs it when none is installed.
func (b *Bot) reportError(ec *ErrorContext) {
	if h := b.opts.ErrorHandler; h != nil {
		func() {
			defer func() { _ = recover() }()
			h(ec)
		}()
		return
	}
	if ec.Panic != nil {
		b.errorf("panic in handler: %v\n%s", ec.Panic, ec.Stack)
		return
	}
	b.errorf("handler error: %v", ec.Err)
}

// updateKey identifies an update for duplicate suppression. Updates
// without a stable identity (start/stop events, payments) are never
// filtered.
func updateKey(u *models.Update) string {
	if u == nil {
		return ""
	}
	switch u.Type {
	case models.UpdateTypeNewMessage, models.UpdateTypeUpdatedMessage:
		if id := u.MessageID(); id != "" {
			return "u:" + string(u.Type) + ":" + id + ":" + u.UpdateTime
		}
	case models.UpdateTypeRemovedMessage:
		if u.RemovedMessageID != "" {
			return "u:" + string(u.Type) + ":" + u.RemovedMessageID + ":" + u.ChatID
		}
	}
	return ""
}
