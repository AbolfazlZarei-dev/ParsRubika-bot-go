package bot

import (
	"context"
	"sync"
	"time"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// meCacheTTL is how long the result of getMe is reused.
const meCacheTTL = 10 * time.Minute

type meCache struct {
	mu   sync.Mutex
	bot  *models.Bot
	when time.Time
}

// Me returns the bot's own profile. The answer is cached for ten
// minutes, so it is cheap to call from handlers.
func (b *Bot) Me(ctx context.Context) (*models.Bot, error) {
	b.me.mu.Lock()
	defer b.me.mu.Unlock()
	if b.me.bot != nil && time.Since(b.me.when) < meCacheTTL {
		cp := *b.me.bot
		return &cp, nil
	}
	m, err := b.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	b.me.bot, b.me.when = m, time.Now()
	cp := *m
	return &cp, nil
}

// Name returns the display name of the bot.
func (b *Bot) Name(ctx context.Context) (string, error) {
	m, err := b.Me(ctx)
	if err != nil {
		return "", err
	}
	return m.BotTitle, nil
}

// Username returns the bot's @username (without the "@").
func (b *Bot) Username(ctx context.Context) (string, error) {
	m, err := b.Me(ctx)
	if err != nil {
		return "", err
	}
	return m.Username, nil
}
