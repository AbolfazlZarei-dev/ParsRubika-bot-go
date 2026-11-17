// Package forms is a ready-made plugin that provides a tiny
// multi-step form helper. Each form is a list of questions asked
// one after another, with per-field validation.
package forms

import (
	"context"
	"sync"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot"
	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// Step is a single step of a form. The callback receives the user's
// answer and returns the next step name (or "" to finish).
type Step func(ctx context.Context, b *bot.Bot, upd *models.Update, answer string) (next string, err error)

// Form is a multi-step conversation scoped to a single chat. Steps
// are registered by name; the first step is named "start".
type Form struct {
	mu      sync.Mutex
	steps   map[string]Step
	current map[string]string // chatID → current step
}

// New constructs a fresh Form.
func New() *Form {
	return &Form{
		steps:   map[string]Step{},
		current: map[string]string{},
	}
}

// Add registers a step under the supplied name.
func (f *Form) Add(name string, step Step) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.steps[name] = step
}

// Start begins the form for a chat at the "start" step.
func (f *Form) Start(chatID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.current[chatID] = "start"
}

// IsActive reports whether a form is in progress for a chat.
func (f *Form) IsActive(chatID string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.current[chatID]
	return ok
}

// Stop ends the form for a chat.
func (f *Form) Stop(chatID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.current, chatID)
}

// Handle is the callback you register on the bot for any message
// that arrives while a form is active.
func (f *Form) Handle(ctx context.Context, b *bot.Bot, upd *models.Update) error {
	chatID := upd.ChatID
	f.mu.Lock()
	stepName, ok := f.current[chatID]
	// Look up the step while still holding the lock to avoid a data
	// race with concurrent Add() calls that write f.steps.
	var step Step
	if ok {
		step, ok = f.steps[stepName]
	}
	f.mu.Unlock()
	if !ok {
		return nil
	}
	answer := ""
	if upd.NewMessage != nil {
		answer = upd.NewMessage.Text
	}
	next, err := step(ctx, b, upd, answer)
	if err != nil {
		return err
	}
	if next == "" {
		f.Stop(chatID)
		return nil
	}
	f.mu.Lock()
	f.current[chatID] = next
	f.mu.Unlock()
	return nil
}
