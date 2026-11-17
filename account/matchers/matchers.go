// Package matchers provides the composable matchers used by the
// account client to filter incoming updates before dispatching them
// to the user's callback.
//
// Each matcher implements the account.Matcher interface. Matchers
// can be combined with the three combinators And, Or and Not.
package matchers

import (
	"context"

	"github.com/AbolfazlZarei-dev/ParsRubika-bot-go/v4/events"
)

// Matcher is the contract implemented by every account-side matcher.
// It is identical to the account.Matcher interface; it is re-declared
// here to avoid an import cycle.
type Matcher interface {
	Kind() string
	Match(ctx context.Context, raw map[string]any) bool
}

// asEvent wraps `raw` into an *events.Event for the property-style
// helpers.
func asEvent(raw map[string]any) *events.Event { return events.New(raw) }
