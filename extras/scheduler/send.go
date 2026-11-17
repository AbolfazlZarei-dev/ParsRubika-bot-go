// Package scheduler wraps the bot's scheduler to provide a simple
// "send later" helper. It schedules a goroutine that waits for the
// requested delay then sends the message. This is the client-side
// emulation of scheduled messages, since the Rubika API does not
// expose a native schedule endpoint.
package scheduler

import (
	"context"
	"time"
)

// Sender is the minimal interface a scheduler client needs (both
// `*bot.Bot` and `*account.Account` satisfy it).
type Sender interface {
	// SendMessage sends a text message and returns the message id
	// (or an error). The exact return type doesn't matter — the
	// scheduler ignores it.
	SendMessage(ctx context.Context, chatID, text string, _ ...any) (any, error)
}

// SendLater schedules `text` to be sent to `chatID` after `delay`.
// The supplied `sender` is used to send the message; pass the bot
// client or the account client.
//
// The returned `cancel` function can be invoked to cancel the
// scheduled send (e.g. when the user types `/cancel` before the
// reminder fires).
func SendLater(parent context.Context, sender Sender, chatID, text string, delay time.Duration) (cancel func()) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-time.After(delay):
			_, _ = sender.SendMessage(ctx, chatID, text)
		}
	}()
	return cancel
}

// SendRecurring schedules `text` to be sent to `chatID` every
// `interval`. The returned `cancel` function stops the recurring
// send.
func SendRecurring(parent context.Context, sender Sender, chatID, text string, interval time.Duration) (cancel func()) {
	ctx, cancel := context.WithCancel(parent)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, _ = sender.SendMessage(ctx, chatID, text)
			}
		}
	}()
	return cancel
}
