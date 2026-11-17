// Package rules provides the composable rules used by the bot client
// to filter incoming updates before dispatching them to the user's
// callback.
package rules

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/bot/models"
)

// Rule is the contract implemented by every bot rule.
type Rule interface {
	// Check returns true if the rule matches the supplied update.
	Check(ctx context.Context, upd *models.Update) bool
}

// ===== Combinators =====================================================

// And matches when ALL sub-rules match (short-circuit).
type And []Rule

// Check implements Rule.
func (a And) Check(ctx context.Context, upd *models.Update) bool {
	for _, r := range a {
		if !r.Check(ctx, upd) {
			return false
		}
	}
	return true
}

// Or matches when ANY sub-rule matches.
type Or []Rule

// Check implements Rule.
func (o Or) Check(ctx context.Context, upd *models.Update) bool {
	for _, r := range o {
		if r.Check(ctx, upd) {
			return true
		}
	}
	return false
}

// Not negates the wrapped rule.
type Not struct{ Inner Rule }

// Check implements Rule.
func (n Not) Check(ctx context.Context, upd *models.Update) bool { return !n.Inner.Check(ctx, upd) }
